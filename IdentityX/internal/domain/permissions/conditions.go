package permissions

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Condition Root - can be any of these types
type Condition struct {
	And *[]Condition `json:"and,omitempty"`
	Or  *[]Condition `json:"or,omitempty"`
	Not *Condition   `json:"not,omitempty"`

	// For leaf conditions (only one of these should be set)
	Path       string      `json:"path,omitempty"`
	Op         string      `json:"op,omitempty"`
	Value      interface{} `json:"value,omitempty"`
	Ref        string      `json:"ref,omitempty"`
	Margin     string      `json:"margin,omitempty"`
	Field      string      `json:"field,omitempty"`
	FieldStart string      `json:"field_start,omitempty"`
	FieldEnd   string      `json:"field_end,omitempty"`

	// Permission check specific
	Action  string  `json:"action,omitempty"`
	Object  string  `json:"object,omitempty"`
	ScopeID *string `json:"scope_id,omitempty"`
}

// Leaf Condition Types:

// Logical Combinators

type LogicalAnd struct {
	And []Condition `json:"and"`
}
type LogicalOr struct {
	Or []Condition `json:"or"`
}
type LogicalNot struct {
	Not Condition `json:"not"`
}

// GraceCondition
// Single-Point Grace (applies to one timestamp)
// Op: "grace_before" | "grace_after" | "grace_around"
type GraceCondition struct {
	Field  string `json:"field"`  // e.g., "resource.start_time", "environment.now"
	Op     string `json:"op"`     // grace_before | grace_after | grace_around
	Margin string `json:"margin"` // ISO 8601 duration: "15m", "1h", "30m", "24h"
}

// grace_before: [field - margin, field]      -- Only before the timestamp
// grace_after:  [field, field + margin]      -- Only after the timestamp
// grace_around: [field - margin, field + margin] -- Both before and after

// GraceDurationCondition
// Duration Grace (applies to a time range [start, end])
// Op: "grace_duration"
// Behavior: [start - margin, end + margin] (around the whole duration)
type GraceDurationCondition struct {
	FieldStart string `json:"field_start"` // e.g., "resource.start_time"
	FieldEnd   string `json:"field_end"`   // e.g., "resource.end_time"
	Op         string `json:"op"`          // "grace_duration"
	Margin     string `json:"margin"`      // ISO 8601 duration
}

// Valid window: [start - margin, end + margin]

// PredicateCondition
// Standard Predicates (Comparison operators)
// Ops: "eq" | "neq" | "gt" | "gte" | "lt" | "lte" |
//
//	"startsWith" | "endsWith" | "contains" | "matches" |
//	"in" | "contains" (for arrays) | "exists"
type PredicateCondition struct {
	Path  string      `json:"path"`            // dot notation: "subject.id", "resource.status"
	Op    string      `json:"op"`              // operator
	Value interface{} `json:"value,omitempty"` // literal value
	Ref   string      `json:"ref,omitempty"`   // reference to another path
}

// Supported Operators:
const (
	// Equality

	OpEq  = "eq"
	OpNeq = "neq"

	// Comparison (numeric/temporal)

	OpGt  = "gt"
	OpGte = "gte"
	OpLt  = "lt"
	OpLte = "lte"

	// String

	OpStartsWith = "startsWith"
	OpEndsWith   = "endsWith"
	OpMatches    = "matches" // regex

	// Array

	OpIn          = "in" // value in array at path
	OpContainsAll = "containsAll"
	OpContainsAny = "containsAny"

	// String and Array

	OpContains = "contains"

	// Existence

	OpExists = "exists"

	// Temporal (Grace)

	OpGraceBefore   = "grace_before"
	OpGraceAfter    = "grace_after"
	OpGraceAround   = "grace_around"
	OpGraceDuration = "grace_duration"
)

type ValidationError struct {
	Field    string `json:"field"`
	Required string `json:"required"`
	Motive   string `json:"motive"` // MACHINE_READABLE_CODE
}

// Examples:

// Example 1: Own resource only
/*
{
  "path": "subject.id",
  "op": "eq",
  "ref": "resource.owner_id"
}
*/

// Example 2: Check-in grace period (15m before start only)
/*
{
  "field": "resource.start_time",
  "op": "grace_before",
  "margin": "15m"
}
*/

// Example 3: Event access with duration grace (30m before start, 30m after end)
/*
{
  "field_start": "resource.start_time",
  "field_end": "resource.end_time",
  "op": "grace_duration",
  "margin": "30m"
}
*/

type Motive struct {
	Code    string
	Field   string
	Message string
	Path    string
}

