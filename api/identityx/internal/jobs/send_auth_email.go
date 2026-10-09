package jobs

import (
	"context"

	"IdentityX/internal/emails"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/email"
	"lib/telemetry"

	"go.uber.org/zap"
)

// SendAuthEmailWorker renders and sends the async verify/reset email.
// The action token is minted and persisted in the request path, so a
// retry after a send failure reuses the same token (single-use
// semantics stay intact). Send errors are returned so the carrier retries
// it; the email client is only ever touched here.
type SendAuthEmailWorker struct {
	emailClient email.Sender
	templates   ports.EmailTemplateRepo
}

func NewSendAuthEmailWorker(emailClient email.Sender, templates ports.EmailTemplateRepo) *SendAuthEmailWorker {
	return &SendAuthEmailWorker{emailClient: emailClient, templates: templates}
}

func (w *SendAuthEmailWorker) Work(ctx context.Context, args emails.SendAuthEmailArgs) error {
	kind := models.EmailTemplateKind(args.TemplateKind)

	tpl, err := emails.ResolveTemplate(ctx, w.templates, args.ProjectID, kind)
	if err != nil {
		return err
	}

	subject, body, err := emails.Render(tpl, emails.Data{
		ActionURL:     emails.ActionURL(args.BaseDomain, kind, args.ProjectID, args.Token),
		ProjectName:   args.ProjectName,
		Expiry:        args.Expiry,
		ProjectDomain: emails.DomainHost(args.BaseDomain),
		Email:         args.ToEmail,
	})
	if err != nil {
		return err
	}

	err = w.emailClient.Send(email.Message{
		To:      []string{args.ToEmail},
		Subject: subject,
		Body:    body,
		HTML:    true,
	})
	if err != nil {
		telemetry.Log().Error("failed to send auth email",
			zap.String("kind", args.TemplateKind),
			zap.String("to", args.ToEmail),
			zap.Error(err),
		)
		return err
	}
	return nil
}
