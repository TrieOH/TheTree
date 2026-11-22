import { createFileRoute } from "@tanstack/solid-router";
import { PageHeader, Card } from "@trieoh/ui-solid";
import { routerLink } from "@/shared/ui/router-link";
import { For, Show } from "solid-js";
import { SERVICES } from "@/features/catalog/services";
import { DOC_CATEGORIES, docsIn } from "@/features/docs/catalog";

export const Route = createFileRoute("/docs/")({
  head: () => ({ meta: [{ title: "Docs · TrieOH dev" }] }),
  component: DocsIndex,
});

function DocsIndex() {
  return (
    <>
      <PageHeader title="Docs" description="Onboarding, service docs, runbooks, decisions and incident reports, from the repo's docs/ folder." />
      <div class="grid gap-3 sm:grid-cols-2">
        <For each={DOC_CATEGORIES}>
          {(category) => (
            <>
              <Card href={`/docs/${category.id}`} renderLink={routerLink} class="p-4">
                <span class="flex items-baseline justify-between gap-3">
                  <span class="font-medium">{category.label}</span>
                  <span class="text-[12px] tabular-nums text-muted-foreground">{docsIn(category.id).length}</span>
                </span>
                <span class="mt-1 block text-[13px] text-muted-foreground">{category.description}</span>
              </Card>
              {/* Service docs come right after onboarding, as in the sidebar. */}
              <Show when={category.id === "onboarding"}>
                <Card href="/docs/services" renderLink={routerLink} class="p-4">
                  <span class="flex items-baseline justify-between gap-3">
                    <span class="font-medium">Services</span>
                    <span class="text-[12px] tabular-nums text-muted-foreground">{SERVICES.length}</span>
                  </span>
                  <span class="mt-1 block text-[13px] text-muted-foreground">How to use each service, and how each one works inside.</span>
                </Card>
              </Show>
            </>
          )}
        </For>
      </div>
    </>
  );
}
