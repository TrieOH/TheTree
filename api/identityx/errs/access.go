package errs

import "errors"

// Access to resources outside the caller's reach, beyond the role checks
// (those return lib/authz.ErrInsufficientRole).
var (
	ErrActorOutsideProject = errors.New("insufficient permissions")
	ErrNotProfileOwner     = errors.New("insufficient permissions")
	ErrActorNotInProject   = errors.New("actor does not belong to this project")
)

// Project and organization membership.
var (
	ErrAddSelfToProject           = errors.New("users can't add themselves to projects")
	ErrAddOwnerToProject          = errors.New("owners can't be added to projects they own")
	ErrAlreadyProjectMember       = errors.New("user is already a member of the project")
	ErrRemoveSelfFromProject      = errors.New("cannot remove yourself from the project")
	ErrRemoveProjectOwner         = errors.New("cannot remove the owner of the project")
	ErrNotProjectMember           = errors.New("user is not a member of the project")
	ErrOrgOwnerAlreadyInProject   = errors.New("owner of the organization is already a member of the project")
	ErrOrgMemberAlreadyInProject  = errors.New("organization member is already a member of the project")
	ErrRemoveOrgOwnerFromProject  = errors.New("cannot remove the owner of the organization from the project")
	ErrRemoveOrgMemberFromProject = errors.New("cannot remove organization member from project")

	ErrAddSelfToOrg      = errors.New("users can't add themselves to organizations")
	ErrAddOwnerToOrg     = errors.New("owners can't be added to organizations they own")
	ErrAlreadyOrgMember  = errors.New("user is already a member of the organization")
	ErrRemoveSelfFromOrg = errors.New("cannot remove yourself from the organization")
	ErrRemoveOrgOwner    = errors.New("cannot remove the owner of the organization")
	ErrNotOrgMember      = errors.New("user is not a member of the organization")
)
