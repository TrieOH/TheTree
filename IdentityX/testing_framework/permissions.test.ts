import { beforeAll, describe, test } from "vitest";
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

const ErrPermissionAlreadyExists = "0_PERM_0008_D";
const ErrPermissionInvalidObject = "0_PERM_0007_D";
const ErrPermissionInvalidAction = "0_PERM_0006_D";
const ErrRequestValidation       = "0_REQ_0004_D";
const ErrRequestUnknownQueryParam = "0_REQ_0003_D";

describe("Permissions", () => {
    let user: any;
    let projectID: string;
    let permissionID: string;
    let anotherPermissionID: string;
    let createProductPermissionID: string;
    let editProductPermissionID: string;

    beforeAll(async () => {
        await createClient().withCredentials("permissions@mail.com", ValidPassword).register();
        user = await createClient().withCredentials("permissions@mail.com", ValidPassword).login();
        const project = await user.post("/projects", { project_name: "permissions testing", metadata: { env: "test" } });
        projectID = project.id;
    });

    test("CreatePermission", async () => {
        const ref = { current: "" };
        const data = await user.post(`/projects/${projectID}/permissions`, {
            object: "document",
            action: "create",
        });
        Validate(data, {
            id:         Store(ref, AnyUUID),
            project_id: AsString(projectID, AnyUUID),
            object:     "document",
            action:     "create",
            created_at: AnyDate,
        });
        permissionID = ref.current as string;
    });

    test("GetPermissionByID", async () => {
        const data = await user.get(`/projects/${projectID}/permissions/${permissionID}`);
        Validate(data, {
            id:         AsString(permissionID, AnyUUID),
            project_id: AsString(projectID, AnyUUID),
            object:     "document",
            action:     "create",
            created_at: AnyDate,
        });
    });

    test("CreateAnotherPermission", async () => {
        const ref = { current: "" };
        const data = await user.post(`/projects/${projectID}/permissions`, {
            object: "event",
            action: "read",
        });
        Validate(data, {
            id:         Store(ref, AnyUUID),
            project_id: AsString(projectID, AnyUUID),
            object:     "event",
            action:     "read",
            created_at: AnyDate,
        });
        anotherPermissionID = ref.current as string;
    });

    test("ListProjectPermissions", async () => {
        const data = await user.get(`/projects/${projectID}/permissions`);
        Validate(data, [
            { id: AsString(anotherPermissionID, AnyUUID), project_id: AsString(projectID, AnyUUID), object: "event",    action: "read",   created_at: AnyDate },
            { id: AsString(permissionID,        AnyUUID), project_id: AsString(projectID, AnyUUID), object: "document", action: "create", created_at: AnyDate },
        ]);
    });

    test("GetPermissionByObject", async () => {
        const data = await user.get(`/projects/${projectID}/permissions?object=event`);
        Validate(data, [
            { id: AsString(anotherPermissionID, AnyUUID), project_id: AsString(projectID, AnyUUID), object: "event", action: "read", created_at: AnyDate },
        ]);
    });

    test("GetPermissionByAction", async () => {
        const data = await user.get(`/projects/${projectID}/permissions?action=create`);
        Validate(data, [
            { id: AsString(permissionID, AnyUUID), project_id: AsString(projectID, AnyUUID), object: "document", action: "create", created_at: AnyDate },
        ]);
    });

    test("CreateCreateProductPermission", async () => {
        const ref = { current: "" };
        const data = await user.post(`/projects/${projectID}/permissions`, {
            object: "product",
            action: "create",
        });
        Validate(data, {
            id:         Store(ref, AnyUUID),
            project_id: AsString(projectID, AnyUUID),
            object:     "product",
            action:     "create",
            created_at: AnyDate,
        });
        createProductPermissionID = ref.current as string;
    });

    test("GetPermissionByActionAgain", async () => {
        const data = await user.get(`/projects/${projectID}/permissions?action=create`);
        Validate(data, [
            { id: AsString(createProductPermissionID, AnyUUID), project_id: AsString(projectID, AnyUUID), object: "product",  action: "create", created_at: AnyDate },
            { id: AsString(permissionID,              AnyUUID), project_id: AsString(projectID, AnyUUID), object: "document", action: "create", created_at: AnyDate },
        ]);
    });

    test("CreateEditProductPermission", async () => {
        const ref = { current: "" };
        const data = await user.post(`/projects/${projectID}/permissions`, {
            object: "product",
            action: "edit",
        });
        Validate(data, {
            id:         Store(ref, AnyUUID),
            project_id: AsString(projectID, AnyUUID),
            object:     "product",
            action:     "edit",
            created_at: AnyDate,
        });
        editProductPermissionID = ref.current as string;
    });

    test("GetPermissionByObjectAgain", async () => {
        const data = await user.get(`/projects/${projectID}/permissions?object=product`);
        Validate(data, [
            { id: AsString(editProductPermissionID,   AnyUUID), project_id: AsString(projectID, AnyUUID), object: "product", action: "edit",   created_at: AnyDate },
            { id: AsString(createProductPermissionID, AnyUUID), project_id: AsString(projectID, AnyUUID), object: "product", action: "create", created_at: AnyDate },
        ]);
    });

    test("GetPermissionByObjectAndAction", async () => {
        const data = await user.get(`/projects/${projectID}/permissions?object=product&action=edit`);
        Validate(data, [
            { id: AsString(editProductPermissionID, AnyUUID), project_id: AsString(projectID, AnyUUID), object: "product", action: "edit", created_at: AnyDate },
        ]);
    });

    test("CreateAdminPermission", async () => {
        const data = await user.post(`/projects/${projectID}/permissions`, {
            object: "user",
            action: "delete",
        });
        Validate(data, {
            id:         AnyUUID,
            project_id: AsString(projectID, AnyUUID),
            object:     "user",
            action:     "delete",
            created_at: AnyDate,
        });
    });

    test("CreateMasterPermission", async () => {
        const data = await user.post(`/projects/${projectID}/permissions`, {
            object: "system",
            action: "admin",
        });
        Validate(data, {
            id:         AnyUUID,
            project_id: AsString(projectID, AnyUUID),
            object:     "system",
            action:     "admin",
            created_at: AnyDate,
        });
    });

    test("CreateDuplicatePermission", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/permissions`, { object: "document", action: "create" }),
            409
        );
        assertErrID(body, ErrPermissionAlreadyExists);
        assertMessage(body, "permission with object(document) and action(create) already exists");
    });

    test("CreatePermissionNoAction", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/permissions`, { object: "document", action: "" }),
            400
        );
        assertErrID(body, ErrRequestValidation);
        assertMessage(body, "Validation failed");
        assertValidationError(body, "action is required");
    });

    test("CreatePermissionNoObject", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/permissions`, { object: "", action: "create" }),
            400
        );
        assertErrID(body, ErrRequestValidation);
        assertMessage(body, "Validation failed");
        assertValidationError(body, "object is required");
    });

    test("CreateEmptyPermission", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/permissions`, { object: "", action: "" }),
            400
        );
        assertErrID(body, ErrRequestValidation);
        assertMessage(body, "Validation failed");
        assertValidationError(body, "object is required", "action is required");
    });

    test("CreateInvalidObjectPermission", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/permissions`, { object: "event:*", action: "read" }),
            400
        );
        assertErrID(body, ErrPermissionInvalidObject);
        assertMessage(body, "invalid permission object: (event:*)");
    });

    test("CreateInvalidActionPermission", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/permissions`, { object: "document", action: "attendance:mark" }),
            400
        );
        assertErrID(body, ErrPermissionInvalidAction);
        assertMessage(body, "invalid permission action: (attendance:mark)");
    });

    test("CreateWildcardObjectPermission", async () => {
        const data = await user.post(`/projects/${projectID}/permissions`, { object: "*", action: "read" });
        Validate(data, {
            id:         AnyUUID,
            project_id: AsString(projectID, AnyUUID),
            object:     "*",
            action:     "read",
            created_at: AnyDate,
        });
    });

    test("CreateWildcardActionPermission", async () => {
        const data = await user.post(`/projects/${projectID}/permissions`, { object: "document", action: "*" });
        Validate(data, {
            id:         AnyUUID,
            project_id: AsString(projectID, AnyUUID),
            object:     "document",
            action:     "*",
            created_at: AnyDate,
        });
    });

    test("CreateFullWildcardPermission", async () => {
        const data = await user.post(`/projects/${projectID}/permissions`, { object: "*", action: "*" });
        Validate(data, {
            id:         AnyUUID,
            project_id: AsString(projectID, AnyUUID),
            object:     "*",
            action:     "*",
            created_at: AnyDate,
        });
    });

    test("NotAllowedQueryParam", async () => {
        const body = await shouldFail(
            user.get(`/projects/${projectID}/permissions?not-allowed=should_deny`),
            400
        );
        assertErrID(body, ErrRequestUnknownQueryParam);
        assertMessage(body, "unknown query parameter: not-allowed");
    });
});