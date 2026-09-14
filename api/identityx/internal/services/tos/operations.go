// Package tos owns the project terms-of-service concern: versioned
// documents, the consent ledger, registration gating, and the
// change-notification flow.
package tos

import (
	"context"
	"time"

	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/errx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// noticePeriod is how long existing users are given between a ToS
// announcement and its effective date. 30 days is the notice period
// customary in the Brazilian consumer legal framework (CDC) and the safe
// harbour for the LGPD art. 9 §6 notification duty.
const noticePeriod = 30 * 24 * time.Hour

type Operations struct {
	tos      ports.TosRepo
	projects ports.ProjectRepo
	actors   ports.ActorRepo
	authz    *authz.Service
	notifier ports.TosNotifier
}

func NewOperations(
	tos ports.TosRepo,
	projects ports.ProjectRepo,
	actors ports.ActorRepo,
	authzSvc *authz.Service,
	notifier ports.TosNotifier,
) *Operations {
	return errx.MustProvide(&Operations{
		tos:      tos,
		projects: projects,
		actors:   actors,
		authz:    authzSvc,
		notifier: notifier,
	})
}

// authorizeAdmin gates a write on the authenticated actor holding an
// admin or owner role on the project; unknown projects surface as 404
// (CheckProject passes the project lookup through).
func (o *Operations) authorizeAdmin(ctx context.Context, projectID uuid.UUID) error {
	return o.authz.CheckProject(ctx, projectID, models.ProjectRoleAdmin)
}

// requireActorInProject denies access to actors outside the project scope:
// project-scoped actors must belong to the given project, platform actors
// have no project and are never reachable through project routes.
func (o *Operations) requireActorInProject(ctx context.Context, actorID, projectID uuid.UUID) error {
	actor, err := o.actors.GetByID(ctx, actorID)
	if err != nil {
		return err
	}
	if actor.ProjectID == nil || *actor.ProjectID != projectID {
		return fun.ErrForbidden("actor does not belong to this project")
	}
	return nil
}
