package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/errx"
	"GoAuth/internal/ports/outbounds"
	"context"
	"encoding/json"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type schemaFieldsRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

func (repo *schemaFieldsRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(transactions.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ outbounds.SchemaFieldsRepository = (*schemaFieldsRepo)(nil)

func NewFieldsRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbounds.SchemaFieldsRepository {
	return &schemaFieldsRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func mapSchemaFieldFromDB(dst *field.Field, src *sqlc.SchemaField) {
	dst.ObjectID = src.ObjectID
	dst.ID = src.ID
	dst.SchemaID = src.SchemaID
	dst.SchemaVersionID = src.SchemaVersionID
	dst.Key = src.Key
	dst.Type = field.Type(src.Type)
	dst.Owner = field.Owner(src.Owner)
	dst.Title = src.Title
	dst.Description = src.Description
	dst.Placeholder = src.Placeholder
	dst.Required = src.Required
	dst.Mutable = src.Mutable
	dst.DefaultValue = src.DefaultValue
	dst.Position = src.Position
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func (repo *schemaFieldsRepo) Create(ctx context.Context, toCreate field.Field) (*field.Field, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.Create",
		trace.WithAttributes(
			attribute.String("field.schema_id", toCreate.SchemaID.String()),
			attribute.String("field.schema_version_id", toCreate.SchemaVersionID.String()),
		),
	)
	defer span.End()

	sqlcSchemaField, err := repo.queries(ctx).CreateSchemaField(ctx, sqlc.CreateSchemaFieldParams{
		Key:             toCreate.Key,
		Type:            sqlc.FieldType(toCreate.Type),
		Owner:           sqlc.FieldOwner(toCreate.Owner),
		Title:           toCreate.Title,
		Description:     toCreate.Description,
		Placeholder:     toCreate.Placeholder,
		Required:        toCreate.Required,
		Mutable:         toCreate.Mutable,
		DefaultValue:    toCreate.DefaultValue,
		Position:        toCreate.Position,
		SchemaVersionID: toCreate.SchemaVersionID,
		SchemaID:        toCreate.SchemaID,
	})
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("field.id", sqlcSchemaField.ID.String()),
		attribute.String("field.object_id", sqlcSchemaField.ObjectID.String()),
	)

	var newSchemaField field.Field
	mapSchemaFieldFromDB(&newSchemaField, &sqlcSchemaField)
	return &newSchemaField, nil
}

func (repo *schemaFieldsRepo) GetByVersionID(ctx context.Context, schemaVersionID uuid.UUID) ([]field.Field, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.GetByVersionID",
		trace.WithAttributes(
			attribute.String("field.schema_version_id", schemaVersionID.String()),
		),
	)
	defer span.End()

	sqlcFields, err := repo.queries(ctx).GetFieldsByVersionID(ctx, schemaVersionID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.Int("count", len(sqlcFields)),
	)

	fields := make([]field.Field, 0, len(sqlcFields))
	for _, sqlcField := range sqlcFields {
		var newSchemaField field.Field
		mapSchemaFieldFromDB(&newSchemaField, &sqlcField)
		fields = append(fields, newSchemaField)
	}

	return fields, nil
}

func (repo *schemaFieldsRepo) ListFromSchema(ctx context.Context, schemaID uuid.UUID) ([]field.Field, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.ListFromSchema",
		trace.WithAttributes(
			attribute.String("field.schema_id", schemaID.String()),
		),
	)
	defer span.End()

	sqlcFields, err := repo.queries(ctx).ListFieldsFromSchema(ctx, schemaID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Int("count", len(sqlcFields)))

	fields := make([]field.Field, 0, len(sqlcFields))
	for _, sqlcField := range sqlcFields {
		var newSchemaField field.Field
		mapSchemaFieldFromDB(&newSchemaField, &sqlcField)
		fields = append(fields, newSchemaField)
	}
	return fields, nil
}

func (repo *schemaFieldsRepo) ListFromVersion(ctx context.Context, schemaID, versionID uuid.UUID) ([]field.Field, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.ListFromVersion",
		trace.WithAttributes(
			attribute.String("field.schema_id", schemaID.String()),
			attribute.String("field.version_id", versionID.String()),
		),
	)
	defer span.End()

	sqlcFields, err := repo.queries(ctx).ListFieldsFromVersion(ctx, sqlc.ListFieldsFromVersionParams{
		SchemaID:        schemaID,
		SchemaVersionID: versionID,
	})
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Int("count", len(sqlcFields)))

	fields := make([]field.Field, 0, len(sqlcFields))
	for _, sqlcField := range sqlcFields {
		var newSchemaField field.Field
		mapSchemaFieldFromDB(&newSchemaField, &sqlcField)
		fields = append(fields, newSchemaField)
	}
	return fields, nil
}

