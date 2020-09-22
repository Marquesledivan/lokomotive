package prometheus //nolint:testpackage

import (
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yamlserializer "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/util/jsonpath"

	"github.com/kinvolk/lokomotive/pkg/components/util"
)

func renderManifests(t *testing.T, configHCL string) map[string]string {
	name := "prometheus-operator"

	component := newComponent()

	body, diagnostics := util.GetComponentBody(configHCL, name)
	if diagnostics.HasErrors() {
		t.Fatalf("Getting component body: %v", diagnostics.Errs())
	}

	diagnostics = component.LoadConfig(body, &hcl.EvalContext{})
	if diagnostics.HasErrors() {
		t.Fatalf("Valid config should not return an error, got: %v", diagnostics)
	}

	ret, err := component.RenderManifests()
	if err != nil {
		t.Fatalf("Rendering manifests with valid config should succeed, got: %s", err)
	}

	return ret
}

//nolint:funlen
func TestConversion(t *testing.T) {
	testCases := []struct {
		Name                   string
		InputConfig            string
		ExpectedConfigFileName string
		Expected               reflect.Value
		JSONPath               string
	}{
		{
			Name: "use external_url param",
			InputConfig: `
component "prometheus-operator" {
  prometheus {
    external_url = "https://prometheus.externalurl.net"
  }
}
`,
			ExpectedConfigFileName: "prometheus-operator/templates/prometheus/prometheus.yaml",
			Expected:               reflect.ValueOf("https://prometheus.externalurl.net"),
			JSONPath:               "{.spec.externalUrl}",
		},
		{
			Name: "no external_url param",
			InputConfig: `
		component "prometheus-operator" {
		  prometheus {
		    ingress {
		      host                       = "prometheus.mydomain.net"
		      class                      = "contour"
		      certmanager_cluster_issuer = "letsencrypt-production"
		    }
		  }
		}
		`,
			ExpectedConfigFileName: "prometheus-operator/templates/prometheus/prometheus.yaml",
			Expected:               reflect.ValueOf("https://prometheus.mydomain.net"),
			JSONPath:               "{.spec.externalUrl}",
		},
		{
			Name: "ingress creation for prometheus",
			InputConfig: `
		component "prometheus-operator" {
		  prometheus {
		    ingress {
		      host                       = "prometheus.mydomain.net"
		      class                      = "contour"
		      certmanager_cluster_issuer = "letsencrypt-production"
		    }
		  }
		}
		`,
			ExpectedConfigFileName: "prometheus-operator/templates/prometheus/ingress.yaml",
			Expected:               reflect.ValueOf("prometheus.mydomain.net"),
			JSONPath:               "{.spec.rules[0].host}",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			m := renderManifests(t, tc.InputConfig)
			if len(m) == 0 {
				t.Fatalf("Rendered manifests shouldn't be empty")
			}

			gotConfig, ok := m[tc.ExpectedConfigFileName]
			if !ok {
				t.Fatalf("Config not found with filename: %q", tc.ExpectedConfigFileName)
			}

			u := getUnstructredObj(t, gotConfig)
			got := getValFromObject(t, tc.JSONPath, u)

			switch got.Kind() { //nolint:exhaustive
			case reflect.Interface:
				switch gotVal := got.Interface().(type) {
				// Add more cases as the expected values become heterogeneous.
				// case bool:
				case string:
					expVal, ok := tc.Expected.Interface().(string)
					if !ok {
						t.Fatalf("expected value is not string")
					}

					if gotVal != expVal {
						t.Fatalf("expected: %s, got: %s", expVal, gotVal)
					}

				default:
					t.Fatalf("Unknown type of the interface object: %T", got.Interface())
				}
			default:
				t.Fatalf("Unknown type of the object extracted from the converted YAML: %v", got.Kind())
			}
		})
	}
}

func getUnstructredObj(t *testing.T, yamlObj string) *unstructured.Unstructured {
	u := &unstructured.Unstructured{}

	// Decode YAML into `unstructured.Unstructured`.
	dec := yamlserializer.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	if _, _, err := dec.Decode([]byte(yamlObj), nil, u); err != nil {
		t.Fatalf("Converting config to unstructured.Unstructured: %v", err)
	}

	return u
}

func getValFromObject(t *testing.T, jp string, obj *unstructured.Unstructured) reflect.Value {
	jPath := jsonpath.New("parse")
	if err := jPath.Parse(jp); err != nil {
		t.Fatalf("Parsing JSONPath: %v", err)
	}

	v, err := jPath.FindResults(obj.Object)
	if err != nil {
		t.Fatalf("Finding results using JSONPath in the YAML file: %v", err)
	}

	if len(v) == 0 || len(v[0]) == 0 {
		t.Fatalf("No result found")
	}

	return v[0][0]
}
