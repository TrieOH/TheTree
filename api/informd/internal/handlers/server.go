// Package handlers implements the generated StrictServerInterface by
// delegating to one subpackage per feature. Auth, validation, and error
// mapping run in the strict middleware stack (see internal/app); the
// handlers here are pure domain logic + fun envelope construction.
package handlers

import (
	"context"

	"Informd/internal/handlers/fields"
	"Informd/internal/handlers/forms"
	"Informd/internal/handlers/namespaces"
	"Informd/internal/handlers/responses"
	"Informd/internal/handlers/steps"
	"Informd/internal/openapi"
	"Informd/internal/services"
)

// Server implements openapi.StrictServerInterface.
type Server struct {
	namespaces *namespaces.Handlers
	forms      *forms.Handlers
	steps      *steps.Handlers
	fields     *fields.Handlers
	responses  *responses.Handlers
}

// NewServer wires the per-feature handlers from the services aggregate.
func NewServer(ops *services.Operations) *Server {
	return &Server{
		namespaces: namespaces.New(ops.Namespaces),
		forms:      forms.New(ops.Forms),
		steps:      steps.New(ops.Steps),
		fields:     fields.New(ops.Fields),
		responses:  responses.New(ops.Responses),
	}
}

// ── StrictServerInterface ────────────────────────────────────────────────

func (s *Server) ListMyForms(ctx context.Context, req openapi.ListMyFormsRequestObject) (openapi.ListMyFormsResponseObject, error) {
	return s.forms.ListMyForms(ctx, req)
}

func (s *Server) CreateForm(ctx context.Context, req openapi.CreateFormRequestObject) (openapi.CreateFormResponseObject, error) {
	return s.forms.CreateForm(ctx, req)
}

func (s *Server) ListMyArchivedForms(ctx context.Context, req openapi.ListMyArchivedFormsRequestObject) (openapi.ListMyArchivedFormsResponseObject, error) {
	return s.forms.ListMyArchivedForms(ctx, req)
}

func (s *Server) ArchiveForm(ctx context.Context, req openapi.ArchiveFormRequestObject) (openapi.ArchiveFormResponseObject, error) {
	return s.forms.ArchiveForm(ctx, req)
}

func (s *Server) GetAnswerableForm(ctx context.Context, req openapi.GetAnswerableFormRequestObject) (openapi.GetAnswerableFormResponseObject, error) {
	return s.forms.GetAnswerableForm(ctx, req)
}

func (s *Server) CloseForm(ctx context.Context, req openapi.CloseFormRequestObject) (openapi.CloseFormResponseObject, error) {
	return s.forms.CloseForm(ctx, req)
}

func (s *Server) GetFullForm(ctx context.Context, req openapi.GetFullFormRequestObject) (openapi.GetFullFormResponseObject, error) {
	return s.forms.GetFullForm(ctx, req)
}

func (s *Server) RemoveFormMember(ctx context.Context, req openapi.RemoveFormMemberRequestObject) (openapi.RemoveFormMemberResponseObject, error) {
	return s.forms.RemoveFormMember(ctx, req)
}

func (s *Server) ListFormMembers(ctx context.Context, req openapi.ListFormMembersRequestObject) (openapi.ListFormMembersResponseObject, error) {
	return s.forms.ListFormMembers(ctx, req)
}

func (s *Server) AddFormMember(ctx context.Context, req openapi.AddFormMemberRequestObject) (openapi.AddFormMemberResponseObject, error) {
	return s.forms.AddFormMember(ctx, req)
}

func (s *Server) OpenForm(ctx context.Context, req openapi.OpenFormRequestObject) (openapi.OpenFormResponseObject, error) {
	return s.forms.OpenForm(ctx, req)
}

func (s *Server) RedraftForm(ctx context.Context, req openapi.RedraftFormRequestObject) (openapi.RedraftFormResponseObject, error) {
	return s.forms.RedraftForm(ctx, req)
}

func (s *Server) SubmitResponse(ctx context.Context, req openapi.SubmitResponseRequestObject) (openapi.SubmitResponseResponseObject, error) {
	return s.responses.SubmitResponse(ctx, req)
}

func (s *Server) GetFormResponseCount(ctx context.Context, req openapi.GetFormResponseCountRequestObject) (openapi.GetFormResponseCountResponseObject, error) {
	return s.forms.GetFormResponseCount(ctx, req)
}

