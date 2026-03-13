import { beforeAll, describe, test } from "vitest";
import { expect } from "vitest";
import { createClient } from "./helpers/index.js";
import { assertErrID, assertMessage, assertValidationError, shouldFail } from "./helpers/assert.js";
import {
    Validate,
    AnyUUID,
    AnyDate,
    AsString,
    Store,
} from "./helpers/validate.js";
import { ValidPassword } from "./fixtures/auth/testdata.js";

const ErrRoleNameAlreadyTaken    = "0_ROLE_0002_D";
const ErrRoleAlreadyGranted      = "0_ROLE_0001_D";
const ErrPermissionAlreadyGranted = "0_PERM_0004_D";
const ErrRequestValidation        = "0_REQ_0004_D";
const ErrRequestUnknownQueryParam = "0_REQ_0003_D";

describe("Roles", () => {
    let user: any;
    let projectID: string;

    let adminRoleID: string;
    let adminCreateDate: string;
    let adminUpdate1Date: string;
    let adminUpdate2Date: string;
    let userRoleID: string;
    let scopelessRoleID: string;

    let createEventPermissionID: string;
    let markAttendancePermissionID: string;
    let attendActivityPermissionID: string;

    let eventScopeID: string;
    let activityScopeID: string;
    let activitySubScopeID: string;

    let projectUserID: string;

    beforeAll(async () => {
        await createClient().withCredentials("roles@mail.com", ValidPassword).register();
        user = await createClient().withCredentials("roles@mail.com", ValidPassword).login();
        const project = await user.post("/projects", { project_name: "roles testing", metadata: { env: "test" } });
        projectID = project.id;

        // Register and login a project user, then fetch their identity ID
        await createClient().withCredentials("roles-user@mail.com", ValidPassword).projectRegister(projectID);
        const projectUser = await createClient().withCredentials("roles-user@mail.com", ValidPassword).projectLogin(projectID);
        const me = await projectUser.get("/sessions/me");
        projectUserID = me.access.sub.id;
    });

    // ── Role CRUD ──────────────────────────────────────────────────────────────

    test("CreateRole", async () => {
        const idRef  = { current: "" };
        const dateRef = { current: "" };
        const data = await user.post(`/projects/${projectID}/roles`, {
            name: "admin",
            description: "can do stuff",
        });
        Validate(data, {
            id:          Store(idRef, AnyUUID),
            project_id:  AsString(projectID, AnyUUID),
            name:        "admin",
            description: "can do stuff",
            created_at:  AnyDate,
            updated_at:  Store(dateRef, AnyDate),
        });
        adminRoleID    = idRef.current as string;
        adminCreateDate = dateRef.current as string;
    });

    test("UpdateRoleDescription", async () => {
        await user.patch(`/projects/${projectID}/roles/${adminRoleID}/description`, {
            description: "can do stuff and more stuff",
        });
    });

    test("GetRoleByID", async () => {
        const dateRef = { current: "" };
        const data = await user.get(`/projects/${projectID}/roles/${adminRoleID}`);
        Validate(data, {
            id:          AsString(adminRoleID, AnyUUID),
            project_id:  AsString(projectID, AnyUUID),
            name:        "admin",
            description: "can do stuff and more stuff",
            created_at:  AnyDate,
            updated_at:  Store(dateRef, AnyDate),
        });
        adminUpdate1Date = dateRef.current as string;
        expect(adminCreateDate).not.toBe(adminUpdate1Date);
    });

    test("UpdateRoleDescriptionAgain", async () => {
        await user.patch(`/projects/${projectID}/roles/${adminRoleID}/description`, {
            description: "can do stuff and more stuff but not that",
        });
    });

    test("GetRoleByName", async () => {
        const dateRef = { current: "" };
        const data = await user.get(`/projects/${projectID}/roles/search?name=admin`);
        Validate(data, {
            id:          AsString(adminRoleID, AnyUUID),
            project_id:  AsString(projectID, AnyUUID),
            name:        "admin",
            description: "can do stuff and more stuff but not that",
            created_at:  AnyDate,
            updated_at:  Store(dateRef, AnyDate),
        });
        adminUpdate2Date = dateRef.current as string;
        expect(adminUpdate1Date).not.toBe(adminUpdate2Date);
    });

    test("CreateDuplicateRole", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/roles`, { name: "admin" }),
            409
        );
        assertErrID(body, ErrRoleNameAlreadyTaken);
        assertMessage(body, "role name already taken");
    });

    test("CreateUserRole", async () => {
        const ref = { current: "" };
        const data = await user.post(`/projects/${projectID}/roles`, { name: "user" });
        Validate(data, {
            id:          Store(ref, AnyUUID),
            project_id:  AsString(projectID, AnyUUID),
            name:        "user",
            description: null,
            created_at:  AnyDate,
            updated_at:  AnyDate,
        });
        userRoleID = ref.current as string;
    });

    test("ListProjectRoles", async () => {
        const data = await user.get(`/projects/${projectID}/roles`);
        Validate(data, [
            {
                id:          AsString(userRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "user",
                description: null,
                created_at:  AnyDate,
                updated_at:  AnyDate,
            },
            {
                id:          AsString(adminRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "admin",
                description: "can do stuff and more stuff but not that",
                created_at:  AnyDate,
                updated_at:  AsString(adminUpdate2Date, AnyDate),
            },
        ]);
    });

    test("CreateRoleNoName", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/roles`, { name: "" }),
            400
        );
        assertErrID(body, ErrRequestValidation);
        assertMessage(body, "Validation failed");
        assertValidationError(body, "name is required");
    });

    test("ForbiddenQueryParam", async () => {
        const body = await shouldFail(
            user.get(`/projects/${projectID}/roles/search?something_else=should_fail`),
            400
        );
        assertErrID(body, ErrRequestUnknownQueryParam);
        assertMessage(body, "unknown query parameter: something_else");
    });

    // ── Permissions ────────────────────────────────────────────────────────────

    test("CreateEventPermission", async () => {
        const ref = { current: "" };
        const data = await user.post(`/projects/${projectID}/permissions`, {
            object: "event",
            action: "create",
        });
        Validate(data, {
            id:         Store(ref, AnyUUID),
            project_id: AsString(projectID, AnyUUID),
            object:     "event",
            action:     "create",
            created_at: AnyDate,
        });
        createEventPermissionID = ref.current as string;
    });

    test("CreateAttendanceMarkPermission", async () => {
        const ref = { current: "" };
        const data = await user.post(`/projects/${projectID}/permissions`, {
            object: "attendance",
            action: "mark",
        });
        Validate(data, {
            id:         Store(ref, AnyUUID),
            project_id: AsString(projectID, AnyUUID),
            object:     "attendance",
            action:     "mark",
            created_at: AnyDate,
        });
        markAttendancePermissionID = ref.current as string;
    });

    test("CreateActivityAttendancePermission", async () => {
        const ref = { current: "" };
        const data = await user.post(`/projects/${projectID}/permissions`, {
            object: "activity",
            action: "attend",
        });
        Validate(data, {
            id:         Store(ref, AnyUUID),
            project_id: AsString(projectID, AnyUUID),
            object:     "activity",
            action:     "attend",
            created_at: AnyDate,
        });
        attendActivityPermissionID = ref.current as string;
    });

    test("AddAdminPermissions", async () => {
        const r1 = await user.post(`/projects/${projectID}/roles/${adminRoleID}/permissions/${createEventPermissionID}`, {});
        assertMessage(r1, "Added permission to role");
        const r2 = await user.post(`/projects/${projectID}/roles/${adminRoleID}/permissions/${markAttendancePermissionID}`, {});
        assertMessage(r2, "Added permission to role");
    });

    test("AddUserPermissions", async () => {
        const r = await user.post(`/projects/${projectID}/roles/${userRoleID}/permissions/${attendActivityPermissionID}`, {});
        assertMessage(r, "Added permission to role");
    });

    test("GetAdminPermissions", async () => {
        const data = await user.get(`/projects/${projectID}/roles/${adminRoleID}/permissions`);
        Validate(data, [
            { id: AsString(markAttendancePermissionID, AnyUUID), project_id: AsString(projectID, AnyUUID), object: "attendance", action: "mark",   created_at: AnyDate },
            { id: AsString(createEventPermissionID,    AnyUUID), project_id: AsString(projectID, AnyUUID), object: "event",      action: "create", created_at: AnyDate },
        ]);
    });

    test("GetUserPermissions", async () => {
        const data = await user.get(`/projects/${projectID}/roles/${userRoleID}/permissions`);
        Validate(data, [
            { id: AsString(attendActivityPermissionID, AnyUUID), project_id: AsString(projectID, AnyUUID), object: "activity", action: "attend", created_at: AnyDate },
        ]);
    });

    test("RemoveAdminPermission", async () => {
        const r = await user.del(`/projects/${projectID}/roles/${adminRoleID}/permissions/${createEventPermissionID}`);
        assertMessage(r, "Removed permission from role");
    });

    test("GetAdminPermissionsAgain", async () => {
        const data = await user.get(`/projects/${projectID}/roles/${adminRoleID}/permissions`);
        Validate(data, [
            { id: AsString(markAttendancePermissionID, AnyUUID), project_id: AsString(projectID, AnyUUID), object: "attendance", action: "mark", created_at: AnyDate },
        ]);
    });

    // ── Scopes ─────────────────────────────────────────────────────────────────

    test("CreateEventScope", async () => {
        const ref = { current: "" };
        const data = await user.post(`/projects/${projectID}/scopes`, {
            name:        "events",
            external_id: null,
        });
        Validate(data, {
            id:          Store(ref, AnyUUID),
            project_id:  AsString(projectID, AnyUUID),
            name:        "events",
            external_id: null,
            type:        "project_scope",
            created_at:  AnyDate,
        });
        eventScopeID = ref.current as string;
    });

    // ── Identity role assignment ───────────────────────────────────────────────

    test("GiveUserAdminRole", async () => {
        const r = await user.post(`/projects/${projectID}/identities/${projectUserID}/roles`, {
            role_id:  adminRoleID,
            scope_id: eventScopeID,
        });
        assertMessage(r, "Added role to user");
    });

    test("GetUserRoles", async () => {
        const data = await user.get(`/projects/${projectID}/identities/${projectUserID}/roles`);
        Validate(data, [
            {
                id:          AsString(adminRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "admin",
                description: "can do stuff and more stuff but not that",
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    AsString(eventScopeID, AnyUUID),
                scope_name:  "events",
                external_id: null,
            },
        ]);
    });

    test("GiveUserUserRole", async () => {
        const r = await user.post(`/projects/${projectID}/identities/${projectUserID}/roles`, {
            role_id:  userRoleID,
            scope_id: eventScopeID,
        });
        assertMessage(r, "Added role to user");
    });

    test("GetUserRolesAgain", async () => {
        const data = await user.get(`/projects/${projectID}/identities/${projectUserID}/roles`);
        Validate(data, [
            {
                id:          AsString(adminRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "admin",
                description: "can do stuff and more stuff but not that",
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    AsString(eventScopeID, AnyUUID),
                scope_name:  "events",
                external_id: null,
            },
            {
                id:          AsString(userRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "user",
                description: null,
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    AsString(eventScopeID, AnyUUID),
                scope_name:  "events",
                external_id: null,
            },
        ]);
    });

    test("CreateScopelessRole", async () => {
        const ref = { current: "" };
        const data = await user.post(`/projects/${projectID}/roles`, {
            name:        "scopeless",
            description: "this role should be project wide",
        });
        Validate(data, {
            id:          Store(ref, AnyUUID),
            project_id:  AsString(projectID, AnyUUID),
            name:        "scopeless",
            description: "this role should be project wide",
            created_at:  AnyDate,
            updated_at:  AnyDate,
        });
        scopelessRoleID = ref.current as string;
    });

    test("GiveUserScopelessRole", async () => {
        const r = await user.post(`/projects/${projectID}/identities/${projectUserID}/roles`, {
            role_id:  scopelessRoleID,
            scope_id: null,
        });
        assertMessage(r, "Added role to user");
    });

    test("GetUserRolesAfterScopeless", async () => {
        const data = await user.get(`/projects/${projectID}/identities/${projectUserID}/roles`);
        Validate(data, [
            {
                id:          AsString(adminRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "admin",
                description: "can do stuff and more stuff but not that",
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    AsString(eventScopeID, AnyUUID),
                scope_name:  "events",
                external_id: null,
            },
            {
                id:          AsString(scopelessRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "scopeless",
                description: "this role should be project wide",
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    null,
                scope_name:  null,
                external_id: null,
            },
            {
                id:          AsString(userRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "user",
                description: null,
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    AsString(eventScopeID, AnyUUID),
                scope_name:  "events",
                external_id: null,
            },
        ]);
    });

    test("TakeUserAdminRole", async () => {
        const r = await user.del(`/projects/${projectID}/identities/${projectUserID}/roles`, {
            role_id:  adminRoleID,
            scope_id: eventScopeID,
        });
        assertMessage(r, "Removed role from user");
    });

    test("GetUserRolesAfterTake", async () => {
        const data = await user.get(`/projects/${projectID}/identities/${projectUserID}/roles`);
        Validate(data, [
            {
                id:          AsString(scopelessRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "scopeless",
                description: "this role should be project wide",
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    null,
                scope_name:  null,
                external_id: null,
            },
            {
                id:          AsString(userRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "user",
                description: null,
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    AsString(eventScopeID, AnyUUID),
                scope_name:  "events",
                external_id: null,
            },
        ]);
    });

    test("GiveUserScopelessRoleOnAScope", async () => {
        const r = await user.post(`/projects/${projectID}/identities/${projectUserID}/roles`, {
            role_id:  scopelessRoleID,
            scope_id: eventScopeID,
        });
        assertMessage(r, "Added role to user");
    });

    test("GetUserRolesAfterScopelessOnScope", async () => {
        const data = await user.get(`/projects/${projectID}/identities/${projectUserID}/roles`);
        Validate(data, [
            {
                id:          AsString(scopelessRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "scopeless",
                description: "this role should be project wide",
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    null,
                scope_name:  null,
                external_id: null,
            },
            {
                id:          AsString(scopelessRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "scopeless",
                description: "this role should be project wide",
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    eventScopeID,
                scope_name:  "events",
                external_id: null,
            },
            {
                id:          AsString(userRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "user",
                description: null,
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    AsString(eventScopeID, AnyUUID),
                scope_name:  "events",
                external_id: null,
            },
        ]);
    });

    test("CreateActivityScope", async () => {
        const ref = { current: "" };
        const data = await user.post(`/projects/${projectID}/scopes`, { name: "activities" });
        Validate(data, {
            id:          Store(ref, AnyUUID),
            project_id:  AsString(projectID, AnyUUID),
            name:        "activities",
            external_id: null,
            type:        "project_scope",
            created_at:  AnyDate,
        });
        activityScopeID = ref.current as string;
    });

    test("GiveUserScopelessRoleOnActivityScope", async () => {
        const r = await user.post(`/projects/${projectID}/identities/${projectUserID}/roles`, {
            role_id:  scopelessRoleID,
            scope_id: activityScopeID,
        });
        assertMessage(r, "Added role to user");
    });

    test("GetUserRolesAfterScopelessOnActivityScope", async () => {
        const data = await user.get(`/projects/${projectID}/identities/${projectUserID}/roles`);
        Validate(data, [
            {
                id:          AsString(scopelessRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "scopeless",
                description: "this role should be project wide",
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    null,
                scope_name:  null,
                external_id: null,
            },
            {
                id:          AsString(scopelessRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "scopeless",
                description: "this role should be project wide",
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    activityScopeID,
                scope_name:  "activities",
                external_id: null,
            },
            {
                id:          AsString(scopelessRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "scopeless",
                description: "this role should be project wide",
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    eventScopeID,
                scope_name:  "events",
                external_id: null,
            },
            {
                id:          AsString(userRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "user",
                description: null,
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    AsString(eventScopeID, AnyUUID),
                scope_name:  "events",
                external_id: null,
            },
        ]);
    });

    test("CreateActivitySubScope", async () => {
        const ref = { current: "" };
        const data = await user.post(`/projects/${projectID}/scopes`, {
            name:        "activities",
            external_id: "123",
        });
        Validate(data, {
            id:          Store(ref, AnyUUID),
            project_id:  AsString(projectID, AnyUUID),
            name:        "activities",
            external_id: "123",
            type:        "project_scope",
            created_at:  AnyDate,
        });
        activitySubScopeID = ref.current as string;
    });

    test("GiveUserScopelessRoleOnActivitySubScope", async () => {
        const r = await user.post(`/projects/${projectID}/identities/${projectUserID}/roles`, {
            role_id:  scopelessRoleID,
            scope_id: activitySubScopeID,
        });
        assertMessage(r, "Added role to user");
    });

    test("GetUserRolesAfterScopelessOnActivitySubScope", async () => {
        const data = await user.get(`/projects/${projectID}/identities/${projectUserID}/roles`);
        Validate(data, [
            {
                id:          AsString(scopelessRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "scopeless",
                description: "this role should be project wide",
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    null,
                scope_name:  null,
                external_id: null,
            },
            {
                id:          AsString(scopelessRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "scopeless",
                description: "this role should be project wide",
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    activityScopeID,
                scope_name:  "activities",
                external_id: null,
            },
            {
                id:          AsString(scopelessRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "scopeless",
                description: "this role should be project wide",
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    activitySubScopeID,
                scope_name:  "activities",
                external_id: "123",
            },
            {
                id:          AsString(scopelessRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "scopeless",
                description: "this role should be project wide",
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    eventScopeID,
                scope_name:  "events",
                external_id: null,
            },
            {
                id:          AsString(userRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "user",
                description: null,
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    AsString(eventScopeID, AnyUUID),
                scope_name:  "events",
                external_id: null,
            },
        ]);
    });

    test("GiveUserDuplicateScopelessRole", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/identities/${projectUserID}/roles`, {
                role_id:  scopelessRoleID,
                scope_id: null,
            }),
            400
        );
        assertErrID(body, ErrRoleAlreadyGranted);
        assertMessage(body, "scopeless role already granted to user");
    });

    test("GiveUserDuplicateScopelessRoleOnEventScope", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/identities/${projectUserID}/roles`, {
                role_id:  scopelessRoleID,
                scope_id: eventScopeID,
            }),
            400
        );
        assertErrID(body, ErrRoleAlreadyGranted);
        assertMessage(body, "scopeless role already granted to user");
    });

    test("GiveUserDuplicateScopelessRoleOnActivityScope", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/identities/${projectUserID}/roles`, {
                role_id:  scopelessRoleID,
                scope_id: activityScopeID,
            }),
            400
        );
        assertErrID(body, ErrRoleAlreadyGranted);
        assertMessage(body, "scopeless role already granted to user");
    });

    test("GiveUserDuplicateScopelessRoleOnActivitySubScope", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/identities/${projectUserID}/roles`, {
                role_id:  scopelessRoleID,
                scope_id: activitySubScopeID,
            }),
            400
        );
        assertErrID(body, ErrRoleAlreadyGranted);
        assertMessage(body, "scopeless role already granted to user");
    });

    // FIXME: DELETE currently returns 200 even if the row didn't exist (no :execrows check yet)
    test("TakeRoleAlreadyTaken", async () => {
        const r = await user.del(`/projects/${projectID}/identities/${projectUserID}/roles`, {
            role_id:  adminRoleID,
            scope_id: eventScopeID,
        });
        assertMessage(r, "Removed role from user");
    });

    // ── By-name endpoints ──────────────────────────────────────────────────────

    test("GiveUserAdminRoleByName", async () => {
        const r = await user.post(`/projects/${projectID}/identities/${projectUserID}/roles/by-name`, {
            role_name: "admin",
            scope_id:  eventScopeID,
        });
        assertMessage(r, "Added role to user");
    });

    test("GetUserRolesByName", async () => {
        const data = await user.get(`/projects/${projectID}/identities/${projectUserID}/roles`);
        Validate(data, [
            {
                id:          AsString(adminRoleID, AnyUUID),
                project_id:  AsString(projectID, AnyUUID),
                name:        "admin",
                description: "can do stuff and more stuff but not that",
                created_at:  AnyDate,
                updated_at:  AnyDate,
                scope_id:    AsString(eventScopeID, AnyUUID),
                scope_name:  "events",
                external_id: null,
            },
        ]);
    });

    test("GiveUserUserRoleByName", async () => {
        const r = await user.post(`/projects/${projectID}/identities/${projectUserID}/roles/by-name`, {
            role_name: "user",
            scope_id:  eventScopeID,
        });
        assertMessage(r, "Added role to user");
    });

    test("GiveUserRoleByNameNotFound", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/identities/${projectUserID}/roles/by-name`, {
                role_name: "nonexistent_role",
                scope_id:  null,
            }),
            404
        );
        assertMessage(body, "role not found by name: nonexistent_role");
    });

    test("GiveUserRoleByNameEmpty", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/identities/${projectUserID}/roles/by-name`, {
                role_name: "",
                scope_id:  null,
            }),
            400
        );
        assertMessage(body, "invalid role name:");
    });

    test("GiveUserScopelessRoleByName", async () => {
        const r = await user.post(`/projects/${projectID}/identities/${projectUserID}/roles/by-name`, {
            role_name: "scopeless",
            scope_id:  null,
        });
        assertMessage(r, "Added role to user");
    });

    test("TakeUserAdminRoleByName", async () => {
        const r = await user.del(`/projects/${projectID}/identities/${projectUserID}/roles/by-name`, {
            role_name: "admin",
            scope_id:  eventScopeID,
        });
        assertMessage(r, "Removed role from user");
    });

    test("TakeUserRoleByNameNotFound", async () => {
        const body = await shouldFail(
            user.del(`/projects/${projectID}/identities/${projectUserID}/roles/by-name`, {
                role_name: "nonexistent_role",
                scope_id:  null,
            }),
            404
        );
        assertMessage(body, "role not found by name: nonexistent_role");
    });

    test("TakeUserRoleByNameEmpty", async () => {
        const body = await shouldFail(
            user.del(`/projects/${projectID}/identities/${projectUserID}/roles/by-name`, {
                role_name: "",
                scope_id:  null,
            }),
            400
        );
        assertMessage(body, "invalid role name:");
    });

    // ── Direct permission assignment ───────────────────────────────────────────

    test("GiveUserDirectEventPermissionNoScope", async () => {
        const r = await user.post(`/projects/${projectID}/identities/${projectUserID}/permissions`, {
            permission_id: createEventPermissionID,
            scope_id:      null,
        });
        assertMessage(r, "Added permission to user");
    });

    test("GiveUserDirectEventPermissionWithEventScope", async () => {
        const r = await user.post(`/projects/${projectID}/identities/${projectUserID}/permissions`, {
            permission_id: createEventPermissionID,
            scope_id:      eventScopeID,
        });
        assertMessage(r, "Added permission to user");
    });

    test("GiveUserDirectDuplicateEventPermissionNoScope", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/identities/${projectUserID}/permissions`, {
                permission_id: createEventPermissionID,
                scope_id:      null,
            }),
            409
        );
        assertErrID(body, ErrPermissionAlreadyGranted);
        assertMessage(body, "user already has this permission in the specified scope");
    });

    test("TakeDirectEventPermissionWithEventScopeFromUser", async () => {
        const r = await user.del(`/projects/${projectID}/identities/${projectUserID}/permissions`, {
            permission_id: createEventPermissionID,
            scope_id:      eventScopeID,
        });
        assertMessage(r, "Removed permission from user");
    });

    // FIXME: DELETE currently returns 200 even if the row didn't exist (no :execrows check yet)
    test("TakeDirectEventPermissionWithEventScopeFromUserAlreadyTaken", async () => {
        const r = await user.del(`/projects/${projectID}/identities/${projectUserID}/permissions`, {
            permission_id: createEventPermissionID,
            scope_id:      eventScopeID,
        });
        assertMessage(r, "Removed permission from user");
    });
});