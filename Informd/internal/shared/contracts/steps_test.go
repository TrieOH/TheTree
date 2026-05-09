package contracts

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStep_HappyPath_OptionalDescription(t *testing.T) {
	formID := uuid.New()
	s, err := NewStep(formID, "step one", nil, 1)
	require.NoError(t, err)

	require.NotNil(t, s)
	assert.Equal(t, formID, s.FormID)
	assert.Equal(t, "step one", s.Title)
	assert.Nil(t, s.Description)
	assert.Equal(t, 1, s.PositionHint)
}

func TestNewStep_HappyPath_WithDescription(t *testing.T) {
	formID := uuid.New()
	s, err := NewStep(formID, "step one", new("description"), 1)
	require.NoError(t, err)

	require.NotNil(t, s)
	assert.Equal(t, formID, s.FormID)
	assert.Equal(t, "step one", s.Title)
	assert.Equal(t, new("description"), s.Description)
	assert.Equal(t, 1, s.PositionHint)
}

func TestNewStep_ValidationErrors(t *testing.T) {
	validForm := uuid.New()

	tests := []struct {
		name         string
		formID       uuid.UUID
		title        string
		description  *string
		positionHint int
		wantErr      bool
	}{
		{
			name:         "zero value formID",
			formID:       uuid.Nil,
			title:        "step one",
			description:  nil,
			positionHint: 1,
			wantErr:      true,
		},
		{
			name:         "no title",
			formID:       validForm,
			title:        "",
			description:  nil,
			positionHint: 1,
			wantErr:      true,
		},
		{
			name:         "position hint zero",
			formID:       validForm,
			title:        "step one",
			description:  nil,
			positionHint: 0,
			wantErr:      true,
		},
		{
			name:         "position hint negative",
			formID:       validForm,
			title:        "step one",
			description:  nil,
			positionHint: -1,
			wantErr:      true,
		},
		{
			name:         "all validation errors",
			formID:       uuid.Nil,
			title:        "",
			description:  nil,
			positionHint: -1,
			wantErr:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := NewStep(tc.formID, tc.title, tc.description, tc.positionHint)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, f)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, f)
			}
		})
	}
}
