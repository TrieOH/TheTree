package schema_version

import (
	"GoAuth/internal/adapters/observability/tracing"
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/authz"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbound"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	usecaseTracer = otel.Tracer("GoAuth.SchemaVersionService")
)

type UseCase struct {
	schemas  outbound.SchemaRepository
	versions outbound.SchemaVersionRepository
	projects outbound.ProjectRepository
}

var _ inbounds.SchemaVersionService = (*UseCase)(nil)

func New(
	schemas outbound.SchemaRepository,
	versions outbound.SchemaVersionRepository,
	projects outbound.ProjectRepository,
) inbounds.SchemaVersionService {
	return &UseCase{
		schemas:  schemas,
		versions: versions,
		projects: projects,
	}
}

func (uc *UseCase) Draft(ctx context.Context, in inbounds.DraftSchemaVersionInput) (*inbounds.DraftSchemaVersionOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaVersionService.Draft")
	defer span.End()

	var err error
	defer func() {
		span.SetAttributes(attribute.Bool("draft.success", err == nil))
	}()

	var principal *authz.Principal
	principal, err = authz.RequirePrincipal(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	tracing.AnnotatePrincipal(span, principal)

	var pid uuid.UUID
	pid, err = uuid.Parse(in.ProjectID)
	if err != nil {
		err = apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var sid uuid.UUID
	sid, err = uuid.Parse(in.SchemaID)
	if err != nil {
		err = apierr.ErrInvalidInput.WithMsg("invalid schema id").WithID(apierr.SchemaInvalidID).WithCause(err)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, pid, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		err = apierr.ErrUnauthorized.WithMsg("cannot draft a schema version for a project you dont own").WithID(apierr.ProjectNotOwnedByPrincipal)
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
		err = apierr.ErrUnauthorized.WithMsg("cannot draft a schema version for a schema you dont own").WithID(apierr.SchemaNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var latest *schema.Version
	latest, err = uc.versions.GetLatest(ctx, sid)

	if err != nil && !apierr.IsNotFound(err) {
		return nil, err
	}

	if err != nil && apierr.IsNotFound(err) {
		firstVersion := schema.Version{
			SchemaID:      sid,
			VersionNumber: 1,
		}

		var newVersion *schema.Version
		newVersion, err = uc.versions.Draft(ctx, firstVersion)
		if err != nil {
			return nil, err
		}

		if err = uc.schemas.SetVersion(ctx, schema.Schema{
			ID:               sid,
			ProjectID:        pid,
			CurrentVersionID: &newVersion.ID,
		}); err != nil {
			return nil, err
		}

		return inbounds.SchemaVersionToOutput(newVersion), nil
	}

	var newVersion *schema.Version
	newVersion, err = uc.versions.Draft(ctx, schema.Version{
		SchemaID:      sid,
		VersionNumber: latest.VersionNumber + 1,
	})
	if err != nil {
		return nil, err
	}

	if err = uc.schemas.SetVersion(ctx, schema.Schema{
		ID:               sid,
		ProjectID:        pid,
		CurrentVersionID: &newVersion.ID,
	}); err != nil {
		return nil, err
	}

	return inbounds.SchemaVersionToOutput(newVersion), nil
}

func (uc *UseCase) Publish(ctx context.Context, in inbounds.PublishSchemaVersionInput) error {
	return nil
}
