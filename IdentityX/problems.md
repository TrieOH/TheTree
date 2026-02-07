# Security Assessment Findings

## Critical Severity


### 1. Missing Email Verification (Account Squatting)
*   **Severity**: Critical
*   **Affected components**: `internal/application/auth/usecase.go` (`Register`, `registerInternal`, `RegisterProjectUser`, `registerProjectUserInternal`)
*   **Description**: The `Register` and `RegisterProjectUser` functions (and their internal counterparts) create active user records and identities (`users.Register`, `projectUsers.Register`, `sessions.CreateIdentity`) *before* email verification. A verification email is sent, but the user account is created as active immediately.
*   **Attack scenario**: An attacker can register an account with *any* email address (e.g., `ceo@target.com`) without owning it. The account for `ceo@target.com` is created in the database. This allows the attacker to squat on legitimate user emails, preventing the real owner from registering. If the login/authorization flow does not strictly enforce `is_verified == true` (which is already a High-severity finding, "Active Status Ignored"), the attacker could potentially log in using the squatted email and a password they chose, leading to account takeover. Even with strict verification checks, the legitimate user is prevented from using the service under their own email.
*   **Why this is dangerous**: Leads to pre-auth account takeover, identity squatting, and denial of service for legitimate users by preventing them from registering with their own email.
*   **Conditions required to exploit**: Public registration endpoints and lack of immediate verification before account activation.
*   **Suggested mitigation (high-level, no patch required)**: Implement a "pending" state for new users. Do not create an active user record or grant an identity until the email verification process is successfully completed.

---

## High Severity

### 1. Active Status Ignored for Projects and Project Users
*   **Severity:** High
*   **Affected components:** `internal/application/auth/usecase.go` (`Login`, `LoginProjectUser`), `internal/database/queries/projects.sql` (`GetProjectByIDExternal`, `GetProjectByIDInternal`, `ListProjects`), `internal/database/queries/project_users.sql` (`GetProjectUserById`, `GetProjectUserByIdInternal`, `ListProjectUsersExternal`, `ListProjectUsersInternal`, `GetProjectUserByEmailExternal`, `GetProjectUserByEmailInternal`)
*   **Description:** The `projects` and `project_users` tables contain `is_active` flags, but these are never checked during authentication (Login), registration, or JWT verification. Specifically, the `Login` and `LoginProjectUser` functions in `internal/application/auth/usecase.go` retrieve the user but do *not* check `IsActive` after retrieval before proceeding to create a new session and issue tokens. Furthermore, critical database queries for retrieving projects and project users also lack `is_active` filtering.
*   **Attack scenario:** An administrator deactivates a project or a specific project user. However, the user can still log in, and their existing or new JWTs are still considered valid by the `resolvePublicKey` and `resolvePrivateKey` logic because the `WHERE` clauses do not filter for `is_active = true`.
*   **Why this is dangerous:** Revocation of access is broken. Deactivated accounts/projects remain fully functional.
*   **Suggested mitigation:** Add `AND is_active = true` to all SQL queries involved in authentication and key resolution. In the use case, explicitly check the `IsActive` field after retrieval.

### 2. Broken Validation for Advanced Field Types (Availability)
*   **Severity:** High
*   **Affected components:** `internal/application/auth/validate_field_type.go` (`validateFieldValue`)
*   **Description:** The `validateFieldValue` function explicitly contains a `FIXME` comment: `// FIXME: Implement other field types when they are implemented in the API. // Currently, types like 'email', 'select', 'radio', 'checkbox' will always return false, // making them unusable for project user registration.` Any attempt to use these common field types (`email`, `select`, `radio`, `checkbox`) will result in `validateFieldValue` returning `false` (because they fall through to the `default` case), making them unusable.
*   **Attack scenario:** If a project administrator configures a schema that includes field types like `email`, `select`, `radio`, or `checkbox` (which are valid types for custom user fields), all attempts by users to register for that project will fail due to validation errors. This causes a direct denial of service for new user registrations for affected projects.
*   **Why this is dangerous:** Leads to a critical availability issue for configured features, effectively causing a denial of service for certain project registrations.
*   **Conditions required to exploit:** A project attempts to use any custom field type other than `string`, `bool`, or `int`.
*   **Suggested mitigation:** Implement proper validation logic for `email` (e.g., regex), and `select`/`radio`/`checkbox` (e.g., check against allowed options).

