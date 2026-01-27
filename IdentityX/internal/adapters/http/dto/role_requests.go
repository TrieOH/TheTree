package dto

type CreateRoleRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
}

type UpdateRoleRequest struct {
	Description *string `json:"description"`
}
