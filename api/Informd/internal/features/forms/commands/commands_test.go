package commands

import (
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/mocks"
	"context"
	"errors"
	authz2 "lib/authz"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
)

type testDeps struct {
	svc        *CommandService
	forms      *mocks.MockFormsRepo
	steps      *mocks.MockStepRepo
	namespaces *mocks.MockNamespaceRepo
	perms      *mocks.MockChecker
}

func newTestDeps(t *testing.T) testDeps {
	t.Helper()
	forms := mocks.NewMockFormsRepo(t)
	steps := mocks.NewMockStepRepo(t)
	namespaces := mocks.NewMockNamespaceRepo(t)
	perms := mocks.NewMockChecker(t)
	tracer := noop.NewTracerProvider().Tracer("")

	svc := NewCommands(forms, steps, namespaces, perms, nil, tracer)
	return testDeps{
		svc:        svc,
		forms:      forms,
		steps:      steps,
		namespaces: namespaces,
		perms:      perms,
	}
}

func ctxWithUser(id uuid.UUID) context.Context {
	sub := &authz2.UserSubject{ID: id}
	return authz2.WithSubject(context.Background(), sub)
}

var errGeneric = errors.New("something went wrong")

func TestCreate_NoNamespace_Success(t *testing.T) {
	d := newTestDeps(t)
	userID := uuid.New()
	formID := uuid.New()
	ctx := ctxWithUser(userID)
	expected := &contracts.Form{ID: formID, Title: "My Form", OwnerID: userID}

	d.perms.EXPECT().
		Require(mock.Anything,
			authz2.Subject("user", userID),
			authz2.Permission("create_form"),
			authz2.Resource("user", userID.String()),
			map[string]any{"subject_id": userID.String()},
		).Return(nil)

	d.forms.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(expected, nil)

	d.perms.EXPECT().
		CreateRelation(mock.Anything, "form:"+formID.String()+"#parent_user@user:"+userID.String()).
		Return(nil)

	result, err := d.svc.Create(ctx, "My Form", nil)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreate_NoNamespace_NoSubjectInContext(t *testing.T) {
	d := newTestDeps(t)

	_, err := d.svc.Create(context.Background(), "My Form", nil)

	require.Error(t, err)
}

func TestCreate_NoNamespace_AuthzDenied(t *testing.T) {
	d := newTestDeps(t)
	userID := uuid.New()
	ctx := ctxWithUser(userID)

	d.perms.EXPECT().
		Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errGeneric)

	_, err := d.svc.Create(ctx, "My Form", nil)

	require.Error(t, err)
}

func TestCreate_NoNamespace_BadFormInput(t *testing.T) {
	d := newTestDeps(t)
	userID := uuid.New()
	ctx := ctxWithUser(userID)

	d.perms.EXPECT().
		Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	_, err := d.svc.Create(ctx, "", nil)

	require.Error(t, err)
}

func TestCreate_NoNamespace_DatabaseError(t *testing.T) {
	d := newTestDeps(t)
	userID := uuid.New()
	ctx := ctxWithUser(userID)

	d.perms.EXPECT().
		Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	d.forms.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(nil, errGeneric)

	_, err := d.svc.Create(ctx, "my form", nil)
	require.Error(t, err)
}

func TestCreate_NoNamespace_CouldntCreateRelation(t *testing.T) {
	d := newTestDeps(t)
	userID := uuid.New()
	formID := uuid.New()
	ctx := ctxWithUser(userID)
	created := &contracts.Form{ID: formID, Title: "my form", OwnerID: userID}

	d.perms.EXPECT().
		Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	d.forms.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(created, nil)

	d.perms.EXPECT().
		CreateRelation(mock.Anything, mock.Anything).
		Return(errGeneric)

	_, err := d.svc.Create(ctx, "my form", nil)
	require.Error(t, err)
}

