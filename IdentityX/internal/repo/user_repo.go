package repo

import (
	"GoAuth/internal/models"
	"GoAuth/internal/sqlc"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type UserRepo interface {
	Register(ctx context.Context, email, password string) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ListUsers() ([]*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(userID string) error
}

type userRepo struct {
	q   *sqlc.Queries
	log *zap.Logger
}

func NewUserRepo(q *sqlc.Queries, l *zap.Logger) UserRepo {
	return &userRepo{
		q:   q,
		log: l,
	}
}

func mapUserFromDB(dst *models.User, src *sqlc.User) error {
	return copier.Copy(dst, src)
}

func (u userRepo) Register(ctx context.Context, email, password string) (*models.User, error) {
	sqlcUser, err := u.q.RegisterUser(ctx, sqlc.RegisterUserParams{
		Email:    email,
		Password: password,
	})

	if err != nil {
		return nil, err
	}

	var user models.User
	if err = mapUserFromDB(&user, &sqlcUser); err != nil {
		u.log.Error("failed to copy user", zap.Error(err))
		return nil, fmt.Errorf("failed to copy user properly: %w", err)
	}

	return &user, nil
}

func (u userRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	sqlcUser, err := u.q.GetUserById(ctx, userID)

	if err != nil {
		return nil, err
	}

	var user models.User
	if err = mapUserFromDB(&user, &sqlcUser); err != nil {
		u.log.Error("failed to copy user", zap.Error(err))
		return nil, fmt.Errorf("failed to copy user properly: %w", err)
	}

	return &user, nil
}

func (u userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	sqlcUser, err := u.q.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	var user models.User
	if err = mapUserFromDB(&user, &sqlcUser); err != nil {
		u.log.Error("failed to copy user", zap.Error(err))
		return nil, fmt.Errorf("failed to copy user properly: %w", err)
	}

	return &user, nil
}

func (u userRepo) ListUsers() ([]*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) UpdateUser(user *models.User) error {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) DeleteUser(userID string) error {
	//TODO implement me
	panic("implement me")
}
