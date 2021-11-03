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

package v1

import (
	v1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

type NetworkingV1Interface interface {
	RESTClient() rest.Interface
	IngressesGetter
	IngressClassesGetter
	NetworkPoliciesGetter
}

// NetworkingV1Client is used to interact with features provided by the networking.k8s.io group.
type NetworkingV1Client struct {
	restClient rest.Interface
	cluster    string
}

func (c *NetworkingV1Client) Ingresses(namespace string) IngressInterface {
	return newIngresses(c, namespace)
}

func (c *NetworkingV1Client) IngressClasses() IngressClassInterface {
	return newIngressClasses(c)
}

func (c *NetworkingV1Client) NetworkPolicies(namespace string) NetworkPolicyInterface {
	return newNetworkPolicies(c, namespace)
}

// NewForConfig creates a new NetworkingV1Client for the given config.
func NewForConfig(c *rest.Config) (*NetworkingV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &NetworkingV1Client{restClient: client}, nil
}

// NewForConfigOrDie creates a new NetworkingV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *NetworkingV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new NetworkingV1Client for the given RESTClient.
func New(c rest.Interface) *NetworkingV1Client {
	return &NetworkingV1Client{restClient: c}
}

// NewWithCluster creates a new NetworkingV1Client for the given RESTClient and cluster.
func NewWithCluster(c rest.Interface, cluster string) *NetworkingV1Client {
	return &NetworkingV1Client{restClient: c, cluster: cluster}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *NetworkingV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
