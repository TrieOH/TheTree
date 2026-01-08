package testing

import (
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Special matcher types
type Matcher interface {
	Match(t *testing.T, val *httpexpect.Value) interface{}
}

// AnyString - matches any non-empty string and returns it
type AnyString struct{}

func (AnyString) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	return val.String().NotEmpty().Raw()
}

// AnyNumber - matches any number and returns it
type AnyNumber struct{}

func (AnyNumber) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	return val.Number().Raw()
}

// AnyUUID - validates UUID format and returns it
type AnyUUID struct{}

func (AnyUUID) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	s := val.String().NotEmpty().Raw()
	require.Len(t, s, 36, "expected UUID format")
	_, err := uuid.Parse(s)
	require.NoError(t, err, "expected valid UUID format, got %q", s)
	return s
}

// NotEmpty - just validates not empty, returns raw value
type NotEmpty struct{}

func (NotEmpty) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	val.NotNull()
	return val.Raw()
}

// Store - captures a value into a variable for later comparison
type Store struct {
	Into *interface{}
}

func (s Store) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	*s.Into = val.Raw()
	return *s.Into
}

// SameAs - compares against a stored value
type SameAs struct {
	Ref interface{}
}

func (s SameAs) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	actual := val.Raw()
	require.Equal(t, s.Ref, actual)
	return actual
}

// ByKey - for arrays of objects, index by a key field
type ByKey struct {
	Key        string                 // the field to use as key (e.g. "key", "id")
	Spec       map[string]interface{} // key -> expected spec
	AllowExtra bool                   // if true, allows keys not in Spec
}

func (b ByKey) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	arr := val.Array()

	// Build map of actual items by key
	// We store the *httpexpect.Value (not Object) so we can pass it to Validate
	actual := make(map[string]*httpexpect.Value)
	for i := 0; i < int(arr.Length().Raw()); i++ {
		item := arr.Value(i) // Keep as Value
		obj := item.Object() // Convert to Object temporarily to read the key
		keyValue := obj.Value(b.Key).String().Raw()
		actual[keyValue] = item // Store the Value, not the Object
	}

	// Check all expected keys exist
	for key := range b.Spec {
		require.Contains(t, actual, key, "missing key %q in array", key)
	}

	// Check no unexpected keys
	if !b.AllowExtra {
		for key := range actual {
			require.Contains(t, b.Spec, key, "unexpected key %q in array", key)
		}
	}

	// Validate each item
	results := make(map[string]interface{})
	for key, spec := range b.Spec {
		itemVal := actual[key]
		results[key] = Validate(t, itemVal, spec)
	}

	return results
}

// Each - validates each array element against the same spec
type Each struct {
	Spec interface{}
}

func (e Each) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	arr := val.Array()
	length := int(arr.Length().Raw())
	results := make([]interface{}, length)

	for i := 0; i < length; i++ {
		item := arr.Value(i)
		result := Validate(t, item, e.Spec)
		results[i] = result
	}

	return results
}

// AtIndex - validates a specific array index
type AtIndex struct {
	Index int
	Spec  interface{}
}

func (a AtIndex) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	arr := val.Array()
	require.Greater(t, int(arr.Length().Raw()), a.Index, "array too short for index %d", a.Index)

	item := arr.Value(a.Index)
	return Validate(t, item, a.Spec)
}

// InOrder - validates array elements in exact positional order
type InOrder struct {
	Specs []interface{}
}

func (io InOrder) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	arr := val.Array()
	arr.Length().IsEqual(len(io.Specs))

	results := make([]interface{}, len(io.Specs))
	for i, spec := range io.Specs {
		item := arr.Value(i)
		results[i] = Validate(t, item, spec)
	}

	return results
}

// HasAtPosition - for arrays of objects, validates that a specific key appears at a position
type HasAtPosition struct {
	Key      string      // field to identify by (e.g. "key")
	Value    string      // the value to find (e.g. "email")
	Position int         // expected position
	Spec     interface{} // spec to validate against
}

func (h HasAtPosition) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	arr := val.Array()

	// Find the item with the matching key
	var foundAt int = -1
	var item *httpexpect.Value

	for i := 0; i < int(arr.Length().Raw()); i++ {
		candidate := arr.Value(i)
		obj := candidate.Object()
		keyVal := obj.Value(h.Key).String().Raw()

		if keyVal == h.Value {
			foundAt = i
			item = candidate
			break
		}
	}

	require.NotEqual(t, -1, foundAt, "item with %s=%s not found", h.Key, h.Value)
	require.Equal(t, h.Position, foundAt, "item with %s=%s found at wrong position", h.Key, h.Value)

	return Validate(t, item, h.Spec)
}

