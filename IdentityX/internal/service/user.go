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

func (s *AuthService) Login(ctx context.Context, req models.LoginUserRequest) *resp.Response {
  req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	dbUser, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		readable := utils.ParseDBError(err)
		if strings.Contains(readable.Error(), "record not found") {
      return resp.Unauthorized("invalid email or password")
		}
		return resp.InternalServerError("error retrieving user").WithTracePrefix("database-error").AddTrace(readable)
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(req.Password))
	if err != nil {
		return resp.Unauthorized("invalid email or password")
	}

	return nil
}
