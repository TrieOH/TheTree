# HTML Report Format

The UX report is rendered as a single self-contained HTML file in the OS temp directory. Tailwind comes from a CDN. No Mermaid here — friction isn't graph-shaped, it's a sequence with a sticking point, and a simple horizontal step strip communicates that better than a flowchart would. Reach for hand-built SVG only for the step strip; nothing else needs a diagram.

## Scaffold

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>UX review — {{product name}} ({{interface}})</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
      .blocking { border-left-color: #dc2626; }
      .costly { border-left-color: #d97706; }
      .rough { border-left-color: #eab308; }
      .polish { border-left-color: #94a3b8; }
      .stuck-step { fill: #fee2e2; stroke: #dc2626; }
    </style>
  </head>
  <body class="bg-stone-50 text-slate-900 font-sans">
    <main class="max-w-5xl mx-auto px-6 py-12 space-y-12">
      <header>...</header>
      <section id="findings" class="space-y-8">...</section>
      <section id="fix-first">...</section>
    </main>
  </body>
</html>
```

## Header

Product name, date, **interface being audited** (UI / API / CLI / docs — say which, prominently, it changes how to read every card), and who was named as the user (e.g. "Lens: a developer integrating payments for the first time"). No introduction paragraph — straight into the findings.

## Finding card

Each finding is one `<article class="border-l-4 rounded bg-white p-4">`, left-border colored by severity (`blocking` red, `costly` orange, `rough` yellow, `polish` slate).

- **Title** — the friction in plain language, from the user's side. "Error doesn't say which field failed validation," not "improve error handling."
- **Severity badge** — `Blocking`, `Costly`, `Rough`, `Polish` — matching the border color. Plus a small interface tag (`ui`, `api`, `cli`, `docs`) if the report ever mixes interfaces (it usually shouldn't, per SKILL.md — but tag anyway in case scope grew mid-run).
- **Who hits this** — the concrete user from step 1, not "users."
- **The step strip** — a small horizontal sequence showing where in the task flow this occurs, with the sticking point highlighted. See pattern below. Only for `Blocking` and `Costly` — skip for `Rough`/`Polish`, a sentence is enough there.
- **What happens** — one to two sentences, from the user's point of view, present tense. "They submit the form, get a 400, and the message just says 'invalid request.'"
- **Fix** — one sentence, plain, what changes.
- **Cost today** — bullets, ≤6 words each, concrete: "avg. 3 support tickets/week", "adds ~10 min to first call", "2 extra round trips to debug." Use a real number if you have one; if you don't, say the qualitative cost plainly rather than inventing a number.
- **Context callout** (if applicable) — one line, amber-tinted box, only when relevant (e.g. "intentional — rate limit is a pricing lever, not a bug").

No paragraphs. If a finding needs more than this to land, the finding isn't specific enough yet — narrow it.

## Step strip pattern

A single row of small boxes connected by thin lines, one box per step in the task (`Sign up → Get API key → First request → Handle response`, or `Land on page → Fill form → Submit → Confirmation`). The step where the friction lives gets the `stuck-step` treatment — red fill, red border, maybe a small "⏸" or "!" glyph — the rest stay neutral (`bg-slate-100 border-slate-300`).

```html
<div class="flex items-center gap-1 my-3">
  <div class="px-3 py-2 rounded bg-slate-100 border border-slate-300 text-xs">Sign up</div>
  <div class="w-4 h-px bg-slate-300"></div>
  <div class="px-3 py-2 rounded bg-slate-100 border border-slate-300 text-xs">Get API key</div>
  <div class="w-4 h-px bg-slate-300"></div>
  <div class="px-3 py-2 rounded stuck-step text-xs font-medium">First request → 400, no detail</div>
  <div class="w-4 h-px bg-slate-300"></div>
  <div class="px-3 py-2 rounded bg-slate-100 border border-slate-300 text-xs opacity-50">Handle response</div>
</div>
```

Keep it to 3-6 steps. If the real flow has more, collapse the ones before and after the sticking point into single boxes ("...earlier steps", "...later steps") rather than drawing all ten.

## Fix first section

One larger card, same left-border treatment (red, since Fix first is almost always a Blocking finding). Names the finding, one sentence why first, anchor link to its card.

## Style guidance

- Lean editorial, not corporate-dashboard. Generous whitespace. Serif optional for the page title only — card content stays sans, it needs to scan fast and read like a real user's account, not a report.
- Colour carries severity and nothing else. Red is reserved for Blocking and the stuck step in the strip — don't decorate elsewhere.
- Order findings `Blocking` → `Costly` → `Rough` → `Polish`. Within a severity, order by where it sits in the user's path — earliest friction first, since that's what determines whether they ever reach the rest.
- The only script is the Tailwind CDN. Static otherwise.

## Tone

Write every "what happens" section as if narrating over someone's shoulder while they try to use the thing — present tense, concrete, no editorializing.

**Use exactly:** friction, blocking, first-attempt, the interface's own vocabulary (endpoint/field/status code for an API; flag/stdout/exit code for a CLI; screen/state for a UI) — match the vocabulary to the interface named in the header, don't default to UI language for an API finding.

**Never substitute:** "UX issue" (say what actually happens) · "improve usability" (say the specific change) · "intuitive" / "user-friendly" as a goal — these aren't falsifiable, describe the concrete behavior instead.

**Phrasings that fit the style:**

- "They get a 400 and the body says `{\"error\": \"invalid request\"}` — no field name."
- "Three clicks to undo something that took one click to do."
- "The `--dry-run` flag silently does nothing on this subcommand."
- "The quickstart example uses an endpoint that was renamed two versions ago."

**Cost bullets** stay concrete: *"avg 3 tickets/week"*, *"adds ~10 min first call"*, *"users retry the same request 4x before giving up"*. If there's no hard number, name the qualitative cost plainly — don't manufacture false precision.

No hedging, no throat-clearing, no "it's worth noting that…". If a sentence could be a bullet, make it a bullet. If a bullet could be cut, cut it.