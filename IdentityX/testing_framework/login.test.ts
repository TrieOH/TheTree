import { beforeAll, describe, test } from "vitest";
import { createClient } from "./helpers/index.js";
import { assertErrID, assertMessage, shouldFail } from "./helpers/assert.js";
import { ValidPassword } from "./fixtures/auth/testdata.js";

// Error IDs
// AuthInvalidCredentials = fail.ID(0, "AUTH", 1, false, ...) → 0_AUTH_0001_D
// SessionRevoked         = fail.ID(0, "SESSION", 0, false, ...) → 0_SESSION_0000_D
const ErrInvalidCredentials = "0_AUTH_0001_D";
const ErrSessionRevoked     = "0_SESSION_0000_D";

// ============================================================================
// LOGIN TESTS
// ============================================================================

describe("Login", () => {
    // Register a shared user once before all login tests
    beforeAll(async () => {
        const client = createClient();
        await client.withCredentials("login@mail.com", ValidPassword).register();
    });

    test("WrongPassword", async () => {
        const client = createClient();
        const body = await shouldFail(
            client.http.post("/auth/login", {
                email: "login@mail.com",
                password: "WrongPass123!",
            }),
            401
        );
        assertErrID(body, ErrInvalidCredentials);
        assertMessage(body, "invalid email or password");
    });

    test("WrongEmail", async () => {
        const client = createClient();
        const body = await shouldFail(
            client.http.post("/auth/login", {
                email: "wrong@mail.com",
                password: ValidPassword,
            }),
            401
        );
        assertErrID(body, ErrInvalidCredentials);
        assertMessage(body, "invalid email or password");
    });

    test("Success", async () => {
        const client = createClient();
        await client.withCredentials("login@mail.com", ValidPassword).login();
    });

    test("Logout", async () => {
        const client = createClient();
        const loggedIn = await client
            .withCredentials("login@mail.com", ValidPassword)
            .login();

        await loggedIn.logout();

        // Try using the revoked session — should be rejected
        const body = await shouldFail(loggedIn.logout(), 401);
        assertErrID(body, ErrSessionRevoked);
        assertMessage(body, "session not found or revoked");
    });
});