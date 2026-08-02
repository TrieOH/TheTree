package oauth_providers

import (
	"IdentityX/models"
	"context"
	"lib/oauth"
	"lib/telemetry"

	"github.com/google/uuid"
)

// providerOrder keeps discovery output deterministic. The registry is a map;
// the login modal renders buttons in this order.
var providerOrder = []string{"google", "github"}

// ListEnabledProviders is the public discovery route: the providers a
// project (or IdentityX itself, when projectID is nil) currently accepts
// logins with. Disabled providers are hidden so the frontend modal never
// offers a login that would be rejected.
func (o *Operations) ListEnabledProviders(ctx context.Context, projectID *uuid.UUID) ([]models.OAuthProvider, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListEnabledProviders")
	defer span.End()

	if projectID == nil {
		var out []models.OAuthProvider
		for _, name := range providerOrder {
			if _, ok := oauth.EnvCredentials(name); ok {
				out = append(out, models.OAuthProvider(name))
			}
		}
		return out, nil
	}

	_, err := o.projects.GetByID(ctx, *projectID)
	if err != nil {
		return nil, err
	}

	rows, err := o.providers.ListByProject(ctx, *projectID)
	if err != nil {
		return nil, err
	}
	var out []models.OAuthProvider
	for _, row := range rows {
		if row.Enabled {
			out = append(out, row.Provider)
		}
	}
	return out, nil
}
