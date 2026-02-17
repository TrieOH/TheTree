package schema_fields

import (
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/version"
	"GoAuth/internal/errx"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"
	"fmt"

	"github.com/MintzyG/fail/v3"
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
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return inbounds.CreateFieldsResult{}, err
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return inbounds.CreateFieldsResult{}, err
	}

	if !isOwner {
		return inbounds.CreateFieldsResult{}, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot create fields for schema versions in a project you don't own").RecordCtx(ctx)
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
		return inbounds.CreateFieldsResult{}, fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot create fields for a schema you don't own").RecordCtx(ctx)
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return inbounds.CreateFieldsResult{}, err
	}

	if latest.VersionNumber != in.VersionNumber {
		return inbounds.CreateFieldsResult{}, fail.New(errx.SchemaVersionMismatch).RecordCtx(ctx)
	}
	if latest.Status != version.StatusDraft {
		return inbounds.CreateFieldsResult{}, fail.New(errx.SchemaVersionNonDraftAddFieldsNotAllowed).RecordCtx(ctx)
	}

	// Validate field types/owners first
	for _, f := range in.Fields {
		if !field.IsValidFieldType(f.Type) {
			return inbounds.CreateFieldsResult{}, fail.New(errx.FIELDInvalidType).WithArgs(f.Type, f.Key).RecordCtx(ctx)
		}
		if !field.IsValidOwnerType(f.Owner) {
			return inbounds.CreateFieldsResult{}, fail.New(errx.FIELDInvalidOwner).WithArgs(f.Owner, f.Key).RecordCtx(ctx)
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
			return inbounds.CreateFieldsResult{}, fail.New(errx.FIELDNotFound).WithArgs(f.Key).RecordCtx(ctx)
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

func (uc *UseCase) EditField(ctx context.Context, in inbounds.EditFieldInput) (*field.Field, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaFieldService.EditField")
	defer span.End()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	var principal *authz.Principal
	var err error
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot edit fields for a project you don't own").RecordCtx(ctx)
	}

	var belongs bool
	belongs, err = schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return nil, err
	}

	if !belongs {
		return nil, fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot edit fields for a schema you don't own").RecordCtx(ctx)
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return nil, err
	}

	if latest.VersionNumber != in.VersionNumber {
		return nil, fail.New(errx.SchemaVersionMismatch).RecordCtx(ctx)
	}

	if latest.Status != version.StatusDraft {
		return nil, fail.New(errx.SchemaVersionNonDraftAddFieldsNotAllowed).WithArgs("editing only allowed on draft versions").RecordCtx(ctx)
	}

	// Verify field exists
	existingField, err := fields.GetByObjectID(ctx, in.FieldObjectID)
	if err != nil {
		return nil, fail.New(errx.FIELDNotFound).WithArgs(in.FieldObjectID).RecordCtx(ctx)
	}

	// Check if field belongs to this version
	if existingField.SchemaVersionID != latest.ID {
		return nil, fail.New(errx.FIELDNotFound).WithArgs("field does not belong to this version").RecordCtx(ctx)
	}

	// Validate key uniqueness if key is being updated
	if in.Key != nil && *in.Key != existingField.Key {
		exists, err := fields.CheckFieldKeyExists(ctx, latest.ID, *in.Key, in.FieldObjectID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fail.New(errx.FIELDKeyAlreadyExists).WithArgs(*in.Key).RecordCtx(ctx)
		}
	}

	// Validate field type if it's being updated
	if in.Type != nil {
		if !field.IsValidFieldType(*in.Type) {
			return nil, fail.New(errx.FIELDInvalidType).WithArgs(*in.Type).RecordCtx(ctx)
		}
	}

	// Build updates map
	updates := make(map[string]interface{})
	if in.Key != nil {
		updates["key"] = *in.Key
	}
	if in.Type != nil {
		updates["type"] = *in.Type
	}
	if in.Title != nil {
		updates["title"] = *in.Title
	}
	if in.Description != nil {
		updates["description"] = in.Description
	}
	if in.Placeholder != nil {
		updates["placeholder"] = in.Placeholder
	}
	if in.Required != nil {
		updates["required"] = *in.Required
	}
	if in.Mutable != nil {
		updates["mutable"] = *in.Mutable
	}
	if in.DefaultValue != nil {
		updates["default_value"] = in.DefaultValue
	}
	if in.Position != nil {
		updates["position"] = *in.Position
	}

	// Update the field
	updatedField, err := fields.UpdateField(ctx, in.FieldObjectID, latest.ID, updates)
	if err != nil {
		return nil, err
	}

	return updatedField, nil
}

