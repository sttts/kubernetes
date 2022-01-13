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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// PersistentVolumeLister helps list PersistentVolumes.
// All objects returned here must be treated as read-only.
type PersistentVolumeLister interface {
	Scoped(scope rest.Scope) PersistentVolumeLister
	// List lists all PersistentVolumes in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.PersistentVolume, err error)
	// Get retrieves the PersistentVolume from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.PersistentVolume, error)
	PersistentVolumeListerExpansion
}

// persistentVolumeLister implements the PersistentVolumeLister interface.
type persistentVolumeLister struct {
	indexer cache.Indexer
	scope   rest.Scope
}

// NewPersistentVolumeLister returns a new PersistentVolumeLister.
func NewPersistentVolumeLister(indexer cache.Indexer) PersistentVolumeLister {
	return &persistentVolumeLister{indexer: indexer}
}

func (s *persistentVolumeLister) Scoped(scope rest.Scope) PersistentVolumeLister {
	return &persistentVolumeLister{
		indexer: s.indexer,
		scope:   scope,
	}
}

// List lists all PersistentVolumes in the indexer.
func (s *persistentVolumeLister) List(selector labels.Selector) (ret []*v1.PersistentVolume, err error) {
	appendFunc := func(m interface{}) {
		ret = append(ret, m.(*v1.PersistentVolume))
	}

	if s.scope == nil {
		err = cache.ListAll(s.indexer, selector, appendFunc)
		return ret, err
	}

	indexValue := s.scope.Name()
	err = cache.ListAllByIndexAndValue(s.indexer, cache.ListAllIndex, indexValue, selector, appendFunc)
	return ret, err
}

// Get retrieves the PersistentVolume from the index for a given name.
func (s *persistentVolumeLister) Get(name string) (*v1.PersistentVolume, error) {
	key := name
	if s.scope != nil {
		key = s.scope.CacheKey(key)
	}
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("persistentvolume"), name)
	}
	return obj.(*v1.PersistentVolume), nil
}
