package service

import (
	"context"
	"net/http"
	"strings"
	"time"

	"GoAuth/internal/models"
	"GoAuth/internal/repository"
	"GoAuth/internal/utils"
	"GoAuth/internal/logs"

	"go.uber.org/zap"
	resp "github.com/MintzyG/GoResponse/response"
	"github.com/spf13/viper"
  "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Register(ctx context.Context, req models.RegisterUserRequest) *resp.Response {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		if strings.Contains(err.Error(), "password length exceeds 72 bytes") {
			return resp.BadRequest("error registering user").WithTracePrefix("error").AddTrace("password exceeds 72 char limit")
		}
		return resp.InternalServerError("error hashing user password").WithTracePrefix("error").AddTrace(err)
	}

	_, err = s.queries.RegisterUser(ctx, repository.RegisterUserParams{
		Email:    req.Email,
		Password: string(hashedPassword),
	})

	if err != nil {
		readable := utils.ParseDBError(err)
		if strings.Contains(readable.Error(), "email is already in use") {
			return resp.Conflict("error registering user").WithTracePrefix("error").AddTrace("email already in use")
		}
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
	accessToken, accessJTI, rs := newAccessToken(dbUser)
	if rs != nil {
		return nil, rs
	}
	tokens.AccessTokenString = accessToken

	agent := r.UserAgent()
	ip := utils.GetClientIP(r)
  expires_at := time.Now().Add(7 * 24 * time.Hour)
	refresh_jti := uuid.New()
	refreshToken, rs := newRefreshToken(accessJTI, refresh_jti, agent, ip, expires_at)
	if rs != nil {
		return nil, rs
	}
	tokens.RefreshTokenString = refreshToken

	_, err = s.queries.CreateUserSession(ctx, repository.CreateUserSessionParams{
		TokenID: refresh_jti,
		IssuedAt: time.Now(),
		UserAgent: agent,
		UserIp: ip,
		ExpiresAt: expires_at,
	})

	if err != nil {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		logs.L().Error("Create User Session Failed",
			zap.String("error_value", err.Error()),
			zap.String("request_id", reqID),
			zap.String("user_id", dbUser.ID.String()),
			zap.String("method", r.Method),
			zap.String("path", logs.NormalizePath(r)),
			zap.String("remote_addr", r.RemoteAddr),
		)
	}

	return &tokens, nil
}

func (s *AuthService) Logout(r *http.Request, ctx context.Context) *resp.Response {
	refresh_token_cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return resp.Unauthorized("missing refresh_token cookie")
	}

	refreshClaims, rs := utils.ParseRefreshToken(refresh_token_cookie.Value, viper.GetString("JWT_SECRET"))
	if rs != nil {
		return rs
	}

	jti, err := uuid.Parse(refreshClaims.ID)
	if err != nil {
		return resp.Unauthorized("invalid token ID")
	}

  err = s.queries.DeleteUserSessionByTokenId(ctx, jti)
	if err != nil {
		userID := r.Header.Get("X-User-ID")
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		logs.L().Error("Delete User Session Failed",
			zap.String("error_value", err.Error()),
			zap.String("request_id", reqID),
			zap.String("user_id", userID),
			zap.String("method", r.Method),
			zap.String("path", logs.NormalizePath(r)),
			zap.String("remote_addr", r.RemoteAddr),
		)
	}

	err = s.queries.BlacklistToken(ctx, repository.BlacklistTokenParams{
		TokenID: jti,
		AccessJti: refreshClaims.Sub.AccessJTI,
		ExpiresAt: refreshClaims.ExpiresAt.Time,
	})
	if err != nil {
		readable := utils.ParseDBError(err)
		if strings.Contains(readable.Error(), "duplicate value") {
			return resp.BadRequest("user already logged out").WithTracePrefix("error").AddTrace("token already blacklisted")
		}
		return resp.InternalServerError("error blacklisting token").WithTracePrefix("database-error").AddTrace(readable)
	}

	return nil
}
