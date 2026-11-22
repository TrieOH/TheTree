// Package errs holds IdentityX's named errors (ADR-0016). It is a leaf:
// it imports only the standard library, so repos, services, authz, jobs,
// handlers and the HTTP mapper (internal/app) can all match on these
// errors without importing each other.
//
// Every error a problem answers is declared in the spec, in problems.yml or
// its feature's spec.yml, and tools/problemgen generates it into
// errs.gen.go along with its Mapper case: x-problem's error and detail for
// a sentinel, a typed error (a field per member, the detail a template
// over them) for a problem with members, and x-shared-errors for the
// not-found errors and check constraints answered with a shared problem.
// The detail is the text a client reads. Declare errors there first.
// Written by hand here are only the errors the spec can't express, answered
// by internal/app/errmap's problem switch.
package errs
