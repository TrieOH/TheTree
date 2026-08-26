package assets

import (
	"bytes"
	"embed"
	"html/template"
)

// WARNING: everything under assets/ is compiled into the binary via embed.
// Only put small text files here (email templates, JSON schemas, examples).
// Do NOT embed images, fonts, PDFs, or anything over a few KB.

//go:embed emails/*.html
var files embed.FS

var templates = template.Must(template.ParseFS(files, "emails/*.html"))

type RequestSignatureEmailData struct {
	SignatoryName string
	EventName     string
	EditionName   string
	Link          string
	ExpiresInDays int
}

func RenderRequestSignatureEmail(data RequestSignatureEmailData) (string, error) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "request_signature.html", data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type SignatureCreatedEmailData struct {
	SignatoryName string
	EventName     string
	EditionName   string
	RevokeLink    string
}

func RenderSignatureCreatedEmail(data SignatureCreatedEmailData) (string, error) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "signature_created.html", data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type CertGrantedEmailData struct {
	AttendeeName string
	EventName    string
	EditionName  string
	CertName     string
	CertLink     string
	VerifyLink   string
}

func RenderCertGrantedEmail(data CertGrantedEmailData) (string, error) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "cert_granted.html", data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type BadgeEmittedEmailData struct {
	AttendeeName string
	EventName    string
	EditionName  string
	BadgeLink    string
}

func RenderBadgeEmittedEmail(data BadgeEmittedEmailData) (string, error) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "badge_emitted.html", data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type TicketGiftedEmailData struct {
	RecipientName  string
	RecipientEmail string
	EventName      string
	EditionName    string
	TicketTypeName string
	ClaimLink      string
}

func RenderTicketGiftedEmail(data TicketGiftedEmailData) (string, error) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "ticket_gifted.html", data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
