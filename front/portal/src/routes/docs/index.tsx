import { Link, createFileRoute } from "@tanstack/solid-router";
import { For } from "solid-js";
import { DOC_CATEGORIES, docsIn } from "@/features/docs/catalog";
import { PageHeader } from "@/widgets/shell/PageHeader";

export const Route = createFileRoute("/docs/")({
  head: () => ({ meta: [{ title: "Docs · TrieOH dev" }] }),
  component: DocsIndex,
});

function DocsIndex() {
  return (
    <>
      <PageHeader title="Docs" description="Decisions, runbooks, incident reports and onboarding, from the repo's docs/ folder." />
      <div class="grid gap-3 sm:grid-cols-2">
        <For each={DOC_CATEGORIES}>
          {(category) => (
            <Link
              to="/docs/$category"
              params={{ category: category.id }}
              class="rounded-xl border border-border bg-card p-4 transition-colors hover:bg-foreground/5"
            >
              <span class="flex items-baseline justify-between gap-3">
                <span class="font-medium">{category.label}</span>
                <span class="text-[12px] tabular-nums text-muted-foreground">{docsIn(category.id).length}</span>
              </span>
              <span class="mt-1 block text-[13px] text-muted-foreground">{category.description}</span>
            </Link>
          )}
        </For>
      </div>
    </>
  );
}
