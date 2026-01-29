package permissions

import (
	"encoding/json"
	"fmt"
	"time"
)

// ValidateCondition validates a condition structure recursively
func ValidateCondition(c Condition) error {
	return validateConditionRecursive(c, "")
}

func validateConditionRecursive(c Condition, path string) error {
	currentPath := path
	if path != "" {
		currentPath = path + "."
	}

	// Check logical combinators first
	if c.And != nil {
		if len(*c.And) == 0 {
			return fmt.Errorf("%sand: AND conditions cannot be empty", currentPath)
		}
		for i, cond := range *c.And {
			if err := validateConditionRecursive(cond, fmt.Sprintf("%sand[%d]", currentPath, i)); err != nil {
				return err
			}
		}
		return nil
	}

	if c.Or != nil {
		if len(*c.Or) == 0 {
			return fmt.Errorf("%sor: OR conditions cannot be empty", currentPath)
		}
		for i, cond := range *c.Or {
			if err := validateConditionRecursive(cond, fmt.Sprintf("%sor[%d]", currentPath, i)); err != nil {
				return err
			}
		}
		return nil
	}

	if c.Not != nil {
		return validateConditionRecursive(*c.Not, fmt.Sprintf("%snot", currentPath))
	}

	return validateLeafCondition(c, currentPath)
}

func validateLeafCondition(c Condition, path string) error {
	if c.Op == "" {
		return fmt.Errorf("%sop: operator is required", path)
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
		return fmt.Errorf("%sop: unsupported operator '%s'", path, c.Op)
	}

	// Validate based on operator group
	if temporalOps[c.Op] {
		return validateTemporalCondition(c, path)
	}
	return validatePredicateCondition(c, path)
}

func validateTemporalCondition(c Condition, path string) error {
	switch c.Op {
	case OpGraceBefore, OpGraceAfter, OpGraceAround:
		if c.Field == "" {
			return fmt.Errorf("%sfield: required for '%s' operator", path, c.Op)
		}
		if c.Margin == "" {
			return fmt.Errorf("%smargin: required for '%s' operator", path, c.Op)
		}
		if _, err := time.ParseDuration(c.Margin); err != nil {
			return fmt.Errorf("%smargin: invalid duration '%s': %w", path, c.Margin, err)
		}
		// Ensure no conflicting fields
		if c.Path != "" || c.Value != nil || c.Ref != "" || c.FieldStart != "" || c.FieldEnd != "" {
			return fmt.Errorf("%s: unexpected fields for grace operator (only field and margin allowed)", path)
		}

	case OpGraceDuration:
		if c.FieldStart == "" {
			return fmt.Errorf("%sfield_start: required for '%s' operator", path, c.Op)
		}
		if c.FieldEnd == "" {
			return fmt.Errorf("%sfield_end: required for '%s' operator", path, c.Op)
		}
		if c.Margin == "" {
			return fmt.Errorf("%smargin: required for '%s' operator", path, c.Op)
		}
		if _, err := time.ParseDuration(c.Margin); err != nil {
			return fmt.Errorf("%smargin: invalid duration '%s': %w", path, c.Margin, err)
		}
		// Ensure no conflicting fields
		if c.Path != "" || c.Value != nil || c.Ref != "" || c.Field != "" {
			return fmt.Errorf("%s: unexpected fields for grace_duration operator (only field_start, field_end, and margin allowed)", path)
		}
	}

	return nil
}

func validatePredicateCondition(c Condition, path string) error {
	if c.Path == "" {
		return fmt.Errorf("%spath: required for predicate operator '%s'", path, c.Op)
	}

	// For 'exists' operator, value/ref are not needed
	if c.Op == OpExists {
		if c.Value != nil || c.Ref != "" {
			return fmt.Errorf("%s: value and ref should not be provided for 'exists' operator", path)
		}
		return nil
	}

	// For other operators, need either Value or Ref
	if c.Value == nil && c.Ref == "" {
		return fmt.Errorf("%s: either value or ref must be provided for operator '%s'", path, c.Op)
	}

	// Check no temporal fields are set
	if c.Field != "" || c.FieldStart != "" || c.FieldEnd != "" || c.Margin != "" {
		return fmt.Errorf("%s: temporal fields (field, field_start, field_end, margin) not allowed for predicate operator", path)
	}

	return nil
}

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

func DecodeCondition(raw *json.RawMessage) (*Condition, error) {
	if raw == nil {
		return nil, nil
	}

	var c Condition
	if err := json.Unmarshal(*raw, &c); err != nil {
		return nil, fmt.Errorf("failed to decode condition: %w", err)
	}

	return &c, nil
}

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