type ConditionContext struct {
	Subject     map[string]interface{}
	Resource    map[string]interface{}
	Environment map[string]interface{}
	Scope       map[string]interface{}
}

func (m Motive) IsEmpty() bool {
	return m.Code == ""
}

type CheckerInput struct {
	EntityID              uuid.UUID
	Action, Object, Depth string
	ProjectID, ScopeID    *uuid.UUID
	Resource              map[string]interface{}
}

func (c Condition) Evaluate(ctx context.Context, evalCtx *ConditionContext) (bool, Motive, error) {
	// Logical combinators
	if c.And != nil {
		return evalAnd(ctx, evalCtx, *c.And)
	}
	if c.Or != nil {
		return evalOr(ctx, evalCtx, *c.Or)
	}
	if c.Not != nil {
		return evalNot(ctx, evalCtx, *c.Not)
	}

	// Leaf node evaluation
	return evalLeaf(evalCtx, c)
}

func evalAnd(ctx context.Context, evalCtx *ConditionContext, conditions []Condition) (bool, Motive, error) {
	for _, c := range conditions {
		ok, motive, err := c.Evaluate(ctx, evalCtx)
		if err != nil {
			return false, motive, err
		}
		if !ok {
			return false, Motive{Code: "AND_FAILED", Message: fmt.Sprintf("condition failed: %v", c)}, nil
		}
	}
	return true, Motive{}, nil
}

func evalOr(ctx context.Context, evalCtx *ConditionContext, conditions []Condition) (bool, Motive, error) {
	for _, c := range conditions {
		ok, motive, err := c.Evaluate(ctx, evalCtx)
		if err != nil {
			return false, motive, err
		}
		if ok {
			return true, Motive{}, nil
		}
	}
	return false, Motive{Code: "OR_FAILED", Message: "no condition in OR chain matched"}, nil
}

func evalNot(ctx context.Context, evalCtx *ConditionContext, condition Condition) (bool, Motive, error) {
	ok, motive, err := condition.Evaluate(ctx, evalCtx)
	if err != nil {
		return false, motive, err
	}
	return !ok, Motive{}, nil
}

func evalLeaf(evalCtx *ConditionContext, c Condition) (bool, Motive, error) {
	switch c.Op {
	// Temporal operators
	case OpGraceBefore, OpGraceAfter, OpGraceAround:
		return evalSingleGrace(evalCtx, c)
	case OpGraceDuration:
		return evalGraceDuration(evalCtx, c)

	// Standard predicates
	default:
		return evalPredicate(evalCtx, c)
	}
}

// Temporal Evaluation
func evalSingleGrace(evalCtx *ConditionContext, c Condition) (bool, Motive, error) {
	fieldValue, err := resolvePath(evalCtx, c.Field)
	if err != nil {
		return false, Motive{Code: "FIELD_MISSING", Field: c.Field, Message: err.Error()}, nil
	}

	timestamp, ok := fieldValue.(time.Time)
	if !ok {
		// Try parsing string
		if tsStr, ok := fieldValue.(string); ok {
			parsed, err := time.Parse(time.RFC3339, tsStr)
			if err != nil {
				return false, Motive{Code: "INVALID_TIME", Field: c.Field, Message: "cannot parse as RFC3339"}, nil
			}
			timestamp = parsed
		} else {
			return false, Motive{Code: "TYPE_ERROR", Field: c.Field, Message: "expected time.Time or RFC3339 string"}, nil
		}
	}

	margin, err := time.ParseDuration(c.Margin)
	if err != nil {
		return false, Motive{Code: "INVALID_MARGIN", Message: fmt.Sprintf("cannot parse duration: %s", c.Margin)}, nil
	}

	now := getNow(evalCtx)

	switch c.Op {
	case OpGraceBefore:
		// Valid: [timestamp - margin, timestamp]
		windowStart := timestamp.Add(-margin)
		return now.After(windowStart) && now.Before(timestamp) || now.Equal(timestamp), Motive{}, nil
	case OpGraceAfter:
		// Valid: [timestamp, timestamp + margin]
		windowEnd := timestamp.Add(margin)
		return (now.After(timestamp) || now.Equal(timestamp)) && now.Before(windowEnd), Motive{}, nil
	case OpGraceAround:
		// Valid: [timestamp - margin, timestamp + margin]
		windowStart := timestamp.Add(-margin)
		windowEnd := timestamp.Add(margin)
		return now.After(windowStart) && now.Before(windowEnd), Motive{}, nil
	}

	return false, Motive{Code: "UNKNOWN_OP", Message: c.Op}, nil
}

