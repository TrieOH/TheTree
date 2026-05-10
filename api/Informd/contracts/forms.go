package contracts

import (
	"fmt"
	"lib/utils"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FormStatus string

const (
	FormStatusDraft    FormStatus = "draft"
	FormStatusOpen     FormStatus = "open"
	FormStatusClosed   FormStatus = "closed"
	FormStatusArchived FormStatus = "archived"
)

type Form struct {
	ID          uuid.UUID  `json:"id"`
	NamespaceID *uuid.UUID `json:"namespace_id"`
	OwnerID     uuid.UUID  `json:"owner_id" validate:"required"`
	Title       string     `json:"title"    validate:"required"`
	Status      FormStatus `json:"status"`
	OpenedAt    *time.Time `json:"opened_at"`
	ClosedAt    *time.Time `json:"closed_at"`
	ArchivedAt  *time.Time `json:"archived_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewForm(namespaceID *uuid.UUID, ownerID uuid.UUID, title string) (*Form, error) {
	f := &Form{
		NamespaceID: namespaceID,
		OwnerID:     ownerID,
		Title:       title,
		Status:      FormStatusDraft,
	}
	if err := validate.Struct(f); err != nil {
		return nil, err
	}
	return f, nil
}

type CreateFormRequest struct {
	Title string `json:"title" validate:"required"`
}

type BulkGetRequest struct {
	IDs []uuid.UUID `json:"ids" validate:"required"`
}

type CreateStepRequest struct {
	Title        string  `json:"title" validate:"required"`
	Description  *string `json:"description"`
	PositionHint int     `json:"position_hint" validate:"required"`
}

func FilterForms(forms []Form, params BulkGetParams) ([]Form, error) {
	if params.FilterKey == "" || params.FilterOp == "" {
		return forms, nil
	}
	key, ok := AllowedFilterKeys[params.FilterKey]
	if !ok {
		return nil, fmt.Errorf("invalid filter_key: %s", params.FilterKey)
	}

	op, ok := AllowedOps[params.FilterOp]
	if !ok {
		return nil, fmt.Errorf("invalid filter_op: %s", params.FilterOp)
	}

	out := make([]Form, 0, len(forms))

	for _, form := range forms {
		if matchesFormFilter(form, key, op, params.FilterValue) {
			out = append(out, form)
		}
	}

	return out, nil
}

func matchesFormFilter(
	form Form,
	key string,
	op string,
	value string,
) bool {
	var field any

	switch key {
	case "status":
		field = form.Status
	case "name":
		field = form.Title
	case "opened_at":
		field = form.OpenedAt
	case "closed_at":
		field = form.ClosedAt
	case "archived_at":
		field = form.ArchivedAt
	default:
		return false
	}

	switch op {
	case "=":
		return fmt.Sprint(field) == value
	case "!=":
		return fmt.Sprint(field) != value
	case "is_null":
		return utils.IsNil(field)
	case "not_null":
		return !utils.IsNil(field)
	case "contains":
		return strings.Contains(strings.ToLower(fmt.Sprint(field)), strings.ToLower(value))
	default:
		return false
	}
}

func SortForms(forms []Form, params BulkGetParams) {
	desc := !strings.EqualFold(params.FilterOrder, "asc")

	slices.SortFunc(forms, func(a, b Form) int {
		if desc {
			return b.CreatedAt.Compare(a.CreatedAt)
		}

		return a.CreatedAt.Compare(b.CreatedAt)
	})
}
