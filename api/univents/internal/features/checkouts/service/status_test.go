package service_test

import (
	"univents/internal/app/errmap"
	"univents/internal/problems"
)

// statusOf is the HTTP status Univents answers err with — what a client
// sees, through its declared problems (ADR-0019).
func statusOf(err error) int {
	return problems.Mapper(errmap.Problem)(err).Status
}
