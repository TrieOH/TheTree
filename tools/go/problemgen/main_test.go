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
		"detail without error": {
			own:  problem("      x-problem: { slug: bad, status: 400, title: Bad, description: Bad., detail: why }\n"),
			spec: spec,
			want: "a detail belongs to a problem's error",
		},
		"error without detail": {
			own:  problem("      x-problem: { slug: bad, status: 400, title: Bad, description: Bad., error: ErrBad }\n"),
			spec: spec,
			want: "needs a detail",
		},
		"detail repeats the title": {
			own:  problem("      x-problem: { slug: bad, status: 400, title: Bad Thing, description: Bad., error: ErrBad, detail: bad thing. }\n"),
			spec: spec,
			want: "repeats the title",
		},
		"detail ends a sentence": {
			own:  problem("      x-problem: { slug: bad, status: 400, title: Bad, description: Bad., error: ErrBad, detail: you did a bad thing. }\n"),
			spec: spec,
			want: "reads like any Go error string",
		},
		"detail capitalized": {
			own:  problem("      x-problem: { slug: bad, status: 400, title: Bad, description: Bad., error: ErrBad, detail: You did a bad thing }\n"),
			spec: spec,
			want: "reads like any Go error string",
		},
		"sentinel on a problem with members": {
			own:  problem("      x-problem: { slug: bad, status: 400, title: Bad, description: Bad., error: ErrBad, detail: why it failed }\n      required: [why]\n      properties:\n        why: { type: string }\n"),
			spec: spec,
			want: "ending in Error",
		},
		"old message key": {
			own:  problem("      x-problem: { slug: bad, status: 400, title: Bad, description: Bad., error: ErrBad, message: why it failed }\n"),
			spec: spec,
			want: "message is now detail",
		},
		"typed error naming no member": {
			own:  problem("      x-problem: { slug: bad, status: 400, title: Bad, description: Bad., error: BadError, detail: 'it failed: {how}' }\n      required: [why]\n      properties:\n        why: { type: string }\n"),
			spec: spec,
			want: "{how} in the detail is not a member",
		},
		"typed error with a stray brace": {
			own:  problem("      x-problem: { slug: bad, status: 400, title: Bad, description: Bad., error: BadError, detail: 'it failed: {why' }\n      required: [why]\n      properties:\n        why: { type: string }\n"),
			spec: spec,
			want: "braces in a detail only name members",
		},
		"error not a Go name": {
			own:  problem("      x-problem: { slug: bad, status: 400, title: Bad, description: Bad., error: badError, detail: why it failed }\n"),
			spec: spec,
			want: "must be an exported Go identifier starting with Err",
		},
		"untyped ref": {
			own:  problem("      x-problem: { slug: bad, status: 400, title: Bad, description: Bad. }\n      required: [p]\n      properties:\n        p: { $ref: '../problems.yml#/components/schemas/Problem' }\n"),
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

// featureSpec is a feature's spec.yml: a path the spec imports, a schema
// of its own, one of the spec's, and a problem.
const featureSpec = `paths:
  /pay:
    post:
      operationId: pay
      x-problems: [card-declined, not-found]
      requestBody:
        content:
          application/json:
            schema: { $ref: '#/components/schemas/PayInput' }
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema: { $ref: '../../../api-spec.yml#/components/schemas/Receipt' }
  /unused:
    get:
      responses:
        "200": { description: ok }
components:
  parameters:
    PayID: { name: pay_id, in: path, required: true, schema: { type: string } }
  schemas:
    PayInput:
      type: object
      properties:
        amount: { $ref: '#/components/schemas/Cents' }
    Cents: { type: integer }
    CardDeclined:
      x-problem: { slug: card-declined, status: 402, title: Card Declined, description: The card was declined. }
`

const featureRoot = `openapi: 3.1.0
paths:
  /pay:
    $ref: './internal/features/pay/spec.yml#/paths/~1pay'
  /unused:
    $ref: './internal/features/pay/spec.yml#/paths/~1unused'
components:
  schemas:
    Receipt: { type: object }
`

func writeFeature(t *testing.T, specPath, body string) {
	t.Helper()
	path := filepath.Join(filepath.Dir(specPath), "internal", "features", "pay", "spec.yml")
	err := os.MkdirAll(filepath.Dir(path), 0o750)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(path, []byte(body), 0o600)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFeatureSpec(t *testing.T) {
	specPath := fixture(t, ownProblems, featureRoot)
	writeFeature(t, specPath, featureSpec)
	bundled, src := generate(t, specPath)

	op := dig(t, bundled, "paths", "/pay", "post")
	if got := dig(t, op, "requestBody", "content", "application/json", "schema", "$ref"); got != "#/components/schemas/PayInput" {
		t.Errorf("own schema ref = %v", got)
	}
	if got := dig(t, op, "responses", "200", "content", "application/json", "schema", "$ref"); got != "#/components/schemas/Receipt" {
		t.Errorf("spec schema ref = %v", got)
	}
	if got := dig(t, bundled, "components", "schemas", "PayInput", "properties", "amount", "$ref"); got != "#/components/schemas/Cents" {
		t.Errorf("nested ref = %v", got)
	}
	if dig(t, bundled, "components", "parameters", "PayID") == nil {
		t.Error("the feature's parameters are not merged")
	}
	if dig(t, op, "responses", "402", "$ref") != "#/components/responses/CardDeclinedProblem" {
		t.Error("the feature's problem is not answered")
	}
	if !strings.Contains(src, "func CardDeclined(detail string)") {
		t.Error("the feature's problem has no constructor")
	}
}

func TestFeatureSpecRejects(t *testing.T) {
	cases := map[string]struct {
		root, feature, want string
	}{
		"path never imported": {
			root:    "openapi: 3.1.0\npaths:\n  /pay:\n    $ref: './internal/features/pay/spec.yml#/paths/~1pay'\ncomponents:\n  schemas:\n    Receipt: { type: object }\n",
			feature: featureSpec,
			want:    "paths /unused is never imported",
		},
		"imported under another path": {
			root:    "openapi: 3.1.0\npaths:\n  /other:\n    $ref: './internal/features/pay/spec.yml#/paths/~1pay'\n",
			feature: "paths:\n  /pay:\n    get: { responses: { \"200\": { description: ok } } }\n",
			want:    "under another path",
		},
		"missing component": {
			root:    "openapi: 3.1.0\npaths:\n  /pay:\n    $ref: './internal/features/pay/spec.yml#/paths/~1pay'\n",
			feature: "paths:\n  /pay:\n    get:\n      responses:\n        \"200\":\n          description: ok\n          content: { application/json: { schema: { $ref: '#/components/schemas/Gone' } } }\n",
			want:    "declares no schemas Gone",
		},
		"schema declared twice": {
			root:    featureRoot,
			feature: strings.Replace(featureSpec, "    Cents: { type: integer }", "    Cents: { type: integer }\n    Receipt: { type: object }", 1),
			want:    "Receipt is also declared by the spec",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			specPath := fixture(t, ownProblems, c.root)
			writeFeature(t, specPath, c.feature)
			dir := filepath.Dir(specPath)
			err := run(specPath, filepath.Join(dir, "out.yml"), filepath.Join(dir, "p", "problems.gen.go"))
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("err = %v, want it to mention %q", err, c.want)
			}
		})
	}
}

