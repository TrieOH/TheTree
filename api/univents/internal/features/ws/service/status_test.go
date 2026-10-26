package service_test

import "univents/internal/app/errmap"

// statusOf is the HTTP status Univents answers err with — what a client
// sees, through the one error map (ADR-0016).
func statusOf(err error) int {
	return errmap.Errors().Status(err)
}
