package permissions

import (
	"GoAuth/internal/errx"
	"github.com/MintzyG/fail/v3"
)

// ObjectMatch checks if a request object matches a permission object.
// Supports global wildcard (*) which matches any object.
func ObjectMatch(permissionObj, requestObj string) bool {
	// Global wildcard matches everything
	if permissionObj == "*" {
		return true
	}
	return permissionObj == requestObj
}

// ActionMatch checks if a request action matches a permission action.
// Supports global wildcard (*) which matches any action.
func ActionMatch(permissionAction, requestAction string) bool {
	// Global wildcard matches everything
	if permissionAction == "*" {
		return true
	}
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
