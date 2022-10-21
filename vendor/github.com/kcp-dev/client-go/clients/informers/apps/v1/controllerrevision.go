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
	"context"
	"time"

	kcpcache "github.com/kcp-dev/apimachinery/pkg/cache"
	kcpinformers "github.com/kcp-dev/apimachinery/third_party/informers"
	"github.com/kcp-dev/logicalcluster/v2"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	upstreamappsv1informers "k8s.io/client-go/informers/apps/v1"
	upstreamappsv1listers "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"

	clientset "github.com/kcp-dev/client-go/clients/clientset/versioned"
	"github.com/kcp-dev/client-go/clients/informers/internalinterfaces"
	appsv1listers "github.com/kcp-dev/client-go/clients/listers/apps/v1"
)

// ControllerRevisionClusterInformer provides access to a shared informer and lister for
// ControllerRevisions.
type ControllerRevisionClusterInformer interface {
	Cluster(logicalcluster.Name) upstreamappsv1informers.ControllerRevisionInformer
	Informer() kcpcache.ScopeableSharedIndexInformer
	Lister() appsv1listers.ControllerRevisionClusterLister
}

type controllerRevisionClusterInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewControllerRevisionClusterInformer constructs a new informer for ControllerRevision type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewControllerRevisionClusterInformer(client clientset.ClusterInterface, resyncPeriod time.Duration, indexers cache.Indexers) kcpcache.ScopeableSharedIndexInformer {
	return NewFilteredControllerRevisionClusterInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredControllerRevisionClusterInformer constructs a new informer for ControllerRevision type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredControllerRevisionClusterInformer(client clientset.ClusterInterface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) kcpcache.ScopeableSharedIndexInformer {
	return kcpinformers.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AppsV1().ControllerRevisions().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AppsV1().ControllerRevisions().Watch(context.TODO(), options)
			},
		},
		&appsv1.ControllerRevision{},
		resyncPeriod,
		indexers,
	)
}

func (f *controllerRevisionClusterInformer) defaultInformer(client clientset.ClusterInterface, resyncPeriod time.Duration) kcpcache.ScopeableSharedIndexInformer {
	return NewFilteredControllerRevisionClusterInformer(client, resyncPeriod, cache.Indexers{
		kcpcache.ClusterIndexName:             kcpcache.ClusterIndexFunc,
		kcpcache.ClusterAndNamespaceIndexName: kcpcache.ClusterAndNamespaceIndexFunc},
		f.tweakListOptions,
	)
}

func (f *controllerRevisionClusterInformer) Informer() kcpcache.ScopeableSharedIndexInformer {
	return f.factory.InformerFor(&appsv1.ControllerRevision{}, f.defaultInformer)
}

func (f *controllerRevisionClusterInformer) Lister() appsv1listers.ControllerRevisionClusterLister {
	return appsv1listers.NewControllerRevisionClusterLister(f.Informer().GetIndexer())
}

func (f *controllerRevisionClusterInformer) Cluster(cluster logicalcluster.Name) upstreamappsv1informers.ControllerRevisionInformer {
	return &controllerRevisionInformer{
		informer: f.Informer().Cluster(cluster),
		lister:   f.Lister().Cluster(cluster),
	}
}

type controllerRevisionInformer struct {
	informer cache.SharedIndexInformer
	lister   upstreamappsv1listers.ControllerRevisionLister
}

func (f *controllerRevisionInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

func (f *controllerRevisionInformer) Lister() upstreamappsv1listers.ControllerRevisionLister {
	return f.lister
}
