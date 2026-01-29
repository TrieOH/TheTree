package testing

import (
	"GoAuth/internal/domain/permissions"
	"testing"
	"time"
)

func testPermissionMatching(t *testing.T, _ *TestSuite) {
	t.Run("TestObjectMatching", func(t *testing.T) {
		tests := []struct {
			name       string
			permission string
			request    string
			want       bool
		}{
			// Global wildcard
			{"GlobalWildcardMatchesAll", "*", "event:123/activity:456", true},
			{"GlobalWildcardMatchesSimple", "*", "anything", true},

			// Exact matches
			{"ExactMatchSimple", "event:123", "event:123", true},
			{"ExactMatchNested", "event:123/activity:456", "event:123/activity:456", true},
			{"ExactMismatchValue", "event:123/activity:456", "event:123/activity:789", false},
			{"ExactMismatchNamespace", "event:123/activity:456", "event:123/session:456", false},

			// Modifier wildcards (name:*)
			{"ModifierWildcardMatches", "event:123/activity:*", "event:123/activity:99", true},
			{"ModifierWildcardDiffNS", "event:123/activity:*", "event:123/session:99", false},
			{"ModifierWildcardMulti", "event:*/activity:1", "event:999/activity:1", true},
			{"ModifierNoFalsePositive", "event:*/activity:1", "event:999/activity:2", false},

			// Segment wildcards (/* at end or middle conceptually, but regex restricts to end)
			{"SegmentWildcard", "event:123/*", "event:123/activity:1", true},
			{"SegmentWildcardAnyNS", "event:123/*", "event:123/anything:here", true},
			{"SegmentWildcardNoDeeper", "event:123/*", "event:123/activity:1/session:1", false},
			{"SegmentWildcardExactFail", "event:123/*", "event:123", false}, // * consumes exactly one segment

			// Recursive wildcards (/** at end)
			{"RecursiveMatchesSelf", "event:123/**", "event:123", true},
			{"RecursiveMatchesChild", "event:123/**", "event:123/activity:1", true},
			{"RecursiveMatchesDeep", "event:123/**", "event:123/activity:1/session:1/user:1", true},
			{"RecursiveMismatchRoot", "event:123/**", "event:456/activity:1", false},

			// Combined patterns
			{"ComplexWildcard", "event:*/activity:**", "event:123/activity:456/session:1", true},
			{"ComplexFail", "event:*/activity:**", "event:123/session:456", false},
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

	time.Sleep(100 * time.Millisecond)

	t.Run("TestActionMatching", func(t *testing.T) {
		tests := []struct {
			name       string
			permission string
			request    string
			want       bool
		}{
			// Global wildcard (*)
			{"WildcardGlobalMatchesSimple", "*", "delete", true},
			{"WildcardGlobalMatchesNested", "*", "attendance:mark", true},
			{"WildcardGlobalMatchesDeeplyNested", "*", "admin:user:create:verify", true},
			{"WildcardGlobalMatchesMultipleSegments", "*", "a:b:c:d:e:f", true},

			// Prefix wildcards (domain:*)
			{"WildcardPrefixMatchesNested", "attendance:*", "attendance:mark", true},
			{"WildcardPrefixMatchesCheck", "attendance:*", "attendance:check", true},
			{"WildcardPrefixNoMatchDeeplyNested", "admin:*", "admin:user:create", false},
			{"WildcardPrefixNoMatchTripleNested", "admin:*", "admin:role:delete:permanent", false},
			{"WildcardPrefixMatchesSingleSegment", "coordination:*", "coordination:start", true},
			{"WildcardPrefixNoMatchDifferentNamespace", "attendance:*", "session:mark", false},
			{"WildcardPrefixNoMatchDifferentAction", "attendance:*", "participation:mark", false},
			{"WildcardPrefixNoMatchEmptyAfterColon", "attendance:*", "attendance:", false},   // Regex should prevent this anyway
			{"WildcardPrefixExactMatchWithoutWildcard", "attendance:*", "attendance", false}, // "attendance" doesn't start with "attendance:"

			// Exact matches
			{"ExactMatchSimple", "edit", "edit", true},
			{"ExactMismatchSimple", "edit", "delete", false},
			{"ExactMatchNested", "attendance:mark", "attendance:mark", true},
			{"ExactMismatchNestedValue", "attendance:mark", "attendance:check", false},
			{"ExactMatchTripleNested", "admin:user:create", "admin:user:create", true},
			{"ExactMismatchTripleNestedEnd", "admin:user:create", "admin:user:delete", false},
			{"ExactMismatchTripleNestedMiddle", "admin:user:create", "admin:role:create", false},
			{"ExactMismatchNestedVersusSimple", "attendance:mark", "attendance", false},
			{"ExactMismatchSimpleVersusNested", "attendance", "attendance:mark", false},

			// Case sensitivity
			{"CaseSensitiveMismatch", "Edit", "edit", false},
			{"CaseSensitiveNestedMismatch", "Attendance:Mark", "attendance:mark", false},

			// Edge cases (valid per regex)
			{"SingleColonExact", "a:b", "a:b", true},
			{"SingleCharSegments", "a:b:c", "a:b:c", true},
			{"WildcardVersusLiteralStar", "*", "*", true},       // Both are global wildcard
			{"LiteralStarDoesntMatchAction", "*", "edit", true}, // Star matches everything

			// Star Segment
			{"StarMatchesSegment", "attendance:*", "attendance:mark", true},
			{"StarMatchesDifferent", "attendance:*", "attendance:check", true},
			{"StarDoesNotMatchNested", "attendance:*", "attendance:mark:verify", false},
			{"StarMiddleMatch", "admin:*:create", "admin:user:create", true},
			{"StarMiddleMismatch", "admin:*:create", "admin:user:delete", false},

			// Double wildcard (**) - recursive
			{"DoubleStarMatchesSelf", "admin:**", "admin", true},
			{"DoubleStarMatchesOne", "admin:**", "admin:create", true},
			{"DoubleStarMatchesDeep", "admin:**", "admin:user:role:create", true},
			{"DoubleStarAfterSpecific", "admin:user:**", "admin:user:create", true},
			{"DoubleStarAfterSpecificDeep", "admin:user:**", "admin:user:role:assign:permanent", true},
			{"DoubleStarWrongPrefix", "admin:**", "user:create", false},

			// Mixed patterns
			{"StarThenDoubleStar", "action:*:**", "action:custom:deep:nested", true},
			{"MultipleStars", "a:*:b:*:c", "a:1:b:2:c", true},
			{"MultipleStarsFail", "a:*:b:*:c", "a:1:b:2:d", false},
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

	time.Sleep(200 * time.Millisecond)
}