### 3. Denial of Service via Uncached DB Lookups for JWT Keys
*   **Severity:** High
*   **Affected components:** `internal/application/keys/keys.go` (`VerifyProject`), `internal/application/tokens/verifier/usecase.go` (`verifyToken`)
*   **Description:** The `verifyToken` function delegates project-specific key verification to `keys.VerifyProject`, which in turn performs a database lookup (`uc.repo.GetProjectKeyByKID`) for every JWT verification with a "project:" prefixed `kid`. This lookup is not cached.
*   **Attack scenario:** An attacker can flood the API with requests containing JWTs (even invalid ones) with random `kid` headers like `project:{random-uuid}:v1`. This forces the server to perform a `SELECT` query on the `key_pair` table for every request, bypassing signature verification (since the key is needed *to* verify). Each such request, valid or not, incurs a database query overhead.
*   **Why this is dangerous:** Database resource exhaustion (CPU/IO), leading to Denial of Service for legitimate users.
*   **Conditions required to exploit**: Public network access.
*   **Suggested mitigation:** Implement an in-memory LRU cache for project public keys with a short TTL (e.g., 5-10 minutes) to reduce database load.

### 4. Lack of Secure Transport (SSL/TLS) Enforcement for Database Connection
**Severity**: High
**Affected components**: `internal/database/database.go` (`WaitForDB`), `DATABASE_URL`
*   **Description**: The `WaitForDB` function connects to the PostgreSQL database using a DSN from the `DATABASE_URL` environment variable. There is no explicit enforcement or check for `sslmode=require` or `sslmode=verify-full` in the connection string. Without strict SSL/TLS enforcement, the database connection could be established over an unencrypted channel, even if the database supports TLS.
*   **Attack scenario**: An attacker positioned on the network path between the application and the database (e.g., via a compromised internal network segment or by ARP spoofing) can intercept all traffic, including authentication credentials, sensitive data, and JWT master keys being transmitted. They can also perform Man-in-the-Middle attacks to alter data in transit.
*   **Why this is dangerous**: Confidentiality and integrity of all database communications are at risk, leading to data exfiltration, manipulation, and credential theft.
*   **Conditions required to exploit**: Network access to intercept traffic between the application and the database.
*   **Suggested mitigation (high-level, no patch required)**: Always enforce `sslmode=verify-full` in the `DATABASE_URL` DSN for production environments. The application should fail to start if a secure connection cannot be established.

### 6. Refresh Token Verification Bypasses Session Revocation
*   **Severity:** High
*   **Affected components:** `internal/application/tokens/verifier/usecase.go` (`VerifyRefreshToken`), `internal/application/auth/usecase.go` (`refreshInternal`)
*   **Description:** The `VerifyRefreshToken` function only checks the cryptographic validity and expiry of the refresh token itself. It does not check the status of the underlying session in the database. The `refreshInternal` function, which consumes the verified refresh token, checks for session revocation later, but this creates a race condition.
*   **Attack scenario:** An attacker obtains a valid, unexpired refresh token for a session that has been subsequently revoked (e.g., through a "logout everywhere" feature or admin action). The attacker can use this token to call the refresh endpoint. If they win the race condition against the database update that marks the session as revoked, `VerifyRefreshToken` will succeed, and they may be able to obtain a new access token before the `refreshInternal` function's own revocation check can stop them.
*   **Why this is dangerous:** It creates a window of opportunity for an attacker to re-establish an authenticated session using a token that should be invalid, undermining the session revocation mechanism.
*   **Suggested mitigation:** The `VerifyRefreshToken` function should not only verify the token's signature but also immediately check the status of the corresponding session in the database. It should return an error if the session is revoked or expired. This eliminates the race condition by ensuring that the token is invalid for all practical purposes as soon as the session is marked as revoked.

---

## Medium Severity

### 1. Project Enumeration via Registration Endpoint
*   **Severity:** Medium
*   **Affected components:** `internal/application/auth/usecase.go`, `internal/adapters/persistence/project_user_repo.go`
*   **Description:** The `RegisterProjectUser` endpoint (`/projects/{project_id}/register`) allows distinguishing between a non-existent `project_id` and an existing one based on the error response.
*   **Attack scenario:** An attacker can iterate through UUIDs or test specific IDs to see if they exist.
*   **Why this is dangerous**: Leaks system state and valid project identifiers.
*   **Conditions required to exploit**: Ability to call the public API.
*   **Suggested mitigation:** In the use case, check if the project exists *before* attempting insertion. If it doesn't exist, return a standard `404 Not Found`.

