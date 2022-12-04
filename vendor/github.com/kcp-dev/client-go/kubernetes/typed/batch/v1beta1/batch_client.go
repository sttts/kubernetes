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

package v1beta1

import (
	"net/http"

	kcpclient "github.com/kcp-dev/apimachinery/pkg/client"
	"github.com/kcp-dev/logicalcluster/v3"

	batchv1beta1 "k8s.io/client-go/kubernetes/typed/batch/v1beta1"
	"k8s.io/client-go/rest"
)

type BatchV1beta1ClusterInterface interface {
	BatchV1beta1ClusterScoper
	CronJobsClusterGetter
}

type BatchV1beta1ClusterScoper interface {
	Cluster(logicalcluster.Path) batchv1beta1.BatchV1beta1Interface
}

type BatchV1beta1ClusterClient struct {
	clientCache kcpclient.Cache[*batchv1beta1.BatchV1beta1Client]
}

func (c *BatchV1beta1ClusterClient) Cluster(path logicalcluster.Path) batchv1beta1.BatchV1beta1Interface {
	if path == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}
	return c.clientCache.ClusterOrDie(path)
}

func (c *BatchV1beta1ClusterClient) CronJobs() CronJobClusterInterface {
	return &cronJobsClusterInterface{clientCache: c.clientCache}
}

// NewForConfig creates a new BatchV1beta1ClusterClient for the given config.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*BatchV1beta1ClusterClient, error) {
	client, err := rest.HTTPClientFor(c)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(c, client)
}

// NewForConfigAndClient creates a new BatchV1beta1ClusterClient for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
func NewForConfigAndClient(c *rest.Config, h *http.Client) (*BatchV1beta1ClusterClient, error) {
	cache := kcpclient.NewCache(c, h, &kcpclient.Constructor[*batchv1beta1.BatchV1beta1Client]{
		NewForConfigAndClient: batchv1beta1.NewForConfigAndClient,
	})
	if _, err := cache.Cluster(logicalcluster.New("root")); err != nil {
		return nil, err
	}
	return &BatchV1beta1ClusterClient{clientCache: cache}, nil
}

// NewForConfigOrDie creates a new BatchV1beta1ClusterClient for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *BatchV1beta1ClusterClient {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}
