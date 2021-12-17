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

// Code generated by client-gen. DO NOT EDIT.

package internalversion

import (
	"fmt"

	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
	exampleinternalversion "k8s.io/code-generator/examples/apiserver/clientset/internalversion/typed/example/internalversion"
	secondexampleinternalversion "k8s.io/code-generator/examples/apiserver/clientset/internalversion/typed/example2/internalversion"
	thirdexampleinternalversion "k8s.io/code-generator/examples/apiserver/clientset/internalversion/typed/example3.io/internalversion"
)

type ClusterInterface interface {
	Cluster(name string) Interface
}

type Cluster struct {
	*scopedClientset
}

// Cluster sets the cluster for a Clientset.
func (c *Cluster) Cluster(name string) Interface {
	return &Clientset{
		scopedClientset: c.scopedClientset,
		cluster:         name,
	}
}

// NewClusterForConfig creates a new Cluster for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewClusterForConfig will generate a rate-limiter in configShallowCopy.
func NewClusterForConfig(c *rest.Config) (*Cluster, error) {
	cs, err := NewForConfig(c)
	if err != nil {
		return nil, err
	}
	return &Cluster{scopedClientset: cs.scopedClientset}, nil
}

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	Example() exampleinternalversion.ExampleInterface
	SecondExample() secondexampleinternalversion.SecondExampleInterface
	ThirdExample() thirdexampleinternalversion.ThirdExampleInterface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*scopedClientset
	cluster string
}

// scopedClientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type scopedClientset struct {
	*discovery.DiscoveryClient
	example       *exampleinternalversion.ExampleClient
	secondExample *secondexampleinternalversion.SecondExampleClient
	thirdExample  *thirdexampleinternalversion.ThirdExampleClient
}

// Example retrieves the ExampleClient
func (c *Clientset) Example() exampleinternalversion.ExampleInterface {
	return exampleinternalversion.NewWithCluster(c.example.RESTClient(), c.cluster)
}

// SecondExample retrieves the SecondExampleClient
func (c *Clientset) SecondExample() secondexampleinternalversion.SecondExampleInterface {
	return secondexampleinternalversion.NewWithCluster(c.secondExample.RESTClient(), c.cluster)
}

// ThirdExample retrieves the ThirdExampleClient
func (c *Clientset) ThirdExample() thirdexampleinternalversion.ThirdExampleInterface {
	return thirdexampleinternalversion.NewWithCluster(c.thirdExample.RESTClient(), c.cluster)
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient.WithCluster(c.cluster)
}

// NewForConfig creates a new Clientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfig will generate a rate-limiter in configShallowCopy.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c

	// share the transport between all clients
	httpClient, err := rest.HTTPClientFor(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	return NewForConfigAndClient(&configShallowCopy, httpClient)
}

// NewForConfigAndClient creates a new Clientset for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfigAndClient will generate a rate-limiter in configShallowCopy.
func NewForConfigAndClient(c *rest.Config, httpClient rest.HTTPClient) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}

	var cs scopedClientset
	var err error
	cs.example, err = exampleinternalversion.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	cs.secondExample, err = secondexampleinternalversion.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	cs.thirdExample, err = thirdexampleinternalversion.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	return &Clientset{scopedClientset: &cs}, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	cs, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs scopedClientset
	cs.example = exampleinternalversion.New(c)
	cs.secondExample = secondexampleinternalversion.New(c)
	cs.thirdExample = thirdexampleinternalversion.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &Clientset{scopedClientset: &cs}
}
