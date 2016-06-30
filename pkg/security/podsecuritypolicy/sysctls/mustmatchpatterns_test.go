/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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
	"testing"

	"k8s.io/kubernetes/pkg/api"
)

func TestValidate(t *testing.T) {
	tests := map[string]struct {
		patterns   []string
		allowed    []string
		disallowed []string
	}{
		// no container requests
		"nil": {
			patterns: nil,
			allowed:  []string{"foo"},
		},
		"empty": {
			patterns:   []string{},
			disallowed: []string{"foo"},
		},
		"without wildcard": {
			patterns:   []string{"a", "a.b"},
			allowed:    []string{"a", "a.b"},
			disallowed: []string{"b"},
		},
		"with catch-all wildcard": {
			patterns: []string{"*"},
			allowed:  []string{"a", "a.b"},
		},
		"with catch-all wildcard and non-wildcard": {
			patterns: []string{"a.b.c", "*"},
			allowed:  []string{"a", "a.b", "a.b.c", "b"},
		},
		"without catch-all wildcard": {
			patterns:   []string{"a.*", "b.*", "c.d.e", "d.e.f.*"},
			allowed:    []string{"a.b", "b.c", "c.d.e", "d.e.f.g.h"},
			disallowed: []string{"a", "b", "c", "c.d", "d.e", "d.e.f"},
		},
	}

	for k, v := range tests {
		strategy, err := NewMustMatchPatterns(v.patterns)
		if err != nil {
			t.Errorf("%s failed: %v", k, err)
			continue
		}

		pod := &api.Pod{
			Spec: api.PodSpec{
				SecurityContext: &api.PodSecurityContext{
					Sysctls: []api.Sysctl{},
				},
			},
		}
		errs := strategy.Validate(pod)
		if len(errs) != 0 {
			t.Errorf("%s: unexpected validaton errors for empty sysctls: %v", k, errs)
		}

		pod.Spec.SecurityContext.Sysctls = []api.Sysctl{}
		for _, s := range v.allowed {
			pod.Spec.SecurityContext.Sysctls = append(pod.Spec.SecurityContext.Sysctls, api.Sysctl{s, "dummy"})
		}
		errs = strategy.Validate(pod)
		if len(errs) != 0 {
			t.Errorf("%s: unexpected validaton errors: %v", k, errs)
		}

		for _, s := range v.disallowed {
			pod.Spec.SecurityContext.Sysctls = []api.Sysctl{{s, "dummy"}}
			errs = strategy.Validate(pod)
			if len(errs) == 0 {
				t.Errorf("%s: expected error for sysctl %q", k, s)
			}
		}
	}
}