### 2. Lack of Account Lockout
*   **Severity:** Medium
*   **Affected components:** `internal/application/auth/usecase.go` (`Login`)
*   **Description:** The authentication logic currently lacks account lockout mechanisms (e.g., tracking failed login attempts and temporarily locking accounts after a certain threshold).
*   **Attack scenario:** An attacker can continuously attempt to brute-force user credentials without any countermeasure to slow them down or prevent further attempts on a specific account. This is particularly risky if a distributed (multi-IP) attack is used to bypass the global IP-based rate limiting.
*   **Why this is dangerous**: Increases susceptibility to brute-force attacks and credential stuffing.
*   **Conditions required to exploit**: Public login endpoints.
*   **Suggested mitigation:** Implement logic to track failed login attempts per user and, after a configured number of failures, temporarily lock the account. This typically involves adding fields to the user model (e.g., `failed_login_attempts`, `locked_until`).

### 3. Potential Information Leakage via Error Traces
*   **Severity:** Medium
*   **Affected components:** `internal/apierr/http_mapper.go`, `internal/adapters/http/middleware/error_to_response.go`
*   **Description:** The error handling middleware attaches the underlying error cause (including raw SQL errors) to the HTTP response via `AddTrace`.
*   **Attack scenario:** An attacker triggers a database error (e.g., FK violation) to see schema details like `violated unique constraint "users_email_key"`.
*   **Why this is dangerous**: Aids reconnaissance and SQL injection crafting.
*   **Suggested mitigation:** Configure error mapper to suppress traces in production.

### 4. Unbounded String Length in Custom Fields and Metadata (DoS)
*   **Severity:** Medium
*   **Affected components:** `internal/application/auth/validate_field_type.go`, `internal/application/auth/validate_custom_fields.go`, `internal/adapters/http/dto/project_requests.go`
*   **Description:** The application blindly accepts any string value for `field.String` and project `Metadata` without checking its length. These are stored in JSONB columns.
*   **Attack scenario:** An attacker registers a project user or creates a project with multi-megabyte strings in custom fields/metadata. Repeated requests can rapidly consume database storage and bandwidth.
*   **Why this is dangerous**: Unbounded input storage.
*   **Suggested mitigation:** Enforce a maximum length (e.g., 4KB) for all string inputs stored in the database.

### 5. Lack of Session Expiry Enforcement in Database (Cleanup)
*   **Severity:** Medium
*   **Affected components:** `internal/database/queries/sessions.sql`, `init.go` (Scheduler)
*   **Description:** While `RevokeExpiredSessions` exists and is scheduled, the accumulation of expired-but-not-deleted sessions could impact performance. The scheduler only runs daily.
*   **Why this is dangerous**: Performance degradation over time.
*   **Suggested mitigation:** Ensure the scheduler runs frequently enough.

### 7. Unrestricted Tenant Self-Registration (Public Admin Access)
*   **Severity:** Medium (Design Concern)
*   **Affected components:** `internal/adapters/http/router/router.go`, `internal/application/auth/usecase.go`
*   **Description:** The `/auth/register` endpoint is publicly accessible and allows anyone to create a "Client" account. Client accounts have the ability to create and manage projects (becoming Tenant Admins).
*   **Attack scenario:** If deployed in a private/enterprise context, an attacker can register an admin account and potentially exploit other vulnerabilities or abuse resource quotas (creating thousands of projects).
*   **Why this is dangerous**: Lack of access control for tenant creation.
*   **Suggested mitigation:** Disable public registration for Client accounts in production (feature flag), or implement an invitation-only flow for new Tenants.

### 8. User Enumeration via Timing Attack
*   **Severity:** Medium
*   **Affected components:** `internal/application/auth/usecase.go` (`Login`)
*   **Description:** The `Login` function returns immediately if `GetUserByEmail` fails (user not found), but executes a slow `bcrypt.CompareHashAndPassword` if the user exists. This creates a measurable timing difference.
*   **Attack scenario:** An attacker can send login requests with various email addresses and measure the response time. Significantly longer response times indicate the email is valid and registered in the system.
*   **Why this is dangerous**: Allows attackers to build a list of valid users for targeted phishing or credential stuffing.
*   **Suggested mitigation:** Implement a "fake" bcrypt comparison or constant-time logic when the user is not found to normalize response times.

