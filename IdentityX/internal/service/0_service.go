package service

import (
	"GoAuth/internal/repo"
)

type AuthService struct {
	userRepo                 repo.UserRepo
	sessionRepo              repo.SessionRepo
	revokedRefreshTokensRepo repo.RevokedRefreshTokensRepo
	projectRepo              repo.ProjectRepo
	projectUserRepo          repo.ProjectUserRepo
}

func NewAuthService(
	userRepo repo.UserRepo,
	sessionRepo repo.SessionRepo,
	revokedRefreshTokensRepo repo.RevokedRefreshTokensRepo,
	projectRepo repo.ProjectRepo,
	projectUserRepo repo.ProjectUserRepo,
) *AuthService {
	return &AuthService{
		userRepo:                 userRepo,
		sessionRepo:              sessionRepo,
		revokedRefreshTokensRepo: revokedRefreshTokensRepo,
		projectRepo:              projectRepo,
		projectUserRepo:          projectUserRepo,
	}
}
