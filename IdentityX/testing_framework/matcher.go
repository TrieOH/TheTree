package testing

import (
	"fmt"
	"testing"

	"github.com/gavv/httpexpect/v2"
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
	// Basic UUID validation (you can make this stricter)
	require.Len(t, s, 36, "expected UUID format")
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
	Key  string                 // the field to use as key (e.g. "key", "id")
	Spec map[string]interface{} // key -> expected spec
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

	// Check no unexpected keys (optional - remove if you want to allow extras)
	for key := range actual {
		require.Contains(t, b.Spec, key, "unexpected key %q in array", key)
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
	results := make([]interface{}, 0)

	for i := 0; i < int(arr.Length().Raw()); i++ {
		item := arr.Value(i)
		result := Validate(t, item, e.Spec)
		results = append(results, result)
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
	for i := 0; i < int(arr.Length().Raw()); i++ {
		candidate := arr.Value(i)
		obj := candidate.Object()
		keyVal := obj.Value(h.Key).String().Raw()

		if keyVal == h.Value {
			// Found it - check position and validate
			require.Equal(t, h.Position, i, "item with %s=%s found at wrong position (expected %d, got %d)",
				h.Key, h.Value, h.Position, i)
			return Validate(t, candidate, h.Spec)
		}
	}

	// Not found
	require.FailNow(t, fmt.Sprintf("item with %s=%s not found in array", h.Key, h.Value))
	return nil // unreachable, but makes linter happy
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

	// Handle primitive values - direct equality
	v := val.Raw()
	require.Equal(t, spec, v, "value mismatch")
	return v
}

// Example usage:
//
// // Capture field IDs for cross-version comparison
// var emailIDv1, emailIDv2 interface{}
//
// spec := map[string]interface{}{
// 	"id":         AnyUUID{},
// 	"project_id": projectID,
// 	"title":      "scti-register-flow",
// 	"flow_id":    "scti-register",
// 	"type":       "context",
// 	"status":     "published",
// 	"created_at": NotEmpty{},
// 	"updated_at": NotEmpty{},
// 	"versions": Each{
// 		Spec: map[string]interface{}{
// 			"id":            AnyUUID{},
// 			"version_number": AnyNumber{},
// 			"schema_id":      AnyUUID{},
// 			"fields": ByKey{
// 				Key: "key",
// 				Spec: map[string]interface{}{
// 					"matricula": map[string]interface{}{
// 						"id":          Store{Into: &emailIDv1}, // capture in v1
// 						"key":         "matricula",
// 						"type":        "string",
// 						"owner":       "user",
// 						"title":       "Numero da Matrícula",
// 						"description": "Sua matrícula da UENF como aparece no sistema acadêmico",
// 						"required":    true,
// 						"mutable":     true,
// 						"position":    0,
// 					},
// 					"curso": map[string]interface{}{
// 						"id":       AnyUUID{},
// 						"key":      "curso",
// 						"type":     "string",
// 						"required": true,
// 					},
// 				},
// 			},
// 		},
// 	},
// }
//
// Validate(t, schema, spec)
//
// // Later you can assert: emailIDv1 == emailIDv2
