package openapi

import (
	"testing"

	"github.com/go-openapi/spec"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/apimachinery/pkg/util/json"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var crd = &apiextensions.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{Name: "plural.group.com"},
	Spec: apiextensions.CustomResourceDefinitionSpec{
		Group: "group.com",
		Scope: apiextensions.ResourceScope("Cluster"),
		Names: apiextensions.CustomResourceDefinitionNames{
			Plural:   "plural",
			Singular: "singular",
			Kind:     "Plural",
			ListKind: "PluralList",
		},
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    "v1",
				Served:  true,
				Storage: true,
			},
			{
				Name:    "v2",
				Served:  true,
				Storage: false,
			},
		},
	},
	Status: apiextensions.CustomResourceDefinitionStatus{
		AcceptedNames: apiextensions.CustomResourceDefinitionNames{
			Plural:   "plural",
			Singular: "singular",
			Kind:     "Plural",
			ListKind: "PluralList",
		},
		StoredVersions: []string{"v1"},
	},
}

func TestNewSpec(t *testing.T) {
	type args struct {
		crd    *apiextensions.CustomResourceDefinition
		v      string
		schema *spec.Schema
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"simple", args{crd, "v1", &spec.Schema{}}, "", false},
		{"schema", args{crd, "v1", &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Properties: map[string]spec.Schema{
					"foo": spec.Schema{SchemaProps: spec.SchemaProps{Type: spec.StringOrArray([]string{"string"})}},
				},
			},
		}}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := newSpec(tt.args.crd, tt.args.v, tt.args.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("newSpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			bs, _ := json.Marshal(spec)
			got := string(bs)
			if got != tt.want {
				t.Errorf("unexpected spec: %s", diff.StringDiff(got, tt.want))
			}
		})
	}
}
