package auth

import (
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/domain/field"
	"encoding/json"
	"math"

	"go.uber.org/zap"
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
			// Reject special float values
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return false
			}
			// JSON numbers default to float64
			if math.Trunc(v) != v {
				return false
			}
			// Check if within int64 range (with safer bounds accounting for float64 precision)
			const maxSafeInt64 float64 = 1<<53 - 1 // 2^53 - 1, largest precise integer in float64
			const minSafeInt64 float64 = -(1 << 53)
			if v > maxSafeInt64 || v < minSafeInt64 {
				logs.L().Info(
					"rejected a value outside the safe range for float precision",
					zap.Float64("value", v),
				)
				return false
			}
			return true
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
