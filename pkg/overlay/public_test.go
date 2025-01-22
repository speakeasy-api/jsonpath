package overlay_test

import (
	"github.com/speakeasy-api/jsonpath/pkg/overlay"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApplyOverlay_JSON(t *testing.T) {
	originalJSON := `{
  "openapi": "3.0.0",
  "info": {
    "title": "Example API",
    "version": "1.0.0",
    "description": "This is an example API"
  },
  "paths": {}
}`
	overlayYAML := `overlay: 1.0.0
info:
  title: Example Overlay
  version: 0.0.1
actions:
- target: $.info.description
  update: Updated description
`
	expectedResult := `{
  "openapi": "3.0.0",
  "info": {
    "title": "Example API",
    "version": "1.0.0",
    "description": "Updated description"
  },
  "paths": {}
}`

	result, err := overlay.ApplyOverlay(originalJSON, overlayYAML)
	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestApplyOverlay_YAML(t *testing.T) {
	originalYAML := `openapi: 3.0.0
info:
  title: Example API
  version: 1.0.0
  description: This is an example API
paths: {}
`
	overlayYAML := `overlay: 1.0.0
info:
  title: Example Overlay
  version: 0.0.1
actions:
- target: $.info.description
  update: Updated description
`
	expectedResult := `openapi: 3.0.0
info:
  title: Example API
  version: 1.0.0
  description: Updated description
paths: {}
`

	result, err := overlay.ApplyOverlay(originalYAML, overlayYAML)
	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}
