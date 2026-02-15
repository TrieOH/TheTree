package scope

import (
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/scopes"
	"GoAuth/internal/errx"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"

	"github.com/MintzyG/fail/v3"
	"go.opentelemetry.io/otel"
)

var (
	usecaseTracer = otel.Tracer("GoAuth.ScopeService")
)

type UseCase struct {
	projects outbounds.ProjectRepository
	scopes   outbounds.ScopeRepository
	tx       inbounds.TxRunner
}

var _ inbounds.ScopeService = (*UseCase)(nil)

func New(
	projects outbounds.ProjectRepository,
	scopes outbounds.ScopeRepository,
	tx inbounds.TxRunner,
) inbounds.ScopeService {
	return &UseCase{
		projects: projects,
		scopes:   scopes,
		tx:       tx,
	}
}

func (uc *UseCase) Create(ctx context.Context, in inbounds.CreateScopeInput) (*inbounds.ScopeOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "ScopeService.Create")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot get scopes for a project you don't own").RecordCtx(ctx)
	}

	if in.Name == "" {
		return nil, fail.New(errx.SCOPEEmptyName).RecordCtx(ctx)
	}

	scope, err := uc.scopes.Create(ctx, scopes.Scope{
		Type:       scopes.ScopeTypeProjectScope,
		ProjectID:  &in.ProjectID,
		Name:       &in.Name,
		ExternalID: in.ExternalID,
	})
	if err != nil {
		return nil, err
	}

	return inbounds.ScopeToScopeOutput(scope), nil
}

func (uc *UseCase) GetByIDExternal(ctx context.Context, in inbounds.GetScopeInput) (*inbounds.ScopeOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "ScopeService.GetByIDExternal")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot get a scope for a project you don't own").RecordCtx(ctx)
	}

	scope, err := uc.scopes.GetByIDExternal(ctx, in.ScopeID, in.ProjectID)
	if err != nil {
		return nil, err
	}

	return inbounds.ScopeToScopeOutput(scope), nil
}

func (uc *UseCase) GetProjectScopesExternal(ctx context.Context, in inbounds.GetScopeInput) ([]inbounds.ScopeOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "ScopeService.GetProjectScopesExternal")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var isOwner bool
	isOwner, err = uc.projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(errx.ProjectNotOwnedByPrincipal).WithArgs("cannot get scopes for a project you don't own").RecordCtx(ctx)
	}

	projectScopes, err := uc.scopes.GetProjectScopes(ctx, in.ProjectID)
	if err != nil {
		return nil, err
	}

	return inbounds.ScopeSliceToScopeSliceOutput(projectScopes), nil
}
