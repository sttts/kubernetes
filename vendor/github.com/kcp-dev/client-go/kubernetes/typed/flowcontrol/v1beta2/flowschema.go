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
	"context"

	kcpclient "github.com/kcp-dev/apimachinery/pkg/client"
	"github.com/kcp-dev/logicalcluster/v3"

	flowcontrolv1beta2 "k8s.io/api/flowcontrol/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	flowcontrolv1beta2client "k8s.io/client-go/kubernetes/typed/flowcontrol/v1beta2"
)

// FlowSchemasClusterGetter has a method to return a FlowSchemaClusterInterface.
// A group's cluster client should implement this interface.
type FlowSchemasClusterGetter interface {
	FlowSchemas() FlowSchemaClusterInterface
}

// FlowSchemaClusterInterface can operate on FlowSchemas across all clusters,
// or scope down to one cluster and return a flowcontrolv1beta2client.FlowSchemaInterface.
type FlowSchemaClusterInterface interface {
	Cluster(logicalcluster.Path) flowcontrolv1beta2client.FlowSchemaInterface
	List(ctx context.Context, opts metav1.ListOptions) (*flowcontrolv1beta2.FlowSchemaList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type flowSchemasClusterInterface struct {
	clientCache kcpclient.Cache[*flowcontrolv1beta2client.FlowcontrolV1beta2Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *flowSchemasClusterInterface) Cluster(path logicalcluster.Path) flowcontrolv1beta2client.FlowSchemaInterface {
	if path == logicalcluster.WildcardPath {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return c.clientCache.ClusterOrDie(path).FlowSchemas()
}

// List returns the entire collection of all FlowSchemas across all clusters.
func (c *flowSchemasClusterInterface) List(ctx context.Context, opts metav1.ListOptions) (*flowcontrolv1beta2.FlowSchemaList, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.WildcardPath).FlowSchemas().List(ctx, opts)
}

// Watch begins to watch all FlowSchemas across all clusters.
func (c *flowSchemasClusterInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.WildcardPath).FlowSchemas().Watch(ctx, opts)
}
