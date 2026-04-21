package auth

import (
	"IdentityX/internal/platform/telemetry"
	_ "IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"IdentityX/internal/shared/validation"
	"encoding/json"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
	"go.uber.org/zap"
)

type Handler struct {
	users CommandService
	redis ports.RedisCacheService
}

func NewHandler(
	users CommandService,
	redis ports.RedisCacheService,
) *Handler {
	return &Handler{
		users: users,
		redis: redis,
	}
}

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,passwd,min=8,max=72"`
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user in the system.
// @Tags auth
// @Accept json
// @Produce json
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

	in := RegisterInput{
		Email:     req.Email,
		Password:  req.Password,
		ProjectID: nil,
	}

	if err := handler.users.Register(r.Context(), in); err != nil {
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
// @Description Authenticates a user and sets authentication cookies.
// @Tags auth
// @Accept json
// @Produce json
// @Param loginInfo body LoginUserRequest true "User login information"
// @Success 200 {object} object "Successful authentication"
// @Header 200 {string} Set-Cookie "access_token cookie for authentication"
// @Header 200 {string} Set-Cookie "refresh_token cookie for authentication"
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

	in := LoginInput{
		Email:    req.Email,
		Password: req.Password,

		Agent:     r.UserAgent(),
		IP:        validation.ClientIPString(validation.GetClientIP(r, validation.HTTPProxyConfig)),
		ProjectID: nil,
	}

	tokens, err := handler.users.Login(r.Context(), in)
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

	err := handler.users.Logout(ctx)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Logged out").Send(w)
}

// Refresh godoc
// @Summary Refreshes user security
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
	tokens, err := handler.users.Refresh(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Refreshed security").WithData(map[string]any{
		"access_token":  tokens.AccessTokenString,
		"refresh_token": tokens.RefreshTokenString,
	}).Send(w)
}

type RegisterProjectUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,passwd,min=8,max=72"`
}

// ProjectRegister godoc
// @Summary Register a new user into a client project
// @Description Registers a new user within a specific client project.
// @Tags auth
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to register user"
// @Param schema_type query string false "Schema type (default: core)"
// @Param flow_id query string false "Flow ID (default: none)"
// @Param version query string false "Version (default: 0)"
// @Param registerInfo body RegisterProjectUserRequest true "User registration information for the project"
// @Success 201 {object} object "User registered successfully in project"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input or missing project ID"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/register [post]
func (handler *Handler) ProjectRegister(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req RegisterProjectUserRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := RegisterInput{
		Email:     req.Email,
		Password:  req.Password,
		ProjectID: &projectID,
	}

	ctx := r.Context()
	if err := handler.users.Register(ctx, in); err != nil {
		resp.Error(err).Send(w)
		return
	}

	resp.Created("Registered user").Send(w)
}

type LoginProjectUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=72"`
}

// ProjectLogin godoc
// @Summary Authenticates a user into a client project
// @Description Authenticates a user within a specific client project and sets authentication cookies.
// @Tags auth
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to login user"
// @Param loginInfo body LoginProjectUserRequest true "User login information for the project"
// @Success 200 {object} object "Successfully authenticated into project"
// @Header 200 {string} Set-Cookie "access_token cookie for authentication"
// @Header 200 {string} Set-Cookie "refresh_token cookie for authentication"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input or missing project ID"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: Invalid credentials"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/login [post]
func (handler *Handler) ProjectLogin(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req LoginProjectUserRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	agent := r.UserAgent()
	ip := validation.ClientIPString(validation.GetClientIP(r, validation.HTTPProxyConfig))

	in := LoginInput{
		Email:     req.Email,
		Password:  req.Password,
		IP:        ip,
		Agent:     agent,
		ProjectID: &projectID,
	}

	ctx := r.Context()
	tokens, err := handler.users.Login(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Logged in").WithData(map[string]any{
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
// @Success 200 {object} object "JSON Web Key Set (JWKS)"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /.well-known/jwks.json [get]
func (handler *Handler) GetJWKS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jwks, err := handler.users.GetJWKS(ctx)
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
		// Just log if writing fails
		telemetry.Log().Error("failed to write JWKS response", zap.Error(err))
	}
}
