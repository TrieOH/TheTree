package app

import (
	"testing"

	spec "payssage"

	"gopkg.in/yaml.v3"
)

// connectRequestSchema is the minimal YAML shape we read from the spec to
// assert the connect payload contract (D7): the caller never supplies the
// provider redirect URL — Payssage owns it.
type connectRequestSchema struct {
	Required   []string             `yaml:"required"`
	Type       yaml.Node            `yaml:"type"`
	Properties map[string]yaml.Node `yaml:"properties"`
}

func TestConnectPayload_NoProviderRedirectURL(t *testing.T) {
	var doc struct {
		Components struct {
			Schemas map[string]connectRequestSchema `yaml:"schemas"`
		} `yaml:"components"`
	}
	err := yaml.Unmarshal(spec.OpenAPISpec, &doc)
	if err != nil {
		t.Fatalf("parse spec: %v", err)
	}

	schema, ok := doc.Components.Schemas["ConnectRequest"]
	if !ok {
		t.Fatal("ConnectRequest schema not found in spec")
	}

	if _, ok := schema.Properties["provider_redirect_url"]; ok {
		t.Fatal("ConnectRequest still accepts provider_redirect_url (D7): the provider redirect URI is Payssage-owned config")
	}
	for _, req := range schema.Required {
		if req == "provider_redirect_url" {
			t.Fatal("ConnectRequest still requires provider_redirect_url (D7)")
		}
	}
	if len(schema.Required) != 2 || schema.Required[0] != "flow" || schema.Required[1] != "final_redirect_url" {
		t.Fatalf("ConnectRequest required = %v, want [flow final_redirect_url]", schema.Required)
	}
}
