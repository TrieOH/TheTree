# Vulnerability Checklist

Reference for the vuln-scanner skill. Check every applicable pattern per file being audited.

---

## Universal (all languages)

### Secrets & Credentials
- [ ] Hardcoded API keys, tokens, passwords, private keys in source
- [ ] Secrets in environment variable defaults (`.env.example` with real values)
- [ ] Secrets committed to config files (`config.yml`, `docker-compose.yml`)
- [ ] Private keys or certs checked into repo
- [ ] Debug/test credentials that could reach production

### Infrastructure / Config
- [ ] Docker images running as root with no `USER` directive
- [ ] Exposed ports that shouldn't be public (`0.0.0.0` vs `127.0.0.1`)
- [ ] Missing resource limits in Docker/compose (memory, CPU)
- [ ] `privileged: true` in compose without justification
- [ ] Volumes mounting sensitive host paths (e.g. `/var/run/docker.sock`)
- [ ] CI/CD secrets printed in logs (`echo $SECRET`)
- [ ] Caddy/nginx configs with overly permissive CORS (`Access-Control-Allow-Origin: *` on auth endpoints)

---

## Go

### Injection
- [ ] SQL built via `fmt.Sprintf` or string concat instead of parameterized queries
- [ ] Shell commands via `exec.Command` with user-controlled input
- [ ] `os.Open` / `ioutil.ReadFile` with unsanitized path (path traversal)
- [ ] Template injection: `template/text` used instead of `template/html` for HTML output
- [ ] SSRF: `http.Get(userInput)` without URL allowlist validation

### Authentication & Authorization
- [ ] JWT parsed without algorithm validation (`alg: none` accepted)
- [ ] JWT secret hardcoded or loaded from insecure source
- [ ] Missing signature verification on incoming webhooks
- [ ] Auth middleware not applied to all routes (check router setup carefully)
- [ ] API keys compared with `==` instead of `subtle.ConstantTimeCompare`
- [ ] Tokens stored in plaintext (should be hashed, e.g. SHA-256)
- [ ] Missing expiry check on tokens/sessions
- [ ] OAuth state parameter not validated (CSRF on OAuth callback)
- [ ] `context` user claims trusted without re-validation from token

### Cryptography
- [ ] `math/rand` used for security-sensitive randomness (use `crypto/rand`)
- [ ] Weak hash: MD5 or SHA1 used for passwords or security tokens
- [ ] Custom crypto implementation instead of stdlib
- [ ] IV/nonce reuse in AES-GCM or AES-CBC
- [ ] Missing HMAC on encrypted data (encryption without authentication)

### Race Conditions & Concurrency
- [ ] Shared mutable state accessed without mutex
- [ ] TOCTOU: check-then-use without holding a lock or using atomic ops
- [ ] Token/key rotation logic with a race window

### Input Validation
- [ ] Missing validation on pagination params (negative offset, huge limit → DoS)
- [ ] Missing max size on file/body uploads (`http.MaxBytesReader` not set)
- [ ] Unvalidated redirect URLs (open redirect)
- [ ] Email/username normalization missing (unicode lookalike attacks)

### Error Handling & Info Disclosure
- [ ] Errors returned to client with internal stack traces or DB schema details
- [ ] Verbose error messages leaking service topology
- [ ] Panic recovery logging sensitive request data

### Dependency / Supply Chain
- [ ] `go.sum` missing or not committed
- [ ] Dependencies with known CVEs (check `govulncheck` output if available)
- [ ] Use of `replace` directives pointing to local/unreviewed forks

---

## TypeScript / JavaScript

### Injection
- [ ] `dangerouslySetInnerHTML` with user data (XSS)
- [ ] `eval()` / `new Function()` with any dynamic input
- [ ] Template literals inserted into `innerHTML`
- [ ] SQL via string concat in any DB client
- [ ] `child_process.exec(userInput)` — shell injection
- [ ] Server-side: `res.redirect(req.query.url)` — open redirect
- [ ] `JSON.parse` on untrusted input without try/catch (denial of service via prototype pollution on some parsers)

### Authentication & Authorization
- [ ] JWT verified client-side only (no server-side verification)
- [ ] `localStorage` used to store tokens (prefer `httpOnly` cookies)
- [ ] CSRF tokens missing on state-changing requests
- [ ] Missing `SameSite` cookie attribute
- [ ] Missing `HttpOnly` / `Secure` on auth cookies
- [ ] Auth state stored in non-protected React context accessible to third-party scripts

### Next.js / SSR specific
- [ ] `getServerSideProps` passing sensitive env vars to the client-side bundle
- [ ] API routes missing auth middleware
- [ ] `next.config.js` with `dangerouslyAllowSVG: true` without `contentDispositionType: attachment`
- [ ] Exposed `/api/` routes not protected with rate limiting

### Cryptography
- [ ] `Math.random()` used for tokens, OTPs, or session IDs
- [ ] Client-side only encryption (security theater)
- [ ] Weak algorithms in `crypto` module (MD5, SHA1)

### Dependencies
- [ ] `npm audit` findings not addressed
- [ ] `package-lock.json` / `yarn.lock` not committed
- [ ] Packages with overly broad permissions (e.g. postinstall scripts)

### Miscellaneous
- [ ] `console.log` printing sensitive data in production code
- [ ] Source maps deployed to production (leaks original source)
- [ ] Debug/internal endpoints reachable in production builds

---

## Python

### Injection
- [ ] `subprocess.call(shell=True)` with user-controlled input
- [ ] f-string or `%`-formatted SQL queries
- [ ] `pickle.loads(userInput)` — arbitrary code execution
- [ ] `yaml.load()` without `Loader=yaml.SafeLoader`
- [ ] `eval()` / `exec()` on any user data
- [ ] Path traversal via `open(user_path)` without normalization

### Authentication
- [ ] `JWT` decoded without verifying signature (`options={"verify_signature": False}`)
- [ ] Password hashing with `hashlib.md5/sha1` instead of `bcrypt`/`argon2`
- [ ] Missing `secrets` module usage for token generation (`random` used instead)

### Miscellaneous
- [ ] Debug mode enabled in Flask/Django in non-dev environments (`DEBUG=True`)
- [ ] Django `SECRET_KEY` hardcoded
- [ ] `ALLOWED_HOSTS = ['*']` in Django settings
- [ ] Unrestricted file uploads with no MIME type validation

---

## Database / Migrations

- [ ] Migration files dropping columns/tables without a rollback plan documented
- [ ] Sensitive columns stored in plaintext (passwords, PII, payment info)
- [ ] Missing indexes on foreign keys used in auth/authz queries (can cause full table scans exploitable via timing)
- [ ] `GRANT ALL` to app DB user instead of least-privilege
- [ ] Backup scripts with credentials embedded in the script

---

## Scoring priority

When multiple vulns exist, prioritize writing files in this order:
1. CRITICAL first
2. Then HIGH
3. Then MEDIUM
4. Then LOW

This ensures the most important findings are written even if the session is interrupted.