// Validate recursively validates a value against a spec
// Returns captured values for further assertions
func Validate(t *testing.T, val *httpexpect.Value, spec interface{}) interface{} {
	t.Helper()

	// Handle Matcher types
	if matcher, ok := spec.(Matcher); ok {
		return matcher.Match(t, val)
	}

	// Handle maps (objects)
	if specMap, ok := spec.(map[string]interface{}); ok {
		obj := val.Object()
		results := make(map[string]interface{})

		for key, expectedVal := range specMap {
			fieldVal := obj.Value(key)
			results[key] = Validate(t, fieldVal, expectedVal)
		}

		return results
	}

	// Handle slices (arrays with positional validation)
	if specSlice, ok := spec.([]interface{}); ok {
		arr := val.Array()
		arr.Length().IsEqual(len(specSlice))
		results := make([]interface{}, len(specSlice))

		for i, expectedVal := range specSlice {
			item := arr.Value(i)
			results[i] = Validate(t, item, expectedVal)
		}

		return results
	}

	// Handle primitive values - direct equality with numeric tolerance
	actual := val.Raw()

	// Special handling for numeric comparisons (JSON decodes all numbers as float64)
	switch expected := spec.(type) {
	case int:
		actualFloat, ok := actual.(float64)
		if ok {
			require.Equal(t, float64(expected), actualFloat, "numeric value mismatch")
			return actual
		}
	case int64:
		actualFloat, ok := actual.(float64)
		if ok {
			require.Equal(t, float64(expected), actualFloat, "numeric value mismatch")
			return actual
		}
	case float64:
		actualFloat, ok := actual.(float64)
		if ok {
			require.Equal(t, expected, actualFloat, "numeric value mismatch")
			return actual
		}
	case float32:
		actualFloat, ok := actual.(float64)
		if ok {
			require.Equal(t, float64(expected), actualFloat, "numeric value mismatch")
			return actual
		}
	case uint:
		actualFloat, ok := actual.(float64)
		if ok {
			require.Equal(t, float64(expected), actualFloat, "numeric value mismatch")
			return actual
		}
	case uint64:
		actualFloat, ok := actual.(float64)
		if ok {
			require.Equal(t, float64(expected), actualFloat, "numeric value mismatch")
			return actual
		}
	}

	// Default: direct equality
	require.Equal(t, spec, actual, "value mismatch")
	return actual
}

// Example usage with the actual test:
//
// t.Run("GetSchemaVerbose", func(t *testing.T) {
// 	authClient := suite.Client(t).Auth(user.auth)
// 	schema := authClient.GET("/projects/" + projectID + "/schemas/" + schemaID + "/verbose").
// 		Expect(http.StatusOK).
// 		Data()
//
// 	// Capture field IDs for cross-version stability checks
// 	var (
// 		matriculaV1ID, matriculaV2ID interface{}
// 		cursoV1ID, cursoV2ID         interface{}
// 	)
//
// 	spec := map[string]interface{}{
// 		"id":                  schemaID,
// 		"project_id":          projectID,
// 		"title":               "scti-register-flow",
// 		"flow_id":             "scti-register",
// 		"type":                "context",
// 		"status":              "published",
// 		"current_version_id":  schemaVersion2ID,
// 		"created_at":          NotEmpty{},
// 		"updated_at":          NotEmpty{},
// 		"versions": InOrder{
// 			Specs: []interface{}{
// 				// Version 2 (newest first in response)
// 				map[string]interface{}{
// 					"id":             AnyUUID{},
// 					"schema_id":      schemaID,
// 					"version_number": 2,
// 					"fields": ByKey{
// 						Key: "key",
// 						Spec: map[string]interface{}{
// 							"matricula": map[string]interface{}{
// 								"id":          Store{Into: &matriculaV2ID},
// 								"key":         "matricula",
// 								"type":        "string",
// 								"owner":       "user",
// 								"title":       "Numero da Matrícula",
// 								"description": "Sua matrícula da UENF como aparece no sistema acadêmico",
// 								"placeholder": "20223200045",
// 								"required":    true,
// 								"mutable":     true,
// 								"position":    0,
// 							},
// 							"curso": map[string]interface{}{
// 								"id":          Store{Into: &cursoV2ID},
// 								"key":         "curso",
// 								"type":        "string",
// 								"owner":       "user",
// 								"title":       "Curso de Matrícula",
// 								"description": "O curso que você está matrículado na UENF",
// 								"placeholder": "Ciência da Computação",
// 								"required":    true,
// 								"mutable":     true,
// 								"position":    1,
// 							},
// 							"periodo": map[string]interface{}{
// 								"id":          AnyUUID{},
// 								"key":         "periodo",
// 								"type":        "int",
// 								"owner":       "user",
// 								"title":       "Período Atual",
// 								"description": "O período da sua matéria mais avançada da grade",
// 								"required":    true,
// 								"mutable":     true,
// 								"position":    2,
// 							},
// 						},
// 					},
// 				},
// 				// Version 1
// 				map[string]interface{}{
// 					"id":             AnyUUID{},
// 					"schema_id":      schemaID,
// 					"version_number": 1,
// 					"fields": ByKey{
// 						Key: "key",
// 						Spec: map[string]interface{}{
// 							"matricula": map[string]interface{}{
// 								"id":          Store{Into: &matriculaV1ID},
// 								"key":         "matricula",
// 								"type":        "string",
// 								"owner":       "user",
// 								"title":       "Numero da Matrícula",
// 								"description": "Sua matrícula da UENF como aparece no sistema acadêmico",
// 								"placeholder": "20223200045",
// 								"required":    true,
// 								"mutable":     true,
// 								"position":    0,
// 							},
// 							"curso": map[string]interface{}{
// 								"id":          Store{Into: &cursoV1ID},
// 								"key":         "curso",
// 								"type":        "string",
// 								"owner":       "user",
// 								"title":       "Curso de Matrícula",
// 								"description": "O curso que você está matrículado na UENF",
// 								"placeholder": "Ciência da Computação",
// 								"required":    true,
// 								"mutable":     true,
// 								"position":    1,
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
//
// 	Validate(t, schema, spec)
//
// 	// Cross-version field ID stability checks
// 	require.Equal(t, matriculaV1ID, matriculaV2ID, "matricula field ID changed between versions")
// 	require.Equal(t, cursoV1ID, cursoV2ID, "curso field ID changed between versions")
// })
