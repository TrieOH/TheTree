package errs

import "errors"

// Form membership.
var (
	ErrAddSelfToForm                 = errors.New("users can't add themselves to forms")
	ErrFormOwnerAlreadyMember        = errors.New("owner of the form is already a member of the form")
	ErrNamespaceOwnerAlreadyInForm   = errors.New("owner of the namespace is already a member of the form")
	ErrNamespaceMemberAlreadyInForm  = errors.New("namespace member is already a member of the form")
	ErrAlreadyFormMember             = errors.New("user is already a member of the form")
	ErrRemoveSelfFromForm            = errors.New("users can't remove themselves from forms")
	ErrRemoveFormOwner               = errors.New("cannot remove owner of the form")
	ErrRemoveNamespaceOwnerFromForm  = errors.New("cannot remove owner of the namespace from form")
	ErrRemoveNamespaceMemberFromForm = errors.New("cannot remove namespace member from form")
	ErrNotFormMember                 = errors.New("user is not a member of the form")
)

// Namespace membership.
var (
	ErrAddSelfToNamespace     = errors.New("users can't add themselves to namespaces")
	ErrAddOwnerToNamespace    = errors.New("owners can't be added to namespaces they own")
	ErrAlreadyNamespaceMember = errors.New("user is already a member of the namespace")
	ErrNotNamespaceMember     = errors.New("user is not a member of the namespace")
)
