import { beforeAll, describe, expect, test } from "vitest";
import { createClient } from "./helpers/index.js";
import {
    assertErrID,
    assertMessage,
    assertValidationError,
    shouldFail,
} from "./helpers/assert.js";
import { Validate, AnyUUID, AnyDate, AnyNumber, AnyString, AsString } from "./helpers/validate.js";
import { ValidPassword, ValidationTests, WeakPasswordTests } from "./fixtures/auth/testdata.js";
import { importSPKI, jwtVerify } from "jose";

// Error IDs
// RequestValidationError     = fail.ID(0, "REQ", 4, false, ...)      → 0_REQ_0004_D
// AuthEmailAlreadyUsed       = fail.ID(0, "AUTH", 0, false, ...)     → 0_AUTH_0000_D
// AuthInvalidCredentials     = fail.ID(0, "AUTH", 1, false, ...)     → 0_AUTH_0001_D
// AuthNotClient              = fail.ID(0, "AUTH", 2, true, ...)      → 0_AUTH_0002_S
// SessionRevoked             = fail.ID(0, "SESSION", 0, false, ...)  → 0_SESSION_0000_D
// SessionSelfRevokeForbidden = fail.ID(0, "SESSION", 1, true, ...)   → 0_SESSION_0001_S
// TokenReuseIdentified       = fail.ID(0, "TOKEN", 18, false, ...)   → 0_TOKEN_0018_D
// SchemaInvalidSchemaType    = fail.ID(0, "SCHEMA", 2, true, ...)    → 0_SCHEMA_0002_S
// SchemaInvalidFlowID        = fail.ID(0, "SCHEMA", 2, false, ...)   → 0_SCHEMA_0002_D
// SchemaMetadataNotAllowed   = fail.ID(0, "SCHEMA", 7, true, ...)    → 0_SCHEMA_0007_S
const ErrRequestValidation     = "0_REQ_0004_D";
const ErrEmailAlreadyUsed      = "0_AUTH_0000_D";
const ErrInvalidCredentials    = "0_AUTH_0001_D";
const ErrAuthNotClient         = "0_AUTH_0002_S";
const ErrSessionRevoked        = "0_SESSION_0000_D";
const ErrSessionSelfRevoke     = "0_SESSION_0001_S";
const ErrTokenReuseIdentified  = "0_TOKEN_0018_D";
const ErrSchemaInvalidType     = "0_SCHEMA_0002_S";
const ErrSchemaInvalidFlowID   = "0_SCHEMA_0002_D";
const ErrSchemaMetadata        = "0_SCHEMA_0007_S";

// ============================================================================
// Ed25519 helpers
// ============================================================================

function decodeBase64Url(s: string): Buffer {
    const base64 = s.replace(/-/g, "+").replace(/_/g, "/");
    const padded = base64.padEnd(base64.length + (4 - (base64.length % 4)) % 4, "=");
    return Buffer.from(padded, "base64");
}

function verifyEd25519Key(x: string): void {
    const bytes = decodeBase64Url(x);
    expect(bytes.length).toBe(32);
}

// Convert raw Ed25519 public key bytes to a SubjectPublicKeyInfo (SPKI) PEM
// so jose can import it for JWT verification
function rawEd25519ToSPKIPem(rawBytes: Buffer): string {
    // DER prefix for Ed25519 SubjectPublicKeyInfo
    const prefix = Buffer.from("302a300506032b6570032100", "hex");
    const der = Buffer.concat([prefix, rawBytes]);
    const b64 = der.toString("base64");
    return `-----BEGIN PUBLIC KEY-----\n${b64}\n-----END PUBLIC KEY-----`;
}

// ============================================================================
// PROJECT USERS TESTS
// ============================================================================

