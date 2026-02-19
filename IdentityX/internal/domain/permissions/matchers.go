package permissions

import (
	"GoAuth/internal/errx"
	"github.com/MintzyG/fail/v3"
)

// ObjectMatch checks if a request object matches a permission object.
// Now uses simple equality - no wildcards or patterns supported.
func ObjectMatch(permissionObj, requestObj string) bool {
	return permissionObj == requestObj
}

// ActionMatch checks if a request action matches a permission action.
// Now uses simple equality - no wildcards or patterns supported.
func ActionMatch(permissionAction, requestAction string) bool {
	return permissionAction == requestAction
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
