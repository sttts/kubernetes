/*
Copyright 2017 The Kubernetes Authors.

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

package customresource

import (
	"fmt"

	"github.com/go-openapi/validate"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	apiserverstorage "k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	utilfeature "k8s.io/apiserver/pkg/util/feature"

	apiservervalidation "k8s.io/apiextensions-apiserver/pkg/apiserver/validation"
	apiextensionsfeatures "k8s.io/apiextensions-apiserver/pkg/features"
)

// customResourceStrategy implements behavior for CustomResources.
type customResourceStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	namespaceScoped bool
	validator       customResourceValidator
}

func NewStrategy(typer runtime.ObjectTyper, namespaceScoped bool, kind schema.GroupVersionKind, validator *validate.SchemaValidator) customResourceStrategy {
	return customResourceStrategy{
		ObjectTyper:     typer,
		NameGenerator:   names.SimpleNameGenerator,
		namespaceScoped: namespaceScoped,
		validator: customResourceValidator{
			namespaceScoped: namespaceScoped,
			kind:            kind,
			validator:       validator,
		},
	}
}

// DefaultGarbageCollectionPolicy returns Orphan because that was the default
// behavior before the server-side garbage collection was implemented.
func (customResourceStrategy) DefaultGarbageCollectionPolicy() rest.GarbageCollectionPolicy {
	return rest.OrphanDependents
}

func (a customResourceStrategy) NamespaceScoped() bool {
	return a.namespaceScoped
}

// PrepareForCreate clears the status of a CustomResource before creation.
func (customResourceStrategy) PrepareForCreate(ctx genericapirequest.Context, obj runtime.Object) {
	if utilfeature.DefaultFeatureGate.Enabled(apiextensionsfeatures.CustomResourceSubResources) {
		customResourceObject := obj.(*unstructured.Unstructured)
		customResource := customResourceObject.UnstructuredContent()

		// create cannot set status
		if _, ok := customResource["status"]; ok {
			customResource["status"] = nil
		}
	}

	accessor, _ := meta.Accessor(obj)
	accessor.SetGeneration(1)
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (customResourceStrategy) PrepareForUpdate(ctx genericapirequest.Context, obj, old runtime.Object) {
	if utilfeature.DefaultFeatureGate.Enabled(apiextensionsfeatures.CustomResourceSubResources) {
		newCustomResourceObject := obj.(*unstructured.Unstructured)
		oldCustomResourceObject := old.(*unstructured.Unstructured)

		newCustomResource := newCustomResourceObject.UnstructuredContent()
		oldCustomResource := oldCustomResourceObject.UnstructuredContent()

		// update is not allowed to set status
		if _, ok := newCustomResource["status"]; ok {
			newCustomResource["status"] = oldCustomResource["status"]
		}

		// Any changes to the spec increment the generation number, any changes to the
		// status should reflect the generation number of the corresponding object. We push
		// the burden of managing the status onto the clients because we can't (in general)
		// know here what version of spec the writer of the status has seen. It may seem like
		// we can at first -- since obj contains spec -- but in the future we will probably make
		// status its own object, and even if we don't, writes may be the result of a
		// read-update-write loop, so the contents of spec may not actually be the spec that
		// the CustomResource has *seen*.

		if _, ok := newCustomResource["spec"]; ok {
			oldSpec := oldCustomResource["spec"]
			newSpec := newCustomResource["spec"]

			if !apiequality.Semantic.DeepEqual(oldSpec, newSpec) {
				oldAccessor, _ := meta.Accessor(oldCustomResourceObject)
				newAccessor, _ := meta.Accessor(newCustomResourceObject)
				newAccessor.SetGeneration(oldAccessor.GetGeneration() + 1)
			}
		}
	}
}

// Validate validates a new CustomResource.
func (a customResourceStrategy) Validate(ctx genericapirequest.Context, obj runtime.Object) field.ErrorList {
	return a.validator.Validate(ctx, obj)
}

// Canonicalize normalizes the object after validation.
func (customResourceStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is false for CustomResources; this means a POST is
// needed to create one.
func (customResourceStrategy) AllowCreateOnUpdate() bool {
	return false
}

// AllowUnconditionalUpdate is the default update policy for CustomResource objects.
func (customResourceStrategy) AllowUnconditionalUpdate() bool {
	return false
}

// ValidateUpdate is the default update validation for an end user updating status.
func (a customResourceStrategy) ValidateUpdate(ctx genericapirequest.Context, obj, old runtime.Object) field.ErrorList {
	return a.validator.ValidateUpdate(ctx, obj, old)
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func (a customResourceStrategy) GetAttrs(obj runtime.Object) (labels.Set, fields.Set, bool, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, nil, false, err
	}
	return labels.Set(accessor.GetLabels()), objectMetaFieldsSet(accessor, a.namespaceScoped), accessor.GetInitializers() != nil, nil
}

// objectMetaFieldsSet returns a fields that represent the ObjectMeta.
func objectMetaFieldsSet(objectMeta metav1.Object, namespaceScoped bool) fields.Set {
	if namespaceScoped {
		return fields.Set{
			"metadata.name":      objectMeta.GetName(),
			"metadata.namespace": objectMeta.GetNamespace(),
		}
	}
	return fields.Set{
		"metadata.name": objectMeta.GetName(),
	}
}

// MatchCustomResourceDefinitionStorage is the filter used by the generic etcd backend to route
// watch events from etcd to clients of the apiserver only interested in specific
// labels/fields.
func (a customResourceStrategy) MatchCustomResourceDefinitionStorage(label labels.Selector, field fields.Selector) apiserverstorage.SelectionPredicate {
	return apiserverstorage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: a.GetAttrs,
	}
}

type customResourceValidator struct {
	namespaceScoped bool
	kind            schema.GroupVersionKind
	validator       *validate.SchemaValidator
}

func (a customResourceValidator) Validate(ctx genericapirequest.Context, obj runtime.Object) field.ErrorList {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return field.ErrorList{field.Invalid(field.NewPath("metadata"), nil, err.Error())}
	}
	typeAccessor, err := meta.TypeAccessor(obj)
	if err != nil {
		return field.ErrorList{field.Invalid(field.NewPath("kind"), nil, err.Error())}
	}
	if typeAccessor.GetKind() != a.kind.Kind {
		return field.ErrorList{field.Invalid(field.NewPath("kind"), typeAccessor.GetKind(), fmt.Sprintf("must be %v", a.kind.Kind))}
	}
	if typeAccessor.GetAPIVersion() != a.kind.Group+"/"+a.kind.Version {
		return field.ErrorList{field.Invalid(field.NewPath("apiVersion"), typeAccessor.GetAPIVersion(), fmt.Sprintf("must be %v", a.kind.Group+"/"+a.kind.Version))}
	}

	customResourceObject, ok := obj.(*unstructured.Unstructured)
	// this will never happen.
	if !ok {
		return field.ErrorList{field.Invalid(field.NewPath(""), customResourceObject, fmt.Sprintf("has type %T. Must be a pointer to an Unstructured type", customResourceObject))}
	}
	customResource := customResourceObject.UnstructuredContent()

	if err = apiservervalidation.ValidateCustomResource(customResource, a.validator); err != nil {
		return field.ErrorList{field.Invalid(field.NewPath(""), customResource, err.Error())}
	}

	return validation.ValidateObjectMetaAccessor(accessor, a.namespaceScoped, validation.NameIsDNSSubdomain, field.NewPath("metadata"))
}

func (a customResourceValidator) ValidateUpdate(ctx genericapirequest.Context, obj, old runtime.Object) field.ErrorList {
	objAccessor, err := meta.Accessor(obj)
	if err != nil {
		return field.ErrorList{field.Invalid(field.NewPath("metadata"), nil, err.Error())}
	}
	oldAccessor, err := meta.Accessor(old)
	if err != nil {
		return field.ErrorList{field.Invalid(field.NewPath("metadata"), nil, err.Error())}
	}
	typeAccessor, err := meta.TypeAccessor(obj)
	if err != nil {
		return field.ErrorList{field.Invalid(field.NewPath("kind"), nil, err.Error())}
	}
	if typeAccessor.GetKind() != a.kind.Kind {
		return field.ErrorList{field.Invalid(field.NewPath("kind"), typeAccessor.GetKind(), fmt.Sprintf("must be %v", a.kind.Kind))}
	}
	if typeAccessor.GetAPIVersion() != a.kind.Group+"/"+a.kind.Version {
		return field.ErrorList{field.Invalid(field.NewPath("apiVersion"), typeAccessor.GetAPIVersion(), fmt.Sprintf("must be %v", a.kind.Group+"/"+a.kind.Version))}
	}

	customResourceObject, ok := obj.(*unstructured.Unstructured)
	// this will never happen.
	if !ok {
		return field.ErrorList{field.Invalid(field.NewPath(""), customResourceObject, fmt.Sprintf("has type %T. Must be a pointer to an Unstructured type", customResourceObject))}
	}
	customResource := customResourceObject.UnstructuredContent()

	if err = apiservervalidation.ValidateCustomResource(customResource, a.validator); err != nil {
		return field.ErrorList{field.Invalid(field.NewPath(""), customResource, err.Error())}
	}

	return validation.ValidateObjectMetaAccessorUpdate(objAccessor, oldAccessor, field.NewPath("metadata"))
}

func (a customResourceValidator) ValidateStatusUpdate(ctx genericapirequest.Context, obj, old runtime.Object) field.ErrorList {
	objAccessor, err := meta.Accessor(obj)
	if err != nil {
		return field.ErrorList{field.Invalid(field.NewPath("metadata"), nil, err.Error())}
	}
	oldAccessor, err := meta.Accessor(old)
	if err != nil {
		return field.ErrorList{field.Invalid(field.NewPath("metadata"), nil, err.Error())}
	}
	typeAccessor, err := meta.TypeAccessor(obj)
	if err != nil {
		return field.ErrorList{field.Invalid(field.NewPath("kind"), nil, err.Error())}
	}
	if typeAccessor.GetKind() != a.kind.Kind {
		return field.ErrorList{field.Invalid(field.NewPath("kind"), typeAccessor.GetKind(), fmt.Sprintf("must be %v", a.kind.Kind))}
	}
	if typeAccessor.GetAPIVersion() != a.kind.Group+"/"+a.kind.Version {
		return field.ErrorList{field.Invalid(field.NewPath("apiVersion"), typeAccessor.GetAPIVersion(), fmt.Sprintf("must be %v", a.kind.Group+"/"+a.kind.Version))}
	}

	customResourceObject, ok := obj.(*unstructured.Unstructured)
	// this will never happen.
	if !ok {
		return field.ErrorList{field.Invalid(field.NewPath(""), customResourceObject, fmt.Sprintf("has type %T. Must be a pointer to an Unstructured type", customResourceObject))}
	}
	customResource := customResourceObject.UnstructuredContent()

	// validate only the status
	customResourceStatus := customResource["status"]
	if err = apiservervalidation.ValidateCustomResource(customResourceStatus, a.validator); err != nil {
		return field.ErrorList{field.Invalid(field.NewPath("status"), customResourceStatus, err.Error())}
	}

	return validation.ValidateObjectMetaAccessorUpdate(objAccessor, oldAccessor, field.NewPath("metadata"))
}

type customResourceDefinitionStorageStatusStrategy struct {
	customResourceStrategy
}

func NewStatusStrategy(strategy customResourceStrategy) customResourceDefinitionStorageStatusStrategy {
	return customResourceDefinitionStorageStatusStrategy{strategy}
}

func (customResourceDefinitionStorageStatusStrategy) PrepareForUpdate(ctx genericapirequest.Context, obj, old runtime.Object) {
	newCustomResourceObject := obj.(*unstructured.Unstructured)
	oldCustomResourceObject := old.(*unstructured.Unstructured)

	newCustomResource := newCustomResourceObject.UnstructuredContent()
	oldCustomResource := oldCustomResourceObject.UnstructuredContent()

	// update is not allowed to set spec and metadata
	newCustomResource["spec"] = oldCustomResource["spec"]
	newCustomResourceObject.SetAnnotations(oldCustomResourceObject.GetAnnotations())
	newCustomResourceObject.SetFinalizers(oldCustomResourceObject.GetFinalizers())
	newCustomResourceObject.SetGeneration(oldCustomResourceObject.GetGeneration())
	newCustomResourceObject.SetLabels(oldCustomResourceObject.GetLabels())
	newCustomResourceObject.SetOwnerReferences(oldCustomResourceObject.GetOwnerReferences())
	newCustomResourceObject.SetSelfLink(oldCustomResourceObject.GetSelfLink())

}

// ValidateUpdate is the default update validation for an end user updating status.
func (a customResourceDefinitionStorageStatusStrategy) ValidateUpdate(ctx genericapirequest.Context, obj, old runtime.Object) field.ErrorList {
	return a.customResourceStrategy.validator.ValidateStatusUpdate(ctx, obj, old)
}