func evalGraceDuration(evalCtx *ConditionContext, c Condition) (bool, Motive, error) {
	// Resolve both timestamps
	startVal, err := resolvePath(evalCtx, c.FieldStart)
	if err != nil {
		return false, Motive{Code: "FIELD_MISSING", Field: c.FieldStart, Message: err.Error()}, nil
	}

	endVal, err := resolvePath(evalCtx, c.FieldEnd)
	if err != nil {
		return false, Motive{Code: "FIELD_MISSING", Field: c.FieldEnd, Message: err.Error()}, nil
	}

	start := toTime(startVal)
	end := toTime(endVal)

	if start.IsZero() || end.IsZero() {
		return false, Motive{Code: "INVALID_TIME", Message: "start or end time could not be parsed"}, nil
	}

	margin, err := time.ParseDuration(c.Margin)
	if err != nil {
		return false, Motive{Code: "INVALID_MARGIN", Message: c.Margin}, nil
	}

	now := getNow(evalCtx)

	// GraceDuration is always [start - margin, end + margin] (around)
	windowStart := start.Add(-margin)
	windowEnd := end.Add(margin)

	return now.After(windowStart) && now.Before(windowEnd), Motive{}, nil
}

// Standard Predicate Evaluation
func evalPredicate(evalCtx *ConditionContext, c Condition) (bool, Motive, error) {
	// Get value from path
	leftVal, err := resolvePath(evalCtx, c.Path)
	if err != nil {
		// For exists check, this is expected
		if c.Op == OpExists {
			return false, Motive{}, nil
		}
		return false, Motive{Code: "PATH_NOT_FOUND", Path: c.Path}, nil
	}

	// If using ref, get right value from ref path
	var rightVal interface{}
	if c.Ref != "" {
		rightVal, err = resolvePath(evalCtx, c.Ref)
		if err != nil {
			return false, Motive{Code: "REF_NOT_FOUND", Path: c.Ref}, nil
		}
	} else {
		rightVal = c.Value
	}

	// Type-aware contains
	if c.Op == OpContains {
		return evalContains(leftVal, rightVal)
	}

	// Standard comparisons
	switch c.Op {
	case OpEq:
		return leftVal == rightVal, Motive{}, nil
	case OpNeq:
		return leftVal != rightVal, Motive{}, nil
	case OpGt, OpGte, OpLt, OpLte:
		return evalComparison(leftVal, rightVal, c.Op)
	case OpStartsWith:
		return evalStringOp(leftVal, rightVal, func(s, substr string) bool {
			return strings.HasPrefix(s, substr)
		})
	case OpEndsWith:
		return evalStringOp(leftVal, rightVal, func(s, substr string) bool {
			return strings.HasSuffix(s, substr)
		})
	case OpMatches:
		return evalRegex(leftVal, rightVal)
	case OpIn:
		return evalIn(leftVal, rightVal)
	case OpExists:
		return true, Motive{}, nil // If we got here, path exists
	}

	return false, Motive{Code: "UNKNOWN_OP", Message: c.Op}, nil
}

// Type-aware contains: string (substring) vs array (membership)
func evalContains(container, item interface{}) (bool, Motive, error) {
	if s, ok := container.(string); ok {
		if substr, ok := item.(string); ok {
			return strings.Contains(s, substr), Motive{}, nil
		}
		return false, Motive{Code: "TYPE_MISMATCH", Message: "string contains requires string operand"}, nil
	}

	if arr, ok := container.([]interface{}); ok {
		target := fmt.Sprintf("%v", item)
		for _, v := range arr {
			if fmt.Sprintf("%v", v) == target {
				return true, Motive{}, nil
			}
		}
		return false, Motive{}, nil
	}

	if arr, ok := container.([]string); ok {
		target := fmt.Sprintf("%v", item)
		for _, v := range arr {
			if v == target {
				return true, Motive{}, nil
			}
		}
		return false, Motive{}, nil
	}

	return false, Motive{Code: "TYPE_ERROR", Message: "contains requires string or array"}, nil
}

