// Copyright (C) 2015 NTT Innovation Institute, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcdv3

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	syn "sync"
	"time"

	"github.com/cloudwan/gohan/metrics"
	"github.com/cloudwan/gohan/sync"
	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	pb "github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/streamrail/concurrent-map"
	"github.com/twinj/uuid"
)

const (
	processPath = "/gohan/cluster/process"
	masterTTL   = 10
)

//Sync is struct for etcd based sync
type Sync struct {
	locks      cmap.ConcurrentMap
	timeout    time.Duration
	etcdClient *etcd.Client
	processID  string
}

func (s *Sync) withTimeout() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), s.timeout)
	return ctx
}

//NewSync initialize new etcd sync
func NewSync(etcdServers []string, timeout time.Duration) (*Sync, error) {
	sync := &Sync{
		locks:   cmap.New(),
		timeout: timeout,
	}
	client, err := etcd.New(
		etcd.Config{
			Endpoints:   etcdServers,
			DialTimeout: timeout,
		},
	)
	if err != nil {
		return nil, err
	}
	sync.etcdClient = client
	hostname, _ := os.Hostname()
	sync.processID = hostname + uuid.NewV4().String()
	return sync, nil
}

//GetProcessID returns processID
func (s *Sync) GetProcessID() string {
	return s.processID
}

func measureTime(timeStarted time.Time, action string) {
	metrics.UpdateTimer(timeStarted, "sync.v3.%s", action)
}

func updateCounter(delta int64, counter string) {
	metrics.UpdateCounter(delta, "sync.v3.%s", counter)
}

//Update sync update sync
//When jsonString is empty, this method do nothing because
//etcd v3 doesn't support directories.
func (s *Sync) Update(key, jsonString string) error {
	defer measureTime(time.Now(), "update")

	var err error
	if jsonString == "" {
		// do nothing, because clientv3 doesn't have directories
		return nil
	}
	_, err = s.etcdClient.Put(s.withTimeout(), key, jsonString)
	if err != nil {
		log.Error(fmt.Sprintf("failed to sync with backend %s", err))
		updateCounter(1, "update.error")
		return err
	}
	return nil
}

//Delete sync update sync
func (s *Sync) Delete(key string, prefix bool) error {
	defer measureTime(time.Now(), "delete")

	opts := []etcd.OpOption{}
	if prefix {
		opts = append(opts, etcd.WithPrefix())
	}
	_, err := s.etcdClient.Delete(s.withTimeout(), key, opts...)
	if err != nil {
		updateCounter(1, "delete.error")
	}
	return err
}

//Fetch data from sync
func (s *Sync) Fetch(key string) (*sync.Node, error) {
	defer measureTime(time.Now(), "fetch")
	dir, err := s.etcdClient.Get(s.withTimeout(), key, etcd.WithPrefix(), etcd.WithSort(etcd.SortByKey, etcd.SortAscend))
	if err != nil {
		updateCounter(1, "fetch.dir.error")
		return nil, err
	}

	sep := "/"
	curr := strings.Count(key, sep)
	if curr == 0 {
		curr = 1
	}
	children := recursiveFetch(curr, dir.Kvs, key, sep)
	if children == nil {
		return nil, fmt.Errorf("Key not found (%s)", key)
	}
	return children, nil
}

