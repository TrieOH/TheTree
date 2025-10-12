package service

import (
	"context"
	"net/http"
	"strings"

	"GoAuth/internal/models"
	"GoAuth/internal/repository"
	"GoAuth/internal/utils"
	resp "github.com/MintzyG/GoResponse/response"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Register(ctx context.Context, req models.RegisterUserRequest) *resp.Response {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return resp.InternalServerError("error hashing user password").WithTracePrefix("error").AddTrace(err)
	}

	_, err = s.queries.RegisterUser(ctx, repository.RegisterUserParams{
		Email:    req.Email,
		Password: string(hashedPassword),
	})

	if err != nil {
		readable := utils.ParseDBError(err)
		return resp.InternalServerError("error registering user").WithTracePrefix("database-error").AddTrace(readable)
	}

	return nil
}

func (s *AuthService) Login(r *http.Request, ctx context.Context, req models.LoginUserRequest) (*models.UserTokens, *resp.Response) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	dbUser, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		readable := utils.ParseDBError(err)
		if strings.Contains(readable.Error(), "record not found") {
			return nil, resp.Unauthorized("invalid email or password")
		}
		return nil, resp.InternalServerError("error retrieving user").WithTracePrefix("database-error").AddTrace(readable)
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(req.Password))
	if err != nil {
		return nil, resp.Unauthorized("invalid email or password")
	}

	var tokens models.UserTokens
	accessToken, rs := newAccessToken(dbUser)
	if rs != nil {
		return nil, rs
	}
	tokens.AccessTokenString = accessToken

	refreshToken, rs := newRefreshToken()
	if rs != nil {
		return nil, rs
	}
	tokens.RefreshTokenString = refreshToken

	return &tokens, nil
}

func (s *AuthService) Logout(r *http.Request, ctx context.Context) *resp.Response {
	access_token_cookie, err := r.Cookie("access_token")
	if err != nil {
		return resp.Unauthorized("missing access_token cookie")
	}

	refresh_token_cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return resp.Unauthorized("missing refresh_token cookie")
	}

	_, rs := parseAccessToken(access_token_cookie.Value, viper.GetString("JWT_SECRET"))
	if rs != nil && !strings.Contains(rs.Message, "invalid token") {
		return rs
	}
	_, rs = parseRefreshToken(refresh_token_cookie.Value, viper.GetString("JWT_SECRET"))
	if rs != nil {
		return rs
	}

	return nil
}
