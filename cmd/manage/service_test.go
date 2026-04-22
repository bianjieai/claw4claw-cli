package manage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertYAMLToJSON(t *testing.T) {
	// Test basic string
	assert.Equal(t, "hello", convertYAMLToJSON("hello"))

	// Test map[interface{}]interface{}
	inputMap := map[interface{}]interface{}{
		"key1": "value1",
		"key2": map[interface{}]interface{}{
			"nestedKey": "nestedValue",
		},
		"key3": []interface{}{
			map[interface{}]interface{}{
				"arrKey": "arrValue",
			},
		},
	}

	expectedMap := map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{
			"nestedKey": "nestedValue",
		},
		"key3": []interface{}{
			map[string]interface{}{
				"arrKey": "arrValue",
			},
		},
	}

	result := convertYAMLToJSON(inputMap)
	assert.Equal(t, expectedMap, result)
}

func TestLoadServiceFromFile(t *testing.T) {
	tmpDir := t.TempDir()

	yamlContent := `
title: "Test YAML Service"
description: "YAML Description"
category: "programming"
price: 10.5
avgResponseMs: 1000
inputSchema:
  type: object
  properties:
    query: { type: string }
  allOf:
    - type: object
      properties:
        id: { type: integer }
outputSchema:
  type: object
  properties:
    result: { type: string }
`
	yamlFile := filepath.Join(tmpDir, "service.yaml")
	require.NoError(t, os.WriteFile(yamlFile, []byte(yamlContent), 0644))

	jsonContent := `{
		"title": "Test JSON Service",
		"description": "JSON Description",
		"category": "design",
		"price": 20.0,
		"avgResponseMs": 2000,
		"inputSchema": {
			"type": "object",
			"properties": {
				"query": { "type": "string" }
			},
			"allOf": [
				{
					"type": "object",
					"properties": {
						"id": { "type": "integer" }
					}
				}
			]
		},
		"outputSchema": {
			"type": "object",
			"properties": {
				"result": { "type": "string" }
			}
		}
	}`
	jsonFile := filepath.Join(tmpDir, "service.json")
	require.NoError(t, os.WriteFile(jsonFile, []byte(jsonContent), 0644))

	t.Run("Load YAML", func(t *testing.T) {
		req, err := loadServiceFromFile(yamlFile)
		require.NoError(t, err)

		assert.Equal(t, "Test YAML Service", req.Title)
		assert.Equal(t, "programming", req.Category)
		assert.Equal(t, 10.5, req.Price)

		require.NotNil(t, req.InputSchema)
		require.NotNil(t, req.OutputSchema)

		// Verify nested maps are correctly typed as map[string]interface{}
		inputType := req.InputSchema["type"].(string)
		assert.Equal(t, "object", inputType)

		inputProps := req.InputSchema["properties"].(map[string]interface{})
		queryProp := inputProps["query"].(map[string]interface{})
		assert.Equal(t, "string", queryProp["type"])

		// Verify array and its nested maps
		allOf := req.InputSchema["allOf"].([]interface{})
		assert.Len(t, allOf, 1)
		allOfFirst := allOf[0].(map[string]interface{})
		assert.Equal(t, "object", allOfFirst["type"])
	})

	t.Run("Load JSON", func(t *testing.T) {
		req, err := loadServiceFromFile(jsonFile)
		require.NoError(t, err)

		assert.Equal(t, "Test JSON Service", req.Title)
		assert.Equal(t, "design", req.Category)
		assert.Equal(t, 20.0, req.Price)

		require.NotNil(t, req.InputSchema)
		require.NotNil(t, req.OutputSchema)

		// Verify nested maps are correctly typed as map[string]interface{}
		inputType := req.InputSchema["type"].(string)
		assert.Equal(t, "object", inputType)

		inputProps := req.InputSchema["properties"].(map[string]interface{})
		queryProp := inputProps["query"].(map[string]interface{})
		assert.Equal(t, "string", queryProp["type"])

		// Verify array and its nested maps
		allOf := req.InputSchema["allOf"].([]interface{})
		assert.Len(t, allOf, 1)
		allOfFirst := allOf[0].(map[string]interface{})
		assert.Equal(t, "object", allOfFirst["type"])
	})
}