func recursiveFetch(curr int, children []*pb.KeyValue, rootKey, sep string) *sync.Node {
	if len(children) == 0 {
		return nil
	}
	if len(children) == 1 {
		return handleSingleChild(curr, children[0], rootKey, sep)
	}

	key := substrN(string(children[0].Key), sep, curr)
	commonChild := make(map[string][]*pb.KeyValue)
	for _, child := range children {
		val := substrN(string(child.Key), sep, curr+1)
		commonChild[val] = append(commonChild[val], child)
	}
	// children nodes has to be alphabetically sorted
	keys := make([]string, 0, len(commonChild))
	for k, v := range commonChild {
		if len(v) != 0 {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	nodes := make([]*sync.Node, 0, len(keys))
	for _, key := range keys {
		v := commonChild[key]
		if node := recursiveFetch(curr+1, v, rootKey, sep); node != nil {
			nodes = append(nodes, node)
		}
	}
	node := &sync.Node{Key: key, Children: nodes}
	return node
}

func handleSingleChild(curr int, child *pb.KeyValue, rootKey, sep string) *sync.Node {
	key := substrN(string(child.Key), sep, curr)
	if string(child.Key) == key {
		if string(child.Key) != rootKey && string(child.Key[:len(rootKey)+1]) != rootKey+sep { // remove invalid keys
			return nil
		}
		return &sync.Node{Key: string(child.Key), Value: string(child.Value), Revision: child.ModRevision}
	}
	nodes := handleSingleChild(curr+1, child, rootKey, sep)
	return &sync.Node{Key: key, Children: []*sync.Node{nodes}}
}

func substrN(s, substr string, n int) string {
	idx := 1
	for i := 0; i < n; i++ {
		tmp := strings.Index(s[idx:], substr)
		if tmp == -1 {
			return s
		}
		idx += tmp + 1
	}
	return s[0 : idx-1]
}

//HasLock checks current process owns lock or not
func (s *Sync) HasLock(path string) bool {
	return s.locks.Has(path)
}

// Lock locks resources on sync
// This call blocks until you can get lock
func (s *Sync) Lock(path string, block bool) (chan struct{}, error) {
	defer measureTime(time.Now(), "lock")
	updateCounter(1, "lock.waiting")
	defer updateCounter(-1, "lock.waiting")

	for {
		var err error
		lease, err := s.etcdClient.Grant(s.withTimeout(), masterTTL)
		var resp *etcd.TxnResponse
		if err == nil {
			cmp := etcd.Compare(etcd.CreateRevision(path), "=", 0)
			put := etcd.OpPut(path, s.processID, etcd.WithLease(lease.ID))
			resp, err = s.etcdClient.Txn(s.withTimeout()).If(cmp).Then(put).Commit()
		}
		if err != nil || !resp.Succeeded {
			msg := fmt.Sprintf("failed to lock path %s", path)
			if err != nil {
				updateCounter(1, "lock.error")
				msg = fmt.Sprintf("failed to lock path %s: %s", path, err)
			}
			log.Notice(msg)

			if !block {
				return nil, errors.New(msg)
			}
			time.Sleep(masterTTL * time.Second)
			continue
		}
		s.locks.Set(path, lease.ID)
		log.Info("Locked %s", path)
		updateCounter(1, "lock.granted")

		//Refresh master token
		lost := make(chan struct{})
		go func() {
			defer s.abortLock(path)
			defer close(lost)
			defer updateCounter(-1, "lock.granted")

			for s.HasLock(path) {
				resp, err := s.etcdClient.KeepAliveOnce(s.withTimeout(), lease.ID)
				if err != nil || resp.TTL <= 0 {
					updateCounter(1, "lock.keepalive.error")
					log.Notice("failed to keepalive lock for %s %s", path, err)
					return
				}
				time.Sleep(masterTTL / 2 * time.Second)
			}
		}()

		return lost, nil
	}
}

func (s *Sync) abortLock(path string) etcd.LeaseID {
	leaseID, ok := s.locks.Get(path)
	if !ok {
		return 0
	}
	s.locks.Remove(path)
	log.Info("Unlocked path %s", path)
	return leaseID.(etcd.LeaseID)
}

//Unlock path
func (s *Sync) Unlock(path string) error {
	defer measureTime(time.Now(), "unlock")

	leaseID := s.abortLock(path)
	if leaseID > 0 {
		s.etcdClient.Revoke(s.withTimeout(), leaseID)

		cmp := etcd.Compare(etcd.Value(path), "=", s.processID)
		del := etcd.OpDelete(path)
		s.etcdClient.Txn(s.withTimeout()).If(cmp).Then(del).Commit()
	}
	return nil
}

func eventsFromNode(action string, kvs []*pb.KeyValue, responseChan chan *sync.Event, stopChan chan bool) {
	for _, kv := range kvs {
		event := &sync.Event{
			Action:   action,
			Key:      string(kv.Key),
			Revision: kv.ModRevision,
		}
		if kv.Value != nil {
			err := json.Unmarshal(kv.Value, &event.Data)
			if err != nil {
				log.Warning("failed to unmarshal watch response value %s: %s", kv.Value, err)
			}
		}
		select {
		case <-stopChan:
			log.Debug("Events from node interrupted by stop")
			return
		case responseChan <- event:
		}
	}
}

//Watch keep watch update under the path
func (s *Sync) Watch(path string, responseChan chan *sync.Event, stopChan chan bool, revision int64) error {
	options := []etcd.OpOption{etcd.WithPrefix(), etcd.WithSort(etcd.SortByModRevision, etcd.SortAscend)}
	if revision != sync.RevisionCurrent {
		options = append(options, etcd.WithMinModRev(revision))
	}
	node, err := s.etcdClient.Get(s.withTimeout(), path, options...)
	if err != nil {
		updateCounter(1, "watch.get.error")
		return err
	}
	eventsFromNode("get", node.Kvs, responseChan, stopChan)
	revision = node.Header.Revision + 1

	ctx, cancel := context.WithCancel(context.Background())
	errors := make(chan error, 1)
	var wg syn.WaitGroup
	wg.Add(1)
	go func() {
		updateCounter(1, "watch.active")
		defer updateCounter(-1, "watch.active")

		defer wg.Done()
		err := func() error {
			rch := s.etcdClient.Watch(ctx, path, etcd.WithPrefix(), etcd.WithRev(revision))

			for wresp := range rch {
				err := wresp.Err()
				if err != nil {
					updateCounter(1, "watch.client_watch.error")
					return err
				}
				for _, ev := range wresp.Events {
					action := "unknown"
					switch ev.Type {
					case etcd.EventTypePut:
						action = "set"
					case etcd.EventTypeDelete:
						action = "delete"
					}
					eventsFromNode(action, []*pb.KeyValue{ev.Kv}, responseChan, stopChan)
				}
			}

			return nil
		}()
		errors <- err
	}()
	defer func() {
		cancel()
		wg.Wait()
	}()

	// since Watch() doesn't close the returning channel even when
	// it gets an error, we need a side channel to see the connection state.
	session, err := concurrency.NewSession(s.etcdClient, concurrency.WithTTL(masterTTL))
	if err != nil {
		updateCounter(1, "watch.session.error")
		return err
	}
	defer session.Close()

	select {
	case <-session.Done():
		return fmt.Errorf("Watch aborted by etcd session close")
	case <-stopChan:
		return nil
	case err := <-errors:
		return err
	}
}

// WatchContext keep watch update under the path until context is canceled
func (s *Sync) WatchContext(ctx context.Context, path string, revision int64) <-chan *sync.Event {
	eventCh := make(chan *sync.Event, 32)
	stopCh := make(chan bool)
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Watch(path, eventCh, stopCh, revision)
	}()
	go func() {
		defer close(eventCh)

		select {
		case <-ctx.Done():
			close(stopCh)
			// don't return without ensuring Watch finished or we risk panic: send on closed channel
			<-errCh
		case err := <-errCh:
			close(stopCh)
			if err != nil {
				select {
				case eventCh <- &sync.Event{Err: err}:
				default:
					log.Debug("Unable to send error: '%s' via response chan. Don't linger.", err)
				}
			}
		}
	}()
	return eventCh
}

// Close closes etcd client
func (s *Sync) Close() {
	defer measureTime(time.Now(), "close")
	s.etcdClient.Close()
}
