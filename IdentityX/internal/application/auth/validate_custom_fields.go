package auth

import (
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/version"
	"context"
	"encoding/json"
	"strings"

	"github.com/MintzyG/fail"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// validateAndConstructMetadata validates custom fields against a schema and returns structured metadata
func (uc *UseCase) validateAndConstructMetadata(
	ctx context.Context,
	span trace.Span,
	projectID uuid.UUID,
	schemaType schema.Type,
	flowID string,
	customFields *json.RawMessage,
) (*json.RawMessage, error) {
	var ok bool
	var err error
	var registerSchema *schema.Schema

	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	if registerSchema, err = schemas.FindByFlowIDAndType(ctx, flowID, schemaType, projectID); err != nil {
		return nil, err
	}

	if err = registerSchema.CanRegister(); err != nil {
		return nil, apierr.FromService(span, err)
	}
	var registerVersion *version.Version
	if registerVersion, err = versions.GetCurrent(ctx, registerSchema.ID); err != nil {
		return nil, err
	}
	if err = registerVersion.CanRegister(); err != nil {
		return nil, apierr.FromService(span, err)
	}

	if ok = registerSchema.IsVersion(registerVersion.ID); !ok {
		return nil, fail.New(apierr.SchemaVersionMismatch)
	}

	var registerFields []field.Field
	if registerFields, err = fields.GetByVersionIDWithRelations(ctx, registerVersion.ID); err != nil {
		return nil, err
	}

	fieldDefs := make(map[string]field.Field)
	for _, f := range registerFields {
		fieldDefs[f.Key] = f
	}

	var custom map[string]any
	if custom, err = parseCustomFields(customFields); err != nil {
		return nil, apierr.FromService(span, err)
	}

	var validated map[string]any
	validated, err = validateFields(custom, fieldDefs, registerFields)
	if err != nil {
		return nil, err
	}

	metadata := make(map[string]any)
	schemaPayload := make(map[string]any)

	schemaPayload["schema_id"] = registerSchema.ID.String()
	schemaPayload["schema_version_id"] = registerVersion.ID.String()
	schemaPayload["fields"] = validated

	flowMap := map[string]any{
		flowID: schemaPayload,
	}

	metadata[string(schemaType)] = flowMap

	marshalledMetadata, err := json.Marshal(metadata)
	if err != nil {
		return nil, fail.New(apierr.ProjectUserErrorEncodingMetadata).With(err)
	}

	rawMetadata := json.RawMessage(marshalledMetadata)
	return &rawMetadata, nil
}

func validateFields(custom map[string]any, fieldDefs map[string]field.Field, registerFields []field.Field) (map[string]any, error) {
	var errored bool
	validationError := fail.New(apierr.FIELDValidationErrorOnSchemaRegister)
	validated := make(map[string]any)

	// Build map of provided values (filter unknown fields, allow null for optional)
	providedValues := make(map[string]any)
	for key, value := range custom {
		// Skip unknown fields silently
		if _, exists := fieldDefs[key]; !exists {
			continue
		}
		providedValues[key] = value // value may be nil
	}

	// Build field ID to key mapping for rule evaluation
	fieldIDToKey := make(map[uuid.UUID]string)
	for _, f := range registerFields {
		fieldIDToKey[f.ObjectID] = f.Key
	}

	// Evaluate each schema field
	for _, f := range registerFields {
		value, wasProvided := providedValues[f.Key]
		isNil := wasProvided && value == nil

		// 1. CHECK VISIBILITY
		// Field is visible if: no visibility rules OR any visibility rule matches
		isVisible := true
		if len(f.VisibilityRules) > 0 {
			isVisible = false
			for _, rule := range f.VisibilityRules {
				if evaluateRuleMatches(rule.DependsOnFieldID, rule.Operator, rule.Value, providedValues, fieldIDToKey) {
					isVisible = true
					break
				}
			}
		}

		// If field is NOT visible, skip entirely (no validation, no error for missing)
		if !isVisible {
			continue
		}

		// 2. CHECK REQUIREMENT
		// Base requirement OR conditional requirement
		isRequired := f.Required

		// Check conditional required rules (these add to the requirement, not replace)
		// e.g., f2 required if f1 is filled
		if len(f.RequiredRules) > 0 {
			for _, rule := range f.RequiredRules {
				if evaluateRuleMatches(rule.DependsOnFieldID, rule.Operator, rule.Value, providedValues, fieldIDToKey) {
					isRequired = true
					break
				}
			}
		}

		// 3. VALIDATE BASED ON REQUIRED STATUS
		if isRequired {
			// Required field must be provided and not nil
			if !wasProvided || isNil {
				errored = true
				_ = validationError.Trace(fail.New(apierr.FORMMissingRequiredField).WithArgs(f.Key).Render().Error())
				continue
			}
		} else {
			// Optional field - if not provided or nil, skip validation
			if !wasProvided || isNil {
				continue
			}
		}

		// 4. TYPE & OPTIONS VALIDATION
		if ok := validateFieldValue(f, value); !ok {
			errored = true
			_ = validationError.Trace(fail.New(apierr.FORMInvalidFieldValue).WithArgs(f.Key, string(f.Type), value).Render().Error())
			continue
		}

		validated[f.Key] = value
	}

	if errored {
		return validated, validationError
	}
	return validated, nil
}

func evaluateRuleMatches(dependsOnFieldID uuid.UUID, operator field.RuleOperator, ruleValue *json.RawMessage, providedValues map[string]any, fieldIDToKey map[uuid.UUID]string) bool {
	depKey, exists := fieldIDToKey[dependsOnFieldID]
	if !exists {
		return false
	}

	actualValue, wasProvided := providedValues[depKey]
	isNil := wasProvided && actualValue == nil

	switch operator {
	case field.RuleOperatorExists:
		return wasProvided && !isNil
	case field.RuleOperatorNotExists:
		return !wasProvided || isNil
	case field.RuleOperatorEquals:
		if !wasProvided || isNil {
			return false
		}
		return compareRuleValue(actualValue, ruleValue) == 0
	case field.RuleOperatorNotEquals:
		if !wasProvided || isNil {
			return true
		}
		return compareRuleValue(actualValue, ruleValue) != 0
	case field.RuleOperatorGt:
		if !wasProvided || isNil {
			return false
		}
		return compareRuleValue(actualValue, ruleValue) > 0
	case field.RuleOperatorGte:
		if !wasProvided || isNil {
			return false
		}
		return compareRuleValue(actualValue, ruleValue) >= 0
	case field.RuleOperatorLt:
		if !wasProvided || isNil {
			return false
		}
		return compareRuleValue(actualValue, ruleValue) < 0
	case field.RuleOperatorLte:
		if !wasProvided || isNil {
			return false
		}
		return compareRuleValue(actualValue, ruleValue) <= 0
	case field.RuleOperatorContains:
		str, ok := actualValue.(string)
		if !ok {
			return false
		}
		var target string
		if ruleValue != nil {
			if err := json.Unmarshal(*ruleValue, &target); err != nil {
				logs.L().Error("error unmarshalling rule", zap.Error(err))
				return false
			}
		}
		return strings.Contains(str, target)
	case field.RuleOperatorIn:
		if !wasProvided || isNil {
			return false
		}
		if ruleValue == nil {
			return false
		}
		var allowed []interface{}
		if err := json.Unmarshal(*ruleValue, &allowed); err != nil {
			var single interface{}
			if err := json.Unmarshal(*ruleValue, &single); err == nil {
				allowed = []interface{}{single}
			}
		}
		for _, a := range allowed {
			if compareValues(actualValue, a) == 0 {
				return true
			}
		}
		return false
	case field.RuleOperatorNotIn:
		if !wasProvided || isNil {
			return true
		}
		if ruleValue == nil {
			return true
		}
		var allowed []interface{}
		if err := json.Unmarshal(*ruleValue, &allowed); err != nil {
			var single interface{}
			if err := json.Unmarshal(*ruleValue, &single); err == nil {
				allowed = []interface{}{single}
			}
		}
		for _, a := range allowed {
			if compareValues(actualValue, a) == 0 {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// compareRuleValue unmarshalls rule JSON value and compares
func compareRuleValue(actual interface{}, ruleValue *json.RawMessage) int {
	if ruleValue == nil {
		if actual == nil {
			return 0
		}
		return 1
	}

	var target interface{}
	if err := json.Unmarshal(*ruleValue, &target); err != nil {
		return -2 // Not comparable
	}

	return compareValues(actual, target)
}

// compareValues compares two values numerically or lexically
// Returns: -1 (a < b), 0 (a == b), 1 (a > b), -2 (not comparable)
func compareValues(a, b interface{}) int {
	// Handle nil
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// Try numeric comparison first
	aNum, aOk := toFloat64(a)
	bNum, bOk := toFloat64(b)
	if aOk && bOk {
		if aNum < bNum {
			return -1
		} else if aNum > bNum {
			return 1
		}
		return 0
	}

	// String comparison
	aStr, aOk := a.(string)
	bStr, bOk := b.(string)
	if aOk && bOk {
		if aStr < bStr {
			return -1
		} else if aStr > bStr {
			return 1
		}
		return 0
	}

	// Bool comparison
	aBool, aOk := a.(bool)
	bBool, bOk := b.(bool)
	if aOk && bOk {
		if aBool == bBool {
			return 0
		}
		if !aBool && bBool {
			return -1
		}
		return 1
	}

	return -2 // Not comparable
}

// toFloat64 attempts to convert a value to float64 for numeric comparison
func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}

func parseCustomFields(customFields *json.RawMessage) (custom map[string]any, err error) {
	if customFields == nil {
		return nil, fail.New(apierr.RequestMissingSchemaCustomFields)
	}
	if err = json.Unmarshal(*customFields, &custom); err != nil {
		return nil, fail.New(apierr.RequestInvalidCustomFieldsJSON).With(err)
	}
	return custom, nil
}
