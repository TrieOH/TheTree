package service

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"GoAuth/internal/logs"
	"GoAuth/internal/models"
	"GoAuth/internal/repository"
	"GoAuth/internal/utils"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
	agent := r.UserAgent()
	ip := utils.GetClientIP(r)

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	refreshJti := uuid.New()

	session, err := s.queries.CreateUserSession(ctx, repository.CreateUserSessionParams{
		TokenID:   refreshJti,
		IssuedAt:  time.Now(),
		UserAgent: agent,
		UserIp:    ip,
		ExpiresAt: expiresAt,
		UserID:    dbUser.ID,
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
			zap.String("path", utils.NormalizePath(r)),
			zap.String("remote_addr", r.RemoteAddr),
		)
	}

	accessToken, accessJTI, rs := newAccessToken(dbUser, ip, agent, session.SessionID)
	if rs != nil {
		return nil, rs
	}
	tokens.AccessTokenString = accessToken

	refreshToken, rs := newRefreshToken(accessJTI, refreshJti, expiresAt)
	if rs != nil {
		return nil, rs
	}
	tokens.RefreshTokenString = refreshToken

	return &tokens, nil
}

func (s *AuthService) Logout(r *http.Request, ctx context.Context) *resp.Response {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		return resp.Unauthorized("missing refresh_token cookie")
	}

	refreshClaims, rs := utils.ParseRefreshToken(refreshTokenCookie.Value, utils.GoAuthPublicKey)
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
			zap.String("path", utils.NormalizePath(r)),
			zap.String("remote_addr", r.RemoteAddr),
		)
	}

	err = s.queries.BlacklistToken(ctx, repository.BlacklistTokenParams{
		TokenID:   jti,
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

func (s *AuthService) Refresh(r *http.Request, ctx context.Context) (*models.UserTokens, *resp.Response) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		return nil, resp.Unauthorized("missing refresh_token cookie")
	}

	refreshToken, rs := utils.ParseRefreshToken(refreshTokenCookie.Value, utils.GoAuthPublicKey)
	if rs != nil {
		return nil, rs
	}

	jti, err := uuid.Parse(refreshToken.ID)
	if err != nil {
		return nil, resp.Unauthorized("couldn't parse refresh JTI")
	}

	blacklisted, err := s.queries.GetRefreshBlacklistById(ctx, jti)
	if err != nil && !strings.Contains(err.Error(), "no rows") {
		return nil, resp.Unauthorized("couldn't fetch refresh token").WithTracePrefix("database-error").AddTrace(err)
	}

	if blacklisted.TokenID == jti {
		return nil, resp.Unauthorized("refresh token is invalidated")
	}

	session, err := s.queries.GetUserSessionByTokenId(ctx, jti)
	if err != nil {
		return nil, resp.Unauthorized("couldn't fetch user session").WithTracePrefix("database-error").AddTrace(err)
	}

	var dbUser repository.User
	var dbProjectUser repository.ProjectUser
	if session.ProjectID == nil {
		dbUser, err = s.queries.GetUserById(ctx, session.UserID)
		if err != nil {
			return nil, resp.Unauthorized("couldn't fetch user from database").WithTracePrefix("database-error").AddTrace(err)
		}

		var tokens models.UserTokens
		agent := r.UserAgent()
		ip := utils.GetClientIP(r)

		newAccessToken, accessJTI, rs := newAccessToken(dbUser, ip, agent, session.SessionID)
		if rs != nil {
			return nil, rs
		}
		tokens.AccessTokenString = newAccessToken

		expiresAt := time.Now().Add(7 * 24 * time.Hour)
		refreshJti := uuid.New()
		newRefreshToken, rs := newRefreshToken(accessJTI, refreshJti, expiresAt)
		if rs != nil {
			return nil, rs
		}
		tokens.RefreshTokenString = newRefreshToken

		_, err = s.queries.UpdateUserSession(ctx, repository.UpdateUserSessionParams{
			IssuedAt:  time.Now(),
			UserAgent: agent,
			UserIp:    ip,
			ExpiresAt: expiresAt,
			TokenID:   refreshJti,
			SessionID: session.SessionID,
		})

		if err != nil {
			log.Printf("Couldn't update user session: %v", err)
		}

		err = s.queries.BlacklistToken(ctx, repository.BlacklistTokenParams{
			TokenID:   jti,
			ExpiresAt: refreshToken.ExpiresAt.Time,
		})

		if err != nil {
			log.Printf("Couldn't blacklist old token: %v", err)
		}

		return &tokens, nil
	} else {
		dbProjectUser, err = s.queries.GetProjectUserByIdInternal(ctx, repository.GetProjectUserByIdInternalParams{
			ID:        session.UserID,
			ProjectID: *session.ProjectID,
		})
		if err != nil {
			return nil, resp.Unauthorized("couldn't fetch user from database").WithTracePrefix("database-error").AddTrace(err)
		}

		var tokens models.UserTokens
		agent := r.UserAgent()
		ip := utils.GetClientIP(r)

		newAccessToken, accessJTI, rs := newProjectAccessToken(dbProjectUser, ip, agent, session.SessionID)
		if rs != nil {
			return nil, rs
		}
		tokens.AccessTokenString = newAccessToken

		expiresAt := time.Now().Add(7 * 24 * time.Hour)
		refreshJti := uuid.New()
		newRefreshToken, rs := newRefreshToken(accessJTI, refreshJti, expiresAt)
		if rs != nil {
			return nil, rs
		}
		tokens.RefreshTokenString = newRefreshToken

		_, err = s.queries.UpdateUserSession(ctx, repository.UpdateUserSessionParams{
			IssuedAt:  time.Now(),
			UserAgent: agent,
			UserIp:    ip,
			ExpiresAt: expiresAt,
			TokenID:   refreshJti,
			SessionID: session.SessionID,
		})

		if err != nil {
			log.Printf("Couldn't update user session: %v", err)
		}

		err = s.queries.BlacklistToken(ctx, repository.BlacklistTokenParams{
			TokenID:   jti,
			ExpiresAt: refreshToken.ExpiresAt.Time,
		})

		if err != nil {
			log.Printf("Couldn't blacklist old token: %v", err)
		}

		return &tokens, nil
	}
}
