/*
Copyright 2018 The Kubernetes Authors.

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

// Package mutating delegates admission checks to dynamically configured
// mutating webhooks.
package mutating

import (
	"context"
	"fmt"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/klog"

	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/api/admissionregistration/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/admission"
	admissionmetrics "k8s.io/apiserver/pkg/admission/metrics"
	webhookerrors "k8s.io/apiserver/pkg/admission/plugin/webhook/errors"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/generic"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/request"
	"k8s.io/apiserver/pkg/admission/plugin/webhook/util"
	"k8s.io/apiserver/pkg/util/webhook"
)

type mutatingDispatcher struct {
	cm     *webhook.ClientManager
	plugin *Plugin
}

func newMutatingDispatcher(p *Plugin) func(cm *webhook.ClientManager) generic.Dispatcher {
	return func(cm *webhook.ClientManager) generic.Dispatcher {
		return &mutatingDispatcher{cm, p}
	}
}

var _ generic.Dispatcher = &mutatingDispatcher{}

func (a *mutatingDispatcher) Dispatch(ctx context.Context, attr *generic.VersionedAttributes, o admission.ObjectInterfaces, relevantHooks []*v1beta1.Webhook) error {
	reinvokeCtx := attr.GetReinvocationContext()
	var webhookReinvokeCtx *webhookReinvokeContext
	if v := reinvokeCtx.Value(PluginName); v != nil {
		webhookReinvokeCtx = v.(*webhookReinvokeContext)
	} else {
		webhookReinvokeCtx = &webhookReinvokeContext{}
		reinvokeCtx.SetValue(PluginName, webhookReinvokeCtx)
	}

	reinvokeUpToHere := func() {
		webhookReinvokeCtx.RequireReinvokingPreviouslyInvokedPlugins()
		reinvokeCtx.SetShouldReinvoke()
	}

	if reinvokeCtx.IsReinvoke() && webhookReinvokeCtx.IsOutputChangedSinceLastWebhookInvocation(attr.Attributes.GetObject()) {
		// If the object has changed, we know the in-tree plugin re-invocations have mutated the object,
		// and we need to reinvoke all eligible webhooks.
		reinvokeUpToHere()
	}
	defer func() {
		webhookReinvokeCtx.SetLastWebhookInvocationOutput(attr.Attributes.GetObject())
	}()

	for _, hook := range relevantHooks {
		if reinvokeCtx.IsReinvoke() && !webhookReinvokeCtx.ShouldReinvokeWebhook(hook.Name) {
			continue
		}
		t := time.Now()
		changed, err := a.callAttrMutatingHook(ctx, hook, attr, o)
		admissionmetrics.Metrics.ObserveWebhook(time.Since(t), err != nil, attr.Attributes, "admit", hook.Name)

		if changed {
			// Patch had changed the object. Prepare to reinvoke all previous webhooks that are eligible for re-invocation.
			reinvokeUpToHere()
		}
		if hook.ReinvocationPolicy != nil && *hook.ReinvocationPolicy == v1beta1.IfNeededReinvocationPolicy {
			webhookReinvokeCtx.AddReinvocableWebhookToPreviouslyInvoked(hook.Name)
		}
		if err == nil {
			continue
		}

		ignoreClientCallFailures := hook.FailurePolicy != nil && *hook.FailurePolicy == v1beta1.Ignore
		if callErr, ok := err.(*webhook.ErrCallingWebhook); ok {
			if ignoreClientCallFailures {
				klog.Warningf("Failed calling webhook, failing open %v: %v", hook.Name, callErr)
				utilruntime.HandleError(callErr)
				continue
			}
			klog.Warningf("Failed calling webhook, failing closed %v: %v", hook.Name, err)
			return apierrors.NewInternalError(err)
		}
		return err
	}

	// convert attr.VersionedObject to the internal version in the underlying admission.Attributes
	if attr.VersionedObject != nil {
		return o.GetObjectConvertor().Convert(attr.VersionedObject, attr.Attributes.GetObject(), nil)
	}

	return nil
}

// note that callAttrMutatingHook updates attr
func (a *mutatingDispatcher) callAttrMutatingHook(ctx context.Context, h *v1beta1.Webhook, attr *generic.VersionedAttributes, o admission.ObjectInterfaces) (bool, error) {
	if attr.IsDryRun() {
		if h.SideEffects == nil {
			return false, &webhook.ErrCallingWebhook{WebhookName: h.Name, Reason: fmt.Errorf("Webhook SideEffects is nil")}
		}
		if !(*h.SideEffects == v1beta1.SideEffectClassNone || *h.SideEffects == v1beta1.SideEffectClassNoneOnDryRun) {
			return false, webhookerrors.NewDryRunUnsupportedErr(h.Name)
		}
	}

	// Currently dispatcher only supports `v1beta1` AdmissionReview
	// TODO: Make the dispatcher capable of sending multiple AdmissionReview versions
	if !util.HasAdmissionReviewVersion(v1beta1.SchemeGroupVersion.Version, h) {
		return false, &webhook.ErrCallingWebhook{WebhookName: h.Name, Reason: fmt.Errorf("webhook does not accept v1beta1 AdmissionReview")}
	}

	// Make the webhook request
	request := request.CreateAdmissionReview(attr)
	client, err := a.cm.HookClient(util.HookClientConfigForWebhook(h))
	if err != nil {
		return false, &webhook.ErrCallingWebhook{WebhookName: h.Name, Reason: err}
	}
	response := &admissionv1beta1.AdmissionReview{}
	r := client.Post().Context(ctx).Body(&request)
	if h.TimeoutSeconds != nil {
		r = r.Timeout(time.Duration(*h.TimeoutSeconds) * time.Second)
	}
	if err := r.Do().Into(response); err != nil {
		return false, &webhook.ErrCallingWebhook{WebhookName: h.Name, Reason: err}
	}

	if response.Response == nil {
		return false, &webhook.ErrCallingWebhook{WebhookName: h.Name, Reason: fmt.Errorf("Webhook response was absent")}
	}

	for k, v := range response.Response.AuditAnnotations {
		key := h.Name + "/" + k
		if err := attr.AddAnnotation(key, v); err != nil {
			klog.Warningf("Failed to set admission audit annotation %s to %s for mutating webhook %s: %v", key, v, h.Name, err)
		}
	}

	if !response.Response.Allowed {
		return false, webhookerrors.ToStatusErr(h.Name, response.Response.Result)
	}

	patchJS := response.Response.Patch
	if len(patchJS) == 0 {
		return false, nil
	}
	patchObj, err := jsonpatch.DecodePatch(patchJS)
	if err != nil {
		return false, apierrors.NewInternalError(err)
	}
	if len(patchObj) == 0 {
		return false, nil
	}

	// if a non-empty patch was provided, and we have no object we can apply it to (e.g. a DELETE admission operation), error
	if attr.VersionedObject == nil {
		return false, apierrors.NewInternalError(fmt.Errorf("admission webhook %q attempted to modify the object, which is not supported for this operation", h.Name))
	}

	jsonSerializer := json.NewSerializer(json.DefaultMetaFactory, o.GetObjectCreater(), o.GetObjectTyper(), false)
	objJS, err := runtime.Encode(jsonSerializer, attr.VersionedObject)
	if err != nil {
		return false, apierrors.NewInternalError(err)
	}
	patchedJS, err := patchObj.Apply(objJS)
	if err != nil {
		return false, apierrors.NewInternalError(err)
	}

	var newVersionedObject runtime.Object
	if _, ok := attr.VersionedObject.(*unstructured.Unstructured); ok {
		// Custom Resources don't have corresponding Go struct's.
		// They are represented as Unstructured.
		newVersionedObject = &unstructured.Unstructured{}
	} else {
		newVersionedObject, err = o.GetObjectCreater().New(attr.GetKind())
		if err != nil {
			return false, apierrors.NewInternalError(err)
		}
	}

	changed := !apiequality.Semantic.DeepEqual(attr.VersionedObject, newVersionedObject)

	// TODO: if we have multiple mutating webhooks, we can remember the json
	// instead of encoding and decoding for each one.
	if _, _, err := jsonSerializer.Decode(patchedJS, nil, newVersionedObject); err != nil {
		return false, apierrors.NewInternalError(err)
	}
	attr.VersionedObject = newVersionedObject
	o.GetObjectDefaulter().Default(attr.VersionedObject)
	return changed, nil
}
