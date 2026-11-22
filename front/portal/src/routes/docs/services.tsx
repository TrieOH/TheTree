import { createFileRoute } from "@tanstack/solid-router";
import { Card, PageHeader } from "@trieoh/ui-solid";
import { For } from "solid-js";
import { SERVICES } from "@/features/catalog/services";
import { serviceDocs } from "@/features/docs/service-docs";
import { routerLink } from "@/shared/ui/router-link";

export const Route = createFileRoute("/docs/services")({
  head: () => ({ meta: [{ title: "Service docs · TrieOH dev" }] }),
  component: ServiceDocsIndex,
});

/** The Docs tab's way into each service's docs, which live under the service. */
function ServiceDocsIndex() {
  return (
    <>
      <PageHeader title="Services" description="Each service's own docs, for the apps that call it and for whoever maintains it." />
      <div class="grid gap-3 sm:grid-cols-2">
        <For each={SERVICES}>
          {(service) => (
            <Card href={`/services/${service.id}/docs`} renderLink={routerLink} class="p-4">
              <span class="flex items-baseline justify-between gap-3">
                <span class="font-medium">{service.name}</span>
                <span class="text-[12px] tabular-nums text-muted-foreground">{serviceDocs(service.id).length}</span>
              </span>
              <span class="mt-1 block text-[13px] text-muted-foreground">{service.description}</span>
            </Card>
          )}
        </For>
      </div>
    </>
  );
}
