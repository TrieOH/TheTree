package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const sharedProblems = `
components:
  schemas:
    Problem:
      type: object
      required: [type, status, title, detail, instance]
      properties:
        type: { type: string }
    BatchItem:
      type: object
      x-problem-go-type: BatchItem
    NotFound:
      x-problem: { slug: not-found, status: 404, title: Not Found, description: Missing. }
    Validation:
      x-problem: { slug: validation, status: 400, title: Validation Failed, description: Bad input. }
    BatchFailed:
      x-problem: { slug: batch-failed, status: 400, title: Batch Failed, description: Items failed. }
      type: object
      required: [items]
      properties:
        items: { type: array, items: { $ref: '#/components/schemas/BatchItem' } }
`

const ownProblems = `
components:
  schemas:
    PaymentRejected:
      x-problem:
        slug: payment-rejected
        status: 400
        title: Payment Rejected
        description: |
          The provider refused
          the payment.
      type: object
      required: [message, retry_after]
      properties:
        message: { type: string, example: card declined }
        retry_after: { type: integer }
`

const spec = `openapi: 3.1.0
paths:
  /pay:
    post:
      x-problems: [payment-rejected, validation, not-found]
      responses:
        "200": { description: ok }
`

// fixture lays out api/problems.yml and api/svc/{api-spec,problems}.yml
// and returns the spec's path.
func fixture(t *testing.T, own, spec string) string {
	t.Helper()
	root := t.TempDir()
	svc := filepath.Join(root, "svc")
	write := func(path, body string) {
		t.Helper()
		err := os.MkdirAll(filepath.Dir(path), 0o750)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(path, []byte(body), 0o600)
		if err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join(root, "problems.yml"), sharedProblems)
	write(filepath.Join(svc, "problems.yml"), own)
	write(filepath.Join(svc, "api-spec.yml"), spec)
	return filepath.Join(svc, "api-spec.yml")
}

func generate(t *testing.T, specPath string) (map[string]any, string) {
	t.Helper()
	dir := filepath.Dir(specPath)
	err := run(specPath, filepath.Join(dir, "out.yml"), filepath.Join(dir, "p", "problems.gen.go"))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(dir, "out.yml")) //nolint:gosec // test temp dir
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(raw), "# "+header) {
		t.Errorf("bundle lacks the generated header:\n%s", raw)
	}
	var bundled map[string]any
	err = yaml.Unmarshal(raw, &bundled)
	if err != nil {
		t.Fatal(err)
	}
	src, err := os.ReadFile(filepath.Join(dir, "p", "problems.gen.go")) //nolint:gosec // test temp dir
	if err != nil {
		t.Fatal(err)
	}
	return bundled, string(src)
}

func dig(t *testing.T, v any, keys ...string) any {
	t.Helper()
	for _, k := range keys {
		m, ok := v.(map[string]any)
		if !ok {
			t.Fatalf("at %s: not a mapping: %#v", k, v)
		}
		v, ok = m[k]
		if !ok {
			t.Fatalf("missing key %s in %v", k, keys)
		}
	}
	return v
}

func TestBundle(t *testing.T) {
	bundled, _ := generate(t, fixture(t, ownProblems, spec))

	if dig(t, bundled, "components", "schemas", "Problem") == nil {
		t.Error("plain schemas are not merged")
	}
	schema := dig(t, bundled, "components", "schemas", "PaymentRejectedProblem")
	ext := dig(t, schema, "allOf").([]any)[1]
	if got := dig(t, ext, "properties", "status", "const"); got != 400 {
		t.Errorf("status const = %v", got)
	}
	if got := dig(t, ext, "properties", "title", "const"); got != "Payment Rejected" {
		t.Errorf("title const = %v", got)
	}
	// A title const must not become an enum constant named after its value.
	if got := dig(t, ext, "properties", "title", "x-go-type"); got != "string" {
		t.Errorf("title x-go-type = %v", got)
	}
	if got := dig(t, ext, "properties", "type", "pattern"); got != "/problems/payment-rejected$" {
		t.Errorf("type pattern = %v", got)
	}
	if dig(t, ext, "properties", "message") == nil {
		t.Error("members are not carried")
	}
	if got := dig(t, bundled, "components", "responses", "PaymentRejectedProblem", "description"); got != "Payment Rejected: The provider refused the payment." {
		t.Errorf("response description = %q", got)
	}

	op := dig(t, bundled, "paths", "/pay", "post")
	if dig(t, op, "x-problems") == nil {
		t.Error("x-problems must stay in the bundle for the writer")
	}
	if got := dig(t, op, "responses", "404", "$ref"); got != "#/components/responses/NotFoundProblem" {
		t.Errorf("404 = %v", got)
	}
	oneOf := dig(t, op, "responses", "400", "content", "application/problem+json", "schema", "oneOf").([]any)
	if len(oneOf) != 2 {
		t.Errorf("400 oneOf = %v, want the two 400 problems", oneOf)
	}
}

