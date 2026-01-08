package auth

import (
	"GoAuth/internal/domain/field"
	"encoding/json"
	"math"
)

// TODO: Implement other field types when they are implemented in the API
func validateFieldType(fieldType field.Type, value any) bool {
	switch fieldType {

	case field.String:
		_, ok := value.(string)
		return ok

	case field.Bool:
		_, ok := value.(bool)
		return ok

	case field.Int:
		switch v := value.(type) {
		case int, int8, int16, int32, int64:
			return true
		case float64:
			// JSON numbers default to float64
			return math.Trunc(v) == v
		case json.Number:
			_, err := v.Int64()
			return err == nil
		default:
			return false
		}

	default:
		return false
	}
}
