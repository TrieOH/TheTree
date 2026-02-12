package dto

type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Service string `json:"service" example:"univents-api"`
	UserID  string `json:"user_id,omitempty" example:"some-uuid"`
}
