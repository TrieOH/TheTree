// Package emails owns the verify/reset email concern end to end: the
// baked-in default templates, template rendering and mandate validation,
// the single-use action token lifecycle, and the async send job payload.
package emails

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"

	"IdentityX/assets"
	"IdentityX/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// Data are the variables available to every email template.
type Data struct {
	ActionURL     string
	ProjectName   string
	Expiry        int
	ProjectDomain string
	Email         string
}

// Template is a resolved email template: the project override when one
// exists, otherwise the built-in default. Source reports which.
type Template struct {
	Kind    models.EmailTemplateKind
	Subject string
	Body    string
	Source  string
}

const (
	SourceDefault  = "default"
	SourceOverride = "override"
)

// DefaultSubject returns the built-in subject for a template kind.
func DefaultSubject(kind models.EmailTemplateKind) string {
	switch kind {
	case models.VerifyEmailTemplateKind:
		return "Verify your email address"
	case models.ResetEmailTemplateKind:
		return "Reset your password"
	default:
		return ""
	}
}

// Default returns the baked-in template for a kind. The caller validates
// the kind; an unknown kind yields an empty template.
func Default(kind models.EmailTemplateKind) Template {
	return Template{
		Kind:    kind,
		Subject: DefaultSubject(kind),
		Body:    defaultBody(kind),
		Source:  SourceDefault,
	}
}

func defaultBody(kind models.EmailTemplateKind) string {
	switch kind {
	case models.VerifyEmailTemplateKind:
		return assets.VerifyEmailBody()
	case models.ResetEmailTemplateKind:
		return assets.ResetEmailBody()
	default:
		return ""
	}
}

// TemplateStore is the subset of the email-template repo the renderer
// needs: resolve the project override for a kind. ports.EmailTemplateRepo
// satisfies it.
type TemplateStore interface {
	GetByProjectAndKind(ctx context.Context, projectID uuid.UUID, kind models.EmailTemplateKind) (*models.EmailTemplate, error)
}

// ResolveTemplate returns the effective template for a project+kind:
// the project's override when present, otherwise the built-in default.
// projectID nil resolves to the default (platform-level emails).
func ResolveTemplate(ctx context.Context, store TemplateStore, projectID *uuid.UUID, kind models.EmailTemplateKind) (Template, error) {
	if projectID != nil {
		override, err := store.GetByProjectAndKind(ctx, *projectID, kind)
		if err == nil {
			return Template{
				Kind:    kind,
				Subject: override.Subject,
				Body:    override.Body,
				Source:  SourceOverride,
			}, nil
		}
		if !fun.Is(err, fun.CodeNotFound) {
			return Template{}, err
		}
	}
	return Default(kind), nil
}

// Render executes the template's subject and body with the data. Both
// subject and body are Go templates sharing the same variables.
func Render(t Template, d Data) (subject, body string, err error) {
	subject, err = renderString("subject", t.Subject, d)
	if err != nil {
		return "", "", fmt.Errorf("email template %s subject: %w", t.Kind, err)
	}
	body, err = renderString("body", t.Body, d)
	if err != nil {
		return "", "", fmt.Errorf("email template %s body: %w", t.Kind, err)
	}
	return subject, body, nil
}

func renderString(name, src string, d Data) (string, error) {
	tpl, err := template.New(name).Parse(src)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, d)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// actionURLVariable is the mandated variable every verify/reset template
// body must include at least once. It is rendered by IdentityX with the
// real link at send time; the project can never supply its own.
const actionURLVariable = "{{.ActionURL}}"

// Validate enforces the template contract before an override is saved:
//   - the body must reference {{.ActionURL}} at least once, and
//   - the template must parse and render with sample data.
//
// The ActionURL mandate is the security-critical piece: the actionable
// link always comes from IdentityX, never from the project.
func Validate(t Template) error {
	if !strings.Contains(t.Body, actionURLVariable) {
		return fun.ErrValidation(
			"email template body must include the {{.ActionURL}} variable at least once",
		)
	}
	_, _, err := Render(t, sentinelData())
	if err != nil {
		return fun.ErrValidation("email template does not render: " + err.Error())
	}
	return nil
}

// sentinelData renders a template with realistic values so conditional
// template logic (e.g. {{if .ProjectName}}) does not hide the ActionURL.
func sentinelData() Data {
	return Data{
		ActionURL:     "https://example.com/auth/verify?token=__SENTINEL__",
		ProjectName:   "Example Project",
		Expiry:        10,
		ProjectDomain: "example.com",
		Email:         "user@example.com",
	}
}

// ActionURL builds the forced link for a verify/reset email:
//
//	{base}/auth/{kind}?project_id={id}&token={token}   (project actor)
//	{base}/auth/{kind}?token={token}                   (platform actor)
//
// The path and query shape are fixed; only the template body around the
// link is project-customizable.
func ActionURL(baseDomain string, kind models.EmailTemplateKind, projectID *uuid.UUID, token string) string {
	var b strings.Builder
	b.WriteString(strings.TrimRight(baseDomain, "/"))
	fmt.Fprintf(&b, "/auth/%s?", kind)
	if projectID != nil {
		fmt.Fprintf(&b, "project_id=%s&", projectID)
	}
	fmt.Fprintf(&b, "token=%s", token)
	return b.String()
}

// DomainHost strips the scheme from a base domain for the ProjectDomain
// template variable.
func DomainHost(baseDomain string) string {
	return strings.TrimPrefix(strings.TrimPrefix(baseDomain, "https://"), "http://")
}

// SendAuthEmailArgs is the River job payload for an async verify/reset
// email. The token is minted and persisted in the request path (retry-safe:
// a River retry reuses the same token instead of minting a second one);
// the worker only resolves the template, renders, and sends.
type SendAuthEmailArgs struct {
	TemplateKind string     `json:"kind"`
	Token        string     `json:"token"`
	ToEmail      string     `json:"to_email"`
	ProjectID    *uuid.UUID `json:"project_id,omitempty"`
	ProjectName  string     `json:"project_name,omitempty"`
	BaseDomain   string     `json:"base_domain"`
	Expiry       int        `json:"expiry"`
}

func (SendAuthEmailArgs) Kind() string { return "auth_email.send" }
