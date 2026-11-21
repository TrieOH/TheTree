import { EmptyState, PageSection, Sparkline, StatusDot, Table, type StatusDotTone } from "@trieoh/ui-solid";
import { For, Show, createMemo } from "solid-js";
import { commitUrl, runUrl, tagUrl, type Deploy } from "../model";
import { dailyDeploys, timeAgo } from "../stats";

const external = { target: "_blank", rel: "noopener" } as const;

const KINDS = [
  { id: "api", title: "APIs", match: (c: string) => !c.endsWith("-ui") && !c.startsWith("pages/") },
  { id: "front", title: "Fronts", match: (c: string) => c.endsWith("-ui") },
  { id: "pages", title: "Pages", match: (c: string) => c.startsWith("pages/") },
];

const TONE: Record<Deploy["status"], { tone: StatusDotTone; label: string }> = {
  success: { tone: "success", label: "Deployed" },
  failure: { tone: "destructive", label: "Last deploy failed" },
  cancelled: { tone: "muted", label: "Last deploy cancelled" },
};

/** What runs in production now, one row per component, with how often each ships. */
export function Production(props: { production: readonly Deploy[]; deploys: readonly Deploy[] }) {
  const groups = createMemo(
    () =>
      KINDS.map((kind) => ({
        ...kind,
        rows: props.production
          .filter((d) => kind.match(d.component))
          .map((d) => ({ ...d, daily: dailyDeploys(props.deploys, d.component) })),
      })).filter((g) => g.rows.length > 0),
    { name: "delivery-production" },
  );

  return (
    <Show
      when={props.production.length > 0}
      fallback={<EmptyState variant="inline" title="No deploys recorded yet." />}
    >
      <p class="mb-6 max-w-[68ch] text-[13px] text-muted-foreground">
        The version each component runs, and its successful deploys per day over the last 30 days.
      </p>
      <For each={groups()}>
        {(group) => (
          <PageSection title={group.title}>
            <Table
              caption={`${group.title} in production`}
              class="text-[13.5px]"
              items={group.rows}
              columns={[
                {
                  id: "component",
                  header: "Component",
                  class: "font-medium",
                  cell: (d) => (
                    <span class="flex items-center gap-2.5">
                      <StatusDot tone={TONE[d.status].tone} label={TONE[d.status].label} size="md" />
                      {group.id === "pages" ? d.component.slice("pages/".length) : d.component}
                    </span>
                  ),
                },
                {
                  id: "version",
                  header: "Version",
                  class: "font-mono text-[12.5px]",
                  title: (d) => `Commit ${d.sha.slice(0, 7)}`,
                  cell: (d) => (
                    <a class="hover:underline" href={tagUrl(d.tag)} {...external}>
                      {d.tag.slice(d.component.length + 1) || d.tag}
                    </a>
                  ),
                },
                {
                  id: "shipped",
                  header: "Shipped",
                  class: "whitespace-nowrap text-muted-foreground",
                  cell: (d) => (
                    <Show when={d.runId} fallback={timeAgo(d.at)}>
                      {(id) => (
                        <a class="hover:underline" href={runUrl(id())} {...external}>
                          {timeAgo(d.at)}
                        </a>
                      )}
                    </Show>
                  ),
                },
                {
                  id: "by",
                  header: "By",
                  class: "text-muted-foreground",
                  cell: (d) => (
                    <a class="hover:underline" href={commitUrl(d.sha)} {...external}>
                      {d.actor}
                    </a>
                  ),
                },
                {
                  id: "cadence",
                  header: "Last 30 days",
                  class: "w-[34%] min-w-40",
                  cell: (d) => (
                    <span class="flex items-center gap-3">
                      <Sparkline
                        type="bar"
                        fill
                        data={d.daily}
                        height={18}
                        label={`${d.component}: deploys per day, last 30 days`}
                        valueFormatter={(n) => `${n} ${n === 1 ? "deploy" : "deploys"}`}
                        color="color-mix(in oklab, var(--muted-foreground) 55%, transparent)"
                        class="min-w-0 flex-1"
                      />
                      <span class="w-6 shrink-0 text-right text-[12px] tabular-nums text-muted-foreground">
                        {d.daily.reduce((a, b) => a + b, 0)}
                      </span>
                    </span>
                  ),
                },
              ]}
            />
          </PageSection>
        )}
      </For>
    </Show>
  );
}
