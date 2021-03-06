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

package util

import (
	"fmt"
	"bytes"
)

// CollectData creates string representing go file
func CollectData(name string, data []string) string {
	var buffer bytes.Buffer
	prefix := "package " + name + "\n"
	buffer.WriteString(prefix)
	for _, element := range data {
		if element != "" {
			buffer.WriteString(fmt.Sprintf("\n%s", element))
		}
	}
	result := buffer.String()
	if result == prefix {
		return ""
	}
	return result
}

func Const(data []string) string {
	var buffer bytes.Buffer
	buffer.WriteString("const (\n")
	for _, value := range data {
		buffer.WriteString(fmt.Sprintf("\t%s\n", value))
	}
	buffer.WriteString(")\n")
	return buffer.String()
}
