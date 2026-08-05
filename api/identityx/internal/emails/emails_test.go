package emails

import (
	"strings"
	"testing"

	"IdentityX/models"

	"github.com/google/uuid"
)

func TestDefaultTemplatesRender(t *testing.T) {
	for _, kind := range models.AllEmailTemplateKinds {
		t.Run(string(kind), func(t *testing.T) {
			tpl := Default(kind)
			if tpl.Source != SourceDefault {
				t.Fatalf("source = %q, want %q", tpl.Source, SourceDefault)
			}
			if tpl.Subject == "" || tpl.Body == "" {
				t.Fatalf("default %s template is empty", kind)
			}
			subject, body, err := Render(tpl, Data{
				ActionURL:     "https://acme.example.com/auth/verify?token=abc",
				ProjectName:   "Acme",
				Expiry:        10,
				ProjectDomain: "acme.example.com",
				Email:         "user@example.com",
			})
			if err != nil {
				t.Fatalf("render default %s: %v", kind, err)
			}
			if !strings.Contains(body, "acme.example.com") {
				t.Errorf("rendered body does not contain the action url: %s", body)
			}
			if subject == "" {
				t.Errorf("rendered subject is empty")
			}
		})
	}
}

func TestValidateEnforcesActionURLMandate(t *testing.T) {
	tests := []struct {
		name string
		body string
		ok   bool
	}{
		{"action url present", `<a href="{{.ActionURL}}">verify</a>`, true},
		{"action url multiple times", `<a href="{{.ActionURL}}">a</a> <a href="{{.ActionURL}}">b</a>`, true},
		{"action url missing", `<p>click the link</p>`, false},
		{"action url typo", `<a href="{{.ActionUrL}}">verify</a>`, false},
		{"broken template", `<a href="{{.ActionURL">verify</a>`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tpl := Template{Kind: models.VerifyEmailTemplateKind, Subject: "s", Body: tt.body, Source: SourceOverride}
			err := Validate(tpl)
			if tt.ok && err != nil {
				t.Fatalf("want valid, got %v", err)
			}
			if !tt.ok && err == nil {
				t.Fatal("want validation error, got nil")
			}
		})
	}
}

func TestActionURLShape(t *testing.T) {
	projectID := uuid.New()
	token := "jwt-token"

	got := ActionURL("https://acme.example.com", models.VerifyEmailTemplateKind, &projectID, token)
	want := "https://acme.example.com/auth/verify?project_id=" + projectID.String() + "&token=" + token
	if got != want {
		t.Fatalf("project url = %q, want %q", got, want)
	}

	got = ActionURL("https://identityx.example.com", models.ResetEmailTemplateKind, nil, token)
	want = "https://identityx.example.com/auth/reset?token=" + token
	if got != want {
		t.Fatalf("platform url = %q, want %q", got, want)
	}
}

func TestDomainHost(t *testing.T) {
	if got := DomainHost("https://acme.example.com"); got != "acme.example.com" {
		t.Fatalf("DomainHost(https) = %q", got)
	}
	if got := DomainHost("acme.example.com"); got != "acme.example.com" {
		t.Fatalf("DomainHost(no scheme) = %q", got)
	}
}
