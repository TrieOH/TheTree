package http

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/utils"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/FastUtilitiesNet/validation"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	uc *auth.UseCase
}

func NewAuthHandler(uc *auth.UseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

// Register godoc
// @Summary Register a new customer
// @Description registers a customer into the system
// @Tags auth
// @Accept json
// @Produce json
// @Param registerInfo body dto.RegisterUserRequest true "register request data"
// @Success 201 {string} string "Registered user"
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	if rs := validation.ValidateInto(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	in := auth.RegisterUserInput{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := ah.uc.Register(r.Context(), in); err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.Created("Registered user").Send(w)
}

// Login godoc
// @Summary Authenticates a customer
// @Description Authenticates a customer of the system
// @Tags auth
// @Accept json
// @Produce json
// @Param loginInfo body dto.LoginUserRequest true "login request data"
// @Success 200 {string} string "Logged in"
// @Header 200 {string} Set-Cookie "access_token cookie for authentication"
// @Header 200 {string} Set-Cookie "refresh_token cookie for authentication"
// @Success 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest
	if rs := validation.ValidateInto(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	in := auth.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	}

	tokens, err := ah.uc.Login(r, r.Context(), in)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	accessCookie := CreateCookie("access_token", tokens.AccessTokenString, tokens.AccessExpiresAt)
	refreshCookie := CreateCookie("refresh_token", tokens.RefreshTokenString, tokens.RefreshExpiresAt)

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	resp.OK("Logged in").Send(w)
}

// Logout godoc
// @Summary Logs out a customer
// @Description Logs out an authenticated customer of the system
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {string} string "Logged out"
// @Header 200 {string} Set-Cookie "clears the access_token cookie"
// @Header 200 {string} Set-Cookie "clears the refresh_token cookie"
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/logout [post]
func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	err := ah.uc.Logout(r.Context())
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	accessCookie := DeleteCookie("access_token")
	refreshCookie := DeleteCookie("refresh_token")

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	resp.OK("Logged out").Send(w)
}

// Refresh godoc
// @Summary Refreshes the user token pair
// @Description Creates a new token pair from a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: refresh_token=yyy"
// @Header 200 {string} Set-Cookie "refresh_token cookie for authentication"
// @Success 200 {string} string "Refreshed tokens"
// @Failure 500 {object} ErrorResponse
// @Router /auth/refresh [post]
func (ah *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		resp.Unauthorized("error getting refresh token").AddTrace(err).Send(w)
		return
	}

	if refreshTokenCookie.Value == "" {
		resp.Unauthorized("missing refresh token value").Send(w)
		return
	}

	in := auth.RefreshInput{
		RefreshCookie: refreshTokenCookie,
		Agent:         r.UserAgent(),
		IP:            utils.GetClientIP(r),
	}

	ctx := r.Context()
	tokens, err := ah.uc.Refresh(ctx, in)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	accessCookie := CreateCookie("access_token", tokens.AccessTokenString, tokens.AccessExpiresAt)
	refreshCookie := CreateCookie("refresh_token", tokens.RefreshTokenString, tokens.RefreshExpiresAt)

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	resp.OK("Refreshed tokens").Send(w)
}

// JWKS godoc
// @Summary Exposes the public key using a JWKS
// @Description Lets users verify the tokens using the public key through a JWKS
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Router /.well-known/jwks.json [get]
func (ah *AuthHandler) JWKS(w http.ResponseWriter, _ *http.Request) {
	jwks := map[string]any{
		"keys": []any{utils.PublicKeyToJWK(utils.GoAuthPublicKey)},
	}

	resp.OK().WithData(jwks).Send(w)
}

// ProjectRegister godoc
// @Summary Register a new user into a client project
// @Description registers a user into the specified project
// @Tags auth
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to register user"
// @Param registerInfo body dto.RegisterProjectUserRequest true "register project user request data"
// @Success 201 {string} string "Registered user"
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/register [post]
func (ah *AuthHandler) ProjectRegister(w http.ResponseWriter, r *http.Request) {
	projectId := chi.URLParam(r, "project_id")
	if projectId == "" {
		resp.BadRequest("missing project id parameter").Send(w)
		return
	}

	var req dto.RegisterProjectUserRequest
	if rs := validation.ValidateInto(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	in := auth.ProjectRegisterInput{
		Email:        req.Email,
		Password:     req.Password,
		CustomFields: req.CustomFields,
		ProjectID:    projectId,
	}

	ctx := r.Context()
	if err := ah.uc.RegisterProjectUser(ctx, in); err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.Created("Registered user").Send(w)
}

// ProjectLogin godoc
// @Summary Authenticates a user into a client project
// @Description Authenticates a user into the specified client project
// @Tags auth
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to login user"
// @Param loginInfo body dto.LoginProjectUserRequest true "login project user request data"
// @Success 200 {string} string "Logged in"
// @Header 200 {string} Set-Cookie "access_token cookie for authentication"
// @Header 200 {string} Set-Cookie "refresh_token cookie for authentication"
// @Success 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/login [post]
func (ah *AuthHandler) ProjectLogin(w http.ResponseWriter, r *http.Request) {
	projectId := chi.URLParam(r, "project_id")
	if projectId == "" {
		resp.BadRequest("missing project id parameter").Send(w)
		return
	}

	var req dto.LoginProjectUserRequest
	if rs := validation.ValidateInto(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	agent := r.UserAgent()
	ip := utils.GetClientIP(r)

	in := auth.ProjectLoginInput{
		Email:     req.Email,
		Password:  req.Password,
		ProjectID: projectId,
		IP:        ip,
		Agent:     agent,
	}

	ctx := r.Context()
	tokens, err := ah.uc.LoginProjectUser(ctx, in)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	accessCookie := CreateCookie("access_token", tokens.AccessTokenString, tokens.AccessExpiresAt)
	refreshCookie := CreateCookie("refresh_token", tokens.RefreshTokenString, tokens.RefreshExpiresAt)

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	resp.OK("Logged in").Send(w)
}
