package users

import (
	_ "IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"IdentityX/internal/shared/validation"
	"encoding/json"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
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

	in := RegisterUserInput{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := handler.users.Register(r.Context(), in); err != nil {
		resp.FromError(err).Send(w)
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

	in := LoginUserInput{
		Email:    req.Email,
		Password: req.Password,

		Agent: r.UserAgent(),
		IP:    validation.ClientIPString(validation.GetClientIP(r, validation.HTTPProxyConfig)),
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
	tokens, err := handler.users.Refresh(ctx, in)
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
	}
}

type RegisterProjectUserRequest struct {
	Email        string           `json:"email" validate:"required,email,max=255"`
	Password     string           `json:"password" validate:"required,passwd,min=8,max=72"`
	CustomFields *json.RawMessage `json:"custom_fields"`
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

	schemaType := r.URL.Query().Get("schema_type")
	flowID := r.URL.Query().Get("flow_id")

	in := ProjectRegisterInput{
		Email:        req.Email,
		Password:     req.Password,
		CustomFields: req.CustomFields,
		ProjectID:    projectID,
		SchemaType:   schemaType,
		FlowID:       flowID,
	}

	if in.CustomFields == nil {
		meta := json.RawMessage("{}")
		in.CustomFields = &meta
	}

	ctx := r.Context()
	if err := handler.users.RegisterProjectUser(ctx, in); err != nil {
		resp.FromError(err).Send(w)
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

	in := ProjectLoginInput{
		Email:     req.Email,
		Password:  req.Password,
		ProjectID: projectID,
		IP:        ip,
		Agent:     agent,
	}

	ctx := r.Context()
	tokens, err := handler.users.LoginProjectUser(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Logged in").WithData(map[string]any{
		"access_token":  tokens.AccessTokenString,
		"refresh_token": tokens.RefreshTokenString,
	}).Send(w)
}

// ProjectLogout godoc
// @Summary Logs out a project user
// @Description Logs out an authenticated project user by revoking their GoAuth session.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: refresh_token=yyy"
// @Success 200 {object} object "Successfully logged out"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/logout [post]
func (handler *Handler) ProjectLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projectID, rs := validation.GetUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil || refreshTokenCookie.Value == "" {
		resp.Unauthorized().WithMsg("missing refresh_token cookie").Send(w)
		return
	}

	in := ProjectLogoutInput{
		ProjectID:          projectID,
		RefreshTokenCookie: refreshTokenCookie,
	}

	err = handler.users.LogoutProjectUser(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Logged out").Send(w)
}

// Verify godoc
// @Summary Verify user email
// @Description Verifies a user's email address using a verification token.
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification Token"
// @Success 200 {object} object "User verified successfully"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Missing or invalid token"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/verify [get]
func (handler *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	token, rs := validation.GetString(r, "token")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.users.Verify(ctx, token)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("user verified, please refresh").Send(w)
}

// ResendVerificationEmail godoc
// @Summary Resend verification email
// @Description Resends the email verification link to the currently authenticated user.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} object "Verification email resent successfully"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/verify/resend [post]
func (handler *Handler) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := handler.users.ResendVerificationEmail(ctx)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("verification email resent successfully").Send(w)
}

type ForgotPasswordRequest struct {
	Email     string     `json:"email" validate:"required,email"`
	ProjectID *uuid.UUID `json:"project_id"`
}

// ForgotPassword godoc
// @Summary Initiate password reset
// @Description Sends a password reset email to the user if the account exists.
// @Tags auth
// @Accept json
// @Produce json
// @Param forgotInfo body ForgotPasswordRequest true "User email and optional project ID"
// @Success 200 {object} object "Forgot password email sent"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/forgot-password [post]
func (handler *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := ForgotPasswordInput{
		Email:     req.Email,
		ProjectID: req.ProjectID,
	}

	ctx := r.Context()
	err := handler.users.ForgotPassword(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("forgot password email sent").Send(w)
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,passwd,min=8,max=72"`
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Resets the user's password using a valid reset token.
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Reset password token"
// @Param resetInfo body ResetPasswordRequest true "New password information"
// @Success 200 {object} object "Password reset successfully"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input or token"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: Invalid or expired token"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /auth/reset-password [post]
func (handler *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	token, rs := validation.GetString(r, "token")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req ResetPasswordRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := ResetPasswordInput{
		Token:       token,
		NewPassword: req.NewPassword,
	}

	ctx := r.Context()
	err := handler.users.ResetPassword(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("password reset successfully").Send(w)
}
