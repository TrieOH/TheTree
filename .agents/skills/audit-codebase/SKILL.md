---
name: audit-codebase
description: Scan a codebase for clear bugs, security holes, and correctness footguns, present them as a visual HTML report ranked by severity, then walk through fixes for whichever one you pick.
disable-model-invocation: true
---

# Audit Codebase

Surface **concrete defects** — bugs that will misbehave, inputs that will break something, and gaps an attacker could walk through. This is not a style review and not an architecture review; skip anything that's merely ugly, and skip anything that's a genuine design trade-off rather than a mistake. A finding earns its place here only if you can point at the exact input, sequence, or actor that turns it into an incident.

This command is _informed_ by the project's domain model:

- The domain language in `CONTEXT.md` gives names to the things being audited; ADRs in `docs/adr/` record decisions this command should not re-litigate as if they were bugs (a documented trade-off is not a finding).
- If `CONTEXT.md` or ADRs don't exist, proceed without them — this skill doesn't require the codebase-design scaffolding.

## Process

### 1. Scope

**Scope before you scan — YAGNI.** An audit pays off fastest on the code that's exposed to untrusted input or that changes often, so put extra weight there.

- If the user named a direction — a module, an endpoint, a recent PR — take it, and skip the inference below.
- Otherwise, walk back a good stretch of the commit history (`git log --oneline`) and prioritize: anything touching auth, payments, request parsing, file/shell/SQL construction, or serialization boundaries; then anything that's changed a lot recently; then everything else if time allows.

Read `CONTEXT.md` and nearby ADRs for the area you're touching, if they exist, so you don't flag an intentional trade-off as a defect.

### 2. Find

Use the Agent tool with `subagent_type=Explore` to walk the scoped code. Don't run a rigid linter checklist — read the code the way an attacker or an unlucky user would, and pull the thread on anything that feels wrong. Categories to keep in mind, not a checklist to mechanically tick:

- **Bugs** — logic that's wrong on some reachable path: off-by-one, wrong operator, inverted condition, race condition, mishandled error, resource leak, nil/nullable dereference, mutation of shared state, floating-point money math.
- **Security** — injection (SQL, shell, template, log), auth/authz gaps (missing check, wrong check, IDOR), secrets in code or logs, unsafe deserialization, SSRF, path traversal, missing input validation at a trust boundary, weak or missing crypto, dependency with a known CVE actually reachable from this code.
- **Correctness footguns** — the kind of thing that works today and breaks under load, scale, or a slightly different input: unbounded queries, missing timeouts, unhandled pagination, silent truncation, timezone/locale assumptions, retry storms, non-idempotent handlers behind at-least-once delivery.

For every candidate finding, apply the **reachability test**: can you name the concrete input or actor that triggers it? "This could theoretically be a problem" is not a finding. "A request with `id=../../etc/passwd` reads this file" is a finding. If you can't complete the sentence, downgrade it to `Speculative` or drop it.

Grade severity honestly:

- **Critical** — exploitable now, or wrong on the common path, with real damage (data loss, auth bypass, financial error, RCE).
- **High** — exploitable or wrong under conditions that will occur in production, not just adversarial edge cases.
- **Medium** — real defect, but needs an unusual trigger or has a mitigating factor already in place.
- **Speculative** — smells wrong, but you couldn't complete the reachability test. Include at most a few of these, clearly marked, never as the headline.

### 3. Present findings as an HTML report

Write a self-contained HTML file to the OS temp directory so nothing lands in the repo. Resolve the temp dir from `$TMPDIR`, falling back to `/tmp` (or `%TEMP%` on Windows), and write to `<tmpdir>/audit-report-<timestamp>.html` so each run gets a fresh file. Open it for the user — `xdg-open <path>` on Linux, `open <path>` on macOS, `start <path>` on Windows — and tell them the absolute path.

See [HTML-REPORT.md](HTML-REPORT.md) for the full scaffold, card layout, diagram patterns, and styling guidance.

End the report with a **Fix first** section: the single finding you'd patch before anything else ships, and why.

Do NOT write the fix yet. After the file is written, ask the user: "Which of these would you like to fix first?"

### 4. Fix loop

Once the user picks a finding, walk through the fix with them directly — no separate interview skill needed, since a bug fix is usually narrower in scope than an architecture change. Confirm:

- The exact reachable trigger, so the fix targets the real cause and not a symptom.
- Whether the fix needs a regression test (it almost always does — write one that fails before the fix and passes after).
- Whether the same mistake is likely repeated elsewhere in the codebase (grep for the same pattern before declaring done).

If the user rejects a finding with a load-bearing reason — it's intentionally permissive, the input is already validated upstream, whatever — don't argue; note it and move on. Only suggest an ADR if the codebase already uses them and the reasoning would genuinely save a future auditor from re-flagging the same thing.