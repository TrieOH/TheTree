package service

import (
	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"
	"context"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

type Deps struct {
	Templates ports.EmailTemplateRepo
	Authz     *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

// authorizeAdmin gates a write on the authenticated actor holding an
// admin or owner role on the project; unknown projects surface as 404
// (CheckProject passes the project lookup through).
func (o *Operations) authorizeAdmin(ctx context.Context, projectID uuid.UUID) error {
	return o.Authz.CheckProject(ctx, projectID, models.ProjectRoleAdmin)
}

func validKind(kind models.EmailTemplateKind) bool {
	return kind == models.VerifyEmailTemplateKind || kind == models.ResetEmailTemplateKind || kind == models.TosEmailTemplateKind
}

func invalidKindErr() error {
	return fun.ErrValidation("email template kind must be one of: verify, reset, tos")
}
