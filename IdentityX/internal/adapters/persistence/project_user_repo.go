package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/ports/outbounds"
	"context"
	"encoding/json"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type projectUserRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

func (repo *projectUserRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(transactions.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ outbounds.ProjectUserRepository = (*projectUserRepo)(nil)

func NewProjectUserRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) outbounds.ProjectUserRepository {
	return &projectUserRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func mapProjectUserFromDB(dst *project_users.ProjectUser, src *sqlc.ProjectUser) {
	dst.ID = src.ID
	dst.ProjectID = src.ProjectID
	dst.Email = src.Email
	dst.PasswordHash = src.PasswordHash
	dst.UserType = src.UserType
	dst.Metadata = &src.Metadata
	dst.IsActive = src.IsActive
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
	dst.LastLoginAt = src.LastLoginAt
	dst.IsVerified = src.IsVerified
	dst.VerifiedAt = src.VerifiedAt
}

func (repo *projectUserRepo) Register(ctx context.Context, toRegister project_users.ProjectUser) (*project_users.ProjectUser, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectUserRepo.Register",
		trace.WithAttributes(
			attribute.String("user.project_id", toRegister.ProjectID.String()),
		),
	)
	defer span.End()

	var metadata json.RawMessage
	if toRegister.Metadata != nil {
		metadata = *toRegister.Metadata
	} else {
		metadata = json.RawMessage("{}")
	}

	sqlcUser, err := repo.queries(ctx).RegisterProjectUser(ctx, sqlc.RegisterProjectUserParams{
		ProjectID:    toRegister.ProjectID,
		Email:        toRegister.Email,
		PasswordHash: toRegister.PasswordHash,
		Metadata:     metadata,
	})
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("project_user.id", sqlcUser.ID.String()),
		attribute.String("project_user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	mapProjectUserFromDB(&toRegister, &sqlcUser)

	return &toRegister, nil
}

func (repo *projectUserRepo) GetByIDExternal(ctx context.Context, projectUserID, projectID, ownerID uuid.UUID) (*project_users.ProjectUser, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectUserRepo.GetByIDExternal",
		trace.WithAttributes(
			attribute.String("project_user.id", projectUserID.String()),
			attribute.String("project_user.project_id", projectID.String()),
			attribute.String("project.owner_id", ownerID.String()),
		),
	)
	defer span.End()

	sqlcUser, err := repo.queries(ctx).GetProjectUserById(ctx, sqlc.GetProjectUserByIdParams{
		ID:        projectUserID,
		ProjectID: projectID,
		OwnerID:   ownerID,
	})
	if err != nil {
		return nil, fail.From(err).WithArgs("project user").RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("project_user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
		attribute.Bool("user.is_active", sqlcUser.IsActive),
	)

	var user project_users.ProjectUser
	mapProjectUserFromDB(&user, &sqlcUser)

	return &user, nil
}

func (repo *projectUserRepo) GetByIDInternal(ctx context.Context, projectUserID, projectID uuid.UUID) (*project_users.ProjectUser, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectUserRepo.GetByIDInternal",
		trace.WithAttributes(
			attribute.String("project_user.project_id", projectID.String()),
			attribute.String("project_user.id", projectUserID.String()),
		),
	)
	defer span.End()

	sqlcUser, err := repo.queries(ctx).GetProjectUserByIdInternal(ctx, sqlc.GetProjectUserByIdInternalParams{
		ID:        projectUserID,
		ProjectID: projectID,
	})
	if err != nil {
		return nil, fail.From(err).WithArgs("project user").RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
		attribute.Bool("user.is_active", sqlcUser.IsActive),
	)

	var user project_users.ProjectUser
	mapProjectUserFromDB(&user, &sqlcUser)

	return &user, nil
}

func (repo *projectUserRepo) GetByEmailExternal(ctx context.Context, projectID uuid.UUID, email string, ownerID uuid.UUID) (*project_users.ProjectUser, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectUserRepo.GetByEmailExternal",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project_user.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcUser, err := repo.queries(ctx).GetProjectUserByEmailExternal(ctx, sqlc.GetProjectUserByEmailExternalParams{
		ProjectID: projectID,
		Email:     email,
		OwnerID:   ownerID,
	})
	if err != nil {
		return nil, fail.From(err).WithArgs("project user").RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("project_user.id", sqlcUser.ID.String()),
		attribute.String("project_user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	var user project_users.ProjectUser
	mapProjectUserFromDB(&user, &sqlcUser)

	return &user, nil
}

func (repo *projectUserRepo) GetByEmailInternal(ctx context.Context, projectID uuid.UUID, email string) (*project_users.ProjectUser, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectUserRepo.GetByEmailInternal",
		trace.WithAttributes(
			attribute.String("project_user.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcUser, err := repo.queries(ctx).GetProjectUserByEmailInternal(ctx, sqlc.GetProjectUserByEmailInternalParams{
		ProjectID: projectID,
		Email:     email,
	})
	if err != nil {
		return nil, fail.From(err).WithArgs("project user").RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("project_user.id", sqlcUser.ID.String()),
		attribute.String("project_user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	var user project_users.ProjectUser
	mapProjectUserFromDB(&user, &sqlcUser)

	return &user, nil
}

func (repo *projectUserRepo) ListExternal(ctx context.Context, projectID, ownerID uuid.UUID) ([]project_users.ProjectUser, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectUserRepo.ListExternal",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcUsers, err := repo.queries(ctx).ListProjectUsersExternal(ctx, sqlc.ListProjectUsersExternalParams{
		ProjectID: projectID,
		OwnerID:   ownerID,
	})
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Int("project_users.count", len(sqlcUsers)))

	users := make([]project_users.ProjectUser, 0, len(sqlcUsers))
	for _, u := range sqlcUsers {
		var user project_users.ProjectUser
		mapProjectUserFromDB(&user, &u)
		users = append(users, user)
	}

	return users, nil
}

