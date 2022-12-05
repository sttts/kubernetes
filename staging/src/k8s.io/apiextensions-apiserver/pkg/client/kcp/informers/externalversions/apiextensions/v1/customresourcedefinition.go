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

package v1

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	kcpcache "github.com/kcp-dev/apimachinery/pkg/cache"
	kcpinformers "github.com/kcp-dev/apimachinery/third_party/informers"
	"github.com/kcp-dev/logicalcluster/v3"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1listers "k8s.io/apiextensions-apiserver/pkg/client/kcp/listers/apiextensions/v1"
	upstreamapiextensionsv1listers "k8s.io/apiextensions-apiserver/pkg/client/listers/apiextensions/v1"
	upstreamapiextensionsv1informers "k8s.io/apiextensions-apiserver/pkg/client/informers/externalversions/apiextensions/v1"
	clientset "k8s.io/apiextensions-apiserver/pkg/client/kcp/clientset/versioned"
	"k8s.io/apiextensions-apiserver/pkg/client/kcp/informers/externalversions/internalinterfaces"
)

// CustomResourceDefinitionClusterInformer provides access to a shared informer and lister for
// CustomResourceDefinitions.
type CustomResourceDefinitionClusterInformer interface {
	Cluster(logicalcluster.Name) upstreamapiextensionsv1informers.CustomResourceDefinitionInformer
	Informer() kcpcache.ScopeableSharedIndexInformer
	Lister() apiextensionsv1listers.CustomResourceDefinitionClusterLister
}

type customResourceDefinitionClusterInformer struct {
	factory internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewCustomResourceDefinitionClusterInformer constructs a new informer for CustomResourceDefinition type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewCustomResourceDefinitionClusterInformer(client clientset.ClusterInterface, resyncPeriod time.Duration, indexers cache.Indexers) kcpcache.ScopeableSharedIndexInformer {
	return NewFilteredCustomResourceDefinitionClusterInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredCustomResourceDefinitionClusterInformer constructs a new informer for CustomResourceDefinition type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCustomResourceDefinitionClusterInformer(client clientset.ClusterInterface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) kcpcache.ScopeableSharedIndexInformer {
	return kcpinformers.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApiextensionsV1().CustomResourceDefinitions().Watch(context.TODO(), options)
			},
		},
		&apiextensionsv1.CustomResourceDefinition{},
		resyncPeriod,
		indexers,
	)
}

func (f *customResourceDefinitionClusterInformer) defaultInformer(client clientset.ClusterInterface, resyncPeriod time.Duration) kcpcache.ScopeableSharedIndexInformer {
	return NewFilteredCustomResourceDefinitionClusterInformer(client, resyncPeriod, cache.Indexers{
			kcpcache.ClusterIndexName: kcpcache.ClusterIndexFunc,
			}, 
		f.tweakListOptions,
	)
}

func (f *customResourceDefinitionClusterInformer) Informer() kcpcache.ScopeableSharedIndexInformer {
	return f.factory.InformerFor(&apiextensionsv1.CustomResourceDefinition{}, f.defaultInformer)
}

func (f *customResourceDefinitionClusterInformer) Lister() apiextensionsv1listers.CustomResourceDefinitionClusterLister {
	return apiextensionsv1listers.NewCustomResourceDefinitionClusterLister(f.Informer().GetIndexer())
}

func (f *customResourceDefinitionClusterInformer) Cluster(cluster logicalcluster.Name) upstreamapiextensionsv1informers.CustomResourceDefinitionInformer {
	return &customResourceDefinitionInformer{
		informer: f.Informer().Cluster(cluster),
		lister:   f.Lister().Cluster(cluster),
	}
}

type customResourceDefinitionInformer struct {
	informer cache.SharedIndexInformer
	lister upstreamapiextensionsv1listers.CustomResourceDefinitionLister
}

func (f *customResourceDefinitionInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

func (f *customResourceDefinitionInformer) Lister() upstreamapiextensionsv1listers.CustomResourceDefinitionLister {
	return f.lister
}


