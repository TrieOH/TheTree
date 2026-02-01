package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/ports/outbounds"
	"context"

	"github.com/MintzyG/fail"
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
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
	return fail.New(apierr.SYSFunctionalityNotImplemented)
}

func (repo *schemaFieldsRepo) Delete(ctx context.Context, fieldID uuid.UUID) error {
	// TODO Implement me!
	return fail.New(apierr.SYSFunctionalityNotImplemented)
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}
	span.SetAttributes(attribute.Int64("cloned_fields", fieldsRows))

	optionsRows, err := repo.queries(ctx).CloneFieldOptions(ctx, sqlc.CloneFieldOptionsParams{
		DraftVersionID:  toVersionID,
		SourceVersionID: fromVersionID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}
	span.SetAttributes(attribute.Int64("cloned_options", optionsRows))

	visRows, err := repo.queries(ctx).CloneVisibilityRules(ctx, sqlc.CloneVisibilityRulesParams{
		DraftVersionID:  toVersionID,
		SourceVersionID: fromVersionID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}
	span.SetAttributes(attribute.Int64("cloned_visibility_rules", visRows))

	reqRows, err := repo.queries(ctx).CloneRequiredRules(ctx, sqlc.CloneRequiredRulesParams{
		DraftVersionID:  toVersionID,
		SourceVersionID: fromVersionID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}
	span.SetAttributes(attribute.Int64("cloned_required_rules", reqRows))

	// FIXME make this a domain error
	if reqRows+visRows+optionsRows+fieldsRows == 0 {
		apiErr := apierr.ErrNotFound.WithMsg("no affected rows").WithID(apierr.FieldNoAffectedRowsOnClone)
		apierr.RecordDomainError(span, apiErr)
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return false, sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return schema.DiffResult{}, sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
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
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}
