package middleware

import (
	"TrieForms/internal/shared/authz"
	"TrieForms/internal/shared/contracts"
	"TrieForms/internal/shared/ports"
	"net/http"
	"strings"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/TrieOH/IdentityX-SDK-Go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

type AuthMiddleware struct {
	idxClient idx.Client
	apiKeys   ports.ApiKeysRepo
	projects  ports.ProjectsRepo
	tracer    trace.Tracer
}

func NewAuthMiddleware(
	idxClient *idx.Client,
	apiKeys ports.ApiKeysRepo,
	projects ports.ProjectsRepo,
	tracer trace.Tracer,
) *AuthMiddleware {
	return &AuthMiddleware{
		idxClient: *idxClient,
		apiKeys:   apiKeys,
		projects:  projects,
		tracer:    tracer,
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

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				resp.Unauthorized().WithMsg("missing access token").WithModule("AuthMW").Send(w)
				return
			}

			_, tokenStr, found := strings.Cut(authHeader, "Bearer ")
			if !found || tokenStr == "" {
				resp.Unauthorized().WithMsg("invalid authorization header").WithModule("AuthMW").Send(w)
				return
			}

			accessClaims, err := mw.idxClient.Tokens.VerifyAccessToken(ctx, tokenStr)
			if err != nil {
				resp.Unauthorized().WithMsg("invalid access token").WithModule("AuthMW").Send(w)
				return
			}

			// Inject subject into context
			subject := authz.UserSubject{
				ID:    accessClaims.Sub.ID,
				Email: accessClaims.Sub.Email,
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

			project, err := mw.projects.GetByID(ctx, matched.ProjectID)
			if err != nil {
				resp.Unauthorized("workspace not found").WithModule("AuthMW").Send(w)
				return
			}

			ctx = authz.WithProject(ctx, project)
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

			if r.Header.Get("X-API-Key") != "" {
				mw.APIKey()(next).ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
				mw.Auth()(next).ServeHTTP(w, r.WithContext(ctx))
				return
			}

			resp.Unauthorized("missing credentials").WithModule("AuthMW").Send(w)
		})
	}
}
