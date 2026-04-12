package middleware

import (
	"errors"
	"net/http"
	"payssage/internal/shared/authz"
	"payssage/internal/shared/contracts"
	"payssage/internal/shared/ports"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

type AuthMiddleware struct {
	gaClient   goauth.Client
	apiKeys    ports.ApiKeysRepo
	workspaces ports.WorkspaceRepo
	tracer     trace.Tracer
}

func NewAuthMiddleware(
	gaClient *goauth.Client,
	apiKeys ports.ApiKeysRepo,
	workspaces ports.WorkspaceRepo,
	tracer trace.Tracer,
) *AuthMiddleware {
	return &AuthMiddleware{
		gaClient:   *gaClient,
		apiKeys:    apiKeys,
		workspaces: workspaces,
		tracer:     tracer,
	}
}

// Auth is a middleware function that checks for valid access and refresh tokens.
// It validates the tokens, checks if the refresh token is revoked, and creates a principal from the tokens.
// The principal is then added to the request context.
func (mw *AuthMiddleware) Auth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := mw.tracer.Start(ctx, "Middleware.Auth")
			defer span.End()

			var err error
			defer func() {
				span.SetAttributes(attribute.Bool("success", err == nil))
			}()

			cookie, err := r.Cookie("svc_session")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					resp.Unauthorized("missing svc_session cookie").WithModule("AuthMW").Send(w)
					return
				}
				resp.Unauthorized("invalid svc_session cookie").WithModule("AuthMW").Send(w)
				return
			}

			sessionData, err := mw.gaClient.Sessions.Get(ctx, cookie.Value)
			if err != nil || sessionData == nil {
				resp.Unauthorized("service session not found").WithModule("AuthMW").Send(w)
				return
			}

			snapshot, err := authz.UnmarshalSnapshot(sessionData)
			if err != nil {
				resp.InternalServerError("invalid session payload").WithModule("AuthMW").Send(w)
				return
			}

			subject := authz.UserSubject{
				ID:    snapshot.UserID,
				Email: snapshot.Email,
			}
			ctx = authz.WithSubject(ctx, &subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (mw *AuthMiddleware) APIKey() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := mw.tracer.Start(ctx, "Middleware.APIKey")
			defer span.End()

			var err error
			defer func() {
				span.SetAttributes(attribute.Bool("success", err == nil))
			}()

			rawKey := r.Header.Get("X-API-Key")
			if rawKey == "" {
				resp.Unauthorized("missing X-API-Key header").WithModule("AuthMW").Send(w)
				return
			}

			// extract prefix (first 11 chars: "tp_" + 8 hex)
			if len(rawKey) < 11 {
				resp.Unauthorized("invalid api key format").WithModule("AuthMW").Send(w)
				return
			}
			prefix := rawKey[:11]

			candidates, err := mw.apiKeys.GetByPrefix(ctx, prefix)
			if err != nil || len(candidates) == 0 {
				resp.Unauthorized("invalid api key").WithModule("AuthMW").Send(w)
				return
			}

			var matched *contracts.APIKey
			for _, candidate := range candidates {
				if err := bcrypt.CompareHashAndPassword([]byte(candidate.KeyHash), []byte(rawKey)); err == nil {
					matched = &candidate
					break
				}
			}
			if matched == nil {
				resp.Unauthorized("invalid api key").WithModule("AuthMW").Send(w)
				return
			}

			workspace, err := mw.workspaces.GetByID(ctx, matched.WorkspaceID)
			if err != nil {
				resp.Unauthorized("workspace not found").WithModule("AuthMW").Send(w)
				return
			}

			ctx = authz.WithWorkspace(ctx, workspace)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (mw *AuthMiddleware) AnyAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := mw.tracer.Start(ctx, "Middleware.AnyAuth")
			defer span.End()

			rawKey := r.Header.Get("X-API-Key")
			if rawKey != "" {
				if len(rawKey) < 11 {
					resp.Unauthorized("invalid api key").WithModule("AuthMW").Send(w)
					return
				}

				candidates, err := mw.apiKeys.GetByPrefix(ctx, rawKey[:11])
				if err != nil || len(candidates) == 0 {
					resp.Unauthorized("invalid api key").WithModule("AuthMW").Send(w)
					return
				}

				var matched *contracts.APIKey
				for _, candidate := range candidates {
					if err := bcrypt.CompareHashAndPassword([]byte(candidate.KeyHash), []byte(rawKey)); err == nil {
						matched = &candidate
						break
					}
				}
				if matched == nil {
					resp.Unauthorized("invalid api key").WithModule("AuthMW").Send(w)
					return
				}

				workspace, err := mw.workspaces.GetByID(ctx, matched.WorkspaceID)
				if err != nil {
					resp.Unauthorized("invalid api key").WithModule("AuthMW").Send(w)
					return
				}

				ctx = authz.WithWorkspace(ctx, workspace)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// fallback: svc_session
			cookie, err := r.Cookie("svc_session")
			if err != nil {
				resp.Unauthorized("missing credentials").WithModule("AuthMW").Send(w)
				return
			}

			sessionData, err := mw.gaClient.Sessions.Get(ctx, cookie.Value)
			if err != nil || sessionData == nil {
				resp.Unauthorized("service session not found").WithModule("AuthMW").Send(w)
				return
			}

			snapshot, err := authz.UnmarshalSnapshot(sessionData)
			if err != nil {
				resp.InternalServerError("invalid session payload").WithModule("AuthMW").Send(w)
				return
			}

			subject := authz.UserSubject{
				ID:    snapshot.UserID,
				Email: snapshot.Email,
			}
			ctx = authz.WithSubject(ctx, &subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
