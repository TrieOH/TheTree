package service

import (
	"GoAuth/internal/repo"
	"GoAuth/internal/sqlc"
)

type AuthService struct {
	userRepo                 repo.UserRepo
	sessionRepo              repo.SessionRepo
	revokedRefreshTokensRepo repo.RevokedRefreshTokensRepo
	queries                  *sqlc.Queries
}

func NewAuthService(userRepo repo.UserRepo, sessionRepo repo.SessionRepo, rrr repo.RevokedRefreshTokensRepo, queries *sqlc.Queries) *AuthService {
	return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo, revokedRefreshTokensRepo: rrr, queries: queries}
}
