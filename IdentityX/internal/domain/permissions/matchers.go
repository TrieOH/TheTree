package permissions

import (
	"GoAuth/internal/errx"
	"strings"

	"github.com/MintzyG/fail/v3"
)

// ObjectMatch checks if a request object matches a permission object pattern.
// Supports:
// - Exact match: event:123/activity:456
// - Modifier wildcard: event:123/activity:*
// - Segment wildcard: event:123/* (matches any namespace:specifier)
// - Global wildcard: *
// - Recursive wildcard (end only): event:123/** (matches any depth)
func ObjectMatch(permissionObj, requestObj string) bool {
	// Global wildcard matches everything
	if permissionObj == "*" {
		return true
	}

	// Normalize and split
	permParts := strings.Split(permissionObj, "/")
	reqParts := strings.Split(requestObj, "/")

	return matchParts(permParts, reqParts, 0, 0)
}

func matchParts(perm, req []string, pIdx, rIdx int) bool {
	// Base case: consumed all permission parts
	if pIdx >= len(perm) {
		// Match only if we also consumed all request parts
		return rIdx >= len(req)
	}

	permPart := perm[pIdx]

	// Recursive wildcard ** - matches zero or more segments
	if permPart == "**" {
		// Try matching 0 segments (end here) up to all remaining
		for i := rIdx; i <= len(req); i++ {
			if matchParts(perm, req, pIdx+1, i) {
				return true
			}
		}
		return false
	}

	// Single segment wildcard * - matches exactly ONE segment of any content
	if permPart == "*" {
		if rIdx >= len(req) {
			return false
		}
		return matchParts(perm, req, pIdx+1, rIdx+1)
	}

	// Must have a request part to compare for non-wildcards
	if rIdx >= len(req) {
		return false
	}

	reqPart := req[rIdx]

	// Check for modifier wildcards with ** (e.g., "activity:**")
	if strings.HasSuffix(permPart, ":**") {
		permNS := permPart[:len(permPart)-3] // Remove ":**"
		if !strings.HasPrefix(reqPart, permNS+":") {
			return false
		}
		// ** consumes this segment AND all remaining request segments
		// Immediately returns true since ** consumes everything remaining
		return true
	}

	// Check for modifier wildcards with * (e.g., "activity:*")
	if strings.HasSuffix(permPart, ":*") {
		permNS := permPart[:len(permPart)-2] // Remove ":*"
		if !strings.HasPrefix(reqPart, permNS+":") {
			return false
		}
		// Single modifier matches, move both pointers forward
		return matchParts(perm, req, pIdx+1, rIdx+1)
	}

	// Exact match required
	if permPart != reqPart {
		return false
	}

	// Continue matching next parts
	return matchParts(perm, req, pIdx+1, rIdx+1)
}

// ActionMatch checks if a request action matches a permission action.
// Supports:
// - Exact match: "edit", "attendance:mark", "admin:user:create"
// - Global wildcard: "*" matches everything
// - Prefix wildcard: "attendance:*" matches "attendance:mark", "attendance:check", etc.
// - "admin:*" matches "admin:create", "admin:delete", "admin:user:create", etc.
func ActionMatch(permissionAction, requestAction string) bool {
	if permissionAction == "*" {
		return true
	}

	// Reject empty request segments immediately
	if strings.HasSuffix(requestAction, ":") || strings.Contains(requestAction, "::") {
		return false
	}

	permParts := strings.Split(permissionAction, ":")
	reqParts := strings.Split(requestAction, ":")

	return matchActionParts(permParts, reqParts, 0, 0)
}

func matchActionParts(perm, req []string, pIdx, rIdx int) bool {
	// Consumed all permission segments
	if pIdx >= len(perm) {
		return rIdx >= len(req)
	}

	pSeg := perm[pIdx]

	// ** matches zero or more segments (recursive)
	if pSeg == "**" {
		for i := rIdx; i <= len(req); i++ {
			if matchActionParts(perm, req, pIdx+1, i) {
				return true
			}
		}
		return false
	}

	// Consumed all request but permission remains (and it's not **)
	if rIdx >= len(req) {
		return false
	}

	// * matches exactly one segment (any content)
	if pSeg == "*" {
		return matchActionParts(perm, req, pIdx+1, rIdx+1)
	}

	// Exact match required
	if pSeg != req[rIdx] {
		return false
	}

	return matchActionParts(perm, req, pIdx+1, rIdx+1)
}

// MustMatch is a validation helper for clean error messages
func MustMatch(permObj, permAction, reqObj, reqAction string) error {
	if !ObjectMatch(permObj, reqObj) {
		return fail.New(errx.PERMissionObjectMismatch).WithArgs(permObj, reqObj)
	}
	if !ActionMatch(permAction, reqAction) {
		return fail.New(errx.PERMissionActionMismatch).WithArgs(permAction, reqAction)
	}
	return nil
}
