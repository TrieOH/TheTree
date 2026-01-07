package schema_fields

import (
	"GoAuth/internal/adapters/observability/tracing"
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/authz"
	"GoAuth/internal/application/transactions"
	"GoAuth/internal/domain/field"
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
	fields   outbound.SchemaFieldsRepository
	projects outbound.ProjectRepository
	tx       transactions.TxRunner
}

var _ inbounds.SchemaFieldsService = (*UseCase)(nil)

func New(
	schemas outbound.SchemaRepository,
	versions outbound.SchemaVersionRepository,
	fields outbound.SchemaFieldsRepository,
	projects outbound.ProjectRepository,
	tx transactions.TxRunner,
) inbounds.SchemaFieldsService {
	return &UseCase{
		schemas:  schemas,
		versions: versions,
		fields:   fields,
		projects: projects,
		tx:       tx,
	}
}

func (uc *UseCase) Create(ctx context.Context, in inbounds.CreateSchemaFieldInput) ([]inbounds.OutputField, error) {
	var out []inbounds.OutputField
	err := uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		out, err = uc.createInternal(ctx, in)
		return err
	})

	return out, err
}

func (uc *UseCase) createInternal(ctx context.Context, in inbounds.CreateSchemaFieldInput) ([]inbounds.OutputField, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaFieldService.Create")
	defer span.End()

	var err error
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
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
	if err != nil {
		return nil, err
	}

	createdFields := make([]field.Field, 0, len(in.Fields))
	for _, f := range in.Fields {
		var created *field.Field
		created, err = uc.fields.Create(ctx, field.Field{
			ID:              uuid.New(),
			SchemaID:        sid,
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
		})
		if err != nil {
			return nil, err
		}
		createdFields = append(createdFields, *created)
	}

	return inbounds.FieldSliceToOutputFieldSlice(createdFields), nil
}
