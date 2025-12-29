package service

import (
	"GoAuth/internal/apierr"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"GoAuth/internal/models"
	"GoAuth/internal/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) RegisterProjectUser(ctx context.Context, projectID string, req models.RegisterProjectUserRequest) error {
	ctx, span := GoAuthServiceTracer.Start(ctx, "ProjectService.RegisterProjectUser",
		trace.WithAttributes(attribute.String("project.id", projectID)),
	)
	defer span.End()

	pid, err := uuid.Parse(projectID)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if len(req.Password) > 72 {
		return apierr.ErrInvalidInput.WithMsg("password length exceeds 72 bytes").WithID(apierr.AuthInvalidPassword)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		hashErr := apierr.ErrInternal.WithMsg("error hashing user password").WithID(apierr.SystemInternalError).WithCause(err)
		span.SetAttributes(attribute.Bool("password.hashing.failed", true))
		apierr.RecordSystemError(span, hashErr)
		return hashErr
	}

	var user *models.ProjectUser
	user, err = s.projectUserRepo.Register(ctx, models.ProjectUser{
		ProjectID:    pid,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Metadata:     &req.CustomFields,
	})
	if apierr.IsConflict(err) {
		return apierr.ErrConflict.WithMsg("error registering user").WithID(apierr.UserAlreadyExists).WithCause(errors.New("email already in use"))
	} else if err != nil {
		return err
	}

	span.SetAttributes(
		attribute.String("user.id", user.ID.String()),
		attribute.Int64("user.created_at", user.CreatedAt.Unix()),
		attribute.String("user.type", user.UserType),
	)

	return nil
}

func (s *AuthService) LoginProjectUser(r *http.Request, ctx context.Context, projectID string, req models.LoginProjectUserRequest) (*models.UserTokens, error) {
	ctx, span := GoAuthServiceTracer.Start(ctx, "ProjectService.LoginProjectUser",
		trace.WithAttributes(attribute.String("project.id", projectID)),
	)
	defer span.End()

	pid, err := uuid.Parse(projectID)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return nil, apiErr
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	user, err := s.projectUserRepo.GetByEmailInternal(ctx, pid, req.Email)
	if apierr.IsNotFound(err) {
		authErr := apierr.ErrUnauthorized.WithMsg("invalid email or password").WithID(apierr.AuthInvalidCredentials).WithCause(err)
		apierr.RecordDomainError(span, authErr)
		return nil, authErr
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		authErr := apierr.ErrUnauthorized.WithMsg("invalid email or password").WithID(apierr.AuthInvalidCredentials).WithCause(err)
		apierr.RecordDomainError(span, authErr)
		return nil, authErr
	}

	var tokens models.UserTokens
	agent := r.UserAgent()
	ip := utils.GetClientIP(r)

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	session, err := s.sessionRepo.Create(ctx, models.Session{
		IssuedAt:  time.Now(),
		UserAgent: agent,
		UserIp:    ip,
		ExpiresAt: expiresAt,
		UserID:    user.ID,
		ProjectID: &pid,
	})
	if err != nil {
		return nil, err
	}

	accessToken, accessJTI, err := newProjectAccessToken(*user, ip, agent, session.SessionID)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}
	tokens.AccessTokenString = accessToken

	refreshToken, err := newRefreshToken(accessJTI, session.TokenID, expiresAt)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}
	tokens.RefreshTokenString = refreshToken

	return &tokens, nil
}
