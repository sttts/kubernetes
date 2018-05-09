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

package conversion

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// NoConversionConverter is a converter that does only set the apiVersion and kind, but without actual conversions.
type NoConversionConverter struct {
	ClusterScoped bool
}

// NewNoConversionConverter create a NoConversionConverter.
func NewNoConversionConverter(clusterScoped bool) NoConversionConverter {
	return NoConversionConverter{
		ClusterScoped: clusterScoped,
	}
}

func (NoConversionConverter) Convert(in, out, context interface{}) error {
	unstructIn, ok := in.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("input type %T in not valid for unstructured conversion", in)
	}

	unstructOut, ok := out.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("output type %T in not valid for unstructured conversion", out)
	}

	// maybe deep copy the map? It is documented in the
	// ObjectConverter interface that this function is not
	// guaranteed to not mutate the input. Or maybe set the input
	// object to nil.
	unstructOut.Object = unstructIn.Object
	return nil
}

func (c NoConversionConverter) ConvertToVersion(in runtime.Object, target runtime.GroupVersioner) (runtime.Object, error) {
	if kind := in.GetObjectKind().GroupVersionKind(); !kind.Empty() {
		gvk, ok := target.KindForGroupVersionKinds([]schema.GroupVersionKind{kind})
		if !ok {
			// TODO: should this be a typed error?
			return nil, fmt.Errorf("%v is unstructured and is not suitable for converting to %q", kind, target)
		}
		in.GetObjectKind().SetGroupVersionKind(gvk)

		if meta.IsListType(in) {
			return in, meta.EachListItem(in, func(item runtime.Object) error {
				item.GetObjectKind().SetGroupVersionKind(gvk)
				return nil
			})
		}
	}

	return in, nil
}

func (c NoConversionConverter) ConvertFieldLabel(version, kind, label, value string) (string, string, error) {
	// We currently only support metadata.namespace and metadata.name.
	switch {
	case label == "metadata.name":
		return label, value, nil
	case !c.ClusterScoped && label == "metadata.namespace":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("field label not supported: %s", label)
	}
}