const namedProblems = `
components:
  schemas:
    AddSelf:
      x-problem:
        slug: add-self
        status: 400
        title: Add Self
        description: The caller tried to add themselves.
        error: ErrAddSelf
        detail: you can't add yourself to the project; ask an admin to add you
`

func TestNamedErrors(t *testing.T) {
	specPath := fixture(t, namedProblems, "openapi: 3.1.0\npaths:\n  /x:\n    post:\n      x-problems: [add-self]\n")
	dir := filepath.Dir(specPath)
	err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module svc\n\ngo 1.27\n"), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	bundled, src := generate(t, specPath)

	for _, want := range []string{`"svc/errs"`, "{errs.ErrAddSelf, AddSelf}"} {
		if !strings.Contains(src, want) {
			t.Errorf("package lacks %s", want)
		}
	}
	x := dig(t, bundled, "components", "schemas", "AddSelfProblem", "x-problem").(map[string]any)
	if _, ok := x["error"]; ok {
		t.Error("the bundle shows the Go error's name")
	}
	errsSrc, err := os.ReadFile(filepath.Join(dir, errsFile)) //nolint:gosec // test temp dir
	if err != nil {
		t.Fatal(err)
	}
	want := `ErrAddSelf = errors.New("you can't add yourself to the project; ask an admin to add you")`
	if !strings.Contains(string(errsSrc), want) {
		t.Errorf("errs lacks %s:\n%s", want, errsSrc)
	}
}

