package programs_test

import (
	"univents/internal/problems"
)

// statusOf is the HTTP status Univents answers err with — what a client
// sees, through its declared problems (ADR-0019).
func statusOf(err error) int {
	return problems.Mapper(nil)(err).Status
}
