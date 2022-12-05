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
	"github.com/kcp-dev/logicalcluster/v3"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	upstreamcorev1informers "k8s.io/client-go/informers/core/v1"
	upstreamcorev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/kcp-dev/client-go/informers/internalinterfaces"
	clientset "github.com/kcp-dev/client-go/kubernetes"
	corev1listers "github.com/kcp-dev/client-go/listers/core/v1"
)

// PersistentVolumeClaimClusterInformer provides access to a shared informer and lister for
// PersistentVolumeClaims.
type PersistentVolumeClaimClusterInformer interface {
	Cluster(logicalcluster.Name) upstreamcorev1informers.PersistentVolumeClaimInformer
	Informer() kcpcache.ScopeableSharedIndexInformer
	Lister() corev1listers.PersistentVolumeClaimClusterLister
}

type persistentVolumeClaimClusterInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewPersistentVolumeClaimClusterInformer constructs a new informer for PersistentVolumeClaim type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewPersistentVolumeClaimClusterInformer(client clientset.ClusterInterface, resyncPeriod time.Duration, indexers cache.Indexers) kcpcache.ScopeableSharedIndexInformer {
	return NewFilteredPersistentVolumeClaimClusterInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredPersistentVolumeClaimClusterInformer constructs a new informer for PersistentVolumeClaim type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredPersistentVolumeClaimClusterInformer(client clientset.ClusterInterface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) kcpcache.ScopeableSharedIndexInformer {
	return kcpinformers.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1().PersistentVolumeClaims().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1().PersistentVolumeClaims().Watch(context.TODO(), options)
			},
		},
		&corev1.PersistentVolumeClaim{},
		resyncPeriod,
		indexers,
	)
}

func (f *persistentVolumeClaimClusterInformer) defaultInformer(client clientset.ClusterInterface, resyncPeriod time.Duration) kcpcache.ScopeableSharedIndexInformer {
	return NewFilteredPersistentVolumeClaimClusterInformer(client, resyncPeriod, cache.Indexers{
		kcpcache.ClusterIndexName:             kcpcache.ClusterIndexFunc,
		kcpcache.ClusterAndNamespaceIndexName: kcpcache.ClusterAndNamespaceIndexFunc},
		f.tweakListOptions,
	)
}

func (f *persistentVolumeClaimClusterInformer) Informer() kcpcache.ScopeableSharedIndexInformer {
	return f.factory.InformerFor(&corev1.PersistentVolumeClaim{}, f.defaultInformer)
}

func (f *persistentVolumeClaimClusterInformer) Lister() corev1listers.PersistentVolumeClaimClusterLister {
	return corev1listers.NewPersistentVolumeClaimClusterLister(f.Informer().GetIndexer())
}

func (f *persistentVolumeClaimClusterInformer) Cluster(cluster logicalcluster.Name) upstreamcorev1informers.PersistentVolumeClaimInformer {
	return &persistentVolumeClaimInformer{
		informer: f.Informer().Cluster(cluster),
		lister:   f.Lister().Cluster(cluster),
	}
}

type persistentVolumeClaimInformer struct {
	informer cache.SharedIndexInformer
	lister   upstreamcorev1listers.PersistentVolumeClaimLister
}

func (f *persistentVolumeClaimInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

func (f *persistentVolumeClaimInformer) Lister() upstreamcorev1listers.PersistentVolumeClaimLister {
	return f.lister
}
