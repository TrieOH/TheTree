package auth

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/authz"
	"GoAuth/internal/application/transactions"
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/domain/revoked_refreshes"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/session"
	"GoAuth/internal/domain/user"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbound"
	"context"
	"encoding/json"
	"errors"
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
	schemas      outbound.SchemaRepository
	versions     outbound.SchemaVersionRepository
	fields       outbound.SchemaFieldsRepository
	projectUsers outbound.ProjectUserRepository
	tx           transactions.TxRunner
}

var _ inbounds.AuthService = (*UseCase)(nil)

func New(
	users outbound.UserRepository,
	sessions outbound.SessionRepository,
	refresh outbound.RevokedRefreshTokenRepository,
	schemas outbound.SchemaRepository,
	versions outbound.SchemaVersionRepository,
	fields outbound.SchemaFieldsRepository,
	projectUsers outbound.ProjectUserRepository,
	tx transactions.TxRunner,
) inbounds.AuthService {
	return &UseCase{
		users:        users,
		sessions:     sessions,
		refresh:      refresh,
		schemas:      schemas,
		versions:     versions,
		fields:       fields,
		projectUsers: projectUsers,
		tx:           tx,
	}
}

// Register handles the business logic for creating a new user.
// It validates the input, hashes the password, and then attempts to create the user in the database.
// It returns an error if the email is already in use or if there is a problem with the database.
func (uc *UseCase) Register(ctx context.Context, in inbounds.RegisterUserInput) error {
	return uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		return uc.registerInternal(ctx, in)
	})
}

func (uc *UseCase) registerInternal(ctx context.Context, in inbounds.RegisterUserInput) error {
	var err error
	ctx, span := usecaseTracer.Start(ctx, "AuthService.Register")
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

// Login handles the business logic for logging in a user.
// It finds the user by email, compares the password, and if successful,
// creates a new session and returns a new set of access and refresh tokens.
func (uc *UseCase) Login(ctx context.Context, in inbounds.LoginUserInput) (*inbounds.UserTokensOutput, error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	var err error
	ctx, span := usecaseTracer.Start(ctx, "AuthService.Login")
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

	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	var sess *session.Session
	sess, err = uc.sessions.Create(ctx, session.Session{
		IssuedAt:  time.Now(),
		UserAgent: in.Agent,
		UserIP:    in.IP,
		ExpiresAt: refreshExpiresAt,
		UserID:    u.ID,
	})

	if err != nil {
		return nil, err
	}

	var accessToken string
	var accessJTI uuid.UUID
	accessExpiresAt := time.Now().Add(15 * time.Minute)
	accessToken, accessJTI, err = newAccessToken(*u, in.IP, in.Agent, sess.SessionID, accessExpiresAt)
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

	return &inbounds.UserTokensOutput{
		AccessTokenString:  accessToken,
		RefreshTokenString: refreshToken,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
	}, nil
}

// Logout handles the business logic for logging out a user.
// It retrieves the principal from the context, deletes the session, and revokes the refresh token.
func (uc *UseCase) Logout(ctx context.Context) error {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.Logout")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("logout.success", err == nil))
		}
	}()

	var principal *authz.Principal
	principal, err = authz.RequirePrincipal(ctx)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return err
	}

	span.SetAttributes(
		attribute.String("user.id", principal.UserID.String()),
		attribute.String("user.type", principal.UserType),
		attribute.String("user.session_id", principal.SessionID.String()),
	)

	if principal.ProjectID != nil {
		span.SetAttributes(
			attribute.String("user.project_id", principal.ProjectID.String()),
		)
	}

	if _, err = uc.sessions.DeleteByFilter(ctx, session.Filter{
		TokenID: &principal.RefreshJTI,
		UserID:  principal.UserID,
	}); err != nil {
		return err
	}

	if err = uc.refresh.Revoke(ctx, revoked_refreshes.RevokedRefreshToken{
		TokenID:   principal.RefreshJTI,
		ExpiresAt: principal.RefreshClaims.ExpiresAt.Time,
	}); err != nil {
		return err
	}

	return nil
}

// Refresh handles the business logic for refreshing a user's tokens.
// It parses the refresh token, checks if it's revoked, and if not,
// determines whether to refresh the tokens for a client or a project user.
func (uc *UseCase) Refresh(ctx context.Context, in inbounds.RefreshInput) (*inbounds.UserTokensOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.Refresh")
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
		tokenErr := apierr.ErrUnauthorized.WithMsg("refresh token revoked").WithID(apierr.TokenRevoked)
		apierr.RecordDomainError(span, tokenErr)
		return nil, tokenErr
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

	var tokens *inbounds.UserTokensOutput
	if sess.ProjectID == nil {
		tokens, err = uc.RefreshClient(ctx, sess, in, jti, refreshToken)
		span.SetAttributes(attribute.String("refresh.flow", "client"))
	} else {
		tokens, err = uc.RefreshProjectUser(ctx, sess, in, jti, refreshToken)
		span.SetAttributes(attribute.String("refresh.flow", "project"))
	}

	return tokens, err
}

