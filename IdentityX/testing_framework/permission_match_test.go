package testing

import (
	"GoAuth/internal/domain/permissions"
	"testing"
)

func testPermissionMatching(t *testing.T, _ *TestSuite) {
	t.Run("TestObjectMatching", func(t *testing.T) {
		tests := []struct {
			name       string
			permission string
			request    string
			want       bool
		}{
			// Exact matches only
			{"ExactMatchSimple", "document", "document", true},
			{"ExactMatchEvent", "event", "event", true},
			{"ExactMatchUser", "user", "user", true},
			{"ExactMismatch", "document", "file", false},
			{"CaseSensitiveMismatch", "Document", "document", false},
			{"EmptyRequest", "document", "", false},
			{"EmptyPermission", "", "document", false},
			{"DifferentNames", "user_profile", "userprofile", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := permissions.ObjectMatch(tt.permission, tt.request)
				if got != tt.want {
					t.Errorf("ObjectMatch(%q, %q) = %v, want %v",
						tt.permission, tt.request, got, tt.want)
				}
			})
		}
	})

	t.Run("TestActionMatching", func(t *testing.T) {
		tests := []struct {
			name       string
			permission string
			request    string
			want       bool
		}{
			// Exact matches only
			{"ExactMatchRead", "read", "read", true},
			{"ExactMatchWrite", "write", "write", true},
			{"ExactMatchDelete", "delete", "delete", true},
			{"ExactMismatch", "read", "write", false},
			{"CaseSensitiveMismatch", "Read", "read", false},
			{"EmptyRequest", "read", "", false},
			{"EmptyPermission", "", "read", false},
			{"DifferentActions", "create_record", "createrecord", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := permissions.ActionMatch(tt.permission, tt.request)
				if got != tt.want {
					t.Errorf("ActionMatch(%q, %q) = %v, want %v",
						tt.permission, tt.request, got, tt.want)
				}
			})
		}
	})
}