func (uc *UseCase) DeleteField(ctx context.Context, in inbounds.DeleteFieldInput) error {
	ctx, span := usecaseTracer.Start(ctx, "SchemaFieldService.DeleteField")
	defer span.End()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	var principal *authz.Principal
	var err error
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		return fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot delete fields for a project you don't own").RecordCtx(ctx)
	}

	var belongs bool
	belongs, err = schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return err
	}

	if !belongs {
		return fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot delete fields for a schema you don't own").RecordCtx(ctx)
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return err
	}

	if latest.VersionNumber != in.VersionNumber {
		return fail.New(errx.SchemaVersionMismatch).RecordCtx(ctx)
	}

	if latest.Status != version.StatusDraft {
		return fail.New(errx.SchemaVersionNonDraftAddFieldsNotAllowed).WithArgs("deletion only allowed on draft versions").RecordCtx(ctx)
	}

	// Verify field exists
	existingField, err := fields.GetByObjectID(ctx, in.FieldObjectID)
	if err != nil {
		return fail.New(errx.FIELDNotFound).WithArgs(in.FieldObjectID).RecordCtx(ctx)
	}

	// Check if field belongs to this version
	if existingField.SchemaVersionID != latest.ID {
		return fail.New(errx.FIELDNotFound).WithArgs("field does not belong to this version").RecordCtx(ctx)
	}

	// Check if other fields have rules that depend on this field
	dependentFields, err := fields.HasDependentRules(ctx, in.FieldObjectID)
	if err != nil {
		return err
	}

	if len(dependentFields) > 0 {
		fieldKeys := make([]string, len(dependentFields))
		for i, f := range dependentFields {
			fieldKeys[i] = f.Key
		}
		return fail.New(errx.FIELDHasDependentRules).WithArgs("field is referenced by other fields", fieldKeys).RecordCtx(ctx)
	}

	// Delete in transaction: options, rules, then field
	err = uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := fields.DeleteFieldOptions(ctx, in.FieldObjectID); err != nil {
			return err
		}
		if err := fields.DeleteFieldVisibilityRules(ctx, in.FieldObjectID); err != nil {
			return err
		}
		if err := fields.DeleteFieldRequiredRules(ctx, in.FieldObjectID); err != nil {
			return err
		}
		if err := fields.DeleteField(ctx, in.FieldObjectID); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) SetFieldOptions(ctx context.Context, in inbounds.SetFieldOptionsInput) ([]field.Option, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaFieldService.SetFieldOptions")
	defer span.End()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	var principal *authz.Principal
	var err error
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot edit options for a project you don't own").RecordCtx(ctx)
	}

	var belongs bool
	belongs, err = schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return nil, err
	}

	if !belongs {
		return nil, fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot edit options for a schema you don't own").RecordCtx(ctx)
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return nil, err
	}

	if latest.VersionNumber != in.VersionNumber {
		return nil, fail.New(errx.SchemaVersionMismatch).RecordCtx(ctx)
	}

	if latest.Status != version.StatusDraft {
		return nil, fail.New(errx.SchemaVersionNonDraftAddFieldsNotAllowed).WithArgs("options editing only allowed on draft versions").RecordCtx(ctx)
	}

	// Verify field exists
	existingField, err := fields.GetByObjectID(ctx, in.FieldObjectID)
	if err != nil {
		return nil, fail.New(errx.FIELDNotFound).WithArgs(in.FieldObjectID).RecordCtx(ctx)
	}

	// Check if field belongs to this version
	if existingField.SchemaVersionID != latest.ID {
		return nil, fail.New(errx.FIELDNotFound).WithArgs("field does not belong to this version").RecordCtx(ctx)
	}

	// Validate field type supports options
	if !existingField.Type.IsOptionType() {
		return nil, fail.New(errx.FIELDInvalidType).WithArgs("field type does not support options").RecordCtx(ctx)
	}

	// Validate unique option values
	valueSet := make(map[string]bool)
	for _, opt := range in.Options {
		if valueSet[opt.Value] {
			return nil, fail.New(errx.FIELDSameKeyForMultipleFields).WithArgs("duplicate option value", opt.Value).RecordCtx(ctx)
		}
		valueSet[opt.Value] = true
	}

	// Convert InputOption to field.Option
	options := make([]field.Option, len(in.Options))
	for i, opt := range in.Options {
		options[i] = field.Option{
			FieldID:  in.FieldObjectID,
			Value:    opt.Value,
			Label:    opt.Label,
			Position: opt.Position,
		}
	}

	// Replace all options
	if err := fields.SetFieldOptions(ctx, in.FieldObjectID, options); err != nil {
		return nil, err
	}

	return options, nil
}

