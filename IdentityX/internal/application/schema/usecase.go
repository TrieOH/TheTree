package schema

import (
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/version"
	"GoAuth/internal/errx"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	usecaseTracer = otel.Tracer("GoAuth.SchemaService")
)

type UseCase struct {
	deps Deps
	tx   inbounds.TxRunner
}

type Deps struct {
	Schemas      outbounds.SchemaRepository
	Versions     outbounds.SchemaVersionRepository
	Fields       outbounds.SchemaFieldsRepository
	Projects     outbounds.ProjectRepository
	ProjectUsers outbounds.ProjectUserRepository
	Cache        outbounds.CacheService
}

var _ inbounds.SchemaService = (*UseCase)(nil)

func New(
	deps Deps,
	tx inbounds.TxRunner,
) inbounds.SchemaService {
	return &UseCase{
		deps: deps,
		tx:   tx,
	}
}

func (uc *UseCase) Draft(ctx context.Context, in inbounds.SchemaServiceInput) (*inbounds.SchemaOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.Draft")
	defer span.End()

	var err error
	defer func() {
		span.SetAttributes(attribute.Bool("draft.success", err == nil))
	}()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas

	var principal *authz.Principal
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	if principal.ProjectID != nil && *principal.ProjectID != in.ProjectID {
		return nil, fail.New(errx.ProjectNotFound).RecordCtx(ctx)
	}

	if in.FlowID == "" {
		return nil, fail.New(errx.SchemaEmptyFlowID).RecordCtx(ctx)
	}

	if in.SchemaType == "" {
		return nil, fail.New(errx.SchemaEmptySchemaType).RecordCtx(ctx)
	}

	in.FlowID = strings.TrimSpace(strings.ToLower(in.FlowID))
	in.SchemaType = strings.TrimSpace(strings.ToLower(in.SchemaType))

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot draft a schema for a project you don't own").RecordCtx(ctx)
	}

	if !schema.IsValidSchemaType(in.SchemaType) {
		return nil, fail.New(errx.SchemaInvalidSchemaType).RecordCtx(ctx)
	}

	// FlowIDs cannot be the same as schema types so if this matches we error out
	if schema.IsValidSchemaType(in.FlowID) {
		return nil, fail.New(errx.SchemaInvalidFlowID).WithArgs("flow id can't be the same as a schema type").RecordCtx(ctx)
	}

	if schema.IsFlowIDReserved(in.FlowID) {
		return nil, fail.New(errx.SchemaFlowIDIsReserved).WithArgs(in.FlowID).RecordCtx(ctx)
	}

	validSchemaType := schema.Type(in.SchemaType)

	var exists bool
	exists, err = schemas.Exists(ctx, schema.Schema{
		FlowID:    in.FlowID,
		ProjectID: in.ProjectID,
		Type:      validSchemaType,
	})
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fail.New(errx.SchemaFlowIDAlreadyExistsInType).RecordCtx(ctx)
	}

	var drafted *schema.Schema
	drafted, err = schemas.Draft(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		Title:     in.Title,
		FlowID:    in.FlowID,
		Type:      validSchemaType,
	})
	if err != nil {
		return nil, err
	}

	return inbounds.SchemaToSchemaOutput(drafted), nil
}