func TestBundleExpandsStatuses(t *testing.T) {
	byStatus := strings.Replace(spec, "[payment-rejected, validation, not-found]", "[400, not-found]", 1)
	bundled, _ := generate(t, fixture(t, ownProblems, byStatus))

	op := dig(t, bundled, "paths", "/pay", "post")
	got := dig(t, op, "x-problems").([]any)
	want := []any{"validation", "batch-failed", "payment-rejected", "not-found"}
	if len(got) != len(want) {
		t.Fatalf("x-problems = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("x-problems = %v, want %v", got, want)
		}
	}
	if got := dig(t, op, "responses", "400", "$ref"); got != "#/components/responses/Status400Problem" {
		t.Errorf("400 = %v, want the shared response of every 400 problem", got)
	}
	if got := dig(t, bundled, "components", "responses", "Status400Problem", "content", "application/problem+json", "schema", "$ref"); got != "#/components/schemas/Status400Problem" {
		t.Errorf("Status400Problem schema = %v, want the named union", got)
	}
	oneOf := dig(t, bundled, "components", "schemas", "Status400Problem", "oneOf").([]any)
	if len(oneOf) != 3 {
		t.Errorf("Status400Problem oneOf = %v, want every 400 problem", oneOf)
	}
}

func TestPackage(t *testing.T) {
	_, src := generate(t, fixture(t, ownProblems, spec))

	for _, want := range []string{
		"func PaymentRejected(detail string, message string, retryAfter int) *fp.Problem",
		`WithType("payment-rejected")`,
		`WithExtension("retry_after", retryAfter)`,
		"func BatchFailed(detail string, items []BatchItem) *fp.Problem",
		"func Mapper(app fp.ErrorMapper) fp.ErrorMapper",
		"Extensions: json.RawMessage(`{\"message\":\"card declined\",\"retry_after\":0}`)",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("package lacks %s", want)
		}
	}
}

func TestRejects(t *testing.T) {
	problem := func(body string) string {
		return "components:\n  schemas:\n    Bad:\n" + body
	}
	cases := map[string]struct{ own, spec, want string }{
		"undeclared slug": {
			own:  ownProblems,
			spec: "openapi: 3.1.0\npaths:\n  /x:\n    get:\n      x-problems: [nope]\n",
			want: "undeclared problem nope",
		},
		"status with no problem": {
			own:  ownProblems,
			spec: "openapi: 3.1.0\npaths:\n  /x:\n    get:\n      x-problems: [418]\n",
			want: "status 418 has no declared problem",
		},
		"status declared twice": {
			own:  ownProblems,
			spec: "openapi: 3.1.0\npaths:\n  /x:\n    get:\n      x-problems: [not-found]\n      responses:\n        \"404\": { description: x }\n",
			want: "404 is declared both",
		},
		"optional member": {
			own:  problem("      x-problem: { slug: bad, status: 400, title: Bad, description: Bad. }\n      properties:\n        why: { type: string }\n"),
			spec: spec,
			want: "member why must be required",
		},
		"standard member": {
			own:  problem("      x-problem: { slug: bad, status: 400, title: Bad, description: Bad. }\n      required: [detail]\n      properties:\n        detail: { type: string }\n"),
			spec: spec,
			want: "detail is a standard member",
		},
		"not an error status": {
			own:  problem("      x-problem: { slug: bad, status: 302, title: Bad, description: Bad. }\n"),
			spec: spec,
			want: "not a 4xx or 5xx",
		},
		"duplicate slug": {
			own:  problem("      x-problem: { slug: not-found, status: 404, title: Bad, description: Bad. }\n"),
			spec: spec,
			want: "slug not-found is declared by both",
		},
		"untyped ref": {
			own:  problem("      x-problem: { slug: bad, status: 400, title: Bad, description: Bad. }\n      required: [p]\n      properties:\n        p: { $ref: '#/components/schemas/Problem' }\n"),
			spec: spec,
			want: "has no x-problem-go-type",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			specPath := fixture(t, c.own, c.spec)
			dir := filepath.Dir(specPath)
			err := run(specPath, filepath.Join(dir, "out.yml"), filepath.Join(dir, "p", "problems.gen.go"))
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("err = %v, want it to mention %q", err, c.want)
			}
		})
	}
}

func TestParam(t *testing.T) {
	for in, want := range map[string]string{
		"message":     "message",
		"retry_after": "retryAfter",
		"seller-id":   "sellerId",
		"Status":      "status",
		"type":        "typeValue",
		"range":       "rangeValue",
		"string":      "stringValue",
	} {
		if got := param(in); got != want {
			t.Errorf("param(%q) = %q, want %q", in, got, want)
		}
	}
}
