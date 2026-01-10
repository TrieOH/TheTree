package auth

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/authz"
	"GoAuth/internal/application/transactions"
	"GoAuth/internal/application/validation"
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/session"
	"GoAuth/internal/domain/user"
	authport "GoAuth/internal/ports/auth"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbound"
	"context"
	"crypto/ed25519"
	"database/sql"
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
	users         outbound.UserRepository
	sessions      outbound.SessionRepository
	schemas       outbound.SchemaRepository
	versions      outbound.SchemaVersionRepository
	fields        outbound.SchemaFieldsRepository
	projects      outbound.ProjectRepository
	projectUsers  outbound.ProjectUserRepository
	tokenVerifier authport.TokenVerifier
	tx            transactions.TxRunner
}

var _ inbounds.AuthService = (*UseCase)(nil)

func New(
	users outbound.UserRepository,
	sessions outbound.SessionRepository,
	schemas outbound.SchemaRepository,
	versions outbound.SchemaVersionRepository,
	fields outbound.SchemaFieldsRepository,
	projects outbound.ProjectRepository,
	projectUsers outbound.ProjectUserRepository,
	tokenVerifier authport.TokenVerifier,
	tx transactions.TxRunner,
) inbounds.AuthService {
	return &UseCase{
		users:         users,
		sessions:      sessions,
		schemas:       schemas,
		versions:      versions,
		fields:        fields,
		projects:      projects,
		projectUsers:  projectUsers,
		tokenVerifier: tokenVerifier,
		tx:            tx,
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
		return apierr.ErrConflict.WithMsg("error registering user").WithID(apierr.AuthEmailAlreadyUsed).WithCause(errors.New("email already in use"))
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
	accessJTI := uuid.New()
	accessExpiresAt := time.Now().Add(15 * time.Minute)
	accessToken, err = newAccessToken(*u, in.IP, in.Agent, accessJTI.String(), "goauth:v1", sess.SessionID, accessExpiresAt)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}

	var refreshToken string
	refreshToken, err = newRefreshToken("goauth:v1", utils.GoAuthPrivateKey, accessJTI, sess.TokenID, refreshExpiresAt)
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
			span.SetAttributes(attribute.Bool("success", err == nil))
		}
	}()

	var principal *authz.Principal
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	var sess *session.Session
	sess, err = uc.sessions.MarkRevokedByID(ctx, principal.UserID, principal.SessionID)
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.String("session.id", sess.SessionID.String()))

	return nil
}

// Refresh handles the business logic for refreshing a user's tokens.
// It parses the refresh token, checks if it's revoked, and if not,
// determines whether to refresh the tokens for a client or a project user.
func (uc *UseCase) Refresh(ctx context.Context, in inbounds.RefreshInput) (*inbounds.UserTokensOutput, error) {
	txOptions := transactions.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	}

	var out *inbounds.UserTokensOutput
	err := uc.tx.WithinTxWithOptions(ctx, txOptions, func(ctx context.Context) error {
		var err error
		out, err = uc.refreshInternal(ctx, in)
		return err
	})

	return out, err
}

func (uc *UseCase) refreshInternal(ctx context.Context, in inbounds.RefreshInput) (*inbounds.UserTokensOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.Refresh")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("refresh.success", err == nil))
		}
	}()

	var refreshToken *auth.RefreshClaims
	refreshToken, err = uc.tokenVerifier.VerifyRefreshToken(ctx, in.RefreshCookie.Value)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var oldJTI *uuid.UUID
	oldJTI, err = validation.RequireRefreshJTI(span, &refreshToken.ID)
	if err != nil {
		return nil, err
	}

	var (
		newRefreshJTI = uuid.New()
		refreshExp    = time.Now().Add(7 * 24 * time.Hour)
	)

	span.SetAttributes(attribute.String("old_token.id", oldJTI.String()))
	span.SetAttributes(attribute.String("new_token.id", newRefreshJTI.String()))

	sess, err := uc.sessions.RotateToken(
		ctx,
		*oldJTI,
		newRefreshJTI,
		refreshExp,
	)
	if err != nil {
		// sql.ErrNoRows → raced / reused / revoked
		tokenErr := apierr.ErrUnauthorized.WithMsg("refresh token is invalid").WithID(apierr.TokenInvalid)
		return nil, tokenErr
	}

	span.SetAttributes(
		attribute.String("session.id", sess.SessionID.String()),
		attribute.String("session.token_id", sess.TokenID.String()),
		attribute.String("session.user_id", sess.UserID.String()),
		attribute.String("session.user_type", sess.UserType),
	)

	if sess.ProjectID == nil {
		return uc.finishClientRefresh(ctx, sess, in, refreshToken.Sub.AccessJTI, newRefreshJTI, refreshExp)
	}

	return uc.finishProjectUserRefresh(ctx, sess, in, refreshToken.Sub.AccessJTI, newRefreshJTI, refreshExp)
}