func (uc *UseCase) Publish(ctx context.Context, in inbounds.SchemaServiceInput) error {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.Publish")
	defer span.End()

	var err error
	defer func() {
		span.SetAttributes(attribute.Bool("publish.success", err == nil))
	}()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions

	var principal *authz.Principal
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	if principal.ProjectID != nil && *principal.ProjectID != in.ProjectID {
		return fail.New(errx.ProjectNotFound).RecordCtx(ctx)
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		return fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot publish a schema for a project you don't own").RecordCtx(ctx)
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
		return fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot publish a schema you don't own").RecordCtx(ctx)
	}

	var toPublish *schema.Schema
	toPublish, err = schemas.FindByID(ctx, in.SchemaID, in.ProjectID)
	if err != nil {
		return err
	}

	if toPublish.Status != schema.StatusDraft {
		if toPublish.Status == schema.StatusPublished {
			err = fail.New(errx.SchemaTryingToPublishPublished).RecordCtx(ctx)
		} else if toPublish.Status == schema.StatusArchived {
			err = fail.New(errx.SchemaTryingToPublishArchived).RecordCtx(ctx)
		} else {
			err = fail.New(errx.SchemaNoValidStatus).WithArgs(toPublish.Status).RecordCtx(ctx)
		}
		return err
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil && !fail.Is(err, errx.SQLNotFound) {
		return err
	}

	if err != nil && fail.Is(err, errx.SQLNotFound) {
		return fail.New(errx.SCHEMANoPublishedVersion).RecordCtx(ctx)
	}

	if latest.VersionNumber == 1 && latest.Status == version.StatusDraft {
		return fail.New(errx.SchemaHasOnlyDraftVersion).RecordCtx(ctx)
	}

	if latest.VersionNumber == 1 && latest.Status == version.StatusArchived {
		return fail.New(errx.SchemaHasOnlyArchivedVersion).RecordCtx(ctx)
	}

	if err = schemas.Publish(ctx, schema.Schema{
		ID:        in.SchemaID,
		ProjectID: in.ProjectID,
	}); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) GetByID(ctx context.Context, in inbounds.SchemaServiceInput) (*inbounds.SchemaOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.GetByID")
	defer span.End()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	if principal.ProjectID != nil && *principal.ProjectID != in.ProjectID {
		return nil, fail.New(errx.ProjectNotFound).RecordCtx(ctx)
	}

	isOwner, err := projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot get a schema from a project you don't own").RecordCtx(ctx)
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
		return nil, fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot get a schema you don't own").RecordCtx(ctx)
	}

	found, err := schemas.FindByID(ctx, in.SchemaID, in.ProjectID)
	if err != nil {
		return nil, err
	}

	return inbounds.SchemaToSchemaOutput(found), nil
}

func (uc *UseCase) GetVerbose(ctx context.Context, in inbounds.SchemaServiceInput) (*inbounds.SchemaVerboseOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.GetVerbose")
	defer span.End()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	if principal.ProjectID != nil && *principal.ProjectID != in.ProjectID {
		return nil, fail.New(errx.ProjectNotFound).RecordCtx(ctx)
	}

	isOwner, err := projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot get a schema from a project you don't own").RecordCtx(ctx)
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
		return nil, fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot get a schema you don't own").RecordCtx(ctx)
	}

	var schemaPart *schema.Schema
	schemaPart, err = schemas.FindByID(ctx, in.SchemaID, in.ProjectID)
	if err != nil {
		return nil, err
	}

	schemaOutput := &inbounds.SchemaVerboseOutput{
		SchemaOutput: inbounds.SchemaOutput{
			ID:               schemaPart.ID,
			ProjectID:        schemaPart.ProjectID,
			Title:            schemaPart.Title,
			FlowID:           schemaPart.FlowID,
			Type:             string(schemaPart.Type),
			CurrentVersionID: schemaPart.CurrentVersionID,
			Status:           string(schemaPart.Status),
			CreatedAt:        schemaPart.CreatedAt,
			UpdatedAt:        schemaPart.UpdatedAt,
		},
	}

	var versionsPart []version.Version
	versionsPart, err = versions.List(ctx, in.SchemaID)
	if err != nil {
		return nil, err
	}

	versionsOutput := make([]inbounds.VersionVerboseOutput, 0, len(versionsPart))
	for _, v := range versionsPart {
		versionOutput := inbounds.VersionVerboseOutput{
			SchemaVersionOutput: inbounds.SchemaVersionOutput{
				ID:               v.ID,
				SchemaID:         v.SchemaID,
				BasedOnVersionID: v.BasedOnVersionID,
				VersionNumber:    v.VersionNumber,
				Status:           v.Status,
				CreatedAt:        v.CreatedAt,
				UpdatedAt:        v.UpdatedAt,
			},
			Fields: nil,
		}
		versionsOutput = append(versionsOutput, versionOutput)
	}

	schemaOutput.Versions = versionsOutput

	var fieldsPart []field.Field
	fieldsPart, err = fields.ListFromSchema(ctx, in.SchemaID)
	if err != nil {
		return nil, err
	}

	for i := range schemaOutput.Versions {
		for _, f := range fieldsPart {
			if f.SchemaVersionID != schemaOutput.Versions[i].ID {
				continue
			}

			schemaOutput.Versions[i].Fields = append(schemaOutput.Versions[i].Fields, inbounds.OutputField{
				ObjectID:        f.ObjectID,
				ID:              f.ID,
				Key:             f.Key,
				SchemaID:        f.SchemaID,
				SchemaVersionID: f.SchemaVersionID,
				Type:            string(f.Type),
				Owner:           string(f.Owner),
				Title:           f.Title,
				Description:     f.Description,
				Placeholder:     f.Placeholder,
				Required:        f.Required,
				Mutable:         f.Mutable,
				DefaultValue:    f.DefaultValue,
				Position:        f.Position,
				CreatedAt:       f.CreatedAt,
				UpdatedAt:       f.UpdatedAt,
			})
		}
	}

	return schemaOutput, nil
}

