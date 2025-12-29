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

func (s *AuthService) Register(ctx context.Context, req models.RegisterUserRequest) error {
	var err error
	ctx, span := GoAuthServiceTracer.Start(ctx, "Auth.Register")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("register.success", err == nil))
		}
	}()

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if len(req.Password) > 72 {
		return apierr.ErrInvalidInput.WithMsg("password length exceeds 72 bytes").WithID(apierr.AuthInvalidPassword)
	}

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		hashErr := apierr.ErrInternal.WithMsg("error hashing user password").WithID(apierr.SystemInternalError).WithCause(err)
		span.SetAttributes(attribute.Bool("password.hashing.failed", true))
		apierr.RecordSystemError(span, hashErr)
		return hashErr
	}

	var user *models.User
	user, err = s.userRepo.Register(ctx, req.Email, string(hashedPassword))
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

func (s *AuthService) Login(r *http.Request, ctx context.Context, req models.LoginUserRequest) (*models.UserTokens, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	var err error
	ctx, span := GoAuthServiceTracer.Start(ctx, "Auth.Login")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("login.success", err == nil))
		}
	}()

	var user *models.User
	user, err = s.userRepo.GetUserByEmail(ctx, req.Email)
	if apierr.IsNotFound(err) {
		authErr := apierr.ErrUnauthorized.WithMsg("invalid email or password").WithID(apierr.AuthInvalidCredentials).WithCause(err)
		apierr.RecordDomainError(span, authErr)
		return nil, authErr
	} else if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("user.id", user.ID.String()),
		attribute.String("user.type", user.UserType),
		attribute.Int64("user.created_at_unix", user.CreatedAt.Unix()),
	)

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

	var session *models.Session
	session, err = s.sessionRepo.Create(ctx, models.Session{
		IssuedAt:  time.Now(),
		UserAgent: agent,
		UserIp:    ip,
		ExpiresAt: expiresAt,
		UserID:    user.ID,
	})

	if err != nil {
		return nil, err
	}

	var accessToken string
	var accessJTI uuid.UUID
	accessToken, accessJTI, err = newAccessToken(*user, ip, agent, session.SessionID)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}
	tokens.AccessTokenString = accessToken

	var refreshToken string
	refreshToken, err = newRefreshToken(accessJTI, session.TokenID, expiresAt)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}
	tokens.RefreshTokenString = refreshToken

	return &tokens, nil
}

func (s *AuthService) Logout(ctx context.Context) error {
	ctx, span := GoAuthServiceTracer.Start(ctx, "Auth.Logout")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("logout.success", err == nil))
		}
	}()

	var accessClaims *models.AccessClaims
	accessClaims, err = models.GetAccessClaims(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return err
	}

	span.SetAttributes(
		attribute.String("user.id", accessClaims.Sub.ID.String()),
		attribute.String("user.type", accessClaims.Sub.UserType),
		attribute.String("user.session_id", accessClaims.Sub.SessionID.String()),
	)

	if accessClaims.Sub.ProjectID != nil {
		span.SetAttributes(
			attribute.String("user.project_id", accessClaims.Sub.ProjectID.String()),
		)
	}

	var refreshClaims *models.RefreshClaims
	refreshClaims, err = models.GetRefreshClaims(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return err
	}

	var jti uuid.UUID
	jti, err = uuid.Parse(refreshClaims.ID)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return apierr.ErrUnauthorized.WithMsg("unable to parse refresh ID").WithID(apierr.TokenInvalidID)
	}

	if _, err = s.sessionRepo.DeleteByFilter(ctx, models.SessionFilter{
		TokenID: &jti,
		UserID:  accessClaims.Sub.ID,
	}); err != nil {
		return err
	}

	if err = s.revokedRefreshTokensRepo.Revoke(ctx, models.RevokedRefreshToken{
		TokenID:   jti,
		ExpiresAt: refreshClaims.ExpiresAt.Time,
	}); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Refresh(payload models.RefreshData, ctx context.Context) (*models.UserTokens, error) {
	ctx, span := GoAuthServiceTracer.Start(ctx, "Auth.Refresh")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("refresh.success", err == nil))
		}
	}()

	var refreshToken *models.RefreshClaims
	refreshToken, err = utils.ParseRefreshToken(payload.RefreshCookie.Value, utils.GoAuthPublicKey)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var jti uuid.UUID
	jti, err = uuid.Parse(refreshToken.ID)
	if err != nil {
		tokenErr := apierr.ErrInvalidInput.WithMsg("unable to parse refresh token ID").WithID(apierr.TokenInvalidID)
		apierr.RecordDomainError(span, tokenErr)
		return nil, tokenErr
	}

	span.SetAttributes(attribute.String("refresh_token.id", jti.String()))

	var isRevoked bool
	isRevoked, err = s.revokedRefreshTokensRepo.IsRevoked(ctx, jti)
	if err != nil {
		return nil, err
	}

	if isRevoked {
		revoked, err := s.revokedRefreshTokensRepo.GetByID(ctx, jti)

		if err != nil {
			return nil, err
		}

		if revoked.TokenID == jti {
			tokenErr := apierr.ErrUnauthorized.WithMsg("refresh token revoked").WithID(apierr.TokenRevoked)
			apierr.RecordDomainError(span, tokenErr)
			return nil, tokenErr
		}
	}

	var session *models.Session
	session, err = s.sessionRepo.GetByTokenID(ctx, jti)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("session.id", session.SessionID.String()),
		attribute.String("session.token_id", session.TokenID.String()),
		attribute.String("session.user_id", session.UserID.String()),
		attribute.String("session.user_type", session.UserType),
	)

	if session.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", session.ProjectID.String()))
	}

	var tokens *models.UserTokens
	if session.ProjectID == nil {
		tokens, err = s.RefreshClient(ctx, session, payload, jti, refreshToken)
		span.SetAttributes(attribute.String("refresh.flow", "client"))
	} else {
		tokens, err = s.RefreshProjectUser(ctx, session, payload, jti, refreshToken)
		span.SetAttributes(attribute.String("refresh.flow", "project"))
	}

	return tokens, err
}