func (uc *UseCase) DeleteFieldOption(ctx context.Context, in inbounds.DeleteFieldOptionInput) error {
	ctx, span := usecaseTracer.Start(ctx, "SchemaFieldService.DeleteFieldOption")
	defer span.End()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	var principal *authz.Principal
	var err error
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		return fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot delete options for a project you don't own").RecordCtx(ctx)
	}

	var belongs bool
	belongs, err = schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return err
	}

	if !belongs {
		return fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot delete options for a schema you don't own").RecordCtx(ctx)
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return err
	}

	if latest.VersionNumber != in.VersionNumber {
		return fail.New(errx.SchemaVersionMismatch).RecordCtx(ctx)
	}

	if latest.Status != version.StatusDraft {
		return fail.New(errx.SchemaVersionNonDraftAddFieldsNotAllowed).WithArgs("option deletion only allowed on draft versions").RecordCtx(ctx)
	}

	// Verify field exists
	existingField, err := fields.GetByObjectID(ctx, in.FieldObjectID)
	if err != nil {
		return fail.New(errx.FIELDNotFound).WithArgs(in.FieldObjectID).RecordCtx(ctx)
	}

	// Check if field belongs to this version
	if existingField.SchemaVersionID != latest.ID {
		return fail.New(errx.FIELDNotFound).WithArgs("field does not belong to this version").RecordCtx(ctx)
	}

	// Get the option to check its value
	option, err := fields.GetOptionByID(ctx, in.OptionID)
	if err != nil {
		return fail.New(errx.FIELDNotFound).WithArgs("option not found").RecordCtx(ctx)
	}

	// Verify option belongs to this field
	if option.FieldID != in.FieldObjectID {
		return fail.New(errx.FIELDNotFound).WithArgs("option does not belong to this field").RecordCtx(ctx)
	}

	// Check if option value is referenced in rules
	isReferenced, err := fields.IsOptionValueReferenced(ctx, in.FieldObjectID, option.Value)
	if err != nil {
		return err
	}

	if isReferenced {
		return fail.New(errx.FIELDHasDependentRules).WithArgs("option value is referenced in field rules", option.Value).RecordCtx(ctx)
	}

	// Delete the option
	if err := fields.DeleteOptionByID(ctx, in.OptionID); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) SetVisibilityRules(ctx context.Context, in inbounds.SetVisibilityRulesInput) ([]field.VisibilityRule, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaFieldService.SetVisibilityRules")
	defer span.End()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	var principal *authz.Principal
	var err error
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot edit visibility rules for a project you don't own").RecordCtx(ctx)
	}

	var belongs bool
	belongs, err = schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return nil, err
	}

	if !belongs {
		return nil, fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot edit visibility rules for a schema you don't own").RecordCtx(ctx)
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return nil, err
	}

	if latest.VersionNumber != in.VersionNumber {
		return nil, fail.New(errx.SchemaVersionMismatch).RecordCtx(ctx)
	}

	if latest.Status != version.StatusDraft {
		return nil, fail.New(errx.SchemaVersionNonDraftAddFieldsNotAllowed).WithArgs("visibility rules editing only allowed on draft versions").RecordCtx(ctx)
	}

	// Verify field exists
	existingField, err := fields.GetByObjectID(ctx, in.FieldObjectID)
	if err != nil {
		return nil, fail.New(errx.FIELDNotFound).WithArgs(in.FieldObjectID).RecordCtx(ctx)
	}

	// Check if field belongs to this version
	if existingField.SchemaVersionID != latest.ID {
		return nil, fail.New(errx.FIELDNotFound).WithArgs("field does not belong to this version").RecordCtx(ctx)
	}

	// Validate all rules
	rules := make([]field.VisibilityRule, len(in.VisibilityRules))
	for i, ruleInput := range in.VisibilityRules {
		if !field.IsValidRuleOperator(ruleInput.Operator) {
			return nil, fail.New(errx.FIELDInvalidType).WithArgs("invalid operator", ruleInput.Operator).RecordCtx(ctx)
		}

		rules[i] = field.VisibilityRule{
			FieldID:  in.FieldObjectID,
			Operator: field.RuleOperator(ruleInput.Operator),
			Value:    ruleInput.Value,
		}

		// If DependsOnFieldKey is provided, we need to resolve it to DependsOnFieldID
		if ruleInput.DependsOnFieldKey != "" {
			// Get all fields in this version to find the matching key
			versionFields, err := fields.ListFromVersion(ctx, in.SchemaID, latest.ID)
			if err != nil {
				return nil, err
			}
			found := false
			for _, f := range versionFields {
				if f.Key == ruleInput.DependsOnFieldKey {
					rules[i].DependsOnFieldID = f.ObjectID
					found = true
					break
				}
			}
			if !found {
				return nil, fail.New(errx.FIELDNotFound).WithArgs("depends_on_field_key not found", ruleInput.DependsOnFieldKey).RecordCtx(ctx)
			}
		}
	}

	// Replace all visibility rules
	if err := fields.SetVisibilityRules(ctx, in.FieldObjectID, rules); err != nil {
		return nil, err
	}

	return rules, nil
}

