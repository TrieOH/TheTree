// Package forms implements the StrictServerInterface methods for the forms
// feature. NOTE: the /forms/{form_id}/asnwerable path keeps its original
// spelling — part of the interface.
package forms

import (
	"context"
	"time"

	"Informd/internal/openapi"
	"Informd/internal/services"
	"Informd/models"
)

const module = "Informd"

type Handlers struct {
	ops *services.Forms
}

func New(ops *services.Forms) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListMyForms(ctx context.Context, _ openapi.ListMyFormsRequestObject) (openapi.ListMyFormsResponseObject, error) {
	forms, err := h.ops.ListForms(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListMyForms200JSONResponse{
		Code: 200, Data: &forms, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CreateForm(ctx context.Context, req openapi.CreateFormRequestObject) (openapi.CreateFormResponseObject, error) {
	form, err := h.ops.Create(ctx, req.Body.Title)
	if err != nil {
		return nil, err
	}
	return openapi.CreateForm201JSONResponse{
		Code: 201, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListMyArchivedForms(ctx context.Context, _ openapi.ListMyArchivedFormsRequestObject) (openapi.ListMyArchivedFormsResponseObject, error) {
	forms, err := h.ops.ListArchivedForms(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListMyArchivedForms200JSONResponse{
		Code: 200, Data: &forms, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetAnswerableForm(ctx context.Context, req openapi.GetAnswerableFormRequestObject) (openapi.GetAnswerableFormResponseObject, error) {
	form, err := h.ops.GetAnswerable(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetAnswerableForm200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetFullForm(ctx context.Context, req openapi.GetFullFormRequestObject) (openapi.GetFullFormResponseObject, error) {
	form, err := h.ops.GetFull(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetFullForm200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListFormMembers(ctx context.Context, req openapi.ListFormMembersRequestObject) (openapi.ListFormMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ListFormMembers200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) AddFormMember(ctx context.Context, req openapi.AddFormMemberRequestObject) (openapi.AddFormMemberResponseObject, error) {
	err := h.ops.AddMember(ctx, models.AddFormMemberInput{
		UserID: req.Body.UserId,
		FormID: req.FormId,
		Role:   req.Body.Role,
	})
	if err != nil {
		return nil, err
	}
	return openapi.AddFormMember201Response{}, nil
}

func (h *Handlers) RemoveFormMember(ctx context.Context, req openapi.RemoveFormMemberRequestObject) (openapi.RemoveFormMemberResponseObject, error) {
	err := h.ops.RemoveMember(ctx, models.RemoveFormMemberInput{
		UserID: req.Body.UserId,
		FormID: req.FormId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveFormMember200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) OpenForm(ctx context.Context, req openapi.OpenFormRequestObject) (openapi.OpenFormResponseObject, error) {
	form, err := h.ops.Open(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.OpenForm200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CloseForm(ctx context.Context, req openapi.CloseFormRequestObject) (openapi.CloseFormResponseObject, error) {
	form, err := h.ops.Close(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.CloseForm200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ArchiveForm(ctx context.Context, req openapi.ArchiveFormRequestObject) (openapi.ArchiveFormResponseObject, error) {
	form, err := h.ops.Archive(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ArchiveForm200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) RedraftForm(ctx context.Context, req openapi.RedraftFormRequestObject) (openapi.RedraftFormResponseObject, error) {
	form, err := h.ops.ReDraft(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.RedraftForm200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetFormResponseCount(ctx context.Context, req openapi.GetFormResponseCountRequestObject) (openapi.GetFormResponseCountResponseObject, error) {
	count, err := h.ops.GetResponseCount(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetFormResponseCount200JSONResponse{
		Code: 200,
		Data: &struct {
			Count int `json:"count"`
		}{Count: count},
		Timestamp: time.Now(), Module: module,
	}, nil
}
