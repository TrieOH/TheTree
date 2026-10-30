package httpserver

import (
	"errors"
	"net/http"
	"strings"
	"sync"

	"lib/fielderr"
	"lib/telemetry"

	fp "github.com/MintzyG/problematic"
	"go.uber.org/zap"
)

// Every error response is an RFC 9457 problem (ADR-0019). The Backend's
// mapper (its generated problems.Mapper) is the only place an error gets a
// status: the strict handlers, the auth, scope, rate-limit and recover
// middlewares and the request binding all write through WriteError, so
// they answer alike.

// fieldError is one request field a problem names, its `errors` member.
type fieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var (
	problemsMu sync.RWMutex
	mapper     fp.ErrorMapper
	apiURL     string
)

// UseProblems installs the Backend's mapper and its public base URL
// (API_URL) for the process. A problem's type is {API_URL}/problems/{slug}.
func UseProblems(m fp.ErrorMapper, baseURL string) {
	problemsMu.Lock()
	mapper = m
	apiURL = strings.TrimSuffix(baseURL, "/")
	problemsMu.Unlock()
}

// internalServerError is the answer to an error the mapper doesn't know,
// and to every error before a mapper is installed.
func internalServerError() *fp.Problem {
	return fp.New(http.StatusInternalServerError).
		WithType("internal-server-error").
		WithTitle("Internal Server Error").
		WithDetail("Internal Server Error")
}

// problemFor is err's problem under the installed mapper, and the base URL
// its type resolves against.
func problemFor(err error) (*fp.Problem, string) {
	problemsMu.RLock()
	m, base := mapper, apiURL
	problemsMu.RUnlock()
	var p *fp.Problem
	if m != nil {
		p = m(err)
	}
	if p == nil {
		p = internalServerError()
	}
	return p, base
}

// WriteError writes err as its problem: the type resolved against API_URL,
// the request path as instance, and the request fields a fielderr.Lister
// in err's chain names as `errors`. A 5xx's cause goes to the logs, never
// the body.
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	p, base := problemFor(err)
	if p.Status >= http.StatusInternalServerError {
		telemetry.Log().Error("request failed", zap.Int("status", p.Status), zap.String("path", r.URL.Path), zap.Error(err))
	}
	p = p.Clone()
	p.Type = problemType(base, p.Type)
	p.Instance = r.URL.Path
	if fields := fieldsOf(err); len(fields) > 0 {
		p = p.WithExtension("errors", fields)
	}
	_ = fp.Write(w, p)
}

// ProblemType resolves a problem slug against API_URL, for problems nested
// in another's extensions (a batch's items) that WriteError doesn't reach.
func ProblemType(slug string) string {
	problemsMu.RLock()
	base := apiURL
	problemsMu.RUnlock()
	return problemType(base, slug)
}

func problemType(base, slug string) string {
	return base + "/problems/" + slug
}

// errorStatus is the status err is answered with.
func errorStatus(err error) int {
	p, _ := problemFor(err)
	return p.Status
}

// fieldsOf lists the request fields a fielderr.Lister in err's chain names.
func fieldsOf(err error) []fieldError {
	l, ok := errors.AsType[fielderr.Lister](err)
	if !ok {
		return nil
	}
	fields := make([]fieldError, len(l.FieldList()))
	for i, f := range l.FieldList() {
		fields[i] = fieldError{Field: f.Field, Message: f.Message}
	}
	return fields
}