// RefreshClient handles the business logic for refreshing a client's tokens.
// It generates a new access and refresh token pair, updates the session, and revokes the old refresh token.
func (uc *UseCase) RefreshClient(
	ctx context.Context,
	sess *session.Session,
	in inbounds.RefreshInput,
	jti uuid.UUID,
	refreshToken *auth.RefreshClaims,
) (
	*inbounds.UserTokensOutput,
	error,
) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.RefreshClient")
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
		UserIP:    in.IP,
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

	return &inbounds.UserTokensOutput{
		AccessTokenString:  newAccessTokenStr,
		RefreshTokenString: newRefreshTokenStr,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
	}, nil
}

// RefreshProjectUser handles the business logic for refreshing a project user's tokens.
// It generates a new access and refresh token pair, updates the session, and revokes the old refresh token.
func (uc *UseCase) RefreshProjectUser(
	ctx context.Context,
	sess *session.Session,
	in inbounds.RefreshInput,
	jti uuid.UUID,
	refreshToken *auth.RefreshClaims,
) (
	*inbounds.UserTokensOutput,
	error,
) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.RefreshProjectUser")
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
		UserIP:    in.IP,
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

	return &inbounds.UserTokensOutput{
		AccessTokenString:  newAccessTokenStr,
		RefreshTokenString: newRefreshTokenStr,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
	}, nil
}

// RegisterProjectUser handles the business logic for creating a new project user.
// It validates the input, hashes the password, and then attempts to create the user in the database.
func (uc *UseCase) RegisterProjectUser(ctx context.Context, in inbounds.ProjectRegisterInput) error {
	return uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		return uc.registerProjectUserInternal(ctx, in)
	})
}

func (uc *UseCase) registerProjectUserInternal(ctx context.Context, in inbounds.ProjectRegisterInput) error {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.RegisterProjectUser",
		trace.WithAttributes(attribute.String("project.id", in.ProjectID)),
	)
	defer span.End()

	if in.FlowID == "" {
		apiErr := apierr.ErrInvalidInput.WithMsg("flow id can't be empty").WithID(apierr.SchemaInvalidFlowID)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}

	if in.SchemaType == "" {
		apiErr := apierr.ErrInvalidInput.WithMsg("schema type can't be empty").WithID(apierr.SchemaInvalidSchemaType)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.FlowID = strings.TrimSpace(strings.ToLower(in.FlowID))
	in.SchemaType = strings.TrimSpace(strings.ToLower(in.SchemaType))

	pid, err := uuid.Parse(in.ProjectID)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}

	if !schema.IsValidSchemaType(in.SchemaType) {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid schema type").WithID(apierr.SchemaInvalidSchemaType)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}

	// FlowIDs cannot be the same as schema types so if this matches we error out
	if schema.IsValidSchemaType(in.FlowID) {
		apiErr := apierr.ErrInvalidInput.WithMsg("flow id can't be the same as a schema type").WithID(apierr.SchemaInvalidFlowID)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}

	if schema.Type(in.SchemaType) == schema.Core && schema.IsFlowIDReserved(in.FlowID) && in.CustomFields != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("custom fields are not allowed for core schema").WithID(apierr.SchemaMetadataNotAllowed)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}

	empty := json.RawMessage(`{}`)
	customFields := &empty

	// Validate and construct metadata for non-core or non-reserved flows
	isCoreWithReservedFlow := schema.Type(in.SchemaType) == schema.Core && schema.IsFlowIDReserved(in.FlowID)
	if !isCoreWithReservedFlow {
		validatedMetadata, err := uc.validateAndConstructMetadata(ctx, span, pid, in.SchemaType, in.FlowID, in.CustomFields)
		if err != nil {
			return err
		}
		if validatedMetadata != nil {
			customFields = validatedMetadata
		}
	}

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
		Metadata:     customFields,
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

// LoginProjectUser handles the business logic for logging in a project user.
// It finds the user by email, compares the password, and if successful,
// creates a new session and returns a new set of access and refresh tokens.
func (uc *UseCase) LoginProjectUser(
	ctx context.Context,
	in inbounds.ProjectLoginInput,
) (
	*inbounds.UserTokensOutput,
	error,
) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.LoginProjectUser",
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
		UserIP:    in.IP,
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

	return &inbounds.UserTokensOutput{
		AccessTokenString:  accessToken,
		RefreshTokenString: refreshToken,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
	}, nil
}
