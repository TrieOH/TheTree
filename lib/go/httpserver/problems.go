package httpserver

import (
	"encoding/json"
	"errors"
	"maps"
	"net/http"

	"lib/database"

	"github.com/go-chi/chi/v5"
)

// errUnknownProblem is a GET /problems/{slug} for a slug nothing declares;
// the shared not-found problem answers it.
var errUnknownProblem = &database.NotFoundError{Named: errors.New("problem type not found")}

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

// problemDocBody is what GET /problems/{slug} answers: the problem type's
// documentation and an example of its body.
type problemDocBody struct {
	Type        string         `json:"type"`
	Status      int            `json:"status"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Example     map[string]any `json:"example"`
}

// serveProblemDocs answers GET /problems/{slug} from docs; an unknown slug
// is the not-found problem.
func serveProblemDocs(docs []ProblemDoc) http.HandlerFunc {
	bySlug := make(map[string]ProblemDoc, len(docs))
	for _, d := range docs {
		bySlug[d.Slug] = d
	}
	return func(w http.ResponseWriter, r *http.Request) {
		d, ok := bySlug[chi.URLParam(r, "slug")]
		if !ok {
			WriteError(w, r, errUnknownProblem)
			return
		}
		typ := ProblemType(d.Slug)
		example := map[string]any{}
		_ = json.Unmarshal(d.Extensions, &example)
		maps.Copy(example, map[string]any{
			"type":     typ,
			"status":   d.Status,
			"title":    d.Title,
			"detail":   d.Title,
			"instance": "/path/of/the/request",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(problemDocBody{
			Type:        typ,
			Status:      d.Status,
			Title:       d.Title,
			Description: d.Description,
			Example:     example,
		})
	}
}
