// Package errs holds Informd's named errors (ADR-0016). It is a leaf: it
// imports only the standard library, so repos, services, authz, handlers
// and the HTTP mapper (internal/app) can all match on these errors without
// importing each other.
//
// A sentinel's text is the message a client reads. What status each error
// is answered with is decided in one place, internal/app's problem switch.
package errs
