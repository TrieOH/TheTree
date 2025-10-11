package service

import (
	"GoAuth/internal/repository"
)

type GreetService struct {
	queries *repository.Queries
}

func NewGreetService(queries *repository.Queries) *GreetService {
	return &GreetService{queries: queries}
}

