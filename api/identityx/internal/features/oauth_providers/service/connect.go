package service

import (
	"context"
	"errors"

	"IdentityX/errs"
	"lib/database"

	"github.com/google/uuid"
)

// Connect starts an OAuth login: it resolves the scope's provider
// credentials, records a one-time login state, and returns the provider's
// authorization URL for the caller to redirect to. Provider policy lives
// in this module, so this is the flow's only contact with "is this
// provider configured/enabled".
//
// acceptedTos is the clickwrap consent flag: a project with terms only
// starts flows whose caller agreed to them, so the consent screen precedes
// the provider redirect (the callback itself has no UI to ask in).
func (o *Operations) Connect(ctx context.Context, provider string, projectID *uuid.UUID, acceptedTos bool) (string, error) {
	p, ok := o.Meta[provider]
	if !ok {
		return "", &errs.UnsupportedProviderError{Provider: provider}
	}

	if projectID != nil {
		// The project must exist before its provider row can resolve; a
		// missing project surfaces as not-found, not as a misconfigured row.
		_, err := o.Projects.GetByID(ctx, *projectID)
		if err != nil {
			return "", err
		}
		_, hasTos, err := o.Tos.CurrentTos(ctx, *projectID)
		if err != nil {
			return "", err
		}
		if hasTos && !acceptedTos {
			return "", errs.ErrTermsNotAccepted
		}
	}

	resolved, err := o.resolveLoginProvider(ctx, provider, projectID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return "", &errs.ProviderNotConfiguredError{Provider: provider, ForProject: true}
		}
		return "", err
	}

	state, err := o.createLoginState(ctx, provider, projectID)
	if err != nil {
		return "", err
	}

	return p.Config(resolved.Creds).AuthCodeURL(state), nil
}
