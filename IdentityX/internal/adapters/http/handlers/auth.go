package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/utils"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type AuthHandler struct {
	auth inbounds.AuthService
}

func NewAuthHandler(uc inbounds.AuthService) *AuthHandler {
	return &AuthHandler{auth: uc}
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

	accessCookie := CreateCookie("access_token", tokens.AccessTokenString, tokens.AccessExpiresAt)
	refreshCookie := CreateCookie("refresh_token", tokens.RefreshTokenString, tokens.RefreshExpiresAt)

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	resp.OK("Logged in").Send(w)
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
	err := handler.auth.Logout(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	accessCookie := DeleteCookie("access_token")
	refreshCookie := DeleteCookie("refresh_token")

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

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

	accessCookie := CreateCookie("access_token", tokens.AccessTokenString, tokens.AccessExpiresAt)
	refreshCookie := CreateCookie("refresh_token", tokens.RefreshTokenString, tokens.RefreshExpiresAt)

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	resp.OK("Refreshed tokens").Send(w)
}

// JWKS godoc
// @Summary Exposes the public key using a JWKS endpoint
// @Description Provides the JSON Web Key Set (JWKS) for verifying JWTs issued by the authentication service.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} object "JSON Web Key Set (JWKS)"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /.well-known/jwks.json [get]
func (handler *AuthHandler) JWKS(w http.ResponseWriter, _ *http.Request) {
	jwks := map[string]any{
		"keys": []any{utils.PublicKeyToJWK(utils.GoAuthPublicKey)},
	}

	resp.OK().WithData(jwks).Send(w)
}

// ProjectRegister godoc
// @Summary Register a new user into a client project
// @Description Registers a new user within a specific client project.
// @Tags auth
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to register user"
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

	accessCookie := CreateCookie("access_token", tokens.AccessTokenString, tokens.AccessExpiresAt)
	refreshCookie := CreateCookie("refresh_token", tokens.RefreshTokenString, tokens.RefreshExpiresAt)

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	resp.OK("Logged in").Send(w)
}
