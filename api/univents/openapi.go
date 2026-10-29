// Package spec exposes the service's OpenAPI 3.1 specification, embedded
// from api-spec.gen.yml (api-spec.yml bundled with its problems by
// tools/problemgen) and served by the harness at /docs/openapi.yml.
package spec

import _ "embed"

// OpenAPISpec is the service's OpenAPI 3.1 specification, bundled.
//
//go:embed api-spec.gen.yml
var OpenAPISpec []byte
