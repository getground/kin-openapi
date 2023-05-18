package openapi3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIssue344(t *testing.T) {
	sl := NewLoader()
	sl.IsExternalRefsAllowed = true

	doc, err := sl.LoadFromFile("testdata/spec.yaml")
	require.NoError(t, err)

	err = doc.Validate(sl.Context)
	require.NoError(t, err)

	require.Equal(t, "string", doc.Components.Schemas.Value("Test").Value.Properties.Value("test").Value.Properties.Value("name").Value.Type)
}
