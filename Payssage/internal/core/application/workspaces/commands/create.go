package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"
	"encoding/json"

	"github.com/TrieOH/goauth-sdk-go"
)

func (uc *CommandService) Create(ctx context.Context, name string) (ws *domain.Workspace, err error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.Create")
	defer span.End()

	ga := uc.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var workspace *domain.Workspace
	workspace, err = domain.NewWorkspace(sub.ID, name)
	if err != nil {
		return nil, err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("workspaces").
		Action("create").
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("workspace").SetMessage("insufficient permissions")
	}

	meta := json.RawMessage(`{"color": "#6a07e3", "icon": "Shield"}`)
	var scope *goauth.Scope
	var idStr = workspace.ID.String()
	scope, err = ga.Scopes.CreateWithParent(ctx, workspace.Name, &idStr, nil, meta)
	if err != nil {
		return nil, err
	}
	workspace.AddScope(scope.ID)

	var created *domain.Workspace
	created, err = uc.workspaces.Create(ctx, *workspace)
	if err != nil {
		return nil, err
	}

	return created, nil
}