func TestCreate_Namespace_Success(t *testing.T) {
	d := newTestDeps(t)
	userID := uuid.New()
	formID := uuid.New()
	ctx := ctxWithUser(userID)
	namespaceID := new(uuid.New())
	ns := &contracts.Namespace{ID: *namespaceID}
	created := &contracts.Form{ID: formID, Title: "my form", NamespaceID: namespaceID, OwnerID: userID}

	d.namespaces.EXPECT().
		GetByID(mock.Anything, *namespaceID).
		Return(ns, nil)

	d.perms.EXPECT().
		Require(mock.Anything,
			authz2.Subject("user", userID),
			authz2.Permission("create_form"),
			authz2.Resource("namespace", namespaceID.String()),
			map[string]any{"subject_id": userID.String()},
		).Return(nil)

	d.forms.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(created, nil)

	d.perms.EXPECT().
		CreateRelation(mock.Anything, "form:"+formID.String()+"#parent_namespace@namespace:"+namespaceID.String()).
		Return(nil)

	result, err := d.svc.Create(ctx, "my form", namespaceID)
	require.NoError(t, err)
	require.Equal(t, created, result)
}

func TestCreate_Namespace_NotFound(t *testing.T) {
	d := newTestDeps(t)
	userID := uuid.New()
	ctx := ctxWithUser(userID)
	namespaceID := new(uuid.New())

	d.namespaces.EXPECT().
		GetByID(mock.Anything, *namespaceID).
		Return(nil, errGeneric)

	_, err := d.svc.Create(ctx, "my form", namespaceID)
	require.Error(t, err)
}

func TestCreate_Namespace_AuthzDenied(t *testing.T) {
	d := newTestDeps(t)
	userID := uuid.New()
	ctx := ctxWithUser(userID)
	namespaceID := new(uuid.New())
	ns := &contracts.Namespace{ID: *namespaceID}

	d.namespaces.EXPECT().
		GetByID(mock.Anything, *namespaceID).
		Return(ns, nil)

	d.perms.EXPECT().
		Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errGeneric)

	_, err := d.svc.Create(ctx, "My Form", namespaceID)

	require.Error(t, err)
}

func TestCreate_Namespace_BadFormInput(t *testing.T) {
	d := newTestDeps(t)
	userID := uuid.New()
	ctx := ctxWithUser(userID)
	namespaceID := new(uuid.New())
	ns := &contracts.Namespace{ID: *namespaceID}

	d.namespaces.EXPECT().
		GetByID(mock.Anything, *namespaceID).
		Return(ns, nil)

	d.perms.EXPECT().
		Require(mock.Anything,
			authz2.Subject("user", userID),
			authz2.Permission("create_form"),
			authz2.Resource("namespace", namespaceID.String()),
			map[string]any{"subject_id": userID.String()},
		).Return(nil)

	_, err := d.svc.Create(ctx, "", namespaceID)

	require.Error(t, err)
}

func TestCreate_Namespace_DatabaseError(t *testing.T) {
	d := newTestDeps(t)
	userID := uuid.New()
	ctx := ctxWithUser(userID)
	namespaceID := new(uuid.New())
	ns := &contracts.Namespace{ID: *namespaceID}

	d.namespaces.EXPECT().
		GetByID(mock.Anything, *namespaceID).
		Return(ns, nil)

	d.perms.EXPECT().
		Require(mock.Anything,
			authz2.Subject("user", userID),
			authz2.Permission("create_form"),
			authz2.Resource("namespace", namespaceID.String()),
			map[string]any{"subject_id": userID.String()},
		).Return(nil)

	d.forms.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(nil, errGeneric)

	_, err := d.svc.Create(ctx, "my form", namespaceID)
	require.Error(t, err)
}

