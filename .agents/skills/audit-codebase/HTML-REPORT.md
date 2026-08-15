# HTML Report Format

The audit is rendered as a single self-contained HTML file in the OS temp directory. Tailwind and Mermaid both come from CDNs. Mermaid handles graph-shaped diagrams (call paths, request flow); hand-built divs and inline SVG handle the more editorial visuals (data flow with a highlighted breach point, trust-boundary crossings). Mix the two — don't lean on Mermaid for everything, it'll start to look generic.

Where the architecture report is organized around **before/after**, this one is organized around **severity**. The diagram's job here isn't to show a transformation — it's to show exactly where the bad input enters and what it touches on the way to causing damage.

## Scaffold

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>Audit — {{repo name}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script type="module">
      import mermaid from "https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.esm.min.mjs";
      mermaid.initialize({ startOnLoad: true, theme: "neutral", securityLevel: "loose" });
    </script>
    <style>
      /* small custom layer for things Tailwind doesn't cover cleanly */
      .trigger { stroke: #dc2626; stroke-width: 2px; }
      .tainted { stroke-dasharray: 4 4; stroke: #dc2626; }
      .critical { background: linear-gradient(135deg, #450a0a, #7f1d1d); }
    </style>
  </head>
  <body class="bg-stone-50 text-slate-900 font-sans">
    <main class="max-w-5xl mx-auto px-6 py-12 space-y-12">
      <header>...</header>
      <section id="summary">...</section>
      <section id="findings" class="space-y-10">...</section>
      <section id="fix-first">...</section>
    </main>
  </body>
</html>
```

## Header

Repo name, date, scope (which paths/commits were audited), and a compact legend: red solid line = tainted/attacker-controlled path, red dashed line = trust boundary crossed without a check, dark card = critical. No introduction paragraph — straight into the summary.

## Summary strip

Before the findings, one row of counts as small stat blocks — not a diagram, just numbers: `Critical · High · Medium · Speculative`. Each count in its severity color, `text-3xl font-serif`, label beneath in `text-xs uppercase tracking-wider`. This is the only place raw counts appear; don't repeat them elsewhere.

## Finding card

Each finding is one `<article>`, bordered on the left with a severity-colored bar (`border-l-4`, color per severity below) instead of a full badge-only treatment — the eye should catch severity from the left edge scanning down the page.

- **Title** — names the defect plainly (e.g. "Unescaped user input in order-note SQL query"). Not the category — the actual thing.
- **Severity badge** — `Critical` (red/`bg-red-600`), `High` (orange/`bg-orange-500`), `Medium` (amber/`bg-amber-400 text-slate-900`), `Speculative` (slate/`bg-slate-300 text-slate-700`). Plus a small category tag: `bug`, `security`, `footgun`.
- **Location** — monospaced `file:line`, `font-mono text-sm`. Link-styled even though it's not clickable — it should read as a precise coordinate, not prose.
- **Trigger** — one sentence, the reachability test made explicit. What input or actor causes this. If you can't write this sentence concretely, the finding shouldn't be here (or belongs in Speculative).
- **Diagram** — the centerpiece for Critical and High findings only. Not required for Medium/Speculative — a code snippet is enough there (see below). See patterns below.
- **Snippet** — the offending code, `<pre class="text-xs font-mono bg-slate-900 text-slate-100 rounded p-3 overflow-x-auto">`, 3-8 lines, no more. Highlight the exact bad line with a `bg-red-900/40` span if feasible.
- **Fix** — one sentence, plain English, what changes. Not a diff — a description.
- **Impact** — bullets, ≤6 words each, concrete. e.g. "Any authenticated user reads any order", "Retry storm on partner outage", "Silent data loss on overflow." Not "could be a problem."
- **ADR/context callout** (if applicable) — one line, amber-tinted box, only when relevant.

No paragraphs of explanation. If the diagram needs a paragraph to be understood, redraw the diagram. If a finding needs a paragraph to justify its severity, the severity is probably wrong.

## Diagram patterns

Use a diagram for Critical/High findings where the *path* matters — where something travels from an untrusted source to a dangerous sink. Skip it when the bug is purely local (a wrong comparison operator doesn't need a diagram; a privilege-escalation path does).

### Mermaid flow (the workhorse for attacker-controlled paths)

Use when the point is "input enters here, and nothing stops it before it reaches there." Style the tainted path red end-to-end with `classDef`; style the point where a check *should* exist but doesn't with a dashed red node.

```html
<div class="rounded-lg border border-red-200 bg-white p-4">
  <pre class="mermaid">
    flowchart LR
      A[HTTP request] -->|"note field, no escaping"| B[OrderHandler]
      B --> C[BuildQuery]
      C -.no parameterization.-> D[(orders table)]
      classDef tainted stroke:#dc2626,stroke-width:2px;
      class A,B,C,D tainted
  </pre>
</div>
```

Sequence diagrams work well for race conditions and TOCTOU bugs — show the two actors' timelines interleaving and mark the window where the invariant breaks.

### Hand-built boxes-and-arrows (trust boundary crossings)

Modules as `<div>`s with borders; the trust boundary itself as a vertical dashed red line (`border-l-2 border-dashed border-red-500`) running down the diagram. Arrows crossing the line without passing through a visible "check" box are the finding — make the missing check conspicuous by its absence (a faded, dashed ghost box labeled "no check here").

### Timeline (race conditions, TOCTOU)

Two horizontal lanes, one per actor/thread, time flowing left to right. Mark the check with a small dot, the act with another dot, and shade the gap between them red — that gap is the window.

### Blast radius (authz gaps, IDOR)

One central node (the vulnerable check), with a fan of reachable resources around it. Before: fan is wide (anything reachable). After (in the Fix): fan is narrowed to what should be reachable, with the excluded resources greyed out.

## Style guidance

- Lean editorial, not corporate-dashboard. Generous whitespace. Serif optional for headings (`font-serif` works well with stone/slate) — but severity badges and code stay sans/mono, never serif.
- Colour carries meaning here more than in a general report: red is reserved for Critical/High and for tainted paths — don't decorate with it elsewhere.
- Keep diagrams ~280px tall — this report is longer (more findings) than an architecture report, so diagrams need to stay compact to avoid an endless scroll.
- Use `text-xs uppercase tracking-wider` for node labels inside diagrams — schematic, not UI chrome.
- Order findings Critical → High → Medium → Speculative, full stop. Don't group by file or category first; severity is the spine of this report.
- The only scripts are the Tailwind CDN and the Mermaid ESM import. Static otherwise — no app code, no interactivity beyond Mermaid's own rendering.

## Fix first section

One larger card, same left-border treatment as a Critical finding. Names the one finding to patch before anything else ships, one sentence why (usually: cheapest fix, worst reachable damage), anchor link to its full card.

## Tone

Plain English, concise, concrete nouns. Say what the bug does, not how bad it feels.

**Use exactly:** trigger, reachable, severity, tainted path, trust boundary, sink, check.

**Never substitute:** "issue" or "concern" for a graded finding (name the severity instead) · "best practice violation" (say what breaks) · "potential vulnerability" without completing the reachability test first — if you can't complete it, it's Speculative, say so plainly rather than hedging in prose.

**Phrasings that fit the style:**

- "User-controlled `note` field reaches the SQL string unescaped."
- "No check between token issuance and token use — 400ms window."
- "Any authenticated user can read any tenant's export by guessing the ID."
- "Fix: parameterize the query."

**Impact bullets** stay concrete and reachable: *"any user reads any order"*, *"retry storm doubles load per outage"*, *"silent overflow drops the last digit"*. Don't write *"security risk"* or *"could cause issues"* — name the actual consequence.

No hedging, no throat-clearing, no "it's worth noting that…". If a sentence could be a bullet, make it a bullet. If a bullet could be cut, cut it. If you're tempted to write "may" or "could potentially," either complete the reachability test and drop the hedge, or mark it Speculative and stop apologizing for it.