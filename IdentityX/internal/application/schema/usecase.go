package schema

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/version"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbound"
	"context"
	"strings"

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
		return nil, apierr.FromService(span, err)
	}

	if in.FlowID == "" {
		return nil, apierr.FromService(span, inbounds.ErrEmptyFlowID{})
	}

	if in.SchemaType == "" {
		return nil, apierr.FromService(span, inbounds.ErrEmptySchemaType{})
	}

	in.FlowID = strings.TrimSpace(strings.ToLower(in.FlowID))
	in.SchemaType = strings.TrimSpace(strings.ToLower(in.SchemaType))

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot draft a schema for a project you don't own"})
	}

	if !schema.IsValidSchemaType(in.SchemaType) {
		return nil, apierr.FromService(span, inbounds.ErrInvalidSchemaType{})
	}

	// FlowIDs cannot be the same as schema types so if this matches we error out
	if schema.IsValidSchemaType(in.FlowID) {
		return nil, apierr.FromService(span, inbounds.ErrInvalidFlowID{Why: "flow id can't be the same as a schema type"})
	}

	if schema.IsFlowIDReserved(in.FlowID) {
		return nil, apierr.FromService(span, inbounds.ErrFlowIDIsReserved{Reserved: in.FlowID})
	}

	validSchemaType := schema.Type(in.SchemaType)

	var exists bool
	exists, err = uc.schemas.Exists(ctx, schema.Schema{
		FlowID:    in.FlowID,
		ProjectID: in.ProjectID,
		Type:      validSchemaType,
	})
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, apierr.FromService(span, inbounds.ErrFlowIDSchemaTypeConflict{})
	}

	var drafted *schema.Schema
	drafted, err = uc.schemas.Draft(ctx, schema.Schema{
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

	var principal *authz.Principal
	principal, err = auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return apierr.FromService(span, err)
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		return apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot publish a schema for a project you don't own"})
	}

	var belongs bool
	belongs, err = uc.schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return err
	}

	if !belongs {
		return apierr.FromService(span, inbounds.ErrSchemaNotOwned{Msg: "cannot publish a schema you don't own"})
	}

	var toPublish *schema.Schema
	toPublish, err = uc.schemas.FindByID(ctx, in.SchemaID, in.ProjectID)
	if err != nil {
		return err
	}

	if toPublish.Status != schema.StatusDraft {
		if toPublish.Status == schema.StatusPublished {
			err = apierr.FromService(span, inbounds.ErrPublishSchemaPublished{})
		} else if toPublish.Status == schema.StatusArchived {
			err = apierr.FromService(span, inbounds.ErrPublishSchemaArchived{})
		} else {
			err = apierr.FromService(span, inbounds.ErrSchemaInvalidStatus{Status: string(toPublish.Status)})
		}
		return err
	}

	var latest *version.Version
	latest, err = uc.versions.GetLatest(ctx, in.SchemaID)
	if err != nil && !apierr.IsNotFound(err) {
		return err
	}

	if err != nil && apierr.IsNotFound(err) {
		return apierr.FromService(span, inbounds.ErrSchemaNoPublishedVersions{Msg: "cannot publish a schema with no versions"})
	}

	if latest.VersionNumber == 1 && latest.Status == version.StatusDraft {
		return apierr.FromService(span, inbounds.ErrSchemaOnlyDraft{Msg: "cannot publish a schema with only draft versions"})
	}

	if latest.VersionNumber == 1 && latest.Status == version.StatusArchived {
		return apierr.FromService(span, inbounds.ErrSchemaOnlyArchived{Msg: "cannot publish a schema with only archived versions"})
	}

	if err = uc.schemas.Publish(ctx, schema.Schema{
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

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	isOwner, err := uc.projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot get a schema from a project you don't own"})
	}

	var belongs bool
	belongs, err = uc.schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return nil, err
	}

	if !belongs {
		return nil, apierr.FromService(span, inbounds.ErrSchemaNotOwned{Msg: "cannot get a schema you don't own"})
	}

	found, err := uc.schemas.FindByID(ctx, in.SchemaID, in.ProjectID)
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
		return nil, apierr.FromService(span, err)
	}

	isOwner, err := uc.projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot get a schema from a project you don't own"})
	}

	var belongs bool
	belongs, err = uc.schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return nil, err
	}

	if !belongs {
		return nil, apierr.FromService(span, inbounds.ErrSchemaNoPublishedVersions{Msg: "cannot get a schema you don't own"})
	}

	var schemaPart *schema.Schema
	schemaPart, err = uc.schemas.FindByID(ctx, in.SchemaID, in.ProjectID)
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
	versionsPart, err = uc.versions.List(ctx, in.SchemaID)
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
	fieldsPart, err = uc.fields.List(ctx, in.SchemaID)
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
