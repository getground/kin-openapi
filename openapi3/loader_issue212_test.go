package openapi3

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIssue212(t *testing.T) {
	spec := `
openapi: 3.0.1
info:
  title: 'test'
  version: 1.0.0
servers:
  - url: /api

paths:
  /available-products:
    get:
      operationId: getAvailableProductCollection
      responses:
        "200":
          description: test
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/AvailableProduct"

components:
  schemas:
    AvailableProduct:
      type: object
      properties:
        id:
          type: string
        type:
          type: string
        name:
          type: string
        media:
          type: object
          properties:
            documents:
              type: array
              items:
                allOf:
                  - $ref: "#/components/schemas/AvailableProduct/properties/previewImage/allOf/0"
                  - type: object
                    properties:
                      uri:
                        type: string
                        pattern: ^\/documents\/[0-9a-f]{64}$
        previewImage:
          allOf:
            - type: object
              required:
                - id
                - uri
              properties:
                id:
                  type: string
                uri:
                  type: string
            - type: object
              properties:
                uri:
                  type: string
                  pattern: ^\/images\/[0-9a-f]{64}$
`

	loader := NewLoader()
	doc, err := loader.LoadFromData([]byte(spec))
	require.NoError(t, err)
	err = doc.Validate(loader.Context)
	require.NoError(t, err)

	schema := NewSchemas()
	schema.Set("id", &SchemaRef{Value: &Schema{Type: "string"}})
	schema.Set("uri", &SchemaRef{Value: &Schema{Type: "string"}})
	expected, err := json.Marshal(&Schema{
		Type:       "object",
		Required:   []string{"id", "uri"},
		Properties: schema,
	},
	)
	require.NoError(t, err)
	got, err := json.Marshal(doc.Components.Schemas.Value("AvailableProduct").Value.Properties.Value("media").Value.Properties.Value("documents").Value.Items.Value.AllOf[0].Value)
	require.NoError(t, err)

	require.Equal(t, expected, got)
}
