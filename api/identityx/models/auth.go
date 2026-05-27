package models

type CredentialType string

const (
	TokenCredentialType  CredentialType = "token"
	ApiKeyCredentialType CredentialType = "api_key"
)

type IDXRegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,passwd,min=8,max=72"`
}

type IDXRegisterInput struct {
	Email    string
	Password string
}

func (r IDXRegisterRequest) ToInput() IDXRegisterInput {
	return IDXRegisterInput{
		Email:    r.Email,
		Password: r.Password,
	}
}

type IDXLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwd,min=8"`
}

type IDXLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r IDXLoginRequest) ToInput() IDXLoginInput {
	return IDXLoginInput{
		Email:    r.Email,
		Password: r.Password,
	}
}

type SetupInput struct {
	Email    string
	Password string
}

func (r IDXLoginRequest) ToSetupInput() SetupInput {
	return SetupInput{
		Email:    r.Email,
		Password: r.Password,
	}
}

type LogoutInput struct {
	AccessToken  string
	RefreshToken string
}