func (uc *UseCase) EditVisibilityRule(ctx context.Context, in inbounds.EditVisibilityRuleInput) (*field.VisibilityRule, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaFieldService.EditVisibilityRule")
	defer span.End()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	var principal *authz.Principal
	var err error
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot edit visibility rules for a project you don't own").RecordCtx(ctx)
	}

	var belongs bool
	belongs, err = schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return nil, err
	}

	if !belongs {
		return nil, fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot edit visibility rules for a schema you don't own").RecordCtx(ctx)
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return nil, err
	}

	if latest.VersionNumber != in.VersionNumber {
		return nil, fail.New(errx.SchemaVersionMismatch).RecordCtx(ctx)
	}

	if latest.Status != version.StatusDraft {
		return nil, fail.New(errx.SchemaVersionNonDraftAddFieldsNotAllowed).WithArgs("visibility rules editing only allowed on draft versions").RecordCtx(ctx)
	}

	// Verify field exists
	existingField, err := fields.GetByObjectID(ctx, in.FieldObjectID)
	if err != nil {
		return nil, fail.New(errx.FIELDNotFound).WithArgs(in.FieldObjectID).RecordCtx(ctx)
	}

	// Check if field belongs to this version
	if existingField.SchemaVersionID != latest.ID {
		return nil, fail.New(errx.FIELDNotFound).WithArgs("field does not belong to this version").RecordCtx(ctx)
	}

	// Verify rule exists and belongs to this field
	existingRule, err := fields.GetVisibilityRuleByID(ctx, in.RuleID)
	if err != nil {
		return nil, fail.New(errx.FIELDNotFound).WithArgs("visibility rule not found").RecordCtx(ctx)
	}

	if existingRule.FieldID != in.FieldObjectID {
		return nil, fail.New(errx.FIELDNotFound).WithArgs("rule does not belong to this field").RecordCtx(ctx)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if in.DependsOnFieldID != nil {
		updates["depends_on_field_id"] = *in.DependsOnFieldID
	}
	if in.Operator != nil {
		if !field.IsValidRuleOperator(*in.Operator) {
			return nil, fail.New(errx.FIELDInvalidType).WithArgs("invalid operator", *in.Operator).RecordCtx(ctx)
		}
		updates["operator"] = *in.Operator
	}
	if in.Value != nil {
		updates["value"] = in.Value
	}

	// Update the rule
	updatedRule, err := fields.UpdateVisibilityRule(ctx, in.RuleID, updates)
	if err != nil {
		return nil, err
	}

	return updatedRule, nil
}

