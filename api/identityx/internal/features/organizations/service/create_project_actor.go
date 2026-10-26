package service

import (
	"IdentityX/errs"
	"IdentityX/models"
	"context"
	"time"
)

func (o *Operations) CreateProjectActor(ctx context.Context, payload models.CreateActorInput) (*models.Actor, error) {
	err := o.Authz.CheckProject(ctx, *payload.ProjectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	if payload.Type != models.HumanActorType && payload.AuthMethod == models.PasswordAuthMethod {
		return nil, errs.ErrPasswordOnlyForHumans
	}

	actor := models.Actor{
		ProjectID:  payload.ProjectID,
		AuthMethod: payload.AuthMethod,
		VerifiedAt: new(time.Now()),
		Email:      payload.Email,
		Type:       payload.Type,
	}

	return o.Actors.Register(ctx, actor)
}