### 9. Volatile Session Binding Information in JWT Claims
**Severity**: Medium
**Affected components**: `internal/application/tokens/issuer/usecase.go` (`NewAccessToken`)
**Description**: The `NewAccessToken` function includes volatile information such as `UserAgent` and `UserIP` directly within the `AccessClaims` of the JWT. JWTs are designed to be stateless and self-contained. Binding them to frequently changing attributes like IP address or User-Agent can lead to legitimate tokens being invalidated prematurely (e.g., if a user switches networks or updates their browser). If these fields are intended for security binding, their presence in a stateless token can create inconsistencies or be bypassed.
**Attack scenario**: A legitimate user's token becomes invalid due to a change in their IP address (e.g., roaming, VPN connection) or User-Agent, leading to a frustrating user experience and forced re-authentication. Alternatively, if these fields are not strictly enforced during verification, their presence offers a false sense of security for session binding. An attacker might attempt to replay a token while spoofing `UserIP` or `UserAgent` if these fields are used for binding and are not properly validated.
*   **Why this is dangerous**: Can degrade user experience through premature token invalidation or provide a false sense of security for session binding, which is better handled by refresh token rotation and explicit session management.
*   **Conditions required to exploit**: Frequent changes in user's IP/User-Agent, or an attacker successfully spoofing these headers if they are used for binding and are not properly validated.
*   **Suggested mitigation (high-level, no patch required)**: Avoid including highly volatile information like `UserAgent` and `UserIP` directly in JWT claims. Session binding should primarily rely on secure cookie flags, IP-based rate limiting, and robust refresh token rotation schemes. If binding to these attributes is necessary for specific security policies, ensure the verification logic is resilient to legitimate changes and prevents spoofing.

### 10. Incorrect Integer Range Validation for `float64` Inputs in Custom Fields
**Severity**: Medium
**Affected components**: `internal/application/auth/validate_field_type.go` (`validateFieldValue` for `field.Int` type)
**Description**: The `validateFieldValue` function, when validating `field.Int` types from `float64` inputs (which JSON numbers default to), uses `maxSafeInt64` (1<<53 - 1) and `minSafeInt64` (-(1 << 53)) for range checking. These bounds are significantly smaller than the actual `int64` range (approx. 2^63 - 1 to -(2^63)).
**Attack scenario**: Legitimate integer values provided as `float64` (common for JSON numbers) that fall between 2^53 and 2^63-1 (or their negative counterparts) will be incorrectly rejected by the validation logic. This can lead to legitimate user inputs being denied, causing functional errors or, in a malicious context, a denial of service by rejecting valid data.
*   **Why this is dangerous**: Leads to a functional bug where valid integer inputs are rejected, impairing legitimate application usage. Could be exploited to cause denial of service for users trying to register with specific (but otherwise valid) integer values in custom fields.
*   **Conditions required to exploit**: User inputting large integer values via JSON for `field.Int` types.
*   **Suggested mitigation (high-level, no patch required)**: Adjust the validation logic for `float64` to `int64` conversion to correctly use `math.MaxInt64` and `math.MinInt64`, potentially with a safer conversion method that accounts for float precision without prematurely rejecting valid `int64` values. If `json.Number` is available, convert to `int64` via that method first.

### 11. Unauthenticated Access to Project JWKS
*   **Severity:** Medium (Information Leakage)
*   **Affected components:** `internal/application/project/usecase.go` (`GetJWKS`)
*   **Description:** The `GetJWKS` function does not perform any authorization checks. It retrieves the public keys for a project based solely on the `projectID` provided in the request.
*   **Attack scenario:** An attacker with a valid account can iterate through UUIDs or use a known `projectID` to retrieve the public keys for any project on the platform. This leaks information about the existence of projects and their cryptographic materials.
*   **Why this is dangerous:** It breaks project isolation and allows attackers to gather information about other tenants.
*   **Suggested mitigation:** Before retrieving the keys, the `GetJWKS` function should verify that the authenticated principal has the right to access the specified project. This could be an ownership check or a check for project membership.

---

## Low / Informational

### 1. Container Runs as Root
*   **Severity:** Low
*   **Affected components:** `Dockerfile`
*   **Description:** The Docker image uses the default `alpine` user, which is `root`.
*   **Attack scenario:** If an attacker compromises the application, they gain root privileges inside the container.
*   **Why this is dangerous**: Violation of Principle of Least Privilege.
*   **Suggested mitigation:** Create and use a non-root user (e.g., `appuser`) in the Dockerfile.

### 2. Denial of Service via Unbounded User-Agent Storage
*   **Severity:** Low
*   **Affected components:** `internal/database/migrations/003_create_sessions.sql`, `internal/adapters/http/auth.go`
*   **Description:** The `sessions` table stores `user_agent` as `TEXT` without a length constraint.
*   **Attack scenario:** An attacker sends a Login request with a multi-megabyte User-Agent string to consume storage.
*   **Why this is dangerous**: Resource exhaustion.
*   **Suggested mitigation:** Truncate the User-Agent string to a reasonable length.

### 3. Potential Race Condition in Revocation Checks (Clock Skew)
*   **Severity:** Low
*   **Affected components:** `internal/database/migrations/003_create_sessions.sql`
*   **Description:** The session table has a constraint `CHECK (revoked_at <= NOW())`. Clock skew between app and DB could cause revocation to fail.
*   **Why this is dangerous**: Revocation reliability compromised.
*   **Suggested mitigation:** Use `revoked_at = NOW()` in the SQL instead of passing Go time.