func (repo *schemaFieldsRepo) Update(ctx context.Context, toUpdate field.Field) error {
	// TODO Implement me!
	return fail.New(errx.SYSFunctionalityNotImplemented)
}

func (repo *schemaFieldsRepo) Delete(ctx context.Context, fieldID uuid.UUID) error {
	// TODO Implement me!
	return fail.New(errx.SYSFunctionalityNotImplemented)
}

func (repo *schemaFieldsRepo) CloneFromTo(ctx context.Context, fromVersionID, toVersionID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.CloneFromTo",
		trace.WithAttributes(
			attribute.String("from_version_id", fromVersionID.String()),
			attribute.String("to_version_id", toVersionID.String()),
		),
	)
	defer span.End()

	fieldsRows, err := repo.queries(ctx).CloneFields(ctx, sqlc.CloneFieldsParams{
		DraftVersionID:  toVersionID,
		SourceVersionID: fromVersionID,
	})
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}
	span.SetAttributes(attribute.Int64("cloned_fields", fieldsRows))

	optionsRows, err := repo.queries(ctx).CloneFieldOptions(ctx, sqlc.CloneFieldOptionsParams{
		DraftVersionID:  toVersionID,
		SourceVersionID: fromVersionID,
	})
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}
	span.SetAttributes(attribute.Int64("cloned_options", optionsRows))

	visRows, err := repo.queries(ctx).CloneVisibilityRules(ctx, sqlc.CloneVisibilityRulesParams{
		DraftVersionID:  toVersionID,
		SourceVersionID: fromVersionID,
	})
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}
	span.SetAttributes(attribute.Int64("cloned_visibility_rules", visRows))

	reqRows, err := repo.queries(ctx).CloneRequiredRules(ctx, sqlc.CloneRequiredRulesParams{
		DraftVersionID:  toVersionID,
		SourceVersionID: fromVersionID,
	})
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}
	span.SetAttributes(attribute.Int64("cloned_required_rules", reqRows))

	if reqRows+visRows+optionsRows+fieldsRows == 0 {
		apiErr := fail.New(errx.FieldNoAffectedRowsOnClone)
		return apiErr
	}

	return nil
}

func (repo *schemaFieldsRepo) DiffVersionsState(ctx context.Context, baseVersionID, draftVersionID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.DiffVersionsState",
		trace.WithAttributes(
			attribute.String("base_version_id", baseVersionID.String()),
			attribute.String("draft_version_id", draftVersionID.String()),
		),
	)
	defer span.End()

	hasChanges, err := repo.queries(ctx).DiffVersionFields(ctx, sqlc.DiffVersionFieldsParams{
		BaseVersionID:  baseVersionID,
		DraftVersionID: draftVersionID,
	})
	if err != nil {
		return false, fail.From(err).RecordCtx(ctx)
	}

	return hasChanges, nil
}

func (repo *schemaFieldsRepo) DiffVersionsFullState(ctx context.Context, baseVersionID, draftVersionID uuid.UUID) (schema.DiffResult, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.DiffVersionsState",
		trace.WithAttributes(
			attribute.String("base_version_id", baseVersionID.String()),
			attribute.String("draft_version_id", draftVersionID.String()),
		),
	)
	defer span.End()

	// sqlc generates named parameter struct: DiffVersionFieldsFullParams
	result, err := repo.queries(ctx).DiffVersionFieldsFull(ctx, sqlc.DiffVersionFieldsFullParams{
		BaseVersionID:  baseVersionID,
		DraftVersionID: draftVersionID,
	})
	if err != nil {
		return schema.DiffResult{}, fail.From(err).RecordCtx(ctx)
	}

	diff := schema.DiffResult{
		FieldsChanged:          result.FieldsChanged,
		OptionsChanged:         result.OptionsChanged,
		VisibilityRulesChanged: result.VisibilityRulesChanged,
		RequiredRulesChanged:   result.RequiredRulesChanged,
	}

	diff.Annotate(span)

	return diff, nil
}

