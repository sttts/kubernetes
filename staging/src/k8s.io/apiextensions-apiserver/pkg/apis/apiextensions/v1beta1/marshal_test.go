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

package v1beta1

import (
	"encoding/json"
	"reflect"
	"testing"
)

type JSONSchemaPropsOrBoolHolder struct {
	JSPoB          JSONSchemaPropsOrBool  `json:"val1"`
	JSPoBOmitEmpty *JSONSchemaPropsOrBool `json:"val2,omitempty"`
}

func TestJSONSchemaPropsOrBoolUnmarshalJSON(t *testing.T) {
	cases := []struct {
		input  string
		result JSONSchemaPropsOrBoolHolder
	}{
		{`{}`, JSONSchemaPropsOrBoolHolder{}},

		{`{"val1": {}}`, JSONSchemaPropsOrBoolHolder{JSPoB: JSONSchemaPropsOrBool{Schema: &JSONSchemaProps{}}}},
		{`{"val1": {"type":"string"}}`, JSONSchemaPropsOrBoolHolder{JSPoB: JSONSchemaPropsOrBool{Schema: &JSONSchemaProps{Type: "string"}}}},
		{`{"val1": false}`, JSONSchemaPropsOrBoolHolder{JSPoB: JSONSchemaPropsOrBool{}}},
		{`{"val1": true}`, JSONSchemaPropsOrBoolHolder{JSPoB: JSONSchemaPropsOrBool{Allows: true}}},

		{`{"val2": {}}`, JSONSchemaPropsOrBoolHolder{JSPoBOmitEmpty: &JSONSchemaPropsOrBool{Schema: &JSONSchemaProps{}}}},
		{`{"val2": {"type":"string"}}`, JSONSchemaPropsOrBoolHolder{JSPoBOmitEmpty: &JSONSchemaPropsOrBool{Schema: &JSONSchemaProps{Type: "string"}}}},
		{`{"val2": false}`, JSONSchemaPropsOrBoolHolder{JSPoBOmitEmpty: &JSONSchemaPropsOrBool{}}},
		{`{"val2": true}`, JSONSchemaPropsOrBoolHolder{JSPoBOmitEmpty: &JSONSchemaPropsOrBool{Allows: true}}},
	}

	for _, c := range cases {
		var result JSONSchemaPropsOrBoolHolder
		if err := json.Unmarshal([]byte(c.input), &result); err != nil {
			t.Errorf("Failed to unmarshal input '%v': %v", c.input, err)
		}
		if !reflect.DeepEqual(result, c.result) {
			t.Errorf("Failed to unmarshal input '%v': expected %+v, got %+v", c.input, c.result, result)
		}
	}
}

func TestStringArrayOrStringMarshalJSON(t *testing.T) {
	cases := []struct {
		input  JSONSchemaPropsOrBoolHolder
		result string
	}{
		{JSONSchemaPropsOrBoolHolder{}, `{"val1":false}`},

		{JSONSchemaPropsOrBoolHolder{JSPoB: JSONSchemaPropsOrBool{Schema: &JSONSchemaProps{}}}, `{"val1":{}}`},
		{JSONSchemaPropsOrBoolHolder{JSPoB: JSONSchemaPropsOrBool{Schema: &JSONSchemaProps{Type: "string"}}}, `{"val1":{"type":"string"}}`},
		{JSONSchemaPropsOrBoolHolder{JSPoB: JSONSchemaPropsOrBool{}}, `{"val1":false}`},
		{JSONSchemaPropsOrBoolHolder{JSPoB: JSONSchemaPropsOrBool{Allows: true}}, `{"val1":true}`},

		{JSONSchemaPropsOrBoolHolder{JSPoBOmitEmpty: &JSONSchemaPropsOrBool{Schema: &JSONSchemaProps{}}}, `{"val1":false,"val2":{}}`},
		{JSONSchemaPropsOrBoolHolder{JSPoBOmitEmpty: &JSONSchemaPropsOrBool{Schema: &JSONSchemaProps{Type: "string"}}}, `{"val1":false,"val2":{"type":"string"}}`},
		{JSONSchemaPropsOrBoolHolder{JSPoBOmitEmpty: &JSONSchemaPropsOrBool{}}, `{"val1":false,"val2":false}`},
		{JSONSchemaPropsOrBoolHolder{JSPoBOmitEmpty: &JSONSchemaPropsOrBool{Allows: true}}, `{"val1":false,"val2":true}`},
	}

	for _, c := range cases {
		result, err := json.Marshal(&c.input)
		if err != nil {
			t.Errorf("Unexpected error marshaling input '%v': %v", c.input, err)
		}
		if string(result) != c.result {
			t.Errorf("Failed to marshal input '%v': expected: %q, got %q", c.input, c.result, string(result))
		}
	}
}
