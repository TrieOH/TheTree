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
	projects  ports.ProjectRepo
	authz     *authz.Service
}

func NewOperations(templates ports.EmailTemplateRepo, projects ports.ProjectRepo, authz *authz.Service) *Operations {
	return &Operations{
		templates: templates,
		projects:  projects,
		authz:     authz,
	}
}

// authorizeAdmin resolves the caller and checks it holds at least the
// admin role on the project; unknown projects surface as 404.
func (o *Operations) authorizeAdmin(ctx context.Context, projectID uuid.UUID) error {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return err
	}
	_, err = o.projects.GetByID(ctx, projectID)
	if err != nil {
		return err
	}
	return o.authz.CheckProject(ctx, ident.Sub.ID, projectID, nil, models.ProjectRoleAdmin)
}

func validKind(kind models.EmailTemplateKind) bool {
	return kind == models.VerifyEmailTemplateKind || kind == models.ResetEmailTemplateKind
}

func invalidKindErr() error {
	return fun.ErrValidation("email template kind must be one of: verify, reset")
}
