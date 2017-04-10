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

package server

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cloudwan/gohan/extension"
	"github.com/cloudwan/gohan/job"

	l "github.com/cloudwan/gohan/log"
	gohan_sync "github.com/cloudwan/gohan/sync"
	"github.com/cloudwan/gohan/util"
)

const (
	//StateUpdateEventName used in etcd path
	StateUpdateEventName = "state_update"
	//MonitoringUpdateEventName used in etcd path
	MonitoringUpdateEventName = "monitoring_update"

	SyncWatchRevisionPrefix = "/gohan/watch/revision"

	processPathPrefix = "/gohan/cluster/process"

	masterTTL = 10
)

//Sync Watch Process
func startSyncWatchProcess(server *Server, chProcessWatch chan *gohan_sync.Event) {
	config := util.GetConfig()
	watch := config.GetStringList("watch/keys", nil)
	events := config.GetStringList("watch/events", nil)
	if watch == nil {
		return
	}
	extensions := map[string]extension.Environment{}
	for _, event := range events {
		path := "sync://" + event
		env, err := server.NewEnvironmentForPath("sync."+event, path)
		if err != nil {
			log.Fatal(err.Error())
		}
		extensions[event] = env
	}

	responseChans := make(map[string]chan *gohan_sync.Event)
	stopChans := make(map[string]chan bool)
	var processList []string

	// wait group to avoid map concurrent access
	wg := &sync.WaitGroup{}
	for idx, path := range watch {
		// new goroutine launch for new watching path
		responseChans[path] = make(chan *gohan_sync.Event)
		stopChans[path] = make(chan bool)

		// goroutine to decide watching keys with the consideration of LB
		wg.Add(1)
		go func(idx int, path string) {
			defer l.LogFatalPanic(log)
			responseChan := responseChans[path]
			stopChan := stopChans[path]
			wg.Done()
			for server.running {
				func() {
					size := len(processList)
					if (size == 0) {
						log.Debug("Wait until at least one process launch")
						time.Sleep(5 * time.Second)
						return
					}
					log.Debug(fmt.Sprintf("Start calculation of priority for path: %s", path))
					// index of the process in this cluster
					num := -1
					for p, v := range processList {
						if v == processPathPrefix + "/" + server.sync.GetProcessID() {
							num = p
							break
						}
					}
					// priority is 0-origin and lower number has higher priority
					priority := (num - (idx % size) + size) % size
					log.Debug(fmt.Sprintf("Calculated priority for path %s is %d", path, priority))

					lockKey := lockPath + "/watch" + path

					var err error
					if priority == 0 {
						// try to acquire lock greedily for the highest priority path
						err = server.sync.Lock(lockKey, false)
					}
					if err != nil {
						// wait based on priority time
						timeChan := time.NewTimer(time.Duration(masterTTL * (priority + 1)) * time.Second).C
						select {
							case <- timeChan:
							case <- stopChan:
								log.Debug("Wait interrupted because cluster change is detected")
								return
						}
						err = server.sync.Lock(lockKey, false)
						if err != nil {
							log.Debug("Can't start watch process because lock is already acquired: %s", err)
							return
						}
					}
					defer server.sync.Unlock(lockKey)

					// start watch process for locked path
					fromRevision := int64(gohan_sync.RevisionCurrent)
					lastSeen, err := server.sync.Fetch(SyncWatchRevisionPrefix + path)
					if err == nil {
						inStore, err := strconv.ParseInt(lastSeen.Value, 10, 64)
						if err == nil {
							log.Info("Using last seen revision `%d` for watching path `%s`", inStore, path)
							fromRevision = inStore
						}
					}

					err = server.sync.Watch(path, responseChan, stopChan, fromRevision)
					if err != nil {
						log.Error(fmt.Sprintf("sync watch error: %s", err))
					}
				}()
			}
		}(idx, path)
		wg.Wait()

		// goroutine to trigger notification event
		wg.Add(1)
		go func(path string) {
			defer l.LogFatalPanic(log)
			responseChan := responseChans[path]
			wg.Done()
			for server.running {
				response := <-responseChan
				err := server.sync.Update(SyncWatchRevisionPrefix+path, strconv.FormatInt(response.Revision, 10))
				if err != nil {
					log.Error("Failed to update revision number for watch path `%s` in sync storage", path)
				}
				server.queue.Add(job.NewJob(
					func() {
						defer l.LogPanic(log)
						for _, event := range events {
							//match extensions
							if strings.HasPrefix(response.Key, "/"+event) {
								env := extensions[event]
								runExtensionOnSync(server, response, env.Clone())
								return
							}
						}
					}),
				)
			}
		}(path)
		wg.Wait()
	}

	// goroutine to wait cluster change then stop sync watch for rebalance
	go func() {
		for server.running {
			processWatchEvent := <- chProcessWatch
			log.Debug(fmt.Sprintf("cluster change detected: %s process %s", processWatchEvent.Action, processWatchEvent.Key))

			// modify gohan process list
			pos := -1
			for p, v := range processList {
				if v == processWatchEvent.Key {
					pos = p
				}
			}
			if processWatchEvent.Action == "delete" {
				// remove detected process from list
				if pos > -1 {
					processList = append((processList)[:pos], (processList)[pos+1:]...)
				}
			} else {
				// add detected process from list
				if pos == -1 {
					processList = append(processList, processWatchEvent.Key)
					sort.Sort(sort.StringSlice(processList))
				}
			}
			log.Debug(fmt.Sprintf("Current cluster consists of following processes: %s", processList))

			// Forcibly stop current watch process then re-acquire lock based on newly calculated priority
			for _, path := range watch {
				stopChans[path] <- true
			}
		}
	}()
}

//Stop Watch Process
func stopSyncWatchProcess(server *Server) {
}

//Run extension on sync
func runExtensionOnSync(server *Server, response *gohan_sync.Event, env extension.Environment) {
	context := map[string]interface{}{
		"action": response.Action,
		"data":   response.Data,
		"key":    response.Key,
	}
	if err := env.HandleEvent("notification", context); err != nil {
		log.Warning(fmt.Sprintf("extension error: %s", err))
		return
	}
	return
}