### 4. Potential JWT Bloat via Custom Fields
*   **Severity:** Low
*   **Affected components:** `internal/application/auth/usecase.go` (Metadata passed to `NewProjectAccessToken`), `internal/application/tokens/issuer/usecase.go` (`NewProjectAccessToken`)
*   **Description:** Project user tokens include the full `Metadata` JSON in the claims (`internal/application/tokens/issuer/usecase.go#NewProjectAccessToken`).
*   **Why this is dangerous**: Tokens could exceed HTTP header size limits.
*   **Suggested mitigation:** Store a hash or fetch metadata from DB.

### 5. Log Spoofing via X-User-ID Header
*   **Severity:** Low
*   **Affected components:** `internal/adapters/http/middleware/logging.go`
*   **Description:** The `RequestID` middleware trusts the `X-User-ID` header from the client for logging.
*   **Attack scenario:** An attacker can send requests with a spoofed ID to confuse audit trails.
*   **Why this is dangerous**: Aids reconnaissance and SQL injection crafting.
*   **Suggested mitigation:** Log the authenticated Principal's ID, not the header.

### 6. Potential PII Leakage in Telemetry (Attributes)
*   **Severity:** Low
*   **Affected components:** `internal/application/auth/usecase.go`, `internal/adapters/persistence/project_repo.go`, etc.
*   **Description:** The application extensively annotates traces with attributes like `project.owner_id`, `user.id`.
*   **Why this is dangerous**: PII leakage to observability platforms.
*   **Suggested mitigation:** Ensure strict access control on traces or hash sensitive identifiers before exporting.

### 7. Lack of Email Normalization for Sub-addressing
*   **Severity:** Low
*   **Affected components:** `internal/application/auth/usecase.go` (`Register`, `Login`)
*   **Description:** Emails are normalized using `strings.ToLower`, but sub-addressing (e.g., `user+alias@gmail.com`) is not handled.
*   **Attack scenario:** A user can register multiple accounts using the same underlying email address, potentially bypassing uniqueness constraints or intended restrictions.
*   **Suggested mitigation:** Perform canonicalization of email addresses (e.g., remove `+` aliases for known providers like Gmail/Outlook) during registration and login.

### 8. Weak Password Policy Enforcement (Complexity)
*   **Severity:** Low
*   **Affected components:** `internal/adapters/http/dto/auth_requests.go`
*   **Description:** The `passwd` validator used in DTOs may not be properly registered or enforced, and the code only requires a minimum length of 8 characters.
*   **Attack scenario:** Users can use easily guessable passwords (e.g., `12345678`), increasing the risk of account takeover via credential stuffing or dictionary attacks.
*   **Suggested mitigation:** Ensure complexity requirements are enforced (uppercase, lowercase, digits, special characters) and verify the `passwd` validator's registration.

## Informational

### 1. Inaccurate `problems.md` References
**Severity**: Informational
**Affected components**: `problems.md`, `internal/application/auth/token_verifier.go` (non-existent), `internal/application/auth/token_issuer.go` (non-existent)
**Description**: The `problems.md` document contained multiple references to `internal/application/auth/token_verifier.go` and `internal/application/auth/token_issuer.go` as affected components for Critical and High severity findings. These files do not exist in the current codebase. The actual implementation for token verification is located at `internal/application/tokens/verifier/usecase.go` and for token issuance at `internal/application/tokens/issuer/usecase.go`.
*   **Why this is dangerous**: Outdated documentation can lead to confusion during security reviews, make it harder to track and remediate vulnerabilities, and give a false sense of security regarding fixed issues or non-existent code paths.
*   **Conditions required to exploit**: None, this is a documentation issue.
*   **Suggested mitigation (high-level, no patch required)**: Update `problems.md` to reflect the correct file paths for affected components. Regularly review and synchronize security documentation with the evolving codebase.

## FIXED

### 1. Baked-in Private Keys in Docker Image
*   **Severity:** Critical
*   **Affected components:** `Dockerfile`, `internal/crypto/ed25519.go` (`GenerateEd25519` function), `init.go`
*   **How it was fixed:** The key generation is no longer part of the `Dockerfile` build process. Instead, `init.go` checks for an existing signing key in the database on application startup. If not found, a new key pair is generated using `internal/crypto/ed25519.go#GenerateEd25519` and stored securely in the database.
*   **What changed:** The key management shifted from a static, build-time process to a dynamic, runtime, and database-persisted one, with automated key rotation. `internal/utils/keys.go` (originally mentioned in `problems.md`) does not exist; the relevant code is `internal/crypto/ed25519.go`.
*   **Why the fix is sufficient:** Keys are no longer exposed in the Docker image layers. They are generated on first run, persisted in the database, and automatically rotated, significantly improving security posture regarding static key exposure.
*   **Any residual risk or assumptions:** Assumes the database itself is secure and the `JWT_MASTER_KEY` is managed securely outside of the application's environment variables (which is a separate problem #3 and #4 in this document).

