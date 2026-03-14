package schema

import (
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/version"
	"GoAuth/internal/errx"
	"context"
	"encoding/json"
	"math"
	"strings"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var (
	validateMetadataTracer = otel.Tracer("validateMetadataTracer")
	validate               = validator.New()
)

func (uc *UseCase) ValidateAndConstructMetadata(
	ctx context.Context,
	projectID uuid.UUID,
	schemaType schema.Type,
	flowID string,
	customFields *json.RawMessage,
) (*json.RawMessage, error) {
	var ok bool
	var err error
	var registerSchema *schema.Schema

	ctx, span := validateMetadataTracer.Start(ctx, "ValidateAndConstructMetadata")
	defer span.End()

	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	if registerSchema, err = schemas.FindByFlowIDAndType(ctx, flowID, schemaType, projectID); err != nil {
		return nil, err
	}

	if err = registerSchema.CanRegister(ctx); err != nil {
		return nil, err
	}
	var registerVersion *version.Version
	if registerVersion, err = versions.GetCurrent(ctx, registerSchema.ID); err != nil {
		return nil, err
	}
	if err = registerVersion.CanRegister(ctx); err != nil {
		return nil, err
	}

	if ok = registerSchema.IsVersion(registerVersion.ID); !ok {
		return nil, fail.New(errx.SchemaVersionMismatch).RecordCtx(ctx)
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
	if custom, err = uc.parseCustomFields(ctx, customFields); err != nil {
		return nil, err
	}

	var validated map[string]any
	validated, err = uc.ValidateFields(ctx, custom, fieldDefs, registerFields)
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
		return nil, fail.New(errx.ProjectUserErrorEncodingMetadata).With(err).RecordCtx(ctx)
	}

	rawMetadata := json.RawMessage(marshalledMetadata)
	return &rawMetadata, nil
}

func (uc *UseCase) ValidateFields(ctx context.Context, custom map[string]any, fieldDefs map[string]field.Field, registerFields []field.Field) (map[string]any, error) {
	ctx, span := validateMetadataTracer.Start(ctx, "ValidateFields")
	defer span.End()

	var errored bool
	validationError := fail.New(errx.FIELDValidationErrorOnSchemaRegister).RecordCtx(ctx)
	validated := make(map[string]any)

	providedValues := make(map[string]any)
	for key, value := range custom {
		if _, exists := fieldDefs[key]; !exists {
			continue
		}
		providedValues[key] = value
	}

	fieldIDToKey := make(map[uuid.UUID]string)
	for _, f := range registerFields {
		fieldIDToKey[f.ObjectID] = f.Key
	}

	for _, f := range registerFields {
		value, wasProvided := providedValues[f.Key]
		isNil := wasProvided && value == nil

		isVisible := true
		if len(f.VisibilityRules) > 0 {
			isVisible = false
			for _, rule := range f.VisibilityRules {
				if uc.evaluateRuleMatches(rule.DependsOnFieldID, rule.Operator, rule.Value, providedValues, fieldIDToKey) {
					isVisible = true
					break
				}
			}
		}

		if !isVisible {
			continue
		}

		isRequired := f.Required
		if len(f.RequiredRules) > 0 {
			for _, rule := range f.RequiredRules {
				if uc.evaluateRuleMatches(rule.DependsOnFieldID, rule.Operator, rule.Value, providedValues, fieldIDToKey) {
					isRequired = true
					break
				}
			}
		}

		if isRequired {
			if !wasProvided || isNil {
				errored = true
				_ = validationError.Trace(fail.New(errx.FORMMissingRequiredField).WithArgs(f.Key).Render().Error()).RecordCtx(ctx)
				continue
			}
		} else {
			if !wasProvided || isNil {
				continue
			}
		}

		if ok := uc.validateFieldValue(f, value); !ok {
			errored = true
			_ = validationError.Trace(fail.New(errx.FORMInvalidFieldValue).WithArgs(f.Key, string(f.Type), value).Render().Error()).RecordCtx(ctx)
			continue
		}

		validated[f.Key] = value
	}

	if errored {
		return validated, validationError
	}
	return validated, nil
}

func (uc *UseCase) UpdateMetadata(ctx context.Context, customFields *json.RawMessage) error {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.UpdateMetadata")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	userID := principal.UserID
	if principal.ProjectID == nil {
		return fail.New(errx.AuthNotProjectUser).RecordCtx(ctx)
	}
	projectID := *principal.ProjectID

	u, err := uc.deps.ProjectUsers.GetByIDInternal(ctx, userID, projectID)
	if err != nil {
		return err
	}

	var input map[string]any
	if err := json.Unmarshal(*customFields, &input); err != nil {
		return fail.New(errx.RequestInvalidCustomFieldsJSON).With(err).RecordCtx(ctx)
	}

	newMetadataMap := make(map[string]any)
	if u.Metadata != nil {
		_ = json.Unmarshal(*u.Metadata, &newMetadataMap)
	}

	schemas, err := uc.deps.Schemas.List(ctx, projectID)
	if err != nil {
		return err
	}

	for _, s := range schemas {
		if s.Status != schema.StatusPublished || s.CurrentVersionID == nil {
			continue
		}

		typeMap, ok := newMetadataMap[string(s.Type)].(map[string]any)
		if !ok {
			continue
		}

		flowMap, ok := typeMap[s.FlowID].(map[string]any)
		if !ok {
			continue
		}

		userFields, _ := flowMap["fields"].(map[string]any)
		if userFields == nil {
			userFields = make(map[string]any)
		}

		// Merge input into userFields
		for k, v := range input {
			userFields[k] = v
		}

		// Validate
		fields, _ := uc.deps.Fields.GetByVersionIDWithRelations(ctx, *s.CurrentVersionID)
		fieldDefs := make(map[string]field.Field)
		for _, f := range fields {
			fieldDefs[f.Key] = f
		}

		validated, err := uc.ValidateFields(ctx, userFields, fieldDefs, fields)
		if err != nil {
			return err
		}

		// Update metadata
		flowMap["fields"] = validated
		flowMap["schema_version_id"] = s.CurrentVersionID.String()
		flowMap["schema_id"] = s.ID.String()

		// Clear cache
		cacheKey := "compat:" + projectID.String() + ":" + s.CurrentVersionID.String() + ":" + userID.String()
		uc.deps.Redis.Delete(ctx, cacheKey)
		uc.deps.Redis.Set(ctx, cacheKey, true, time.Hour)
	}

	marshalled, err := json.Marshal(newMetadataMap)
	if err != nil {
		return err
	}

	raw := json.RawMessage(marshalled)
	err = uc.deps.ProjectUsers.UpdateMetadata(ctx, userID, projectID, &raw)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) evaluateRuleMatches(dependsOnFieldID uuid.UUID, operator field.RuleOperator, ruleValue *json.RawMessage, providedValues map[string]any, fieldIDToKey map[uuid.UUID]string) bool {
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
		return uc.compareRuleValue(actualValue, ruleValue) == 0
	case field.RuleOperatorNotEquals:
		if !wasProvided || isNil {
			return true
		}
		return uc.compareRuleValue(actualValue, ruleValue) != 0
	case field.RuleOperatorGt:
		if !wasProvided || isNil {
			return false
		}
		return uc.compareRuleValue(actualValue, ruleValue) > 0
	case field.RuleOperatorGte:
		if !wasProvided || isNil {
			return false
		}
		return uc.compareRuleValue(actualValue, ruleValue) >= 0
	case field.RuleOperatorLt:
		if !wasProvided || isNil {
			return false
		}
		return uc.compareRuleValue(actualValue, ruleValue) < 0
	case field.RuleOperatorLte:
		if !wasProvided || isNil {
			return false
		}
		return uc.compareRuleValue(actualValue, ruleValue) <= 0
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
			if uc.compareValues(actualValue, a) == 0 {
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
			if uc.compareValues(actualValue, a) == 0 {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (uc *UseCase) compareRuleValue(actual interface{}, ruleValue *json.RawMessage) int {
	if ruleValue == nil {
		if actual == nil {
			return 0
		}
		return 1
	}

	var target interface{}
	if err := json.Unmarshal(*ruleValue, &target); err != nil {
		return -2
	}

	return uc.compareValues(actual, target)
}

func (uc *UseCase) compareValues(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	aNum, aOk := uc.toFloat64(a)
	bNum, bOk := uc.toFloat64(b)
	if aOk && bOk {
		if aNum < bNum {
			return -1
		} else if aNum > bNum {
			return 1
		}
		return 0
	}

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

	return -2
}

func (uc *UseCase) toFloat64(v interface{}) (float64, bool) {
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

func (uc *UseCase) parseCustomFields(ctx context.Context, customFields *json.RawMessage) (custom map[string]any, err error) {
	if customFields == nil {
		return nil, fail.New(errx.RequestMissingSchemaCustomFields).RecordCtx(ctx)
	}
	if err = json.Unmarshal(*customFields, &custom); err != nil {
		return nil, fail.New(errx.RequestInvalidCustomFieldsJSON).With(err).RecordCtx(ctx)
	}
	return custom, nil
}

func (uc *UseCase) validateFieldValue(f field.Field, value any) bool {
	if !uc.validateFieldType(f.Type, value) {
		return false
	}
	if f.Type.IsOptionType() {
		return uc.validateOptionValue(f, value)
	}
	return true
}

func (uc *UseCase) validateFieldType(fieldType field.Type, value any) bool {
	switch fieldType {
	case field.String:
		_, ok := value.(string)
		return ok
	case field.Email:
		str, ok := value.(string)
		if !ok {
			return false
		}
		return validate.Var(str, "required,email") == nil
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
				return false
			}
			return true
		case json.Number:
			_, err := v.Int64()
			return err == nil
		default:
			return false
		}
	case field.Select, field.Radio:
		_, ok := value.(string)
		return ok
	case field.Checkbox:
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

func (uc *UseCase) validateOptionValue(f field.Field, value any) bool {
	if len(f.Options) == 0 {
		return false
	}
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
		switch v := value.(type) {
		case []any:
			for _, item := range v {
				strItem, ok := item.(string)
				if !ok || !allowed[strItem] {
					return false
				}
			}
			return len(v) > 0
		case []string:
			for _, item := range v {
				if !allowed[item] {
					return false
				}
			}
			return len(v) > 0
		case string:
			return allowed[v]
		case bool:
			return true
		default:
			return false
		}
	}
	return false
}
