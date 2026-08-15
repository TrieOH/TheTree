---
name: improve-user-experience
description: Scan how the product is actually experienced by whoever consumes it — end users, API integrators, CLI users, or docs readers — for friction points, then present them as a visual HTML report ranked by how much they hurt, and grill deeper on whichever one gets picked.
disable-model-invocation: true
---

# Improve User Experience

Surface **experience friction** — the places where whoever consumes this product has to work harder than they should to get what they came for. "User" here is not assumed to mean "person looking at a UI." A user is whoever is on the other side of the interface: someone clicking through a screen, someone integrating against an API, someone piping output through a CLI, someone reading the docs to figure out which endpoint to call. Pick the right lens for what's actually being audited — don't force a UI vocabulary onto an API, and don't force an API vocabulary onto a UI.

This is not a feature-discovery pass (see `/find-new-features` for that) and not a bug audit (see `/audit-codebase`). The product can be working exactly as designed and still have bad experience — that gap is what this skill finds.

- The domain language in `CONTEXT.md` names the things the user is trying to do; use it instead of internal names.
- ADRs in `docs/adr/` may record why an experience is deliberately shaped the way it is — check before flagging something as friction that was actually a considered trade-off.
- If `CONTEXT.md` or ADRs don't exist, proceed without them.

## Process

### 1. Identify the interface and the user

Before scanning, be explicit about what's being audited — the fix and the vocabulary differ a lot by interface:

- If the user named a direction — a specific flow, endpoint, screen, command — take it.
- Otherwise, ask, or infer from the codebase: is the primary interface a UI, an API, a CLI, an SDK, or docs? Products often have more than one; if so, scope to one interface per run rather than mixing vocabularies in a single report.

Name who's on the other end concretely — not "the user," but "a developer integrating payments for the first time," "a support agent looking up an order," "a script calling this endpoint on a cron." The friction that matters depends entirely on who's hitting it and in what state (first time vs. hundredth time, calm vs. mid-incident).

### 2. Find friction

Use the Agent tool with `subagent_type=Explore` to walk the interface the way that user actually would — not the way the code is organized. Then, or in place of code exploration when there's someone to ask, run the `/grilling` skill to interview the user (or whoever plays the user's role — support, a real integrator, whoever's available) about where they get stuck, what they have to re-read, what they work around. Combine both where possible: code exploration finds structural friction, grilling finds the friction nobody thought to name in a ticket.

What counts as friction depends on the interface:

- **UI** — unclear next step, error states that don't say what to do, a flow that takes more steps than the task needs, state that resets when it shouldn't, feedback that arrives too late or not at all.
- **API / SDK** — inconsistent naming or shapes across endpoints, undocumented required fields, errors that don't say which field or why, pagination or auth that surprises on the first attempt, a happy path that requires reading five endpoints' docs to assemble.
- **CLI** — flags that don't compose, unhelpful `--help`, silent failure, output that's hard to pipe or parse, a common task that takes many invocations instead of one.
- **Docs** — the answer exists but can't be found, the example doesn't match the current API, missing the one example that would have saved the read.

For every candidate, apply the **first-attempt test**: would someone doing this for the first time, without asking for help, get it right? If yes, it's probably not friction worth reporting even if it could theoretically be smoother — this skill is for real friction, not polish for its own sake.

Grade severity honestly:

- **Blocking** — the user can't complete the task without outside help (support ticket, Slack message, reading source code).
- **Costly** — completable, but takes meaningfully longer or more attempts than it should.
- **Rough** — noticeable friction, low cost — a bad first impression more than a real obstacle.
- **Polish** — true but minor; include only a few, never as a headline finding.

### 3. Present findings as an HTML report

Write a self-contained HTML file to the OS temp directory so nothing lands in the repo. Resolve the temp dir from `$TMPDIR`, falling back to `/tmp` (or `%TEMP%` on Windows), and write to `<tmpdir>/ux-report-<timestamp>.html` so each run gets a fresh file. Open it for the user — `xdg-open <path>` on Linux, `open <path>` on macOS, `start <path>` on Windows — and tell them the absolute path.

See [HTML-REPORT.md](HTML-REPORT.md) for the full scaffold, card layout, and styling guidance.

End the report with a **Fix first** section: the single finding you'd address before anything else, and why (usually: most blocking, cheapest fix, or the one that shows up earliest in the user's path).

Do NOT design the fix yet. After the file is written, ask the user: "Which of these would you like to go deeper on?"

### 4. Grill the fix

Once the user picks a finding, run `/grilling` to work through the fix: what the user actually needs at that point (not just what's currently missing), what the smallest change is that clears the blocker, whether the fix could introduce friction somewhere else in the same flow, how you'd confirm it actually helped (a metric, a follow-up with a real user, a before/after walkthrough).

Side effects happen inline as decisions crystallize:

- **Naming a concept not yet in `CONTEXT.md`?** Add it. Create the file lazily if it doesn't exist.
- **User rejects a finding with a load-bearing reason** (it's intentional, the cost is acceptable, it's already being reworked)? Note it, offer an ADR only if the codebase already uses them and a future pass would otherwise re-flag the same thing.
- **Ready to build the fix?** Hand off rather than continuing to design inline — `/implement` for straightforward changes, `/to-spec` or `/to-tickets` if it's bigger than one sitting.