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

package fuzzer

import (
	"reflect"
	"strings"

	"github.com/google/gofuzz"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
)

// Funcs returns the fuzzer functions for the apiextensions apis.
func Funcs(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(obj *apiextensions.CustomResourceDefinitionSpec, c fuzz.Continue) {
			c.FuzzNoCustom(obj)

			// match our defaulter
			if len(obj.Scope) == 0 {
				obj.Scope = apiextensions.NamespaceScoped
			}
			if len(obj.Names.Singular) == 0 {
				obj.Names.Singular = strings.ToLower(obj.Names.Kind)
			}
			if len(obj.Names.ListKind) == 0 && len(obj.Names.Kind) > 0 {
				obj.Names.ListKind = obj.Names.Kind + "List"
			}
		},
		func(obj *apiextensions.JSONSchemaProps, c fuzz.Continue) {
			// we cannot use c.FuzzNoCustom because of the interface{} fields. So let's loop with reflection.
			vobj := reflect.ValueOf(obj).Elem()
			tobj := reflect.TypeOf(obj).Elem()
			for i := 0; i < tobj.NumField(); i++ {
				field := tobj.Field(i)
				switch field.Name {
				case "Default", "Enum", "Example", "Items", "AdditionalProperties", "AdditionalItems", "Schema":
					continue
				default:
					isValue := true
					switch field.Type.Kind() {
					case reflect.Interface, reflect.Map, reflect.Slice, reflect.Ptr:
						isValue = false
					}
					if isValue || c.Intn(5) == 0 {
						c.Fuzz(vobj.Field(i).Addr().Interface())
					}
				}
			}
			if c.RandBool() {
				obj.Default = `{"some": {"json": "test"}, "string": 42}` // some valid json
			}
			if c.RandBool() {
				obj.Enum = []interface{}{c.Uint64(), c.RandString(), c.RandBool()}
			}
			if c.RandBool() {
				obj.Example = "foobarbaz"
			}
			if c.RandBool() {
				c.Fuzz(obj.Items)
			}
			if c.RandBool() {
				c.Fuzz(obj.AdditionalProperties)
			}
			if c.RandBool() {
				c.Fuzz(obj.AdditionalItems)
			}
			if c.RandBool() {
				obj.Schema = apiextensions.JSONSchemaURL("example.com") // some valid url
			}
		},
	}
}
