import { beforeAll, describe, expect, test } from "vitest";
import { createClient } from "./helpers/index.js";
import { assertErrID, assertMessage, shouldFail } from "./helpers/assert.js";
import {
    Validate,
    AnyUUID,
    AnyDate,
    AnyNumber,
    AnyString,
    AsString,
} from "./helpers/validate.js";
import { ValidPassword } from "./fixtures/auth/testdata.js";

// AuthAlreadyVerified = fail.ID(0, "AUTH", 4, true, ...) → 0_AUTH_0004_S
const ErrAuthAlreadyVerified = "0_AUTH_0004_S";

const MAILHOG_URL = process.env.MAILHOG_URL ?? "http://localhost:8025";

interface MailHogResponse {
    items: Array<{
        Content: {
            Body: string;
        };
    }>;
}

async function getLatestVerificationLink(): Promise<string> {
    const res = await fetch(`${MAILHOG_URL}/api/v2/messages`);
    if (!res.ok) throw new Error(`MailHog request failed: ${res.status}`);

    const mh: MailHogResponse = await res.json();
    if (!mh.items || mh.items.length === 0) throw new Error("no emails found");

    const body = mh.items[0].Content.Body;
    const match = body.match(/href="([^"]+)"/);
    if (!match || match.length < 2) throw new Error("verification link not found in email");

    return match[1];
}

function sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}

