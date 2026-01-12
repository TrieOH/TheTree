package schema

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/application/validation"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/version"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbound"
	"context"
	"strings"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	usecaseTracer = otel.Tracer("GoAuth.SchemaService")
)

type UseCase struct {
	schemas  outbound.SchemaRepository
	versions outbound.SchemaVersionRepository
	fields   outbound.SchemaFieldsRepository
	projects outbound.ProjectRepository
}

var _ inbounds.SchemaService = (*UseCase)(nil)

func New(
	schemas outbound.SchemaRepository,
	versions outbound.SchemaVersionRepository,
	fields outbound.SchemaFieldsRepository,
	projects outbound.ProjectRepository,
) inbounds.SchemaService {
	return &UseCase{
		schemas:  schemas,
		versions: versions,
		fields:   fields,
		projects: projects,
	}
}

func (uc *UseCase) Draft(ctx context.Context, in inbounds.SchemaServiceInput) (*inbounds.SchemaOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.Draft")
	defer span.End()

	var err error
	defer func() {
		span.SetAttributes(attribute.Bool("draft.success", err == nil))
	}()

	var principal *authz.Principal
	principal, err = auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	if in.FlowID == "" {
		err = apierr.ErrInvalidInput.WithMsg("flow id can't be empty").WithID(apierr.SchemaInvalidFlowID)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	if in.SchemaType == "" {
		err = apierr.ErrInvalidInput.WithMsg("schema type can't be empty").WithID(apierr.SchemaInvalidSchemaType)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	in.FlowID = strings.TrimSpace(strings.ToLower(in.FlowID))
	in.SchemaType = strings.TrimSpace(strings.ToLower(in.SchemaType))

	var pid uuid.UUID
	pid, err = validation.RequireProjectID(span, &in.ProjectID)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, pid, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		err = apierr.ErrUnauthorized.WithMsg("cannot draft a schema for a project you don't own").WithID(apierr.ProjectNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	if !schema.IsValidSchemaType(in.SchemaType) {
		err = apierr.ErrInvalidInput.WithMsg("invalid schema type").WithID(apierr.SchemaInvalidSchemaType)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	// FlowIDs cannot be the same as schema types so if this matches we error out
	if schema.IsValidSchemaType(in.FlowID) {
		err = apierr.ErrInvalidInput.WithMsg("flow id can't be the same as a schema type").WithID(apierr.SchemaInvalidFlowID)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	if schema.IsFlowIDReserved(in.FlowID) {
		err = apierr.ErrInvalidInput.WithMsg("flow id can't be the reserved keyword '" + string(in.FlowID) + "'").WithID(apierr.SchemaFlowIDIsReserved)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	validSchemaType := schema.Type(in.SchemaType)

	var exists bool
	exists, err = uc.schemas.Exists(ctx, schema.Schema{
		FlowID:    in.FlowID,
		ProjectID: pid,
		Type:      validSchemaType,
	})
	if err != nil {
		return nil, err
	}

	if exists {
		err = apierr.ErrConflict.WithMsg("schema with this flow ID already exists in this type").WithID(apierr.SchemaFlowIDAlreadyExistsInType)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var drafted *schema.Schema
	drafted, err = uc.schemas.Draft(ctx, schema.Schema{
		ProjectID: pid,
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

	var principal *authz.Principal
	principal, err = auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	in.SchemaID = strings.TrimSpace(strings.ToLower(in.SchemaID))

	var sid uuid.UUID
	sid, err = validation.RequireSchemaID(span, &in.SchemaID)
	if err != nil {
		return err
	}

	var pid uuid.UUID
	pid, err = validation.RequireProjectID(span, &in.ProjectID)
	if err != nil {
		return err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, pid, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema for a project you don't own").WithID(apierr.ProjectNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return err
	}

	var belongs bool
	belongs, err = uc.schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: pid,
		ID:        sid,
	})
	if err != nil {
		return err
	}

	if !belongs {
		err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema you don't own").WithID(apierr.SchemaNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return err
	}

	var toPublish *schema.Schema
	toPublish, err = uc.schemas.FindByID(ctx, sid, pid)
	if err != nil {
		return err
	}

	if toPublish.Status != schema.StatusDraft {
		if toPublish.Status == schema.StatusPublished {
			err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema that isn't a draft").WithID(apierr.SchemaTryingToPublishPublished)
			apierr.RecordDomainError(span, err)
		} else if toPublish.Status == schema.StatusArchived {
			err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema that isn't a draft").WithID(apierr.SchemaTryingToPublishArchived)
			apierr.RecordDomainError(span, err)
		} else {
			err = apierr.ErrInternal.WithMsg("CATASTROPHIC: schema found with no valid status").WithID(apierr.SchemaNoValidType)
			apierr.RecordSystemError(span, err)
		}
		return err
	}

	var latest *version.Version
	latest, err = uc.versions.GetLatest(ctx, sid)
	if err != nil && !apierr.IsNotFound(err) {
		return err
	}

	if err != nil && apierr.IsNotFound(err) {
		err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema with no versions").WithID(apierr.SchemaNoPublishedVersion)
		apierr.RecordDomainError(span, err)
		return err
	}

	if latest.VersionNumber == 1 && latest.Status == version.StatusDraft {
		err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema with only draft versions").WithID(apierr.SchemaHasOnlyDraftVersion)
		apierr.RecordDomainError(span, err)
		return err
	}

	if latest.VersionNumber == 1 && latest.Status == version.StatusArchived {
		err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema with only archived versions").WithID(apierr.SchemaHasOnlyArchivedVersion)
		apierr.RecordDomainError(span, err)
		return err
	}

	if err = uc.schemas.Publish(ctx, schema.Schema{
		ID:        sid,
		ProjectID: pid,
	}); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) GetByID(ctx context.Context, in inbounds.SchemaServiceInput) (*inbounds.SchemaOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.GetByID")
	defer span.End()

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	sid, err := validation.RequireSchemaID(span, &in.SchemaID)
	if err != nil {
		return nil, err
	}

	pid, err := validation.RequireProjectID(span, &in.ProjectID)
	if err != nil {
		return nil, err
	}

	isOwner, err := uc.projects.IsOwnerOf(ctx, pid, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		err = apierr.ErrUnauthorized.WithMsg("cannot get a schema from a project you don't own").WithID(apierr.ProjectNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var belongs bool
	belongs, err = uc.schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: pid,
		ID:        sid,
	})
	if err != nil {
		return nil, err
	}

	if !belongs {
		err = apierr.ErrUnauthorized.WithMsg("cannot get a schema you don't own").WithID(apierr.SchemaNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	found, err := uc.schemas.FindByID(ctx, sid, pid)
	if err != nil {
		return nil, err
	}

	return inbounds.SchemaToSchemaOutput(found), nil
}

func (uc *UseCase) GetVerbose(ctx context.Context, in inbounds.SchemaServiceInput) (*inbounds.SchemaVerboseOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.GetVerbose")
	defer span.End()

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	sid, err := validation.RequireSchemaID(span, &in.SchemaID)
	if err != nil {
		return nil, err
	}

	pid, err := validation.RequireProjectID(span, &in.ProjectID)
	if err != nil {
		return nil, err
	}

	isOwner, err := uc.projects.IsOwnerOf(ctx, pid, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		err = apierr.ErrUnauthorized.WithMsg("cannot get a schema from a project you don't own").WithID(apierr.ProjectNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var belongs bool
	belongs, err = uc.schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: pid,
		ID:        sid,
	})
	if err != nil {
		return nil, err
	}

	if !belongs {
		err = apierr.ErrUnauthorized.WithMsg("cannot get a schema you don't own").WithID(apierr.SchemaNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var schemaPart *schema.Schema
	schemaPart, err = uc.schemas.FindByID(ctx, sid, pid)
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
	versionsPart, err = uc.versions.List(ctx, sid)
	if err != nil {
		return nil, err
	}

	versionsOutput := make([]inbounds.VersionVerboseOutput, 0, len(versionsPart))
	for _, version := range versionsPart {
		versionOutput := inbounds.VersionVerboseOutput{
			SchemaVersionOutput: inbounds.SchemaVersionOutput{
				ID:               version.ID,
				SchemaID:         version.SchemaID,
				BasedOnVersionID: version.BasedOnVersionID,
				VersionNumber:    version.VersionNumber,
				Status:           version.Status,
				CreatedAt:        version.CreatedAt,
				UpdatedAt:        version.UpdatedAt,
			},
			Fields: nil,
		}
		versionsOutput = append(versionsOutput, versionOutput)
	}

	schemaOutput.Versions = versionsOutput

	var fieldsPart []field.Field
	fieldsPart, err = uc.fields.List(ctx, sid)
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
