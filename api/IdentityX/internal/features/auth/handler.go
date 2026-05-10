package auth

import (
	"IdentityX/contracts"
	"encoding/json"
	"lib/telemetry"
	"net/http"
	"time"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/MintzyG/fun/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"go.uber.org/zap"
)

type Handler struct {
	commands CommandService
	queries  QueryService
}

func NewHandler(
	commands CommandService,
	queries QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

type ProjectIDParam struct {
	ProjectID string `fun_query:"project_id"`
}

func RegisterAuthRoutes(
	r *chi.Mux,
	h *Handler,
	disableRateLimit bool,
	jwt func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		if disableRateLimit {
			r.Use(httprate.Limit(5, 1*time.Minute, httprate.WithKeyFuncs(httprate.KeyByRealIP)))
		}

		r.Group(func(r chi.Router) {
			r.Use(middlewares.WithParams[ProjectIDParam](true))
			r.Post("/auth/register", h.Register)
			r.Post("/auth/login", h.Login)
			r.Get("/.well-known/jwks.json", h.GetJWKS)
		})

		r.Post("/auth/refresh", h.Refresh)
		r.With(jwt).Post("/auth/logout", h.Logout)
	})
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user in the system, optionally scoped to a project.
// @Tags auth
// @Accept json
// @Produce json
// @Param project_id query string false "Project UUID to scope the registration"
// @Param registerInfo body contracts.RegisterUserRequest true "User registration information"
// @Success 201 {object} object "User registered successfully"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/register [post]
func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	projectID := req.Query("project_id").UUIDPtr()
	var payload contracts.RegisterUserRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err := handler.commands.Register(r.Context(), payload.ToInput(projectID))
	if fun.Bail(w, err) {
		return
	}
	fun.Created("Registered user").Send(w)
}

// Login godoc
// @Summary Authenticate a user
// @Description Authenticates a user, optionally scoped to a project.
// @Tags auth
// @Accept json
// @Produce json
// @Param project_id query string false "Project UUID to scope the login"
// @Param loginInfo body contracts.LoginUserRequest true "User login information"
// @Success 200 {object} object "Access and Refresh tokens"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input provided"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: Invalid credentials"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/login [post]
func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	projectID := req.Query("project_id").UUIDPtr()
	var payload contracts.LoginUserRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	tokens, err := handler.commands.Login(r.Context(), payload.ToInput(projectID, r.UserAgent(), req.ClientIP()))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, tokens.ToResponse())
}

// Logout godoc
// @Summary Logs out a user
// @Description Logs out an authenticated user by clearing their session cookies.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} object "Successfully logged out"
// @Header 200 {string} Set-Cookie "clears the access_token cookie"
// @Header 200 {string} Set-Cookie "clears the refresh_token cookie"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/logout [post]
func (handler *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	err := handler.commands.Logout(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.OK("Logged out").Send(w)
}

// Refresh godoc
// @Summary Refreshes user tokens
// @Description Creates a new access and refresh token pair using a valid refresh token.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: refresh_token=yyy"
// @Success 200 {object} object "Tokens refreshed successfully"
// @Header 200 {string} Set-Cookie "access_token cookie for authentication"
// @Header 200 {string} Set-Cookie "refresh_token cookie for authentication"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Missing or invalid refresh token"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: Invalid or expired refresh token"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/refresh [post]
func (handler *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	refreshTokenCookie, err := req.Cookie("refresh_token").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	tokens, err := handler.commands.Refresh(r.Context(), contracts.ToRefreshInput(refreshTokenCookie, r.UserAgent(), req.ClientIP()))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, tokens.ToResponse())
}

// GetJWKS godoc
// @Summary Exposes the public key using a JWKS endpoint
// @Description Provides the JSON Web Key Set (JWKS) for verifying JWTs issued by the authentication service.
// @Tags auth
// @Accept json
// @Produce json
// @Param project_id query string false "Project UUID to scope the JWKS"
// @Success 200 {object} object "JSON Web Key Set (JWKS)"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /.well-known/jwks.json [get]
func (handler *Handler) GetJWKS(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	projectID := req.Query("project_id").UUIDPtr()
	jwks, err := handler.queries.GetJWKS(r.Context(), projectID)
	if fun.Bail(w, err) {
		return
	}
	data, err := json.Marshal(jwks)
	if fun.Bail(w, err) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=7200")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(data); err != nil {
		telemetry.Log().Error("failed to write JWKS response", zap.Error(err))
	}
}
