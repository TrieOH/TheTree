# Performance Optimizations

## 1. Auth Middleware Session Caching

- **Affected Components**: `AuthMiddleware`, `SessionRepository`, Postgres
- **Current Behavior**: Every authenticated request triggers a DB query `GetUserSessionByTokenID`.
- **Inefficiency**: 
    - Adds ~1-5ms (best case) to *every* request latency.
    - Scales linearly with RPS, putting significant load on the primary DB for read-heavy workloads.
    - Redundant: Session validity rarely changes within seconds.
- **Evidence**: `internal/adapters/http/middleware/auth.go` calls `mw.sessions.GetByTokenID` unconditionally.
- **Expected Improvement**: 
    - Reduce DB ops for auth by 95%+ (depending on cache TTL).
    - Reduce latency by removing DB round-trip.
- **Tradeoffs**: 
    - Cache invalidation complexity (revocation must clear cache).
    - Slight consistency delay if using TTL-only (revocation might take `TTL` to propagate).
- **Suggested Approach**: 
    - Introduce a `SessionCache` (Redis/In-memory).
    - Cache key: `session:valid:{token_id}`.
    - Write-through or invalidation on `Revoke...` calls.
    - **Crucial**: Ensure graceful fallback to DB if cache is unavailable to avoid single point of failure.

## 2. Project Public Key Caching

- **Affected Components**: `TokenVerifier`, `ProjectRepository`
- **Current Behavior**: Project-scoped tokens (Project Users) trigger a DB lookup for the project's public key during signature verification. This happens twice per request (Access + Refresh token).
- **Inefficiency**: 
    - Public keys are immutable or change very rarely.
    - Fetching them from DB on every request is wasteful.
    - Adds 2 extra DB round-trips for project user requests.
- **Evidence**: `internal/application/auth/token_verifier.go` calls `uc.projects.GetPublicKeyByID` inside `resolvePublicKey`.
- **Expected Improvement**: 
    - Eliminate 2 DB calls per project-user request.
- **Tradeoffs**: 
    - Memory usage for cache.
    - Need invalidation if key rotation is implemented (currently TODO).
