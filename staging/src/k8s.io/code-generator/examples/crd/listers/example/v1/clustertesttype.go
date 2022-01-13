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

package v1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	v1 "k8s.io/code-generator/examples/crd/apis/example/v1"
)

// ClusterTestTypeLister helps list ClusterTestTypes.
// All objects returned here must be treated as read-only.
type ClusterTestTypeLister interface {
	Scoped(scope rest.Scope) ClusterTestTypeLister
	// List lists all ClusterTestTypes in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.ClusterTestType, err error)
	// Get retrieves the ClusterTestType from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.ClusterTestType, error)
	ClusterTestTypeListerExpansion
}

// clusterTestTypeLister implements the ClusterTestTypeLister interface.
type clusterTestTypeLister struct {
	indexer cache.Indexer
	scope   rest.Scope
}

// NewClusterTestTypeLister returns a new ClusterTestTypeLister.
func NewClusterTestTypeLister(indexer cache.Indexer) ClusterTestTypeLister {
	return &clusterTestTypeLister{indexer: indexer}
}

func (s *clusterTestTypeLister) Scoped(scope rest.Scope) ClusterTestTypeLister {
	return &clusterTestTypeLister{
		indexer: s.indexer,
		scope:   scope,
	}
}

// List lists all ClusterTestTypes in the indexer.
func (s *clusterTestTypeLister) List(selector labels.Selector) (ret []*v1.ClusterTestType, err error) {
	appendFunc := func(m interface{}) {
		ret = append(ret, m.(*v1.ClusterTestType))
	}

	if s.scope == nil {
		err = cache.ListAll(s.indexer, selector, appendFunc)
		return ret, err
	}

	indexValue := s.scope.Name()
	err = cache.ListAllByIndexAndValue(s.indexer, cache.ListAllIndex, indexValue, selector, appendFunc)
	return ret, err
}

// Get retrieves the ClusterTestType from the index for a given name.
func (s *clusterTestTypeLister) Get(name string) (*v1.ClusterTestType, error) {
	key := name
	if s.scope != nil {
		key = s.scope.CacheKey(key)
	}
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("clustertesttype"), name)
	}
	return obj.(*v1.ClusterTestType), nil
}
