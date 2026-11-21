import { createFileRoute } from "@tanstack/solid-router";
import { PageHeader, Card } from "@trieoh/ui-solid";
import { routerLink } from "@/shared/ui/router-link";
import { For } from "solid-js";
import { DOC_CATEGORIES, docsIn } from "@/features/docs/catalog";

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
            <Card href={`/docs/${category.id}`} renderLink={routerLink} class="p-4">
              <span class="flex items-baseline justify-between gap-3">
                <span class="font-medium">{category.label}</span>
                <span class="text-[12px] tabular-nums text-muted-foreground">{docsIn(category.id).length}</span>
              </span>
              <span class="mt-1 block text-[13px] text-muted-foreground">{category.description}</span>
            </Card>
          )}
        </For>
      </div>
    </>
  );
}
