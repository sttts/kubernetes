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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	v1alpha1 "k8s.io/sample-apiserver/pkg/apis/wardle/v1alpha1"
)

// FlunderLister helps list Flunders.
// All objects returned here must be treated as read-only.
type FlunderLister interface {
	Scoped(scope rest.Scope) FlunderLister
	// List lists all Flunders in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Flunder, err error)
	// Flunders returns an object that can list and get Flunders.
	Flunders(namespace string) FlunderNamespaceLister
	FlunderListerExpansion
}

// flunderLister implements the FlunderLister interface.
type flunderLister struct {
	indexer cache.Indexer
	scope   rest.Scope
}

// NewFlunderLister returns a new FlunderLister.
func NewFlunderLister(indexer cache.Indexer) FlunderLister {
	return &flunderLister{indexer: indexer}
}

func (s *flunderLister) Scoped(scope rest.Scope) FlunderLister {
	return &flunderLister{
		indexer: s.indexer,
		scope:   scope,
	}
}

// List lists all Flunders in the indexer.
func (s *flunderLister) List(selector labels.Selector) (ret []*v1alpha1.Flunder, err error) {
	appendFunc := func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Flunder))
	}

	if s.scope == nil {
		err = cache.ListAll(s.indexer, selector, appendFunc)
		return ret, err
	}

	indexValue := s.scope.Name()
	err = cache.ListAllByIndexAndValue(s.indexer, cache.ListAllIndex, indexValue, selector, appendFunc)
	return ret, err
}

// Flunders returns an object that can list and get Flunders.
func (s *flunderLister) Flunders(namespace string) FlunderNamespaceLister {
	return flunderNamespaceLister{indexer: s.indexer, namespace: namespace, scope: s.scope}
}

// FlunderNamespaceLister helps list and get Flunders.
// All objects returned here must be treated as read-only.
type FlunderNamespaceLister interface {
	// List lists all Flunders in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Flunder, err error)
	// Get retrieves the Flunder from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Flunder, error)
	FlunderNamespaceListerExpansion
}

// flunderNamespaceLister implements the FlunderNamespaceLister
// interface.
type flunderNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
	scope     rest.Scope
}

// List lists all Flunders in the indexer for a given namespace.
func (s flunderNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Flunder, err error) {
	indexValue := s.namespace
	if s.scope != nil {
		indexValue = s.scope.CacheKey(s.namespace)
	}
	err = cache.ListAllByIndexAndValue(s.indexer, cache.NamespaceIndex, indexValue, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Flunder))
	})
	return ret, err
}

// Get retrieves the Flunder from the indexer for a given namespace and name.
func (s flunderNamespaceLister) Get(name string) (*v1alpha1.Flunder, error) {
	key := s.namespace + "/" + name
	if s.scope != nil {
		key = s.scope.CacheKey(key)
	}
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("flunder"), name)
	}
	return obj.(*v1alpha1.Flunder), nil
}