// GetByVersionIDWithRelations returns fields with options and rules populated
func (repo *schemaFieldsRepo) GetByVersionIDWithRelations(ctx context.Context, schemaVersionID uuid.UUID) ([]field.Field, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.GetByVersionIDWithRelations",
		trace.WithAttributes(attribute.String("field.schema_version_id", schemaVersionID.String())),
	)
	defer span.End()

	// Get base fields
	fields, err := repo.GetByVersionID(ctx, schemaVersionID)
	if err != nil {
		return nil, err
	}

	if len(fields) == 0 {
		return fields, nil
	}

	// Collect field IDs (using ObjectID which is the PK)
	fieldIDs := make([]uuid.UUID, len(fields))
	for i, f := range fields {
		fieldIDs[i] = f.ObjectID
	}

	// Batch fetch options
	options, err := repo.GetOptionsByFieldIDs(ctx, fieldIDs)
	if err != nil {
		return nil, err
	}

	// Batch fetch visibility rules
	visRules, err := repo.GetVisibilityRulesByFieldIDs(ctx, fieldIDs)
	if err != nil {
		return nil, err
	}

	// Batch fetch required rules
	reqRules, err := repo.GetRequiredRulesByFieldIDs(ctx, fieldIDs)
	if err != nil {
		return nil, err
	}

	// Map relations to fields
	optionsMap := make(map[uuid.UUID][]field.Option)
	for _, opt := range options {
		optionsMap[opt.FieldID] = append(optionsMap[opt.FieldID], opt)
	}

	visRulesMap := make(map[uuid.UUID][]field.VisibilityRule)
	for _, rule := range visRules {
		visRulesMap[rule.FieldID] = append(visRulesMap[rule.FieldID], rule)
	}

	reqRulesMap := make(map[uuid.UUID][]field.RequiredRule)
	for _, rule := range reqRules {
		reqRulesMap[rule.FieldID] = append(reqRulesMap[rule.FieldID], rule)
	}

	// Attach to fields
	for i := range fields {
		fid := fields[i].ObjectID
		fields[i].Options = optionsMap[fid]
		fields[i].VisibilityRules = visRulesMap[fid]
		fields[i].RequiredRules = reqRulesMap[fid]
	}

	return fields, nil
}

// ListFromVersionWithRelations returns fields with options and rules for a specific schema version
func (repo *schemaFieldsRepo) ListFromVersionWithRelations(ctx context.Context, schemaID, versionID uuid.UUID) ([]field.Field, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.ListFromVersionWithRelations",
		trace.WithAttributes(
			attribute.String("field.schema_id", schemaID.String()),
			attribute.String("field.version_id", versionID.String()),
		),
	)
	defer span.End()

	// This could be optimized to use the relations method directly since versionID uniquely identifies fields
	return repo.GetByVersionIDWithRelations(ctx, versionID)
}

// CreateOption creates a field option
func (repo *schemaFieldsRepo) CreateOption(ctx context.Context, option field.Option) (*field.Option, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.CreateOption",
		trace.WithAttributes(attribute.String("option.field_id", option.FieldID.String())),
	)
	defer span.End()

	sqlcOpt, err := repo.queries(ctx).CreateFieldOption(ctx, sqlc.CreateFieldOptionParams{
		FieldID:  option.FieldID,
		Value:    option.Value,
		Label:    option.Label,
		Position: option.Position,
	})
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	return &field.Option{
		ID:       sqlcOpt.ID,
		FieldID:  sqlcOpt.FieldID,
		Value:    sqlcOpt.Value,
		Label:    sqlcOpt.Label,
		Position: sqlcOpt.Position,
	}, nil
}

// GetOptionsByFieldIDs batch fetches options for multiple fields
func (repo *schemaFieldsRepo) GetOptionsByFieldIDs(ctx context.Context, fieldIDs []uuid.UUID) ([]field.Option, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.GetOptionsByFieldIDs",
		trace.WithAttributes(attribute.Int("field.count", len(fieldIDs))),
	)
	defer span.End()

	sqlcOpts, err := repo.queries(ctx).GetFieldOptions(ctx, fieldIDs)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	options := make([]field.Option, len(sqlcOpts))
	for i, o := range sqlcOpts {
		options[i] = field.Option{
			ID:       o.ID,
			FieldID:  o.FieldID,
			Value:    o.Value,
			Label:    o.Label,
			Position: o.Position,
		}
	}

	return options, nil
}

// CreateVisibilityRule creates a visibility rule
func (repo *schemaFieldsRepo) CreateVisibilityRule(ctx context.Context, rule field.VisibilityRule) (*field.VisibilityRule, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.CreateVisibilityRule",
		trace.WithAttributes(attribute.String("rule.field_id", rule.FieldID.String())),
	)
	defer span.End()

	sqlcRule, err := repo.queries(ctx).CreateVisibilityRule(ctx, sqlc.CreateVisibilityRuleParams{
		FieldID:          rule.FieldID,
		DependsOnFieldID: rule.DependsOnFieldID,
		Operator:         sqlc.RuleOperator(rule.Operator),
		Value:            rule.Value,
	})
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	return &field.VisibilityRule{
		ID:               sqlcRule.ID,
		FieldID:          sqlcRule.FieldID,
		DependsOnFieldID: sqlcRule.DependsOnFieldID,
		Operator:         field.RuleOperator(sqlcRule.Operator),
		Value:            sqlcRule.Value,
		CreatedAt:        sqlcRule.CreatedAt,
	}, nil
}

