# ADR Format

ADRs live in `docs/adr/` as `ADR-NNNN.md`, numbered in order. To pick the next number, take the highest in `docs/adr/` and add one.

## Template

```md
---
number: ADR-0023
title: Short sentence saying what we decided
status: accepted
date: 2026-11-21
summary: One sentence a reader can stop at.
scope:
  - lib/go/authz
  - api/*/internal/app
---

## Why

What was wrong or missing, in plain words. Facts and numbers, not adjectives.

## Decision

What we do now. Lead with the rule, then the details someone needs to follow it.

## Consequences

What this costs, what it changes, and what a future reader might trip over.
```

An ADR can be short. The value is recording *that* we decided and *why*, not filling in sections.

## The preamble

The YAML block at the top is what the dev portal reads (`front/portal/mdx/remark-adr.ts`). Keep it valid YAML, and don't repeat the title as a `#` heading in the body.

| Field | Required | Value |
|---|---|---|
| `number` | yes | `ADR-NNNN`, same as the file name |
| `title` | yes | the decision as a short sentence |
| `status` | yes | `proposed`, `accepted`, `superseded` or `deprecated` |
| `date` | yes | when it was decided, `YYYY-MM-DD` |
| `summary` | yes | one sentence |
| `scope` | no | list of paths it applies to |
| `updated` | no | last time the decision was changed in place |
| `supersedes` | no | list of ADRs this one replaces |
| `superseded_by` | no | the ADR that replaced this one (set `status: superseded`) |
| `amends` | no | list of ADRs this one changes in part |
| `amended_by` | no | list of ADRs that changed this one in part |

Keep the links both ways: if a new ADR amends or supersedes an old one, update the old one's preamble too.

## Optional sections

Add them only when they help:

- **Rejected**: alternatives worth remembering, so nobody proposes them again without the context.
- **Examples** or **Rollout**: when the decision is easier to follow with code or steps.
- **Changes**: dated notes when the decision was later changed in place. Write the Decision as it stands now, and list what changed here.

## Writing style

- Plain sentences. Say it the way you'd explain it to a teammate.
- No em dashes, no filler, no hype words, no "this ADR". Say what happens, not how important it is.
- Name real files, functions and numbers. Skip line counts unless they're the point.

## When to offer an ADR

All three of these must be true:

1. **Hard to reverse**: the cost of changing your mind later is meaningful
2. **Surprising without context**: a future reader will look at the code and wonder "why on earth did they do it this way?"
3. **The result of a real trade-off**: there were genuine alternatives and you picked one for specific reasons

If a decision is easy to reverse, skip it: you'll just reverse it. If it's not surprising, nobody will wonder why. If there was no real alternative, there's nothing to record beyond "we did the obvious thing."

### What qualifies

- **Architectural shape.** "We're using a monorepo." "The write model is event-sourced, the read model is projected into Postgres."
- **Integration patterns between contexts.** "Ordering and Billing communicate via domain events, not synchronous HTTP."
- **Technology choices that carry lock-in.** Database, message bus, auth provider, deployment target. Not every library: just the ones that would take a quarter to swap out.
- **Boundary and scope decisions.** "Customer data is owned by the Customer context; other contexts reference it by ID only." The explicit no-s are as valuable as the yes-s.
- **Deliberate deviations from the obvious path.** "We're using manual SQL instead of an ORM because X." Anything where a reasonable reader would assume the opposite. These stop the next engineer from "fixing" something that was deliberate.
- **Constraints not visible in the code.** "We can't use AWS because of compliance requirements." "Response times must be under 200ms because of the partner API contract."
- **Rejected alternatives when the rejection is non-obvious.** If you considered GraphQL and picked REST for subtle reasons, record it; otherwise someone will suggest GraphQL again in six months.
