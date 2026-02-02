package schema

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
	"strings"

	"github.com/MintzyG/fail"
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
	Schemas  outbounds.SchemaRepository
	Versions outbounds.SchemaVersionRepository
	Fields   outbounds.SchemaFieldsRepository
	Projects outbounds.ProjectRepository
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
	principal, err = auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	if in.FlowID == "" {
		return nil, fail.New(apierr.SchemaEmptyFlowID)
	}

	if in.SchemaType == "" {
		return nil, fail.New(apierr.SchemaEmptySchemaType)
	}

	in.FlowID = strings.TrimSpace(strings.ToLower(in.FlowID))
	in.SchemaType = strings.TrimSpace(strings.ToLower(in.SchemaType))

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(apierr.ProjectNotOwnedByPrincipal).WithArgs("cannot draft a schema for a project you don't own")
	}

	if !schema.IsValidSchemaType(in.SchemaType) {
		return nil, fail.New(apierr.SchemaInvalidSchemaType)
	}

	// FlowIDs cannot be the same as schema types so if this matches we error out
	if schema.IsValidSchemaType(in.FlowID) {
		return nil, fail.New(apierr.SchemaInvalidFlowID).WithArgs("flow id can't be the same as a schema type")
	}

	if schema.IsFlowIDReserved(in.FlowID) {
		return nil, fail.New(apierr.SchemaFlowIDIsReserved).WithArgs(in.FlowID)
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
		return nil, fail.New(apierr.SchemaFlowIDAlreadyExistsInType)
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
	principal, err = auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return apierr.FromService(span, err)
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		return fail.New(apierr.ProjectNotOwnedByPrincipal).WithArgs("cannot publish a schema for a project you don't own")
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
		return fail.New(apierr.SchemaNotOwnedByPrincipal).WithArgs("cannot publish a schema you don't own")
	}

	var toPublish *schema.Schema
	toPublish, err = schemas.FindByID(ctx, in.SchemaID, in.ProjectID)
	if err != nil {
		return err
	}

	if toPublish.Status != schema.StatusDraft {
		if toPublish.Status == schema.StatusPublished {
			err = fail.New(apierr.SchemaTryingToPublishPublished)
		} else if toPublish.Status == schema.StatusArchived {
			err = fail.New(apierr.SchemaTryingToPublishArchived)
		} else {
			err = fail.New(apierr.SchemaNoValidStatus).WithArgs(toPublish.Status)
		}
		return err
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil && !fail.Is(err, apierr.SQLNotFound) {
		return err
	}

	if err != nil && fail.Is(err, apierr.SQLNotFound) {
		return fail.New(apierr.SCHEMANoPublishedVersion)
	}

	if latest.VersionNumber == 1 && latest.Status == version.StatusDraft {
		return fail.New(apierr.SchemaHasOnlyDraftVersion)
	}

	if latest.VersionNumber == 1 && latest.Status == version.StatusArchived {
		return fail.New(apierr.SchemaHasOnlyArchivedVersion)
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

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	isOwner, err := projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(apierr.ProjectNotOwnedByPrincipal).WithArgs("cannot get a schema from a project you don't own")
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
		return nil, fail.New(apierr.SchemaNotOwnedByPrincipal).WithArgs("cannot get a schema you don't own")
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

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	isOwner, err := projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(apierr.ProjectNotOwnedByPrincipal).WithArgs("cannot get a schema from a project you don't own")
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
		return nil, fail.New(apierr.SchemaNotOwnedByPrincipal).WithArgs("cannot get a schema you don't own")
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

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	isOwner, err := projects.IsOwnerOf(ctx, projectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(apierr.ProjectNotOwnedByPrincipal).WithArgs("cannot get schema IDs from a project you don't own")
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

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	isOwner, err := projects.IsOwnerOf(ctx, projectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(apierr.ProjectNotOwnedByPrincipal).WithArgs("cannot get schemas from a project you don't own")
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
	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	isOwner, err := projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, fail.New(apierr.ProjectNotOwnedByPrincipal).WithArgs("cannot get form for a project you don't own")
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
			return nil, fail.New(apierr.SchemaNotOwnedByPrincipal).WithArgs("cannot get form for a schema you don't own")
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
			return nil, fail.New(apierr.SCHEMANoPublishedVersion)
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

	// Ensure version is published (for forms we don't want to expose drafts)
	if v.Status != version.StatusPublished {
		return nil, apierr.FromService(span, inbounds.ErrVersionNotPublished{})
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