### 2. Inadequate Validation and Secure Sourcing of ED25519 Keys
*   **Severity:** Critical
*   **Affected components:** `init.go`, `internal/crypto/ed25519.go` (`GenerateEd25519` function), `key_pair` table
*   **How it was fixed:** The key management system was refactored. Keys are no longer loaded from environment variables (`JWT_PRIVATE_KEY`, `JWT_PUBLIC_KEY`) or from files. Instead, `init.go` generates new keys using `internal/crypto/ed25519.go#GenerateEd25519` (which uses a cryptographically secure random number generator) if they do not exist in the database. These keys are then stored in the `key_pair` table.
*   **What changed:** The sourcing of keys shifted from potentially insecure environment variables/files to runtime generation and secure database storage. The generation process itself uses strong cryptographic primitives.
*   **Why the fix is sufficient:** This mitigates risks associated with weak or guessable keys, and keys loaded from insecure file system locations. The keys are now securely generated and managed within the application's persistent storage (the database).
*   **Any residual risk or assumptions:** The security of the keys now relies heavily on the database's security and the secure management of the `JWT_MASTER_KEY` (addressed in issues #2 and #3).

### 3. Global Exposure of ED25519 Signing Keys
*   **Severity:** High
*   **Affected components:** `internal/application/keys/keys.go` (`SignGoAuth`, `SignProject`), `init.go`
*   **How it was fixed:** The key management system was refactored. Private keys are no longer stored in global variables. Instead, they are retrieved from the database, used for the signing operation within the scope of a single function, and then zeroed out from memory using `defer zero(priv)`.
*   **What changed:** The lifecycle of private keys in memory is now short-lived and function-scoped, which significantly reduces the attack surface for accidental leakage via logs, memory dumps, or other side-channels.
*   **Why the fix is sufficient:** This mitigates the risk of global key exposure by adopting a much more secure key handling pattern.
*   **Any residual risk or assumptions:** None beyond the security of the Go runtime and the underlying operating system.

### 4. Unversioned Project Public Keys Impair Key Rotation
*   **Severity:** Medium
*   **Affected components:** `init.go` (`tryRotateProjectKeys`), `internal/database/queries/keys.sql` (`RotateSigningKeysForProject`, `ListActivePublicKeysForProject`)
*   **How it was fixed:** The system now supports key versioning and rotation. The `kid` for project keys includes a unique `ulid` component, allowing for multiple keys per project. The `tryRotateProjectKeys` function in `init.go` periodically checks for expiring keys, and the `RotateSigningKeysForProject` query transitions the old key to a `rotated` state, allowing it to be used for verification but not for signing. The JWKS endpoint correctly serves both `active` and `rotated` keys.
*   **What changed:** A complete key rotation mechanism has been implemented, allowing for graceful key migration and reducing the impact of a key compromise.
*   **Why the fix is sufficient:** The system can now securely rotate keys without invalidating existing tokens, which is a fundamental security control.
*   **Any residual risk or assumptions:** The rotation period (currently 7 days for new keys) should be reviewed to ensure it aligns with the project's security requirements.

### 5. Project Private Keys Encrypted via DB Session Key (SQLi Risk)
*   **Severity:** High (Critical context)
*   **Affected components:** `internal/database/database.go` (`SetJWTMasterKey`), `internal/database/migrations/009_create_keys.sql` (`key_pair` table, `private_key` column)
*   **How it was fixed:** The `SetJWTMasterKey` function has been removed from `internal/database/database.go`. Key encryption/decryption is now handled at the application layer, using a key managed outside the database session. This eliminates the database session as a vector for key exfiltration via SQL injection.
*   **What changed:** The responsibility for encryption of private keys shifted from the database (using session variables) to the application layer.
*   **Why the fix is sufficient:** By performing encryption/decryption in the application, the `JWT_MASTER_KEY` is no longer exposed to the database session, preventing attackers with SQL injection capabilities from easily decrypting private keys using database functions.
*   **Any residual risk or assumptions:** Assumes the application-layer key management is secure.

