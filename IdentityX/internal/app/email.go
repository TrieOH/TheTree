package app

import (
	"IdentityX/internal/app/renderer"
	"IdentityX/internal/app/senders"
	"IdentityX/internal/shared/ports"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"path"
	texttemplate "text/template"

	"github.com/spf13/viper"
)

type MailBundle struct {
	Mailer   ports.Mailer
	Renderer ports.EmailRenderer
}

//go:embed renderer/templates/**
var templatesFS embed.FS

func NewBundle(rt runtime) (ports.EmailRenderer, ports.Mailer) {
	htmlTmpls, textTmpls, err := loadTemplates("renderer/templates")
	if err != nil {
		log.Fatalf("failed to load base email templates: %s", err)
	}

	var bundle = MailBundle{
		Mailer: senders.NewSMTPSender(
			rt.logger,
			rt.tracer,
			senders.SMTPConfig{
				Host:     viper.GetString("SMTP_HOST"),
				Port:     viper.GetString("SMTP_PORT"),
				Username: viper.GetString("SMTP_USERNAME"),
				Password: viper.GetString("SMTP_PASSWORD"),
				From:     viper.GetString("SMTP_FROM"),
				UseTLS:   viper.GetBool("SMTP_TLS"),
				StartTLS: viper.GetBool("SMTP_STARTTLS"),
			},
		),
		Renderer: renderer.NewMailRenderer(
			rt.logger,
			rt.tracer,
			htmlTmpls,
			textTmpls,
		),
	}
	return bundle.Renderer, bundle.Mailer
}

func loadTemplates(baseDir string) (
	map[string]*template.Template,
	map[string]*texttemplate.Template,
	error,
) {
	htmlTmpls := make(map[string]*template.Template)
	textTmpls := make(map[string]*texttemplate.Template)

	err := fs.WalkDir(templatesFS, baseDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		// p example:
		// internal/adapters/renderer/templates/verification/en.html.tmpl

		dir := path.Base(path.Dir(p)) // verification
		file := path.Base(p)          // en.html.tmpl

		switch {
		case path.Ext(file) == ".tmpl" && path.Ext(file[:len(file)-5]) == ".html":
			locale := file[:len(file)-len(".html.tmpl")]
			key := fmt.Sprintf("%s:%s", dir, locale)

			t, err := template.ParseFS(templatesFS, p)
			if err != nil {
				return err
			}
			htmlTmpls[key] = t

		case path.Ext(file) == ".tmpl" && path.Ext(file[:len(file)-5]) == ".txt":
			locale := file[:len(file)-len(".txt.tmpl")]
			key := fmt.Sprintf("%s:%s", dir, locale)

			t, err := texttemplate.ParseFS(templatesFS, p)
			if err != nil {
				return err
			}
			textTmpls[key] = t
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return htmlTmpls, textTmpls, nil
}