func evalComparison(left, right interface{}, op string) (bool, Motive, error) {
	leftF, leftOk := toFloat(left)
	rightF, rightOk := toFloat(right)

	if leftOk && rightOk {
		switch op {
		case OpGt:
			return leftF > rightF, Motive{}, nil
		case OpGte:
			return leftF >= rightF, Motive{}, nil
		case OpLt:
			return leftF < rightF, Motive{}, nil
		case OpLte:
			return leftF <= rightF, Motive{}, nil
		}
	}

	leftT, leftOk := toTime(left), !toTime(left).IsZero()
	rightT, rightOk := toTime(right), !toTime(right).IsZero()

	if leftOk && rightOk {
		switch op {
		case OpGt:
			return leftT.After(rightT), Motive{}, nil
		case OpGte:
			return leftT.After(rightT) || leftT.Equal(rightT), Motive{}, nil
		case OpLt:
			return leftT.Before(rightT), Motive{}, nil
		case OpLte:
			return leftT.Before(rightT) || leftT.Equal(rightT), Motive{}, nil
		}
	}

	return false, Motive{Code: "TYPE_MISMATCH", Message: "comparison requires numeric or time types"}, nil
}

func evalStringOp(left, right interface{}, op func(string, string) bool) (bool, Motive, error) {
	leftS, ok := left.(string)
	if !ok {
		leftS = fmt.Sprintf("%v", left)
	}
	rightS, ok := right.(string)
	if !ok {
		rightS = fmt.Sprintf("%v", right)
	}
	return op(leftS, rightS), Motive{}, nil
}

func evalRegex(left, right interface{}) (bool, Motive, error) {
	pattern, ok := right.(string)
	if !ok {
		return false, Motive{Code: "INVALID_REGEX", Message: "pattern must be string"}, nil
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, Motive{Code: "INVALID_REGEX", Message: err.Error()}, nil
	}

	leftS := fmt.Sprintf("%v", left)
	return re.MatchString(leftS), Motive{}, nil
}

// evalIn is a type-aware reverse contains
func evalIn(item, container interface{}) (bool, Motive, error) {
	return evalContains(container, item)
}

func resolvePath(evalCtx *ConditionContext, path string) (interface{}, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}

	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid path: %s", path)
	}

	var source map[string]interface{}
	switch parts[0] {
	case "subject":
		source = evalCtx.Subject
	case "resource":
		source = evalCtx.Resource
	case "environment":
		source = evalCtx.Environment
	case "scope":
		source = evalCtx.Scope
	default:
		return nil, fmt.Errorf("unknown source: %s", parts[0])
	}

	current := source
	for i := 1; i < len(parts); i++ {
		part := parts[i]

		// Check for array indexing syntax: key[index]
		openBracket := strings.Index(part, "[")
		closeBracket := strings.Index(part, "]")

		if openBracket != -1 && closeBracket != -1 && openBracket < closeBracket {
			// Parse array access
			key := part[:openBracket]
			indexStr := part[openBracket+1 : closeBracket]

			// Get the array/slice from the map
			val, ok := current[key]
			if !ok {
				return nil, fmt.Errorf("key not found: %s", key)
			}

			// Parse index
			var index int
			if _, err := fmt.Sscanf(indexStr, "%d", &index); err != nil {
				return nil, fmt.Errorf("invalid array index '%s' in path: %s", indexStr, path)
			}

			// Handle different array/slice types
			switch arr := val.(type) {
			case []interface{}:
				if index < 0 || index >= len(arr) {
					return nil, fmt.Errorf("array index out of bounds: %d (array length: %d)", index, len(arr))
				}
				val = arr[index]
			case []string:
				if index < 0 || index >= len(arr) {
					return nil, fmt.Errorf("array index out of bounds: %d (array length: %d)", index, len(arr))
				}
				val = arr[index]
			default:
				return nil, fmt.Errorf("cannot index into type %T at %s", val, key)
			}

			// If this is the last part, return the value
			if i == len(parts)-1 {
				return val, nil
			}

			// Otherwise, continue traversing from this value
			if nextMap, ok := val.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return nil, fmt.Errorf("cannot traverse into type %T at %s", val, part)
			}
		} else {
			// Regular map key access
			val, ok := current[part]
			if !ok {
				return nil, fmt.Errorf("key not found: %s", part)
			}

			if i == len(parts)-1 {
				return val, nil
			}

			// Continue traversing
			if nextMap, ok := val.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return nil, fmt.Errorf("cannot traverse into type %T at %s", val, part)
			}
		}
	}

	return current, nil
}

func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

func toTime(v interface{}) time.Time {
	switch t := v.(type) {
	case time.Time:
		return t
	case string:
		parsed, _ := time.Parse(time.RFC3339, t)
		return parsed
	default:
		return time.Time{}
	}
}

func getNow(evalCtx *ConditionContext) time.Time {
	if evalCtx.Environment != nil {
		if nowVal, ok := evalCtx.Environment["now"]; ok {
			return toTime(nowVal)
		}
	}
	return time.Now().UTC()
}
