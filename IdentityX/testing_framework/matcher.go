package testing

import (
	"testing"
	"time"

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

// AnyDate - validates RFC3339/ISO8601 date format and returns the string
type AnyDate struct{}

func (AnyDate) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	s := val.String().NotEmpty().Raw()
	_, err := time.Parse(time.RFC3339, s)
	if err != nil {
		// Try RFC3339Nano as fallback (includes nanoseconds)
		_, err = time.Parse(time.RFC3339Nano, s)
		require.NoError(t, err, "expected valid RFC3339 date format, got %q", s)
	}
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
// Optionally validates with a Matcher first
type Store struct {
	Into    *interface{}
	Matcher Matcher // optional: if provided, validates before storing
}

func (s Store) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	var result interface{}

	if s.Matcher != nil {
		result = s.Matcher.Match(t, val)
	} else {
		result = val.Raw()
	}

	*s.Into = result
	return result
}

// StoreString - stores directly into a *string variable
type StoreString struct {
	Into    *string
	Matcher Matcher // optional: if provided, validates before storing
}

func (s StoreString) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	var result interface{}

	if s.Matcher != nil {
		result = s.Matcher.Match(t, val)
	} else {
		result = val.Raw()
	}

	*s.Into = result.(string)
	return result
}

// StoreInt - stores directly into an *int variable
type StoreInt struct {
	Into    *int
	Matcher Matcher // optional: if provided, validates before storing
}

func (s StoreInt) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	var result interface{}

	if s.Matcher != nil {
		result = s.Matcher.Match(t, val)
	} else {
		result = val.Raw()
	}

	// Handle JSON numeric type (float64)
	switch v := result.(type) {
	case float64:
		*s.Into = int(v)
	case int:
		*s.Into = v
	default:
		require.FailNow(t, "expected numeric value for StoreInt")
	}

	return result
}

// StoreBool - stores directly into a *bool variable
type StoreBool struct {
	Into    *bool
	Matcher Matcher // optional: if provided, validates before storing
}

func (s StoreBool) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	var result interface{}

	if s.Matcher != nil {
		result = s.Matcher.Match(t, val)
	} else {
		result = val.Raw()
	}

	*s.Into = result.(bool)
	return result
}

// AsString - validates type is string and optionally matches with a Matcher
type AsString struct {
	Value   string  // expected value (if not using Matcher)
	Matcher Matcher // optional: for dynamic validation like AnyUUID, AnyString, etc.
}

func (a AsString) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	s := val.String().Raw() // This validates it's a string

	if a.Matcher != nil {
		// Re-wrap as Value for matcher
		return a.Matcher.Match(t, val)
	}

	require.Equal(t, a.Value, s, "string value mismatch")
	return s
}

// AsInt - validates type is number and optionally matches with a Matcher
type AsInt struct {
	Value   int     // expected value (if not using Matcher)
	Matcher Matcher // optional: for dynamic validation
}

func (a AsInt) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	n := val.Number().Raw() // This validates it's a number

	if a.Matcher != nil {
		return a.Matcher.Match(t, val)
	}

	require.Equal(t, float64(a.Value), n, "numeric value mismatch")
	return n
}

// AsBool - validates type is boolean and optionally checks value
type AsBool struct {
	Value   bool    // expected value (if not using Matcher)
	Matcher Matcher // optional: for dynamic validation
}

func (a AsBool) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	b := val.Boolean().Raw() // This validates it's a boolean

	if a.Matcher != nil {
		return a.Matcher.Match(t, val)
	}

	require.Equal(t, a.Value, b, "boolean value mismatch")
	return b
}

// AsNull - validates that value is null
type AsNull struct{}

func (AsNull) Match(t *testing.T, val *httpexpect.Value) interface{} {
	t.Helper()
	val.IsNull()
	return nil
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
	actual := make(map[string]*httpexpect.Value)
	for i := 0; i < int(arr.Length().Raw()); i++ {
		item := arr.Value(i)
		obj := item.Object()
		keyValue := obj.Value(b.Key).String().Raw()
		actual[keyValue] = item
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
		result := validate(t, item, e.Spec, false)
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
	return validate(t, item, a.Spec, false)
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
		results[i] = validate(t, item, spec, false)
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

	return validate(t, item, h.Spec, false)
}

// Validate recursively validates a value against a spec
// Returns captured values for further assertions
func Validate(t *testing.T, val *httpexpect.Value, spec interface{}) interface{} {
	t.Helper()
	return validate(t, val, spec, false)
}

// ValidateExact is like Validate but fails if the actual data has extra fields not in spec
func ValidateExact(t *testing.T, val *httpexpect.Value, spec interface{}) interface{} {
	t.Helper()
	return validate(t, val, spec, true)
}

func validate(t *testing.T, val *httpexpect.Value, spec interface{}, exact bool) interface{} {
	t.Helper()

	// Handle Matcher types
	if matcher, ok := spec.(Matcher); ok {
		return matcher.Match(t, val)
	}

	// Handle maps (objects)
	if specMap, ok := spec.(map[string]interface{}); ok {
		obj := val.Object()
		results := make(map[string]interface{})

		// Check for extra keys if exact mode
		if exact {
			rawObj := obj.Raw()
			for key := range rawObj {
				require.Contains(t, specMap, key, "unexpected field %q in response", key)
			}
		}

		for key, expectedVal := range specMap {
			fieldVal := obj.Value(key)
			results[key] = validate(t, fieldVal, expectedVal, exact)
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
			results[i] = validate(t, item, expectedVal, exact)
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
