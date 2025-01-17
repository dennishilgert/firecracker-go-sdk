// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package firecracker

import "strings"

// kernelArgs serializes+deserializes kernel boot parameters from/into a map.
// Kernel docs: https://www.kernel.org/doc/Documentation/admin-guide/kernel-parameters.txt
//
// "key=value" will result in map["key"] = &"value"
// "key=" will result in map["key"] = &""
// "key" will result in map["key"] = nil
type kernelArgs map[string]*string

// serialize the kernelArgs back to a string that can be provided
// to the kernel
func (kargs kernelArgs) String() string {
	var fields []string
	var initField *string
	for key, value := range kargs {
		field := key
		if value != nil {
			field += "=" + *value
		}
		if key == "init" {
			initField = &field
			continue
		}
		fields = append(fields, field)
	}

	kernelArgsString := strings.Join(fields, " ")

	// add init as last field because otherwise this would screw up the entire kernel args
	if initField != nil {
		kernelArgsString += " " + *initField
	}

	return kernelArgsString
}

// deserialize the provided string to a kernelArgs map
func parseKernelArgs(rawString string) kernelArgs {
	argMap := make(map[string]*string)
	var initArgs *string
	if strings.Contains(rawString, "init=") {
		initSplit := strings.SplitN(rawString, "init=", 2)
		rawString = strings.TrimSpace(initSplit[0])
		initArgs = &initSplit[1]
	}
	for _, kv := range strings.Fields(rawString) {
		// only split into up to 2 fields (before and after the first "=")
		kvSplit := strings.SplitN(kv, "=", 2)

		key := kvSplit[0]

		var value *string
		if len(kvSplit) == 2 {
			value = &kvSplit[1]
		}

		argMap[key] = value
	}

	// add init separately because otherways this would screw up the entire kernel args
	if initArgs != nil {
		argMap["init"] = initArgs
	}

	return argMap
}
