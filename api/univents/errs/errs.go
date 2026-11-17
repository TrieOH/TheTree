// Package errs holds Univents' named errors (ADR-0016). It is a leaf: it
// imports only the standard library, so repos, services, jobs, handlers
// and the HTTP mapper (internal/app) can all match on these errors without
// importing each other.
//
// A problem the API answers names its sentinel in problems.yml or in its
// feature's spec.yml (x-problem's error and message): tools/problemgen
// generates it into errs.gen.go, its message the detail a client reads,
// and maps it to its problem. Written by hand here are the errors that
// carry data (typed, mapped by internal/app/errmap's problem switch), the
// not-found errors and the schema's check constraints, answered with the
// shared problems. Request fields at fault go on through lib/fielderr.
package errs