func TestTypedErrors(t *testing.T) {
	own := `
components:
  schemas:
    SellerMismatch:
      x-problem:
        slug: seller-mismatch
        status: 409
        title: Seller Mismatch
        description: The seller is elsewhere.
        error: SellerMismatchError
        detail: 'seller {seller_id} is not ours (got "{provider}"), 100% sure'
      type: object
      required: [seller_id, provider]
      properties:
        seller_id: { type: string, format: uuid }
        provider: { type: string }
`
	specPath := fixture(t, own, "openapi: 3.1.0\npaths:\n  /x:\n    post:\n      x-problems: [seller-mismatch]\n")
	dir := filepath.Dir(specPath)
	err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module svc\n\ngo 1.27\n"), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	bundled, src := generate(t, specPath)

	for _, want := range []string{
		`"uuid"`,
		"func SellerMismatch(detail string, sellerId uuid.UUID, provider string) *fp.Problem",
		"errors.AsType[*errs.SellerMismatchError](err)",
		"return SellerMismatch(e.Error(), e.SellerID, e.Provider)",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("package lacks %s:\n%s", want, src)
		}
	}
	if strings.Contains(src, "{errs.SellerMismatchError") {
		t.Error("a typed error is listed as a sentinel")
	}
	x := dig(t, bundled, "components", "schemas", "SellerMismatchProblem", "x-problem").(map[string]any)
	if _, ok := x["detail"]; ok {
		t.Error("the bundle shows the detail template")
	}
	errsSrc, err := os.ReadFile(filepath.Join(dir, errsFile)) //nolint:gosec // test temp dir
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"type SellerMismatchError struct {\n\tSellerID uuid.UUID\n\tProvider string\n}",
		`return fmt.Sprintf("seller %v is not ours (got \"%v\"), 100%% sure", e.SellerID, e.Provider)`,
	} {
		if !strings.Contains(string(errsSrc), want) {
			t.Errorf("errs lacks %s:\n%s", want, errsSrc)
		}
	}
}

func TestField(t *testing.T) {
	for in, want := range map[string]string{
		"provider":      "Provider",
		"seller_id":     "SellerID",
		"intent_status": "IntentStatus",
		"callback-url":  "CallbackURL",
	} {
		if got := field(in); got != want {
			t.Errorf("field(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSharedErrors(t *testing.T) {
	own := ownProblems + `
x-shared-errors:
  not-found:
    ErrFormNotFound: form not found
  validation:
    ErrInvalidStatus: "status must be one of: draft, open"
`
	specPath := fixture(t, own, spec)
	dir := filepath.Dir(specPath)
	err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module svc\n\ngo 1.27\n"), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	bundled, src := generate(t, specPath)

	for _, want := range []string{"{errs.ErrFormNotFound, NotFound}", "{errs.ErrInvalidStatus, Validation}"} {
		if !strings.Contains(src, want) {
			t.Errorf("package lacks %s", want)
		}
	}
	if _, ok := bundled["x-shared-errors"]; ok {
		t.Error("the bundle carries x-shared-errors")
	}
	errsSrc, err := os.ReadFile(filepath.Join(dir, errsFile)) //nolint:gosec // test temp dir
	if err != nil {
		t.Fatal(err)
	}
	want := `ErrFormNotFound = errors.New("form not found")`
	if !strings.Contains(string(errsSrc), want) {
		t.Errorf("errs lacks %s:\n%s", want, errsSrc)
	}
}

func TestSharedErrorsRejects(t *testing.T) {
	for name, tc := range map[string]struct{ block, want string }{
		"unknown slug":  {"  gone:\n    ErrGone: it is gone", "no shared problem"},
		"own problem":   {"  payment-rejected:\n    ErrRejected: it was rejected", "no shared problem"},
		"with members":  {"  batch-failed:\n    ErrBatch: some items failed", "typed error"},
		"bad name":      {"  not-found:\n    FormNotFound: form not found", "exported Go identifier"},
		"bad detail":    {"  not-found:\n    ErrFormNotFound: Form not found.", "lowercase"},
		"named twice":   {"  not-found:\n    ErrFormNotFound: form not found\n  validation:\n    ErrFormNotFound: form is invalid", "named by both"},
		"not a mapping": {"  not-found: [ErrFormNotFound]", "maps each error"},
	} {
		t.Run(name, func(t *testing.T) {
			specPath := fixture(t, ownProblems+"\nx-shared-errors:\n"+tc.block+"\n", spec)
			dir := filepath.Dir(specPath)
			err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module svc\n\ngo 1.27\n"), 0o600)
			if err != nil {
				t.Fatal(err)
			}
			err = run(specPath, filepath.Join(dir, "out.yml"), filepath.Join(dir, "p", "problems.gen.go"))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err = %v, want %q", err, tc.want)
			}
		})
	}
}

func TestSharedProblemNamesNoError(t *testing.T) {
	specPath := fixture(t, ownProblems, spec)
	root := filepath.Dir(filepath.Dir(specPath))
	shared := strings.Replace(sharedProblems, "description: Missing. }", "description: Missing., error: ErrGone, detail: it is gone for good }", 1)
	err := os.WriteFile(filepath.Join(root, "problems.yml"), []byte(shared), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Dir(specPath)
	err = run(specPath, filepath.Join(dir, "out.yml"), filepath.Join(dir, "p", "problems.gen.go"))
	if err == nil || !strings.Contains(err.Error(), "names no error") {
		t.Fatalf("err = %v", err)
	}
}

func TestCapitalized(t *testing.T) {
	for in, want := range map[string]bool{
		"This failed":           true,
		"this failed":           false,
		"IdentityX is down":     false,
		"OAuth state expired":   false,
		"API key expired":       false,
		"a path parameter is X": false,
	} {
		if got := capitalized(in); got != want {
			t.Errorf("capitalized(%q) = %v, want %v", in, got, want)
		}
	}
}
