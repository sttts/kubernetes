/*
Copyright 2014 The Kubernetes Authors All rights reserved.

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
package genericapiserver

import (
	"encoding/json"

	"github.com/golang/glog"
	"golang.org/x/net/context"

	"k8s.io/kubernetes/pkg/api/rest"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/storage"
)

type storageWithAuditing struct {
	storage.Interface

	resourcePrefix string
}

// AuditingDecorator decorates the given storageInterface with auditing output.
func AuditingDecorator(
	storageInterface storage.Interface,
	capacity int,
	objectType runtime.Object,
	resourcePrefix string,
	scopeStrategy rest.NamespaceScopedStrategy,
	newListFunc func() runtime.Object,
) storage.Interface {
	return &storageWithAuditing{
		Interface:      storageInterface,
		resourcePrefix: resourcePrefix,
	}
}

// GuaranteedUpdate implements storage.Interface.GuaranteedUpdate.
func (s *storageWithAuditing) GuaranteedUpdate(ctx context.Context, key string, ptrToType runtime.Object, ignoreNotFound bool, precondtions *storage.Preconditions, tryUpdate storage.UpdateFunc) error {
	var orig, modified []byte
	var err error
	auditedUpdate := func(input runtime.Object, res storage.ResponseMeta) (output runtime.Object, ttl *uint64, err error) {
		orig, err = json.Marshal(input)
		if err != nil {
			return
		}

		output, ttl, err = tryUpdate(input, res)
		if err != nil {
			return
		}

		modified, err = json.Marshal(output)

		return
	}

	err = s.Interface.GuaranteedUpdate(ctx, key, ptrToType, ignoreNotFound, precondtions, auditedUpdate)
	if err != nil {
		return err
	}

	glog.V(2).Infof("Updating %s: %s -> %s", key, string(orig), string(modified))
	return nil
}
