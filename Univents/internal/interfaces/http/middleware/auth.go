package middleware

import (
	"errors"
	"net/http"
	"univents/internal/shared/authz"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AuthMiddleware struct {
	gaClient goauth.Client
	tracer   trace.Tracer
}

func NewAuthMiddleware(
	gaClient *goauth.Client,
	tracer trace.Tracer,
) *AuthMiddleware {
	return &AuthMiddleware{
		gaClient: *gaClient,
		tracer:   tracer,
	}
}

// Auth is a middleware that validates the Authorization header Bearer token.
// It injects the subject into the request context if valid.
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

			// Get the service session cookie set by the frontend
			cookie, err := r.Cookie("svc_session")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					resp.Unauthorized().WithMsg("missing service session cookie").WithModule("AuthMW").Send(w)
					return
				}
				resp.Unauthorized().WithMsg("invalid service session cookie").WithModule("AuthMW").Send(w)
				return
			}

			sessionID := cookie.Value

			// Lookup session in cache
			sessionData, err := mw.gaClient.Sessions.Get(ctx, sessionID)
			if err != nil {
				resp.Unauthorized("service session not found").WithModule("AuthMW").Send(w)
				return
			}

			// Unmarshal payload
			snapshot, err := authz.UnmarshalSnapshot(sessionData)
			if err != nil {
				resp.InternalServerError("invalid session payload").WithModule("AuthMW").Send(w)
				return
			}

			// Inject subject into context
			subject := authz.UserSubject{
				ID:    snapshot.UserID,
				Email: snapshot.Email,
			}
			ctx = authz.WithSubject(ctx, &subject)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
