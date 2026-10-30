// Package errs holds Univents' named errors (ADR-0016). It is a leaf: it
// imports only the standard library, so repos, services, jobs, handlers
// and the HTTP mapper (internal/app) can all match on these errors without
// importing each other.
//
// A sentinel's text is the message a client reads. Errors that carry data
// are typed; request fields go on through lib/fielderr. What status each
// error is answered with is decided in one place, internal/app/errmap's problem switch.
package errs
