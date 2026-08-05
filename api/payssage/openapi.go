// Package spec exposes the service's OpenAPI 3.1 specification, embedded
// from api-spec.yml and served by the harness at /docs/openapi.yml.
package spec

import _ "embed"

// OpenAPISpec is the service's OpenAPI 3.1 specification (api-spec.yml).
//
//go:embed api-spec.yml
var OpenAPISpec []byte
