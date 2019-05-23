/*
Copyright 2014 The Kubernetes Authors.

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

package admission

import (
	"fmt"
	"strings"
	"sync"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apiserver/pkg/authentication/user"
)

type attributesRecord struct {
	kind        schema.GroupVersionKind
	namespace   string
	name        string
	resource    schema.GroupVersionResource
	subresource string
	operation   Operation
	options     runtime.Object
	dryRun      bool
	object      runtime.Object
	oldObject   runtime.Object
	userInfo    user.Info

	// other elements are always accessed in single goroutine.
	// But ValidatingAdmissionWebhook add annotations concurrently.
	annotations     map[string]string
	annotationsLock sync.RWMutex

	reinvocationContext ReinvocationContext
}

func NewAttributesRecord(object runtime.Object, oldObject runtime.Object, kind schema.GroupVersionKind, namespace, name string, resource schema.GroupVersionResource, subresource string, operation Operation, operationOptions runtime.Object, dryRun bool, userInfo user.Info) Attributes {
	return &attributesRecord{
		kind:                kind,
		namespace:           namespace,
		name:                name,
		resource:            resource,
		subresource:         subresource,
		operation:           operation,
		options:             operationOptions,
		dryRun:              dryRun,
		object:              object,
		oldObject:           oldObject,
		userInfo:            userInfo,
		reinvocationContext: newReinvocationContext(),
	}
}

func (record *attributesRecord) GetKind() schema.GroupVersionKind {
	return record.kind
}

func (record *attributesRecord) GetNamespace() string {
	return record.namespace
}

func (record *attributesRecord) GetName() string {
	return record.name
}

func (record *attributesRecord) GetResource() schema.GroupVersionResource {
	return record.resource
}

func (record *attributesRecord) GetSubresource() string {
	return record.subresource
}

func (record *attributesRecord) GetOperation() Operation {
	return record.operation
}

func (record *attributesRecord) GetOperationOptions() runtime.Object {
	return record.options
}

func (record *attributesRecord) IsDryRun() bool {
	return record.dryRun
}

func (record *attributesRecord) GetObject() runtime.Object {
	return record.object
}

func (record *attributesRecord) GetOldObject() runtime.Object {
	return record.oldObject
}

func (record *attributesRecord) GetUserInfo() user.Info {
	return record.userInfo
}

// getAnnotations implements privateAnnotationsGetter.It's a private method used
// by WithAudit decorator.
func (record *attributesRecord) getAnnotations() map[string]string {
	record.annotationsLock.RLock()
	defer record.annotationsLock.RUnlock()

	if record.annotations == nil {
		return nil
	}
	cp := make(map[string]string, len(record.annotations))
	for key, value := range record.annotations {
		cp[key] = value
	}
	return cp
}

func (record *attributesRecord) AddAnnotation(key, value string) error {
	if err := checkKeyFormat(key); err != nil {
		return err
	}

	record.annotationsLock.Lock()
	defer record.annotationsLock.Unlock()

	if record.annotations == nil {
		record.annotations = make(map[string]string)
	}
	if v, ok := record.annotations[key]; ok && v != value {
		return fmt.Errorf("admission annotations are not allowd to be overwritten, key:%q, old value: %q, new value:%q", key, record.annotations[key], value)
	}
	record.annotations[key] = value
	return nil
}

func (record *attributesRecord) GetReinvocationContext() ReinvocationContext {
	return record.reinvocationContext
}

func newReinvocationContext() *reinvocationContext {
	return &reinvocationContext{previouslyInvokedReinvocableWebhooks: sets.NewString(), reinvokeWebhooks: sets.NewString()}
}

type reinvocationContext struct {
	// isReinvoke is true when admission plugins are being reinvoked
	isReinvoke bool
	// lastWebhookOutput holds the result of the last webhook admission plugin call
	lastWebhookOutput runtime.Object
	// previouslyInvokedReinvocableWebhooks holds the set of webhooks that have been invoked and
	// should be reinvoked if a later mutation occurs
	previouslyInvokedReinvocableWebhooks sets.String
	// reinvokeWebhooks holds the set of webhooks that should be reinvoked
	reinvokeWebhooks sets.String
	// reinvokeInTree indicates in-tree plugins should be reinvoked
	reinvokeInTree bool
}

func (rc *reinvocationContext) IsReinvoke() bool {
	return rc.isReinvoke
}

func (rc *reinvocationContext) SetIsReinvoke() {
	rc.isReinvoke = true
}

func (rc *reinvocationContext) ShouldReinvoke() bool {
	return rc.reinvokeInTree || len(rc.reinvokeWebhooks) > 0
}

func (rc *reinvocationContext) IsOutputChangedSinceLastWebhookInvocation(object runtime.Object) bool {
	return !apiequality.Semantic.DeepEqual(rc.lastWebhookOutput, object)
}

func (rc *reinvocationContext) SetLastWebhookInvocationOutput(object runtime.Object) {
	if object == nil {
		rc.lastWebhookOutput = nil
		return
	}
	rc.lastWebhookOutput = object.DeepCopyObject()
}

func (rc *reinvocationContext) ShouldInvokeWebhook(webhook string) bool {
	return !rc.isReinvoke || rc.reinvokeWebhooks.Has(webhook)
}

func (rc *reinvocationContext) AddReinvocableWebhookToPreviouslyInvoked(webhook string) {
	rc.previouslyInvokedReinvocableWebhooks.Insert(webhook)
}

func (rc *reinvocationContext) RequireReinvokingPreviouslyInvokedPlugins() {
	if len(rc.previouslyInvokedReinvocableWebhooks) > 0 {
		for s := range rc.previouslyInvokedReinvocableWebhooks {
			rc.reinvokeWebhooks.Insert(s)
		}
		rc.previouslyInvokedReinvocableWebhooks = sets.NewString()
	}
	rc.reinvokeInTree = true
}

func checkKeyFormat(key string) error {
	parts := strings.Split(key, "/")
	if len(parts) != 2 {
		return fmt.Errorf("annotation key has invalid format, the right format is a DNS subdomain prefix and '/' and key name. (e.g. 'podsecuritypolicy.admission.k8s.io/admit-policy')")
	}
	if msgs := validation.IsQualifiedName(key); len(msgs) != 0 {
		return fmt.Errorf("annotation key has invalid format %s. A qualified name like 'podsecuritypolicy.admission.k8s.io/admit-policy' is required.", strings.Join(msgs, ","))
	}
	return nil
}