func (uc *UseCase) DeleteVisibilityRule(ctx context.Context, in inbounds.DeleteVisibilityRuleInput) error {
	ctx, span := usecaseTracer.Start(ctx, "SchemaFieldService.DeleteVisibilityRule")
	defer span.End()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	var principal *authz.Principal
	var err error
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		return fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot delete visibility rules for a project you don't own").RecordCtx(ctx)
	}

	var belongs bool
	belongs, err = schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return err
	}

	if !belongs {
		return fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot delete visibility rules for a schema you don't own").RecordCtx(ctx)
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return err
	}

	if latest.VersionNumber != in.VersionNumber {
		return fail.New(errx.SchemaVersionMismatch).RecordCtx(ctx)
	}

	if latest.Status != version.StatusDraft {
		return fail.New(errx.SchemaVersionNonDraftAddFieldsNotAllowed).WithArgs("visibility rules deletion only allowed on draft versions").RecordCtx(ctx)
	}

	// Verify field exists
	existingField, err := fields.GetByObjectID(ctx, in.FieldObjectID)
	if err != nil {
		return fail.New(errx.FIELDNotFound).WithArgs(in.FieldObjectID).RecordCtx(ctx)
	}

	// Check if field belongs to this version
	if existingField.SchemaVersionID != latest.ID {
		return fail.New(errx.FIELDNotFound).WithArgs("field does not belong to this version").RecordCtx(ctx)
	}

	// Verify rule exists and belongs to this field
	existingRule, err := fields.GetVisibilityRuleByID(ctx, in.RuleID)
	if err != nil {
		return fail.New(errx.FIELDNotFound).WithArgs("visibility rule not found").RecordCtx(ctx)
	}

	if existingRule.FieldID != in.FieldObjectID {
		return fail.New(errx.FIELDNotFound).WithArgs("rule does not belong to this field").RecordCtx(ctx)
	}

	// Delete the rule
	if err := fields.DeleteVisibilityRuleByID(ctx, in.RuleID); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) SetRequiredRules(ctx context.Context, in inbounds.SetRequiredRulesInput) ([]field.RequiredRule, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaFieldService.SetRequiredRules")
	defer span.End()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	var principal *authz.Principal
	var err error
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot edit required rules for a project you don't own").RecordCtx(ctx)
	}

	var belongs bool
	belongs, err = schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return nil, err
	}

	if !belongs {
		return nil, fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot edit required rules for a schema you don't own").RecordCtx(ctx)
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return nil, err
	}

	if latest.VersionNumber != in.VersionNumber {
		return nil, fail.New(errx.SchemaVersionMismatch).RecordCtx(ctx)
	}

	if latest.Status != version.StatusDraft {
		return nil, fail.New(errx.SchemaVersionNonDraftAddFieldsNotAllowed).WithArgs("required rules editing only allowed on draft versions").RecordCtx(ctx)
	}

	// Verify field exists
	existingField, err := fields.GetByObjectID(ctx, in.FieldObjectID)
	if err != nil {
		return nil, fail.New(errx.FIELDNotFound).WithArgs(in.FieldObjectID).RecordCtx(ctx)
	}

	// Check if field belongs to this version
	if existingField.SchemaVersionID != latest.ID {
		return nil, fail.New(errx.FIELDNotFound).WithArgs("field does not belong to this version").RecordCtx(ctx)
	}

	// Validate all rules
	rules := make([]field.RequiredRule, len(in.RequiredRules))
	for i, ruleInput := range in.RequiredRules {
		if !field.IsValidRuleOperator(ruleInput.Operator) {
			return nil, fail.New(errx.FIELDInvalidType).WithArgs("invalid operator", ruleInput.Operator).RecordCtx(ctx)
		}

		rules[i] = field.RequiredRule{
			FieldID:  in.FieldObjectID,
			Operator: field.RuleOperator(ruleInput.Operator),
			Value:    ruleInput.Value,
		}

		// If DependsOnFieldKey is provided, we need to resolve it to DependsOnFieldID
		if ruleInput.DependsOnFieldKey != "" {
			// Get all fields in this version to find the matching key
			versionFields, err := fields.ListFromVersion(ctx, in.SchemaID, latest.ID)
			if err != nil {
				return nil, err
			}
			found := false
			for _, f := range versionFields {
				if f.Key == ruleInput.DependsOnFieldKey {
					rules[i].DependsOnFieldID = f.ObjectID
					found = true
					break
				}
			}
			if !found {
				return nil, fail.New(errx.FIELDNotFound).WithArgs("depends_on_field_key not found", ruleInput.DependsOnFieldKey).RecordCtx(ctx)
			}
		}
	}

	// Replace all required rules
	if err := fields.SetRequiredRules(ctx, in.FieldObjectID, rules); err != nil {
		return nil, err
	}

	return rules, nil
}

func (uc *UseCase) EditRequiredRule(ctx context.Context, in inbounds.EditRequiredRuleInput) (*field.RequiredRule, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaFieldService.EditRequiredRule")
	defer span.End()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	var principal *authz.Principal
	var err error
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot edit required rules for a project you don't own").RecordCtx(ctx)
	}

	var belongs bool
	belongs, err = schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return nil, err
	}

	if !belongs {
		return nil, fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot edit required rules for a schema you don't own").RecordCtx(ctx)
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return nil, err
	}

	if latest.VersionNumber != in.VersionNumber {
		return nil, fail.New(errx.SchemaVersionMismatch).RecordCtx(ctx)
	}

	if latest.Status != version.StatusDraft {
		return nil, fail.New(errx.SchemaVersionNonDraftAddFieldsNotAllowed).WithArgs("required rules editing only allowed on draft versions").RecordCtx(ctx)
	}

	// Verify field exists
	existingField, err := fields.GetByObjectID(ctx, in.FieldObjectID)
	if err != nil {
		return nil, fail.New(errx.FIELDNotFound).WithArgs(in.FieldObjectID).RecordCtx(ctx)
	}

	// Check if field belongs to this version
	if existingField.SchemaVersionID != latest.ID {
		return nil, fail.New(errx.FIELDNotFound).WithArgs("field does not belong to this version").RecordCtx(ctx)
	}

	// Verify rule exists and belongs to this field
	existingRule, err := fields.GetRequiredRuleByID(ctx, in.RuleID)
	if err != nil {
		return nil, fail.New(errx.FIELDNotFound).WithArgs("required rule not found").RecordCtx(ctx)
	}

	if existingRule.FieldID != in.FieldObjectID {
		return nil, fail.New(errx.FIELDNotFound).WithArgs("rule does not belong to this field").RecordCtx(ctx)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if in.DependsOnFieldID != nil {
		updates["depends_on_field_id"] = *in.DependsOnFieldID
	}
	if in.Operator != nil {
		if !field.IsValidRuleOperator(*in.Operator) {
			return nil, fail.New(errx.FIELDInvalidType).WithArgs("invalid operator", *in.Operator).RecordCtx(ctx)
		}
		updates["operator"] = *in.Operator
	}
	if in.Value != nil {
		updates["value"] = in.Value
	}

	// Update the rule
	updatedRule, err := fields.UpdateRequiredRule(ctx, in.RuleID, updates)
	if err != nil {
		return nil, err
	}

	return updatedRule, nil
}