func (uc *UseCase) GetIDsFromProjectID(ctx context.Context, projectID uuid.UUID) ([]uuid.UUID, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.GetIDsFromProjectID",
		trace.WithAttributes(attribute.String("projectID", projectID.String())),
	)
	defer span.End()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	if principal.ProjectID != nil && *principal.ProjectID != projectID {
		return nil, fail.New(errx.ProjectNotFound).RecordCtx(ctx)
	}

	isOwner, err := projects.IsOwnerOf(ctx, projectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot get schema IDs from a project you don't own").RecordCtx(ctx)
	}

	IDs, err := schemas.GetIDsFromProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("schema.count", len(IDs)))

	return IDs, nil
}

func (uc *UseCase) List(ctx context.Context, projectID uuid.UUID) ([]inbounds.SchemaOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.List",
		trace.WithAttributes(attribute.String("projectID", projectID.String())),
	)
	defer span.End()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	if principal.ProjectID != nil && *principal.ProjectID != projectID {
		return nil, fail.New(errx.ProjectNotFound).RecordCtx(ctx)
	}

	isOwner, err := projects.IsOwnerOf(ctx, projectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot get schemas from a project you don't own").RecordCtx(ctx)
	}

	schemasOutput, err := schemas.List(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return inbounds.SchemaSliceToSchemaOutputSlice(schemasOutput), nil
}

func (uc *UseCase) GetLatestForm(ctx context.Context, in inbounds.SchemaServiceInput) (*inbounds.FormOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.GetLatestForm")
	defer span.End()

	return uc.getForm(ctx, in, 0, span) // 0 indicates "current/latest published"
}

func (uc *UseCase) GetFormByVersion(ctx context.Context, in inbounds.SchemaServiceInput, versionNumber int) (*inbounds.FormOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.GetFormByVersion",
		trace.WithAttributes(attribute.Int("version", versionNumber)),
	)
	defer span.End()

	return uc.getForm(ctx, in, versionNumber, span)
}

