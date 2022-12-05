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

package v1beta2

import (
	"net/http"

	kcpclient "github.com/kcp-dev/apimachinery/pkg/client"
	"github.com/kcp-dev/logicalcluster/v3"

	appsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	"k8s.io/client-go/rest"
)

type AppsV1beta2ClusterInterface interface {
	AppsV1beta2ClusterScoper
	StatefulSetsClusterGetter
	DeploymentsClusterGetter
	DaemonSetsClusterGetter
	ReplicaSetsClusterGetter
	ControllerRevisionsClusterGetter
}

type AppsV1beta2ClusterScoper interface {
	Cluster(logicalcluster.Path) appsv1beta2.AppsV1beta2Interface
}

type AppsV1beta2ClusterClient struct {
	clientCache kcpclient.Cache[*appsv1beta2.AppsV1beta2Client]
}

func (c *AppsV1beta2ClusterClient) Cluster(name logicalcluster.Path) appsv1beta2.AppsV1beta2Interface {
	if name == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}
	return c.clientCache.ClusterOrDie(name)
}

func (c *AppsV1beta2ClusterClient) StatefulSets() StatefulSetClusterInterface {
	return &statefulSetsClusterInterface{clientCache: c.clientCache}
}

func (c *AppsV1beta2ClusterClient) Deployments() DeploymentClusterInterface {
	return &deploymentsClusterInterface{clientCache: c.clientCache}
}

func (c *AppsV1beta2ClusterClient) DaemonSets() DaemonSetClusterInterface {
	return &daemonSetsClusterInterface{clientCache: c.clientCache}
}

func (c *AppsV1beta2ClusterClient) ReplicaSets() ReplicaSetClusterInterface {
	return &replicaSetsClusterInterface{clientCache: c.clientCache}
}

func (c *AppsV1beta2ClusterClient) ControllerRevisions() ControllerRevisionClusterInterface {
	return &controllerRevisionsClusterInterface{clientCache: c.clientCache}
}

// NewForConfig creates a new AppsV1beta2ClusterClient for the given config.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*AppsV1beta2ClusterClient, error) {
	client, err := rest.HTTPClientFor(c)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(c, client)
}

// NewForConfigAndClient creates a new AppsV1beta2ClusterClient for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
func NewForConfigAndClient(c *rest.Config, h *http.Client) (*AppsV1beta2ClusterClient, error) {
	cache := kcpclient.NewCache(c, h, &kcpclient.Constructor[*appsv1beta2.AppsV1beta2Client]{
		NewForConfigAndClient: appsv1beta2.NewForConfigAndClient,
	})
	if _, err := cache.Cluster(logicalcluster.New("root")); err != nil {
		return nil, err
	}
	return &AppsV1beta2ClusterClient{clientCache: cache}, nil
}

// NewForConfigOrDie creates a new AppsV1beta2ClusterClient for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *AppsV1beta2ClusterClient {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}
