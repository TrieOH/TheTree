// Package forms is the forms feature: its repo, handlers and operations.
// Build wires it.
package forms

import (
	"Informd/errs"
	"Informd/internal/authz"
	"Informd/models"
	"Informd/ports"
	"context"
	"errors"
	"lib/database"
	idx "sdk/identityx"
	"time"
	"uuid"
)

type Deps struct {
	Forms      ports.FormsRepo
	Steps      ports.StepRepo
	Namespaces ports.NamespaceRepo
	Fields     ports.FieldsRepo
	Answers    ports.AnswerRepo
	Responses  ports.ResponseRepo
	Responders ports.ResponderRepo
	Authz      *authz.Service
	// Tx opens the transactions multi-step writes run in; threaded from
	// boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) AddMember(ctx context.Context, payload models.AddFormMemberInput) (err error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return errs.ErrAddSelfToForm
	}

	var form *models.Form
	form, err = o.Forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	if payload.UserID == form.OwnerID {
		return errs.ErrFormOwnerAlreadyMember
	}
	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Forms.GetMember(ctx, payload.UserID, form.ID)
	if !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err == nil {
		return errs.ErrAlreadyFormMember
	}

	newMember := models.FormMember{
		UserID:  payload.UserID,
		FormID:  payload.FormID,
		Role:    payload.Role,
		AddedAt: time.Now(),
		AddedBy: ident.Sub.ID,
	}

	return o.Forms.AddMember(ctx, newMember)
}

func (o *Operations) Archive(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusClosed {
		return nil, errs.ErrArchiveNotClosed
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Forms.Archive(ctx, formID)
}

func (o *Operations) Close(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusOpen {
		return nil, errs.ErrCloseNotOpen
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Forms.Close(ctx, formID)
}

func (o *Operations) Create(ctx context.Context, title string) (*models.Form, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	form, err := models.NewForm(nil, ident.Sub.ID, ident.Sub.ID, title)
	if err != nil {
		return nil, err
	}

	var created *models.Form
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = o.Forms.Create(ctx, *form)
		if err != nil {
			return err
		}

		return o.Forms.AddMember(ctx, models.FormMember{
			UserID:  ident.Sub.ID,
			FormID:  created.ID,
			Role:    models.FormMemberRoleOwner,
			AddedAt: time.Now(),
			AddedBy: ident.Sub.ID,
		})
	})
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (o *Operations) GetAnswerable(ctx context.Context, formID uuid.UUID) (*models.FormAnswerable, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusOpen {
		return nil, errs.ErrFormNotAnswerable
	}

	steps, err := o.Steps.List(ctx, formID)
	if err != nil {
		return nil, err
	}
	fields, err := o.Fields.ListByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}

	// index fields by step ID
	fieldsByStep := make(map[uuid.UUID][]models.FieldAnswerable)
	for _, f := range fields {
		var selectConfig *models.FieldSelectConfig
		if f.Type == models.FieldTypeSelect {
			selectConfig, err = o.Fields.GetSelectConfig(ctx, f.ID)
			if err != nil {
				return nil, err
			}
		}

		fieldsByStep[f.StepID] = append(fieldsByStep[f.StepID], models.FieldAnswerable{
			Field:             f,
			FieldSelectConfig: selectConfig,
		})
	}

	// assemble steps
	fullSteps := make([]models.StepAnswerable, len(steps))
	for i, s := range steps {
		fullSteps[i] = models.StepAnswerable{
			Step:   s,
			Fields: fieldsByStep[s.ID],
		}
	}

	return &models.FormAnswerable{
		Form:  *form,
		Steps: fullSteps,
	}, nil
}

func (o *Operations) GetFull(ctx context.Context, formID uuid.UUID) (*models.FullForm, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	steps, err := o.Steps.List(ctx, formID)
	if err != nil {
		return nil, err
	}
	fields, err := o.Fields.ListByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}
	answers, err := o.Answers.GetByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}
	responses, err := o.Responses.ListByForm(ctx, formID)
	if err != nil {
		return nil, err
	}
	responders, err := o.Responders.GetByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}

	// index responders by ID
	responderByID := make(map[uuid.UUID]models.Responder, len(responders))
	for _, r := range responders {
		responderByID[r.ID] = r
	}

	// index responses by ID
	responseByID := make(map[uuid.UUID]models.Response, len(responses))
	for _, r := range responses {
		responseByID[r.ID] = r
	}

	// index answers by field ID
	answersByField := make(map[uuid.UUID][]models.FullAnswer)
	for _, a := range answers {
		email := "anonymous"
		if response, ok := responseByID[a.ResponseID]; ok {
			if response.ResponderID != nil {
				if r, ok := responderByID[*response.ResponderID]; ok {
					email = r.Email
				}
			} else if response.Email != nil {
				email = *response.Email
			}
		}
		if a.FieldID != nil {
			answersByField[*a.FieldID] = append(answersByField[*a.FieldID], models.FullAnswer{
				Answer:    a,
				Responder: email,
			})
		}
	}

	// index fields by step ID
	fieldsByStep := make(map[uuid.UUID][]models.FullField)
	for _, f := range fields {
		fieldsByStep[f.StepID] = append(fieldsByStep[f.StepID], models.FullField{
			Field:   f,
			Answers: answersByField[f.ID],
		})
	}

	// assemble steps
	fullSteps := make([]models.FullStep, len(steps))
	for i, s := range steps {
		fullSteps[i] = models.FullStep{
			Step:   s,
			Fields: fieldsByStep[s.ID],
		}
	}

	return &models.FullForm{
		Form:  *form,
		Steps: fullSteps,
	}, nil
}

func (o *Operations) GetResponseCount(ctx context.Context, formID uuid.UUID) (int, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return 0, err
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleMember)
	if err != nil {
		return 0, err
	}

	return o.Forms.ResponsesCount(ctx, formID)
}

func (o *Operations) ListMembers(ctx context.Context, formID uuid.UUID) ([]models.FormMember, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	members, err := o.Forms.ListDirectMembers(ctx, form.ID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (o *Operations) ListForms(ctx context.Context) (forms []models.Form, err error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	forms, err = o.Forms.ListMine(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}

func (o *Operations) ListArchivedForms(ctx context.Context) (forms []models.Form, err error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	forms, err = o.Forms.ListMineArchived(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}

func (o *Operations) Open(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusDraft {
		return nil, errs.ErrOpenNotDraft
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Forms.Open(ctx, formID)
}

func (o *Operations) ReDraft(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusOpen {
		return nil, errs.ErrRedraftNotOpen
	}

	count, err := o.Forms.ResponsesCount(ctx, formID)
	if err != nil {
		return nil, err
	}

	if count != 0 {
		return nil, errs.ErrRedraftWithResponses
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Forms.ReDraft(ctx, formID)
}

func (o *Operations) RemoveMember(ctx context.Context, payload models.RemoveFormMemberInput) (err error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return errs.ErrRemoveSelfFromForm
	}

	var form *models.Form
	form, err = o.Forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	if payload.UserID == form.OwnerID {
		return errs.ErrRemoveFormOwner
	}
	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Forms.GetMember(ctx, payload.UserID, form.ID)
	if !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err != nil {
		return errs.ErrNotFormMember
	}

	return o.Forms.RemoveMember(ctx, payload.UserID, payload.FormID)
}
