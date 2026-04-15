package renderer

import (
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"bytes"
	"context"
	"html/template"
	texttemplate "text/template"

	"github.com/MintzyG/fail/v3"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type MailRenderer struct {
	logs      *zap.Logger
	tracer    trace.Tracer
	htmlTmpls map[string]*template.Template
	textTmpls map[string]*texttemplate.Template
}

var _ ports.EmailRenderer = (*MailRenderer)(nil)

func NewMailRenderer(
	logs *zap.Logger,
	tracer trace.Tracer,
	htmlTmpls map[string]*template.Template,
	textTmpls map[string]*texttemplate.Template,
) ports.EmailRenderer {
	return &MailRenderer{
		logs:      logs,
		tracer:    tracer,
		htmlTmpls: htmlTmpls,
		textTmpls: textTmpls,
	}
}

func (mr *MailRenderer) Verification(ctx context.Context, data ports.VerificationEmailData) (ports.Email, error) {
	ctx, span := mr.tracer.Start(ctx, "email.render.verification")
	defer span.End()

	key := "verification:" + data.Locale

	subject, textBody, htmlBody, err := mr.render(ctx, key, struct {
		UserID string
		Email  string
		Link   template.URL
	}{
		UserID: data.UserID.String(),
		Email:  data.Email,
		Link:   template.URL(viper.GetString("API_URL") + "/auth/verify?token=" + data.Token),
	})

	if err != nil {
		return ports.Email{}, fail.New(errx.SYSRenderingEmailFailed).With(err).WithArgs("verification")
	}

	return ports.Email{
		To:       data.Email,
		Subject:  subject,
		TextBody: textBody,
		HTMLBody: htmlBody,
	}, nil
}

func (mr *MailRenderer) PasswordReset(ctx context.Context, data ports.PasswordResetEmailData) (ports.Email, error) {
	ctx, span := mr.tracer.Start(ctx, "email.render.password_reset")
	defer span.End()

	key := "password_reset:" + data.Locale

	subject, textBody, htmlBody, err := mr.render(ctx, key, map[string]any{
		"UserID": data.UserID,
		"Email":  data.Email,
		"Link":   template.URL(viper.GetString("API_URL") + "/reset?token=" + data.Token),
	})

	if err != nil {
		return ports.Email{}, err
	}

	return ports.Email{
		To:       data.Email,
		Subject:  subject,
		TextBody: textBody,
		HTMLBody: htmlBody,
	}, nil
}

func (mr *MailRenderer) render(
	ctx context.Context,
	key string,
	data any,
) (subject, textBody, htmlBody string, err error) {
	textTmpl, ok := mr.textTmpls[key]
	if !ok {
		return "", "", "", fail.New(errx.EMAILTemplateNotFound).WithArgs(key, "text").RecordCtx(ctx)
	}

	htmlTmpl, ok := mr.htmlTmpls[key]
	if !ok {
		return "", "", "", fail.New(errx.EMAILTemplateNotFound).WithArgs(key, "html").RecordCtx(ctx)
	}

	var subjectBuf, textBuf, htmlBuf bytes.Buffer

	// subject (named template)
	if err = textTmpl.ExecuteTemplate(&subjectBuf, "subject", data); err != nil {
		return
	}

	// text body (root template)
	if err = textTmpl.Execute(&textBuf, data); err != nil {
		return
	}

	// html body
	if err = htmlTmpl.Execute(&htmlBuf, data); err != nil {
		return
	}

	return subjectBuf.String(), textBuf.String(), htmlBuf.String(), nil
}
