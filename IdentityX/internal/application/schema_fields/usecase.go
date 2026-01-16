package schema_fields

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

func (uc *UseCase) Create(ctx context.Context, in inbounds.SchemaFieldInput) ([]inbounds.OutputField, error) {
	var out []inbounds.OutputField
	err := uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		out, err = uc.createInternal(ctx, in)
		return err
	})

	return out, err
}

func (uc *UseCase) createInternal(ctx context.Context, in inbounds.SchemaFieldInput) ([]inbounds.OutputField, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaFieldService.Create")
	defer span.End()

	var err error
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
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
		return nil, apierr.FromService(span, inbounds.ErrNotProjectOwner{Msg: "cannot create fields for schema versions in a project you don't own"})
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
		return nil, apierr.FromService(span, inbounds.ErrSchemaNotOwned{Msg: "cannot create fields for a schema you don't own"})
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return nil, err
	}

	if latest.VersionNumber != in.VersionNumber {
		return nil, apierr.FromService(span, inbounds.ErrSchemaVersionMismatchLatest{})
	}

	if latest.Status != version.StatusDraft {
		return nil, apierr.FromService(span, inbounds.ErrAddFieldsToNonDraftVersion{})
	}

	createdFields := make([]field.Field, 0, len(in.Fields))
	for _, f := range in.Fields {
		if !field.IsValidFieldType(f.Type) {
			return nil, apierr.FromService(span, inbounds.ErrInvalidFieldType{Type: f.Type, Key: f.Key})
		}
		if !field.IsValidOwnerType(f.Owner) {
			return nil, apierr.FromService(span, inbounds.ErrInvalidFieldOwner{Key: f.Key, Owner: f.Owner})
		}
		var created *field.Field
		created, err = fields.Create(ctx, field.Field{
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
		})
		if err != nil {
			return nil, err
		}
		createdFields = append(createdFields, *created)
	}

	return inbounds.FieldSliceToOutputFieldSlice(createdFields), nil
}
