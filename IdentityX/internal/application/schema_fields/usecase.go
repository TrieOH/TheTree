package schema_fields

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/version"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"
	"fmt"

	"github.com/MintzyG/fail"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	usecaseTracer = otel.Tracer("GoAuth.SchemaFieldsService")
)

type UseCase struct {
	deps Deps
	tx   inbounds.TxRunner
}

type Deps struct {
	Schemas  outbounds.SchemaRepository
	Versions outbounds.SchemaVersionRepository
	Fields   outbounds.SchemaFieldsRepository
	Projects outbounds.ProjectRepository
}

var _ inbounds.SchemaFieldsService = (*UseCase)(nil)

func New(
	deps Deps,
	tx inbounds.TxRunner,
) inbounds.SchemaFieldsService {
	return &UseCase{
		deps: deps,
		tx:   tx,
	}
}

func (uc *UseCase) Create(ctx context.Context, in inbounds.SchemaFieldInput) (inbounds.CreateFieldsResult, error) {
	var result inbounds.CreateFieldsResult
	err := uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		result, err = uc.createInternal(ctx, in)
		return err
	})
	return result, err
}

func (uc *UseCase) createInternal(ctx context.Context, in inbounds.SchemaFieldInput) (out inbounds.CreateFieldsResult, err error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaFieldService.Create")
	defer span.End()

	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	var warnings []error

	var principal *authz.Principal
	principal, err = auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return inbounds.CreateFieldsResult{}, apierr.FromService(span, err)
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return inbounds.CreateFieldsResult{}, err
	}

	if !isOwner {
		return inbounds.CreateFieldsResult{}, fail.New(apierr.ProjectNotOwnedByPrincipal).WithArgs("cannot create fields for schema versions in a project you don't own")
	}

	var belongs bool
	belongs, err = schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return inbounds.CreateFieldsResult{}, err
	}

	if !belongs {
		return inbounds.CreateFieldsResult{}, fail.New(apierr.SchemaNotOwnedByPrincipal).WithArgs("cannot create fields for a schema you don't own")
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return inbounds.CreateFieldsResult{}, err
	}

	if latest.VersionNumber != in.VersionNumber {
		return inbounds.CreateFieldsResult{}, fail.New(apierr.SchemaVersionMismatch)
	}
	if latest.Status != version.StatusDraft {
		return inbounds.CreateFieldsResult{}, apierr.FromService(span, inbounds.ErrAddFieldsToNonDraftVersion{})
	}

	// Validate field types/owners first
	for _, f := range in.Fields {
		if !field.IsValidFieldType(f.Type) {
			return inbounds.CreateFieldsResult{}, fail.New(apierr.FIELDInvalidType).WithArgs(f.Type, f.Key)
		}
		if !field.IsValidOwnerType(f.Owner) {
			return inbounds.CreateFieldsResult{}, fail.New(apierr.FIELDInvalidOwner).WithArgs(f.Owner, f.Key)
		}
	}

	// 1. Batch create all fields
	fieldsToCreate := make([]field.Field, len(in.Fields))
	for i, f := range in.Fields {
		fieldsToCreate[i] = field.Field{
			SchemaID:        in.SchemaID,
			SchemaVersionID: latest.ID,
			Key:             f.Key,
			Type:            field.Type(f.Type),
			Owner:           field.Owner(f.Owner),
			Title:           f.Title,
			Description:     f.Description,
			Placeholder:     f.Placeholder,
			Required:        f.Required,
			Mutable:         f.Mutable,
			DefaultValue:    f.DefaultValue,
			Position:        f.Position,
		}
	}

	if err = fields.CreateBatch(ctx, fieldsToCreate); err != nil {
		return inbounds.CreateFieldsResult{}, err
	}

	// 2. Re-fetch created fields
	var createdFields []field.Field
	createdFields, err = fields.ListFromVersion(ctx, in.SchemaID, latest.ID)
	if err != nil {
		return inbounds.CreateFieldsResult{}, err
	}

	// 3. Build key->ObjectID map
	fieldKeyToID := make(map[string]uuid.UUID, len(createdFields))
	inputKeys := make(map[string]int, len(in.Fields))

	for i, f := range in.Fields {
		inputKeys[f.Key] = i
	}

	for _, f := range createdFields {
		fieldKeyToID[f.Key] = f.ObjectID
	}

	// 4. Prepare batch options and rules (with operator validation)
	var optionsBatch []field.Option
	var visRulesBatch []field.VisibilityRule
	var reqRulesBatch []field.RequiredRule

	for _, f := range in.Fields {
		fieldID, ok := fieldKeyToID[f.Key]
		if !ok {
			return inbounds.CreateFieldsResult{}, fail.New(apierr.FIELDNotFound).WithArgs(f.Key)
		}

		// Collect options (no validation needed here)
		for _, opt := range f.Options {
			optionsBatch = append(optionsBatch, field.Option{
				FieldID:  fieldID,
				Value:    opt.Value,
				Label:    opt.Label,
				Position: opt.Position,
			})
		}

		// Collect visibility rules with operator validation
		for _, rule := range f.VisibilityRules {
			if !field.IsValidRuleOperator(rule.Operator) {
				warnings = append(warnings, inbounds.ValidationWarning{
					FieldKey: f.Key,
					RuleType: "visibility",
					Operator: rule.Operator,
					Message:  fmt.Sprintf("Invalid operator '%s' for visibility rule, skipping", rule.Operator),
				})
				continue // Skip invalid rule
			}

			dependsOnID, ok := fieldKeyToID[rule.DependsOnFieldKey]
			if !ok {
				warnings = append(warnings, inbounds.ValidationWarning{
					FieldKey: f.Key,
					RuleType: "visibility",
					Operator: rule.Operator,
					Message:  fmt.Sprintf("Depends on field key '%s' not found, skipping rule", rule.DependsOnFieldKey),
				})
				continue
			}

			visRulesBatch = append(visRulesBatch, field.VisibilityRule{
				FieldID:          fieldID,
				DependsOnFieldID: dependsOnID,
				Operator:         field.RuleOperator(rule.Operator),
				Value:            rule.Value,
			})
		}

		// Collect required rules with operator validation
		for _, rule := range f.RequiredRules {
			if !field.IsValidRuleOperator(rule.Operator) {
				warnings = append(warnings, inbounds.ValidationWarning{
					FieldKey: f.Key,
					RuleType: "required",
					Operator: rule.Operator,
					Message:  fmt.Sprintf("Invalid operator '%s' for required rule, skipping", rule.Operator),
				})
				continue
			}

			dependsOnID, ok := fieldKeyToID[rule.DependsOnFieldKey]
			if !ok {
				warnings = append(warnings, inbounds.ValidationWarning{
					FieldKey: f.Key,
					RuleType: "required",
					Operator: rule.Operator,
					Message:  fmt.Sprintf("Depends on field key '%s' not found, skipping rule", rule.DependsOnFieldKey),
				})
				continue
			}

			reqRulesBatch = append(reqRulesBatch, field.RequiredRule{
				FieldID:          fieldID,
				DependsOnFieldID: dependsOnID,
				Operator:         field.RuleOperator(rule.Operator),
				Value:            rule.Value,
			})
		}
	}

	// 5. Batch insert relations
	if len(optionsBatch) > 0 {
		if err = fields.CreateOptionsBatch(ctx, optionsBatch); err != nil {
			return inbounds.CreateFieldsResult{}, err
		}
	}
	if len(visRulesBatch) > 0 {
		if err = fields.CreateVisibilityRulesBatch(ctx, visRulesBatch); err != nil {
			return inbounds.CreateFieldsResult{}, err
		}
	}
	if len(reqRulesBatch) > 0 {
		if err = fields.CreateRequiredRulesBatch(ctx, reqRulesBatch); err != nil {
			return inbounds.CreateFieldsResult{}, err
		}
	}

	// 6. Attach relations to return objects
	for i := range createdFields {
		if _, isNew := inputKeys[createdFields[i].Key]; !isNew {
			continue
		}

		createdFields[i].Options = make([]field.Option, 0)
		createdFields[i].VisibilityRules = make([]field.VisibilityRule, 0)
		createdFields[i].RequiredRules = make([]field.RequiredRule, 0)

		fieldID := createdFields[i].ObjectID
		for _, opt := range optionsBatch {
			if opt.FieldID == fieldID {
				createdFields[i].Options = append(createdFields[i].Options, opt)
			}
		}
		for _, rule := range visRulesBatch {
			if rule.FieldID == fieldID {
				createdFields[i].VisibilityRules = append(createdFields[i].VisibilityRules, rule)
			}
		}
		for _, rule := range reqRulesBatch {
			if rule.FieldID == fieldID {
				createdFields[i].RequiredRules = append(createdFields[i].RequiredRules, rule)
			}
		}
	}

	// Filter to only return newly created fields
	var resultFields []field.Field
	for _, f := range createdFields {
		if _, isNew := inputKeys[f.Key]; isNew {
			resultFields = append(resultFields, f)
		}
	}

	return inbounds.CreateFieldsResult{
		Fields:   inbounds.FieldSliceToOutputFieldSlice(resultFields),
		Warnings: warnings,
	}, nil
}