// GetVisibilityRulesByFieldIDs batch fetches visibility rules
func (repo *schemaFieldsRepo) GetVisibilityRulesByFieldIDs(ctx context.Context, fieldIDs []uuid.UUID) ([]field.VisibilityRule, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.GetVisibilityRulesByFieldIDs",
		trace.WithAttributes(attribute.Int("field.count", len(fieldIDs))),
	)
	defer span.End()

	sqlcRules, err := repo.queries(ctx).GetFieldVisibilityRules(ctx, fieldIDs)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	rules := make([]field.VisibilityRule, len(sqlcRules))
	for i, r := range sqlcRules {
		rules[i] = field.VisibilityRule{
			ID:               r.ID,
			FieldID:          r.FieldID,
			DependsOnFieldID: r.DependsOnFieldID,
			Operator:         field.RuleOperator(r.Operator),
			Value:            r.Value,
			CreatedAt:        r.CreatedAt,
		}
	}

	return rules, nil
}

// CreateRequiredRule creates a required rule
func (repo *schemaFieldsRepo) CreateRequiredRule(ctx context.Context, rule field.RequiredRule) (*field.RequiredRule, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.CreateRequiredRule",
		trace.WithAttributes(attribute.String("rule.field_id", rule.FieldID.String())),
	)
	defer span.End()

	sqlcRule, err := repo.queries(ctx).CreateRequiredRule(ctx, sqlc.CreateRequiredRuleParams{
		FieldID:          rule.FieldID,
		DependsOnFieldID: rule.DependsOnFieldID,
		Operator:         sqlc.RuleOperator(rule.Operator),
		Value:            rule.Value,
	})
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	return &field.RequiredRule{
		ID:               sqlcRule.ID,
		FieldID:          sqlcRule.FieldID,
		DependsOnFieldID: sqlcRule.DependsOnFieldID,
		Operator:         field.RuleOperator(sqlcRule.Operator),
		Value:            sqlcRule.Value,
		CreatedAt:        sqlcRule.CreatedAt,
	}, nil
}

// GetRequiredRulesByFieldIDs batch fetches required rules
func (repo *schemaFieldsRepo) GetRequiredRulesByFieldIDs(ctx context.Context, fieldIDs []uuid.UUID) ([]field.RequiredRule, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.GetRequiredRulesByFieldIDs",
		trace.WithAttributes(attribute.Int("field.count", len(fieldIDs))),
	)
	defer span.End()

	sqlcRules, err := repo.queries(ctx).GetFieldRequiredRules(ctx, fieldIDs)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	rules := make([]field.RequiredRule, len(sqlcRules))
	for i, r := range sqlcRules {
		rules[i] = field.RequiredRule{
			ID:               r.ID,
			FieldID:          r.FieldID,
			DependsOnFieldID: r.DependsOnFieldID,
			Operator:         field.RuleOperator(r.Operator),
			Value:            r.Value,
			CreatedAt:        r.CreatedAt,
		}
	}

	return rules, nil
}

func (repo *schemaFieldsRepo) CreateBatch(ctx context.Context, toCreate []field.Field) error {
	if len(toCreate) == 0 {
		return nil
	}

	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.CreateBatch",
		trace.WithAttributes(attribute.Int("count", len(toCreate))),
	)
	defer span.End()

	// Convert domain options to sqlc params
	params := make([]sqlc.CreateSchemaFieldsBatchParams, len(toCreate))
	for i, f := range toCreate {
		params[i] = sqlc.CreateSchemaFieldsBatchParams{
			SchemaID:        f.SchemaID,
			SchemaVersionID: f.SchemaVersionID,
			Key:             f.Key,
			Type:            sqlc.FieldType(f.Type),
			Owner:           sqlc.FieldOwner(f.Owner),
			Title:           f.Title,
			Description:     f.Description,
			Placeholder:     f.Placeholder,
			Required:        f.Required,
			Mutable:         f.Mutable,
			DefaultValue:    f.DefaultValue,
			Position:        f.Position,
		}
	}

	// Bulk insert using COPY FROM (very efficient)
	if _, err := repo.queries(ctx).CreateSchemaFieldsBatch(ctx, params); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