func (s *AuthService) RefreshClient(ctx context.Context, session *models.Session, payload models.RefreshData, jti uuid.UUID, refreshToken *models.RefreshClaims) (*models.UserTokens, error) {
	ctx, span := GoAuthServiceTracer.Start(ctx, "Auth.RefreshClient")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("refresh_client.success", err == nil))
		}
	}()

	span.SetAttributes(
		attribute.String("refresh_client.session.token_id", session.TokenID.String()),
		attribute.String("refresh_client.session.id", session.SessionID.String()),
		attribute.String("refresh_client.session.user_id", session.UserID.String()),
		attribute.String("refresh_client.session.user_type", session.UserType),
	)

	if session.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", session.ProjectID.String()))
	}

	var user *models.User
	user, err = s.userRepo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	var tokens models.UserTokens

	var newAccessTokenStr string
	var accessJTI uuid.UUID
	newAccessTokenStr, accessJTI, err = newAccessToken(*user, payload.IP, payload.Agent, session.SessionID)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}
	tokens.AccessTokenString = newAccessTokenStr

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	refreshJti := uuid.New()

	var newRefreshTokenStr string
	newRefreshTokenStr, err = newRefreshToken(accessJTI, refreshJti, expiresAt)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}
	tokens.RefreshTokenString = newRefreshTokenStr

	if err = s.sessionRepo.Update(ctx, models.Session{
		IssuedAt:  time.Now(),
		UserAgent: payload.Agent,
		UserIp:    payload.IP,
		ExpiresAt: expiresAt,
		TokenID:   refreshJti,
		SessionID: session.SessionID,
	}); err != nil {
		return nil, err
	}

	if err = s.revokedRefreshTokensRepo.Revoke(ctx, models.RevokedRefreshToken{
		TokenID:   jti,
		ExpiresAt: refreshToken.ExpiresAt.Time,
	}); err != nil {
		return nil, err
	}

	span.AddEvent("refresh_client.revoked",
		trace.WithAttributes(
			attribute.Int64("revoked.expires_at", refreshToken.ExpiresAt.Time.Unix()),
			attribute.String("revoked.token_id", jti.String()),
		),
	)

	return &tokens, nil
}

func (s *AuthService) RefreshProjectUser(ctx context.Context, session *models.Session, payload models.RefreshData, jti uuid.UUID, refreshToken *models.RefreshClaims) (*models.UserTokens, error) {
	ctx, span := GoAuthServiceTracer.Start(ctx, "Auth.RefreshProjectUser")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("refresh_project_user.success", err == nil))
		}
	}()

	span.SetAttributes(
		attribute.String("refresh_project_user.session.token_id", session.TokenID.String()),
		attribute.String("refresh_project_user.session.id", session.SessionID.String()),
		attribute.String("refresh_project_user.session.user_id", session.UserID.String()),
		attribute.String("refresh_project_user.session.user_type", session.UserType),
	)

	if session.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", session.ProjectID.String()))
	}

	var projectUser *models.ProjectUser
	projectUser, err = s.projectUserRepo.GetByIDInternal(ctx, session.UserID, *session.ProjectID)
	if err != nil {
		return nil, err
	}

	var tokens models.UserTokens

	var newAccessTokenStr string
	var accessJTI uuid.UUID
	newAccessTokenStr, accessJTI, err = newProjectAccessToken(*projectUser, payload.IP, payload.Agent, session.SessionID)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}
	tokens.AccessTokenString = newAccessTokenStr

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	refreshJti := uuid.New()

	var newRefreshTokenStr string
	newRefreshTokenStr, err = newRefreshToken(accessJTI, refreshJti, expiresAt)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}
	tokens.RefreshTokenString = newRefreshTokenStr

	if err = s.sessionRepo.Update(ctx, models.Session{
		IssuedAt:  time.Now(),
		UserAgent: payload.Agent,
		UserIp:    payload.IP,
		ExpiresAt: expiresAt,
		TokenID:   refreshJti,
		SessionID: session.SessionID,
	}); err != nil {
		return nil, err
	}

	if err = s.revokedRefreshTokensRepo.Revoke(ctx, models.RevokedRefreshToken{
		TokenID:   jti,
		ExpiresAt: refreshToken.ExpiresAt.Time,
	}); err != nil {
		return nil, err
	}

	span.AddEvent("refresh_project_user.revoked",
		trace.WithAttributes(
			attribute.Int64("revoked.expires_at", refreshToken.ExpiresAt.Time.Unix()),
			attribute.String("revoked.token_id", jti.String()),
		),
	)

	return &tokens, nil
}
