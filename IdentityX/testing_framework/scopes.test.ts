import { beforeAll, describe, test } from "vitest";
import { createClient } from "./helpers/index.js";
import { assertErrID, assertMessage, shouldFail } from "./helpers/assert.js";
import {
    Validate,
    AnyUUID,
    AnyDate,
    AsString,
    Store,
} from "./helpers/validate.js";
import { dbQuery } from "./helpers/db.js";
import { ValidPassword } from "./fixtures/auth/testdata.js";

// SCOPEEmptyName        = fail.ID(0, "SCOPE", 0, true, ...) → 0_SCOPE_0000_S
// SCOPEDuplicateSibling = fail.ID(0, "SCOPE", 6, true, ...) → 0_SCOPE_0006_S
const ErrScopeEmptyName        = "0_SCOPE_0000_S";
const ErrScopeDuplicateSibling = "0_SCOPE_0006_S";

describe("Scopes", () => {
    let user: any;
    let projectID: string;
    let scopeID: string;

    beforeAll(async () => {
        await createClient().withCredentials("scopes@mail.com", ValidPassword).register();
        user = await createClient().withCredentials("scopes@mail.com", ValidPassword).login();
        const project = await user.post("/projects", { project_name: "ScopeProject", metadata: { env: "test" } });
        projectID = project.id;
    });

    test("CreateScope", async () => {
        const data = await user.post(`/projects/${projectID}/scopes`, {
            name: "events",
            external_id: "event1",
        });
        const ref = { current: "" };
        Validate(data, {
            id:          Store(ref, AnyUUID),
            project_id:  AsString(projectID, AnyUUID),
            name:        "events",
            external_id: "event1",
            type:        "project_scope",
            created_at:  AnyDate,
        });
        scopeID = ref.current as string;
    });

    test("GetScope", async () => {
        const data = await user.get(`/projects/${projectID}/scopes/${scopeID}`);
        Validate(data, {
            id:          AsString(scopeID, AnyUUID),
            project_id:  AsString(projectID, AnyUUID),
            name:        "events",
            external_id: "event1",
            type:        "project_scope",
            created_at:  AnyDate,
        });
    });

    test("CreateScopeNoName", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/scopes`, { external_id: "event1" }),
            400
        );
        assertErrID(body, ErrScopeEmptyName);
        assertMessage(body, "scope name cannot be empty");
    });

    test("CreateScopeExternalIDAlreadyInName", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/scopes`, { name: "events", external_id: "event1" }),
            409
        );
        assertErrID(body, ErrScopeDuplicateSibling);
        assertMessage(body, "scope with same name already exists under this parent");
    });

    test("CreateScopeExistingNameNoID", async () => {
        const data = await user.post(`/projects/${projectID}/scopes`, { name: "events" });
        Validate(data, {
            id:          AnyUUID,
            project_id:  AsString(projectID, AnyUUID),
            name:        "events",
            external_id: null,
            type:        "project_scope",
            created_at:  AnyDate,
        });
    });

    test("CreateScopeExistingNameAndNewID", async () => {
        const data = await user.post(`/projects/${projectID}/scopes`, {
            name: "events",
            external_id: "event2",
        });
        Validate(data, {
            id:          AnyUUID,
            project_id:  AsString(projectID, AnyUUID),
            name:        "events",
            external_id: "event2",
            type:        "project_scope",
            created_at:  AnyDate,
        });
    });

    test("GetAllProjectScopes", async () => {
        const data = await user.get(`/projects/${projectID}/scopes`);
        Validate(data, [
            { id: AnyUUID, project_id: AsString(projectID, AnyUUID), name: "events", external_id: "event1", type: "project_scope", created_at: AnyDate },
            { id: AnyUUID, project_id: AsString(projectID, AnyUUID), name: "events", external_id: null,     type: "project_scope", created_at: AnyDate },
            { id: AnyUUID, project_id: AsString(projectID, AnyUUID), name: "events", external_id: "event2", type: "project_scope", created_at: AnyDate },
        ]);
    });

    // CreateGlobalScopeError    — SKIPPED: tests raw DB constraint violations (unique/check),
    //                             not expressible via HTTP endpoints. See SKIPPED.md.
    // CreateProjectRootScopeError — SKIPPED: same reason.
    // CreateInvalidScopeType    — SKIPPED: same reason.

    test("CheckProjectRootScope", async () => {
        const result = await dbQuery<{ project_id: string; name: string | null; external_id: string | null; type: string }>(
            "SELECT project_id, name, external_id, type FROM scopes WHERE project_id = $1 AND type = 'project_root'",
            [projectID]
        );
        const rows = result.rows;
        if (rows.length !== 1) throw new Error(`expected 1 project_root scope, got ${rows.length}`);
        const root = rows[0];
        if (root.project_id !== projectID) throw new Error(`expected project_id ${projectID}, got ${root.project_id}`);
        if (root.name !== null) throw new Error(`expected name null, got ${root.name}`);
        if (root.external_id !== null) throw new Error(`expected external_id null, got ${root.external_id}`);
        if (root.type !== "project_root") throw new Error(`expected type project_root, got ${root.type}`);
    });
});