### 6. Direct SQL Injection Vulnerability in `SetJWTMasterKey` via Insufficient Escaping
*   **Severity**: Critical
*   **Affected components**: `internal/database/database.go` (`SetJWTMasterKey`), `init.go`
*   **How it was fixed:** The `SetJWTMasterKey` function, which was susceptible to SQL injection due to manual string escaping, has been entirely removed from the codebase.
*   **What changed:** The vulnerable code path has been eliminated.
*   **Why the fix is sufficient:** The direct SQL injection vulnerability no longer exists as the function that introduced it has been removed.
*   **Any residual risk or assumptions:** None directly related to this specific vulnerability.

### 7. Misleading CORS Log & Potential Insecure Default
*   **Severity:** Medium
*   **Affected components:** `internal/adapters/http/router/router.go` (`GetCORSOptions`)
*   **How it was fixed:** The `GetCORSOptions` function now explicitly checks if `CORS_ALLOWED_ORIGINS` is set. If it's `nil` (meaning the environment variable was not configured), the application will `log.Fatalf`, forcing a secure configuration. This prevents the `go-chi/cors` library from potentially defaulting to an insecure `*` (all origins) policy.
*   **What changed:** The application now fails fast on missing CORS origin configuration, removing the ambiguity and forcing administrators to make an explicit, secure choice.
*   **Why the fix is sufficient:** This eliminates the false sense of security and ensures that CORS is either explicitly configured or the application does not start, preventing unintended wide-open access.
*   **Any residual risk or assumptions:** None, provided the administrator configures `CORS_ALLOWED_ORIGINS` appropriately for their environment.

### 8. Unbounded Request Body Size (DoS)
*   **Severity:** Medium
*   **Affected components:** `internal/adapters/http/router/router.go` (`CreateRouter` middleware stack)
*   **How it was fixed:** A `middleware.MaxBodySize(1 << 20)` (1 MB) middleware was added to the `CreateRouter` function in `internal/adapters/http/router/router.go`. This explicitly limits the maximum size of incoming request bodies.
*   **What changed:** The application now enforces a maximum request body size, preventing attackers from sending excessively large payloads.
*   **Why the fix is sufficient:** This mitigation directly addresses the Denial of Service vulnerability by limiting resource consumption associated with large request bodies.
*   **Any residual risk or assumptions:** The 1MB limit is a reasonable default, but could be configurable if specific use cases require larger payloads.

### 9. Lack of Global Rate Limiting
*   **Severity:** Medium
*   **Affected components:** `internal/adapters/http/router/router.go` (`CreateRouter` middleware stack)
*   **How it was fixed:** The `CreateRouter` function in `internal/adapters/http/router/router.go` now implements an IP-based rate-limiting middleware (`httprate.Limit`) with a default of 400 requests per minute. This mitigates brute-force and resource exhaustion attacks from a single IP address.
*   **What changed:** A global, IP-based rate limit has been applied to all incoming requests.
*   **Why the fix is sufficient:** This provides a strong first layer of defense against volumetric attacks and simple brute-force attempts.
*   **Any residual risk or assumptions:** This does not prevent distributed (multi-IP) brute-force attacks against a single account. Account-level lockout logic is still required for that, as tracked in a separate issue.

### 10. Logic Flaw in Session Listing
*   **Severity:** Low
*   **Affected components:** `internal/database/queries/sessions.sql` (`ListSessions`)
*   **How it was fixed:** The original `ListUserSessions` query was replaced with a `ListSessions` query. This new query correctly joins `sessions` with `session_identities` and filters on the identity `type` (which represents the `user_type`).
*   **What changed:** The query now correctly constrains session lookups by both the entity ID and the user type, preventing collisions between different user types (e.g., `client` vs. `project`) that might share a UUID.
*   **Why the fix is sufficient:** This eliminates the risk of a user from one context (e.g., a client) accidentally viewing or affecting the sessions of a user from another context (e.g., a project user).
*   **Any residual risk or assumptions:** None. The fix is correct.

### 11. Mismatched Key ID (kid) in JWKS vs Tokens
*   **Severity:** Low (Interop/Usability)
*   **Affected components:** `internal/domain/key/key.go`, `internal/application/tokens/issuer/usecase.go`, `internal/adapters/persistence/key_repo.go`, `internal/database/queries/keys.sql`
*   **How it was fixed:** The system was refactored to ensure consistent `kid` values across token issuance and JWKS publication. The `PublicKeyToJWK` function in `internal/domain/key/key.go` now directly uses the `kid` from the `PublicKey` struct. Token issuance (`internal/application/tokens/issuer/usecase.go`) retrieves the active `kid` dynamically from the database via the `keys` service. The `kid` values are managed in the `key_pair` table, ensuring that the `kid` embedded in issued tokens always matches the `kid` published in the JWKS.
*   **What changed:** Hardcoded `kid` values have been replaced with dynamic retrieval from the database, ensuring cryptographic key identifiers are consistent.
*   **Why the fix is sufficient:** This resolves the interoperability issue by ensuring that clients can correctly use the JWKS to verify tokens, as the `kid` in the token will always map to a public key in the JWKS.
*   **Any residual risk or assumptions:** None, as long as the key management system correctly updates and stores `kid` values.

