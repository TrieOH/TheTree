package models

import (
	"testing"
	"uuid"

	"lib/validator"

	"github.com/stretchr/testify/require"
)

// The operations build these models directly and run them through the
// shared validator before they're written; these cases pin the tags.
func TestModelValidation(t *testing.T) {
	validator.SetupValidator()
	id := uuid.New()

	tests := []struct {
		name    string
		model   any
		wantErr bool
	}{
		{"form", &Form{OwnerID: id, CreatedBy: id, Title: "a form"}, false},
		{"form in a namespace", &Form{NamespaceID: new(uuid.New()), OwnerID: id, CreatedBy: id, Title: "a form"}, false},
		{"form in the zero namespace", &Form{NamespaceID: new(uuid.Nil()), OwnerID: id, CreatedBy: id, Title: "a form"}, false},
		{"form without owner", &Form{CreatedBy: id, Title: "a form"}, true},
		{"form without title", &Form{OwnerID: id, CreatedBy: id}, true},

		{"namespace", &Namespace{OwnerID: id, Name: "team"}, false},
		{"namespace without owner", &Namespace{Name: "team"}, true},
		{"namespace without name", &Namespace{OwnerID: id}, true},

		{"step", &Step{FormID: id, Title: "step one", PositionHint: 1}, false},
		{"step without form", &Step{Title: "step one", PositionHint: 1}, true},
		{"step without title", &Step{FormID: id, PositionHint: 1}, true},
		{"step at position zero", &Step{FormID: id, Title: "step one"}, true},
		{"step at a negative position", &Step{FormID: id, Title: "step one", PositionHint: -1}, true},

		{"field", &Field{StepID: id, Key: "email", Title: "Email", PositionHint: 1, Type: FieldTypeString}, false},
		{"field without key", &Field{StepID: id, Title: "Email", PositionHint: 1, Type: FieldTypeString}, true},
		{"field without type", &Field{StepID: id, Key: "email", Title: "Email", PositionHint: 1}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.Validate(tc.model)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
