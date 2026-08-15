# HTML Report Format

The feature report is rendered as a single self-contained HTML file in the OS temp directory. Tailwind comes from a CDN. This report is lighter on diagrams than the architecture or audit ones — the point isn't to show a structure, it's to show a set of bare ideas clearly enough that the user can pick one fast. Reach for a small hand-built visual only when a candidate genuinely benefits from one (a before/after of a user flow, a rough mock of a screen); most cards need none.

## Scaffold

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>Feature ideas — {{product name}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
      .now { border-left-color: #059669; }
      .later { border-left-color: #d97706; }
      .longshot { border-left-color: #64748b; }
    </style>
  </head>
  <body class="bg-stone-50 text-slate-900 font-sans">
    <main class="max-w-5xl mx-auto px-6 py-12 space-y-12">
      <header>...</header>
      <section id="candidates" class="space-y-6">...</section>
      <section id="try-first">...</section>
    </main>
  </body>
</html>
```

## Header

Product name, date, and one line naming who the brainstorm was grounded in (e.g. "Grounded in: existing checkout users, competitor gap vs. X"). No introduction paragraph — straight into the candidates.

## Candidate card

Cards here are **denser and shorter** than an architecture or audit card — this is a scan-and-pick list, not a deep-dive report. Each candidate is one `<article class="border-l-4 rounded bg-white p-4">`, left-border colored by timeframe (`now` emerald, `later` amber, `longshot` slate).

- **Title** — the idea in one line, user-facing language, not a ticket title. "Let users export their order history as CSV," not "Add export endpoint."
- **Timeframe badge** — `Now`, `Later`, `Long shot` — matching the left border color.
- **Who it's for** — one short phrase. A named user segment, not "users."
- **The idea** — one to two sentences, plain. What the user would actually see or do differently.
- **Why now/later/long shot** — one sentence naming the real reason: a dependency, an unvalidated assumption, or just "nothing's blocking this."
- **Signal** — bullets, ≤6 words each — the concrete thing that suggested this idea (a support ticket pattern, a competitor feature, a workaround users already do). Not vibes — cite the actual signal from the grill.
- **Assumption** (Long shot only) — one line, amber-tinted box, naming the bet this depends on.

No paragraphs. If a candidate needs more than the fields above to make sense, it's not bare enough yet — that's what the deepening grill in step 4 is for, not the report.

## Optional small visual

Use only when it earns its place:

- **Flow sketch** — 3-4 boxes in a row (`flex items-center gap-2`, small arrow characters or thin SVG lines between them) showing the user's current path vs. the proposed one. Good for "this removes N steps" ideas.
- **Rough mock** — a very simple bordered rectangle with placeholder text blocks (`bg-slate-100 h-3 rounded w-full`) suggesting a screen's shape. Good for "this needs a new surface" ideas. Not a real UI — just enough to make the idea land visually. Keep it under ~200px tall.

Don't force one of these onto every card. Most candidates are a sentence and a "why" — that's enough.

## Try first section

One larger card, same left-border treatment (emerald, since this should almost always be a `Now`). Candidate name, one sentence on why this one first, anchor link to its card.

## Style guidance

- Lean editorial, not corporate-dashboard. Generous whitespace. Serif optional for the page title only (`font-serif`) — card titles stay sans, they need to scan fast.
- Colour carries the timeframe, nothing else — don't decorate beyond the three timeframe colors plus the amber assumption callout.
- This report is meant to be read top to bottom in under a minute on first pass. If it's taking longer, the cards are too dense — cut, don't add.
- Order candidates `Now` → `Later` → `Long shot`. Within a group, order by how strong the signal is, strongest first.
- The only script is the Tailwind CDN. Static otherwise — no interactivity needed, this isn't a diagram-heavy report.

## Tone

Plain, user-facing language throughout — even though the audience reading the report is the builder, the candidates should read the way a user would describe wanting them, not the way an engineer would describe building them.

**Use exactly:** the product's own domain terms from `CONTEXT.md` for anything user-facing (if `CONTEXT.md` calls it "checkout," say checkout, not "the payment flow").

**Never substitute:** internal module/service names for user-facing concepts · "leverage" a feature (say what it does) · "synergy," "unlock," "empower" — say the concrete thing a user does differently.

**Phrasings that fit the style:**

- "Let returning users skip the address step entirely."
- "Three support tickets this month asked for this exact export."
- "Blocked on: Payssage doesn't support partial refunds yet."
- "Bet: enough users hit the free-tier ceiling for this to matter."

**Signal bullets** cite the real thing that surfaced the idea: *"3 tickets, same ask"*, *"competitor shipped this in March"*, *"users already do this manually"*. Don't write *"users would probably like this"* — that's not a signal, that's a hope. If there's no real signal, say so plainly and mark it Long shot rather than inventing one.

No hedging, no throat-clearing, no "it's worth noting that…". If a sentence could be a bullet, make it a bullet. If a bullet could be cut, cut it.