package email_templates

import (
	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"
	"context"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

type Operations struct {
	templates ports.EmailTemplateRepo
	authz     *authz.Service
}

func NewOperations(templates ports.EmailTemplateRepo, authz *authz.Service) *Operations {
	return &Operations{
		templates: templates,
		authz:     authz,
	}
}

// authorizeAdmin gates a write on the authenticated actor holding an
// admin or owner role on the project; unknown projects surface as 404
// (CheckProject passes the project lookup through).
func (o *Operations) authorizeAdmin(ctx context.Context, projectID uuid.UUID) error {
	return o.authz.CheckProject(ctx, projectID, models.ProjectRoleAdmin)
}

func validKind(kind models.EmailTemplateKind) bool {
	return kind == models.VerifyEmailTemplateKind || kind == models.ResetEmailTemplateKind
}

func invalidKindErr() error {
	return fun.ErrValidation("email template kind must be one of: verify, reset")
}
