import { Link } from "@tanstack/solid-router";
import { PageSection } from "@trieoh/ui-solid";
import { For, Show } from "solid-js";
import type { ServiceId } from "@/features/catalog/services";
import { AUDIENCES, decisionsFor, groupsOf, incidentsFor } from "../service-docs";
import type { DocLink } from "../catalog";
import { DocStatus } from "../DocStatus";
import { NotConnected } from "@/widgets/shell/NotConnected";

const rowLink =
  "block rounded-sm py-3 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring";

/**
 * A service's written docs, consumer then maintainer, under their groups.
 * Both sections always show, so a missing half reads as "not written yet".
 */
export function ServiceDocIndex(props: { service: ServiceId }) {
  return (
    <For each={AUDIENCES}>
      {(audience) => (
        <PageSection title={audience.label} aside={audience.description}>
          <Show
            when={groupsOf(props.service, audience.id).length > 0}
            fallback={
              <NotConnected source={`docs/services/${props.service}/${audience.id}/`}>
                No {audience.label.toLowerCase()} pages yet.
              </NotConnected>
            }
          >
            <div class="grid gap-6">
              <For each={groupsOf(props.service, audience.id)}>
                {(group) => (
                  <section>
                    <Show when={group.name}>
                      <h3 class="mb-1 text-[13px] font-medium text-muted-foreground">{group.name}</h3>
                    </Show>
                    <ul class="divide-y divide-border border-y border-border">
                      <For each={group.docs}>
                        {(doc) => (
                          <li>
                            <Link
                              to="/services/$service/docs/$slug"
                              params={{ service: props.service, slug: doc.slug }}
                              class={`group ${rowLink}`}
                            >
                              <span class="block font-medium group-hover:text-primary">{doc.meta.title ?? doc.slug}</span>
                              <Show when={doc.meta.summary}>
                                <span class="mt-0.5 block text-[13px] text-muted-foreground">{doc.meta.summary}</span>
                              </Show>
                            </Link>
                          </li>
                        )}
                      </For>
                    </ul>
                  </section>
                )}
              </For>
            </div>
          </Show>
        </PageSection>
      )}
    </For>
  );
}

/** ADRs and incident reports that name the service, found from their headers. */
export function ServiceRelated(props: { service: ServiceId }) {
  return (
    <>
      <Show when={decisionsFor(props.service).length > 0}>
        <PageSection title="Decisions" aside="ADRs whose scope names this service">
          <RelatedList docs={decisionsFor(props.service)} />
        </PageSection>
      </Show>
      <Show when={incidentsFor(props.service).length > 0}>
        <PageSection title="Incident reports" aside="reports that list this service">
          <RelatedList docs={incidentsFor(props.service)} />
        </PageSection>
      </Show>
    </>
  );
}

function RelatedList(props: { docs: readonly DocLink[] }) {
  return (
    <ul class="divide-y divide-border border-y border-border">
      <For each={props.docs}>
        {(doc) => (
          <li>
            <Link
              to="/docs/$category/$slug"
              params={{ category: doc.category, slug: doc.slug }}
              class={`group flex items-baseline gap-3 ${rowLink}`}
            >
              <Show when={doc.meta.number}>
                <span class="shrink-0 font-mono text-[12px] text-muted-foreground">{doc.meta.number}</span>
              </Show>
              <span class="min-w-0 flex-1 group-hover:text-primary">{doc.meta.title ?? doc.slug}</span>
              <DocStatus status={doc.meta.status} />
            </Link>
          </li>
        )}
      </For>
    </ul>
  );
}
