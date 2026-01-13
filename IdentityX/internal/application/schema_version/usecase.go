package schema_version

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/application/transactions"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/version"
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
	principal, err = auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
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
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return nil, err
	}

	if !belongs {
		err = apierr.ErrUnauthorized.WithMsg("cannot draft a schema version for a schema you don't own").WithID(apierr.SchemaNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var latest *version.Version
	latest, err = uc.versions.GetLatestForUpdate(ctx, in.SchemaID)

	if err != nil && !apierr.IsNotFound(err) {
		return nil, err
	}

	if apierr.IsNotFound(err) {
		newVersion := &version.Version{
			SchemaID:      in.SchemaID,
			VersionNumber: 1,
		}

		newVersion, err = uc.versions.Draft(ctx, *newVersion)
		if err != nil {
			return nil, err
		}

		if err = uc.schemas.SetVersion(ctx, schema.Schema{
			ID:               in.SchemaID,
			ProjectID:        in.ProjectID,
			CurrentVersionID: &newVersion.ID,
		}); err != nil {
			return nil, err
		}

		return inbounds.SchemaVersionToOutput(newVersion), nil
	}

	if latest.Status != version.StatusPublished {
		err = apierr.ErrBadRequest.WithMsg("new versions can only be drafted from published versions").WithID(apierr.SchemaVersionDraftOnNonPublished)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var newVersionDraft *version.Version
	newVersionDraft, err = uc.versions.CopyOnDraft(ctx, latest.ID)
	if err != nil {
		return nil, err
	}

	err = uc.fields.CloneFromTo(ctx, latest.ID, newVersionDraft.ID)
	if err != nil {
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
		err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema version for a project you don't own").WithID(apierr.ProjectNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return err
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
		err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema version for a schema you don't own").WithID(apierr.SchemaNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return err
	}

	var latest *version.Version
	latest, err = uc.versions.GetLatest(ctx, in.SchemaID)
	if err != nil && !apierr.IsNotFound(err) {
		return err
	}

	if err != nil && apierr.IsNotFound(err) {
		err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema version draft that doesn't exist").WithID(apierr.SchemaVersionDraftDoesntExist)
		apierr.RecordDomainError(span, err)
		return err
	}

	if latest.Status != version.StatusDraft {
		if latest.Status == version.StatusPublished {
			err = apierr.ErrUnauthorized.WithMsg("cannot publish a schema version that isn't a draft").WithID(apierr.SchemaVersionTryingToPublishPublished)
			apierr.RecordDomainError(span, err)
		} else if latest.Status == version.StatusArchived {
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

		if err = uc.versions.Publish(ctx, version.Version{
			SchemaID: in.SchemaID,
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

	if err = uc.versions.Publish(ctx, version.Version{
		SchemaID: in.SchemaID,
		ID:       latest.ID,
	}); err != nil {
		return err
	}

	if err := uc.schemas.SetVersion(ctx, schema.Schema{
		ID:               in.SchemaID,
		ProjectID:        in.ProjectID,
		CurrentVersionID: &latest.ID,
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
