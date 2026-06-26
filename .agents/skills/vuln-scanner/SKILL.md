---
name: vuln-scanner
description: Security vulnerability scanner for Go/TS/Python monorepos. Use this skill whenever the user asks to: scan for vulnerabilities, audit security, find security issues, check for security bugs, review code for vulns, find CVEs, check for injection flaws, audit authentication/authorization code, or anything involving "security scan", "vuln scan", "audit my code". Also trigger when the user says things like "is my code secure?", "check for XSS/SQLi/SSRF/auth issues", or "review my codebase for problems". Each finding is written as a structured .md file under vulns/ at the monorepo root. Always use this skill before attempting a manual security review.
---

# Vulnerability Scanner

Scans a monorepo for security vulnerabilities and writes structured findings to `vulns/` at the root.

## Workflow

### 1. Orient to the repo

```bash
ls -la .
cat go.work 2>/dev/null || true
find . -name "go.mod" -not -path "*/vendor/*" | head -20
find . -name "package.json" -not -path "*/node_modules/*" | head -20
```

Map each service: note its type (Go/TS/Python), its path, and its name.

### 2. Enumerate files to audit

Collect all source files, skipping noise:

```bash
find . \
  -type f \
  \( -name "*.go" -o -name "*.ts" -o -name "*.tsx" -o -name "*.js" -o -name "*.py" \) \
  -not -path "*/vendor/*" \
  -not -path "*/node_modules/*" \
  -not -path "*/.git/*" \
  -not -path "*/dist/*" \
  -not -path "*/build/*" \
  | sort
```

Also check config/infra files:
```bash
find . -type f \( -name "*.yml" -o -name "*.yaml" -o -name "*.env*" -o -name "Dockerfile*" -o -name "*.conf" \) \
  -not -path "*/.git/*" -not -path "*/node_modules/*"
```

### 3. Audit each file

Read files in logical groups (per service or per layer). For each file, check against the full checklist in `references/vuln-checklist.md`.

**Read `references/vuln-checklist.md` before starting.** It contains the full pattern list per language and vulnerability class.

Be thorough — config files, migration files, and test files frequently contain hardcoded secrets or dangerous patterns.

### 4. Write findings

Create `vulns/` first:
```bash
mkdir -p vulns/
```

For **each** vulnerability found, create one file named:
```
vulns/<SEVERITY>-<service-name>-<short-slug>.md
```

Examples:
- `vulns/CRITICAL-identityx-jwt-alg-none.md`
- `vulns/HIGH-informd-sqli-form-query.md`
- `vulns/MEDIUM-payssage-ssrf-webhook-url.md`

**File contents — exact format:**

```markdown
**Vulnerability Details**
- **Service:** SERVICE_TYPE/SERVICE_NAME
- **File:** PATH/TO/FILE.EXT
- **Line:** 123
- **Severity:** [LOW|MEDIUM|HIGH|CRITICAL]
- **Title:** Short descriptive title
- **Impact:** What an attacker could achieve
- **Description:** What the vulnerable code does and why it's a problem
- **Proposed fix:** Concrete code change or architectural fix
- **Reproduction Steps:** Step-by-step to trigger the vulnerability
```

`SERVICE_TYPE` = language/framework (`Go`, `TypeScript`, `Python`).
`SERVICE_NAME` = service directory name (`identityx`, `informd`, `payssage`, `univents`).

### 5. Summarize

After all files are written, print a summary table:

```
## Vulnerability Scan Summary

| Severity | Count |
|----------|-------|
| CRITICAL | N     |
| HIGH     | N     |
| MEDIUM   | N     |
| LOW      | N     |
| **Total**| N     |

Files written to vulns/:
- vulns/CRITICAL-...md
- ...
```

Then ask the user if they want to walk through any finding or start fixing.

---

## Severity Guide

| Severity | When to use                                                                       |
|----------|-----------------------------------------------------------------------------------|
| CRITICAL | RCE, auth bypass, secret/key exposure, full data exfiltration                     |
| HIGH     | Privilege escalation, IDOR, SQLi, SSRF, broken JWT, stored XSS                    |
| MEDIUM   | Reflected XSS, open redirect, missing rate limit, weak crypto in non-auth context |
| LOW      | Info disclosure, verbose errors, missing security headers, unused debug endpoints |

---

## Slug naming

Short, kebab-case, includes vuln class:
`jwt-alg-none` · `sqli-user-query` · `ssrf-webhook` · `hardcoded-secret` · `idor-form-id` · `missing-authn` · `path-traversal` · `weak-hash` · `xss-template` · `open-redirect` · `race-condition` · `insecure-deserialize`

---

## Reference files

- `references/vuln-checklist.md` — full pattern checklist per language. **Read before auditing.**