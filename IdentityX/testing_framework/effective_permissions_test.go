package testing

import (
	"GoAuth/internal/errx"
	"net/http"
	"testing"
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

	// Simplified permission names
	eventParticipatePermission specType
	activityAttendPermission   specType
	productBuyPermission       specType
	ticketBuyPermission        specType
	attendanceMarkPermission   specType
	attendanceCheckPermission  specType
	dashboardAccessPermission  specType
	activityCreatePermission   specType
	activityEditPermission     specType
	activityDeletePermission   specType
	eventCreatePermission      specType
	administratePermission     specType
)

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
					"name":        "event1",
					"external_id": "1",
					"parent_id":   eventsScopeID,
				}).
				Expect(http.StatusCreated).
				HasMessage("Scope Created").
				RequireDataValue()

			event1Scope = map[string]interface{}{
				"id":          StoreString{Into: &event1ScopeID, Matcher: AnyUUID{}},
				"parent_id":   AsString{Value: eventsScopeID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "event1",
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
					"parent_id":   event1ScopeID,
				}).
				Expect(http.StatusCreated).
				HasMessage("Scope Created").
				RequireDataValue()

			activitiesScope = map[string]interface{}{
				"id":          StoreString{Into: &activitiesScopeID, Matcher: AnyUUID{}},
				"parent_id":   AsString{Value: event1ScopeID, Matcher: AnyUUID{}},
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
					"description": "can attend activities",
				}).
				Expect(http.StatusCreated).
				HasMessage("Role Created").
				RequireDataValue()

			participantRole = map[string]interface{}{
				"id":          StoreString{Into: &participantRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "participant",
				"description": "can attend activities",
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
					"description": "can check attendance",
				}).
				Expect(http.StatusCreated).
				HasMessage("Role Created").
				RequireDataValue()

			staffRole = map[string]interface{}{
				"id":          StoreString{Into: &staffRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "staff",
				"description": "can check attendance",
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
					"description": "can create and edit activities",
				}).
				Expect(http.StatusCreated).
				HasMessage("Role Created").
				RequireDataValue()

			adminRole = map[string]interface{}{
				"id":          StoreString{Into: &adminRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "admin",
				"description": "can create and edit activities",
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
					"description": "can do anything",
				}).
				Expect(http.StatusCreated).
				HasMessage("Role Created").
				RequireDataValue()

			ownerRole = map[string]interface{}{
				"id":          StoreString{Into: &ownerRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "owner",
				"description": "can do anything",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},
			}

			Validate(t, val, map[string]interface{}(ownerRole))
		})
	})

	var (
		eventParticipatePermissionID string
		activityAttendPermissionID   string
		productBuyPermissionID       string
		ticketBuyPermissionID        string
		attendanceMarkPermissionID   string
		attendanceCheckPermissionID  string
		dashboardAccessPermissionID  string
		activityCreatePermissionID   string
		activityEditPermissionID     string
		activityDeletePermissionID   string
		eventCreatePermissionID      string
		administratePermissionID     string
	)

	t.Run("CreatePermissions", func(t *testing.T) {
		t.Run("CreateParticipantPermissions", func(t *testing.T) {
			t.Run("CreateEventParticipatePermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object": "event",
						"action": "participate",
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				eventParticipatePermission = map[string]interface{}{
					"id":         StoreString{Into: &eventParticipatePermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event",
					"action":     "participate",
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(eventParticipatePermission))
			})

			t.Run("CreateActivityAttendPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object": "activity",
						"action": "attend",
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				activityAttendPermission = map[string]interface{}{
					"id":         StoreString{Into: &activityAttendPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "activity",
					"action":     "attend",
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(activityAttendPermission))
			})

			t.Run("CreateProductBuyPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object": "product",
						"action": "buy",
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				productBuyPermission = map[string]interface{}{
					"id":         StoreString{Into: &productBuyPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "product",
					"action":     "buy",
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(productBuyPermission))
			})

			t.Run("CreateTicketBuyPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object": "ticket",
						"action": "buy",
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				ticketBuyPermission = map[string]interface{}{
					"id":         StoreString{Into: &ticketBuyPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "ticket",
					"action":     "buy",
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(ticketBuyPermission))
			})
		})

		t.Run("CreateStaffPermissions", func(t *testing.T) {
			t.Run("CreateAttendanceMarkPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object": "attendance",
						"action": "mark",
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				attendanceMarkPermission = map[string]interface{}{
					"id":         StoreString{Into: &attendanceMarkPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "attendance",
					"action":     "mark",
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(attendanceMarkPermission))
			})

			t.Run("CreateAttendanceCheckPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object": "attendance",
						"action": "check",
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				attendanceCheckPermission = map[string]interface{}{
					"id":         StoreString{Into: &attendanceCheckPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "attendance",
					"action":     "check",
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(attendanceCheckPermission))
			})

			t.Run("CreateDashboardAccessPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object": "dashboard",
						"action": "access",
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				dashboardAccessPermission = map[string]interface{}{
					"id":         StoreString{Into: &dashboardAccessPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "dashboard",
					"action":     "access",
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(dashboardAccessPermission))
			})
		})

		t.Run("CreateAdminPermissions", func(t *testing.T) {
			t.Run("CreateActivityCreatePermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object": "activity",
						"action": "create",
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				activityCreatePermission = map[string]interface{}{
					"id":         StoreString{Into: &activityCreatePermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "activity",
					"action":     "create",
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(activityCreatePermission))
			})

			t.Run("CreateActivityEditPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object": "activity",
						"action": "edit",
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				activityEditPermission = map[string]interface{}{
					"id":         StoreString{Into: &activityEditPermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "activity",
					"action":     "edit",
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(activityEditPermission))
			})

			t.Run("CreateActivityDeletePermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object": "activity",
						"action": "delete",
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				activityDeletePermission = map[string]interface{}{
					"id":         StoreString{Into: &activityDeletePermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "activity",
					"action":     "delete",
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(activityDeletePermission))
			})
		})

		t.Run("CreateOwnerPermissions", func(t *testing.T) {
			t.Run("CreateEventCreatePermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object": "event",
						"action": "create",
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				eventCreatePermission = map[string]interface{}{
					"id":         StoreString{Into: &eventCreatePermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event",
					"action":     "create",
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(eventCreatePermission))
			})

			t.Run("CreateAdministratePermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				val := authClient.POST("/projects/" + projectID + "/permissions").
					WithBody(map[string]interface{}{
						"object": "event",
						"action": "administrate",
					}).
					Expect(http.StatusCreated).
					HasMessage("Permission Created").
					RequireDataValue()

				administratePermission = map[string]interface{}{
					"id":         StoreString{Into: &administratePermissionID, Matcher: AnyUUID{}},
					"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
					"object":     "event",
					"action":     "administrate",
					"created_at": AnyDate{},
				}

				Validate(t, val, map[string]interface{}(administratePermission))
			})
		})
	})

	t.Run("RegisterPermissionsToRoles", func(t *testing.T) {
		t.Run("RegisterPermissionsForParticipantRole", func(t *testing.T) {
			t.Run("RegisterActivityAttendPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				authClient.POST("/projects/" + projectID + "/roles/" + participantRoleID + "/permissions/" + activityAttendPermissionID).
					Expect(http.StatusOK).
					HasMessage("Added permission to role")
			})
			t.Run("RegisterEventParticipatePermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				authClient.POST("/projects/" + projectID + "/roles/" + participantRoleID + "/permissions/" + eventParticipatePermissionID).
					Expect(http.StatusOK).
					HasMessage("Added permission to role")
			})
			t.Run("RegisterProductBuyPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				authClient.POST("/projects/" + projectID + "/roles/" + participantRoleID + "/permissions/" + productBuyPermissionID).
					Expect(http.StatusOK).
					HasMessage("Added permission to role")
			})
		})

		t.Run("RegisterPermissionsForStaffRole", func(t *testing.T) {
			t.Run("RegisterAttendanceCheckPermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				authClient.POST("/projects/" + projectID + "/roles/" + staffRoleID + "/permissions/" + attendanceCheckPermissionID).
					Expect(http.StatusOK).
					HasMessage("Added permission to role")
			})
		})

		t.Run("RegisterPermissionsForAdminRole", func(t *testing.T) {
			t.Run("RegisterActivityCreatePermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				authClient.POST("/projects/" + projectID + "/roles/" + adminRoleID + "/permissions/" + activityCreatePermissionID).
					Expect(http.StatusOK).
					HasMessage("Added permission to role")
			})
			t.Run("RegisterAdministratePermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				authClient.POST("/projects/" + projectID + "/roles/" + adminRoleID + "/permissions/" + administratePermissionID).
					Expect(http.StatusOK).
					HasMessage("Added permission to role")
			})
		})

		t.Run("RegisterPermissionsForOwnerRole", func(t *testing.T) {
			t.Run("RegisterEventCreatePermission", func(t *testing.T) {
				authClient := suite.NewClient(t).WithAuth(user.auth)
				authClient.POST("/projects/" + projectID + "/roles/" + ownerRoleID + "/permissions/" + eventCreatePermissionID).
					Expect(http.StatusOK).
					HasMessage("Added permission to role")
			})
		})
	})

	var participant, staff, admin, owner *Client
	t.Run("CreateUsers", func(t *testing.T) {
		t.Run("CreateParticipant", func(t *testing.T) {
			participant = suite.NewClient(t).WithCredentials("participant@mail.com", ValidPassword).
				ProjectRegister(projectID).
				ProjectLogin(projectID)
		})
		t.Run("CreateStaff", func(t *testing.T) {
			staff = suite.NewClient(t).WithCredentials("staff@mail.com", ValidPassword).
				ProjectRegister(projectID).
				ProjectLogin(projectID)
		})
		t.Run("CreateAdmin", func(t *testing.T) {
			admin = suite.NewClient(t).WithCredentials("admin@mail.com", ValidPassword).
				ProjectRegister(projectID).
				ProjectLogin(projectID)
		})
		t.Run("CreateOwner", func(t *testing.T) {
			owner = suite.NewClient(t).WithCredentials("owner@mail.com", ValidPassword).
				ProjectRegister(projectID).
				ProjectLogin(projectID)
		})
	})

	t.Run("AssignRolesAndPermissions", func(t *testing.T) {
		t.Run("AssignOwnerRole", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *owner.projectUserID + "/roles").
				WithBody(map[string]interface{}{
					"role_id":  ownerRoleID,
					"scope_id": eventsScopeID,
				}).
				Expect(http.StatusOK).
				HasMessage("Added role to user")
		})

		t.Run("AssignAdminRole", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *admin.projectUserID + "/roles").
				WithBody(map[string]interface{}{
					"role_id":  adminRoleID,
					"scope_id": event1ScopeID,
				}).
				Expect(http.StatusOK).
				HasMessage("Added role to user")

			// Give admin extra permissions
			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *admin.projectUserID + "/permissions").
				WithBody(map[string]interface{}{
					"permission_id": activityEditPermissionID,
					"scope_id":      event1ScopeID,
				}).
				Expect(http.StatusOK).
				HasMessage("Added permission to user")

			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *admin.projectUserID + "/permissions").
				WithBody(map[string]interface{}{
					"permission_id": activityDeletePermissionID,
					"scope_id":      event1ScopeID,
				}).
				Expect(http.StatusOK).
				HasMessage("Added permission to user")
		})

		t.Run("AssignStaffRole", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *staff.projectUserID + "/roles").
				WithBody(map[string]interface{}{
					"role_id":  staffRoleID,
					"scope_id": event1ScopeID,
				}).
				Expect(http.StatusOK).
				HasMessage("Added role to user")

			// Give staff extra permissions
			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *staff.projectUserID + "/permissions").
				WithBody(map[string]interface{}{
					"permission_id": dashboardAccessPermissionID,
					"scope_id":      eventsScopeID,
				}).
				Expect(http.StatusOK).
				HasMessage("Added permission to user")

			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *staff.projectUserID + "/permissions").
				WithBody(map[string]interface{}{
					"permission_id": attendanceMarkPermissionID,
					"scope_id":      activitiesScopeID,
				}).
				Expect(http.StatusOK).
				HasMessage("Added permission to user")
		})

		t.Run("AssignParticipantRole", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *participant.projectUserID + "/roles").
				WithBody(map[string]interface{}{
					"role_id":  participantRoleID,
					"scope_id": event1ScopeID,
				}).
				Expect(http.StatusOK).
				HasMessage("Added role to user")

			// Give participant ticket buy permission
			suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + *participant.projectUserID + "/permissions").
				WithBody(map[string]interface{}{
					"permission_id": ticketBuyPermissionID,
					"scope_id":      nil,
				}).
				Expect(http.StatusOK).
				HasMessage("Added permission to user")
		})
	})

	t.Run("VerifyEffectivePermissions", func(t *testing.T) {
		t.Run("VerifyParticipantEffectivePermissions", func(t *testing.T) {
			val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*participant.projectUserID+"/permissions").
				WithQuery("scope_id", event1ScopeID).
				Expect(http.StatusOK).
				RequireDataValue()

			// Ordered alphabetically by object, then action
			spec := []interface{}{
				map[string]interface{}(activityAttendPermission),
				map[string]interface{}(eventParticipatePermission),
				map[string]interface{}(productBuyPermission),
				map[string]interface{}(ticketBuyPermission),
			}

			ValidateExact(t, val, spec)
		})

		t.Run("VerifyStaffEffectivePermissions", func(t *testing.T) {
			val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*staff.projectUserID+"/permissions").
				WithQuery("scope_id", activitiesScopeID).
				Expect(http.StatusOK).
				RequireDataValue()

			// Staff gets permissions from role + directly assigned
			// Ordered alphabetically by object, then action
			spec := []interface{}{
				map[string]interface{}(attendanceCheckPermission),
				map[string]interface{}(attendanceMarkPermission),
				map[string]interface{}(dashboardAccessPermission),
			}

			Validate(t, val, spec)
		})

		t.Run("VerifyAdminEffectivePermissions", func(t *testing.T) {
			val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*admin.projectUserID+"/permissions").
				WithQuery("scope_id", event1ScopeID).
				Expect(http.StatusOK).
				RequireDataValue()

			// Admin gets permissions from role + directly assigned
			// Ordered alphabetically by object, then action
			spec := []interface{}{
				map[string]interface{}(activityCreatePermission),
				map[string]interface{}(activityDeletePermission),
				map[string]interface{}(activityEditPermission),
				map[string]interface{}(administratePermission),
			}

			ValidateExact(t, val, spec)
		})

		t.Run("VerifyOwnerEffectivePermissions", func(t *testing.T) {
			val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/identities/"+*owner.projectUserID+"/permissions").
				WithQuery("scope_id", eventsScopeID).
				Expect(http.StatusOK).
				RequireDataValue()

			spec := []interface{}{
				map[string]interface{}(eventCreatePermission),
			}

			ValidateExact(t, val, spec)
		})
	})

	t.Run("VerifyAuthorizationChecks", func(t *testing.T) {
		t.Run("AllowScenario_ParticipantAttendsActivity", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID,
					"entity_id":  *participant.projectUserID,
					"object":     "activity",
					"action":     "attend",
				}).
				Expect(http.StatusOK).
				HasMessage("Permission Granted")
		})

		t.Run("AllowScenario_ParticipantBuysTicket", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID,
					"entity_id":  *participant.projectUserID,
					"object":     "ticket",
					"action":     "buy",
				}).
				Expect(http.StatusOK).
				HasMessage("Permission Granted")
		})

		t.Run("AllowScenario_StaffChecksAttendance", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   activitiesScopeID,
					"entity_id":  *staff.projectUserID,
					"object":     "attendance",
					"action":     "check",
				}).
				Expect(http.StatusOK).
				HasMessage("Permission Granted")
		})

		t.Run("AllowScenario_StaffMarksAttendance", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   activitiesScopeID,
					"entity_id":  *staff.projectUserID,
					"object":     "attendance",
					"action":     "mark",
				}).
				Expect(http.StatusOK).
				HasMessage("Permission Granted")
		})

		t.Run("AllowScenario_AdminCreatesActivity", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID,
					"entity_id":  *admin.projectUserID,
					"object":     "activity",
					"action":     "create",
				}).
				Expect(http.StatusOK).
				HasMessage("Permission Granted")
		})

		t.Run("AllowScenario_AdminEditsActivity", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID,
					"entity_id":  *admin.projectUserID,
					"object":     "activity",
					"action":     "edit",
				}).
				Expect(http.StatusOK).
				HasMessage("Permission Granted")
		})

		t.Run("AllowScenario_OwnerCreatesEvent", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   eventsScopeID,
					"entity_id":  *owner.projectUserID,
					"object":     "event",
					"action":     "create",
				}).
				Expect(http.StatusOK).
				HasMessage("Permission Granted")
		})

		t.Run("DenyScenario_ParticipantCannotCreateActivity", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   event1ScopeID,
					"entity_id":  *participant.projectUserID,
					"object":     "activity",
					"action":     "create",
				}).
				Expect(http.StatusForbidden).
				HasErrID(errx.PERMissionInsufficient)
		})

		t.Run("DenyScenario_StaffCannotDeleteActivity", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   activitiesScopeID,
					"entity_id":  *staff.projectUserID,
					"object":     "activity",
					"action":     "delete",
				}).
				Expect(http.StatusForbidden).
				HasErrID(errx.PERMissionInsufficient)
		})

		t.Run("DenyScenario_AdminCannotCreateEvent", func(t *testing.T) {
			suite.NewClient(t).WithAuth(user.auth).POST("/authz/check").
				WithBody(map[string]interface{}{
					"project_id": projectID,
					"scope_id":   eventsScopeID,
					"entity_id":  *admin.projectUserID,
					"object":     "event",
					"action":     "create",
				}).
				Expect(http.StatusForbidden).
				HasErrID(errx.PERMissionInsufficient)
		})
	})
}
