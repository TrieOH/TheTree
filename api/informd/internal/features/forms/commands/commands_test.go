package commands

import (
	"Informd/generated/mocks"
	"Informd/models"
	"context"
	"errors"
	"lib/authz"
	"testing"

	"github.com/google/uuid"
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
	sub := &authz.UserSubject{ID: id}
	return authz.WithSubject(context.Background(), sub)
}

var errGeneric = errors.New("something went wrong")

func formMatcher(ownerID uuid.UUID, namespaceID *uuid.UUID, title string) interface{} {
	return mock.MatchedBy(func(f models.Form) bool {
		nsMatch := (namespaceID == nil && f.NamespaceID == nil) || (namespaceID != nil && f.NamespaceID != nil && *f.NamespaceID == *namespaceID)
		return f.Title == title &&
			f.OwnerID == ownerID &&
			nsMatch &&
			f.Status == models.FormStatusDraft
	})
}

func expectRequire(d testDeps, err error) {
	d.perms.EXPECT().Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(err)
}

func TestCreate_NoNamespace(t *testing.T) {
	userID := uuid.New()
	formID := uuid.New()
	form := &models.Form{ID: formID, Title: "My Form", OwnerID: userID}

	tests := []struct {
		name    string
		title   string
		setup   func(d testDeps)
		wantErr bool
	}{
		{
			name:  "success",
			title: "My Form",
			setup: func(d testDeps) {
				expectRequire(d, nil)
				d.forms.EXPECT().Create(mock.Anything, formMatcher(userID, nil, "My Form")).Return(form, nil)
				d.perms.EXPECT().CreateRelation(mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "no subject in context",
			title:   "My Form",
			setup:   func(d testDeps) {},
			wantErr: true,
		},
		{
			name:  "authz denied",
			title: "My Form",
			setup: func(d testDeps) {
				expectRequire(d, errGeneric)
			},
			wantErr: true,
		},
		{
			name:  "bad input",
			title: "",
			setup: func(d testDeps) {
				expectRequire(d, nil)
			},
			wantErr: true,
		},
		{
			name:  "db error",
			title: "My Form",
			setup: func(d testDeps) {
				expectRequire(d, nil)
				d.forms.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, errGeneric)
			},
			wantErr: true,
		},
		{
			name:  "could not create relation",
			title: "My Form",
			setup: func(d testDeps) {
				expectRequire(d, nil)
				d.forms.EXPECT().Create(mock.Anything, mock.Anything).Return(form, nil)
				d.perms.EXPECT().CreateRelation(mock.Anything, mock.Anything).Return(errGeneric)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newTestDeps(t)
			tt.setup(d)
			ctx := ctxWithUser(userID)
			if tt.name == "no subject in context" {
				ctx = context.Background()
			}
			_, err := d.svc.Create(ctx, tt.title, nil)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCreate_Namespace(t *testing.T) {
	userID := uuid.New()
	formID := uuid.New()
	form := &models.Form{ID: formID, Title: "My Direct Form", OwnerID: userID}
	namespaceID := uuid.New()
	namespace := &models.Namespace{ID: namespaceID}

	tests := []struct {
		name    string
		title   string
		setup   func(d testDeps)
		wantErr bool
	}{
		{
			name:  "success",
			title: "My Direct Form",
			setup: func(d testDeps) {
				d.namespaces.EXPECT().GetByID(mock.Anything, namespaceID).Return(namespace, nil)
				d.perms.EXPECT().Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				d.forms.EXPECT().Create(mock.Anything, formMatcher(userID, &namespaceID, "My Direct Form")).Return(form, nil)
				d.perms.EXPECT().CreateRelation(mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "no subject in context",
			title:   "My Form",
			setup:   func(d testDeps) {},
			wantErr: true,
		},
		{
			name:  "namespace not found",
			title: "My Direct Form",
			setup: func(d testDeps) {
				d.namespaces.EXPECT().GetByID(mock.Anything, namespaceID).Return(nil, errGeneric)
			},
			wantErr: true,
		},
		{
			name:  "permission denied",
			title: "My Direct Form",
			setup: func(d testDeps) {
				d.namespaces.EXPECT().GetByID(mock.Anything, namespaceID).Return(namespace, nil)
				d.perms.EXPECT().Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errGeneric)
			},
			wantErr: true,
		},
		{
			name:  "bad input",
			title: "",
			setup: func(d testDeps) {
				d.namespaces.EXPECT().GetByID(mock.Anything, namespaceID).Return(namespace, nil)
				d.perms.EXPECT().Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: true,
		},
		{
			name:  "db error",
			title: "My Direct Form",
			setup: func(d testDeps) {
				d.perms.EXPECT().Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				d.namespaces.EXPECT().GetByID(mock.Anything, namespaceID).Return(namespace, nil)
				d.forms.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, errGeneric)
			},
			wantErr: true,
		},
		{
			name:  "could not create relation",
			title: "My Direct Form",
			setup: func(d testDeps) {
				d.namespaces.EXPECT().GetByID(mock.Anything, namespaceID).Return(namespace, nil)
				d.perms.EXPECT().Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				d.forms.EXPECT().Create(mock.Anything, mock.Anything).Return(form, nil)
				d.perms.EXPECT().CreateRelation(mock.Anything, mock.Anything).Return(errGeneric)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newTestDeps(t)
			tt.setup(d)
			ctx := ctxWithUser(userID)
			if tt.name == "no subject in context" {
				ctx = context.Background()
			}
			_, err := d.svc.Create(ctx, tt.title, &namespaceID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_CreateStep(t *testing.T) {
	userID := uuid.New()
	formID := uuid.New()
	form, _ := models.NewForm(new(uuid.New()), userID, "My Form")
	step, _ := models.NewStep(formID, "Step 1", new("My Description"), 1)
	tests := []struct {
		name         string
		title        string
		description  *string
		positionHint int
		setup        func(d testDeps)
		wantErr      bool
	}{
		{
			name:         "success",
			title:        "Step 1",
			description:  new("My Description"),
			positionHint: 1,
			setup: func(d testDeps) {
				d.perms.EXPECT().Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				d.forms.EXPECT().GetByID(mock.Anything, formID).Return(form, nil)
				d.steps.EXPECT().Create(mock.Anything, mock.Anything).Return(step, nil)
			},
			wantErr: false,
		},
		{
			name:         "no subject in context",
			title:        "Step 1",
			description:  new("My Description"),
			positionHint: 1,
			setup:        func(d testDeps) {},
			wantErr:      true,
		},
		{
			name:         "insufficient permission",
			title:        "Step 1",
			description:  new("My Description"),
			positionHint: 1,
			setup: func(d testDeps) {
				d.forms.EXPECT().GetByID(mock.Anything, formID).Return(form, nil)
				d.perms.EXPECT().Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errGeneric)
			},
			wantErr: true,
		},
		{
			name:         "form not found",
			title:        "Step 1",
			description:  new("My Description"),
			positionHint: 1,
			setup: func(d testDeps) {
				d.forms.EXPECT().GetByID(mock.Anything, formID).Return(nil, errGeneric)
			},
			wantErr: true,
		},
		{
			name:         "empty title",
			title:        "",
			description:  new("My Description"),
			positionHint: 1,
			setup: func(d testDeps) {
				d.forms.EXPECT().GetByID(mock.Anything, formID).Return(form, nil)
				d.perms.EXPECT().Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: true,
		},
		{
			name:         "position hint zero",
			title:        "Step 1",
			description:  new("My Description"),
			positionHint: 0,
			setup: func(d testDeps) {
				d.forms.EXPECT().GetByID(mock.Anything, formID).Return(form, nil)
				d.perms.EXPECT().Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: true,
		},
		{
			name:         "position hint negative",
			title:        "Step 1",
			description:  new("My Description"),
			positionHint: -1,
			setup: func(d testDeps) {
				d.forms.EXPECT().GetByID(mock.Anything, formID).Return(form, nil)
				d.perms.EXPECT().Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: true,
		},
		{
			name:         "db error",
			title:        "Step 1",
			description:  new("My Description"),
			positionHint: 1,
			setup: func(d testDeps) {
				d.forms.EXPECT().GetByID(mock.Anything, formID).Return(form, nil)
				d.perms.EXPECT().Require(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				d.steps.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, errGeneric)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newTestDeps(t)
			tt.setup(d)
			ctx := ctxWithUser(userID)
			if tt.name == "no subject in context" {
				ctx = context.Background()
			}
			_, err := d.svc.CreateStep(ctx, formID, models.CreateStepRequest{
				Title:        tt.title,
				Description:  tt.description,
				PositionHint: tt.positionHint,
			})
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