func TestCreate_Namespace_CouldntCreateRelation(t *testing.T) {
	d := newTestDeps(t)
	userID := uuid.New()
	formID := uuid.New()
	ctx := ctxWithUser(userID)
	namespaceID := new(uuid.New())
	ns := &contracts.Namespace{ID: *namespaceID}
	created := &contracts.Form{ID: formID, Title: "my form", OwnerID: userID, NamespaceID: namespaceID}

	d.namespaces.EXPECT().
		GetByID(mock.Anything, *namespaceID).
		Return(ns, nil)

	d.perms.EXPECT().
		Require(mock.Anything,
			authz2.Subject("user", userID),
			authz2.Permission("create_form"),
			authz2.Resource("namespace", namespaceID.String()),
			map[string]any{"subject_id": userID.String()},
		).Return(nil)

	d.forms.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(created, nil)

	d.perms.EXPECT().
		CreateRelation(mock.Anything, mock.Anything).
		Return(errGeneric)

	_, err := d.svc.Create(ctx, "my form", namespaceID)
	require.Error(t, err)
}

func TestCreateStep_Success(t *testing.T) {
	d := newTestDeps(t)
	userID := uuid.New()
	formID := uuid.New()
	ctx := ctxWithUser(userID)
	expected := &contracts.Form{ID: formID, Title: "My Form", OwnerID: userID}

	d.perms.EXPECT().
		Require(mock.Anything,
			authz2.Subject("user", userID),
			authz2.Permission("create_form"),
			authz2.Resource("user", userID.String()),
			map[string]any{"subject_id": userID.String()},
		).Return(nil)

	d.forms.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(expected, nil)

	d.perms.EXPECT().
		CreateRelation(mock.Anything, "form:"+formID.String()+"#parent_user@user:"+userID.String()).
		Return(nil)

	result, err := d.svc.Create(ctx, "My Form", nil)

	require.NoError(t, err)
	assert.Equal(t, expected, result)

	payload := contracts.CreateStepRequest{
		Title:        "Step 1",
		Description:  new("Please fill me out"),
		PositionHint: 1,
	}

	d.forms.EXPECT().
		GetByID(mock.Anything, formID).
		Return(expected, nil)

	d.perms.EXPECT().
		Require(mock.Anything,
			authz2.Subject("user", userID),
			authz2.Permission("edit"),
			authz2.Resource("user", userID.String()),
			map[string]any{"form": formID.String()},
		).Return(nil)

	_, err = d.svc.CreateStep(ctx, result.ID, payload)
	require.NoError(t, err)
}

func TestCreateStep_Fail_InsufficientPermissions(t *testing.T) {
	d := newTestDeps(t)
	userID := uuid.New()
	formID := uuid.New()
	ctx := ctxWithUser(userID)
	expected := &contracts.Form{ID: formID, Title: "My Form", OwnerID: userID}

	d.perms.EXPECT().
		Require(mock.Anything,
			authz2.Subject("user", userID),
			authz2.Permission("create_form"),
			authz2.Resource("user", userID.String()),
			map[string]any{"subject_id": userID.String()},
		).Return(nil)

	d.forms.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(expected, nil)

	d.perms.EXPECT().
		CreateRelation(mock.Anything, "form:"+formID.String()+"#parent_user@user:"+userID.String()).
		Return(nil)

	result, err := d.svc.Create(ctx, "My Form", nil)

	require.NoError(t, err)
	assert.Equal(t, expected, result)

	payload := contracts.CreateStepRequest{
		Title:        "Step 1",
		Description:  new("Please fill me out"),
		PositionHint: 1,
	}

	d.forms.EXPECT().
		GetByID(mock.Anything, formID).
		Return(expected, nil)

	d.perms.EXPECT().
		Require(mock.Anything,
			authz2.Subject("user", userID),
			authz2.Permission("edit"),
			authz2.Resource("form", formID.String()),
			nil,
		).Return(errGeneric)

	_, err = d.svc.CreateStep(ctx, result.ID, payload)
	require.Error(t, err)
}
