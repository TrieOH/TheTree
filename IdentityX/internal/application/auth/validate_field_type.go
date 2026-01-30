package auth

import (
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/domain/field"
	"encoding/json"
	"math"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var validate = validator.New()

func isValidEmail(email string) bool {
	// Uses RFC 5322 compliant validation with internationalized email support
	return validate.Var(email, "required,email") == nil
}

// validateFieldValue validates a value against field type and options
func validateFieldValue(f field.Field, value any) bool {
	// First check basic type validation
	if !validateFieldType(f.Type, value) {
		return false
	}

	// For option types (select, radio, checkbox), validate the value is allowed
	if f.Type.IsOptionType() {
		return validateOptionValue(f, value)
	}

	return true
}

func validateFieldType(fieldType field.Type, value any) bool {
	switch fieldType {

	case field.String:
		_, ok := value.(string)
		return ok
	case field.Email:
		str, ok := value.(string)
		if !ok {
			return false
		}
		return isValidEmail(str)

	case field.Bool:
		_, ok := value.(bool)
		return ok

	case field.Int:
		switch v := value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32:
			return true
		case uint64:
			return v <= math.MaxInt64
		case float32:
			if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
				return false
			}
			if v > math.MaxInt32 || v < math.MinInt32 {
				return false
			}
			if v != float32(int32(v)) {
				return false
			}
			return true
		case float64:
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return false
			}
			if math.Trunc(v) != v {
				return false
			}
			const maxSafeInt64 float64 = 1<<53 - 1
			const minSafeInt64 float64 = -(1 << 53)
			if v > maxSafeInt64 || v < minSafeInt64 {
				logs.L().Debug(
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

	// FIXME implement multiple select
	case field.Select, field.Radio:
		// These should be string values (single selection)
		_, ok := value.(string)
		return ok

	case field.Checkbox:
		// FIXME implement checkbox correctly (single check or list)
		// Checkbox is an array of bools
		switch value.(type) {
		case bool:
			return true
		case []any, []string:
			return true
		default:
			return false
		}

	default:
		return false
	}
}

// validateOptionValue checks if the submitted value is in the allowed options
func validateOptionValue(f field.Field, value any) bool {
	if len(f.Options) == 0 {
		return false
	}

	// Build allowed set
	allowed := make(map[string]bool)
	for _, opt := range f.Options {
		allowed[opt.Value] = true
	}

	switch f.Type {
	case field.Select, field.Radio:
		strVal, ok := value.(string)
		if !ok {
			return false
		}
		return allowed[strVal]

	case field.Checkbox:
		// Handle array of values for multi-checkbox
		switch v := value.(type) {
		case []any:
			for _, item := range v {
				strItem, ok := item.(string)
				if !ok || !allowed[strItem] {
					return false
				}
			}
			return len(v) > 0 // Must select at least one
		case []string:
			for _, item := range v {
				if !allowed[item] {
					return false
				}
			}
			return len(v) > 0
		case string:
			// Single checkbox as string value
			return allowed[v]
		case bool:
			// Boolean checkbox (true/false type)
			return true
		default:
			return false
		}
	}

	return false
}