- **Suggested Approach**: 
    - In-memory LRU cache for Project Public Keys with a reasonable TTL (e.g., 5-10 minutes) or until service restart (since keys don't rotate yet).

## 3. Offload Encryption to Application

- **Affected Components**: `ProjectRepository`, Postgres
- **Current Behavior**: Private keys are encrypted/decrypted using `pgp_sym_encrypt/decrypt` within Postgres.
- **Inefficiency**: 
    - Consumes DB CPU, which is the hardest resource to scale.
    - Requires setting a session variable `app.jwt_master_key` which has security implications (SQLi).
- **Evidence**: `internal/database/queries/projects.sql` uses `pgp_sym_decrypt`.
- **Expected Improvement**: 
    - Reduced DB CPU usage.
    - Improved security posture (key not available in SQL session).
- **Tradeoffs**: 
    - Application must manage encryption (Go is fast at this).
    - Data migration required to re-encrypt existing keys (if format changes).
- **Suggested Approach**: 
    - Use Go's `crypto/cipher` (AES-GCM or NaCl) to encrypt/decrypt keys.
    - Store as `BYTEA` in DB.

## 4. Optimize Session Query Columns

- **Affected Components**: `SessionRepository`, `AuthMiddleware`
- **Current Behavior**: `GetUserSessionByTokenID` selects `*` (all columns), including `user_agent` (text) and `user_ip`.
- **Inefficiency**: 
    - Fetches unnecessary data for the validation check (only `revoked_at`, `expires_at`, `session_id`, `user_id`, `project_id` are strictly needed).
    - Increases network serialization/deserialization cost slightly.
- **Evidence**: `internal/database/queries/sessions.sql` does `SELECT *`.
- **Expected Improvement**: 
    - Minor reduction in network bandwidth and memory allocs.
- **Tradeoffs**: 
    - Requires maintaining a separate query/struct for validation.
- **Suggested Approach**: 
    - Create `GetUserSessionValidity` query selecting only necessary fields.

## 5. Remove Unnecessary Transactions

- **Affected Components**: `AuthService.Register`, `AuthService.RegisterProjectUser`, `TxRunner`
- **Current Behavior**: `Register` wraps a single `INSERT` statement in a `BeginTx` ... `Commit` block.
- **Inefficiency**: 
    - `BeginTx` and `Commit` add 2 extra round-trips to the database.
    - Increases connection hold time.
    - Unnecessary since a single SQL statement is implicitly atomic.
- **Evidence**: `internal/application/auth/usecase.go` calls `uc.tx.WithinTx` for `Register`.
- **Expected Improvement**: 
    - Reduced latency (minus ~2 RTTs) for registration.
    - Reduced DB connection contention.
- **Tradeoffs**: 
    - None for single-statement operations.
- **Suggested Approach**: 
    - Remove `WithinTx` wrapper for single-operation use cases.

## 6. Pagination for List Endpoints

- **Affected Components**: `ProjectRepository`, `SessionRepository`, `ProjectHandler`, `SessionHandler`
- **Current Behavior**: `ListProjects` and `ListUserSessions` fetch *all* records matching the user ID.
- **Inefficiency**: 
    - **Unbounded Query**: If a user has thousands of projects/sessions, this consumes excessive memory and CPU (app & DB) and network bandwidth.
    - Vulnerable to DoS (Denial of Service) via resource exhaustion.
- **Evidence**: `internal/database/queries/projects.sql` lacks `LIMIT` and `OFFSET`.
- **Expected Improvement**: 
    - Constant memory usage regardless of dataset size.
    - Predictable latency.
- **Tradeoffs**: 
    - API breaking change (response format might need to include pagination metadata like `next_page`, `total`).
    - Client complexity (must handle pages).
- **Suggested Approach**: 
    - Add `limit` and `offset` (or cursor) parameters to `List` queries and API endpoints.
    - Enforce a hard maximum limit (e.g., 100).

## 7. HTTP Response Compression (Gzip)

- **Affected Components**: `Router`, `Chi Middleware`
- **Current Behavior**: HTTP responses (JSON) are sent uncompressed.
- **Inefficiency**: 
    - JSON is text-heavy and compresses very well (often 70-90% reduction).
    - Uncompressed responses waste network bandwidth and increase latency for clients on slow networks.
- **Evidence**: `internal/adapters/http/router/router.go` does not include `middleware.Compress`.
- **Expected Improvement**: 
    - Significantly reduced payload size.
    - Faster transfer times.
- **Tradeoffs**: 
    - Minor CPU increase for compression (usually negligible compared to network gains).
- **Suggested Approach**: 
    - Add `r.Use(middleware.Compress(5))` to the middleware stack.

## 8. Migrate to pgx Driver

- **Affected Components**: `Database`, `sqlc`
- **Current Behavior**: The application uses `github.com/lib/pq` (maintenance mode).
- **Inefficiency**: 
    - `lib/pq` is slower and has fewer features than `jackc/pgx`.
    - Parsing and serialization in `pq` are less optimized.
- **Evidence**: `go.mod` and `router.go` import `github.com/lib/pq`.
- **Expected Improvement**: 
    - Better performance (throughput and memory allocation).
    - Access to modern Postgres features (binary protocol for more types).
- **Tradeoffs**: 
    - Significant refactoring effort (changing driver imports, potentially error handling).
    - **Note**: Switching to `pgxpool` native interface requires refactoring `TxRunner` and all Repositories to use `pgx.Tx` instead of `*sql.Tx`. Using `pgx` via `stdlib` is a safer first step.
- **Suggested Approach**: 
    - Switch `sqlc` configuration to use `pgx/v5`.
    - Update `database/database.go` to use `pgxpool` or `stdlib` wrapper.

## 9. Optimize GetVerbose Field Assembly

- **Affected Components**: `SchemaService.GetVerbose`
- **Current Behavior**: Iterates through all versions, and for each version, iterates through *all* fields of the schema to find matching fields.
- **Inefficiency**: 
    - **O(V * F)** complexity, where V is versions and F is total fields.
    - If a schema has many versions, this quadratic-like behavior wastes CPU.
- **Evidence**: `internal/application/schema/usecase.go` nested loops.
- **Expected Improvement**: 
    - **O(F + V)** complexity.
- **Tradeoffs**: 
    - None.
- **Suggested Approach**: 
    - Pre-group fields by `SchemaVersionID` into a map: `map[uuid.UUID][]Field`.
    - Iterate versions and pick from map.

## 10. GetVerbose Endpoint Scalability

- **Affected Components**: `SchemaService.GetVerbose`, `SchemaRepository`, `FieldsRepository`
- **Current Behavior**: Fetches *all* versions and *all* fields for a schema.
- **Inefficiency**: 
    - Unbounded payload growth as schema history grows.
    - Heavy DB query (`ListFieldsFromSchema` returns all history).
- **Evidence**: `internal/application/schema/usecase.go` calls `versions.List` and `fields.List`.
- **Expected Improvement**: 
    - Constant payload size.
- **Tradeoffs**: 
    - Changing API behavior (pagination or limit to latest N versions).
- **Suggested Approach**: 
    - default `GetVerbose` to return only the `current` and `draft` versions.
    - Add `?include_history=true` or pagination for versions.

## 11. Batch Data Cleanup Jobs

- **Affected Components**: `SessionCleanup` Cron Job
- **Current Behavior**: `RevokeExpiredSessions` and `DeleteRevokedSessions` operate on the entire table (all matching rows).
- **Inefficiency**: 
    - Can cause massive transaction logs (WAL) and lock the table if millions of rows are deleted at once.
    - Spikes in DB load daily.
- **Evidence**: `init.go` schedules cleanup daily without batching.
- **Expected Improvement**: 
    - Smooth DB load profile.
    - Prevent transaction timeouts.
- **Tradeoffs**: 
    - Slightly more complex job logic.
- **Suggested Approach**: 
    - Update cleanup queries to use `LIMIT` and run in a loop until no more rows to delete.
    - Run job more frequently (e.g., hourly) to reduce batch size.

## 12. Separate Migrations from Startup

- **Affected Components**: `init.go`, `database/database.go`
- **Current Behavior**: `init()` calls `database.RunMigrations` which applies schema changes on every startup.
- **Inefficiency**: 
    - Slows down application startup (Cold Start).
    - Risky in distributed environments (concurrent migration attempts).
    - Requires application pods to have `ALTER TABLE` permissions.
- **Evidence**: `init.go` calls `RunMigrations`.
- **Expected Improvement**: 
    - Faster startup.
    - Safer deployment process.
- **Tradeoffs**: 
    - Requires operational change (run migration job before deployment).
- **Suggested Approach**: 
    - Remove `RunMigrations` from `init()`.
    - Create a separate entrypoint or command for migrations.

## 13. Limit Request Body Size

- **Affected Components**: `validation/validation.go`
- **Current Behavior**: `ValidateInto` reads the entire request body via `json.Decode` without an explicit size limit.
- **Inefficiency**: 
    - Vulnerable to DoS (Memory Exhaustion) via large payloads.
    - Wastes CPU parsing junk data.
- **Evidence**: `internal/adapters/http/validation/validation.go` uses `json.NewDecoder(r.Body)`.
- **Expected Improvement**: 
    - improved stability and security.
- **Tradeoffs**: 
    - None (legitimate requests are small).
- **Suggested Approach**: 
    - Wrap `r.Body` with `http.MaxBytesReader(w, r.Body, Limit)`.

## 14. Optimize Project Users List Index

- **Affected Components**: `ProjectUserRepository`, Postgres
- **Current Behavior**: `ListProjectUsersInternal` filters by `project_id` and sorts by `created_at DESC`.
- **Inefficiency**: 
    - Uses `idx_project_users_project_id` but requires a sort operation.
- **Evidence**: `internal/database/queries/project_users.sql` uses `ORDER BY created_at DESC`.
- **Expected Improvement**: 
    - Index-only scan or pre-sorted index scan.
- **Tradeoffs**: 
    - Index maintenance cost.
- **Suggested Approach**: 
    - Create `CREATE INDEX idx_project_users_pid_created ON project_users (project_id, created_at DESC)`.

## 15. Add Missing Index on projects(owner_id)

- **Affected Components**: `ProjectRepository`, Postgres
- **Current Behavior**: `ListProjects` queries `projects` by `owner_id`.
- **Inefficiency**: 
    - No index on `owner_id` causes a Sequential Scan on the `projects` table.
    - As `projects` grows, this endpoint becomes slower linearly.
- **Evidence**: `internal/database/migrations/004_create_projects.sql` does not create an index on `owner_id`.
- **Expected Improvement**: 
    - O(log N) lookup instead of O(N).
- **Tradeoffs**: 
    - Index maintenance cost.
- **Suggested Approach**: 
    - `CREATE INDEX idx_projects_owner_id ON projects(owner_id);`

## 16. Pin Tool Versions in Dockerfile

- **Affected Components**: `Dockerfile`
- **Current Behavior**: Uses `go install ...@latest`.
- **Inefficiency**: 
    - Non-deterministic builds.
    - Potential future breakages if tools update.
- **Evidence**: `Dockerfile` lines 5-7.
- **Expected Improvement**: 
    - Reproducible builds.
- **Tradeoffs**: 
    - Maintenance (manual updates).
- **Suggested Approach**: 
    - Pin to specific versions (e.g. `goose@v3.15.0`).

## 17. Optimize Docker Build (Binaries vs Source)

- **Affected Components**: `Dockerfile`
- **Current Behavior**: Compiles `sqlc`, `goose`, `swag` from source using `go install`.
- **Inefficiency**: 
    - Significantly slower build time (compiling these tools takes minutes).
    - Increases CI/CD cost.
- **Evidence**: `Dockerfile` uses `go install`.
- **Expected Improvement**: 
    - Faster build times.
- **Tradeoffs**: 
    - Architecture dependency (need to fetch correct binary for linux/amd64 or arm64).
- **Suggested Approach**: 
    - Download pre-compiled binaries from GitHub Releases via `wget/curl`.

## 18. Remove Baked-in Keys from Docker Image

- **Affected Components**: `Dockerfile`
- **Current Behavior**: Generates Ed25519 keys inside the image build process if missing.
- **Inefficiency**: 
    - **Security Risk**: Keys are part of the image layer.
    - **Operational Issue**: Rebuilding the image rotates the keys, invalidating all tokens.
- **Evidence**: `Dockerfile` runs `openssl genpkey` and `COPY`s keys.
- **Expected Improvement**: 
    - Stateless images.
    - Persistent identity across deployments.
- **Tradeoffs**: 
    - Requires external secret management (K8s secrets, Vault, or Env Vars).
- **Suggested Approach**: 
    - Remove key generation from Dockerfile.
    - Fail startup if keys are missing from volume/env.

## 19. Pin Database Image Version

- **Affected Components**: `docker-compose.yml`
- **Current Behavior**: `postgres:18` (or generic latest).
- **Inefficiency**: 
    - Unpredictable upgrades.
    - Postgres 18 is not stable (as of typical context).
- **Evidence**: `docker-compose.yml` uses `postgres:18`.
- **Expected Improvement**: 
    - Stability.
- **Tradeoffs**: 
    - None.
- **Suggested Approach**: 
    - Use a stable tag like `postgres:16-alpine`.

## 20. Support Key Loading from Environment Variables

- **Affected Components**: `utils/keys.go`, `init.go`
- **Current Behavior**: `LoadEd25519Keys` expects file paths.
- **Inefficiency**: 
    - Forces saving secrets to disk in containerized environments, which is less secure and less convenient than direct Env Var injection.
- **Evidence**: `internal/utils/keys.go` uses `os.ReadFile`.
- **Expected Improvement**: 
    - Cloud-native best practice compliance.
    - No sensitive files on disk.
- **Tradeoffs**: 
    - Code change in `utils`.
- **Suggested Approach**: 
    - Modify `LoadEd25519Keys` to detect if input is a PEM block vs a Path, or add `LoadEd25519KeysFromContent`.

## 21. Cursor-Based Pagination

- **Affected Components**: `ProjectRepository`, `SessionRepository`
- **Current Behavior**: `List` endpoints use implicit "offset" logic (currently unbounded, but even with `OFFSET` it's slow deep in).
- **Inefficiency**: 
    - `OFFSET N` scans and drops N rows. O(N).
- **Evidence**: `ListProjects` queries by `owner_id` ordered by time.
- **Expected Improvement**: 
    - O(1) fetch time regardless of depth.
- **Tradeoffs**: 
    - API complexity (client must send cursor).
- **Suggested Approach**: 
    - Implement keyset pagination (`WHERE (created_at, id) < (cursor_time, cursor_id) LIMIT N`) to handle timestamp collisions.

## 22. High-Performance JSON Library

- **Affected Components**: `Validation`, `Handlers`
- **Current Behavior**: Uses standard `encoding/json`.
- **Inefficiency**: 
    - Reflection-based, relatively slow for high-throughput services.
- **Evidence**: Imports `encoding/json`.
- **Expected Improvement**: 
    - 2x-5x faster JSON serialization/deserialization.
    - Lower memory allocation.
- **Tradeoffs**: 
    - External dependency.
- **Suggested Approach**: 
    - Replace with `github.com/goccy/go-json` or `github.com/segmentio/encoding`.

## 23. Configure HTTP Server Timeouts

- **Affected Components**: `main.go`
- **Current Behavior**: Uses `http.ListenAndServe` with default `http.Server` configuration (unlimited timeouts).
- **Inefficiency**: 
    - Vulnerable to Slowloris attacks and resource leaks from hung connections.
- **Evidence**: `main.go` calls `http.ListenAndServe`.
- **Expected Improvement**: 
    - Robustness against slow clients.
- **Tradeoffs**: 
    - Legitimate slow uploads might be cut off (configure appropriately).
- **Suggested Approach**: 
    - Instantiate `http.Server` with `ReadTimeout`, `WriteTimeout`, `IdleTimeout`, `ReadHeaderTimeout`.

## 24. Configure Database Connection Pool

- **Affected Components**: `database/database.go`
- **Current Behavior**: Uses default `sql.Open` settings (Unlimited Open, 2 Idle).
- **Inefficiency**: 
    - Risk of DB overload (too many connections).
    - High connection churn (too few idle connections).
- **Evidence**: `internal/database/database.go` calls `sql.Open` without setting pool stats.
- **Expected Improvement**: 
    - Predictable DB load.
    - Lower latency for burst traffic (reusing idle conns).
- **Tradeoffs**: 
    - Tuning required based on deployment resources.
- **Suggested Approach**: 
    - `db.SetMaxOpenConns(N)`
    - `db.SetMaxIdleConns(M)`
    - `db.SetConnMaxLifetime(Duration)`

## 25. Tune OTel SQL Tracing

- **Affected Components**: `database/database.go`
- **Current Behavior**: Uses `otelsql.WithSQLCommenter(true)`.
- **Inefficiency**: 
    - Adds tracing overhead to every query.
    - SQL Commenter increases query size and network traffic slightly.
- **Evidence**: `internal/database/database.go` registers `otelsql` driver.
- **Expected Improvement**: 
    - Reduced CPU/Network overhead if tracing is not needed at 100%.
- **Tradeoffs**: 
    - Less observability.
- **Suggested Approach**: 
    - Make `WithSQLCommenter` configurable via Env Var.
    - Verify sampling rate.

## 26. Bulk Insert for Schema Fields

- **Affected Components**: `SchemaFieldsService.Create`
- **Current Behavior**: Iterates and inserts schema fields one by one.
- **Inefficiency**: 
    - N+1 DB round trips where N is number of fields (can be 50+).
- **Evidence**: `internal/application/schema_fields/usecase.go` loop calling `uc.fields.Create`.
- **Expected Improvement**: 
    - 1 DB round trip.
- **Tradeoffs**: 
    - Complexity of constructing bulk query or using `CopyFrom`.
- **Suggested Approach**: 
    - Use `sqlc`'s `CopyFrom` or manual batch insert.

## 27. Inject Build Version into Telemetry

- **Affected Components**: `telemetry/otel.go`
- **Current Behavior**: Hardcoded `semconv.ServiceVersion("dev")`.
- **Inefficiency**: 
    - Harder to correlate performance regressions with specific builds/commits.
- **Evidence**: `internal/infrastructure/telemetry/otel.go` uses string literal.
- **Expected Improvement**: 
    - Better observability.
- **Tradeoffs**: 
    - Build process change (ldflags).
- **Suggested Approach**: 
    - Use `ldflags` to inject version variable and use it in `otel.go`.

---

# Summary of Priorities

## High Impact (10x+ Gains in Hot Paths)
- **1. Auth Middleware Session Caching**: Critical for scaling auth.
- **6. Pagination for List Endpoints**: Critical for stability and avoiding O(N) blowups.
- **14. Optimize Project Users List Index**: Avoids sort on hot path.
- **15. Add Missing Index on projects(owner_id)**: Fixes linear scan on project listing.
- **21. Cursor-Based Pagination**: Enables deep pagination without performance hit.
- **24. Configure Database Connection Pool**: Essential for production stability.

## Medium Impact (Efficiency & Latency)
- **2. Project Public Key Caching**: Reduces DB load for project users.
- **7. HTTP Response Compression**: Reduces network usage significantly.
- **8. Migrate to pgx Driver**: General performance lift.
- **9. Optimize GetVerbose Field Assembly**: Fixes quadratic algo in schema fetch.
- **22. High-Performance JSON Library**: Faster parsing for JSON-heavy APIs.
- **26. Bulk Insert for Schema Fields**: Optimizes schema creation.

## Stability & Operational Excellence
- **3. Offload Encryption to Application**: Removes DB CPU bottleneck & improves security.
- **11. Batch Data Cleanup Jobs**: Prevents maintenance spikes.
- **12. Separate Migrations from Startup**: Safer deployments.
- **13. Limit Request Body Size**: Anti-DoS.
- **18. Remove Baked-in Keys**: Corrects operational flaw.
- **19. Pin Database Image Version**: Stability.
- **23. Configure HTTP Server Timeouts**: Anti-Slowloris.
