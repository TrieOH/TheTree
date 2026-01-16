package schema_version

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/version"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
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
	deps Deps
	tx   inbounds.TxRunner
}

type Deps struct {
	Schemas  outbounds.SchemaRepository
	Versions outbounds.SchemaVersionRepository
	Fields   outbounds.SchemaFieldsRepository
	Projects outbounds.ProjectRepository
}

var _ inbounds.SchemaVersionService = (*UseCase)(nil)

func New(
	deps Deps,
	tx inbounds.TxRunner,
) inbounds.SchemaVersionService {
	return &UseCase{
		deps: deps,
		tx:   tx,
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

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	var principal *authz.Principal
	principal, err = auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot draft a schema version for a project you don't own"})
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
		return nil, apierr.FromService(span, inbounds.ErrSchemaNotOwned{Msg: "cannot draft a schema version for a schema you don't own"})
	}

	var latest *version.Version
	latest, err = versions.GetLatestForUpdate(ctx, in.SchemaID)

	if err != nil && !apierr.IsNotFound(err) {
		return nil, err
	}

	if apierr.IsNotFound(err) {
		newVersion := &version.Version{
			SchemaID:      in.SchemaID,
			VersionNumber: 1,
		}

		newVersion, err = versions.Draft(ctx, *newVersion)
		if err != nil {
			return nil, err
		}

		if err = schemas.SetVersion(ctx, schema.Schema{
			ID:               in.SchemaID,
			ProjectID:        in.ProjectID,
			CurrentVersionID: &newVersion.ID,
		}); err != nil {
			return nil, err
		}

		return inbounds.SchemaVersionToOutput(newVersion), nil
	}

	if latest.Status != version.StatusPublished {
		return nil, apierr.FromService(span, inbounds.ErrDraftVersionOnNonPublished{})
	}

	var newVersionDraft *version.Version
	newVersionDraft, err = versions.CopyOnDraft(ctx, latest.ID)
	if err != nil {
		return nil, err
	}

	err = fields.CloneFromTo(ctx, latest.ID, newVersionDraft.ID)
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

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

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
		return apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot publish a schema version for a project you don't own"})
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
		return apierr.FromService(span, inbounds.ErrSchemaNotOwned{Msg: "cannot publish a schema version for a schema you don't own"})
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil && !apierr.IsNotFound(err) {
		return err
	}

	if err != nil && apierr.IsNotFound(err) {
		return apierr.FromService(span, inbounds.ErrPublishNonExistentVersionDraft{})
	}

	if latest.Status != version.StatusDraft {
		if latest.Status == version.StatusPublished {
			err = apierr.FromService(span, inbounds.ErrPublishVersionPublished{})
		} else if latest.Status == version.StatusArchived {
			err = apierr.FromService(span, inbounds.ErrPublishVersionArchived{})
		} else {
			err = apierr.FromService(span, inbounds.ErrPublishVersionInvalidStatus{})
		}
		return err
	}

	if latest.BasedOnVersionID == nil {
		if err = uc.validateVersionHasFields(ctx, span, latest.ID); err != nil {
			return err
		}

		if err = versions.Publish(ctx, version.Version{
			SchemaID: in.SchemaID,
			ID:       latest.ID,
		}); err != nil {
			return err
		}

		return nil
	}

	var hasChanges bool
	hasChanges, err = fields.DiffVersionsState(ctx, *latest.BasedOnVersionID, latest.ID)
	if err != nil {
		return err
	}

	if !hasChanges {
		return apierr.FromService(span, inbounds.ErrPublishVersionNoChanges{})
	}

	if err = uc.validateVersionHasFields(ctx, span, latest.ID); err != nil {
		return err
	}

	if err = versions.Publish(ctx, version.Version{
		SchemaID: in.SchemaID,
		ID:       latest.ID,
	}); err != nil {
		return err
	}

	if err := schemas.SetVersion(ctx, schema.Schema{
		ID:               in.SchemaID,
		ProjectID:        in.ProjectID,
		CurrentVersionID: &latest.ID,
	}); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) validateVersionHasFields(ctx context.Context, span trace.Span, versionID uuid.UUID) error {
	fields := uc.deps.Fields
	foundFields, err := fields.GetByVersionID(ctx, versionID)
	if err != nil && !apierr.IsNotFound(err) {
		return err
	}

	if apierr.IsNotFound(err) || len(foundFields) == 0 {
		return apierr.FromService(span, inbounds.ErrPublishVersionNoFields{})
	}
	return nil
}
