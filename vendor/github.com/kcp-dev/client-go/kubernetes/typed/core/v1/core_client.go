//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The KCP Authors.

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

// Code generated by kcp code-generator. DO NOT EDIT.

package v1

import (
	"net/http"

	kcpclient "github.com/kcp-dev/apimachinery/pkg/client"
	"github.com/kcp-dev/logicalcluster/v3"

	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

type CoreV1ClusterInterface interface {
	CoreV1ClusterScoper
	PersistentVolumesClusterGetter
	PersistentVolumeClaimsClusterGetter
	PodsClusterGetter
	PodTemplatesClusterGetter
	ReplicationControllersClusterGetter
	ServicesClusterGetter
	ServiceAccountsClusterGetter
	EndpointsClusterGetter
	NodesClusterGetter
	NamespacesClusterGetter
	EventsClusterGetter
	LimitRangesClusterGetter
	ResourceQuotasClusterGetter
	SecretsClusterGetter
	ConfigMapsClusterGetter
	ComponentStatusesClusterGetter
}

type CoreV1ClusterScoper interface {
	Cluster(logicalcluster.Path) corev1.CoreV1Interface
}

type CoreV1ClusterClient struct {
	clientCache kcpclient.Cache[*corev1.CoreV1Client]
}

func (c *CoreV1ClusterClient) Cluster(path logicalcluster.Path) corev1.CoreV1Interface {
	if path == logicalcluster.WildcardPath {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}
	return c.clientCache.ClusterOrDie(path)
}

func (c *CoreV1ClusterClient) PersistentVolumes() PersistentVolumeClusterInterface {
	return &persistentVolumesClusterInterface{clientCache: c.clientCache}
}

func (c *CoreV1ClusterClient) PersistentVolumeClaims() PersistentVolumeClaimClusterInterface {
	return &persistentVolumeClaimsClusterInterface{clientCache: c.clientCache}
}

func (c *CoreV1ClusterClient) Pods() PodClusterInterface {
	return &podsClusterInterface{clientCache: c.clientCache}
}

func (c *CoreV1ClusterClient) PodTemplates() PodTemplateClusterInterface {
	return &podTemplatesClusterInterface{clientCache: c.clientCache}
}

func (c *CoreV1ClusterClient) ReplicationControllers() ReplicationControllerClusterInterface {
	return &replicationControllersClusterInterface{clientCache: c.clientCache}
}

func (c *CoreV1ClusterClient) Services() ServiceClusterInterface {
	return &servicesClusterInterface{clientCache: c.clientCache}
}

func (c *CoreV1ClusterClient) ServiceAccounts() ServiceAccountClusterInterface {
	return &serviceAccountsClusterInterface{clientCache: c.clientCache}
}

func (c *CoreV1ClusterClient) Endpoints() EndpointsClusterInterface {
	return &endpointsClusterInterface{clientCache: c.clientCache}
}

func (c *CoreV1ClusterClient) Nodes() NodeClusterInterface {
	return &nodesClusterInterface{clientCache: c.clientCache}
}

func (c *CoreV1ClusterClient) Namespaces() NamespaceClusterInterface {
	return &namespacesClusterInterface{clientCache: c.clientCache}
}

func (c *CoreV1ClusterClient) Events() EventClusterInterface {
	return &eventsClusterInterface{clientCache: c.clientCache}
}

func (c *CoreV1ClusterClient) LimitRanges() LimitRangeClusterInterface {
	return &limitRangesClusterInterface{clientCache: c.clientCache}
}

func (c *CoreV1ClusterClient) ResourceQuotas() ResourceQuotaClusterInterface {
	return &resourceQuotasClusterInterface{clientCache: c.clientCache}
}

func (c *CoreV1ClusterClient) Secrets() SecretClusterInterface {
	return &secretsClusterInterface{clientCache: c.clientCache}
}

func (c *CoreV1ClusterClient) ConfigMaps() ConfigMapClusterInterface {
	return &configMapsClusterInterface{clientCache: c.clientCache}
}

func (c *CoreV1ClusterClient) ComponentStatuses() ComponentStatusClusterInterface {
	return &componentStatusesClusterInterface{clientCache: c.clientCache}
}

// NewForConfig creates a new CoreV1ClusterClient for the given config.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*CoreV1ClusterClient, error) {
	client, err := rest.HTTPClientFor(c)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(c, client)
}

// NewForConfigAndClient creates a new CoreV1ClusterClient for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
func NewForConfigAndClient(c *rest.Config, h *http.Client) (*CoreV1ClusterClient, error) {
	cache := kcpclient.NewCache(c, h, &kcpclient.Constructor[*corev1.CoreV1Client]{
		NewForConfigAndClient: corev1.NewForConfigAndClient,
	})
	if _, err := cache.Cluster(logicalcluster.New("root")); err != nil {
		return nil, err
	}
	return &CoreV1ClusterClient{clientCache: cache}, nil
}

// NewForConfigOrDie creates a new CoreV1ClusterClient for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *CoreV1ClusterClient {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}
