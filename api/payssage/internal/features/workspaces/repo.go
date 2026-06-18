package workspaces

import (
	"context"
	"payssage/ports"

	"lib/database"
	"payssage/internal/database/sqlc"
	"payssage/internal/shared/errx"
	"payssage/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type workspaceRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.WorkspaceRepo = (*workspaceRepo)(nil)

func NewWorkspaceRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.WorkspaceRepo {
	return &workspaceRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *workspaceRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapWorkspaceFromDB(src *sqlc.Workspace) *models.Workspace {
	return &models.Workspace{
		ID:        src.ID,
		ScopeID:   src.ScopeID,
		UserID:    src.UserID,
		Name:      src.Name,
		Sandbox:   src.Sandbox,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}
}

func (repo *workspaceRepo) Create(ctx context.Context, toCreate models.Workspace) (*models.Workspace, error) {
	ctx, span := repo.tracer.Start(ctx, "WorkspaceRepo.Create")
	defer span.End()

	sqlcWorkspace, err := repo.queries(ctx).CreateWorkspace(ctx, sqlc.CreateWorkspaceParams{
		ScopeID: toCreate.ScopeID,
		ID:      toCreate.ID,
		UserID:  toCreate.UserID,
		Name:    toCreate.Name,
	})
	if err != nil {
		return nil, errx.FromDB(err, "workspace")
	}

	return mapWorkspaceFromDB(&sqlcWorkspace), nil
}

func (repo *workspaceRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Workspace, error) {
	ctx, span := repo.tracer.Start(ctx, "WorkspaceRepo.GetByID")
	defer span.End()

	sqlcWorkspace, err := repo.queries(ctx).GetWorkspaceByID(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "workspace")
	}

	return mapWorkspaceFromDB(&sqlcWorkspace), nil
}
func (repo *workspaceRepo) GetByName(ctx context.Context, name string, userID uuid.UUID) (*models.Workspace, error) {
	ctx, span := repo.tracer.Start(ctx, "WorkspaceRepo.GetByName")
	defer span.End()

	sqlcWorkspace, err := repo.queries(ctx).GetWorkspaceByName(ctx, sqlc.GetWorkspaceByNameParams{
		UserID: userID,
		Name:   name,
	})
	if err != nil {
		return nil, errx.FromDB(err, "workspace")
	}

	return mapWorkspaceFromDB(&sqlcWorkspace), nil
}

func (repo *workspaceRepo) List(ctx context.Context, userID uuid.UUID) ([]models.Workspace, error) {
	ctx, span := repo.tracer.Start(ctx, "WorkspaceRepo.List")
	defer span.End()

	sqlcWorkspaces, err := repo.queries(ctx).ListWorkspacesByUser(ctx, userID)
	if err != nil {
		return nil, errx.FromDB(err, "workspace")
	}

	out := make([]models.Workspace, 0, len(sqlcWorkspaces))
	for _, workspace := range sqlcWorkspaces {
		out = append(out, *mapWorkspaceFromDB(&workspace))
	}
	return out, nil
}

func (repo *workspaceRepo) EnableSandbox(ctx context.Context, id uuid.UUID) (*models.Workspace, error) {
	ctx, span := repo.tracer.Start(ctx, "WorkspaceRepo.EnableSandbox")
	defer span.End()

	sqlcWorkspace, err := repo.queries(ctx).EnableSandbox(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "workspace")
	}

	return mapWorkspaceFromDB(&sqlcWorkspace), nil
}

func (repo *workspaceRepo) DisableSandbox(ctx context.Context, id uuid.UUID) (*models.Workspace, error) {
	ctx, span := repo.tracer.Start(ctx, "WorkspaceRepo.DisableSandbox")
	defer span.End()

	sqlcWorkspace, err := repo.queries(ctx).DisableSandbox(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "workspace")
	}

	return mapWorkspaceFromDB(&sqlcWorkspace), nil
}
