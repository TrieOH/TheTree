package service

import (
	"GoAuth/internal/repository"
)

type AuthService struct {
	queries *repository.Queries
}

func NewAuthService(queries *repository.Queries) *AuthService {
	return &AuthService{queries: queries}
}