func (s *Server) ListSteps(ctx context.Context, req openapi.ListStepsRequestObject) (openapi.ListStepsResponseObject, error) {
	return s.steps.ListSteps(ctx, req)
}

func (s *Server) CreateStep(ctx context.Context, req openapi.CreateStepRequestObject) (openapi.CreateStepResponseObject, error) {
	return s.steps.CreateStep(ctx, req)
}

func (s *Server) BulkEditSteps(ctx context.Context, req openapi.BulkEditStepsRequestObject) (openapi.BulkEditStepsResponseObject, error) {
	return s.steps.BulkEditSteps(ctx, req)
}

func (s *Server) ListFields(ctx context.Context, req openapi.ListFieldsRequestObject) (openapi.ListFieldsResponseObject, error) {
	return s.fields.ListFields(ctx, req)
}

func (s *Server) CreateField(ctx context.Context, req openapi.CreateFieldRequestObject) (openapi.CreateFieldResponseObject, error) {
	return s.fields.CreateField(ctx, req)
}

func (s *Server) BulkEditFields(ctx context.Context, req openapi.BulkEditFieldsRequestObject) (openapi.BulkEditFieldsResponseObject, error) {
	return s.fields.BulkEditFields(ctx, req)
}

func (s *Server) DeleteField(ctx context.Context, req openapi.DeleteFieldRequestObject) (openapi.DeleteFieldResponseObject, error) {
	return s.fields.DeleteField(ctx, req)
}

func (s *Server) GetSelectConfig(ctx context.Context, req openapi.GetSelectConfigRequestObject) (openapi.GetSelectConfigResponseObject, error) {
	return s.fields.GetSelectConfig(ctx, req)
}

func (s *Server) EditSelectConfig(ctx context.Context, req openapi.EditSelectConfigRequestObject) (openapi.EditSelectConfigResponseObject, error) {
	return s.fields.EditSelectConfig(ctx, req)
}

func (s *Server) ListNamespaces(ctx context.Context, req openapi.ListNamespacesRequestObject) (openapi.ListNamespacesResponseObject, error) {
	return s.namespaces.ListNamespaces(ctx, req)
}

func (s *Server) CreateNamespace(ctx context.Context, req openapi.CreateNamespaceRequestObject) (openapi.CreateNamespaceResponseObject, error) {
	return s.namespaces.CreateNamespace(ctx, req)
}

func (s *Server) ListNamespaceForms(ctx context.Context, req openapi.ListNamespaceFormsRequestObject) (openapi.ListNamespaceFormsResponseObject, error) {
	return s.namespaces.ListNamespaceForms(ctx, req)
}

func (s *Server) CreateNamespaceForm(ctx context.Context, req openapi.CreateNamespaceFormRequestObject) (openapi.CreateNamespaceFormResponseObject, error) {
	return s.namespaces.CreateNamespaceForm(ctx, req)
}

func (s *Server) ListNamespaceArchivedForms(ctx context.Context, req openapi.ListNamespaceArchivedFormsRequestObject) (openapi.ListNamespaceArchivedFormsResponseObject, error) {
	return s.namespaces.ListNamespaceArchivedForms(ctx, req)
}

func (s *Server) ArchiveFormNamespaced(ctx context.Context, req openapi.ArchiveFormNamespacedRequestObject) (openapi.ArchiveFormNamespacedResponseObject, error) {
	return s.namespaces.ArchiveFormNamespaced(ctx, req)
}

func (s *Server) CloseFormNamespaced(ctx context.Context, req openapi.CloseFormNamespacedRequestObject) (openapi.CloseFormNamespacedResponseObject, error) {
	return s.namespaces.CloseFormNamespaced(ctx, req)
}

func (s *Server) GetFullFormNamespaced(ctx context.Context, req openapi.GetFullFormNamespacedRequestObject) (openapi.GetFullFormNamespacedResponseObject, error) {
	return s.namespaces.GetFullFormNamespaced(ctx, req)
}

func (s *Server) RemoveFormMemberNamespaced(ctx context.Context, req openapi.RemoveFormMemberNamespacedRequestObject) (openapi.RemoveFormMemberNamespacedResponseObject, error) {
	return s.namespaces.RemoveFormMemberNamespaced(ctx, req)
}

func (s *Server) ListFormMembersNamespaced(ctx context.Context, req openapi.ListFormMembersNamespacedRequestObject) (openapi.ListFormMembersNamespacedResponseObject, error) {
	return s.namespaces.ListFormMembersNamespaced(ctx, req)
}

