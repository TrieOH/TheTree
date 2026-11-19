package spec

import (
	"os"
	"testing"

	"lib/edgeroutes"
)

// The Worker in cloudflare/ starts the container only for paths whose first
// segment this spec serves. Its committed list must follow api-spec.yml;
// `just gen` rewrites it.
func TestEdgeRoutesMatchSpec(t *testing.T) {
	spec, err := os.ReadFile("api-spec.yml")
	if err != nil {
		t.Fatal(err)
	}
	err = edgeroutes.Verify(spec, ".")
	if err != nil {
		t.Fatal(err)
	}
}