describe("Verification", () => {
    let user: any;
    let projectID: string;
    let initialLink: URL;

    beforeAll(async () => {
        await createClient().withCredentials("verification@mail.com", ValidPassword).register();
        user = await createClient().withCredentials("verification@mail.com", ValidPassword).login();
        const project = await user.post("/projects", { project_name: "verification", metadata: { env: "test" } });
        projectID = project.id;

        await sleep(100);
        const link = await getLatestVerificationLink();
        if (!link) throw new Error("Verification link is empty");
        initialLink = new URL(link);
    });

    test("ResendVerificationEmail", async () => {
        await user.post("/auth/verify/resend");

        await sleep(100);
        const link2 = await getLatestVerificationLink();
        if (!link2) throw new Error("Verification link is empty after resend");

        const u2 = new URL(link2);
        expect(u2.searchParams.get("token")).not.toBe(initialLink.searchParams.get("token"));
    });

    test("VerifyUser", async () => {
        const path = initialLink.pathname;
        const token = initialLink.searchParams.get("token");
        const res = await user.http.post(`${path}?token=${token}`);
        expect(res.data?.message ?? res.data?.data?.message).toMatch(/user verified/i);
    });

    test("VerifyUserAgain", async () => {
        // Same token — server should still return 200
        const path = initialLink.pathname;
        const token = initialLink.searchParams.get("token");
        const res = await user.http.post(`${path}?token=${token}`);
        expect(res.data?.message ?? res.data?.data?.message).toMatch(/user verified/i);
    });

    test("ResendVerificationEmailNotAllowed", async () => {
        const body = await shouldFail(user.post("/auth/verify/resend"), 403);
        assertErrID(body, ErrAuthAlreadyVerified);
        assertMessage(body, "user already verified");
    });

    test("SessionInfoBeforeRefreshed", async () => {
        // Use the original auth cookies (token still says is_verified: false)
        const data = await user.get("/sessions/me");
        Validate(data, {
            refresh_expire_date: AnyNumber,
            access: {
                iss: "GoAuth",
                exp: AnyNumber,
                iat: AnyNumber,
                jti: AnyUUID,
                sub: {
                    id:          AnyUUID,
                    email:       "verification@mail.com",
                    project_id:  null,
                    user_type:   "client",
                    metadata:    null,
                    session_id:  AnyUUID,
                    user_agent:  AnyString,
                    user_ip:     AnyString,
                    is_verified: false,
                    verified_at: null,
                },
            },
        });
    });

    test("SessionInfoRefreshed", async () => {
        // Fresh login — token now reflects verified state
        const freshUser = await createClient()
            .withCredentials("verification@mail.com", ValidPassword)
            .login();

        const data = await freshUser.get("/sessions/me");
        Validate(data, {
            refresh_expire_date: AnyNumber,
            access: {
                iss: "GoAuth",
                exp: AnyNumber,
                iat: AnyNumber,
                jti: AnyUUID,
                sub: {
                    id:          AnyUUID,
                    email:       "verification@mail.com",
                    project_id:  null,
                    user_type:   "client",
                    metadata:    null,
                    session_id:  AnyUUID,
                    user_agent:  AnyString,
                    user_ip:     AnyString,
                    is_verified: true,
                    verified_at: AnyDate,
                },
            },
        });
    });

    // ── Project user verification ──────────────────────────────────────────────

    describe("ProjectUser", () => {
        let projectUser: any;
        let projInitialLink: URL;

        beforeAll(async () => {
            // Register + login as project user
            await createClient().http.post(`/projects/${projectID}/register`, {
                email: "verification-project@mail.com",
                password: ValidPassword,
            });
            projectUser = await createClient()
                .withCredentials("verification-project@mail.com", ValidPassword)
                .projectLogin(projectID);

            await sleep(100);
            const link = await getLatestVerificationLink();
            if (!link) throw new Error("Project user verification link is empty");
            projInitialLink = new URL(link);
        });

        test("ResendVerificationEmailProjectUser", async () => {
            await projectUser.post("/auth/verify/resend");

            await sleep(100);
            const link3 = await getLatestVerificationLink();
            if (!link3) throw new Error("Project user verification link is empty after resend");

            const u3 = new URL(link3);
            expect(u3.searchParams.get("token")).not.toBe(projInitialLink.searchParams.get("token"));
        });

        test("VerifyProjectUser", async () => {
            const path = projInitialLink.pathname;
            const token = projInitialLink.searchParams.get("token");
            const res = await projectUser.http.post(`${path}?token=${token}`);
            expect(res.data?.message ?? res.data?.data?.message).toMatch(/user verified/i);
        });

        test("VerifyProjectUserAgain", async () => {
            const path = projInitialLink.pathname;
            const token = projInitialLink.searchParams.get("token");
            const res = await projectUser.http.post(`${path}?token=${token}`);
            expect(res.data?.message ?? res.data?.data?.message).toMatch(/user verified/i);
        });

        test("ResendVerificationEmailNotAllowedProjectUser", async () => {
            const body = await shouldFail(projectUser.post("/auth/verify/resend"), 403);
            assertErrID(body, ErrAuthAlreadyVerified);
            assertMessage(body, "user already verified");
        });

        test("ProjUserSessionInfoBeforeRefreshed", async () => {
            const data = await projectUser.get("/sessions/me");
            Validate(data, {
                refresh_expire_date: AnyNumber,
                access: {
                    iss: AsString(projectID, AnyUUID),
                    exp: AnyNumber,
                    iat: AnyNumber,
                    jti: AnyUUID,
                    sub: {
                        id:          AnyUUID,
                        email:       "verification-project@mail.com",
                        project_id:  AsString(projectID, AnyUUID),
                        user_type:   "project",
                        metadata:    {},
                        session_id:  AnyUUID,
                        user_agent:  AnyString,
                        user_ip:     AnyString,
                        is_verified: false,
                        verified_at: null,
                    },
                },
            });
        });

        test("ProjUserSessionInfoRefreshed", async () => {
            const freshProjUser = await createClient()
                .withCredentials("verification-project@mail.com", ValidPassword)
                .projectLogin(projectID);

            const data = await freshProjUser.get("/sessions/me");
            Validate(data, {
                refresh_expire_date: AnyNumber,
                access: {
                    iss: AsString(projectID, AnyUUID),
                    exp: AnyNumber,
                    iat: AnyNumber,
                    jti: AnyUUID,
                    sub: {
                        id:          AnyUUID,
                        email:       "verification-project@mail.com",
                        project_id:  AsString(projectID, AnyUUID),
                        user_type:   "project",
                        metadata:    {},
                        session_id:  AnyUUID,
                        user_agent:  AnyString,
                        user_ip:     AnyString,
                        is_verified: true,
                        verified_at: AnyDate,
                    },
                },
            });
        });
    });
});