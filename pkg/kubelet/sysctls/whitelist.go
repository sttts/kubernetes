/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package sysctls

import (
	"strings"
)

var whitelist = map[string]struct{}{
	"kernel.msgmax":          {},
	"kernel.msgmnb":          {},
	"kernel.msgmni":          {},
	"kernel.sem":             {},
	"kernel.shmall":          {},
	"kernel.shmmax":          {},
	"kernel.shmmni":          {},
	"kernel.shm_rmid_forced": {},
}

var whitelistPrefixes = []string{
	"net.",
	"fs.mqueue.",
}

// Whitelisted checks that a sysctl is whitelisted because it is known
// to be namespaced by the Linux kernel. Note that being whitelisted is required, but not
// sufficient: the container runtime might have a stricter check and refuse to launch a pod.
func Whitelisted(val string) bool {
	if _, found := whitelist[val]; found {
		return true
	}
	for _, p := range whitelistPrefixes {
		if strings.HasPrefix(val, p) {
			return true
		}
	}
	return false
}
