package permissions

import (
	"GoAuth/internal/apierr"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MintzyG/fail/v3"
)

// FIXME make these accept ctx and record to span

// ValidateCondition validates a condition structure recursively
func ValidateCondition(c Condition) error {
	return validateConditionRecursive(c, "")
}

func validateConditionRecursive(c Condition, path string) error {
	// Check logical combinators first
	if c.And != nil {
		if len(*c.And) == 0 {
			location := "and"
			if path != "" {
				location = fmt.Sprintf("%sand", path)
			}
			return fail.New(apierr.PERMissionLogicalConditionValidationError).WithArgs(location, "AND")
		}
		for i, cond := range *c.And {
			childPath := fmt.Sprintf("%sand[%d]", path, i)
			if err := validateConditionRecursive(cond, childPath); err != nil {
				return err
			}
		}
		return nil
	}

	if c.Or != nil {
		if len(*c.Or) == 0 {
			location := "or"
			if path != "" {
				location = fmt.Sprintf("%sor", path)
			}
			return fail.New(apierr.PERMissionLogicalConditionValidationError).WithArgs(location, "OR")
		}
		for i, cond := range *c.Or {
			childPath := fmt.Sprintf("%sor[%d]", path, i)
			if err := validateConditionRecursive(cond, childPath); err != nil {
				return err
			}
		}
		return nil
	}

	if c.Not != nil {
		childPath := fmt.Sprintf("%snot", path)
		return validateConditionRecursive(*c.Not, childPath)
	}

	// Leaf node validation
	return validateLeafCondition(c, path)
}

func formatPathPrefix(path string) string {
	if path == "" {
		return ""
	}
	return path + "."
}

func validateLeafCondition(c Condition, path string) error {
	if c.Op == "" {
		return validationError(path, "op: operator is required")
	}

	// Define valid operator groups
	temporalOps := map[string]bool{
		OpGraceBefore: true, OpGraceAfter: true, OpGraceAround: true, OpGraceDuration: true,
	}
	predicateOps := map[string]bool{
		OpEq: true, OpNeq: true, OpGt: true, OpGte: true, OpLt: true, OpLte: true,
		OpStartsWith: true, OpEndsWith: true, OpMatches: true,
		OpIn: true, OpContainsAll: true, OpContainsAny: true, OpContains: true,
		OpExists: true,
	}

	if !temporalOps[c.Op] && !predicateOps[c.Op] {
		return validationError(path, fmt.Sprintf("op: unsupported operator '%s'", c.Op))
	}

	// Validate based on operator group
	if temporalOps[c.Op] {
		return validateTemporalCondition(c, path)
	}
	return validatePredicateCondition(c, path)
}

func validationError(path, msg string) error {
	if path == "" {
		return fail.New(apierr.PERMissionConditionValidationError).WithArgs("'.'").Trace(msg)
	}
	return fail.New(apierr.PERMissionConditionValidationError).WithArgs(path).Trace(msg)
}

func validateTemporalCondition(c Condition, path string) error {
	switch c.Op {
	case OpGraceBefore, OpGraceAfter, OpGraceAround:
		if c.Field == "" {
			return validationError(path, fmt.Sprintf("field: required for '%s' operator", c.Op))
		}
		if c.Margin == "" {
			return validationError(path, fmt.Sprintf("margin: required for '%s' operator", c.Op))
		}
		if _, err := time.ParseDuration(c.Margin); err != nil {
			return validationError(path, fmt.Sprintf("margin: invalid duration '%s': %v", c.Margin, err))
		}
		// Ensure no conflicting fields
		if c.Path != "" || c.Value != nil || c.Ref != "" || c.FieldStart != "" || c.FieldEnd != "" {
			return validationError(path, "unexpected fields for grace operator (only field and margin allowed)")
		}

	case OpGraceDuration:
		if c.FieldStart == "" {
			return validationError(path, fmt.Sprintf("field_start: required for '%s' operator", c.Op))
		}
		if c.FieldEnd == "" {
			return validationError(path, fmt.Sprintf("field_end: required for '%s' operator", c.Op))
		}
		if c.Margin == "" {
			return validationError(path, fmt.Sprintf("margin: required for '%s' operator", c.Op))
		}
		if _, err := time.ParseDuration(c.Margin); err != nil {
			return validationError(path, fmt.Sprintf("margin: invalid duration '%s': %v", c.Margin, err))
		}
		// Ensure no conflicting fields
		if c.Path != "" || c.Value != nil || c.Ref != "" || c.Field != "" {
			return validationError(path, "unexpected fields for grace_duration operator (only field_start, field_end, and margin allowed)")
		}
	}

	return nil
}

func validatePredicateCondition(c Condition, path string) error {
	if c.Path == "" {
		return validationError(path, fmt.Sprintf("path: required for predicate operator '%s'", c.Op))
	}

	// For 'exists' operator, value/ref are not needed
	if c.Op == OpExists {
		if c.Value != nil || c.Ref != "" {
			return validationError(path, "value and ref should not be provided for 'exists' operator")
		}
		return nil
	}

	// For other operators, need either Value or Ref
	if c.Value == nil && c.Ref == "" {
		return validationError(path, "either value or ref must be provided for operator '"+c.Op+"'")
	}

	// Check no temporal fields are set
	if c.Field != "" || c.FieldStart != "" || c.FieldEnd != "" || c.Margin != "" {
		// Match the exact expected error message from the test
		return validationError(path, "temporal fields not allowed for predicate operator")
	}

	return nil
}

// EncodeCondition encodes a Condition into a json.RawMessage for database storage
func EncodeCondition(c *Condition) (*json.RawMessage, error) {
	if c == nil {
		return nil, nil
	}

	data, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to encode condition: %w", err)
	}

	raw := json.RawMessage(data)
	return &raw, nil
}

// DecodeCondition decodes a json.RawMessage into a Condition
func DecodeCondition(raw *json.RawMessage) (*Condition, error) {
	if raw == nil {
		return nil, nil
	}

	// Check for JSON null
	if string(*raw) == "null" {
		return nil, nil
	}

	var c Condition
	if err := json.Unmarshal(*raw, &c); err != nil {
		return nil, fmt.Errorf("failed to decode condition: %w", err)
	}

	return &c, nil
}

// DecodeAndValidateCondition decodes and validates a condition in one step
func DecodeAndValidateCondition(raw *json.RawMessage) (*Condition, error) {
	c, err := DecodeCondition(raw)
	if err != nil {
		return nil, err
	}

	if c != nil {
		if err := ValidateCondition(*c); err != nil {
			return nil, fmt.Errorf("condition validation failed: %w", err)
		}
	}

	return c, nil
}
