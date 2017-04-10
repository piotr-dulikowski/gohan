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
	"time"

	gohan_sync "github.com/cloudwan/gohan/sync"
)

const processPath = "/gohan/cluster/process"

//Watch Process Availability
func startProcessWatchProcess(server *Server) chan *gohan_sync.Event {
	responseChan := make(chan *gohan_sync.Event)

	go func() {
		for server.running {
			fromRevision := int64(gohan_sync.RevisionCurrent)
			// register self process to the cluster
			lockKey := processPath + "/" + server.sync.GetProcessID()
			server.sync.Lock(lockKey, true)
			// start process watch
			err := server.sync.Watch(processPath, responseChan, make(chan bool), fromRevision)
			if err != nil {
				log.Error(fmt.Sprintf("process watch error: %s", err))
				time.Sleep(5 * time.Second)
			}
			server.sync.Unlock(lockKey)
		}
	}()
	return responseChan
}

//Stop Watch Process
func stopProcessWatchProcess(server *Server) {
}
