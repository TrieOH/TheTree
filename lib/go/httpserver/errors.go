package httpserver

import (
	"errors"
	"net/http"
	"reflect"
	"sync"

	"lib/authz"
	"lib/database"
	"lib/fielderr"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
)

// ErrorMap is a Backend's table from its named errors to HTTP (ADR-0016):
// the only place a Backend decides what status an error is. Lookup walks
// the error chain from the outermost link inward; the first link the table
// matches wins, so the most specific error maps, and the client reads that
// link's own message, never the %w context around it. An error nothing
// matches is a generic 500 whose cause goes to the logs, not the body.
//
// The shared errors (lib/database, lib/authz) are mapped after the
// Backend's entries, so a Backend only lists its own errors.
type ErrorMap []ErrorMapping

// ErrorMapping maps one named error, or one error type, to a status.
type ErrorMapping struct {
	// status reports the status a link maps to, and whether it matched.
	status     func(link error) (int, bool)
	validation bool
}

// Is maps the sentinel target to status.
func Is(target error, status int) ErrorMapping {
	return ErrorMapping{status: func(link error) (int, bool) {
		return status, sameError(link, target)
	}}
}

// As maps every error of type T to status.
func As[T error](status int) ErrorMapping {
	return AsStatus(func(T) int { return status })
}

// AsStatus maps every error of type T to the status statusOf picks from
// it, for errors whose status depends on their data (an upstream
// provider's status code, say). The decision still lives in the map.
func AsStatus[T error](statusOf func(T) int) ErrorMapping {
	return ErrorMapping{status: func(link error) (int, bool) {
		// One link of the chain, not the chain: resolve walks it.
		t, ok := link.(T) //nolint:errorlint // matching a single link on purpose
		if !ok {
			return 0, false
		}
		return statusOf(t), true
	}}
}

// Validation marks a 400 as a validation failure on the wire
// (VALIDATION_ERROR instead of BAD_REQUEST).
func (m ErrorMapping) Validation() ErrorMapping {
	m.validation = true
	return m
}

var sharedErrors = ErrorMap{
	Is(authz.ErrInsufficientRole, http.StatusForbidden),
	// A missing row is a 404 that reads as the resource's own error, even
	// when the Backend's map has no line for it.
	As[*database.NotFoundError](http.StatusNotFound),
	Is(database.ErrNotFound, http.StatusNotFound),
	Is(database.ErrUniqueViolation, http.StatusConflict),
	Is(database.ErrForeignKeyViolation, http.StatusBadRequest),
	Is(database.ErrCheckViolation, http.StatusBadRequest).Validation(),
	Is(database.ErrExclusionViolation, http.StatusConflict),
	As[*database.BatchError](http.StatusBadRequest).Validation(),
}

var (
	errorMapMu sync.RWMutex
	errorMap   ErrorMap
)

// UseErrorMap installs the Backend's error map for the process. Error
// responses keep the fun envelope until the move to fun/problems: the map
// is registered as fun's AppErrorMapper, so the strict handlers' errors and
// the errors the fun middlewares write go through the same table.
func UseErrorMap(m ErrorMap) {
	full := append(append(ErrorMap{}, m...), sharedErrors...)
	errorMapMu.Lock()
	errorMap = full
	errorMapMu.Unlock()
	fun.RegisterAppErrorMapper(func(err error) *fun.AppError {
		status, code, msg := full.resolve(err)
		if code == fun.CodeInternal {
			telemetry.Log().Error("unmapped error", zap.Int("status", status), zap.Error(err))
		}
		return &fun.AppError{Code: code, Message: msg, Fields: fieldsOf(err), Err: err}
	})
}

// fieldsOf lists the request fields a fielderr.Error in err's chain names.
func fieldsOf(err error) []fun.FieldError {
	fe, ok := errors.AsType[*fielderr.Error](err)
	if !ok {
		return nil
	}
	out := make([]fun.FieldError, len(fe.Fields))
	for i, f := range fe.Fields {
		out[i] = fun.FieldError{Field: f.Field, Message: f.Message}
	}
	return out
}

// WriteError writes err as an error response, through the Backend's map.
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	fun.Error(err).SendWithCtx(r.Context(), w)
}

// errorStatus is the status err is answered with: the root-fun error's own
// status while a Backend still returns those, else the map's.
func errorStatus(err error) int {
	if _, ok := errors.AsType[*fun.AppError](err); ok {
		return fun.Error(err).Code
	}
	errorMapMu.RLock()
	m := errorMap
	errorMapMu.RUnlock()
	if m == nil {
		return http.StatusInternalServerError
	}
	status, _, _ := m.resolve(err)
	return status
}

// Status is the status err is answered with under m (shared entries
// included); a Backend's tests use it to check its map covers its errors.
func (m ErrorMap) Status(err error) int {
	if _, ok := errors.AsType[*fun.AppError](err); ok {
		return fun.Error(err).Code // a lib error still on root fun
	}
	status, _, _ := append(append(ErrorMap{}, m...), sharedErrors...).resolve(err)
	return status
}

// resolve finds the outermost link of err the map matches.
func (m ErrorMap) resolve(err error) (int, fun.ErrorCode, string) {
	queue := []error{err}
	for len(queue) > 0 {
		link := queue[0]
		queue = queue[1:]
		for _, mapping := range m {
			status, ok := mapping.status(link)
			if !ok {
				continue
			}
			code := codeFor(status)
			if mapping.validation && status == http.StatusBadRequest {
				code = fun.CodeValidation
			}
			return status, code, link.Error()
		}
		switch u := link.(type) { //nolint:errorlint // walking the chain by hand
		case interface{ Unwrap() error }:
			next := u.Unwrap()
			if next != nil {
				queue = append(queue, next)
			}
		case interface{ Unwrap() []error }:
			queue = append(queue, u.Unwrap()...)
		}
	}
	return http.StatusInternalServerError, fun.CodeInternal, "internal error"
}

// sameError compares by identity, guarding links whose dynamic type is not
// comparable (== on those panics).
func sameError(link, target error) bool {
	t := reflect.TypeOf(link)
	return t != nil && t.Comparable() && link == target //nolint:errorlint // identity of one link
}

func codeFor(status int) fun.ErrorCode {
	switch status {
	case http.StatusBadRequest:
		return fun.CodeBadRequest
	case http.StatusUnauthorized:
		return fun.CodeUnauthorized
	case http.StatusPaymentRequired:
		return fun.CodePaymentRequired
	case http.StatusForbidden:
		return fun.CodeForbidden
	case http.StatusNotFound:
		return fun.CodeNotFound
	case http.StatusConflict:
		return fun.CodeConflict
	case http.StatusUnprocessableEntity:
		return fun.CodeUnprocessableEntity
	case http.StatusTooManyRequests:
		return fun.CodeTooManyRequests
	case http.StatusBadGateway:
		return fun.CodeBadGateway
	case http.StatusServiceUnavailable:
		return fun.CodeServiceUnavailable
	}
	return fun.CodeInternal
}
