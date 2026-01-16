package email

import (
	"GoAuth/internal/adapters/email/renderer"
	"GoAuth/internal/adapters/email/senders"
	"GoAuth/internal/infrastructure"
	"GoAuth/internal/ports/outbounds"
)

type MailBundle struct {
	Mailer   outbounds.Mailer
	Renderer outbounds.EmailRenderer
}

func NewBundle(infra infrastructure.Infra) MailBundle {
	return MailBundle{
		Mailer:   senders.NewSMTPSender(infra.Logger, infra.Tracer),
		Renderer: renderer.NewMailRenderer(infra.Logger, infra.Tracer),
	}
}
