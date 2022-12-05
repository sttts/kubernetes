/*
Copyright The Kubernetes Authors.

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


//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by kcp code-generator. DO NOT EDIT.

package clientset

import (
	"fmt"
	"net/http"

	kcpclient "github.com/kcp-dev/apimachinery/pkg/client"
	"github.com/kcp-dev/logicalcluster/v3"

	client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/flowcontrol"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/client/kcp/clientset/versioned/typed/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/client/kcp/clientset/versioned/typed/apiextensions/v1beta1"
)

type ClusterInterface interface {
	Cluster(logicalcluster.Path) client.Interface
	Discovery() discovery.DiscoveryInterface
	ApiextensionsV1() apiextensionsv1.ApiextensionsV1ClusterInterface
	ApiextensionsV1beta1() apiextensionsv1beta1.ApiextensionsV1beta1ClusterInterface
}

// ClusterClientset contains the clients for groups.
type ClusterClientset struct {
	*discovery.DiscoveryClient
	clientCache kcpclient.Cache[*client.Clientset]
	apiextensionsV1 *apiextensionsv1.ApiextensionsV1ClusterClient
	apiextensionsV1beta1 *apiextensionsv1beta1.ApiextensionsV1beta1ClusterClient
}

// Discovery retrieves the DiscoveryClient
func (c *ClusterClientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}


// ApiextensionsV1 retrieves the ApiextensionsV1ClusterClient.  
func (c *ClusterClientset) ApiextensionsV1() apiextensionsv1.ApiextensionsV1ClusterInterface {
	return c.apiextensionsV1
}

// ApiextensionsV1beta1 retrieves the ApiextensionsV1beta1ClusterClient.  
func (c *ClusterClientset) ApiextensionsV1beta1() apiextensionsv1beta1.ApiextensionsV1beta1ClusterInterface {
	return c.apiextensionsV1beta1
}
// Cluster scopes this clientset to one cluster.
func (c *ClusterClientset) Cluster(path logicalcluster.Path) client.Interface {
	if path == logicalcluster.WildcardPath {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}
	return c.clientCache.ClusterOrDie(path)
}

// NewForConfig creates a new ClusterClientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable, 
// NewForConfig will generate a rate-limiter in configShallowCopy.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*ClusterClientset, error) {
	configShallowCopy := *c

	if configShallowCopy.UserAgent == "" {
		configShallowCopy.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	// share the transport between all clients
	httpClient, err := rest.HTTPClientFor(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	return NewForConfigAndClient(&configShallowCopy, httpClient)
}

// NewForConfigAndClient creates a new ClusterClientset for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfigAndClient will generate a rate-limiter in configShallowCopy.
func NewForConfigAndClient(c *rest.Config, httpClient *http.Client) (*ClusterClientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}

	cache := kcpclient.NewCache(c, httpClient, &kcpclient.Constructor[*client.Clientset]{
		NewForConfigAndClient: client.NewForConfigAndClient,
	})
	if _, err := cache.Cluster(logicalcluster.New("root")); err != nil {
		return nil, err
	}

	var cs ClusterClientset
	cs.clientCache = cache
	var err error
    cs.apiextensionsV1, err = apiextensionsv1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
    cs.apiextensionsV1beta1, err = apiextensionsv1beta1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new ClusterClientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *ClusterClientset {
	cs, err := NewForConfig(c)
	if err!=nil {
		panic(err)
	}
	return cs
}