func (uc *UseCase) getForm(ctx context.Context, in inbounds.SchemaServiceInput, versionNumber int, span trace.Span) (*inbounds.FormOutput, error) {
	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fieldsRepo := uc.deps.Fields

	// Auth check
	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	if principal.ProjectID != nil && *principal.ProjectID != in.ProjectID {
		return nil, fail.New(errx.ProjectNotFound).RecordCtx(ctx)
	}

	isOwner, err := projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot get form for a project you don't own").RecordCtx(ctx)
	}

	// Get schema - either by ID or by FlowID+Type
	var s *schema.Schema
	if in.SchemaID != uuid.Nil {
		belongs, err := schemas.BelongsToProject(ctx, schema.Schema{
			ProjectID: in.ProjectID,
			ID:        in.SchemaID,
		})
		if err != nil {
			return nil, err
		}
		if !belongs {
			return nil, fail.New(errx.SchemaNotOwnedByPrincipal).WithArgs("cannot get form for a schema you don't own").RecordCtx(ctx)
		}

		s, err = schemas.FindByID(ctx, in.SchemaID, in.ProjectID)
		if err != nil {
			return nil, err
		}
	} else {
		// Lookup by FlowID + Type (e.g., flow_id="login", type="core")
		s, err = schemas.FindByFlowIDAndType(ctx, in.FlowID, schema.Type(in.SchemaType), in.ProjectID)
		if err != nil {
			return nil, err
		}
	}

	// Get the specific version
	var v *version.Version
	if versionNumber == 0 {
		// Get current published version
		if s.CurrentVersionID == nil {
			return nil, fail.New(errx.SCHEMANoPublishedVersion).RecordCtx(ctx)
		}
		v, err = versions.GetByID(ctx, *s.CurrentVersionID)
		if err != nil {
			return nil, err
		}
	} else {
		v, err = versions.GetByVersionNumber(ctx, s.ID, versionNumber)
		if err != nil {
			return nil, err
		}
	}

	// Get all fields for this version
	// Note: You'll need to ensure your fields repository loads options and rules
	domainFields, err := fieldsRepo.ListFromVersionWithRelations(ctx, s.ID, v.ID)
	if err != nil {
		return nil, err
	}

	// Map to output
	form := &inbounds.FormOutput{
		SchemaID:      s.ID,
		Title:         s.Title,
		FlowID:        s.FlowID,
		SchemaType:    string(s.Type),
		VersionID:     v.ID,
		VersionNumber: v.VersionNumber,
		Status:        string(v.Status),
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
		Fields:        make([]inbounds.FormField, 0, len(domainFields)),
	}

	for _, f := range domainFields {
		ff := inbounds.FormField{
			ID:              f.ID,
			ObjectID:        f.ObjectID,
			Key:             f.Key,
			Type:            string(f.Type),
			Owner:           string(f.Owner),
			Title:           f.Title,
			Description:     f.Description,
			Placeholder:     f.Placeholder,
			Required:        f.Required,
			Mutable:         f.Mutable,
			DefaultValue:    f.DefaultValue,
			Position:        f.Position,
			CreatedAt:       f.CreatedAt,
			UpdatedAt:       f.UpdatedAt,
			Options:         make([]inbounds.FormOption, 0),
			VisibilityRules: make([]inbounds.FormRule, 0),
			RequiredRules:   make([]inbounds.FormRule, 0),
		}

		// Map options if present (select, radio, checkbox)
		for _, opt := range f.Options {
			ff.Options = append(ff.Options, inbounds.FormOption{
				ID:       opt.ID,
				Value:    opt.Value,
				Label:    opt.Label,
				Position: opt.Position,
			})
		}

		// Map rules
		for _, r := range f.VisibilityRules {
			ff.VisibilityRules = append(ff.VisibilityRules, inbounds.FormRule{
				ID:               r.ID,
				DependsOnFieldID: r.DependsOnFieldID,
				Operator:         string(r.Operator),
				Value:            r.Value,
			})
		}
		for _, r := range f.RequiredRules {
			ff.RequiredRules = append(ff.RequiredRules, inbounds.FormRule{
				ID:               r.ID,
				DependsOnFieldID: r.DependsOnFieldID,
				Operator:         string(r.Operator),
				Value:            r.Value,
			})
		}

		form.Fields = append(form.Fields, ff)
	}

	return form, nil
}

