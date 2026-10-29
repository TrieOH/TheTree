package forms

import (
	"Informd/errs"
	"Informd/internal/sqlc"
	"Informd/models"
	"Informd/ports"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.FormsRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrFormNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapForm(src sqlc.Form) models.Form {
	return models.Form{
		ID:          src.ID,
		NamespaceID: src.NamespaceID,
		OwnerID:     src.OwnerID,
		Title:       src.Name,
		Status:      models.FormStatus(src.Status),
		OpenedAt:    src.OpenedAt,
		ClosedAt:    src.ClosedAt,
		ArchivedAt:  src.ArchivedAt,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
	}
}

func mapFormMember(src sqlc.FormMember) models.FormMember {
	return models.FormMember{
		UserID:  src.UserID,
		FormID:  src.FormID,
		Role:    models.FormMemberRole(src.Role),
		AddedAt: src.AddedAt,
		AddedBy: src.AddedBy,
	}
}

func (repo *Repo) AddMember(ctx context.Context, toCreate models.FormMember) error {
	err := repo.db(ctx).AddFormMember(ctx, sqlc.AddFormMemberParams{
		UserID:  toCreate.UserID,
		FormID:  toCreate.FormID,
		Role:    string(toCreate.Role),
		AddedAt: toCreate.AddedAt,
		AddedBy: toCreate.AddedBy,
	})
	return repo.dbe(err)
}

func (repo *Repo) Archive(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	sqlcForm, err := repo.db(ctx).ArchiveForm(ctx, formID)
	return database.One(sqlcForm, err, repo.dbe, mapForm)
}

func (repo *Repo) Close(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	sqlcForm, err := repo.db(ctx).CloseForm(ctx, formID)
	return database.One(sqlcForm, err, repo.dbe, mapForm)
}

func (repo *Repo) Create(ctx context.Context, toCreate models.Form) (*models.Form, error) {
	sqlcForm, err := repo.db(ctx).CreateForm(ctx, sqlc.CreateFormParams{
		NamespaceID: toCreate.NamespaceID,
		CreatedBy:   toCreate.CreatedBy,
		OwnerID:     toCreate.OwnerID,
		Name:        toCreate.Title,
		Status:      string(toCreate.Status),
	})
	return database.One(sqlcForm, err, repo.dbe, mapForm)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Form, error) {
	sqlcForm, err := repo.db(ctx).GetFormByID(ctx, id)
	return database.One(sqlcForm, err, repo.dbe, mapForm)
}

func (repo *Repo) GetMember(ctx context.Context, userID, formID uuid.UUID) (*models.FormMember, error) {
	sqlcMember, err := repo.db(ctx).GetFormMember(ctx, sqlc.GetFormMemberParams{
		UserID: userID,
		FormID: formID,
	})
	return database.One(sqlcMember, err, repo.dbe, mapFormMember)
}

func (repo *Repo) ListDirectMembers(ctx context.Context, formID uuid.UUID) ([]models.FormMember, error) {
	sqlcMembers, err := repo.db(ctx).ListDirectFormMembers(ctx, formID)
	return database.Many(sqlcMembers, err, repo.dbe, mapFormMember)
}

func (repo *Repo) ListMineArchived(ctx context.Context, userID uuid.UUID) ([]models.Form, error) {
	sqlcForms, err := repo.db(ctx).ListMyArchivedForms(ctx, userID)
	return database.Many(sqlcForms, err, repo.dbe, mapForm)
}

func (repo *Repo) ListMine(ctx context.Context, userID uuid.UUID) ([]models.Form, error) {
	sqlcForms, err := repo.db(ctx).ListMyForms(ctx, userID)
	return database.Many(sqlcForms, err, repo.dbe, mapForm)
}

func (repo *Repo) ListFromNamespaceArchived(ctx context.Context, namespaceID uuid.UUID) ([]models.Form, error) {
	sqlcForms, err := repo.db(ctx).ListNamespaceArchivedForms(ctx, &namespaceID)
	return database.Many(sqlcForms, err, repo.dbe, mapForm)
}

func (repo *Repo) ListFromNamespace(ctx context.Context, namespaceID uuid.UUID) ([]models.Form, error) {
	sqlcForms, err := repo.db(ctx).ListNamespaceForms(ctx, &namespaceID)
	return database.Many(sqlcForms, err, repo.dbe, mapForm)
}

func (repo *Repo) Open(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	sqlcForm, err := repo.db(ctx).OpenForm(ctx, formID)
	return database.One(sqlcForm, err, repo.dbe, mapForm)
}

func (repo *Repo) ReDraft(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	sqlcForm, err := repo.db(ctx).DraftForm(ctx, formID)
	return database.One(sqlcForm, err, repo.dbe, mapForm)
}

func (repo *Repo) RemoveMember(ctx context.Context, userID, formID uuid.UUID) error {
	err := repo.db(ctx).RemoveFormMember(ctx, sqlc.RemoveFormMemberParams{
		UserID: userID,
		FormID: formID,
	})
	return repo.dbe(err)
}

func (repo *Repo) ResponsesCount(ctx context.Context, formID uuid.UUID) (int, error) {
	count, err := repo.db(ctx).CountFormResponses(ctx, formID)
	if err != nil {
		return 0, repo.dbe(err)
	}
	return int(count), nil
}

func (repo *Repo) GetRole(ctx context.Context, actorID, formID uuid.UUID) (models.FormMemberRole, error) {
	member, err := repo.GetMember(ctx, actorID, formID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}
