package commands

import (
	"IdentityX/models"
	"context"
	"lib/authz"

	"github.com/google/uuid"
)

func (s *CommandService) RevokeAPIKey(ctx context.Context, keyID uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "ApiKeys.Revoke")
	defer span.End()

	sub, err := models.RequireSubject(ctx)
	if err != nil {
		return err
	}

	if err = s.perms.Require(ctx,
		authz.Subject("user", sub.ID),
		authz.Permission("revoke"),
		authz.Resource("api_key", keyID.String()),
		nil,
	); err != nil {
		return err
	}

	if _, err := s.apiKeys.Revoke(ctx, keyID, sub.ID); err != nil {
		return err
	}

	if err = s.perms.DeleteRelation(ctx,
		"api_key:"+keyID.String()+"#parent_user@user:"+sub.ID.String(),
	); err != nil {
		return err
	}

	return nil
}
