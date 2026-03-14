package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/errx"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"encoding/json"
	"net/http"
	"strings"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
)

type AuthHandler struct {
	auth   inbounds.AuthService
	schema inbounds.SchemaService
	redis  outbounds.RedisCacheService
}

func NewAuthHandler(uc inbounds.AuthService, schema inbounds.SchemaService, redis outbounds.RedisCacheService) *AuthHandler {
	return &AuthHandler{
		auth:   uc,
		schema: schema,
		redis:  redis,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user in the system.
// @Tags auth
// @Accept json
// @Produce json
// @Param registerInfo body dto.RegisterUserRequest true "User registration information"
// @Success 201 {object} object "User registered successfully"
// @Failure 400 {object} ErrorResponse "Bad Request: Invalid input"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /auth/register [post]
func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.RegisterUserInput{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := handler.auth.Register(r.Context(), in); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created("Registered user").Send(w)
}

// Login godoc
// @Summary Authenticate a user
// @Description Authenticates a user and sets authentication cookies.
// @Tags auth
// @Accept json
// @Produce json
// @Param loginInfo body dto.LoginUserRequest true "User login information"
// @Success 200 {object} object "Successful authentication"
// @Header 200 {string} Set-Cookie "access_token cookie for authentication"
// @Header 200 {string} Set-Cookie "refresh_token cookie for authentication"
// @Failure 400 {object} ErrorResponse "Bad Request: Invalid input provided"
// @Failure 401 {object} ErrorResponse "Unauthorized: Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /auth/login [post]
func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,

		Agent: r.UserAgent(),
		IP:    ClientIPString(GetClientIP(r, HTTPProxyConfig)),
	}

	tokens, err := handler.auth.Login(r.Context(), in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Logged in").WithData(map[string]any{
		"is_up_to_date": tokens.IsUpToDate,
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
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /auth/logout [post]
func (handler *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	svcCookie, err := r.Cookie("svc_session")
	if err != nil || svcCookie.Value == "" {
		resp.Unauthorized().WithMsg("missing svc_session cookie").Send(w)
		return
	}

	key := "svc_session:" + svcCookie.Value
	data, found, err := handler.redis.GetAny(ctx, key)
	if err != nil || !found {
		resp.Unauthorized().WithMsg("invalid service session").Send(w)
		return
	}

	bytesData, ok := data.([]byte)
	if !ok {
		resp.Unauthorized().WithMsg("invalid session type").Send(w)
		return
	}

	var snapshot authz.ServiceSnapshot
	if err := json.Unmarshal(bytesData, &snapshot); err != nil {
		_ = handler.redis.Delete(ctx, key)
		resp.Unauthorized().WithMsg("failed to unmarshal session").Send(w)
		return
	}

	err = handler.auth.Logout(ctx, snapshot)
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
// @Failure 400 {object} ErrorResponse "Bad Request: Missing or invalid refresh token"
// @Failure 401 {object} ErrorResponse "Unauthorized: Invalid or expired refresh token"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /auth/refresh [post]
func (handler *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		resp.Unauthorized("error getting refresh token").AddTrace(err).Send(w)
		return
	}

	if refreshTokenCookie.Value == "" {
		resp.Unauthorized("missing refresh token value").Send(w)
		return
	}

	in := inbounds.RefreshInput{
		RefreshCookie: refreshTokenCookie,
		Agent:         r.UserAgent(),
		IP:            ClientIPString(GetClientIP(r, HTTPProxyConfig)),
	}

	ctx := r.Context()
	tokens, err := handler.auth.Refresh(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Refreshed tokens").WithData(map[string]any{
		"is_up_to_date": tokens.IsUpToDate,
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
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /.well-known/jwks.json [get]
func (handler *AuthHandler) GetJWKS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jwks, err := handler.auth.GetJWKS(ctx)
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
// @Param registerInfo body dto.RegisterProjectUserRequest true "User registration information for the project"
// @Success 201 {object} object "User registered successfully in project"
// @Failure 400 {object} ErrorResponse "Bad Request: Invalid input or missing project ID"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/register [post]
func (handler *AuthHandler) ProjectRegister(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.RegisterProjectUserRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	schemaType := r.URL.Query().Get("schema_type")
	flowID := r.URL.Query().Get("flow_id")

	in := inbounds.ProjectRegisterInput{
		Email:        req.Email,
		Password:     req.Password,
		CustomFields: req.CustomFields,
		ProjectID:    projectID,
		SchemaType:   schemaType,
		FlowID:       flowID,
	}

	ctx := r.Context()
	if err := handler.auth.RegisterProjectUser(ctx, in); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created("Registered user").Send(w)
}

// ProjectLogin godoc
// @Summary Authenticates a user into a client project
// @Description Authenticates a user within a specific client project and sets authentication cookies.
// @Tags auth
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to login user"
// @Param loginInfo body dto.LoginProjectUserRequest true "User login information for the project"
// @Success 200 {object} object "Successfully authenticated into project"
// @Header 200 {string} Set-Cookie "access_token cookie for authentication"
// @Header 200 {string} Set-Cookie "refresh_token cookie for authentication"
// @Failure 400 {object} ErrorResponse "Bad Request: Invalid input or missing project ID"
// @Failure 401 {object} ErrorResponse "Unauthorized: Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/login [post]
func (handler *AuthHandler) ProjectLogin(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.LoginProjectUserRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	agent := r.UserAgent()
	ip := ClientIPString(GetClientIP(r, HTTPProxyConfig))

	in := inbounds.ProjectLoginInput{
		Email:     req.Email,
		Password:  req.Password,
		ProjectID: projectID,
		IP:        ip,
		Agent:     agent,
	}

	ctx := r.Context()
	tokens, err := handler.auth.LoginProjectUser(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Logged in").WithData(map[string]any{
		"is_up_to_date": tokens.IsUpToDate,
		"access_token":  tokens.AccessTokenString,
		"refresh_token": tokens.RefreshTokenString,
	}).Send(w)
}

// Verify godoc
// @Summary Verify user email
// @Description Verifies a user's email address using a verification token.
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification Token"
// @Success 200 {object} object "User verified successfully"
// @Failure 400 {object} ErrorResponse "Bad Request: Missing or invalid token"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /auth/verify [get]
func (handler *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	token, rs := getString(r, "token")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.auth.Verify(ctx, token)
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
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /auth/verify/resend [post]
func (handler *AuthHandler) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := handler.auth.ResendVerificationEmail(ctx)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("verification email resent successfully").Send(w)
}

// GetUpgradeForm godoc
// @Summary Get upgrade forms for outdated user schemas
// @Description Returns a list of forms that the user needs to complete to update their metadata to the latest schema version.
// @Tags auth
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {array} inbounds.FormResponse "List of upgrade forms"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /projects/{project_id}/upgrade-form [get]
func (handler *AuthHandler) GetUpgradeForm(w http.ResponseWriter, r *http.Request) {
	forms, err := handler.schema.GetUpgradeForm(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(forms).Send(w)
}

// UpdateMetadata godoc
// @Summary Update user metadata
// @Description Updates the user's metadata for a project and validates it against the latest schema version.
// @Tags auth
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param metadata body dto.UpdateMetadataRequest true "Metadata update information"
// @Success 200 {object} object "Metadata updated successfully"
// @Failure 400 {object} ErrorResponse "Bad Request: Invalid input"
// @Failure 403 {object} ErrorResponse "Forbidden: Validation failed"
// @Router /projects/{project_id}/metadata [post]
func (handler *AuthHandler) UpdateMetadata(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateMetadataRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	err := handler.schema.UpdateMetadata(r.Context(), req.CustomFields)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Metadata updated successfully").Send(w)
}

// ForgotPassword godoc
// @Summary Initiate password reset
// @Description Sends a password reset email to the user if the account exists.
// @Tags auth
// @Accept json
// @Produce json
// @Param forgotInfo body dto.ForgotPasswordRequest true "User email and optional project ID"
// @Success 200 {object} object "Forgot password email sent"
// @Failure 400 {object} ErrorResponse "Bad Request: Invalid input"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /auth/forgot-password [post]
func (handler *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ForgotPasswordRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.ForgotPasswordInput{
		Email:     req.Email,
		ProjectID: req.ProjectID,
	}

	ctx := r.Context()
	err := handler.auth.ForgotPassword(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("forgot password email sent").Send(w)
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Resets the user's password using a valid reset token.
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Reset password token"
// @Param resetInfo body dto.ResetPasswordRequest true "New password information"
// @Success 200 {object} object "Password reset successfully"
// @Failure 400 {object} ErrorResponse "Bad Request: Invalid input or token"
// @Failure 401 {object} ErrorResponse "Unauthorized: Invalid or expired token"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /auth/reset-password [post]
func (handler *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	token, rs := getString(r, "token")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.ResetPasswordRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.ResetPasswordInput{
		Token:       token,
		NewPassword: req.NewPassword,
	}

	ctx := r.Context()
	err := handler.auth.ResetPassword(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("password reset successfully").Send(w)
}

// Exchange godoc
// @Summary Exchange global access token
// @Description Exchanges a global access token for a project-scoped session snapshot and tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Global Access Token"
// @Success 200 {object} inbounds.ExchangeOutput "Token exchanged successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/exchange [post]
func (h *AuthHandler) Exchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if len(authHeader) < 7 || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		resp.Unauthorized().WithMsg("missing or invalid Authorization header").Send(w)
		return
	}

	globalAccess := strings.TrimSpace(authHeader[7:])

	out, err := h.auth.Exchange(ctx, globalAccess)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("exchanged").WithData(out).Send(w)
}