### 12. Incomplete Session Revocation Logic (Collision Risk)
*   **Severity:** Medium
*   **Affected components:** `internal/database/queries/sessions.sql` (`RevokeOtherSessions`, `RevokeAllSessions`)
*   **How it was fixed:** The `RevokeOtherSessions` and `RevokeAllSessions` queries in `internal/database/queries/sessions.sql` were modified to explicitly join with `session_identities` and filter by `i.type = $1`. This ensures that revocation operations are constrained by the `user_type` in addition to the `user_id`.
*   **What changed:** Session revocation queries now prevent cross-account denial-of-service by ensuring that sessions are only revoked for the correct user type.
*   **Why the fix is sufficient:** This eliminates the risk of UUID collisions between different user types (client vs. project) leading to unintended session revocations.
*   **Any residual risk or assumptions:** None.

### 13. Missing Refresh Token Reuse Detection
*   **Severity:** Medium
*   **Affected components:** `internal/application/auth/usecase.go` (`refreshInternal`), `internal/adapters/persistence/session_repo.go`
*   **How it was fixed:** In `internal/application/auth/usecase.go#refreshInternal`, the system actively detects refresh token reuse (`if sess.TokenID != oldJTI`). Upon detection, the entire refresh token family is immediately revoked by calling `sessions.MarkRevokedByFamilyID(ctx, sess.FamilyID)`.
*   **What changed:** The core security vulnerability related to refresh token reuse has been mitigated. An attacker cannot use a stolen or replayed refresh token to obtain a new access token, as all related sessions are immediately invalidated upon first detection of reuse.
*   **Why the fix is sufficient:** This directly addresses the attack scenario by preventing an attacker from re-establishing an authenticated session with a replayed refresh token. While a `FIXME` comment related to auditing remains, the critical security flaw of allowing continued access via reused tokens is resolved. Enhanced auditing remains a valuable improvement for incident response but is not essential for preventing the immediate security compromise.
*   **Any residual risk or assumptions:** The `FIXME` comment indicates that robust auditing and incident response mechanisms for such events are still areas for potential enhancement but the primary security concern of preventing session re-establishment is handled.

### 14. JWT Algorithm Confusion (EdDSA -> HS256)
*   **Severity:** Critical
*   **Affected components:** `internal/application/tokens/verifier/usecase.go` (`verifyToken` function)
*   **How it was fixed:** The `verifyToken` function now explicitly validates that the `alg` header and the `token.Method.Alg()` match the expected `EdDSA` algorithm before proceeding with key verification.
*   **What changed:** Added explicit checks: `if alg != jwt.SigningMethodEdDSA.Alg()` and `if token.Method.Alg() != jwt.SigningMethodEdDSA.Alg()`.
*   **Why the fix is sufficient:** This prevents the attack where an attacker forces the server to use `HS256` (HMAC) verification using the public key as the secret.
*   **Any residual risk or assumptions:** None.

### 15. Session Revocation Race Condition (Potential)
*   **Severity**: High
*   **Affected components**: `init.go`, `gocron.NewJob`, `queries.RevokeExpiredSessions`
*   **How it was fixed:** The `refreshInternal` function in `internal/application/auth/usecase.go` explicitly checks if the session is expired (`sess.ExpiresAt.Before(now)`) or revoked (`sess.RevokedAt != nil`) at the point of use. Additionally, the cleanup job runs hourly instead of daily.
*   **What changed:** Application-level validation ensures that expired sessions are rejected immediately during the refresh flow, independent of the background cleanup job.
*   **Why the fix is sufficient:** This eliminates the TOCTOU window where an expired session could be used before the cleanup job runs. The background job is now purely for database hygiene, not security enforcement.
*   **Any residual risk or assumptions:** None.

### 16. Missing Password Reset Functionality
*   **Severity:** Medium (Availability/Process)
*   **Affected components:** `internal/ports/inbounds/auth_service_interface.go`, `internal/application/auth/usecase.go`
*   **How it was fixed:** Implemented a full password reset flow including `ForgotPassword` (issuing signed JWT tokens via email) and `ResetPassword` (verifying tokens, updating passwords, and invalidating all active sessions).
*   **What changed:** Users can now securely recover their accounts. The implementation includes protection against token reuse and ensures session invalidation upon password change.
*   **Why the fix is sufficient:** The standard self-service recovery path is now available and follows security best practices (signed tokens, atomicity via transactions, session revocation).