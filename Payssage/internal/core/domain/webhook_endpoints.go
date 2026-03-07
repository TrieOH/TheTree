package domain

import (
	"TriePayments/internal/shared/errx"
	"TriePayments/internal/shared/validation"
	"strings"
	"time"

	"github.com/google/uuid"
)

type WebhookEndpoint struct {
	ID          uuid.UUID  `json:"id"`
	ScopeID     uuid.UUID  `json:"scope_id"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	URL         string     `json:"url"`
	Secret      string     `json:"-"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func NewWebhookEndpoint(workspaceID uuid.UUID, url, secret string) (*WebhookEndpoint, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, errx.Internal("product").SetMessage("error generating uuid").SetCause(err)
	}

	w := &WebhookEndpoint{
		ID:          id,
		WorkspaceID: workspaceID,
		URL:         url,
		Secret:      secret,
		CreatedAt:   time.Now(),
	}

	if err := w.validate(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *WebhookEndpoint) validate() error {
	return validation.Run(
		validation.RequireUUID("webhook_endpoint", "workspace_id", w.WorkspaceID),
		validation.RequireString("webhook_endpoint", "url", w.URL),
		validation.RequireString("webhook_endpoint", "secret", w.Secret),
		validation.Assert("webhook_endpoint", strings.HasPrefix(w.URL, "http://") || strings.HasPrefix(w.URL, "https://"), "url must be a valid http/https URL"),
	)
}

func (w *WebhookEndpoint) AddScope(scopeID uuid.UUID) {
	w.ScopeID = scopeID
}
