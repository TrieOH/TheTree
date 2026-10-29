package httpserver

import "encoding/json"

// ProblemDoc documents one declared problem type (ADR-0019). Each
// Backend's generated problems package lists its own in Docs, and the
// harness serves them at GET /problems/{slug}.
type ProblemDoc struct {
	Slug        string
	Status      int
	Title       string
	Description string
	// Extensions is an example of the problem's extension members, as a
	// JSON object.
	Extensions json.RawMessage
}
