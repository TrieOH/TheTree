package senders

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/ports/outbounds"
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type SMTPSender struct {
	logs   *zap.Logger
	tracer trace.Tracer
}

var _ outbounds.Mailer = (*SMTPSender)(nil)

func NewSMTPSender(logs *zap.Logger, tracer trace.Tracer) outbounds.Mailer {
	return &SMTPSender{
		logs:   logs,
		tracer: tracer,
	}
}

func (s *SMTPSender) Send(ctx context.Context, toSend outbounds.Email) error {
	return apierr.ErrNotImplemented
}
