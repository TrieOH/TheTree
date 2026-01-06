package schema

import (
	"GoAuth/internal/adapters/observability/tracing"
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/authz"
	"GoAuth/internal/domain/schema"
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
	projects outbound.ProjectRepository
}

var _ inbounds.SchemaService = (*UseCase)(nil)

func New(
	schemas outbound.SchemaRepository,
	projects outbound.ProjectRepository,
) inbounds.SchemaService {
	return &UseCase{
		schemas:  schemas,
		projects: projects,
	}
}

func (uc *UseCase) Draft(ctx context.Context, in inbounds.DraftSchemaInput) (*inbounds.SchemaOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.Draft")
	defer span.End()

	var err error
	defer func() {
		span.SetAttributes(attribute.Bool("draft.success", err == nil))
	}()

	in.FlowID = strings.TrimSpace(strings.ToLower(in.FlowID))

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

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, pid, principal.UserID)
	if err != nil {
		err = apierr.ErrUnauthorized.WithMsg("error checking project ownership").WithID(apierr.ProjectOwnershipCheckFailed).WithCause(err)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	if !isOwner {
		err = apierr.ErrUnauthorized.WithMsg("cannot draft a schema for a project you dont own").WithID(apierr.ProjectNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	isValidType := schema.IsValidSchemaType(in.SchemaType)
	if !isValidType {
		err = apierr.ErrInvalidInput.WithMsg("invalid schema type").WithID(apierr.SchemaInvalidSchemaType)
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

func (uc *UseCase) GetByID(ctx context.Context, in inbounds.GetSchemaByIDInput) (*inbounds.SchemaOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaService.GetByID")
	defer span.End()

	principal, err := authz.RequirePrincipal(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	tracing.AnnotatePrincipal(span, principal)

	sid, err := uuid.Parse(in.SchemaID)
	if err != nil {
		err = apierr.ErrInvalidInput.WithMsg("invalid schema id").WithID(apierr.SchemaInvalidID).WithCause(err)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	pid, err := uuid.Parse(in.ProjectID)
	if err != nil {
		err = apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	isOwner, err := uc.projects.IsOwnerOf(ctx, pid, principal.UserID)
	if err != nil {
		err = apierr.ErrUnauthorized.WithMsg("error checking project ownership").WithID(apierr.ProjectOwnershipCheckFailed).WithCause(err)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	if !isOwner {
		err = apierr.ErrUnauthorized.WithMsg("cannot get a schema from a project you dont own").WithID(apierr.ProjectNotOwnedByPrincipal)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	found, err := uc.schemas.FindByID(ctx, sid, pid)
	if err != nil {
		return nil, err
	}

	return inbounds.SchemaToSchemaOutput(found), nil
}
