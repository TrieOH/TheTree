import { afterAll, beforeAll, describe, expect, test } from "vitest";
import { createClient } from "./helpers/index.js";
import { assertErrID, assertMessage, shouldFail } from "./helpers/assert.js";
import { Validate, AnyUUID, AnyString, AnyNumber } from "./helpers/validate.js";
import { dbQuery, closeDB } from "./helpers/db.js";
import { ValidPassword } from "./fixtures/auth/testdata.js";

// Error IDs
// SessionSelfRevokeForbidden = fail.ID(0, "SESSION", 1, true, ...) → 0_SESSION_0001_S
// SessionRevoked             = fail.ID(0, "SESSION", 0, false, ...) → 0_SESSION_0000_D
const ErrSessionSelfRevokeForbidden = "0_SESSION_0001_S";
const ErrSessionRevoked             = "0_SESSION_0000_D";

afterAll(async () => {
    await closeDB();
});

// ============================================================================
// SESSIONS TESTS
// ============================================================================

describe("Sessions", () => {
    let sharedUser: Awaited<ReturnType<typeof createClient>>;

    beforeAll(async () => {
        await createClient().withCredentials("sessions@mail.com", ValidPassword).register();
        sharedUser = await createClient().withCredentials("sessions@mail.com", ValidPassword).login();
    });

    test("ListSessions", async () => {
        const sessions = await sharedUser.get("/sessions");
        expect(sessions).toHaveLength(1);
    });

    test("MultipleLoginsSessions", async () => {
        await createClient().withCredentials("sessions@mail.com", ValidPassword).login();
        await createClient().withCredentials("sessions@mail.com", ValidPassword).login();
        const user4 = await createClient().withCredentials("sessions@mail.com", ValidPassword).login();

        const sessions: any[] = await user4.get("/sessions");
        expect(sessions).toHaveLength(4);

        const currentSessionID = sessions[0].session_id;
        const oldestSessionID  = sessions[3].session_id;

        // Can't revoke current session
        const forbiddenBody = await shouldFail(user4.del(`/sessions/${currentSessionID}`), 403);
        assertErrID(forbiddenBody, ErrSessionSelfRevokeForbidden);
        assertMessage(forbiddenBody, "cannot revoke the currently active session");

        // Revoke oldest session
        await user4.del(`/sessions/${oldestSessionID}`);

        // Should have 3 sessions now
        const updated: any[] = await user4.get("/sessions");
        expect(updated).toHaveLength(3);
    });

    test("RevokeOtherSessions", async () => {
        await createClient().withCredentials("revoke-others@mail.com", ValidPassword).register();
        const user = await createClient().withCredentials("revoke-others@mail.com", ValidPassword).login();

        await createClient().withCredentials("revoke-others@mail.com", ValidPassword).login();
        await createClient().withCredentials("revoke-others@mail.com", ValidPassword).login();

        await user.del("/sessions/others");

        const sessions: any[] = await user.get("/sessions");
        expect(sessions).toHaveLength(1);
    });

    test("SessionInfo", async () => {
        await createClient().withCredentials("session-me@mail.com", ValidPassword).register();
        const user = await createClient().withCredentials("session-me@mail.com", ValidPassword).login();

        const data = await user.get("/sessions/me");

        Validate(data, {
            refresh_expire_date: AnyNumber,
            access: {
                iss: "GoAuth",
                exp: AnyNumber,
                iat: AnyNumber,
                jti: AnyUUID,
                sub: {
                    id:         AnyUUID,
                    email:      "session-me@mail.com",
                    project_id: null,
                    user_type:  "client",
                    metadata:   null,
                    session_id: AnyUUID,
                    user_agent: AnyString,
                    user_ip:    AnyString,
                },
            },
        });
    });

    test("RevokeAllSessions", async () => {
        await createClient().withCredentials("revoke-all@mail.com", ValidPassword).register();
        const user = await createClient().withCredentials("revoke-all@mail.com", ValidPassword).login();

        await createClient().withCredentials("revoke-all@mail.com", ValidPassword).login();
        await createClient().withCredentials("revoke-all@mail.com", ValidPassword).login();

        await user.del("/sessions");

        const body = await shouldFail(user.get("/sessions"), 401);
        assertErrID(body, ErrSessionRevoked);
        assertMessage(body, "session not found or revoked");
    });

    test("ExpiredSessionNotListed", async () => {
        await createClient().withCredentials("expired@mail.com", ValidPassword).register();
        const user = await createClient().withCredentials("expired@mail.com", ValidPassword).login();

        const { rows: [{ id: userID }] } = await dbQuery<{ id: string }>(
            "SELECT id FROM users WHERE email = 'expired@mail.com'"
        );

        const { rows: [{ id: identityID }] } = await dbQuery<{ id: string }>(
            `INSERT INTO identities (type, entity_id)
       VALUES ('client', $1)
       ON CONFLICT (type, entity_id) DO UPDATE SET type = 'client'
       RETURNING id`,
            [userID]
        );

        await dbQuery(
            `INSERT INTO sessions (identity_id, issued_at, user_agent, user_ip, expires_at, created_at, updated_at, user_type)
       VALUES ($1, NOW() - INTERVAL '2 days', 'Expired Agent', '127.0.0.1', NOW() - INTERVAL '1 day', NOW(), NOW(), 'client')`,
            [identityID]
        );

        const sessions: any[] = await user.get("/sessions");
        expect(sessions).toHaveLength(1);
    });

    test("RevokedSessionNotListed", async () => {
        await createClient().withCredentials("revoked@mail.com", ValidPassword).register();
        const user = await createClient().withCredentials("revoked@mail.com", ValidPassword).login();

        const { rows: [{ id: userID }] } = await dbQuery<{ id: string }>(
            "SELECT id FROM users WHERE email = 'revoked@mail.com'"
        );

        const { rows: [{ id: identityID }] } = await dbQuery<{ id: string }>(
            `INSERT INTO identities (type, entity_id)
       VALUES ('client', $1)
       ON CONFLICT (type, entity_id) DO UPDATE SET type = 'client'
       RETURNING id`,
            [userID]
        );

        await dbQuery(
            `INSERT INTO sessions (identity_id, issued_at, user_agent, user_ip, revoked_at, created_at, updated_at, expires_at, user_type)
       VALUES ($1, NOW() - INTERVAL '2 days', 'Revoked Agent', '127.0.0.1', NOW() - INTERVAL '1 day', NOW(), NOW(), NOW() + INTERVAL '1 day', 'client')`,
            [identityID]
        );

        const sessions: any[] = await user.get("/sessions");
        expect(sessions).toHaveLength(1);
    });

    test("SessionLeakage", async () => {
        await createClient().withCredentials("leakage-client@mail.com", ValidPassword).register();
        await createClient().withCredentials("leakage-project-owner@mail.com", ValidPassword).register();

        const owner = await createClient()
            .withCredentials("leakage-project-owner@mail.com", ValidPassword)
            .login();

        const project = await owner.post("/projects", {
            project_name: "Leakage Test Project",
            metadata: { env: "test" },
        });
        const projectID = project.id;

        // Register and log in client user as project user
        await createClient().http.post(`/projects/${projectID}/register`, {
            email: "leakage-client@mail.com",
            password: ValidPassword,
        });

        const clientSession = await createClient()
            .withCredentials("leakage-client@mail.com", ValidPassword)
            .login();

        // Create a project user session (separate identity type)
        await createClient().http.post(`/projects/${projectID}/login`, {
            email: "leakage-client@mail.com",
            password: ValidPassword,
        });

        // Client should only see their own session, not the project user's
        const sessions: any[] = await clientSession.get("/sessions");
        expect(sessions).toHaveLength(1);
        expect(sessions[0].user_type).toBe("client");
    });
});