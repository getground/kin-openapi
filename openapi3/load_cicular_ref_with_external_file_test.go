//go:build go1.16
// +build go1.16

package openapi3_test

import (
	"embed"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/getkin/kin-openapi/openapi3"
)

//go:embed testdata/circularRef/*
var circularResSpecs embed.FS

func TestLoadCircularRefFromFile(t *testing.T) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, uri *url.URL) ([]byte, error) {
		return circularResSpecs.ReadFile(uri.Path)
	}

	got, err := loader.LoadFromFile("testdata/circularRef/base.yml")
	require.NoError(t, err)

	innerSchema := openapi3.NewSchemas()
	innerSchema.Set("id", &openapi3.SchemaRef{
		Value: &openapi3.Schema{Type: "string"},
	})

	schemas := openapi3.NewSchemas()
	schemas.Set("foo2", &openapi3.SchemaRef{
		Ref: "other.yml#/components/schemas/Foo2", // reference to an external file
		Value: &openapi3.Schema{
			Properties: innerSchema,
		},
	})
	foo := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Properties: schemas,
		},
	}
	bar := &openapi3.SchemaRef{Value: &openapi3.Schema{Properties: openapi3.NewSchemas()}}
	// circular reference
	bar.Value.Properties.Set("foo", &openapi3.SchemaRef{Ref: "#/components/schemas/Foo", Value: foo.Value})
	foo.Value.Properties.Set("bar", &openapi3.SchemaRef{Ref: "#/components/schemas/Bar", Value: bar.Value})

	wantSchema := openapi3.NewSchemas()
	wantSchema.Set("Foo", foo)
	wantSchema.Set("Bar", bar)

	want := &openapi3.T{
		OpenAPI: "3.0.3",
		Info: &openapi3.Info{
			Title:   "Recursive cyclic refs example",
			Version: "1.0",
		},
		Components: &openapi3.Components{
			Schemas: wantSchema,
		},
	}

	jsoner := func(doc *openapi3.T) string {
		data, err := json.Marshal(doc)
		require.NoError(t, err)
		return string(data)
	}
	require.JSONEq(t, jsoner(want), jsoner(got))
}
