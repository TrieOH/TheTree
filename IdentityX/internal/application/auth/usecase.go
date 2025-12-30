package auth

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/domain/revoked_refreshes"
	"GoAuth/internal/domain/session"
	"GoAuth/internal/domain/user"
	"GoAuth/internal/ports/outbound"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"GoAuth/internal/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

var (
	usecaseTracer = otel.Tracer("auth_usecase")
)

type UseCase struct {
	users        outbound.UserRepository
	refresh      outbound.RevokedRefreshTokenRepository
	sessions     outbound.SessionRepository
	projectUsers outbound.ProjectUserRepository
}

func New(
	users outbound.UserRepository,
	sessions outbound.SessionRepository,
	refresh outbound.RevokedRefreshTokenRepository,
	projectUsers outbound.ProjectUserRepository,
) *UseCase {
	return &UseCase{
		users:        users,
		sessions:     sessions,
		refresh:      refresh,
		projectUsers: projectUsers,
	}
}

func (uc *UseCase) Register(ctx context.Context, in RegisterUserInput) error {
	var err error
	ctx, span := usecaseTracer.Start(ctx, "Auth.Create")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("register.success", err == nil))
		}
	}()

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	if len(in.Password) > 72 {
		return apierr.ErrInvalidInput.WithMsg("password length exceeds 72 bytes").WithID(apierr.AuthInvalidPassword)
	}

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		hashErr := apierr.ErrInternal.WithMsg("error hashing user password").WithID(apierr.SystemInternalError).WithCause(err)
		span.SetAttributes(attribute.Bool("password.hashing.failed", true))
		apierr.RecordSystemError(span, hashErr)
		return hashErr
	}

	var u *user.User
	u, err = uc.users.Register(ctx, in.Email, string(hashedPassword))
	if apierr.IsConflict(err) {
		return apierr.ErrConflict.WithMsg("error registering user").WithID(apierr.UserAlreadyExists).WithCause(errors.New("email already in use"))
	} else if err != nil {
		return err
	}

	span.SetAttributes(
		attribute.String("user.id", u.ID.String()),
		attribute.Int64("user.created_at", u.CreatedAt.Unix()),
		attribute.String("user.type", u.UserType),
	)

	return nil
}

func (uc *UseCase) Login(r *http.Request, ctx context.Context, in LoginUserInput) (*UserTokensOutput, error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	var err error
	ctx, span := usecaseTracer.Start(ctx, "Auth.Login")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("login.success", err == nil))
		}
	}()

	var u *user.User
	u, err = uc.users.GetUserByEmail(ctx, in.Email)
	if apierr.IsNotFound(err) {
		authErr := apierr.ErrUnauthorized.WithMsg("invalid email or password").WithID(apierr.AuthInvalidCredentials).WithCause(err)
		apierr.RecordDomainError(span, authErr)
		return nil, authErr
	} else if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("user.id", u.ID.String()),
		attribute.String("user.type", u.UserType),
		attribute.Int64("user.created_at_unix", u.CreatedAt.Unix()),
	)

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password))
	if err != nil {
		authErr := apierr.ErrUnauthorized.WithMsg("invalid email or password").WithID(apierr.AuthInvalidCredentials).WithCause(err)
		apierr.RecordDomainError(span, authErr)
		return nil, authErr
	}

	agent := r.UserAgent()
	ip := utils.GetClientIP(r)

	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	var sess *session.Session
	sess, err = uc.sessions.Create(ctx, session.Session{
		IssuedAt:  time.Now(),
		UserAgent: agent,
		UserIp:    ip,
		ExpiresAt: refreshExpiresAt,
		UserID:    u.ID,
	})

	if err != nil {
		return nil, err
	}

	var accessToken string
	var accessJTI uuid.UUID
	accessExpiresAt := time.Now().Add(15 * time.Minute)
	accessToken, accessJTI, err = newAccessToken(*u, ip, agent, sess.SessionID, accessExpiresAt)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}

	var refreshToken string
	refreshToken, err = newRefreshToken(accessJTI, sess.TokenID, refreshExpiresAt)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}

	return &UserTokensOutput{
		AccessTokenString:  accessToken,
		RefreshTokenString: refreshToken,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
	}, nil
}