func (repo *projectUserRepo) ListInternal(ctx context.Context, projectID uuid.UUID) ([]project_users.ProjectUser, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectUserRepo.ListInternal",
		trace.WithAttributes(
			attribute.String("project.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcUsers, err := repo.queries(ctx).ListProjectUsersInternal(ctx, projectID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Int("project_users.count", len(sqlcUsers)))

	users := make([]project_users.ProjectUser, 0, len(sqlcUsers))
	for _, u := range sqlcUsers {
		var user project_users.ProjectUser
		mapProjectUserFromDB(&user, &u)
		users = append(users, user)
	}

	return users, nil
}

func (repo *projectUserRepo) Update(ctx context.Context, toUpdate project_users.ProjectUser, ownerID uuid.UUID) (*project_users.ProjectUser, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectUserRepo.Update",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project.project_id", toUpdate.ProjectID.String()),
			attribute.String("project_user.id", toUpdate.ID.String()),
		),
	)
	defer span.End()

	sqlcUser, err := repo.queries(ctx).UpdateProjectUser(ctx, sqlc.UpdateProjectUserParams{
		ID:           toUpdate.ID,
		ProjectID:    toUpdate.ProjectID,
		Email:        toUpdate.Email,
		PasswordHash: toUpdate.PasswordHash,
		OwnerID:      ownerID,
	})
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("project_user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	mapProjectUserFromDB(&toUpdate, &sqlcUser)

	return &toUpdate, nil
}

func (repo *projectUserRepo) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ProjectUserRepo.UpdateLastLogin",
		trace.WithAttributes(
			attribute.String("project_user.id", id.String()),
		),
	)
	defer span.End()

	if err := repo.queries(ctx).UpdateLastLoginProjectUser(ctx, id); err != nil {
		return fail.From(err).WithArgs("project user").RecordCtx(ctx)
	}
	return nil
}

