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
	"context"

	kcpclient "github.com/kcp-dev/apimachinery/pkg/client"
	"github.com/kcp-dev/logicalcluster/v3"

	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	batchv1beta1client "k8s.io/client-go/kubernetes/typed/batch/v1beta1"
)

// CronJobsClusterGetter has a method to return a CronJobClusterInterface.
// A group's cluster client should implement this interface.
type CronJobsClusterGetter interface {
	CronJobs() CronJobClusterInterface
}

// CronJobClusterInterface can operate on CronJobs across all clusters,
// or scope down to one cluster and return a CronJobsNamespacer.
type CronJobClusterInterface interface {
	Cluster(logicalcluster.Path) CronJobsNamespacer
	List(ctx context.Context, opts metav1.ListOptions) (*batchv1beta1.CronJobList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type cronJobsClusterInterface struct {
	clientCache kcpclient.Cache[*batchv1beta1client.BatchV1beta1Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *cronJobsClusterInterface) Cluster(name logicalcluster.Path) CronJobsNamespacer {
	if name == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return &cronJobsNamespacer{clientCache: c.clientCache, name: name}
}

// List returns the entire collection of all CronJobs across all clusters.
func (c *cronJobsClusterInterface) List(ctx context.Context, opts metav1.ListOptions) (*batchv1beta1.CronJobList, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).CronJobs(metav1.NamespaceAll).List(ctx, opts)
}

// Watch begins to watch all CronJobs across all clusters.
func (c *cronJobsClusterInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).CronJobs(metav1.NamespaceAll).Watch(ctx, opts)
}

// CronJobsNamespacer can scope to objects within a namespace, returning a batchv1beta1client.CronJobInterface.
type CronJobsNamespacer interface {
	Namespace(string) batchv1beta1client.CronJobInterface
}

type cronJobsNamespacer struct {
	clientCache kcpclient.Cache[*batchv1beta1client.BatchV1beta1Client]
	name        logicalcluster.Path
}

func (n *cronJobsNamespacer) Namespace(namespace string) batchv1beta1client.CronJobInterface {
	return n.clientCache.ClusterOrDie(n.name).CronJobs(namespace)
}
