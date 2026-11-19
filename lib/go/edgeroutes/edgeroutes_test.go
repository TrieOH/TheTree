package edgeroutes

import (
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"lib/httpserver"

	"github.com/go-chi/chi/v5"
)

const spec = `
openapi: 3.1.0
paths:
  /projects/{project_id}/actors:
    get: {}
  /projects:
    post: {}
  /.well-known/jwks.json:
    get: {}
  /auth/refresh:
    post: {}
  /health:
    get: {}
`

func TestSegments(t *testing.T) {
	got, err := Segments([]byte(spec))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{".well-known", "auth", "docs", "health", "problems", "projects"}
	if !slices.Equal(got, want) {
		t.Fatalf("Segments = %v, want %v", got, want)
	}
}

func TestSegmentsRejectsParamFirstSegment(t *testing.T) {
	_, err := Segments([]byte("paths:\n  /{tenant}/things:\n    get: {}\n"))
	if err == nil {
		t.Fatal("want an error for a path whose first segment is a parameter")
	}
}

func TestSegmentsRejectsEmptySpec(t *testing.T) {
	_, err := Segments([]byte("openapi: 3.1.0\n"))
	if err == nil {
		t.Fatal("want an error for a spec with no paths")
	}
}

func TestRender(t *testing.T) {
	out, err := Render([]byte(spec))
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.HasPrefix(s, "// Code generated") {
		t.Fatalf("missing generated header:\n%s", s)
	}
	if !strings.Contains(s, "\t\".well-known\",\n\t\"auth\",\n") {
		t.Fatalf("segments not rendered in order:\n%s", s)
	}
}

func TestVerify(t *testing.T) {
	dir := t.TempDir()
	err := os.MkdirAll(filepath.Join(dir, filepath.Dir(FileName)), 0o750)
	if err != nil {
		t.Fatal(err)
	}
	err = Verify([]byte(spec), dir)
	if err == nil {
		t.Fatal("want an error when the file is missing")
	}
	err = Write([]byte(spec), dir)
	if err != nil {
		t.Fatal(err)
	}
	err = Verify([]byte(spec), dir)
	if err != nil {
		t.Fatalf("freshly written file should verify: %v", err)
	}
	changed := strings.Replace(spec, "/auth/refresh", "/sessions/refresh", 1)
	err = Verify([]byte(changed), dir)
	if err == nil {
		t.Fatal("want an error after the spec gains a segment")
	}
}

// TestServedMatchesTheRouterFrame pins served to what httpserver.NewRouter
// mounts on every backend outside the spec.
func TestServedMatchesTheRouterFrame(t *testing.T) {
	mux, _ := httpserver.NewRouter(httpserver.Config{
		AppName:     "test",
		OpenAPISpec: []byte("openapi: 3.1.0"),
		Problems:    []httpserver.ProblemDoc{{Slug: "internal", Status: 500, Title: "Internal"}},
	})
	set := map[string]bool{}
	err := chi.Walk(mux, func(_, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		first, _, _ := strings.Cut(strings.TrimPrefix(route, "/"), "/")
		set[first] = true
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	frame := make([]string, 0, len(set))
	for s := range set {
		frame = append(frame, s)
	}
	slices.Sort(frame)
	if !slices.Equal(frame, served) {
		t.Fatalf("router frame serves %v, served lists %v", frame, served)
	}
}
