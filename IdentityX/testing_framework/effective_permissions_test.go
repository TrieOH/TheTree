package testing

import (
	"GoAuth/internal/apierr"
	"net/http"
	"testing"
	"time"
)

type specType map[string]interface{}

var (
	eventsScope     specType
	event1Scope     specType
	activitiesScope specType

	participantRole specType
	staffRole       specType
	adminRole       specType
	ownerRole       specType

	eventParticipantPermission    specType
	freeActivity1AttendPermission specType
	freeActivity2AttendPermission specType
	paidActivityAttendPermission  specType
	buyEventProductsPermission    specType
	buyTicketsPermission          specType

	activity1AttendanceMarkPermission      specType
	activity2AttendanceCheckPermission     specType
	allActivitiesAttendanceMarkPermission  specType
	allActivitiesAttendanceCheckPermission specType
	coordinateEventPermission              specType
	coordinatorDashboardPermission         specType

	createActivityPermission    specType
	editActivityPermission      specType
	deleteActivityPermission    specType
	assignEventRolePermission   specType
	administrateEventPermission specType

	fullEventAccessPermission specType
	createEventPermission     specType
)

// FIXME Sub scopes are possibly missing get permissions from parent scope with nil external_id
func testEffectivePermissions(t *testing.T, suite *TestSuite) {
	client := suite.NewClient(t)
	user := client.WithCredentials("effective@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("effective perms testing")

	projectID := user.projectID
	var eventsScopeID, event1ScopeID, activitiesScopeID string
	t.Run("CreateScopes", func(t *testing.T) {
		t.Run("CreateEventsScope", func(t *testing.T) {
			authClient := user.WithT(t)
			val := authClient.POST("/projects/" + projectID + "/scopes").
				WithBody(map[string]interface{}{
					"name":        "events",
					"external_id": nil,
				}).
				Expect(http.StatusCreated).
				HasMessage("Scope Created").
				RequireDataValue()

			eventsScope = map[string]interface{}{
				"id":          StoreString{Into: &eventsScopeID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "events",
				"external_id": nil,
				"type":        "project_scope",
				"created_at":  AnyDate{},
			}

			Validate(t, val, map[string]interface{}(eventsScope))
		})

		t.Run("CreateEvent1Scope", func(t *testing.T) {
			authClient := user.WithT(t)
			val := authClient.POST("/projects/" + projectID + "/scopes").
				WithBody(map[string]interface{}{
					"name":        "events",
					"external_id": "1",
				}).
				Expect(http.StatusCreated).
				HasMessage("Scope Created").
				RequireDataValue()

			event1Scope = map[string]interface{}{
				"id":          StoreString{Into: &event1ScopeID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "events",
				"external_id": "1",
				"type":        "project_scope",
				"created_at":  AnyDate{},
			}

			Validate(t, val, map[string]interface{}(event1Scope))
		})

		t.Run("CreateActivitiesScope", func(t *testing.T) {
			authClient := user.WithT(t)
			val := authClient.POST("/projects/" + projectID + "/scopes").
				WithBody(map[string]interface{}{
					"name":        "activities",
					"external_id": nil,
				}).
				Expect(http.StatusCreated).
				HasMessage("Scope Created").
				RequireDataValue()

			activitiesScope = map[string]interface{}{
				"id":          StoreString{Into: &activitiesScopeID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "activities",
				"external_id": nil,
				"type":        "project_scope",
				"created_at":  AnyDate{},
			}

			Validate(t, val, map[string]interface{}(activitiesScope))
		})
	})

	var participantRoleID, staffRoleID, adminRoleID, ownerRoleID string
	t.Run("CreateRoles", func(t *testing.T) {
		t.Run("CreateParticipantRole", func(t *testing.T) {
			authClient := user.WithT(t)
			val := authClient.POST("/projects/" + projectID + "/roles").
				WithBody(map[string]interface{}{
					"name":        "participant",
					"description": "can attend activities and participate in workshops",
				}).
				Expect(http.StatusCreated).
				HasMessage("Role Created").
				RequireDataValue()

			participantRole = map[string]interface{}{
				"id":          StoreString{Into: &participantRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "participant",
				"description": "can attend activities and participate in workshops",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},
			}

			Validate(t, val, map[string]interface{}(participantRole))
		})

		t.Run("CreateStaffRole", func(t *testing.T) {
			authClient := user.WithT(t)
			val := authClient.POST("/projects/" + projectID + "/roles").
				WithBody(map[string]interface{}{
					"name":        "staff",
					"description": "can mark and check attendance but cannot attend activities and participate in workshops",
				}).
				Expect(http.StatusCreated).
				HasMessage("Role Created").
				RequireDataValue()

			staffRole = map[string]interface{}{
				"id":          StoreString{Into: &staffRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "staff",
				"description": "can mark and check attendance but cannot attend activities and participate in workshops",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},
			}

			Validate(t, val, map[string]interface{}(staffRole))
		})

		t.Run("CreateAdminRole", func(t *testing.T) {
			authClient := user.WithT(t)
			val := authClient.POST("/projects/" + projectID + "/roles").
				WithBody(map[string]interface{}{
					"name":        "admin",
					"description": "can create edit and delete activities or products and can mark and check attendance but cannot attend activities and participate in workshops",
				}).
				Expect(http.StatusCreated).
				HasMessage("Role Created").
				RequireDataValue()

			adminRole = map[string]interface{}{
				"id":          StoreString{Into: &adminRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "admin",
				"description": "can create edit and delete activities or products and can mark and check attendance but cannot attend activities and participate in workshops",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},
			}

			Validate(t, val, map[string]interface{}(adminRole))
		})

		t.Run("CreateOwnerRole", func(t *testing.T) {
			authClient := user.WithT(t)
			val := authClient.POST("/projects/" + projectID + "/roles").
				WithBody(map[string]interface{}{
					"name":        "owner",
					"description": "can do anything in the event",
				}).
				Expect(http.StatusCreated).
				HasMessage("Role Created").
				RequireDataValue()

			ownerRole = map[string]interface{}{
				"id":          StoreString{Into: &ownerRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "owner",
				"description": "can do anything in the event",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},
			}

			Validate(t, val, map[string]interface{}(ownerRole))
		})
	})

	var event1ParticipatePermissionID string
	var freeActivity1AttendPermissionID, freeActivity2AttendPermissionID string
	var paidActivityAttendPermissionID string
	var buyEventProductsPermissionID, buyTicketsPermissionID string

	var activity1AttendanceMarkPermissionID, activity2AttendanceCheckPermissionID string
	var allActivitiesAttendanceMarkPermissionID, allActivitiesAttendanceCheckPermissionID string
	var coordinateEventPermissionID, coordinatorDashboardPermissionID string

	var createActivityPermissionID, editActivityPermissionID, deleteActivityPermissionID, assignEventRolePermissionID string
	var administrateEventPermissionID string

	var createEventPermissionID, fullEventAccessPermissionID string
	t.Run("CreatePermissions", func(t *testing.T) {
		t.Run("CreateParticipantPermissions", func(t *testing.T) {
			t.Run("CreateEventParticipantPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1",
						"action":     "participate",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				eventParticipantPermission = map[string]interface{}{
					"id":         StoreString{Into: &event1ParticipatePermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1",
					"action":     "participate",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(eventParticipantPermission))
			})

			// User receives free activities attend permissions for all activities as soon as they register to the event
			t.Run("CreateFreeUserActivityAttendPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1/activity:1",
						"action":     "attend",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				freeActivity1AttendPermission = map[string]interface{}{
					"id":         StoreString{Into: &freeActivity1AttendPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1/activity:1",
					"action":     "attend",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(freeActivity1AttendPermission))

				val = authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1/activity:2",
						"action":     "attend",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				freeActivity2AttendPermission = map[string]interface{}{
					"id":         StoreString{Into: &freeActivity2AttendPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1/activity:2",
					"action":     "attend",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(freeActivity2AttendPermission))
			})

			// Simulates a free user that bought an activity token and spent it on activity 2
			t.Run("CreatePaidActivityAttendPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1/activity:3",
						"action":     "attend",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				paidActivityAttendPermission = map[string]interface{}{
					"id":         StoreString{Into: &paidActivityAttendPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1/activity:3",
					"action":     "attend",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(paidActivityAttendPermission))
			})

			t.Run("CreateBuyEventProductsPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1/products:*",
						"action":     "buy",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				buyEventProductsPermission = map[string]interface{}{
					"id":         StoreString{Into: &buyEventProductsPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1/products:*",
					"action":     "buy",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(buyEventProductsPermission))
			})

			t.Run("CreateBuyTicketsPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "events:*",
						"action":     "tickets:buy",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				buyTicketsPermission = map[string]interface{}{
					"id":         StoreString{Into: &buyTicketsPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "events:*",
					"action":     "tickets:buy",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(buyTicketsPermission))
			})
		})

		t.Run("CreateStaffPermissions", func(t *testing.T) {
			t.Run("CreateCoordinateEventPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1",
						"action":     "coordinate",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				coordinateEventPermission = map[string]interface{}{
					"id":         StoreString{Into: &coordinateEventPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1",
					"action":     "coordinate",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(coordinateEventPermission))
			})
			t.Run("CreateActivity1AttendanceMarkPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1/activity:1",
						"action":     "attendance:mark",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				activity1AttendanceMarkPermission = map[string]interface{}{
					"id":         StoreString{Into: &activity1AttendanceMarkPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1/activity:1",
					"action":     "attendance:mark",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(activity1AttendanceMarkPermission))
			})

			t.Run("CreateActivity2AttendanceCheckPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1/activity:2",
						"action":     "attendance:check",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				activity2AttendanceCheckPermission = map[string]interface{}{
					"id":         StoreString{Into: &activity2AttendanceCheckPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1/activity:2",
					"action":     "attendance:check",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(activity2AttendanceCheckPermission))
			})

			t.Run("CreateAllActivitiesAttendanceMarkPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1/activity:*",
						"action":     "attendance:mark",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				allActivitiesAttendanceMarkPermission = map[string]interface{}{
					"id":         StoreString{Into: &allActivitiesAttendanceMarkPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1/activity:*",
					"action":     "attendance:mark",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(allActivitiesAttendanceMarkPermission))
			})

			t.Run("CreateAllActivitiesAttendanceCheckPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1/activity:*",
						"action":     "attendance:check",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				allActivitiesAttendanceCheckPermission = map[string]interface{}{
					"id":         StoreString{Into: &allActivitiesAttendanceCheckPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1/activity:*",
					"action":     "attendance:check",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(allActivitiesAttendanceCheckPermission))
			})

			t.Run("CreateCoordinatorDashboardPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:*",
						"action":     "coordinator_dashboard:access",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				coordinatorDashboardPermission = map[string]interface{}{
					"id":         StoreString{Into: &coordinatorDashboardPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:*",
					"action":     "coordinator_dashboard:access",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(coordinatorDashboardPermission))
			})
		})

		t.Run("CreateAdminPermissions", func(t *testing.T) {
			t.Run("CreateAdministrateEventPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1",
						"action":     "administrate",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				administrateEventPermission = map[string]interface{}{
					"id":         StoreString{Into: &administrateEventPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1",
					"action":     "administrate",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(administrateEventPermission))
			})
			t.Run("CreateCreateActivityPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1/activity:*",
						"action":     "create",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				createActivityPermission = map[string]interface{}{
					"id":         StoreString{Into: &createActivityPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1/activity:*",
					"action":     "create",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(createActivityPermission))
			})

			t.Run("CreateEditActivityPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1/activity:*",
						"action":     "edit",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				editActivityPermission = map[string]interface{}{
					"id":         StoreString{Into: &editActivityPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1/activity:*",
					"action":     "edit",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(editActivityPermission))
			})

			t.Run("CreateDeleteActivityPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1/activity:*",
						"action":     "delete",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				deleteActivityPermission = map[string]interface{}{
					"id":         StoreString{Into: &deleteActivityPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1/activity:*",
					"action":     "delete",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(deleteActivityPermission))
			})

			t.Run("CreateAssignEventRolePermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1/*",
						"action":     "role:assign",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				assignEventRolePermission = map[string]interface{}{
					"id":         StoreString{Into: &assignEventRolePermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1/*",
					"action":     "role:assign",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(assignEventRolePermission))
			})
		})

		t.Run("CreateOwnerPermissions", func(t *testing.T) {
			t.Run("CreateEventPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:*",
						"action":     "create",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				createEventPermission = map[string]interface{}{
					"id":         StoreString{Into: &createEventPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:*",
					"action":     "create",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(createEventPermission))
			})

			t.Run("CreateFullEventAccessPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object":     "event:1",
						"action":     "*",
						"conditions": nil,
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				fullEventAccessPermission = map[string]interface{}{
					"id":         StoreString{Into: &fullEventAccessPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event:1",
					"action":     "*",
					"conditions": nil,
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(fullEventAccessPermission))
			})
		})
	})

	t.Run("RegisterPermissionsToRoles", func(t *testing.T) {
		t.Run("RegisterPermissionsForParticipantRole", func(t *testing.T) {
			t.Run("RegisterAttendFreeActivity1Permission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				authClient.POST("/projects/" + projectID + "/roles/" + participantRoleID + "/permissions/" + freeActivity1AttendPermissionID).
					Expect(http.StatusOK).
					HasMessage("Added permission to role")
			})
			t.Run("RegisterAttendFreeActivity2Permission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				authClient.POST("/projects/" + projectID + "/roles/" + participantRoleID + "/permissions/" + freeActivity2AttendPermissionID).
					Expect(http.StatusOK).
					HasMessage("Added permission to role")
			})
			t.Run("RegisterParticipateOnEventPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				authClient.POST("/projects/" + projectID + "/roles/" + participantRoleID + "/permissions/" + event1ParticipatePermissionID).
					Expect(http.StatusOK).
					HasMessage("Added permission to role")
			})
			t.Run("RegisterBuyEventProductsPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				authClient.POST("/projects/" + projectID + "/roles/" + participantRoleID + "/permissions/" + buyEventProductsPermissionID).
					Expect(http.StatusOK).
					HasMessage("Added permission to role")
			})
		})
		t.Run("RegisterPermissionsForStaffRole", func(t *testing.T) {
			t.Run("RegisterCoordinateEventPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				authClient.POST("/projects/" + projectID + "/roles/" + staffRoleID + "/permissions/" + coordinateEventPermissionID).
					Expect(http.StatusOK).
					HasMessage("Added permission to role")
			})
		})
		t.Run("RegisterPermissionsForAdminRole", func(t *testing.T) {
			t.Run("RegisterAdministrateEventPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				authClient.POST("/projects/" + projectID + "/roles/" + adminRoleID + "/permissions/" + administrateEventPermissionID).
					Expect(http.StatusOK).
					HasMessage("Added permission to role")
			})
			t.Run("RegisterCreateActivityPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				authClient.POST("/projects/" + projectID + "/roles/" + adminRoleID + "/permissions/" + createActivityPermissionID).
					Expect(http.StatusOK).
					HasMessage("Added permission to role")
			})
		})
		t.Run("RegisterPermissionsForOwnerRole", func(t *testing.T) {
			t.Run("RegisterCreateEventPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				authClient.POST("/projects/" + projectID + "/roles/" + ownerRoleID + "/permissions/" + createEventPermissionID).
					Expect(http.StatusOK).
					HasMessage("Added permission to role")
			})
		})
	})

	var participant1, participant2, staff1, staff2, staff3, untrustedAdmin, trustedAdmin, owner *Client
	t.Run("CreateUsers", func(t *testing.T) {
		t.Run("CreateParticipant1", func(t *testing.T) {
			participant1 = suite.NewClient(t).WithCredentials("participant1@mail.com", ValidPassword).
				ProjectRegister(projectID).
				ProjectLogin(projectID)
		})
		t.Run("CreateParticipant2", func(t *testing.T) {
			participant2 = suite.NewClient(t).WithCredentials("participant2@mail.com", ValidPassword).
				ProjectRegister(projectID).
				ProjectLogin(projectID)
		})
		t.Run("CreateStaff1", func(t *testing.T) {
			staff1 = suite.NewClient(t).WithCredentials("staff1@mail.com", ValidPassword).
				ProjectRegister(projectID).
				ProjectLogin(projectID)
		})
		t.Run("CreateStaff2", func(t *testing.T) {
			staff2 = suite.NewClient(t).WithCredentials("staff2@mail.com", ValidPassword).
				ProjectRegister(projectID).
				ProjectLogin(projectID)
		})
		t.Run("CreateStaff3", func(t *testing.T) {
			staff3 = suite.NewClient(t).WithCredentials("staff3@mail.com", ValidPassword).
				ProjectRegister(projectID).
				ProjectLogin(projectID)
		})
		t.Run("CreateUntrustedAdmin", func(t *testing.T) {
			untrustedAdmin = suite.NewClient(t).WithCredentials("untrustedAdmin@mail.com", ValidPassword).
				ProjectRegister(projectID).
				ProjectLogin(projectID)
		})
		t.Run("CreateTrustedAdmin", func(t *testing.T) {
			trustedAdmin = suite.NewClient(t).WithCredentials("trustedAdmin@mail.com", ValidPassword).
				ProjectRegister(projectID).
				ProjectLogin(projectID)
		})
		t.Run("CreateOwner", func(t *testing.T) {
			owner = suite.NewClient(t).WithCredentials("eventOwner@mail.com", ValidPassword).
				ProjectRegister(projectID).
				ProjectLogin(projectID)
		})
	})

	t.Run("TestFullExampleUniventsFlow", func(t *testing.T) {
		t.Run("EventCreatorRegistersAccount", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *owner.projectUserID + "/roles").
				WithBody(map[string]interface{}{
					"role_id":  ownerRoleID,
					"scope_id": eventsScopeID,
				}).
				Expect(http.StatusOK).
				HasMessage("Added role to user")
		})
		t.Run("EventCreatorCreatesEvent", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *owner.projectUserID + "/permissions").
				WithBody(map[string]interface{}{
					"permission_id": fullEventAccessPermissionID,
					"scope_id":      event1ScopeID,
				}).
				Expect(http.StatusOK).
				HasMessage("Added permission to user")
		})
		t.Run("EventCreatorPromotesAUntrustedAdmin", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *untrustedAdmin.projectUserID + "/roles").
				WithBody(map[string]interface{}{
					"role_id":  adminRoleID,
					"scope_id": event1ScopeID,
				}).
				Expect(http.StatusOK).
				HasMessage("Added role to user")
		})
		t.Run("EventCreatorPromotesATrustedAdmin", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *trustedAdmin.projectUserID + "/roles").
				WithBody(map[string]interface{}{
					"role_id":  adminRoleID,
					"scope_id": event1ScopeID,
				}).
				Expect(http.StatusOK).
				HasMessage("Added role to user")

			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *trustedAdmin.projectUserID + "/permissions").
				WithBody(map[string]interface{}{
					"permission_id": editActivityPermissionID,
					"scope_id":      event1ScopeID,
				}).
				Expect(http.StatusOK).
				HasMessage("Added permission to user")

			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *trustedAdmin.projectUserID + "/permissions").
				WithBody(map[string]interface{}{
					"permission_id": deleteActivityPermissionID,
					"scope_id":      event1ScopeID,
				}).
				Expect(http.StatusOK).
				HasMessage("Added permission to user")

			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *trustedAdmin.projectUserID + "/permissions").
				WithBody(map[string]interface{}{
					"permission_id": assignEventRolePermissionID,
					"scope_id":      event1ScopeID,
				}).
				Expect(http.StatusOK).
				HasMessage("Added permission to user")
		})
		t.Run("TrustedAdminCreatesStaff", func(t *testing.T) {
			t.Run("TrustedAdminAddsStaffMember1", func(t *testing.T) {
				suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *staff1.projectUserID + "/roles").
					WithBody(map[string]interface{}{
						"role_id":  staffRoleID,
						"scope_id": event1ScopeID,
					}).
					Expect(http.StatusOK).
					HasMessage("Added role to user")

				suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *staff1.projectUserID + "/permissions").
					WithBody(map[string]interface{}{
						"permission_id": coordinatorDashboardPermissionID,
						"scope_id":      eventsScopeID,
					}).
					Expect(http.StatusOK).
					HasMessage("Added permission to user")

				suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *staff1.projectUserID + "/permissions").
					WithBody(map[string]interface{}{
						"permission_id": activity1AttendanceMarkPermissionID,
						"scope_id":      activitiesScopeID,
					}).
					Expect(http.StatusOK).
					HasMessage("Added permission to user")

				suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *staff1.projectUserID + "/permissions").
					WithBody(map[string]interface{}{
						"permission_id": activity2AttendanceCheckPermissionID,
						"scope_id":      activitiesScopeID,
					}).
					Expect(http.StatusOK).
					HasMessage("Added permission to user")
			})

			t.Run("TrustedAdminAddsStaffMember2", func(t *testing.T) {
				suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *staff2.projectUserID + "/roles").
					WithBody(map[string]interface{}{
						"role_id":  staffRoleID,
						"scope_id": event1ScopeID,
					}).
					Expect(http.StatusOK).
					HasMessage("Added role to user")

				suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *staff2.projectUserID + "/permissions").
					WithBody(map[string]interface{}{
						"permission_id": coordinatorDashboardPermissionID,
						"scope_id":      eventsScopeID,
					}).
					Expect(http.StatusOK).
					HasMessage("Added permission to user")

				suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *staff2.projectUserID + "/permissions").
					WithBody(map[string]interface{}{
						"permission_id": activity1AttendanceMarkPermissionID,
						"scope_id":      activitiesScopeID,
					}).
					Expect(http.StatusOK).
					HasMessage("Added permission to user")

				suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *staff2.projectUserID + "/permissions").
					WithBody(map[string]interface{}{
						"permission_id": allActivitiesAttendanceCheckPermissionID,
						"scope_id":      activitiesScopeID,
					}).
					Expect(http.StatusOK).
					HasMessage("Added permission to user")
			})

			t.Run("TrustedAdminAddsStaffMember3", func(t *testing.T) {
				suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *staff3.projectUserID + "/roles").
					WithBody(map[string]interface{}{
						"role_id":  staffRoleID,
						"scope_id": event1ScopeID,
					}).
					Expect(http.StatusOK).
					HasMessage("Added role to user")

				suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *staff3.projectUserID + "/permissions").
					WithBody(map[string]interface{}{
						"permission_id": coordinatorDashboardPermissionID,
						"scope_id":      eventsScopeID,
					}).
					Expect(http.StatusOK).
					HasMessage("Added permission to user")

				suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *staff3.projectUserID + "/permissions").
					WithBody(map[string]interface{}{
						"permission_id": allActivitiesAttendanceMarkPermissionID,
						"scope_id":      activitiesScopeID,
					}).
					Expect(http.StatusOK).
					HasMessage("Added permission to user")

				suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *staff3.projectUserID + "/permissions").
					WithBody(map[string]interface{}{
						"permission_id": allActivitiesAttendanceCheckPermissionID,
						"scope_id":      activitiesScopeID,
					}).
					Expect(http.StatusOK).
					HasMessage("Added permission to user")
			})

			t.Run("ParticipantsRegisterToTheEvent", func(t *testing.T) {
				t.Run("Participant1RegistersAsAFreeParticipant", func(t *testing.T) {
					suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *participant1.projectUserID + "/roles").
						WithBody(map[string]interface{}{
							"role_id":  participantRoleID,
							"scope_id": event1ScopeID,
						}).
						Expect(http.StatusOK).
						HasMessage("Added role to user")

					suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *participant1.projectUserID + "/permissions").
						WithBody(map[string]interface{}{
							"permission_id": buyTicketsPermissionID,
							"scope_id":      nil,
						}).
						Expect(http.StatusOK).
						HasMessage("Added permission to user")
				})
				t.Run("Participant2RegistersAndBuysActivity3", func(t *testing.T) {
					suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *participant2.projectUserID + "/roles").
						WithBody(map[string]interface{}{
							"role_id":  participantRoleID,
							"scope_id": event1ScopeID,
						}).
						Expect(http.StatusOK).
						HasMessage("Added role to user")

					suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *participant2.projectUserID + "/permissions").
						WithBody(map[string]interface{}{
							"permission_id": buyTicketsPermissionID,
							"scope_id":      nil,
						}).
						Expect(http.StatusOK).
						HasMessage("Added permission to user")

					suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *participant2.projectUserID + "/permissions").
						WithBody(map[string]interface{}{
							"permission_id": paidActivityAttendPermissionID,
							"scope_id":      event1ScopeID,
						}).
						Expect(http.StatusOK).
						HasMessage("Added permission to user")
				})
			})
		})
	})

	t.Run("VerifyEffectivePermissions", func(t *testing.T) {
		t.Run("VerifyParticipantsEffectivePermissions", func(t *testing.T) {
			t.Run("VerifyParticipant1EffectivePermissions", func(t *testing.T) {
				val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*participant1.projectUserID+"/permissions").
					WithQuery("scope_id", event1ScopeID).
					Expect(http.StatusOK).
					RequireDataValue()

				spec := []interface{}{
					map[string]interface{}(eventParticipantPermission),
					map[string]interface{}(freeActivity1AttendPermission),
					map[string]interface{}(freeActivity2AttendPermission),
					map[string]interface{}(buyEventProductsPermission),
					map[string]interface{}(buyTicketsPermission),
				}

				ValidateExact(t, val, spec)
			})

			t.Run("VerifyParticipant2EffectivePermissions", func(t *testing.T) {
				val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*participant2.projectUserID+"/permissions").
					WithQuery("scope_id", event1ScopeID).
					Expect(http.StatusOK).
					RequireDataValue()

				spec := []interface{}{
					map[string]interface{}(eventParticipantPermission),
					map[string]interface{}(freeActivity1AttendPermission),
					map[string]interface{}(freeActivity2AttendPermission),
					map[string]interface{}(paidActivityAttendPermission),
					map[string]interface{}(buyEventProductsPermission),
					map[string]interface{}(buyTicketsPermission),
				}

				ValidateExact(t, val, spec)
			})

			t.Run("VerifyStaff1EffectivePermissions", func(t *testing.T) {
				val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*staff1.projectUserID+"/permissions").
					WithQuery("scope_id", eventsScopeID).
					Expect(http.StatusOK).
					RequireDataValue()

				spec := []interface{}{
					map[string]interface{}(coordinatorDashboardPermission),
				}

				ValidateExact(t, val, spec)

				val = suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*staff1.projectUserID+"/permissions").
					WithQuery("scope_id", event1ScopeID).
					Expect(http.StatusOK).
					RequireDataValue()

				spec = []interface{}{
					map[string]interface{}(coordinatorDashboardPermission),
					map[string]interface{}(coordinateEventPermission),
				}

				ValidateExact(t, val, spec)

				val = suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*staff1.projectUserID+"/permissions").
					WithQuery("scope_id", activitiesScopeID).
					Expect(http.StatusOK).
					RequireDataValue()

				spec = []interface{}{
					map[string]interface{}(activity1AttendanceMarkPermission),
					map[string]interface{}(activity2AttendanceCheckPermission),
				}

				ValidateExact(t, val, spec)
			})

			t.Run("VerifyStaff2EffectivePermissions", func(t *testing.T) {
				val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*staff2.projectUserID+"/permissions").
					WithQuery("scope_id", eventsScopeID).
					Expect(http.StatusOK).
					RequireDataValue()

				spec := []interface{}{
					map[string]interface{}(coordinatorDashboardPermission),
				}

				ValidateExact(t, val, spec)

				val = suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*staff2.projectUserID+"/permissions").
					WithQuery("scope_id", event1ScopeID).
					Expect(http.StatusOK).
					RequireDataValue()

				spec = []interface{}{
					map[string]interface{}(coordinatorDashboardPermission),
					map[string]interface{}(coordinateEventPermission),
				}

				ValidateExact(t, val, spec)

				val = suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*staff2.projectUserID+"/permissions").
					WithQuery("scope_id", activitiesScopeID).
					Expect(http.StatusOK).
					RequireDataValue()

				spec = []interface{}{
					map[string]interface{}(allActivitiesAttendanceCheckPermission),
					map[string]interface{}(activity1AttendanceMarkPermission),
				}

				ValidateExact(t, val, spec)
			})

			t.Run("VerifyStaff3EffectivePermissions", func(t *testing.T) {
				val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*staff3.projectUserID+"/permissions").
					WithQuery("scope_id", eventsScopeID).
					Expect(http.StatusOK).
					RequireDataValue()

				spec := []interface{}{
					map[string]interface{}(coordinatorDashboardPermission),
				}

				ValidateExact(t, val, spec)

				val = suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*staff3.projectUserID+"/permissions").
					WithQuery("scope_id", event1ScopeID).
					Expect(http.StatusOK).
					RequireDataValue()

				spec = []interface{}{
					map[string]interface{}(coordinatorDashboardPermission),
					map[string]interface{}(coordinateEventPermission),
				}

				ValidateExact(t, val, spec)

				val = suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*staff3.projectUserID+"/permissions").
					WithQuery("scope_id", activitiesScopeID).
					Expect(http.StatusOK).
					RequireDataValue()

				spec = []interface{}{
					map[string]interface{}(allActivitiesAttendanceCheckPermission),
					map[string]interface{}(allActivitiesAttendanceMarkPermission),
				}

				ValidateExact(t, val, spec)
			})

			t.Run("VerifyUntrustedAdminEffectivePermissions", func(t *testing.T) {
				val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*untrustedAdmin.projectUserID+"/permissions").
					WithQuery("scope_id", event1ScopeID).
					Expect(http.StatusOK).
					RequireDataValue()

				spec := []interface{}{
					map[string]interface{}(administrateEventPermission),
					map[string]interface{}(createActivityPermission),
				}

				ValidateExact(t, val, spec)
			})

			t.Run("VerifyTrustedAdminEffectivePermissions", func(t *testing.T) {
				val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*trustedAdmin.projectUserID+"/permissions").
					WithQuery("scope_id", event1ScopeID).
					Expect(http.StatusOK).
					RequireDataValue()

				spec := []interface{}{
					map[string]interface{}(administrateEventPermission),
					map[string]interface{}(assignEventRolePermission),
					map[string]interface{}(createActivityPermission),
					map[string]interface{}(deleteActivityPermission),
					map[string]interface{}(editActivityPermission),
				}

				ValidateExact(t, val, spec)
			})

			t.Run("VerifyOwnerEffectivePermissions", func(t *testing.T) {
				val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*owner.projectUserID+"/permissions").
					WithQuery("scope_id", eventsScopeID).
					Expect(http.StatusOK).
					RequireDataValue()

				spec := []interface{}{
					map[string]interface{}(createEventPermission),
				}

				ValidateExact(t, val, spec)

				val = suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*owner.projectUserID+"/permissions").
					WithQuery("scope_id", event1ScopeID).
					Expect(http.StatusOK).
					RequireDataValue()

				spec = []interface{}{
					map[string]interface{}(createEventPermission),
					map[string]interface{}(fullEventAccessPermission),
				}

				ValidateExact(t, val, spec)
			})
		})
	})

	t.Run("VerifyAuthorizationChecks", func(t *testing.T) {
		t.Run("AllowScenario_EventAdminEditsEvent", func(t *testing.T) {
			// Owner has fullEventAccessPermission (action: *) on event1Scope
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID,
					"entity_id":  *owner.projectUserID,
					"object":     "event:1",
					"action":     "edit",
				}).
				Expect(http.StatusOK).
				HasMessage("Permission Granted")
		})

		t.Run("AllowScenario_TrustedAdminEditsActivity", func(t *testing.T) {
			// Trusted admin has editActivityPermission on event:1/activity:*
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID,
					"entity_id":  *trustedAdmin.projectUserID,
					"object":     "event:1/activity:999",
					"action":     "edit",
				}).
				Expect(http.StatusOK).
				HasMessage("Permission Granted")
		})

		t.Run("AllowScenario_StaffCoordinatesEvent", func(t *testing.T) {
			// Staff has coordinateEventPermission on event:1
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID,
					"entity_id":  *staff1.projectUserID,
					"object":     "event:1",
					"action":     "coordinate",
				}).
				Expect(http.StatusOK).
				HasMessage("Permission Granted")
		})

		t.Run("AllowScenario_ParticipantAttendsActivity", func(t *testing.T) {
			// Participant1 has freeActivity1AttendPermission via role
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID,
					"entity_id":  *participant1.projectUserID,
					"object":     "event:1/activity:1",
					"action":     "attend",
				}).
				Expect(http.StatusOK).
				HasMessage("Permission Granted")
		})

		t.Run("AllowScenario_ParticipantBuysTickets", func(t *testing.T) {
			// Both participants have buyTicketsPermission at nil scope (inherited anywhere)
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID, // Inherited from nil scope
					"entity_id":  *participant1.projectUserID,
					"object":     "events:1",
					"action":     "tickets:buy",
				}).
				Expect(http.StatusOK).
				HasMessage("Permission Granted")
		})

		t.Run("AllowScenario_StaffAccessesDashboardViaWildcard", func(t *testing.T) {
			// Staff has coordinatorDashboardPermission on event:*
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID,
					"entity_id":  *staff1.projectUserID,
					"object":     "event:1",
					"action":     "coordinator_dashboard:access",
				}).
				Expect(http.StatusOK).
				HasMessage("Permission Granted")
		})

		t.Run("AllowScenario_StaffMarksAttendanceWithSpecificPermission", func(t *testing.T) {
			// Staff1 has activity1AttendanceMarkPermission specifically for activity:1
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   activitiesScopeID,
					"entity_id":  *staff1.projectUserID,
					"object":     "event:1/activity:1",
					"action":     "attendance:mark",
				}).
				Expect(http.StatusOK).
				HasMessage("Permission Granted")
		})

		t.Run("AllowScenario_Staff3MarksAllAttendanceWithWildcard", func(t *testing.T) {
			// Staff3 has allActivitiesAttendanceMarkPermission (activity:*)
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   activitiesScopeID,
					"entity_id":  *staff3.projectUserID,
					"object":     "event:1/activity:999", // Any activity ID
					"action":     "attendance:mark",
				}).
				Expect(http.StatusOK).
				HasMessage("Permission Granted")
		})

		t.Run("DenyScenario_RandomUserEditsEvent", func(t *testing.T) {
			// Participant1 has no edit permissions on event
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID,
					"entity_id":  *participant1.projectUserID,
					"object":     "event:1",
					"action":     "edit",
				}).
				Expect(http.StatusForbidden).
				HasErrID(apierr.ID(apierr.PERMissionInsufficient.String())).
				HasMessage("Permission Denied").
				TraceContains("insufficient permissions")
		})

		t.Run("DenyScenario_ParticipantCannotCoordinate", func(t *testing.T) {
			// Participant1 has participantRole, no coordinate permission
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID,
					"entity_id":  *participant1.projectUserID,
					"object":     "event:1",
					"action":     "coordinate",
				}).
				Expect(http.StatusForbidden).
				HasErrID(apierr.ID(apierr.PERMissionInsufficient.String())).
				HasMessage("Permission Denied").
				TraceContains("insufficient permissions")
		})

		t.Run("DenyScenario_StaffCannotAdministrate", func(t *testing.T) {
			// Staff1 has staffRole, not adminRole
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID,
					"entity_id":  *staff1.projectUserID,
					"object":     "event:1",
					"action":     "administrate",
				}).
				Expect(http.StatusForbidden).
				HasErrID(apierr.ID(apierr.PERMissionInsufficient.String())).
				HasMessage("Permission Denied").
				TraceContains("insufficient permissions")
		})

		t.Run("DenyScenario_UntrustedAdminCannotDeleteActivity", func(t *testing.T) {
			// UntrustedAdmin has createActivityPermission but NOT deleteActivityPermission
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID,
					"entity_id":  *untrustedAdmin.projectUserID,
					"object":     "event:1/activity:1",
					"action":     "delete",
				}).
				Expect(http.StatusForbidden).
				HasErrID(apierr.ID(apierr.PERMissionInsufficient.String())).
				HasMessage("Permission Denied").
				TraceContains("insufficient permissions")
		})

		t.Run("DenyScenario_ParticipantCannotAccessUnpaidActivity", func(t *testing.T) {
			// Participant1 has freeActivity1 and freeActivity2, not activity:3
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID,
					"entity_id":  *participant1.projectUserID,
					"object":     "event:1/activity:3", // Paid activity
					"action":     "attend",
				}).
				Expect(http.StatusForbidden).
				HasErrID(apierr.ID(apierr.PERMissionInsufficient.String())).
				HasMessage("Permission Denied").
				TraceContains("insufficient permissions")
		})

		t.Run("DenyScenario_WrongScopeDenial", func(t *testing.T) {
			// Staff1 has permissions on event1Scope, not eventsScope
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   eventsScopeID, // Master scope (events:*), not event:1
					"entity_id":  *staff1.projectUserID,
					"object":     "event:1",
					"action":     "coordinate",
				}).
				Expect(http.StatusForbidden).
				HasErrID(apierr.ID(apierr.PERMissionInsufficient.String())).
				HasMessage("Permission Denied").
				TraceContains("insufficient permissions")
		})
	})

	time.Sleep(100 * time.Millisecond)
}