func (uc *UseCase) GetUpgradeForm(ctx context.Context) ([]inbounds.FormResponse, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.GetUpgradeForm")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	userID := principal.UserID
	if principal.ProjectID == nil {
		return []inbounds.FormResponse{}, nil
	}
	projectID := *principal.ProjectID

	u, err := uc.deps.ProjectUsers.GetByIDInternal(ctx, userID, projectID)
	if err != nil {
		return nil, err
	}

	schemas, err := uc.deps.Schemas.List(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var metadata map[string]any
	if u.Metadata != nil {
		_ = json.Unmarshal(*u.Metadata, &metadata)
	}

	var responses []inbounds.FormResponse
	for _, s := range schemas {
		if s.Status != schema.StatusPublished || s.CurrentVersionID == nil {
			continue
		}

		// Check if user is compatible
		isCompatible := false
		if metadata != nil {
			if typeMap, ok := metadata[string(s.Type)].(map[string]any); ok {
				if flowMap, ok := typeMap[s.FlowID].(map[string]any); ok {
					metadataVersionIDStr, _ := flowMap["schema_version_id"].(string)
					if metadataVersionIDStr == s.CurrentVersionID.String() {
						isCompatible = true
					} else {
						userFields, _ := flowMap["fields"].(map[string]any)
						fields, _ := uc.deps.Fields.GetByVersionIDWithRelations(ctx, *s.CurrentVersionID)
						fieldDefs := make(map[string]field.Field)
						for _, f := range fields {
							fieldDefs[f.Key] = f
						}
						_, err := uc.ValidateFields(ctx, userFields, fieldDefs, fields)
						if err == nil {
							isCompatible = true
						}
					}
				}
			}
		}

		if !isCompatible {
			// Generate form
			form, err := uc.getFormWithValues(ctx, s, metadata)
			if err != nil {
				continue
			}
			responses = append(responses, *form)
		}
	}

	return responses, nil
}

func (uc *UseCase) getFormWithValues(ctx context.Context, s schema.Schema, metadata map[string]any) (*inbounds.FormResponse, error) {
	v, err := uc.deps.Versions.GetCurrent(ctx, s.ID)
	if err != nil {
		return nil, err
	}

	fields, err := uc.deps.Fields.GetByVersionIDWithRelations(ctx, v.ID)
	if err != nil {
		return nil, err
	}

	// Get user's current values for this schema
	var userValues map[string]any
	if metadata != nil {
		if typeMap, ok := metadata[string(s.Type)].(map[string]any); ok {
			if flowMap, ok := typeMap[s.FlowID].(map[string]any); ok {
				userValues, _ = flowMap["fields"].(map[string]any)
			}
		}
	}

	formFields := make([]inbounds.FormField, len(fields))
	for i, f := range fields {
		ff := inbounds.FormField{
			ID:           f.ID,
			ObjectID:     f.ObjectID,
			Key:          f.Key,
			Type:         string(f.Type),
			Owner:        string(f.Owner),
			Title:        f.Title,
			Description:  f.Description,
			Placeholder:  f.Placeholder,
			Required:     f.Required,
			Mutable:      f.Mutable,
			DefaultValue: f.DefaultValue,
			Position:     f.Position,
			CreatedAt:    f.CreatedAt,
			UpdatedAt:    f.UpdatedAt,
		}

		// Inject current value into DefaultValue if it exists
		if val, ok := userValues[f.Key]; ok {
			marshalled, _ := json.Marshal(val)
			raw := json.RawMessage(marshalled)
			ff.DefaultValue = &raw
		}

		// Map options
		ff.Options = make([]inbounds.FormOption, len(f.Options))
		for j, opt := range f.Options {
			ff.Options[j] = inbounds.FormOption{
				ID:       opt.ID,
				Value:    opt.Value,
				Label:    opt.Label,
				Position: opt.Position,
			}
		}

		// Map rules
		ff.VisibilityRules = make([]inbounds.FormRule, len(f.VisibilityRules))
		for j, rule := range f.VisibilityRules {
			ff.VisibilityRules[j] = inbounds.FormRule{
				ID:               rule.ID,
				DependsOnFieldID: rule.DependsOnFieldID,
				Operator:         string(rule.Operator),
				Value:            rule.Value,
			}
		}

		ff.RequiredRules = make([]inbounds.FormRule, len(f.RequiredRules))
		for j, rule := range f.RequiredRules {
			ff.RequiredRules[j] = inbounds.FormRule{
				ID:               rule.ID,
				DependsOnFieldID: rule.DependsOnFieldID,
				Operator:         string(rule.Operator),
				Value:            rule.Value,
			}
		}

		formFields[i] = ff
	}

	return &inbounds.FormResponse{
		SchemaID:      s.ID,
		Title:         s.Title,
		FlowID:        s.FlowID,
		SchemaType:    string(s.Type),
		VersionID:     v.ID,
		VersionNumber: v.VersionNumber,
		Fields:        formFields,
	}, nil
}

func (uc *UseCase) CheckSchemaCompatibility(ctx context.Context, userID, projectID uuid.UUID) (bool, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.CheckSchemaCompatibility")
	defer span.End()

	// 1. Get Project's current schemas
	schemas, err := uc.deps.Schemas.List(ctx, projectID)
	if err != nil {
		return false, err
	}

	if len(schemas) == 0 {
		return true, nil
	}

	isUpToDate := true
	for _, s := range schemas {
		if s.Status != schema.StatusPublished {
			continue
		}

		if s.CurrentVersionID == nil {
			continue
		}

		// Cache check
		cacheKey := "compat:" + projectID.String() + ":" + s.CurrentVersionID.String() + ":" + userID.String()
		if val, ok := uc.deps.Cache.Get(ctx, cacheKey); ok {
			if compat, ok := val.(bool); ok {
				if !compat {
					isUpToDate = false
				}
				continue
			}
		}

		// Deep Validation
		compat, err := uc.deepValidateCompatibility(ctx, userID, projectID, s)
		if err != nil {
			return false, err
		}

		// Store in cache (1 hour TTL)
		uc.deps.Cache.Set(ctx, cacheKey, compat, time.Hour)

		if !compat {
			isUpToDate = false
		}
	}

	return isUpToDate, nil
}

func (uc *UseCase) deepValidateCompatibility(ctx context.Context, userID, projectID uuid.UUID, s schema.Schema) (bool, error) {
	u, err := uc.deps.ProjectUsers.GetByIDInternal(ctx, userID, projectID)
	if err != nil {
		return false, err
	}

	if u.Metadata == nil {
		return false, nil
	}

	var metadata map[string]any
	if err := json.Unmarshal(*u.Metadata, &metadata); err != nil {
		return false, nil
	}

	typeMap, ok := metadata[string(s.Type)].(map[string]any)
	if !ok {
		return false, nil
	}

	flowMap, ok := typeMap[s.FlowID].(map[string]any)
	if !ok {
		return false, nil
	}

	userFields, ok := flowMap["fields"].(map[string]any)
	if !ok {
		return false, nil
	}

	// Fetch current fields
	fields, err := uc.deps.Fields.GetByVersionIDWithRelations(ctx, *s.CurrentVersionID)
	if err != nil {
		return false, err
	}

	fieldDefs := make(map[string]field.Field)
	for _, f := range fields {
		fieldDefs[f.Key] = f
	}

	// Run validation logic
	_, err = uc.ValidateFields(ctx, userFields, fieldDefs, fields)
	if err != nil {
		return false, nil
	}

	return true, nil
}