func (uc *UseCase) finishClientRefresh(
	ctx context.Context,
	sess *session.Session,
	in inbounds.RefreshInput,
	oldAccessJTI uuid.UUID,
	refreshJTI uuid.UUID,
	refreshExpiresAt time.Time,
) (*inbounds.UserTokensOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.finishClientRefresh")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("success", err == nil))
		}
	}()

	var u *user.User
	u, err = uc.users.GetUserByID(ctx, sess.UserID)
	if err != nil {
		return nil, err
	}

	newAccessJTI := uuid.New()
	if oldAccessJTI.String() == newAccessJTI.String() {
		err = apierr.ErrConflict.WithMsg("new access token ID matched old one, please retry").WithID(apierr.TokenAccessIDMatched)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var accessTokenStr string
	accessExpiresAt := time.Now().Add(15 * time.Minute)
	accessTokenStr, err = newAccessToken(*u, in.IP, in.Agent, newAccessJTI.String(), "goauth:v1", sess.SessionID, accessExpiresAt)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}

	var refreshTokenStr string
	refreshTokenStr, err = newRefreshToken("goauth:v1", utils.GoAuthPrivateKey, newAccessJTI, refreshJTI, refreshExpiresAt)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}

	return &inbounds.UserTokensOutput{
		AccessTokenString:  accessTokenStr,
		RefreshTokenString: refreshTokenStr,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
	}, nil
}

func (uc *UseCase) finishProjectUserRefresh(
	ctx context.Context,
	sess *session.Session,
	in inbounds.RefreshInput,
	oldAccessJTI uuid.UUID,
	refreshJTI uuid.UUID,
	refreshExpiresAt time.Time,
) (*inbounds.UserTokensOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.finishProjectUserRefresh")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("success", err == nil))
		}
	}()

	projectUser, err := uc.projectUsers.GetByIDInternal(ctx, sess.UserID, *sess.ProjectID)
	if err != nil {
		return nil, err
	}

	var privKey string
	privKey, err = uc.projects.GetPrivateKeyByIDInternal(ctx, *sess.ProjectID)
	if err != nil {
		return nil, err
	}

	var decodedKey ed25519.PrivateKey
	decodedKey, err = utils.ParseEd25519PrivateKey(privKey)
	if err != nil {
		return nil, err
	}

	newAccessJTI := uuid.New()
	if oldAccessJTI.String() == newAccessJTI.String() {
		err = apierr.ErrConflict.WithMsg("new access token ID matched old one, please retry").WithID(apierr.TokenAccessIDMatched)
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	var keyID = "project:" + sess.ProjectID.String() + ":v1"
	var accessTokenStr string
	accessExpiresAt := time.Now().Add(15 * time.Minute)
	accessTokenStr, err = newProjectAccessToken(*projectUser, in.IP, in.Agent, newAccessJTI.String(), keyID, sess.SessionID, accessExpiresAt, decodedKey)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}

	var refreshTokenStr string
	refreshTokenStr, err = newRefreshToken(keyID, decodedKey, newAccessJTI, refreshJTI, refreshExpiresAt)
	if err != nil {
		apierr.RecordSystemError(span, err)
		return nil, err
	}

	return &inbounds.UserTokensOutput{
		AccessTokenString:  accessTokenStr,
		RefreshTokenString: refreshTokenStr,
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

	pid, err := validation.RequireProjectID(span, &in.ProjectID)
	if err != nil {
		return err
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
		validatedMetadata, err := uc.validateAndConstructMetadata(ctx, span, *pid, in.SchemaType, in.FlowID, in.CustomFields)
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
		ProjectID:    *pid,
		Email:        in.Email,
		PasswordHash: string(hashedPassword),
		Metadata:     customFields,
	})
	if apierr.IsConflict(err) {
		return apierr.ErrConflict.WithMsg("error registering user").WithID(apierr.AuthEmailAlreadyUsed).WithCause(errors.New("email already in use"))
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

	pid, err := validation.RequireProjectID(span, &in.ProjectID)
	if err != nil {
		return nil, err
	}

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	usr, err := uc.projectUsers.GetByEmailInternal(ctx, *pid, in.Email)
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
		ProjectID: pid,
	})
	if err != nil {
		return nil, err
	}

	var privKey string
	privKey, err = uc.projects.GetPrivateKeyByIDInternal(ctx, *sess.ProjectID)
	if err != nil {
		return nil, err
	}

	var decodedKey ed25519.PrivateKey
	decodedKey, err = utils.ParseEd25519PrivateKey(privKey)
	if err != nil {
		return nil, err
	}

	keyID := "project:" + sess.ProjectID.String() + ":v1"
	accessJTI := uuid.New()
	accessExpiresAt := time.Now().Add(15 * time.Minute)
	accessToken, err := newProjectAccessToken(*usr, in.IP, in.Agent, accessJTI.String(), keyID, sess.SessionID, accessExpiresAt, decodedKey)
	if err != nil {
		apierr.RecordDomainError(span, err)
		return nil, err
	}

	refreshToken, err := newRefreshToken(keyID, decodedKey, accessJTI, sess.TokenID, refreshExpiresAt)
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
