package organizations

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"
	"time"

	"github.com/MintzyG/fun"
)

func (o *Operations) CreateProjectActor(ctx context.Context, payload models.CreateActorInput) (*models.Actor, error) {
	ctx, span := telemetry.StartSpan(ctx, "CreateProjectActor")
	defer span.End()
	err := o.authz.CheckProject(ctx, *payload.ProjectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	if payload.Type != models.HumanActorType && payload.AuthMethod == models.PasswordAuthMethod {
		return nil, fun.ErrBadRequest("only human actors can authenticate via password")
	}

	actor := models.Actor{
		ProjectID:  payload.ProjectID,
		AuthMethod: payload.AuthMethod,
		VerifiedAt: new(time.Now()),
		Email:      payload.Email,
		Type:       payload.Type,
	}

	return o.actors.Register(ctx, actor)
}
