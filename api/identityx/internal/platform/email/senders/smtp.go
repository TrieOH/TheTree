package senders

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"net/textproto"

	"IdentityX/internal/shared/ports"

	"github.com/MintzyG/fun"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string

	// transport flags
	UseTLS   bool // implicit TLS (smtps, 465)
	StartTLS bool // STARTTLS (587)
}

type SMTPSender struct {
	logs   *zap.Logger
	tracer trace.Tracer
	queue  chan ports.Email
	cfg    SMTPConfig
}

var _ ports.Mailer = (*SMTPSender)(nil)

func NewSMTPSender(logs *zap.Logger, tracer trace.Tracer, cfg SMTPConfig) ports.Mailer {
	s := &SMTPSender{
		logs:   logs,
		tracer: tracer,
		queue:  make(chan ports.Email, 100),
		cfg:    cfg,
	}
	go s.worker()
	return s
}

func (s *SMTPSender) Send(ctx context.Context, email ports.Email) error {
	select {
	case s.queue <- email:
		return nil

	case <-ctx.Done():
		return ctx.Err()

	// queue is full → backpressure
	default:
		return fun.ErrInternal("BaseSMTPSender is unavailable")
	}
}

func (s *SMTPSender) worker() {
	// FIXME Worker goroutine has no graceful shutdown mechanism. If the application stops, queued emails will be lost.
	for email := range s.queue {
		ctx, span := s.tracer.Start(context.Background(), "email.smtp.send")

		if err := s.sendSMTP(ctx, email); err != nil {
			s.logs.Error(
				"failed to send email",
				zap.String("to", email.To),
				zap.Error(err),
			)
		}

		span.End()
	}
}

func buildMIME(email ports.Email, from string) ([]byte, error) {
	var buf bytes.Buffer

	writer := multipart.NewWriter(&buf)
	boundary := writer.Boundary()

	// Headers
	headers := map[string]string{
		"From":         from,
		"To":           email.To,
		"Subject":      email.Subject,
		"MIME-Version": "1.0",
		"Content-Type": `multipart/alternative; boundary="` + boundary + `"`,
	}

	for k, v := range headers {
		buf.WriteString(k + ": " + v + "\r\n")
	}
	buf.WriteString("\r\n")

	// ---- text/plain part
	textPart, err := writer.CreatePart(
		textproto.MIMEHeader{
			"Content-Type": {"text/plain; charset=utf-8"},
		},
	)
	if err != nil {
		return nil, err
	}

	if _, err := textPart.Write([]byte(email.TextBody)); err != nil {
		return nil, err
	}

	// ---- text/html part
	htmlPart, err := writer.CreatePart(
		textproto.MIMEHeader{
			"Content-Type": {"text/html; charset=utf-8"},
		},
	)
	if err != nil {
		return nil, err
	}

	if _, err := htmlPart.Write([]byte(email.HTMLBody)); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *SMTPSender) sendSMTP(ctx context.Context, email ports.Email) error {
	ctx, span := s.tracer.Start(ctx, "smtp.send")
	defer span.End()

	msg, err := buildMIME(email, s.cfg.From)
	if err != nil {
		return err
	}

	client, err := dialSMTP(s.cfg)
	if err != nil {
		return err
	}
	defer client.Quit()

	// ---- auth (only if configured)
	if s.cfg.Username != "" {
		auth := smtp.PlainAuth(
			"",
			s.cfg.Username,
			s.cfg.Password,
			s.cfg.Host,
		)

		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(auth); err != nil {
				return err
			}
		}
	}

	if err := client.Mail(s.cfg.From); err != nil {
		return err
	}

	if err := client.Rcpt(email.To); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	if _, err := w.Write(msg); err != nil {
		return err
	}

	return w.Close()
}

func dialSMTP(cfg SMTPConfig) (*smtp.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	// ---- implicit TLS (SMTPS)
	if cfg.UseTLS {
		tlsCfg := &tls.Config{
			ServerName: cfg.Host,
			MinVersion: tls.VersionTLS12,
		}

		conn, err := tls.Dial("tcp", addr, tlsCfg)
		if err != nil {
			return nil, err
		}

		return smtp.NewClient(conn, cfg.Host)
	}

	// ---- plaintext SMTP
	client, err := smtp.Dial(addr)
	if err != nil {
		return nil, err
	}

	// ---- optional STARTTLS
	if cfg.StartTLS {
		if ok, _ := client.Extension("STARTTLS"); ok {
			tlsCfg := &tls.Config{
				ServerName: cfg.Host,
				MinVersion: tls.VersionTLS12,
			}

			if err := client.StartTLS(tlsCfg); err != nil {
				client.Close()
				return nil, err
			}
		}
	}

	return client, nil
}