func (uc *UseCase) Logout(ctx context.Context) error {
	ctx, span := usecaseTracer.Start(ctx, "Auth.Logout")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("logout.success", err == nil))
		}
	}()

	var accessClaims *auth.AccessClaims
	accessClaims, err = auth.GetAccessClaims(ctx)
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

	var refreshClaims *auth.RefreshClaims
	refreshClaims, err = auth.GetRefreshClaims(ctx)
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

	if _, err = uc.sessions.DeleteByFilter(ctx, session.SessionFilter{
		TokenID: &jti,
		UserID:  accessClaims.Sub.ID,
	}); err != nil {
		return err
	}

	if err = uc.refresh.Revoke(ctx, revoked_refreshes.RevokedRefreshToken{
		TokenID:   jti,
		ExpiresAt: refreshClaims.ExpiresAt.Time,
	}); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) Refresh(ctx context.Context, in RefreshInput) (*UserTokensOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "Auth.Refresh")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("refresh.success", err == nil))
		}
	}()

	var refreshToken *auth.RefreshClaims
	refreshToken, err = utils.ParseRefreshToken(in.RefreshCookie.Value, utils.GoAuthPublicKey)
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
	isRevoked, err = uc.refresh.IsRevoked(ctx, jti)
	if err != nil {
		return nil, err
	}

	if isRevoked {
		revoked, err := uc.refresh.GetByID(ctx, jti)

		if err != nil {
			return nil, err
		}

		if revoked.TokenID == jti {
			tokenErr := apierr.ErrUnauthorized.WithMsg("refresh token revoked").WithID(apierr.TokenRevoked)
			apierr.RecordDomainError(span, tokenErr)
			return nil, tokenErr
		}
	}

	var sess *session.Session
	sess, err = uc.sessions.GetByTokenID(ctx, jti)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("session.id", sess.SessionID.String()),
		attribute.String("session.token_id", sess.TokenID.String()),
		attribute.String("session.user_id", sess.UserID.String()),
		attribute.String("session.user_type", sess.UserType),
	)

	if sess.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", sess.ProjectID.String()))
	}

	var tokens *UserTokensOutput
	if sess.ProjectID == nil {
		tokens, err = uc.RefreshClient(ctx, sess, in, jti, refreshToken)
		span.SetAttributes(attribute.String("refresh.flow", "client"))
	} else {
		tokens, err = uc.RefreshProjectUser(ctx, sess, in, jti, refreshToken)
		span.SetAttributes(attribute.String("refresh.flow", "project"))
	}

	return tokens, err
}