func (repo *schemaFieldsRepo) CreateOptionsBatch(ctx context.Context, options []field.Option) error {
	if len(options) == 0 {
		return nil
	}

	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.CreateOptionsBatch",
		trace.WithAttributes(attribute.Int("count", len(options))),
	)
	defer span.End()

	// Convert domain options to sqlc params
	params := make([]sqlc.CreateFieldOptionsBatchParams, len(options))
	for i, opt := range options {
		params[i] = sqlc.CreateFieldOptionsBatchParams{
			FieldID:  opt.FieldID,
			Value:    opt.Value,
			Label:    opt.Label,
			Position: opt.Position,
		}
	}

	// Bulk insert using COPY FROM (very efficient)
	if _, err := repo.queries(ctx).CreateFieldOptionsBatch(ctx, params); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

func (repo *schemaFieldsRepo) CreateVisibilityRulesBatch(ctx context.Context, rules []field.VisibilityRule) error {
	if len(rules) == 0 {
		return nil
	}

	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.CreateVisibilityRulesBatch",
		trace.WithAttributes(attribute.Int("count", len(rules))),
	)
	defer span.End()

	params := make([]sqlc.CreateVisibilityRulesBatchParams, len(rules))
	for i, rule := range rules {
		params[i] = sqlc.CreateVisibilityRulesBatchParams{
			FieldID:          rule.FieldID,
			DependsOnFieldID: rule.DependsOnFieldID,
			Operator:         sqlc.RuleOperator(rule.Operator),
			Value:            rule.Value,
		}
	}

	if _, err := repo.queries(ctx).CreateVisibilityRulesBatch(ctx, params); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

func (repo *schemaFieldsRepo) CreateRequiredRulesBatch(ctx context.Context, rules []field.RequiredRule) error {
	if len(rules) == 0 {
		return nil
	}

	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.CreateRequiredRulesBatch",
		trace.WithAttributes(attribute.Int("count", len(rules))),
	)
	defer span.End()

	params := make([]sqlc.CreateRequiredRulesBatchParams, len(rules))
	for i, rule := range rules {
		params[i] = sqlc.CreateRequiredRulesBatchParams{
			FieldID:          rule.FieldID,
			DependsOnFieldID: rule.DependsOnFieldID,
			Operator:         sqlc.RuleOperator(rule.Operator),
			Value:            rule.Value,
		}
	}

	if _, err := repo.queries(ctx).CreateRequiredRulesBatch(ctx, params); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

// GetByObjectID retrieves a field by its object_id
func (repo *schemaFieldsRepo) GetByObjectID(ctx context.Context, objectID uuid.UUID) (*field.Field, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.GetByObjectID",
		trace.WithAttributes(attribute.String("field.object_id", objectID.String())),
	)
	defer span.End()

	sqlcField, err := repo.queries(ctx).GetFieldByObjectID(ctx, objectID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var f field.Field
	mapSchemaFieldFromDB(&f, &sqlcField)
	return &f, nil
}

// UpdateField updates a field with partial updates
func (repo *schemaFieldsRepo) UpdateField(ctx context.Context, objectID uuid.UUID, schemaVersionID uuid.UUID, updates map[string]interface{}) (*field.Field, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.UpdateField",
		trace.WithAttributes(attribute.String("field.object_id", objectID.String())),
	)
	defer span.End()

	// First, get the current field to preserve unchanged values
	current, err := repo.GetByObjectID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	// Build params with current values, then override with updates
	params := sqlc.UpdateFieldParams{
		ObjectID:        objectID,
		SchemaVersionID: schemaVersionID,
		Key:             current.Key,
		Type:            sqlc.FieldType(current.Type),
		Title:           current.Title,
		Description:     current.Description,
		Placeholder:     current.Placeholder,
		Required:        current.Required,
		Mutable:         current.Mutable,
		DefaultValue:    current.DefaultValue,
		Position:        current.Position,
	}

	// Override with provided updates
	if key, ok := updates["key"].(string); ok {
		params.Key = key
	}
	if fieldType, ok := updates["type"].(string); ok {
		params.Type = sqlc.FieldType(fieldType)
	}
	if title, ok := updates["title"].(string); ok {
		params.Title = title
	}
	if desc, ok := updates["description"].(*string); ok {
		params.Description = desc
	}
	if placeholder, ok := updates["placeholder"].(*string); ok {
		params.Placeholder = placeholder
	}
	if required, ok := updates["required"].(bool); ok {
		params.Required = required
	}
	if mutable, ok := updates["mutable"].(bool); ok {
		params.Mutable = mutable
	}
	if defaultValue, ok := updates["default_value"].(*json.RawMessage); ok {
		params.DefaultValue = defaultValue
	}
	if position, ok := updates["position"].(int); ok {
		params.Position = position
	}

	sqlcField, err := repo.queries(ctx).UpdateField(ctx, params)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var f field.Field
	mapSchemaFieldFromDB(&f, &sqlcField)
	return &f, nil
}

