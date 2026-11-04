import { Link, createFileRoute } from "@tanstack/solid-router";
import { For } from "solid-js";
import { docsIn, type DocCategoryId } from "@/features/docs/catalog";
import { NotConnected, Section } from "@/widgets/shell/PageHeader";

export const Route = createFileRoute("/")({
  component: Home,
});

function greeting(): string {
  const hour = new Date().getHours();
  if (hour < 12) return "Good morning";
  if (hour < 18) return "Good afternoon";
  return "Good evening";
}

const QUICK_LINKS = [
  { category: "decisions", label: "Decisions" },
  { category: "onboarding", label: "Run locally" },
  { category: "runbooks", label: "Deploy runbook" },
] as const;

const chip =
  "inline-flex h-8 items-center rounded-full border border-border px-3 text-[13px] text-muted-foreground transition-colors hover:bg-foreground/5 hover:text-foreground";

function Home() {
  const decisions = docsIn("decisions");

  return (
    <>
      <section class="mb-12">
        <h1 class="font-heading text-[32px] font-semibold leading-tight tracking-tight">{greeting()}</h1>
        <p class="mt-2 text-muted-foreground">Docs, services and the state of TheTree in one place.</p>
        <div class="mt-6 flex flex-wrap gap-2">
          <Link to="/services/$service" params={{ service: "identityx" }} class={chip}>
            IdentityX API
          </Link>
          <For each={QUICK_LINKS}>
            {(link) => (
              <Link to="/docs/$category" params={{ category: link.category }} class={chip}>
                {link.label}
              </Link>
            )}
          </For>
        </div>
      </section>

      <Section title="Needs attention">
        <NotConnected source="Linear, CI and Cloudflare">
          Nothing to flag yet. Urgent issues, failing CI on main and degraded services show up here once
          those sources are connected.
        </NotConnected>
      </Section>

      <Section
        title="Decisions"
        aside={<Link to="/docs/$category" params={{ category: "decisions" }}>All decisions</Link>}
      >
        <ul class="divide-y divide-border border-y border-border">
          <For each={decisions.slice(0, 3)}>
            {(doc) => (
              <li>
                <Link
                  to="/docs/$category/$slug"
                  params={{ category: "decisions", slug: doc.slug }}
                  class="flex items-baseline justify-between gap-4 py-3 hover:text-foreground"
                >
                  <span class="min-w-0">
                    <span class="block font-medium">{doc.meta.title}</span>
                    <span class="text-[13px] text-muted-foreground">{doc.meta.number}</span>
                  </span>
                  <span class="shrink-0 text-[13px] tabular-nums text-muted-foreground">{doc.meta.date}</span>
                </Link>
              </li>
            )}
          </For>
        </ul>
      </Section>

      <Section title="New here?">
        <div class="grid gap-3 sm:grid-cols-3">
          <StartCard category="onboarding" title="Run TheTree locally" note="mise, just and the four APIs" />
          <StartCard category="runbooks" title="How we ship" note="Branches, commits and deploys" />
          <StartCard
            category="decisions"
            title="Why things are this way"
            note={`${decisions.length} decision records`}
          />
        </div>
      </Section>
    </>
  );
}

function StartCard(props: { category: DocCategoryId; title: string; note: string }) {
  return (
    <Link
      to="/docs/$category"
      params={{ category: props.category }}
      class="rounded-xl border border-border bg-card p-4 transition-colors hover:bg-foreground/5"
    >
      <span class="block font-medium">{props.title}</span>
      <span class="mt-1 block text-[13px] text-muted-foreground">{props.note}</span>
    </Link>
  );
}