func (s *Server) AddFormMemberNamespaced(ctx context.Context, req openapi.AddFormMemberNamespacedRequestObject) (openapi.AddFormMemberNamespacedResponseObject, error) {
	return s.namespaces.AddFormMemberNamespaced(ctx, req)
}

func (s *Server) OpenFormNamespaced(ctx context.Context, req openapi.OpenFormNamespacedRequestObject) (openapi.OpenFormNamespacedResponseObject, error) {
	return s.namespaces.OpenFormNamespaced(ctx, req)
}

func (s *Server) RedraftFormNamespaced(ctx context.Context, req openapi.RedraftFormNamespacedRequestObject) (openapi.RedraftFormNamespacedResponseObject, error) {
	return s.namespaces.RedraftFormNamespaced(ctx, req)
}

func (s *Server) GetFormResponseCountNamespaced(ctx context.Context, req openapi.GetFormResponseCountNamespacedRequestObject) (openapi.GetFormResponseCountNamespacedResponseObject, error) {
	return s.namespaces.GetFormResponseCountNamespaced(ctx, req)
}

func (s *Server) ListStepsNamespaced(ctx context.Context, req openapi.ListStepsNamespacedRequestObject) (openapi.ListStepsNamespacedResponseObject, error) {
	return s.steps.ListStepsNamespaced(ctx, req)
}

func (s *Server) CreateStepNamespaced(ctx context.Context, req openapi.CreateStepNamespacedRequestObject) (openapi.CreateStepNamespacedResponseObject, error) {
	return s.steps.CreateStepNamespaced(ctx, req)
}

func (s *Server) BulkEditStepsNamespaced(ctx context.Context, req openapi.BulkEditStepsNamespacedRequestObject) (openapi.BulkEditStepsNamespacedResponseObject, error) {
	return s.steps.BulkEditStepsNamespaced(ctx, req)
}

func (s *Server) ListFieldsNamespaced(ctx context.Context, req openapi.ListFieldsNamespacedRequestObject) (openapi.ListFieldsNamespacedResponseObject, error) {
	return s.fields.ListFieldsNamespaced(ctx, req)
}

func (s *Server) CreateFieldNamespaced(ctx context.Context, req openapi.CreateFieldNamespacedRequestObject) (openapi.CreateFieldNamespacedResponseObject, error) {
	return s.fields.CreateFieldNamespaced(ctx, req)
}

func (s *Server) BulkEditFieldsNamespaced(ctx context.Context, req openapi.BulkEditFieldsNamespacedRequestObject) (openapi.BulkEditFieldsNamespacedResponseObject, error) {
	return s.fields.BulkEditFieldsNamespaced(ctx, req)
}

func (s *Server) DeleteFieldNamespaced(ctx context.Context, req openapi.DeleteFieldNamespacedRequestObject) (openapi.DeleteFieldNamespacedResponseObject, error) {
	return s.fields.DeleteFieldNamespaced(ctx, req)
}

func (s *Server) GetSelectConfigNamespaced(ctx context.Context, req openapi.GetSelectConfigNamespacedRequestObject) (openapi.GetSelectConfigNamespacedResponseObject, error) {
	return s.fields.GetSelectConfigNamespaced(ctx, req)
}

func (s *Server) EditSelectConfigNamespaced(ctx context.Context, req openapi.EditSelectConfigNamespacedRequestObject) (openapi.EditSelectConfigNamespacedResponseObject, error) {
	return s.fields.EditSelectConfigNamespaced(ctx, req)
}

func (s *Server) RemoveNamespaceMember(ctx context.Context, req openapi.RemoveNamespaceMemberRequestObject) (openapi.RemoveNamespaceMemberResponseObject, error) {
	return s.namespaces.RemoveNamespaceMember(ctx, req)
}

func (s *Server) ListNamespaceMembers(ctx context.Context, req openapi.ListNamespaceMembersRequestObject) (openapi.ListNamespaceMembersResponseObject, error) {
	return s.namespaces.ListNamespaceMembers(ctx, req)
}

func (s *Server) AddNamespaceMember(ctx context.Context, req openapi.AddNamespaceMemberRequestObject) (openapi.AddNamespaceMemberResponseObject, error) {
	return s.namespaces.AddNamespaceMember(ctx, req)
}