// DeleteField deletes a field by its object_id
func (repo *schemaFieldsRepo) DeleteField(ctx context.Context, objectID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.DeleteField",
		trace.WithAttributes(attribute.String("field.object_id", objectID.String())),
	)
	defer span.End()

	if err := repo.queries(ctx).DeleteField(ctx, objectID); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

// CheckFieldKeyExists checks if a field key already exists in a version
func (repo *schemaFieldsRepo) CheckFieldKeyExists(ctx context.Context, versionID uuid.UUID, key string, excludeObjectID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.CheckFieldKeyExists",
		trace.WithAttributes(
			attribute.String("version_id", versionID.String()),
			attribute.String("key", key),
		),
	)
	defer span.End()

	exists, err := repo.queries(ctx).CheckFieldKeyExists(ctx, sqlc.CheckFieldKeyExistsParams{
		SchemaVersionID: versionID,
		Key:             key,
		ObjectID:        excludeObjectID,
	})
	if err != nil {
		return false, fail.From(err).RecordCtx(ctx)
	}

	return exists, nil
}

// HasDependentRules checks if other fields have rules that depend on this field
func (repo *schemaFieldsRepo) HasDependentRules(ctx context.Context, fieldObjectID uuid.UUID) ([]field.Field, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.HasDependentRules",
		trace.WithAttributes(attribute.String("field.object_id", fieldObjectID.String())),
	)
	defer span.End()

	rows, err := repo.queries(ctx).HasDependentRules(ctx, fieldObjectID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	dependentFields := make([]field.Field, len(rows))
	for i, row := range rows {
		dependentFields[i] = field.Field{
			ObjectID: row.FieldObjectID,
			Key:      row.FieldKey,
		}
	}

	return dependentFields, nil
}

// DeleteFieldOptions deletes all options for a field
func (repo *schemaFieldsRepo) DeleteFieldOptions(ctx context.Context, fieldID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.DeleteFieldOptions",
		trace.WithAttributes(attribute.String("field_id", fieldID.String())),
	)
	defer span.End()

	if err := repo.queries(ctx).DeleteFieldOptions(ctx, fieldID); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

// DeleteFieldVisibilityRules deletes all visibility rules for a field
func (repo *schemaFieldsRepo) DeleteFieldVisibilityRules(ctx context.Context, fieldID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.DeleteFieldVisibilityRules",
		trace.WithAttributes(attribute.String("field_id", fieldID.String())),
	)
	defer span.End()

	if err := repo.queries(ctx).DeleteVisibilityRules(ctx, fieldID); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

// DeleteFieldRequiredRules deletes all required rules for a field
func (repo *schemaFieldsRepo) DeleteFieldRequiredRules(ctx context.Context, fieldID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.DeleteFieldRequiredRules",
		trace.WithAttributes(attribute.String("field_id", fieldID.String())),
	)
	defer span.End()

	if err := repo.queries(ctx).DeleteRequiredRules(ctx, fieldID); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

// GetOptionByID retrieves an option by its ID
func (repo *schemaFieldsRepo) GetOptionByID(ctx context.Context, optionID uuid.UUID) (*field.Option, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.GetOptionByID",
		trace.WithAttributes(attribute.String("option_id", optionID.String())),
	)
	defer span.End()

	sqlcOption, err := repo.queries(ctx).GetOptionByID(ctx, optionID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	return &field.Option{
		ID:       sqlcOption.ID,
		FieldID:  sqlcOption.FieldID,
		Value:    sqlcOption.Value,
		Label:    sqlcOption.Label,
		Position: sqlcOption.Position,
	}, nil
}

// SetFieldOptions replaces all options for a field
func (repo *schemaFieldsRepo) SetFieldOptions(ctx context.Context, fieldID uuid.UUID, options []field.Option) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.SetFieldOptions",
		trace.WithAttributes(
			attribute.String("field_id", fieldID.String()),
			attribute.Int("option_count", len(options)),
		),
	)
	defer span.End()

	// Delete existing options
	if err := repo.queries(ctx).SetFieldOptions(ctx, fieldID); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	// Insert new options one by one (since we need to handle the :exec query)
	for _, opt := range options {
		if err := repo.queries(ctx).CreateFieldOptionForSet(ctx, sqlc.CreateFieldOptionForSetParams{
			FieldID:  fieldID,
			Value:    opt.Value,
			Label:    opt.Label,
			Position: opt.Position,
		}); err != nil {
			return fail.From(err).RecordCtx(ctx)
		}
	}

	return nil
}

