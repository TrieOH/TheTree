import { beforeAll, describe, expect, test } from "vitest";
import { createClient } from "./helpers/index.js";
import { assertErrID, shouldFail } from "./helpers/assert.js";
import { ValidPassword } from "./fixtures/auth/testdata.js";

// SessionRevoked       = fail.ID(0, "SESSION", 0, false, ...) → 0_SESSION_0000_D
// AuthInvalidCredentials = fail.ID(0, "AUTH", 1, false, ...)  → 0_AUTH_0001_D
// AuthTokenAlreadyUsed = fail.ID(0, "AUTH", 7, true, ...)     → 0_AUTH_0007_S
const ErrSessionRevoked        = "0_SESSION_0000_D";
const ErrAuthInvalidCredentials = "0_AUTH_0001_D";
const ErrAuthTokenAlreadyUsed  = "0_AUTH_0007_S";

const MAILHOG_URL = process.env.MAILHOG_URL ?? "http://localhost:8025";

interface MHMessage {
    Content: { Body: string };
    Raw: { To: string[] };
}
interface MHResponse {
    items: MHMessage[];
}

async function getLatestResetPasswordLink(toEmail: string): Promise<string> {
    const res = await fetch(`${MAILHOG_URL}/api/v2/messages`);
    if (!res.ok) throw new Error(`MailHog request failed: ${res.status}`);

    const mh: MHResponse = await res.json();
    if (!mh.items || mh.items.length === 0) throw new Error("no emails found");

    const hrefRe = /href="([^"]+)"/;
    const resetRe = /\/reset\?token=/;

    for (const item of mh.items) {
        const isForEmail = item.Raw.To.includes(toEmail);
        if (!isForEmail) continue;

        const match = item.Content.Body.match(hrefRe);
        if (match && match[1] && resetRe.test(match[1])) {
            return match[1];
        }
    }

    throw new Error(`reset password link not found for ${toEmail}`);
}

function sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}

describe("ResetPassword", () => {
    test("GlobalUserResetPassword", async () => {
        const userEmail = "global-reset@mail.com";
        const newPassword = "NewSecurePassword123!";

        await createClient().withCredentials(userEmail, ValidPassword).register();
        const userSession = await createClient().withCredentials(userEmail, ValidPassword).login();

        // Request forgot password
        await userSession.http.post("/auth/forgot-password", { email: userEmail });

        await sleep(500);

        const link = await getLatestResetPasswordLink(userEmail);
        const u = new URL(link);
        const token = u.searchParams.get("token");

        // Reset password
        await userSession.http.post(`/auth/reset-password?token=${token}`, { new_password: newPassword });

        // Old session should be revoked
        const revokedBody = await shouldFail(userSession.get("/sessions/me"), 401);
        assertErrID(revokedBody, ErrSessionRevoked);

        // Old password should fail
        const badLoginBody = await shouldFail(
            createClient().http.post("/auth/login", { email: userEmail, password: ValidPassword }),
            401
        );
        assertErrID(badLoginBody, ErrAuthInvalidCredentials);

        // New password should succeed
        await createClient().withCredentials(userEmail, newPassword).login();
    });

    test("ProjectUserResetPassword", async () => {
        const projectOwnerEmail = "project-owner-reset@mail.com";
        const projectUserEmail = "project-user-reset@mail.com";
        const projectUserNewPassword = "ProjectNewPass123!";

        await createClient().withCredentials(projectOwnerEmail, ValidPassword).register();
        const projectOwner = await createClient().withCredentials(projectOwnerEmail, ValidPassword).login();
        const project = await projectOwner.post("/projects", { project_name: "ResetProject", metadata: { env: "test" } });
        const projectID = project.id;

        // Register and login project user
        await createClient().http.post(`/projects/${projectID}/register`, {
            email: projectUserEmail,
            password: ValidPassword,
        });
        const projectUserSession = await createClient()
            .withCredentials(projectUserEmail, ValidPassword)
            .projectLogin(projectID);

        // Request forgot password for project user
        await createClient().http.post("/auth/forgot-password", {
            email: projectUserEmail,
            project_id: projectID,
        });

        await sleep(500);

        const link = await getLatestResetPasswordLink(projectUserEmail);
        const u = new URL(link);
        const token = u.searchParams.get("token");

        // Reset password
        await createClient().http.post(`/auth/reset-password?token=${token}`, {
            new_password: projectUserNewPassword,
        });

        // Old session should be revoked
        const revokedBody = await shouldFail(projectUserSession.get("/sessions/me"), 401);
        assertErrID(revokedBody, ErrSessionRevoked);

        // Old password should fail
        const badLoginBody = await shouldFail(
            createClient().http.post(`/projects/${projectID}/login`, {
                email: projectUserEmail,
                password: ValidPassword,
            }),
            401
        );
        assertErrID(badLoginBody, ErrAuthInvalidCredentials);

        // New password should succeed
        await createClient()
            .withCredentials(projectUserEmail, projectUserNewPassword)
            .projectLogin(projectID);
    });

    test("TokenReusePrevention", async () => {
        const email = "reuse-token@mail.com";
        await createClient().withCredentials(email, ValidPassword).register();

        await createClient().http.post("/auth/forgot-password", { email });

        await sleep(500);

        const link = await getLatestResetPasswordLink(email);
        const u = new URL(link);
        const token = u.searchParams.get("token");

        // First use — should succeed
        await createClient().http.post(`/auth/reset-password?token=${token}`, { new_password: "NewPassword1!" });

        // Second use — should fail
        const body = await shouldFail(
            createClient().http.post(`/auth/reset-password?token=${token}`, { new_password: "NewPassword2!" }),
            403
        );
        assertErrID(body, ErrAuthTokenAlreadyUsed);
    });
});