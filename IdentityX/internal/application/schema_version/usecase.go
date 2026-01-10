package schema_version

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/authz"
	"GoAuth/internal/application/transactions"
	"GoAuth/internal/application/validation"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbound"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	usecaseTracer = otel.Tracer("GoAuth.SchemaVersionService")
)

type UseCase struct {
	schemas  outbound.SchemaRepository
	versions outbound.SchemaVersionRepository
	fields   outbound.SchemaFieldsRepository
	projects outbound.ProjectRepository
	tx       transactions.TxRunner
}

var _ inbounds.SchemaVersionService = (*UseCase)(nil)

func New(
	schemas outbound.SchemaRepository,
	versions outbound.SchemaVersionRepository,
	fields outbound.SchemaFieldsRepository,
	projects outbound.ProjectRepository,
	tx transactions.TxRunner,
) inbounds.SchemaVersionService {
	return &UseCase{
		schemas:  schemas,
		versions: versions,
		fields:   fields,
		projects: projects,
		tx:       tx,
	}
}

func (uc *UseCase) Draft(ctx context.Context, in inbounds.SchemaVersionServiceInput) (*inbounds.SchemaVersionOutput, error) {
	var out *inbounds.SchemaVersionOutput
	err := uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		out, err = uc.draftInternal(ctx, in)
		return err
	})

	return out, err
}

func (uc *UseCase) draftInternal(ctx context.Context, in inbounds.SchemaVersionServiceInput) (*inbounds.SchemaVersionOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaVersionService.Draft")
	defer span.End()

	var err error
	defer func() {
		span.SetAttributes(attribute.Bool("draft.success", err == nil))
	}()

	var principal *authz.Principal
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var pid *uuid.UUID
	pid, err = validation.RequireProjectID(span, &in.ProjectID)
	if err != nil {
		return nil, err
	}

	var sid *uuid.UUID
	sid, err = validation.RequireSchemaID(span, &in.SchemaID)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *pid, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		err = apierr.ErrUnauthorized.WithMsg("cannot draft a schema version for a project you don't own").WithID(apierr.ProjectNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var belongs bool
	belongs, err = uc.schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: *pid,
		ID:        *sid,
	})
	if err != nil {
		return nil, err
	}

	if !belongs {
		err = apierr.ErrUnauthorized.WithMsg("cannot draft a schema version for a schema you don't own").WithID(apierr.SchemaNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var latest *schema.Version
	latest, err = uc.versions.GetLatestForUpdate(ctx, *sid)

	if err != nil && !apierr.IsNotFound(err) {
		return nil, err
	}

	if apierr.IsNotFound(err) {
		newVersion := &schema.Version{
			SchemaID:      *sid,
			VersionNumber: 1,
		}

		newVersion, err = uc.versions.Draft(ctx, *newVersion)
		if err != nil {
			return nil, err
		}

		if err = uc.schemas.SetVersion(ctx, schema.Schema{
			ID:               *sid,
			ProjectID:        *pid,
			CurrentVersionID: &newVersion.ID,
		}); err != nil {
			return nil, err
		}

		return inbounds.SchemaVersionToOutput(newVersion), nil
	}

	if latest.Status != schema.VersionStatusPublished {
		err = apierr.ErrUnauthorized.WithMsg("new versions can only be drafted from published versions").WithID(apierr.SchemaVersionDraftOnNonPublished)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var newVersionDraft *schema.Version
	newVersionDraft, err = uc.versions.CopyOnDraft(ctx, latest.ID)
	if err != nil {
		return nil, err
	}

	err = uc.fields.CloneFromTo(ctx, latest.ID, newVersionDraft.ID)
	if err != nil {
		return nil, err
	}

	if err = uc.schemas.SetVersion(ctx, schema.Schema{
		ID:               *sid,
		ProjectID:        *pid,
		CurrentVersionID: &newVersionDraft.ID,
	}); err != nil {
		return nil, err
	}

	return inbounds.SchemaVersionToOutput(newVersionDraft), nil
}

func (uc *UseCase) Publish(ctx context.Context, in inbounds.SchemaVersionServiceInput) error {
	ctx, span := usecaseTracer.Start(ctx, "SchemaVersionService.Publish")
	defer span.End()

	var err error
	defer func() {
		span.SetAttributes(attribute.Bool("publish.success", err == nil))
	}()

	var principal *authz.Principal
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	var pid *uuid.UUID
	pid, err = validation.RequireProjectID(span, &in.ProjectID)
	if err != nil {
		return err
	}

	var sid *uuid.UUID
	sid, err = validation.RequireSchemaID(span, &in.SchemaID)
	if err != nil {
		return err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, *pid, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema version for a project you don't own").WithID(apierr.ProjectNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return err
	}

	var belongs bool
	belongs, err = uc.schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: *pid,
		ID:        *sid,
	})
	if err != nil {
		return err
	}

	if !belongs {
		err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema version for a schema you don't own").WithID(apierr.SchemaNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return err
	}

	var latest *schema.Version
	latest, err = uc.versions.GetLatest(ctx, *sid)
	if err != nil && !apierr.IsNotFound(err) {
		return err
	}

	if err != nil && apierr.IsNotFound(err) {
		err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema version draft that doesn't exist").WithID(apierr.SchemaVersionDraftDoesntExist)
		apierr.RecordDomainError(span, err)
		return err
	}

	if latest.Status != schema.VersionStatusDraft {
		if latest.Status == schema.VersionStatusPublished {
			err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema version that isn't a draft").WithID(apierr.SchemaVersionTryingToPublishPublished)
			apierr.RecordDomainError(span, err)
		} else if latest.Status == schema.VersionStatusArchived {
			err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema version that isn't a draft").WithID(apierr.SchemaVersionTryingToPublishArchived)
			apierr.RecordDomainError(span, err)
		} else {
			err = apierr.ErrInternal.WithMsg("CATASTROPHIC: schema version found with no valid status").WithID(apierr.SchemaVersionNoValidType)
			apierr.RecordSystemError(span, err)
		}
		return err
	}

	if latest.BasedOnVersionID == nil {
		if err = uc.validateVersionHasFields(ctx, span, latest.ID); err != nil {
			return err
		}

		if err = uc.versions.Publish(ctx, schema.Version{
			SchemaID: *sid,
			ID:       latest.ID,
		}); err != nil {
			return err
		}

		return nil
	}

	var hasChanges bool
	hasChanges, err = uc.fields.DiffVersionsState(ctx, *latest.BasedOnVersionID, latest.ID)
	if err != nil {
		return err
	}

	if !hasChanges {
		err = apierr.ErrInvalidInput.WithMsg("cannot publish a version with no changes").WithID(apierr.SchemaVersionNoChanges)
		apierr.RecordDomainError(span, err)
		return err
	}

	if err = uc.validateVersionHasFields(ctx, span, latest.ID); err != nil {
		return err
	}

	if err = uc.versions.Publish(ctx, schema.Version{
		SchemaID: *sid,
		ID:       latest.ID,
	}); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) validateVersionHasFields(ctx context.Context, span trace.Span, versionID uuid.UUID) error {
	fields, err := uc.fields.GetByVersionID(ctx, versionID)
	if err != nil && !apierr.IsNotFound(err) {
		return err
	}

	if apierr.IsNotFound(err) || len(fields) == 0 {
		err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema version with no fields").WithID(apierr.SchemaVersionPublishWithNoFields)
		apierr.RecordDomainError(span, err)
		return err
	}
	return nil
}