// DeleteOptionByID deletes a single option by its ID
func (repo *schemaFieldsRepo) DeleteOptionByID(ctx context.Context, optionID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.DeleteOptionByID",
		trace.WithAttributes(attribute.String("option_id", optionID.String())),
	)
	defer span.End()

	if err := repo.queries(ctx).DeleteOptionByID(ctx, optionID); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

// IsOptionValueReferenced checks if an option value is referenced in any rules
func (repo *schemaFieldsRepo) IsOptionValueReferenced(ctx context.Context, fieldID uuid.UUID, optionValue string) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.IsOptionValueReferenced",
		trace.WithAttributes(
			attribute.String("field_id", fieldID.String()),
			attribute.String("option_value", optionValue),
		),
	)
	defer span.End()

	isReferenced, err := repo.queries(ctx).CheckOptionValueInRules(ctx, optionValue)
	if err != nil {
		return false, fail.From(err).RecordCtx(ctx)
	}

	if isReferenced == nil {
		return false, nil
	}

	return *isReferenced, nil
}

// GetVisibilityRuleByID retrieves a visibility rule by its ID
func (repo *schemaFieldsRepo) GetVisibilityRuleByID(ctx context.Context, ruleID uuid.UUID) (*field.VisibilityRule, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.GetVisibilityRuleByID",
		trace.WithAttributes(attribute.String("rule_id", ruleID.String())),
	)
	defer span.End()

	sqlcRule, err := repo.queries(ctx).GetVisibilityRuleByID(ctx, ruleID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	return &field.VisibilityRule{
		ID:               sqlcRule.ID,
		FieldID:          sqlcRule.FieldID,
		DependsOnFieldID: sqlcRule.DependsOnFieldID,
		Operator:         field.RuleOperator(sqlcRule.Operator),
		Value:            sqlcRule.Value,
		CreatedAt:        sqlcRule.CreatedAt,
	}, nil
}

// SetVisibilityRules replaces all visibility rules for a field
func (repo *schemaFieldsRepo) SetVisibilityRules(ctx context.Context, fieldID uuid.UUID, rules []field.VisibilityRule) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.SetVisibilityRules",
		trace.WithAttributes(
			attribute.String("field_id", fieldID.String()),
			attribute.Int("rule_count", len(rules)),
		),
	)
	defer span.End()

	// Delete existing rules
	if err := repo.queries(ctx).SetVisibilityRules(ctx, fieldID); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	// Insert new rules
	for _, rule := range rules {
		if _, err := repo.queries(ctx).CreateVisibilityRuleForSet(ctx, sqlc.CreateVisibilityRuleForSetParams{
			FieldID:          fieldID,
			DependsOnFieldID: rule.DependsOnFieldID,
			Operator:         sqlc.RuleOperator(rule.Operator),
			Value:            rule.Value,
		}); err != nil {
			return fail.From(err).RecordCtx(ctx)
		}
	}

	return nil
}

// UpdateVisibilityRule updates a visibility rule with partial updates
func (repo *schemaFieldsRepo) UpdateVisibilityRule(ctx context.Context, ruleID uuid.UUID, updates map[string]interface{}) (*field.VisibilityRule, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.UpdateVisibilityRule",
		trace.WithAttributes(attribute.String("rule_id", ruleID.String())),
	)
	defer span.End()

	// First, get the current rule to preserve unchanged values
	current, err := repo.GetVisibilityRuleByID(ctx, ruleID)
	if err != nil {
		return nil, err
	}

	// Build params with current values, then override with updates
	params := sqlc.UpdateVisibilityRuleParams{
		ID:               ruleID,
		DependsOnFieldID: current.DependsOnFieldID,
		Operator:         sqlc.RuleOperator(current.Operator),
		Value:            current.Value,
	}

	// Override with provided updates
	if dependsOnFieldID, ok := updates["depends_on_field_id"].(uuid.UUID); ok {
		params.DependsOnFieldID = dependsOnFieldID
	}
	if operator, ok := updates["operator"].(string); ok {
		params.Operator = sqlc.RuleOperator(operator)
	}
	if value, ok := updates["value"].(*json.RawMessage); ok {
		params.Value = value
	}

	sqlcRule, err := repo.queries(ctx).UpdateVisibilityRule(ctx, params)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	return &field.VisibilityRule{
		ID:               sqlcRule.ID,
		FieldID:          sqlcRule.FieldID,
		DependsOnFieldID: sqlcRule.DependsOnFieldID,
		Operator:         field.RuleOperator(sqlcRule.Operator),
		Value:            sqlcRule.Value,
		CreatedAt:        sqlcRule.CreatedAt,
	}, nil
}

