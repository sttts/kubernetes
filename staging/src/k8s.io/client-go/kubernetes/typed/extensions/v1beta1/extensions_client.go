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

package v1beta1

import (
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

type ExtensionsV1beta1Interface interface {
	RESTClient() rest.Interface
	DaemonSetsGetter
	ScopedDaemonSetsGetter
	DeploymentsGetter
	ScopedDeploymentsGetter
	IngressesGetter
	ScopedIngressesGetter
	NetworkPoliciesGetter
	ScopedNetworkPoliciesGetter
	PodSecurityPoliciesGetter
	ScopedPodSecurityPoliciesGetter
	ReplicaSetsGetter
	ScopedReplicaSetsGetter
}

// ExtensionsV1beta1Client is used to interact with features provided by the extensions group.
type ExtensionsV1beta1Client struct {
	restClient rest.Interface
	cluster    string
}

func (c *ExtensionsV1beta1Client) DaemonSets(namespace string) DaemonSetInterface {
	return newDaemonSets(c, nil, namespace)
}

func (c *ExtensionsV1beta1Client) ScopedDaemonSets(scope rest.Scope, namespace string) DaemonSetInterface {
	return newDaemonSets(c, scope, namespace)
}

func (c *ExtensionsV1beta1Client) Deployments(namespace string) DeploymentInterface {
	return newDeployments(c, nil, namespace)
}

func (c *ExtensionsV1beta1Client) ScopedDeployments(scope rest.Scope, namespace string) DeploymentInterface {
	return newDeployments(c, scope, namespace)
}

func (c *ExtensionsV1beta1Client) Ingresses(namespace string) IngressInterface {
	return newIngresses(c, nil, namespace)
}

func (c *ExtensionsV1beta1Client) ScopedIngresses(scope rest.Scope, namespace string) IngressInterface {
	return newIngresses(c, scope, namespace)
}

func (c *ExtensionsV1beta1Client) NetworkPolicies(namespace string) NetworkPolicyInterface {
	return newNetworkPolicies(c, nil, namespace)
}

func (c *ExtensionsV1beta1Client) ScopedNetworkPolicies(scope rest.Scope, namespace string) NetworkPolicyInterface {
	return newNetworkPolicies(c, scope, namespace)
}

func (c *ExtensionsV1beta1Client) PodSecurityPolicies() PodSecurityPolicyInterface {
	return newPodSecurityPolicies(c, nil)
}

func (c *ExtensionsV1beta1Client) ScopedPodSecurityPolicies(scope rest.Scope) PodSecurityPolicyInterface {
	return newPodSecurityPolicies(c, scope)
}

func (c *ExtensionsV1beta1Client) ReplicaSets(namespace string) ReplicaSetInterface {
	return newReplicaSets(c, nil, namespace)
}

func (c *ExtensionsV1beta1Client) ScopedReplicaSets(scope rest.Scope, namespace string) ReplicaSetInterface {
	return newReplicaSets(c, scope, namespace)
}

// NewForConfig creates a new ExtensionsV1beta1Client for the given config.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*ExtensionsV1beta1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	httpClient, err := rest.HTTPClientFor(&config)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(&config, httpClient)
}

// NewForConfigAndClient creates a new ExtensionsV1beta1Client for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
func NewForConfigAndClient(c *rest.Config, h rest.HTTPClient) (*ExtensionsV1beta1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientForConfigAndClient(&config, h)
	if err != nil {
		return nil, err
	}
	return &ExtensionsV1beta1Client{restClient: client}, nil
}

// NewForConfigOrDie creates a new ExtensionsV1beta1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *ExtensionsV1beta1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new ExtensionsV1beta1Client for the given RESTClient.
func New(c rest.Interface) *ExtensionsV1beta1Client {
	return &ExtensionsV1beta1Client{restClient: c}
}

// NewWithCluster creates a new ExtensionsV1beta1Client for the given RESTClient and cluster.
func NewWithCluster(c rest.Interface, cluster string) *ExtensionsV1beta1Client {
	return &ExtensionsV1beta1Client{restClient: c, cluster: cluster}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1beta1.SchemeGroupVersion
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
func (c *ExtensionsV1beta1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
