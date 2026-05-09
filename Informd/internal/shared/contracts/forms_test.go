package contracts

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewForm_Namespaced_HappyPath(t *testing.T) {
	ownerID := uuid.New()
	namespaceID := new(uuid.New())
	f, err := NewForm(namespaceID, ownerID, "namespaced form")
	require.NoError(t, err)

	require.NotNil(t, f)
	assert.Equal(t, ownerID, f.OwnerID)
	assert.Equal(t, namespaceID, f.NamespaceID)
	assert.Equal(t, "namespaced form", f.Title)
	assert.Equal(t, FormStatusDraft, f.Status)
	assert.Nil(t, f.OpenedAt)
	assert.Nil(t, f.ClosedAt)
	assert.Nil(t, f.ArchivedAt)
	assert.True(t, f.CreatedAt.IsZero())
	assert.True(t, f.UpdatedAt.IsZero())
}

func TestNewForm_Not_Namespaced_HappyPath(t *testing.T) {
	ownerID := uuid.New()
	f, err := NewForm(nil, ownerID, "not namespaced form")
	require.NoError(t, err)

	require.NotNil(t, f)
	assert.Equal(t, ownerID, f.OwnerID)
	assert.Nil(t, f.NamespaceID)
	assert.Equal(t, "not namespaced form", f.Title)
	assert.Equal(t, FormStatusDraft, f.Status)
	assert.Nil(t, f.OpenedAt)
	assert.Nil(t, f.ClosedAt)
	assert.Nil(t, f.ArchivedAt)
	assert.True(t, f.CreatedAt.IsZero())
	assert.True(t, f.UpdatedAt.IsZero())
}

func TestNewForm_ValidationErrors(t *testing.T) {
	validOwner := uuid.New()
	validNamespace := new(uuid.New())

	tests := []struct {
		name        string
		ownerID     uuid.UUID
		namespaceID *uuid.UUID
		formName    string
		wantErr     bool
	}{
		{
			name:        "no ownerID",
			ownerID:     uuid.Nil,
			namespaceID: validNamespace,
			formName:    "form",
			wantErr:     true,
		},
		{
			name:        "no name",
			ownerID:     validOwner,
			namespaceID: validNamespace,
			formName:    "",
			wantErr:     true,
		},
		{
			name:        "owner and name empty",
			ownerID:     uuid.Nil,
			namespaceID: validNamespace,
			formName:    "",
			wantErr:     true,
		},
		{
			name:        "zero value namespaceID",
			ownerID:     validOwner,
			namespaceID: new(uuid.Nil),
			formName:    "form",
			wantErr:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := NewForm(tc.namespaceID, tc.ownerID, tc.formName)
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
