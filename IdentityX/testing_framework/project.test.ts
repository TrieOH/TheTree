import { beforeAll, describe, expect, test } from "vitest";
import { createClient } from "./helpers/index.js";
import { assertErrID, assertMessage, assertValidationError, shouldFail } from "./helpers/assert.js";
import { Validate, AnyUUID, AnyDate, AsString } from "./helpers/validate.js";
import { ValidPassword } from "./fixtures/auth/testdata.js";

// Error IDs
// RequestValidationError = fail.ID(0, "REQ", 4, false, ...) → 0_REQ_0004_D
// SQLNotFound            = fail.ID(0, "SQL", 0, false, ...) → 0_SQL_0000_D
// ProjectNotFound        = fail.ID(0, "PROJECT", 0, true, ...) → 0_PROJECT_0000_S
const ErrRequestValidation = "0_REQ_0004_D";
const ErrSQLNotFound       = "0_SQL_0000_D";
const ErrProjectNotFound   = "0_PROJECT_0000_S";

// ============================================================================
// Ed25519 JWK verification
// ============================================================================

function verifyEd25519JWKKey(x: string): void {
    // x is base64url-encoded raw public key bytes (32 bytes for Ed25519)
    const base64 = x.replace(/-/g, "+").replace(/_/g, "/");
    const padded = base64.padEnd(base64.length + (4 - (base64.length % 4)) % 4, "=");
    const bytes  = Buffer.from(padded, "base64");

    expect(bytes.length).toBe(32);
}

// ============================================================================
// PROJECTS TESTS
// ============================================================================

describe("Projects", () => {
    let user: Awaited<ReturnType<typeof createClient.prototype.login>>;
    let projectID: string;

    beforeAll(async () => {
        await createClient().withCredentials("projects@mail.com", ValidPassword).register();
        user = await createClient().withCredentials("projects@mail.com", ValidPassword).login();
    });

    test("CreateProject", async () => {
        const data = await user.post("/projects", {
            project_name: "Test Project",
            metadata: { env: "test" },
        });

        Validate(data, {
            id:           AnyUUID,
            owner_id:     AnyUUID,
            project_name: "Test Project",
            is_active:    true,
            metadata:     { env: "test" },
            created_at:   AnyDate,
            updated_at:   AnyDate,
        });

        projectID = data.id;
    });

    test("ValidationCreateProject", async () => {
        const body = await shouldFail(
            user.post("/projects", { project_name: "", metadata: { env: "test" } }),
            400
        );
        assertErrID(body, ErrRequestValidation);
        assertValidationError(body, "project_name is required");
    });

    test("ListProjects", async () => {
        const data = await user.get("/projects");

        Validate(data, [
            {
                id:           AsString(projectID, AnyUUID),
                owner_id:     AnyUUID,
                project_name: "Test Project",
                is_active:    true,
                metadata:     { env: "test" },
                created_at:   AnyDate,
                updated_at:   AnyDate,
            },
        ]);
    });

    test("GetProject", async () => {
        const data = await user.get(`/projects/${projectID}`);

        Validate(data, {
            id:           AsString(projectID, AnyUUID),
            owner_id:     AnyUUID,
            project_name: "Test Project",
            is_active:    true,
            metadata:     { env: "test" },
            created_at:   AnyDate,
            updated_at:   AnyDate,
        });
    });

    test("ListProjectUsers", async () => {
        const projectUserEmail = "projectuser@mail.com";

        await createClient().http.post(`/projects/${projectID}/register`, {
            email:    projectUserEmail,
            password: ValidPassword,
        });

        const data = await user.get(`/projects/${projectID}/users`);

        Validate(data, [
            {
                id:         AnyUUID,
                project_id: AsString(projectID, AnyUUID),
                email:      projectUserEmail,
                user_type:  "project",
                is_active:  true,
                created_at: AnyDate,
                updated_at: AnyDate,
            },
        ]);
    });

    test("CrossUserAccess", async () => {
        await createClient().withCredentials("attacker@mail.com", ValidPassword).register();
        const attacker = await createClient().withCredentials("attacker@mail.com", ValidPassword).login();

        // Try to GET
        const getBody = await shouldFail(attacker.get(`/projects/${projectID}`), 404);
        assertErrID(getBody, ErrSQLNotFound);
        assertMessage(getBody, "project not found");

        // Try to PATCH
        const patchBody = await shouldFail(
            attacker.patch(`/projects/${projectID}`, { project_name: "Hacked" }),
            404
        );
        assertErrID(patchBody, ErrSQLNotFound);
        assertMessage(patchBody, "project not found");

        // Try to DELETE
        const deleteBody = await shouldFail(attacker.del(`/projects/${projectID}`), 404);
        assertErrID(deleteBody, ErrProjectNotFound);
        assertMessage(deleteBody, "project not found");

        // Ensure owner can still access it
        const data = await user.get(`/projects/${projectID}`);
        expect(data.id).toBe(projectID);
    });

    test("UpdateProject", async () => {
        const data = await user.patch(`/projects/${projectID}`, {
            project_name: "Updated Project",
            metadata:     { env: "prod" },
        });

        Validate(data, {
            id:           AsString(projectID, AnyUUID),
            owner_id:     AnyUUID,
            project_name: "Updated Project",
            is_active:    true,
            metadata:     { env: "prod" },
            created_at:   AnyDate,
            updated_at:   AnyDate,
        });
    });

    test("GetProjectJWKS", async () => {
        const res = await user.http.get(`/projects/${projectID}/.well-known/jwks.json`, {
            headers: { Cookie: `access_token=${user.auth!.accessToken}; refresh_token=${user.auth!.refreshToken}` },
        });

        const keys: any[] = res.data.keys;
        expect(keys.length).toBeGreaterThan(0);

        // Verify the first key is a valid Ed25519 public key
        const x = keys[0].x as string;
        expect(typeof x).toBe("string");
        expect(x.length).toBeGreaterThan(0);
        verifyEd25519JWKKey(x);

        // Unauthenticated access should be denied
        const body = await shouldFail(
            createClient().http.get(`/projects/${projectID}/.well-known/jwks.json`),
            401
        );
        expect(body).toBeDefined();
    });

    test("DeleteProject", async () => {
        await user.del(`/projects/${projectID}`);

        const remaining = await user.get("/projects");
        expect(remaining).toHaveLength(0);
    });
});