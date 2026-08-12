package store

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ServeStoreStream is the raw `GET /editions/{edition_id}/store/stream`
// route (public, D10): the storefront's live stock. Unknown editions get a
// plain 404; known ones stream the snapshot then deltas (hijacked
// close-delimited SSE — the harness WriteTimeout does not apply).
func (h *Handlers) ServeStoreStream(w http.ResponseWriter, r *http.Request) {
	editionID, err := uuid.Parse(chi.URLParam(r, "edition_id"))
	if err != nil {
		http.Error(w, "invalid edition_id", http.StatusBadRequest)
		return
	}
	h.ops.ServeStream(w, r, editionID)
}
