// Package assets embeds the built-in email templates into the binary.
// They are the default templates: projects may override them per kind via
// the email_templates table, and the override (or the default) is rendered
// at send time.
package assets

import (
	_ "embed"
	"html/template"
)

var (
	//go:embed emails/verify.html
	verifyEmailHTML string

	//go:embed emails/reset.html
	resetEmailHTML string
)

// Defaults are validated at init: a broken embedded template fails boot,
// not the first send.
func init() {
	for _, src := range []string{verifyEmailHTML, resetEmailHTML} {
		template.Must(template.New("default").Parse(src))
	}
}

// VerifyEmailBody returns the built-in verify-email template source.
func VerifyEmailBody() string { return verifyEmailHTML }

// ResetEmailBody returns the built-in reset-email template source.
func ResetEmailBody() string { return resetEmailHTML }
