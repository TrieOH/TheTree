package service

import (
  "context"
	"strings"

	"GoAuth/internal/models"
	"GoAuth/internal/repository"
	"GoAuth/internal/utils"
	resp "github.com/MintzyG/GoResponse/response"
	// "github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Register(ctx context.Context, req models.RegisterUserRequest) *resp.Response {
  req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return resp.InternalServerError("error hashing user password").WithTracePrefix("error").AddTrace(err)
	}

	_, err = s.queries.RegisterUser(ctx, repository.RegisterUserParams{
    Email: req.Email,
    Password: string(hashedPassword),
	})

	if err != nil {
		readable := utils.ParseDBError(err)
		return resp.InternalServerError("error registering user").WithTracePrefix("database-error").AddTrace(readable)
	}

	return nil
}