func (uc *UseCase) RefreshClient(
	ctx context.Context,
	sess *session.Session,
	in RefreshInput,
	jti uuid.UUID,
	refreshToken *auth.RefreshClaims,
) (
	*UserTokensOutput,
	error,
) {
	ctx, span := usecaseTracer.Start(ctx, "Auth.RefreshClient")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("refresh_client.success", err == nil))
		}
	}()

	span.SetAttributes(
		attribute.String("refresh_client.session.token_id", sess.TokenID.String()),
		attribute.String("refresh_client.session.id", sess.SessionID.String()),
		attribute.String("refresh_client.session.user_id", sess.UserID.String()),
		attribute.String("refresh_client.session.user_type", sess.UserType),
	)

	if sess.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", sess.ProjectID.String()))
	}

	var u *user.User
	u, err = uc.users.GetUserByID(ctx, sess.UserID)
	if err != nil {
		return nil, err
	}

	var newAccessTokenStr string
	var accessJTI uuid.UUID
	accessExpiresAt := time.Now().Add(15 * time.Minute)
	newAccessTokenStr, accessJTI, err = newAccessToken(*u, in.IP, in.Agent, sess.SessionID, accessExpiresAt)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}

	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)
	refreshJti := uuid.New()

	var newRefreshTokenStr string
	newRefreshTokenStr, err = newRefreshToken(accessJTI, refreshJti, refreshExpiresAt)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}

	if err = uc.sessions.Update(ctx, session.Session{
		IssuedAt:  time.Now(),
		UserAgent: in.Agent,
		UserIp:    in.IP,
		ExpiresAt: refreshExpiresAt,
		TokenID:   refreshJti,
		SessionID: sess.SessionID,
	}); err != nil {
		return nil, err
	}

	if err = uc.refresh.Revoke(ctx, revoked_refreshes.RevokedRefreshToken{
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

	return &UserTokensOutput{
		AccessTokenString:  newAccessTokenStr,
		RefreshTokenString: newRefreshTokenStr,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
	}, nil
}

func (uc *UseCase) RefreshProjectUser(
	ctx context.Context,
	sess *session.Session,
	in RefreshInput,
	jti uuid.UUID,
	refreshToken *auth.RefreshClaims,
) (
	*UserTokensOutput,
	error,
) {
	ctx, span := usecaseTracer.Start(ctx, "Auth.RefreshProjectUser")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("refresh_project_user.success", err == nil))
		}
	}()

	span.SetAttributes(
		attribute.String("refresh_project_user.session.token_id", sess.TokenID.String()),
		attribute.String("refresh_project_user.session.id", sess.SessionID.String()),
		attribute.String("refresh_project_user.session.user_id", sess.UserID.String()),
		attribute.String("refresh_project_user.session.user_type", sess.UserType),
	)

	if sess.ProjectID != nil {
		span.SetAttributes(attribute.String("session.project_id", sess.ProjectID.String()))
	}

	var projectUser *project_users.ProjectUser
	projectUser, err = uc.projectUsers.GetByIDInternal(ctx, sess.UserID, *sess.ProjectID)
	if err != nil {
		return nil, err
	}

	var newAccessTokenStr string
	var accessJTI uuid.UUID
	accessExpiresAt := time.Now().Add(15 * time.Minute)
	newAccessTokenStr, accessJTI, err = newProjectAccessToken(*projectUser, in.IP, in.Agent, sess.SessionID, accessExpiresAt)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}

	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)
	refreshJti := uuid.New()

	var newRefreshTokenStr string
	newRefreshTokenStr, err = newRefreshToken(accessJTI, refreshJti, refreshExpiresAt)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}

	if err = uc.sessions.Update(ctx, session.Session{
		IssuedAt:  time.Now(),
		UserAgent: in.Agent,
		UserIp:    in.IP,
		ExpiresAt: refreshExpiresAt,
		TokenID:   refreshJti,
		SessionID: sess.SessionID,
	}); err != nil {
		return nil, err
	}

	if err = uc.refresh.Revoke(ctx, revoked_refreshes.RevokedRefreshToken{
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

	return &UserTokensOutput{
		AccessTokenString:  newAccessTokenStr,
		RefreshTokenString: newRefreshTokenStr,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
	}, nil
}

func (uc *UseCase) RegisterProjectUser(ctx context.Context, in ProjectRegisterInput) error {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.RegisterProjectUser",
		trace.WithAttributes(attribute.String("project.id", in.ProjectID)),
	)
	defer span.End()

	pid, err := uuid.Parse(in.ProjectID)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	if len(in.Password) > 72 {
		return apierr.ErrInvalidInput.WithMsg("password length exceeds 72 bytes").WithID(apierr.AuthInvalidPassword)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		hashErr := apierr.ErrInternal.WithMsg("error hashing user password").WithID(apierr.SystemInternalError).WithCause(err)
		span.SetAttributes(attribute.Bool("password.hashing.failed", true))
		apierr.RecordSystemError(span, hashErr)
		return hashErr
	}

	var usr *project_users.ProjectUser
	usr, err = uc.projectUsers.Register(ctx, project_users.ProjectUser{
		ProjectID:    pid,
		Email:        in.Email,
		PasswordHash: string(hashedPassword),
		Metadata:     &in.CustomFields,
	})
	if apierr.IsConflict(err) {
		return apierr.ErrConflict.WithMsg("error registering user").WithID(apierr.UserAlreadyExists).WithCause(errors.New("email already in use"))
	} else if err != nil {
		return err
	}

	span.SetAttributes(
		attribute.String("user.id", usr.ID.String()),
		attribute.Int64("user.created_at", usr.CreatedAt.Unix()),
		attribute.String("user.type", usr.UserType),
	)

	return nil
}

func (uc *UseCase) LoginProjectUser(
	ctx context.Context,
	in ProjectLoginInput,
) (
	*UserTokensOutput,
	error,
) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.LoginProjectUser",
		trace.WithAttributes(attribute.String("project.id", in.ProjectID)),
	)
	defer span.End()

	pid, err := uuid.Parse(in.ProjectID)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return nil, apiErr
	}

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	usr, err := uc.projectUsers.GetByEmailInternal(ctx, pid, in.Email)
	if apierr.IsNotFound(err) {
		authErr := apierr.ErrUnauthorized.WithMsg("invalid email or password").WithID(apierr.AuthInvalidCredentials).WithCause(err)
		apierr.RecordDomainError(span, authErr)
		return nil, authErr
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(in.Password))
	if err != nil {
		authErr := apierr.ErrUnauthorized.WithMsg("invalid email or password").WithID(apierr.AuthInvalidCredentials).WithCause(err)
		apierr.RecordDomainError(span, authErr)
		return nil, authErr
	}

	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)
	sess, err := uc.sessions.Create(ctx, session.Session{
		IssuedAt:  time.Now(),
		UserAgent: in.Agent,
		UserIp:    in.IP,
		ExpiresAt: refreshExpiresAt,
		UserID:    usr.ID,
		ProjectID: &pid,
	})
	if err != nil {
		return nil, err
	}

	accessExpiresAt := time.Now().Add(15 * time.Minute)
	accessToken, accessJTI, err := newProjectAccessToken(*usr, in.IP, in.Agent, sess.SessionID, accessExpiresAt)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	refreshToken, err := newRefreshToken(accessJTI, sess.TokenID, refreshExpiresAt)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	return &UserTokensOutput{
		AccessTokenString:  accessToken,
		RefreshTokenString: refreshToken,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
	}, nil
}