func (uc *UseCase) DeleteRequiredRule(ctx context.Context, in inbounds.DeleteRequiredRuleInput) error {
	ctx, span := usecaseTracer.Start(ctx, "SchemaFieldService.DeleteRequiredRule")
	defer span.End()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	var principal *authz.Principal
	var err error
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		return fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot delete required rules for a project you don't own").RecordCtx(ctx)
	}

	var belongs bool
	belongs, err = schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return err
	}

	if !belongs {
		return fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot delete required rules for a schema you don't own").RecordCtx(ctx)
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return err
	}

	if latest.VersionNumber != in.VersionNumber {
		return fail.New(errx.SchemaVersionMismatch).RecordCtx(ctx)
	}

	if latest.Status != version.StatusDraft {
		return fail.New(errx.SchemaVersionNonDraftAddFieldsNotAllowed).WithArgs("required rules deletion only allowed on draft versions").RecordCtx(ctx)
	}

	// Verify field exists
	existingField, err := fields.GetByObjectID(ctx, in.FieldObjectID)
	if err != nil {
		return fail.New(errx.FIELDNotFound).WithArgs(in.FieldObjectID).RecordCtx(ctx)
	}

	// Check if field belongs to this version
	if existingField.SchemaVersionID != latest.ID {
		return fail.New(errx.FIELDNotFound).WithArgs("field does not belong to this version").RecordCtx(ctx)
	}

	// Verify rule exists and belongs to this field
	existingRule, err := fields.GetRequiredRuleByID(ctx, in.RuleID)
	if err != nil {
		return fail.New(errx.FIELDNotFound).WithArgs("required rule not found").RecordCtx(ctx)
	}

	if existingRule.FieldID != in.FieldObjectID {
		return fail.New(errx.FIELDNotFound).WithArgs("rule does not belong to this field").RecordCtx(ctx)
	}

	// Delete the rule
	if err := fields.DeleteRequiredRuleByID(ctx, in.RuleID); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) BatchUpdateFields(ctx context.Context, in inbounds.BatchUpdateFieldsInput) (inbounds.BatchUpdateFieldsResult, error) {
	var result inbounds.BatchUpdateFieldsResult
	err := uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		result, err = uc.batchUpdateFieldsInternal(ctx, in)
		return err
	})
	return result, err
}