// DeleteVisibilityRuleByID deletes a visibility rule by its ID
func (repo *schemaFieldsRepo) DeleteVisibilityRuleByID(ctx context.Context, ruleID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.DeleteVisibilityRuleByID",
		trace.WithAttributes(attribute.String("rule_id", ruleID.String())),
	)
	defer span.End()

	if err := repo.queries(ctx).DeleteVisibilityRuleByID(ctx, ruleID); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

// GetRequiredRuleByID retrieves a required rule by its ID
func (repo *schemaFieldsRepo) GetRequiredRuleByID(ctx context.Context, ruleID uuid.UUID) (*field.RequiredRule, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.GetRequiredRuleByID",
		trace.WithAttributes(attribute.String("rule_id", ruleID.String())),
	)
	defer span.End()

	sqlcRule, err := repo.queries(ctx).GetRequiredRuleByID(ctx, ruleID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	return &field.RequiredRule{
		ID:               sqlcRule.ID,
		FieldID:          sqlcRule.FieldID,
		DependsOnFieldID: sqlcRule.DependsOnFieldID,
		Operator:         field.RuleOperator(sqlcRule.Operator),
		Value:            sqlcRule.Value,
		CreatedAt:        sqlcRule.CreatedAt,
	}, nil
}

// SetRequiredRules replaces all required rules for a field
func (repo *schemaFieldsRepo) SetRequiredRules(ctx context.Context, fieldID uuid.UUID, rules []field.RequiredRule) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.SetRequiredRules",
		trace.WithAttributes(
			attribute.String("field_id", fieldID.String()),
			attribute.Int("rule_count", len(rules)),
		),
	)
	defer span.End()

	// Delete existing rules
	if err := repo.queries(ctx).SetRequiredRules(ctx, fieldID); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	// Insert new rules
	for _, rule := range rules {
		if _, err := repo.queries(ctx).CreateRequiredRuleForSet(ctx, sqlc.CreateRequiredRuleForSetParams{
			FieldID:          fieldID,
			DependsOnFieldID: rule.DependsOnFieldID,
			Operator:         sqlc.RuleOperator(rule.Operator),
			Value:            rule.Value,
		}); err != nil {
			return fail.From(err).RecordCtx(ctx)
		}
	}

	return nil
}

// UpdateRequiredRule updates a required rule with partial updates
func (repo *schemaFieldsRepo) UpdateRequiredRule(ctx context.Context, ruleID uuid.UUID, updates map[string]interface{}) (*field.RequiredRule, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.UpdateRequiredRule",
		trace.WithAttributes(attribute.String("rule_id", ruleID.String())),
	)
	defer span.End()

	// First, get the current rule to preserve unchanged values
	current, err := repo.GetRequiredRuleByID(ctx, ruleID)
	if err != nil {
		return nil, err
	}

	// Build params with current values, then override with updates
	params := sqlc.UpdateRequiredRuleParams{
		ID:               ruleID,
		DependsOnFieldID: current.DependsOnFieldID,
		Operator:         sqlc.RuleOperator(current.Operator),
		Value:            current.Value,
	}

	// Override with provided updates
	if dependsOnFieldID, ok := updates["depends_on_field_id"].(uuid.UUID); ok {
		params.DependsOnFieldID = dependsOnFieldID
	}
	if operator, ok := updates["operator"].(string); ok {
		params.Operator = sqlc.RuleOperator(operator)
	}
	if value, ok := updates["value"].(*json.RawMessage); ok {
		params.Value = value
	}

	sqlcRule, err := repo.queries(ctx).UpdateRequiredRule(ctx, params)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	return &field.RequiredRule{
		ID:               sqlcRule.ID,
		FieldID:          sqlcRule.FieldID,
		DependsOnFieldID: sqlcRule.DependsOnFieldID,
		Operator:         field.RuleOperator(sqlcRule.Operator),
		Value:            sqlcRule.Value,
		CreatedAt:        sqlcRule.CreatedAt,
	}, nil
}

// DeleteRequiredRuleByID deletes a required rule by its ID
func (repo *schemaFieldsRepo) DeleteRequiredRuleByID(ctx context.Context, ruleID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.DeleteRequiredRuleByID",
		trace.WithAttributes(attribute.String("rule_id", ruleID.String())),
	)
	defer span.End()

	if err := repo.queries(ctx).DeleteRequiredRuleByID(ctx, ruleID); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}
