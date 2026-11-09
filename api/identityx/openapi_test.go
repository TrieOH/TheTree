package spec

import (
	"os"
	"testing"

	"lib/edgeroutes"
)

// The Worker in cloudflare/ starts the container only for paths whose first
// segment this spec serves. Its committed list must follow api-spec.yml;
// UPDATE_ROUTES=1 (`just gen-routes`) rewrites it.
func TestEdgeRoutesMatchSpec(t *testing.T) {
	spec, err := os.ReadFile("api-spec.yml")
	if err != nil {
		t.Fatal(err)
	}
	err = edgeroutes.Verify(spec, ".", os.Getenv("UPDATE_ROUTES") != "")
	if err != nil {
		t.Fatal(err)
	}
}