func (repo *projectUserRepo) Delete(ctx context.Context, projectUserID, projectID, ownerID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ProjectUserRepo.Delete",
		trace.WithAttributes(
			attribute.String("project.project_id", projectID.String()),
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project_user.id", projectUserID.String()),
		),
	)
	defer span.End()

	if err := repo.queries(ctx).DeleteProjectUser(ctx, sqlc.DeleteProjectUserParams{
		ID:        projectUserID,
		ProjectID: projectID,
		OwnerID:   ownerID,
	}); err != nil {
		return fail.From(err).WithArgs("project user").RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Bool("project_user.deleted", true))

	return nil
}

func (repo *projectUserRepo) UpdateMetadata(ctx context.Context, userID, projectID uuid.UUID, metadata *json.RawMessage) error {
	ctx, span := repo.tracer.Start(ctx, "ProjectUserRepo.UpdateMetadata",
		trace.WithAttributes(
			attribute.String("project.project_id", projectID.String()),
			attribute.String("project_user.id", userID.String()),
		),
	)
	defer span.End()

	if err := repo.queries(ctx).UpdateProjectUserMetadata(ctx, sqlc.UpdateProjectUserMetadataParams{
		ID:        userID,
		ProjectID: projectID,
		Metadata:  *metadata,
	}); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

func (repo *projectUserRepo) UpdateSubContext(ctx context.Context, userID, projectID uuid.UUID, subContext json.RawMessage) error {
	ctx, span := repo.tracer.Start(ctx, "ProjectUserRepo.UpdateSubContext",
		trace.WithAttributes(
			attribute.String("project.project_id", projectID.String()),
			attribute.String("project_user.id", userID.String()),
		),
	)
	defer span.End()

	// Get current user metadata
	user, err := repo.GetByIDInternal(ctx, userID, projectID)
	if err != nil {
		return err
	}

	// Parse current metadata
	var metadataMap map[string]any
	if user.Metadata != nil && len(*user.Metadata) > 0 {
		if err := json.Unmarshal(*user.Metadata, &metadataMap); err != nil {
			return fail.From(err).RecordCtx(ctx)
		}
	} else {
		metadataMap = make(map[string]any)
	}

	// Parse new sub-context
	var subContextMap map[string]any
	if err := json.Unmarshal(subContext, &subContextMap); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	// Merge sub-context into metadata
	metadataMap["sub-context"] = subContextMap

	// Marshal updated metadata
	updatedMetadata, err := json.Marshal(metadataMap)
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	rawMetadata := json.RawMessage(updatedMetadata)
	if err := repo.queries(ctx).UpdateProjectUserMetadata(ctx, sqlc.UpdateProjectUserMetadataParams{
		ID:        userID,
		ProjectID: projectID,
		Metadata:  rawMetadata,
	}); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Bool("sub_context.updated", true))
	return nil
}

func (repo *projectUserRepo) Verify(ctx context.Context, userID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectUserRepo.Verify",
		trace.WithAttributes(
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	wasVerified, err := repo.queries(ctx).VerifyProjectUser(ctx, userID)
	if err != nil {
		return false, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Bool("user.was_already_verified", !wasVerified))

	return !wasVerified, nil
}

func (repo *projectUserRepo) BelongsToProject(ctx context.Context, userID, projectID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectUserRepo.BelongsToProject",
		trace.WithAttributes(
			attribute.String("user.id", userID.String()),
			attribute.String("user.project_id", projectID.String()),
		),
	)
	defer span.End()

	belongs, err := repo.queries(ctx).ProjectUserBelongsToProject(ctx, sqlc.ProjectUserBelongsToProjectParams{
		ID:        userID,
		ProjectID: projectID,
	})
	if err != nil {
		return false, fail.From(err).RecordCtx(ctx)
	}

	return belongs, nil
}

func (repo *projectUserRepo) ResetPassword(ctx context.Context, userID uuid.UUID, passwordHash []byte) error {
	ctx, span := repo.tracer.Start(ctx, "ProjectUserRepo.ResetPassword",
		trace.WithAttributes(
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	err := repo.queries(ctx).ResetProjectUserPassword(ctx, sqlc.ResetProjectUserPasswordParams{
		PasswordHash: string(passwordHash),
		ID:           userID,
	})
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}
