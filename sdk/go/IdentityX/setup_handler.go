package idx

import (
	"encoding/json"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// SetupHandler is an http.Handler that serves the IdentityX project-setup
// lifecycle. Mount it at the route you configure as your setup URL in
// IdentityX.
//
//	handler := idx.NewSetupHandler(client)
//	mux.Handle("/idx/setup", handler)
//
// GET returns {"setup_complete": true|false} so callers can check readiness.
// POST receives a JSON body with api_key and project_id, persists credentials
// via the attached CredentialHandler, reconstructs the inner HTTP client, and
// marks the SDK as ready.
type SetupHandler struct {
	client *Client
}

// NewSetupHandler returns an http.Handler that processes the IdentityX
// project-setup callback and calls client.Setup on success.
func NewSetupHandler(client *Client) http.Handler {
	return &SetupHandler{client: client}
}

type setupPayload struct {
	APIKey    string    `json:"api_key"`
	ProjectID uuid.UUID `json:"project_id"`
}

func (h *SetupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fun.Respond(w, map[string]bool{
			"setup_complete": h.client.IsSetupComplete(),
		})
		return

	case http.MethodPost:
		// continue below

	default:
		fun.MethodNotAllowed("only GET and POST are accepted").Send(w)
		return
	}

	var p setupPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		fun.BadRequest("invalid JSON body").Send(w)
		return
	}

	if p.APIKey == "" || p.ProjectID == uuid.Nil {
		fun.BadRequest("api_key and project_id are required").Send(w)
		return
	}

	if err := h.client.Setup(p.APIKey, p.ProjectID); err != nil {
		fun.InternalServerError(err.Error()).Send(w)
		return
	}

	fun.Respond(w, map[string]string{"status": "ok"})
}
