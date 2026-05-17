package app

import (
	"context"
	"lib/database"
	"lib/errx"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	idx "git.trieoh.com/TrieOH/IdentityX-SDK-Go"
	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupValidator() *validator.Validate {
	var v = validator.New()
	if err := v.RegisterValidation("uuid7", func(fl validator.FieldLevel) bool {
		vv := fl.Field().String()

		u, err := uuid.Parse(vv)
		if err != nil {
			return false
		}

		return u.Version() == 7
	}); err != nil {
		errx.Exit(err, "failed to register uuid7 validator")
	}

	// Use JSON field names for better API responses
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return fld.Name
		}
		return name
	})

	return v
}

func SetupFUN(module string) {
	fun.SetConfig(fun.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        module,
	})

	v := SetupValidator()
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
		"chk_fields_type":             "field type must be one of: string, email, int, float, bool, select",
		"chk_fields_key_format":       "field key must start with a letter or underscore and contain only lowercase letters, digits and underscores",
		"chk_fields_select_behaviour": "select fields must have a behaviour of checkbox or radio",
		"chk_fields_select_options":   "select fields must have a non-empty options array",
		"chk_select_type":             "select_type must be one of: string, email, int, float, bool, select",
		"uniq_key_per_step":           "a field with this key already exists in this step",

		// namespaces
		"uniq_namespace_name_per_user":    "a namespace with this name already exists",
		"chk_valid_namespace_member_role": "invalid namespace member role",

		// responders
		"uniq_responder_email_on_system": "a responder with this email already exists on this system",
	})
}

func SetupCron(db *pgxpool.Pool) gocron.Scheduler {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		errx.Exit(err, "Failed to create Scheduler")
	}

	//_ = database.NewPGXTxRunner(db)
	// Call Jobs in jobs.go

	go scheduler.Start()
	log.Println("Started the cron Scheduler")
	return scheduler
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
