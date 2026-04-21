package auth

import (
	"IdentityX/internal/platform/telemetry"
	_ "IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/validation"
	"encoding/json"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
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

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,passwd,min=8,max=72"`
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user in the system, optionally scoped to a project.
// @Tags auth
// @Accept json
// @Produce json
// @Param project_id query string false "Project UUID to scope the registration"
// @Param registerInfo body RegisterUserRequest true "User registration information"
// @Success 201 {object} object "User registered successfully"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/register [post]
func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	var projectID *uuid.UUID
	if raw := r.URL.Query().Get("project_id"); raw != "" {
		id, err := validation.ParseUUID(raw, "project_id")
		if err != nil {
			resp.FromError(err).Send(w)
			return
		}
		projectID = &id
	}

	if err := handler.commands.Register(r.Context(), RegisterInput{
		Email:     req.Email,
		Password:  req.Password,
		ProjectID: projectID,
	}); err != nil {
		resp.Error(err).Send(w)
		return
	}

	resp.Created("Registered user").Send(w)
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=72"`
}

// Login godoc
// @Summary Authenticate a user
// @Description Authenticates a user, optionally scoped to a project.
// @Tags auth
// @Accept json
// @Produce json
// @Param project_id query string false "Project UUID to scope the login"
// @Param loginInfo body LoginUserRequest true "User login information"
// @Success 200 {object} object "Successful authentication"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input provided"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: Invalid credentials"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/login [post]
func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginUserRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	var projectID *uuid.UUID
	if raw := r.URL.Query().Get("project_id"); raw != "" {
		id, err := validation.ParseUUID(raw, "project_id")
		if err != nil {
			resp.FromError(err).Send(w)
			return
		}
		projectID = &id
	}

	tokens, err := handler.commands.Login(r.Context(), LoginInput{
		Email:     req.Email,
		Password:  req.Password,
		Agent:     r.UserAgent(),
		IP:        validation.ClientIPString(validation.GetClientIP(r, validation.HTTPProxyConfig)),
		ProjectID: projectID,
	})
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Logged in").WithData(map[string]any{
		"access_token":  tokens.AccessTokenString,
		"refresh_token": tokens.RefreshTokenString,
	}).Send(w)
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
	ctx := r.Context()

	err := handler.commands.Logout(ctx)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Logged out").Send(w)
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
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		resp.Unauthorized("error getting refresh token").AddTrace(err).Send(w)
		return
	}

	if refreshTokenCookie.Value == "" {
		resp.Unauthorized("missing refresh token value").Send(w)
		return
	}

	in := RefreshInput{
		RefreshCookie: refreshTokenCookie,
		Agent:         r.UserAgent(),
		IP:            validation.ClientIPString(validation.GetClientIP(r, validation.HTTPProxyConfig)),
	}

	ctx := r.Context()
	tokens, err := handler.commands.Refresh(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Refreshed tokens").WithData(map[string]any{
		"access_token":  tokens.AccessTokenString,
		"refresh_token": tokens.RefreshTokenString,
	}).Send(w)
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
	ctx := r.Context()

	var projectID *uuid.UUID
	if raw := r.URL.Query().Get("project_id"); raw != "" {
		id, err := validation.ParseUUID(raw, "project_id")
		if err != nil {
			resp.FromError(err).Send(w)
			return
		}
		projectID = &id
	}

	jwks, err := handler.queries.GetJWKS(ctx, projectID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	data, err := json.Marshal(jwks)
	if err != nil {
		apiErr := fail.New(errx.SYSJWKSEncodingFailed).With(err).RecordCtx(ctx)
		resp.FromError(apiErr).Send(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=7200")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(data); err != nil {
		telemetry.Log().Error("failed to write JWKS response", zap.Error(err))
	}
}