func (uc *UseCase) batchUpdateFieldsInternal(ctx context.Context, in inbounds.BatchUpdateFieldsInput) (out inbounds.BatchUpdateFieldsResult, err error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaFieldService.BatchUpdateFields")
	defer span.End()

	defer func() {
		span.SetAttributes(attribute.Bool("batch_update.success", err == nil))
	}()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	var principal *authz.Principal
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return inbounds.BatchUpdateFieldsResult{}, err
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return inbounds.BatchUpdateFieldsResult{}, err
	}

	if !isOwner {
		return inbounds.BatchUpdateFieldsResult{}, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot batch update fields for a project you don't own").RecordCtx(ctx)
	}

	var belongs bool
	belongs, err = schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return inbounds.BatchUpdateFieldsResult{}, err
	}

	if !belongs {
		return inbounds.BatchUpdateFieldsResult{}, fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot batch update fields for a schema you don't own").RecordCtx(ctx)
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return inbounds.BatchUpdateFieldsResult{}, err
	}

	if latest.VersionNumber != in.VersionNumber {
		return inbounds.BatchUpdateFieldsResult{}, fail.New(errx.SchemaVersionMismatch).RecordCtx(ctx)
	}

	if latest.Status != version.StatusDraft {
		return inbounds.BatchUpdateFieldsResult{}, fail.New(errx.SchemaVersionNonDraftAddFieldsNotAllowed).WithArgs("batch update only allowed on draft versions").RecordCtx(ctx)
	}

	var warnings []error

	newFields := make([]inbounds.InputField, 0)
	existingFieldsUpdates := make(map[uuid.UUID]map[string]interface{})

	existingFieldsInVersion, err := fields.ListFromVersion(ctx, in.SchemaID, latest.ID)
	if err != nil {
		return inbounds.BatchUpdateFieldsResult{}, err
	}

	existingFieldMap := make(map[uuid.UUID]field.Field, len(existingFieldsInVersion))
	for _, f := range existingFieldsInVersion {
		existingFieldMap[f.ObjectID] = f
	}

	keyToObjectID := make(map[string]uuid.UUID)
	for _, f := range existingFieldsInVersion {
		keyToObjectID[f.Key] = f.ObjectID
	}

	for _, batchField := range in.Fields {
		if batchField.ObjectID == nil {
			if batchField.Key == nil || *batchField.Key == "" {
				return inbounds.BatchUpdateFieldsResult{}, fail.New(errx.FIELDInvalidCharactersInKey).WithArgs("field key is required").RecordCtx(ctx)
			}

			if batchField.Type == nil || *batchField.Type == "" {
				return inbounds.BatchUpdateFieldsResult{}, fail.New(errx.FIELDInvalidType).WithArgs("UNSET", *batchField.Key).RecordCtx(ctx)
			}

			if !field.IsValidFieldType(*batchField.Type) {
				return inbounds.BatchUpdateFieldsResult{}, fail.New(errx.FIELDInvalidType).WithArgs(*batchField.Type, *batchField.Key).RecordCtx(ctx)
			}

			owner := "user"
			if batchField.Mutable != nil && !*batchField.Mutable {
				owner = "system"
			}

			required := false
			if batchField.Required != nil {
				required = *batchField.Required
			}

			mutable := true
			if batchField.Mutable != nil {
				mutable = *batchField.Mutable
			}

			position := 0
			if batchField.Position != nil {
				position = *batchField.Position
			}

			newField := inbounds.InputField{
				Key:          *batchField.Key,
				Type:         *batchField.Type,
				Owner:        owner,
				Title:        "",
				Description:  batchField.Description,
				Placeholder:  batchField.Placeholder,
				Required:     required,
				Mutable:      mutable,
				DefaultValue: batchField.DefaultValue,
				Position:     position,
			}

			if batchField.Title != nil {
				newField.Title = *batchField.Title
			}

			for _, opt := range batchField.Options {
				newField.Options = append(newField.Options, inbounds.InputOption{
					Value:    opt.Value,
					Label:    opt.Label,
					Position: opt.Position,
				})
			}

			for _, rule := range batchField.VisibilityRules {
				newField.VisibilityRules = append(newField.VisibilityRules, inbounds.InputVisibilityRule{
					DependsOnFieldKey: rule.DependsOnFieldKey,
					Operator:          rule.Operator,
					Value:             rule.Value,
				})
			}

			for _, rule := range batchField.RequiredRules {
				newField.RequiredRules = append(newField.RequiredRules, inbounds.InputRequiredRule{
					DependsOnFieldKey: rule.DependsOnFieldKey,
					Operator:          rule.Operator,
					Value:             rule.Value,
				})
			}

			newFields = append(newFields, newField)
		} else {
			existingField, exists := existingFieldMap[*batchField.ObjectID]
			if !exists {
				return inbounds.BatchUpdateFieldsResult{}, fail.New(errx.FIELDNotFound).WithArgs(*batchField.ObjectID).RecordCtx(ctx)
			}

			if existingField.SchemaVersionID != latest.ID {
				return inbounds.BatchUpdateFieldsResult{}, fail.New(errx.FIELDNotFound).WithArgs("field does not belong to this version").RecordCtx(ctx)
			}

			updates := make(map[string]interface{})

			if batchField.Key != nil && *batchField.Key != existingField.Key {
				if _, keyExists := keyToObjectID[*batchField.Key]; keyExists && keyToObjectID[*batchField.Key] != *batchField.ObjectID {
					return inbounds.BatchUpdateFieldsResult{}, fail.New(errx.FIELDKeyAlreadyExists).WithArgs(*batchField.Key).RecordCtx(ctx)
				}
				updates["key"] = *batchField.Key
				keyToObjectID[*batchField.Key] = *batchField.ObjectID
			}

			if batchField.Title != nil {
				updates["title"] = *batchField.Title
			}

			if batchField.Description != nil {
				desc := *batchField.Description
				updates["description"] = &desc
			}

			if batchField.Placeholder != nil {
				placeholder := *batchField.Placeholder
				updates["placeholder"] = &placeholder
			}

			if batchField.Type != nil {
				if !field.IsValidFieldType(*batchField.Type) {
					return inbounds.BatchUpdateFieldsResult{}, fail.New(errx.FIELDInvalidType).WithArgs(*batchField.Type, existingField.Key).RecordCtx(ctx)
				}
				updates["type"] = *batchField.Type
			}

			if batchField.Required != nil {
				updates["required"] = *batchField.Required
			}

			if batchField.Mutable != nil {
				updates["mutable"] = *batchField.Mutable
			}

			if batchField.DefaultValue != nil {
				updates["default_value"] = batchField.DefaultValue
			}

			if batchField.Position != nil {
				updates["position"] = *batchField.Position
			}

			if len(updates) > 0 {
				existingFieldsUpdates[*batchField.ObjectID] = updates
			}
		}
	}

	if len(newFields) > 0 || len(existingFieldsUpdates) > 0 {
		newFieldKeys := make(map[string]bool)
		for _, f := range newFields {
			if f.Key != "" {
				if _, keyExists := keyToObjectID[f.Key]; keyExists {
					return inbounds.BatchUpdateFieldsResult{}, fail.New(errx.FIELDKeyAlreadyExists).WithArgs(f.Key).RecordCtx(ctx)
				}
				if newFieldKeys[f.Key] {
					return inbounds.BatchUpdateFieldsResult{}, fail.New(errx.FIELDSameKeyForMultipleFields).RecordCtx(ctx)
				}
				newFieldKeys[f.Key] = true
			}
		}

		for objectID, updates := range existingFieldsUpdates {
			if newKey, hasKey := updates["key"].(string); hasKey {
				if existingOID, keyExists := keyToObjectID[newKey]; keyExists && existingOID != objectID {
					return inbounds.BatchUpdateFieldsResult{}, fail.New(errx.FIELDKeyAlreadyExists).WithArgs(newKey).RecordCtx(ctx)
				}
			}
		}

		tempOffset := 10000
		originalPositions := make(map[uuid.UUID]int)

		i := 0
		for objectID := range existingFieldsUpdates {
			if updates, ok := existingFieldsUpdates[objectID]; ok {
				if pos, hasPos := updates["position"].(int); hasPos {
					originalPositions[objectID] = pos
					existingFieldsUpdates[objectID]["position"] = tempOffset + i
				}
			}
			i++
		}

		for objectID, updates := range existingFieldsUpdates {
			_, err = fields.UpdateField(ctx, objectID, latest.ID, updates)
			if err != nil {
				return inbounds.BatchUpdateFieldsResult{}, err
			}
		}

		newFieldIDs := make(map[string]uuid.UUID)
		if len(newFields) > 0 {
			fieldsToCreate := make([]field.Field, len(newFields))
			for i, f := range newFields {
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
					Position:        tempOffset + len(existingFieldsUpdates) + i,
				}
			}

			if err = fields.CreateBatch(ctx, fieldsToCreate); err != nil {
				return inbounds.BatchUpdateFieldsResult{}, err
			}

			var created []field.Field
			created, err = fields.ListFromVersion(ctx, in.SchemaID, latest.ID)
			if err != nil {
				return inbounds.BatchUpdateFieldsResult{}, err
			}

			existingKeys := make(map[string]bool)
			for _, ef := range existingFieldsInVersion {
				existingKeys[ef.Key] = true
			}

			for _, cf := range created {
				if existingKeys[cf.Key] {
					continue
				}
				newFieldIDs[cf.Key] = cf.ObjectID
			}

			for i, f := range newFields {
				targetPos := f.Position
				if targetPos == 0 {
					targetPos = len(existingFieldsInVersion) + i
				}
				objID := newFieldIDs[f.Key]
				_, err = fields.UpdateField(ctx, objID, latest.ID, map[string]interface{}{"position": targetPos})
				if err != nil {
					return inbounds.BatchUpdateFieldsResult{}, err
				}
			}
		}

		for objectID, originalPos := range originalPositions {
			_, err = fields.UpdateField(ctx, objectID, latest.ID, map[string]interface{}{"position": originalPos})
			if err != nil {
				return inbounds.BatchUpdateFieldsResult{}, err
			}
		}
	}

	updatedFields, err := fields.ListFromVersion(ctx, in.SchemaID, latest.ID)
	if err != nil {
		return inbounds.BatchUpdateFieldsResult{}, err
	}

	return inbounds.BatchUpdateFieldsResult{
		Fields:   inbounds.FieldSliceToOutputFieldSlice(updatedFields),
		Warnings: warnings,
	}, nil
}
