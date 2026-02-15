package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/ports/inbounds"
	"encoding/json"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type SubContextHandler struct {
	subContext inbounds.SubContextService
}

func NewSubContextHandler(uc inbounds.SubContextService) *SubContextHandler {
	return &SubContextHandler{subContext: uc}
}

// AddSubContext godoc
// @Summary Add sub-context data to a project user
// @Description Adds or updates sub-context metadata for a specific project user. Sub-context is admin-controlled metadata.
// @Tags sub-context
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project"
// @Param data body dto.AddSubContextRequest true "Sub-context data to add"
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Success 200 {object} object "Sub-context added successfully"
// @Failure 400 {object} ErrorResponse "Bad Request: Invalid input"
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 403 {object} ErrorResponse "Forbidden: Only clients can manage sub-context"
// @Failure 404 {object} ErrorResponse "Not Found: Project or user not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/sub-context [post]
func (handler *SubContextHandler) AddSubContext(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.AddSubContextRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	if err := handler.subContext.Add(r.Context(), projectID, req.UserID, req.Data); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Sub-context added successfully").Send(w)
}

// RemoveSubContext godoc
// @Summary Remove sub-context keys from a project user
// @Description Removes specified keys from the sub-context metadata of a project user. Supports nested keys using dot notation.
// @Tags sub-context
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project"
// @Param data body dto.RemoveSubContextRequest true "Keys to remove from sub-context"
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Success 200 {object} object "Sub-context keys removed successfully"
// @Failure 400 {object} ErrorResponse "Bad Request: Invalid input"
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 403 {object} ErrorResponse "Forbidden: Only clients can manage sub-context"
// @Failure 404 {object} ErrorResponse "Not Found: Project or user not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/sub-context [delete]
func (handler *SubContextHandler) RemoveSubContext(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.RemoveSubContextRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	if err := handler.subContext.Remove(r.Context(), projectID, req.UserID, req.Keys); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Sub-context keys removed successfully").Send(w)
}

// GetSubContext godoc
// @Summary Get sub-context for a project user
// @Description Retrieves the sub-context metadata for a specific project user.
// @Tags sub-context
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project"
// @Param user_id path string true "ID of the user"
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Success 200 {object} object "Sub-context retrieved successfully"
// @Failure 400 {object} ErrorResponse "Bad Request: Invalid project or user ID"
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 403 {object} ErrorResponse "Forbidden: Only clients can view sub-context"
// @Failure 404 {object} ErrorResponse "Not Found: Project or user not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/users/{user_id}/sub-context [get]
func (handler *SubContextHandler) GetSubContext(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	userID, rs := getUUID(r, "user_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	subContext, err := handler.subContext.Get(r.Context(), projectID, userID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	var result map[string]any
	if err := json.Unmarshal(subContext, &result); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(map[string]any{
		"sub_context": result,
	}).Send(w)
}
