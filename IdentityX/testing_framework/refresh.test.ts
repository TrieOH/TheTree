import { afterAll, describe, expect, test } from "vitest";
import { createClient } from "./helpers/index.js";
import { assertErrID, assertMessage, shouldFail } from "./helpers/assert.js";
import { dbQuery, closeDB } from "./helpers/db.js";
import { ValidPassword } from "./fixtures/auth/testdata.js";

// Error IDs
// TokenReuseIdentified = fail.ID(0, "TOKEN", 18, false, ...) → 0_TOKEN_0018_D
// SessionNotFound      = fail.ID(0, "SESSION", 0, true, ...)  → 0_SESSION_0000_S
const ErrTokenReuseIdentified = "0_TOKEN_0018_D";
const ErrSessionNotFound      = "0_SESSION_0000_S";

afterAll(async () => {
    await closeDB();
});

// ============================================================================
// REFRESH TESTS
// ============================================================================

describe("Refresh", () => {
    test("RefreshSuccess", async () => {
        await createClient().withCredentials("refresh_success@mail.com", ValidPassword).register();
        const user = await createClient().withCredentials("refresh_success@mail.com", ValidPassword).login();

        const oldAccess  = user.auth!.accessToken;
        const oldRefresh = user.auth!.refreshToken;

        const refreshed = await user.refresh();

        expect(refreshed.auth!.accessToken).not.toBe(oldAccess);
        expect(refreshed.auth!.refreshToken).not.toBe(oldRefresh);
    });

    test("UseOldTokenError", async () => {
        await createClient().withCredentials("old_token_user@mail.com", ValidPassword).register();
        const oldAuth = await createClient().withCredentials("old_token_user@mail.com", ValidPassword).login();

        // Refresh but discard the new auth — old tokens are now invalidated
        await oldAuth.refresh();

        // Try to use the old access token
        const body = await shouldFail(oldAuth.get("/sessions"), 401);
        assertErrID(body, ErrTokenReuseIdentified);
        assertMessage(body, "refresh token reuse not allowed");
    });

    test("RefreshRevokedToken", async () => {
        await createClient().withCredentials("revoked_refresh_user@mail.com", ValidPassword).register();
        const user = await createClient().withCredentials("revoked_refresh_user@mail.com", ValidPassword).login();
        const refreshToken = user.auth!.refreshToken;

        // Manually revoke the session in the database
        await dbQuery(`
      UPDATE sessions
      SET revoked_at = NOW()
      WHERE token_id = (
        SELECT token_id FROM sessions
        WHERE identity_id = (
          SELECT id FROM identities
          WHERE entity_id = (
            SELECT id FROM users WHERE email = 'revoked_refresh_user@mail.com'
          )
        ) LIMIT 1
      )
    `);

        // Attempt to use the revoked refresh token
        const denied = createClient();
        const body = await shouldFail(
            denied.http.post("/auth/refresh", {}, {
                headers: { Cookie: `refresh_token=${refreshToken}` },
            }),
            401
        );
        assertErrID(body, ErrSessionNotFound);
        assertMessage(body, "session not found or revoked");
    });

    test("ConcurrentRefresh", async () => {
        await createClient().withCredentials("concurrent_refresh@mail.com", ValidPassword).register();
        const user = await createClient().withCredentials("concurrent_refresh@mail.com", ValidPassword).login();
        const refreshToken = user.auth!.refreshToken;

        const concurrency = 5;

        const results = await Promise.all(
            Array.from({ length: concurrency }, () =>
                createClient().http.post("/auth/refresh", {}, {
                    headers: { Cookie: `refresh_token=${refreshToken}` },
                    validateStatus: () => true,
                }).then(res => res.status)
            )
        );

        const successCount = results.filter(s => s === 200).length;
        const failCount    = results.filter(s => s === 401).length;

        expect(successCount).toBe(1);
        expect(failCount).toBe(concurrency - 1);
    });
});