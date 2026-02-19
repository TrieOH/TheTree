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
			// Exact matches
			{"ExactMatchSimple", "document", "document", true},
			{"ExactMatchEvent", "event", "event", true},
			{"ExactMatchUser", "user", "user", true},
			{"ExactMismatch", "document", "file", false},
			{"CaseSensitiveMismatch", "Document", "document", false},
			{"EmptyRequest", "document", "", false},
			{"EmptyPermission", "", "document", false},
			{"DifferentNames", "user_profile", "userprofile", false},

			// Global wildcard (*) matches anything
			{"WildcardMatchesAny", "*", "document", true},
			{"WildcardMatchesEvent", "*", "event", true},
			{"WildcardMatchesAnything", "*", "anything_at_all", true},
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
			// Exact matches
			{"ExactMatchRead", "read", "read", true},
			{"ExactMatchWrite", "write", "write", true},
			{"ExactMatchDelete", "delete", "delete", true},
			{"ExactMismatch", "read", "write", false},
			{"CaseSensitiveMismatch", "Read", "read", false},
			{"EmptyRequest", "read", "", false},
			{"EmptyPermission", "", "read", false},
			{"DifferentActions", "create_record", "createrecord", false},

			// Global wildcard (*) matches anything
			{"WildcardMatchesRead", "*", "read", true},
			{"WildcardMatchesWrite", "*", "write", true},
			{"WildcardMatchesAnyAction", "*", "any_action", true},
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
