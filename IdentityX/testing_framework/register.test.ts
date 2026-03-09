import { describe, test } from "vitest";
import { createClient } from "./helpers/index.js";
import {
    assertErrID,
    assertMessage,
    assertValidationError,
    shouldFail,
} from "./helpers/assert.js";
import {
    ValidPassword,
    ValidationTests,
    WeakPasswordTests,
} from "./fixtures/auth/testdata.js";

// Error IDs - mirror Go errx constants
const ErrRequestValidation = "0_REQ_0004_D";
const ErrEmailAlreadyUsed  = "0_AUTH_0000_D";

// ============================================================================
// REGISTER TESTS
// ============================================================================

describe("Register", () => {

    describe("Validation", () => {
        for (const spec of ValidationTests) {
            test(spec.name, async () => {
                const client = createClient();
                const body = await shouldFail(
                    client.http.post("/auth/register", { email: spec.email, password: spec.pass }),
                    400
                );
                assertErrID(body, ErrRequestValidation);
                assertValidationError(body, ...spec.errors);
            });
        }
    });

    describe("WeakPasswords", () => {
        for (const [i, spec] of WeakPasswordTests.entries()) {
            test(spec.name, async () => {
                const client = createClient();
                const body = await shouldFail(
                    client.http.post("/auth/register", {
                        email: `weak${i}@mail.com`,
                        password: spec.password,
                    }),
                    400
                );
                assertErrID(body, ErrRequestValidation);
                assertValidationError(body, ...spec.errors);
            });
        }
    });

    test("PasswordTooLong", async () => {
        const client = createClient();
        const longPass = "A1@" + "a".repeat(70); // 73 chars > 72 limit
        const body = await shouldFail(
            client.http.post("/auth/register", {
                email: "longpass@mail.com",
                password: longPass,
            }),
            400
        );
        assertErrID(body, ErrRequestValidation);
        assertValidationError(body, "password must be at most 72 characters long");
    });

    // FIXME: Current validation package does not consider trailing spaces as a
    // valid email. Once fixed, switch email to "  MixedCase@Example.Com  ".
    test("EmailNormalization", async () => {
        const client = createClient();
        const email           = "MixedCase@Example.Com";
        const normalizedEmail = "mixedcase@example.com";

        // Register with mixed-case email
        await client.http.post("/auth/register", { email, password: ValidPassword });

        // Login with normalized email
        const loginNorm = await client.http.post("/auth/login", {
            email: normalizedEmail,
            password: ValidPassword,
        });
        assertMessage(loginNorm.data, "Logged in");

        // Login with original mixed-case email (normalization on input)
        const loginOrig = await client.http.post("/auth/login", {
            email,
            password: ValidPassword,
        });
        assertMessage(loginOrig.data, "Logged in");

        // Re-register with normalized email — should conflict
        const body = await shouldFail(
            client.http.post("/auth/register", {
                email: normalizedEmail,
                password: ValidPassword,
            }),
            409
        );
        assertErrID(body, ErrEmailAlreadyUsed);
        assertMessage(body, "email already in use");
    });

    test("Success", async () => {
        const client = createClient();
        await client.withCredentials("new@mail.com", ValidPassword).register();
    });

    test("DuplicateEmail", async () => {
        const client = createClient();
        const email = "duplicate@mail.com";

        await client.withCredentials(email, ValidPassword).register();

        const body = await shouldFail(
            client.http.post("/auth/register", { email, password: ValidPassword }),
            409
        );
        assertErrID(body, ErrEmailAlreadyUsed);
        assertMessage(body, "email already in use");
    });
});