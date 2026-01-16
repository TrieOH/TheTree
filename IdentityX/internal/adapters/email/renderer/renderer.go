package renderer

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/ports/outbounds"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type MailRenderer struct {
	logs   *zap.Logger
	tracer trace.Tracer
}

var _ outbounds.EmailRenderer = (*MailRenderer)(nil)

func NewMailRenderer(logs *zap.Logger, tracer trace.Tracer) outbounds.EmailRenderer {
	return &MailRenderer{
		logs:   logs,
		tracer: tracer,
	}
}

func (mr *MailRenderer) Verification(data outbounds.VerificationEmailData) (outbounds.Email, error) {
	return outbounds.Email{}, apierr.ErrNotImplemented
}

func (mr *MailRenderer) PasswordReset(data outbounds.PasswordResetEmailData) (outbounds.Email, error) {
	return outbounds.Email{}, apierr.ErrNotImplemented
}
