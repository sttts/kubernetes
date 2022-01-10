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
	v1beta1 "k8s.io/api/storage/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// VolumeAttachmentLister helps list VolumeAttachments.
// All objects returned here must be treated as read-only.
type VolumeAttachmentLister interface {
	Scoped(scope rest.Scope) VolumeAttachmentLister
	// List lists all VolumeAttachments in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.VolumeAttachment, err error)
	// Get retrieves the VolumeAttachment from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta1.VolumeAttachment, error)
	VolumeAttachmentListerExpansion
}

// volumeAttachmentLister implements the VolumeAttachmentLister interface.
type volumeAttachmentLister struct {
	indexer cache.Indexer
	scope   rest.Scope
}

// NewVolumeAttachmentLister returns a new VolumeAttachmentLister.
func NewVolumeAttachmentLister(indexer cache.Indexer) VolumeAttachmentLister {
	return &volumeAttachmentLister{indexer: indexer}
}

func (s *volumeAttachmentLister) Scoped(scope rest.Scope) VolumeAttachmentLister {
	return &volumeAttachmentLister{
		indexer: s.indexer,
		scope:   scope,
	}
}

// List lists all VolumeAttachments in the indexer.
func (s *volumeAttachmentLister) List(selector labels.Selector) (ret []*v1beta1.VolumeAttachment, err error) {
	var indexValue string
	if s.scope != nil {
		indexValue = s.scope.Name()
	}
	err = cache.ListAllByIndexAndValue(s.indexer, cache.ListAllIndex, indexValue, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.VolumeAttachment))
	})
	return ret, err
}

// Get retrieves the VolumeAttachment from the index for a given name.
func (s *volumeAttachmentLister) Get(name string) (*v1beta1.VolumeAttachment, error) {
	key := name
	if s.scope != nil {
		key = s.scope.CacheKey(key)
	}
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("volumeattachment"), name)
	}
	return obj.(*v1beta1.VolumeAttachment), nil
}
