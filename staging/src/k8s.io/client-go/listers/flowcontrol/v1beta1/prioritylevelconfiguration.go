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

package v1beta1

import (
	v1beta1 "k8s.io/api/flowcontrol/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// PriorityLevelConfigurationLister helps list PriorityLevelConfigurations.
// All objects returned here must be treated as read-only.
type PriorityLevelConfigurationLister interface {
	Scoped(scope rest.Scope) PriorityLevelConfigurationLister
	// List lists all PriorityLevelConfigurations in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.PriorityLevelConfiguration, err error)
	// Get retrieves the PriorityLevelConfiguration from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta1.PriorityLevelConfiguration, error)
	PriorityLevelConfigurationListerExpansion
}

// priorityLevelConfigurationLister implements the PriorityLevelConfigurationLister interface.
type priorityLevelConfigurationLister struct {
	indexer cache.Indexer
	scope   rest.Scope
}

// NewPriorityLevelConfigurationLister returns a new PriorityLevelConfigurationLister.
func NewPriorityLevelConfigurationLister(indexer cache.Indexer) PriorityLevelConfigurationLister {
	return &priorityLevelConfigurationLister{indexer: indexer}
}

func (s *priorityLevelConfigurationLister) Scoped(scope rest.Scope) PriorityLevelConfigurationLister {
	return &priorityLevelConfigurationLister{
		indexer: s.indexer,
		scope:   scope,
	}
}

// List lists all PriorityLevelConfigurations in the indexer.
func (s *priorityLevelConfigurationLister) List(selector labels.Selector) (ret []*v1beta1.PriorityLevelConfiguration, err error) {
	appendFunc := func(m interface{}) {
		ret = append(ret, m.(*v1beta1.PriorityLevelConfiguration))
	}

	if s.scope == nil {
		err = cache.ListAll(s.indexer, selector, appendFunc)
		return ret, err
	}

	indexValue := s.scope.Name()
	err = cache.ListAllByIndexAndValue(s.indexer, cache.ListAllIndex, indexValue, selector, appendFunc)
	return ret, err
}

// Get retrieves the PriorityLevelConfiguration from the index for a given name.
func (s *priorityLevelConfigurationLister) Get(name string) (*v1beta1.PriorityLevelConfiguration, error) {
	key := name
	if s.scope != nil {
		key = s.scope.CacheKey(key)
	}
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("prioritylevelconfiguration"), name)
	}
	return obj.(*v1beta1.PriorityLevelConfiguration), nil
}
