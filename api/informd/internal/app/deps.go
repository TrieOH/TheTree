package app

import (
	"context"
	"lib/validator"
	"net/http"
	"time"

	"lib/database"
	"lib/errx"

	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	fm "github.com/MintzyG/fun/middlewares"
	"github.com/go-chi/chi/v5"
)

func SetupFUN(module string) {
	fun.SetConfig(fun.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        module,
	})

	v := validator.SetupValidator()
	bind.SetValidator(v)
	fun.SetPathParamFunc(func(r *http.Request, key string) string {
		return chi.URLParam(r, key)
	})
}

func SetupConstraintMessages() {
	database.SetConstraintErrorRegistry(database.ConstraintRegistry{
		// forms
		"chk_forms_valid_status":       "status must be one of: draft, open, closed, archived",
		"chk_forms_valid_status_state": "opened_at, closed_at or archived_at must be set when status is open, closed or archived",
		"uniq_form_name_per_namespace": "a form with this name already exists in this namespace",
		"uniq_name_per_user":           "an API key with this name already exists",

		// fields
		"chk_fields_type":       "field type must be one of: string, email, int, float, bool, date, time, datetime, select, file, phone, url",
		"chk_fields_key_format": "field key must start with a letter or underscore and contain only lowercase letters, digits and underscores",
		"uniq_key_per_step":     "a field with this key already exists in this step",

		// select
		"chk_select_behaviour":  "select fields must have a behaviour of checkbox or radio",
		"chk_select_options":    "select fields must have a non-empty options array",
		"chk_select_value_type": "select_type must be one of: email, int, float, date, time, datetime, phone, url",

		// namespaces
		"uniq_namespace_name_per_user":    "a namespace with this name already exists",
		"chk_valid_namespace_member_role": "invalid namespace member role",

		// response
		"uniq_answer_per_field_per_response": "this field already has an answer in this response",

		// responders
		"uniq_responder_email_on_system": "a responder with this email already exists on this system",
		"chk_invite_not_both_named":      "an invite cannot have both responder_id and email set",
		"form_invites_token_key":         "an invite with this token already exists",
	})
}

func SetupIdentityX(cfg Config) *idx.Client {
	client, err := idx.NewClient(idx.Config{
		BaseURL:   cfg.IdxURL,
		APIKey:    cfg.IdxAPIKey,
		ProjectID: cfg.IdxProjectID,
		Debug:     true,
	})
	if err != nil {
		errx.Exit(err, "error creating identity_x client")
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.Tokens.GetJWKS(ctx, true)
		if err != nil {
			errx.Exit(err, "error fetching initial JWKS")
		}
	}()
	return client
}

func setupAuthMiddlewares() *fm.Middleware[*idx.AccessClaims] {
	keyFunc := func(ctx context.Context, tokenStr string) (*idx.AccessClaims, error) {
		return app.idxClient.Tokens.VerifyAccessToken(ctx, tokenStr)
	}

	jwtHook := func(ctx context.Context, claims *idx.AccessClaims) (context.Context, error) {
		return idx.WithIdentity(ctx, &idx.Identity{
			Sub: idx.Subject{
				ID:           claims.Sub.ID,
				ProjectID:    claims.Sub.ProjectID,
				Email:        claims.Sub.Email,
				Type:         claims.Sub.Type,
				Capabilities: claims.Sub.Capabilities,
				Metadata:     claims.Sub.Metadata,
			},
			Cred: idx.Credential{
				Type: "token",
			},
		}), nil
	}

	apiKeyHook := func(ctx context.Context, rawKey string) (context.Context, error) {
		return nil, fun.ErrNotImplemented("api keys are not yet supported")
	}

	return fm.New[*idx.AccessClaims](keyFunc, jwtHook, apiKeyHook)
}
