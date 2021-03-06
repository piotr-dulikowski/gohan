// Copyright (C) 2017 NTT Innovation Institute, Inc.
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


// Sample Gohan Extension

// You need to register event handler

gohan_register_handler("notification", function(context){
  var currentDate = Date.now();
  var writePath = "/test/" + currentDate;

  if ("pref" in context.data) {
    writePath = "/test/" + context.data["pref"] + "/" + currentDate;
  }
  gohan_sync_update(writePath, context.key);

  // wait 5 sec
  gohan_sleep(5000);
});
