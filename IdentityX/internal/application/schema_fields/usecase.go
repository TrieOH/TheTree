package schema_fields

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/application/transactions"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/version"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbound"
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	usecaseTracer = otel.Tracer("GoAuth.SchemaFieldsService")
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
		err = apierr.ErrUnauthorized.WithMsg("cannot create fields for schema versions in a project you don't own").WithID(apierr.ProjectNotOwnedByPrincipal)
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
		err = apierr.ErrUnauthorized.WithMsg("cannot create fields for a schema you don't own").WithID(apierr.SchemaNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var latest *version.Version
	latest, err = uc.versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return nil, err
	}

	if latest.VersionNumber != in.VersionNumber {
		err = apierr.ErrInvalidInput.WithMsg("version number does not match latest version").WithID(apierr.SchemaVersionMismatch)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	if latest.Status != version.StatusDraft {
		err = apierr.ErrConflict.WithMsg("cannot add fields to a non-draft version").WithID(apierr.SchemaVersionNotDraft)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	createdFields := make([]field.Field, 0, len(in.Fields))
	for _, f := range in.Fields {
		if !field.IsValidFieldType(f.Type) {
			err = apierr.ErrInvalidInput.WithMsg("invalid field type (" + f.Type + ") for field: " + f.Key).WithID(apierr.FieldInvalidType)
			apierr.RecordDomainError(span, err)
			return nil, err
		}
		if !field.IsValidOwnerType(f.Owner) {
			err = apierr.ErrInvalidInput.WithMsg("invalid owner type (" + f.Owner + ") for field: " + f.Key).WithID(apierr.FieldInvalidOwner)
			apierr.RecordDomainError(span, err)
			return nil, err
		}
		var created *field.Field
		created, err = uc.fields.Create(ctx, field.Field{
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
