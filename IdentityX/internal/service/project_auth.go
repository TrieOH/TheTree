package service

import (
	"context"
	"net/http"
	"strings"
	"time"

	"GoAuth/internal/logs"
	"GoAuth/internal/models"
	"GoAuth/internal/utils"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) RegisterProjectUser(ctx context.Context, projectId string, req models.RegisterProjectUserRequest) *resp.Response {
	parsedProjectId, err := uuid.Parse(projectId)
	if err != nil {
		return resp.BadRequest("Invalid project ID")
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		if strings.Contains(err.Error(), "password length exceeds 72 bytes") {
			return resp.BadRequest("error registering user").WithTracePrefix("error").AddTrace("password exceeds 72 char limit")
		}
		return resp.InternalServerError("error hashing user password").WithTracePrefix("error").AddTrace(err)
	}

	_, err = s.projectUserRepo.Register(ctx, models.ProjectUser{
		ProjectID: parsedProjectId,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Metadata:  &req.CustomFields,
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

func (s *AuthService) LoginProjectUser(r *http.Request, ctx context.Context, projectId string, req models.LoginProjectUserRequest) (*models.UserTokens, *resp.Response) {
	parsedProjectId, err := uuid.Parse(projectId)
	if err != nil {
		return nil, resp.BadRequest("Invalid project ID")
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	user, err := s.projectUserRepo.GetByEmailInternal(ctx, parsedProjectId, req.Email)
	if err != nil {
		readable := utils.ParseDBError(err)
		if strings.Contains(readable.Error(), "record not found") {
			return nil, resp.Unauthorized("invalid email or password")
		}
		return nil, resp.InternalServerError("error retrieving user").WithTracePrefix("database-error").AddTrace(readable)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, resp.Unauthorized("invalid email or password")
	}

	var tokens models.UserTokens
	agent := r.UserAgent()
	ip := utils.GetClientIP(r)

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	refreshJti := uuid.New()

	session, err := s.sessionRepo.Create(ctx, models.Session{
		TokenID:   refreshJti,
		IssuedAt:  time.Now(),
		UserAgent: agent,
		UserIp:    ip,
		ExpiresAt: expiresAt,
		UserID:    user.ID,
		ProjectID: &parsedProjectId,
	})

	if err != nil {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		logs.L().Error("Create User Session Failed",
			zap.Error(err),
			zap.String("session_token_id", refreshJti.String()),
			zap.String("request_id", reqID),
			zap.String("user_id", user.ID.String()),
			zap.String("method", r.Method),
			zap.String("path", utils.NormalizePath(r)),
			zap.String("remote_addr", r.RemoteAddr),
		)
	}

	accessToken, accessJTI, rs := newProjectAccessToken(*user, ip, agent, session.SessionID)
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
