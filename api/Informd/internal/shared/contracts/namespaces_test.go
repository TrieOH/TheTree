package contracts

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNamespace_HappyPath(t *testing.T) {
	ownerID := uuid.New()
	f, err := NewNamespace(ownerID, "namespace")
	require.NoError(t, err)

	require.NotNil(t, f)
	assert.Equal(t, ownerID, f.OwnerID)
	assert.Equal(t, "namespace", f.Name)
	assert.True(t, f.CreatedAt.IsZero())
	assert.True(t, f.UpdatedAt.IsZero())
}

func TestNewNamespace_ValidationErrors(t *testing.T) {
	validOwner := uuid.New()

	tests := []struct {
		name          string
		ownerID       uuid.UUID
		namespaceName string
		wantErr       bool
	}{
		{
			name:          "zero value ownerID",
			ownerID:       uuid.Nil,
			namespaceName: "namespace",
			wantErr:       true,
		},
		{
			name:          "no name",
			ownerID:       validOwner,
			namespaceName: "",
			wantErr:       true,
		},
		{
			name:          "all required values empty",
			ownerID:       uuid.Nil,
			namespaceName: "",
			wantErr:       true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := NewNamespace(tc.ownerID, tc.namespaceName)
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