describe("Project Users", () => {
    let ownerUser: any;
    let projectID: string;
    let projectUserEmail: string;
    let projectUserPassword: string;

    beforeAll(async () => {
        // Create owner and project
        await createClient().withCredentials("client@mail.com", ValidPassword).register();
        ownerUser = await createClient().withCredentials("client@mail.com", ValidPassword).login();
        const project = await ownerUser.post("/projects", {
            project_name: "test project",
            metadata: { env: "test" },
        });
        projectID = project.id;

        // Register project user with same email as client (different environment)
        projectUserEmail = "client@mail.com";
        projectUserPassword = ValidPassword;
        await createClient().http.post(`/projects/${projectID}/register`, {
            email: projectUserEmail,
            password: projectUserPassword,
        });
    });

    // --------------------------------------------------------------------------
    // REGISTER
    // --------------------------------------------------------------------------

    describe("ProjectUsersRegister", () => {
        test("WrongFormatIDProjectRegister", async () => {
            const body = await shouldFail(
                createClient().http.post("/projects/wrong-format/register", {
                    email: "client@mail.com",
                    password: ValidPassword,
                }),
                400
            );
            assertErrID(body, ErrRequestValidation);
            assertMessage(body, "Validation failed");
        });

        test("InvalidProjectRegister", async () => {
            const fakeID = "00000000-0000-7000-8000-000000000000";
            const body = await shouldFail(
                createClient().http.post(`/projects/${fakeID}/register`, {
                    email: "client@mail.com",
                    password: ValidPassword,
                }),
                400
            );
            assertMessage(body, "can't register on a non existant project");
        });

        describe("ValidationProjectRegister", () => {
            for (const spec of ValidationTests) {
                test(spec.name, async () => {
                    const body = await shouldFail(
                        createClient().http.post(`/projects/${projectID}/register`, {
                            email: spec.email,
                            password: spec.pass,
                        }),
                        400
                    );
                    assertErrID(body, ErrRequestValidation);
                    assertValidationError(body, ...spec.errors);
                });
            }
        });

        describe("WeakPasswordValidationProjectRegister", () => {
            for (const [i, spec] of WeakPasswordTests.entries()) {
                test(spec.name, async () => {
                    const body = await shouldFail(
                        createClient().http.post(`/projects/${projectID}/register`, {
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

        test("DuplicateEmailProjectRegister", async () => {
            const body = await shouldFail(
                createClient().http.post(`/projects/${projectID}/register`, {
                    email: "client@mail.com",
                    password: ValidPassword,
                }),
                409
            );
            assertErrID(body, ErrEmailAlreadyUsed);
            assertMessage(body, "email already in use");
        });

        test("InvalidSchemaTypeRegister", async () => {
            const body = await shouldFail(
                createClient().http.post(
                    `/projects/${projectID}/register?schema_type=invalid`,
                    { email: "invalid_schema@email.com", password: ValidPassword }
                ),
                400
            );
            assertErrID(body, ErrSchemaInvalidType);
            assertMessage(body, "invalid schema type");
        });

        test("FlowIDSameAsTypeRegister", async () => {
            const body = await shouldFail(
                createClient().http.post(
                    `/projects/${projectID}/register?schema_type=context&flow_id=context`,
                    { email: "flow_same_as_type@email.com", password: ValidPassword }
                ),
                400
            );
            assertErrID(body, ErrSchemaInvalidFlowID);
            assertMessage(body, "flow id can't be the same as a schema type");
        });

        test("MetadataRegisterOnCoreDenied", async () => {
            const body = await shouldFail(
                createClient().http.post(`/projects/${projectID}/register`, {
                    email: "metadata_denied@email.com",
                    password: ValidPassword,
                    custom_fields: { curso: "Ciência da Computação" },
                }),
                400
            );
            assertErrID(body, ErrSchemaMetadata);
            assertMessage(body, "custom fields are not allowed for core schema");
        });

        test("SuccessProjectRegister", async () => {
            await createClient().http.post(`/projects/${projectID}/register`, {
                email: "new@mail.com",
                password: ValidPassword,
            });
        });
    });

    // --------------------------------------------------------------------------
    // LOGIN
    // --------------------------------------------------------------------------

    describe("ProjectUsersLogin", () => {
        test("WrongPassword", async () => {
            const body = await shouldFail(
                createClient().http.post(`/projects/${projectID}/login`, {
                    email: projectUserEmail,
                    password: "WrongPass123!",
                }),
                401
            );
            assertErrID(body, ErrInvalidCredentials);
            assertMessage(body, "invalid email or password");
        });

        test("WrongEmail", async () => {
            const body = await shouldFail(
                createClient().http.post(`/projects/${projectID}/login`, {
                    email: "wrong@mail.com",
                    password: projectUserPassword,
                }),
                401
            );
            assertErrID(body, ErrInvalidCredentials);
            assertMessage(body, "invalid email or password");
        });

        test("Success", async () => {
            await createClient()
                .withCredentials(projectUserEmail, projectUserPassword)
                .projectLogin(projectID);
        });

        test("Logout", async () => {
            const loggedIn = await createClient()
                .withCredentials(projectUserEmail, projectUserPassword)
                .projectLogin(projectID);

            await loggedIn.logout();

            const body = await shouldFail(loggedIn.logout(), 401);
            assertErrID(body, ErrSessionRevoked);
            assertMessage(body, "session not found or revoked");
        });
    });

    // --------------------------------------------------------------------------
    // SESSIONS
    // --------------------------------------------------------------------------

    describe("ProjectUsersSession", () => {
        let sessionUser: any;
        const sessionEmail = "sessions@mail.com";

        beforeAll(async () => {
            await createClient().http.post(`/projects/${projectID}/register`, {
                email: sessionEmail,
                password: ValidPassword,
            });
            sessionUser = await createClient()
                .withCredentials(sessionEmail, ValidPassword)
                .projectLogin(projectID);
        });

        test("ListSessions", async () => {
            const sessions = await sessionUser.get("/sessions");
            expect(sessions).toHaveLength(1);
        });

        test("MultipleLoginsSessions", async () => {
            await createClient().withCredentials(sessionEmail, ValidPassword).projectLogin(projectID);
            await createClient().withCredentials(sessionEmail, ValidPassword).projectLogin(projectID);
            const user4 = await createClient().withCredentials(sessionEmail, ValidPassword).projectLogin(projectID);

            const sessions: any[] = await user4.get("/sessions");
            expect(sessions).toHaveLength(4);

            const currentSessionID = sessions[0].session_id;
            const oldestSessionID  = sessions[3].session_id;

            const forbiddenBody = await shouldFail(user4.del(`/sessions/${currentSessionID}`), 403);
            assertErrID(forbiddenBody, ErrSessionSelfRevoke);
            assertMessage(forbiddenBody, "cannot revoke the currently active session");

            await user4.del(`/sessions/${oldestSessionID}`);

            const updated: any[] = await user4.get("/sessions");
            expect(updated).toHaveLength(3);
        });

        test("RevokeOtherSessions", async () => {
            const email = "revoke-others-project@mail.com";
            await createClient().http.post(`/projects/${projectID}/register`, {
                email,
                password: ValidPassword,
            });
            const revokeOthers = await createClient()
                .withCredentials(email, ValidPassword)
                .projectLogin(projectID);

            await createClient().withCredentials(email, ValidPassword).projectLogin(projectID);
            await createClient().withCredentials(email, ValidPassword).projectLogin(projectID);

            await revokeOthers.del("/sessions/others");

            const sessions: any[] = await revokeOthers.get("/sessions");
            expect(sessions).toHaveLength(1);
        });

        test("SessionInfo", async () => {
            const email = "session-me@mail.com";
            await createClient().http.post(`/projects/${projectID}/register`, {
                email,
                password: ValidPassword,
            });
            const infoUser = await createClient()
                .withCredentials(email, ValidPassword)
                .projectLogin(projectID);

            const data = await infoUser.get("/sessions/me");

            Validate(data, {
                refresh_expire_date: AnyNumber,
                access: {
                    iss: AsString(projectID, AnyUUID),
                    exp: AnyNumber,
                    iat: AnyNumber,
                    jti: AnyUUID,
                    sub: {
                        id:         AnyUUID,
                        email:      email,
                        project_id: AsString(projectID, AnyUUID),
                        user_type:  "project",
                        metadata:   {},
                        session_id: AnyUUID,
                        user_agent: AnyString,
                        user_ip:    AnyString,
                    },
                },
            });
        });

        test("RevokeAllSessions", async () => {
            const email = "revoke-all@mail.com";
            await createClient().http.post(`/projects/${projectID}/register`, {
                email,
                password: ValidPassword,
            });
            const revoked = await createClient()
                .withCredentials(email, ValidPassword)
                .projectLogin(projectID);

            await createClient().withCredentials(email, ValidPassword).projectLogin(projectID);
            await createClient().withCredentials(email, ValidPassword).projectLogin(projectID);

            await revoked.del("/sessions");

            const body = await shouldFail(revoked.get("/sessions"), 401);
            assertErrID(body, ErrSessionRevoked);
            assertMessage(body, "session not found or revoked");
        });
    });

    // --------------------------------------------------------------------------
    // REFRESH REUSE
    // --------------------------------------------------------------------------

    test("ProjectUserRefreshReuse", async () => {
        const email = "refresh@mail.com";
        await createClient().http.post(`/projects/${projectID}/register`, {
            email,
            password: ValidPassword,
        });
        const refreshUser = await createClient()
            .withCredentials(email, ValidPassword)
            .projectLogin(projectID);

        const oldAuth = refreshUser.auth!;
        const refreshed = await refreshUser.refresh();

        expect(refreshed.auth!.accessToken).not.toBe(oldAuth.accessToken);
        expect(refreshed.auth!.refreshToken).not.toBe(oldAuth.refreshToken);

        // Old tokens should be rejected
        const oldClient = createClient().withAuth(oldAuth);
        const body = await shouldFail(oldClient.get("/sessions"), 401);
        assertErrID(body, ErrTokenReuseIdentified);
        assertMessage(body, "refresh token reuse not allowed");
    });

    // --------------------------------------------------------------------------
    // PROJECT ACCESS DENIED FOR PROJECT USERS
    // --------------------------------------------------------------------------

    describe("ProjectUsersProjects", () => {
        let nestedUser: any;

        beforeAll(async () => {
            const email = "nested_creator@mail.com";
            await createClient().http.post(`/projects/${projectID}/register`, {
                email,
                password: ValidPassword,
            });
            nestedUser = await createClient()
                .withCredentials(email, ValidPassword)
                .projectLogin(projectID);
        });

        const expectForbidden = async (promise: Promise<any>) => {
            const body = await shouldFail(promise, 403);
            assertErrID(body, ErrAuthNotClient);
            assertMessage(body, "only clients can access this endpoint");
        };

        test("CreateProject", async () => {
            await expectForbidden(
                nestedUser.post("/projects", { project_name: "Test Project", metadata: { env: "test" } })
            );
        });

        test("ListProjects", async () => {
            await expectForbidden(nestedUser.get("/projects"));
        });

        test("GetProject", async () => {
            await expectForbidden(nestedUser.get(`/projects/${projectID}`));
        });

        test("UpdateProject", async () => {
            await expectForbidden(
                nestedUser.patch(`/projects/${projectID}`, {
                    project_name: "Updated Project",
                    metadata: { env: "prod" },
                })
            );
        });

        test("GetProjectJWKS", async () => {
            const body = await shouldFail(
                nestedUser.get(`/projects/${projectID}/.well-known/jwks.json`),
                403
            );
            assertErrID(body, ErrAuthNotClient);
        });

        test("DeleteProject", async () => {
            await expectForbidden(nestedUser.del(`/projects/${projectID}`));
        });
    });

    // --------------------------------------------------------------------------
    // CRYPTOGRAPHIC ISOLATION
    // --------------------------------------------------------------------------

    test("CryptographicIsolation", async () => {
        // Create a fresh owner and project for isolation test
        await createClient().withCredentials("crypto@mail.com", ValidPassword).register();
        const cryptoOwner = await createClient().withCredentials("crypto@mail.com", ValidPassword).login();
        const cryptoProject = await cryptoOwner.post("/projects", {
            project_name: "Crypto Project",
            metadata: { env: "test" },
        });
        const cryptoProjectID = cryptoProject.id;

        // Register and login a project user
        await createClient().http.post(`/projects/${cryptoProjectID}/register`, {
            email: "user@crypto.com",
            password: ValidPassword,
        });
        const cryptoUser = await createClient()
            .withCredentials("user@crypto.com", ValidPassword)
            .projectLogin(cryptoProjectID);

        const accessToken = cryptoUser.auth!.accessToken;

        // 1. Get Project JWKS
        const projectJwksRes = await cryptoOwner.http.get(
            `/projects/${cryptoProjectID}/.well-known/jwks.json`,
            { headers: { Cookie: `access_token=${cryptoOwner.auth!.accessToken}; refresh_token=${cryptoOwner.auth!.refreshToken}` } }
        );
        const projectXBase64: string = projectJwksRes.data.keys[0].x;
        const projectKeyBytes = decodeBase64Url(projectXBase64);
        verifyEd25519Key(projectXBase64);

        // 2. Get Global JWKS
        const globalJwksRes = await createClient().http.get("/.well-known/jwks.json");
        const globalXBase64: string = globalJwksRes.data.keys[0].x;
        const globalKeyBytes = decodeBase64Url(globalXBase64);
        verifyEd25519Key(globalXBase64);

        // 3. Token should verify with project key
        const projectPem = rawEd25519ToSPKIPem(projectKeyBytes);
        const projectKey = await importSPKI(projectPem, "EdDSA");
        await expect(
            jwtVerify(accessToken, projectKey)
        ).resolves.toBeDefined();

        // 4. Token must NOT verify with master key
        const masterPem = rawEd25519ToSPKIPem(globalKeyBytes);
        const masterKey = await importSPKI(masterPem, "EdDSA");
        await expect(
            jwtVerify(accessToken, masterKey)
        ).rejects.toThrow();
    });
});