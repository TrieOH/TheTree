---
name: find-new-features
description: Grill the user to surface raw feature and improvement ideas for the product, present them as a visual HTML report ranked by user value, then grill deeper on whichever one they pick.
disable-model-invocation: true
---

# Find New Features

Surface **feature and improvement candidates** for the product's users — not implementation tasks, not internal tooling, not architecture. The test throughout is the user's: would someone using this product notice and want the thing being proposed?

This command brackets the `/grilling` interview on both sides — once to generate raw candidates, once to deepen the one that's picked. Where `/improve-codebase-architecture` grills *after* presenting candidates (to decide how to build the one you chose), this grills *before* as well (to decide what's even worth proposing) — brainstorming produces sharper candidates than staring at the codebase alone does.

- The domain language in `CONTEXT.md` names the things users actually do with the product; use it instead of internal names (`the Order intake module` internally might just be "checkout" to a user — talk like the user).
- ADRs in `docs/adr/` may record why something was deliberately *not* built — check before proposing it as new.
- If `CONTEXT.md` or ADRs don't exist, proceed without them.

## Process

### 1. Ground the brainstorm

Before grilling, get enough context that the ideas are about *this* product, not generic SaaS suggestions:

- If the user named a direction — a user segment, a part of the product, a competitor gap — take it, and let it frame the interview.
- Otherwise, skim `CONTEXT.md` for the domain model, skim recent `docs/adr/` for what's already been decided against, and glance at recently closed issues/tickets if the tracker is reachable, so the interview starts from the real shape of the product instead of a blank page.

### 2. Grill — raw ideas

Run the `/grilling` skill to interview the user toward a wide set of **bare ideas** — not to converge on one thing yet. Push on:

- Who is this for — which existing user, or which user you don't have yet?
- What can't they do today that they clearly want to?
- What do they currently do around the product (workarounds, exports, other tools stitched in) that the product could absorb?
- What would make an existing feature worth using more often, not just usable?
- What's the version of this product that would make a specific named competitor's feature irrelevant?

The interview should end with a rough list of candidate ideas, still bare — a sentence each, not fleshed out. Don't let the interview collapse onto one favorite; that convergence happens in step 4, with the user, in front of the report.

### 3. Present candidates as an HTML report

Write a self-contained HTML file to the OS temp directory so nothing lands in the repo. Resolve the temp dir from `$TMPDIR`, falling back to `/tmp` (or `%TEMP%` on Windows), and write to `<tmpdir>/feature-report-<timestamp>.html` so each run gets a fresh file. Open it for the user — `xdg-open <path>` on Linux, `open <path>` on macOS, `start <path>` on Windows — and tell them the absolute path.

See [HTML-REPORT.md](HTML-REPORT.md) for the full scaffold, card layout, and styling guidance.

For every candidate, grade honestly rather than generously:

- **Now** — clear user value, fits the product as it exists today, no architecture rework implied.
- **Later** — clear user value, but blocked on something (a dependency, a missing primitive, a decision the team hasn't made).
- **Long shot** — real value *if* a bet pays off (new user segment, unproven demand); flag the assumption it depends on.

End the report with a **Try first** section: the single candidate you'd validate or build first, and why (usually: highest value for the least unlocked-dependency risk).

Do NOT design the feature yet. After the file is written, ask the user: "Which of these would you like to go deeper on?"

### 4. Grill — deepen the pick

Once the user picks a candidate, run `/grilling` again — this time to converge, not to diverge. Walk the decision tree: which user it's really for, what the smallest version looks like, what it explicitly doesn't do yet, what existing feature it might cannibalize or conflict with, how you'd know it worked after shipping.

Side effects happen inline as decisions crystallize:

- **Naming something not yet in `CONTEXT.md`?** Add the term. Create the file lazily if it doesn't exist.
- **User decides against a candidate for a load-bearing reason** (not just "not now")? Offer an ADR the same way `/improve-codebase-architecture` does — only when a future brainstorm would otherwise re-propose the same thing.
- **Ready to move from idea to plan?** Hand off to `/to-spec` or `/to-tickets` rather than continuing to design inline here — this skill's job ends at a well-understood idea, not a shipped feature.