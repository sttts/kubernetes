/*
Copyright 2018 The Kubernetes Authors.

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

package fake

import (
	"sync/atomic"

	"k8s.io/apiserver/pkg/admission/plugin/webhook/config"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/testcerts"
	"k8s.io/client-go/rest"
)

func NewAuthenticationInfoResolver(count *int32) config.AuthenticationInfoResolver {
	return &authenticationInfoResolver{
		restConfig: &rest.Config{
			TLSClientConfig: rest.TLSClientConfig{
				CAData:   testcerts.CACert,
				CertData: testcerts.ClientCert,
				KeyData:  testcerts.ClientKey,
			},
		},
		cachedCount: count,
	}
}

type authenticationInfoResolver struct {
	restConfig  *rest.Config
	cachedCount *int32
}

func (a *authenticationInfoResolver) ClientConfigFor(server string) (*rest.Config, error) {
	atomic.AddInt32(a.cachedCount, 1)
	return a.restConfig, nil
}

func (a *authenticationInfoResolver) ClientConfigForService(serviceName, serviceNamespace string) (*rest.Config, error) {
	atomic.AddInt32(a.cachedCount, 1)
	return a.restConfig, nil
}
