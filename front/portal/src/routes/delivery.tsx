import { createFileRoute } from "@tanstack/solid-router";
import { Tabs, cn } from "@trieoh/ui-solid";
import { For, Match, Show, Switch, createMemo, createSignal } from "solid-js";
import { loadDelivery } from "@/features/delivery/api";
import {
  commitUrl,
  pullUrl,
  runUrl,
  tagUrl,
  type Deploy,
  type PullRequest,
  type WorkflowRun,
} from "@/features/delivery/model";
import {
  changeFailureRate,
  deploysPerWeek,
  formatDuration,
  prLeadTime,
  timeAgo,
  timeToRestore,
  typicalDuration,
} from "@/features/delivery/stats";
import { Kpi, KpiRow } from "@/widgets/shell/Kpi";
import { PageHeader, Section } from "@/widgets/shell/PageHeader";

export const Route = createFileRoute("/delivery")({
  loader: () => loadDelivery(),
  head: () => ({ meta: [{ title: "Delivery · TrieOH dev" }] }),
  component: DeliveryPage,
});

const TABS = [
  { id: "production", label: "Production" },
  { id: "pulls", label: "Pull requests" },
  { id: "ci", label: "CI" },
  { id: "tags", label: "Tags" },
] as const;

type Tab = (typeof TABS)[number]["id"];

const external = { target: "_blank", rel: "noopener" } as const;

function DeliveryPage() {
  const data = Route.useLoaderData();
  const [tab, setTab] = createSignal<Tab>("production", { name: "delivery-tab" });
  const cfr = () => changeFailureRate(data().deploys);

  return (
    <>
      <PageHeader
        title="Delivery"
        description="What's running, what's about to ship and how CI is doing. From GitHub and the deploy workflows."
      />

      <KpiRow>
        <Kpi label="Deploys / week" note="last 4 weeks">{deploysPerWeek(data().deploys).toFixed(1)}</Kpi>
        <Kpi label="PR lead time" note="median, 30 days">{formatDuration(prLeadTime(data().pulls))}</Kpi>
        <Kpi label="Change failure rate" note="30 days">{cfr() === null ? "—" : `${Math.round(cfr()! * 100)}%`}</Kpi>
        <Kpi label="Time to restore" note="median, 30 days">{formatDuration(timeToRestore(data().deploys))}</Kpi>
      </KpiRow>

      <div class="mt-8">
        <Tabs ariaLabel="Delivery views" items={TABS} value={tab()} onChange={setTab} />
      </div>

      <div class="mt-6">
        <Switch>
          <Match when={tab() === "production"}>
            <Production deploys={data().production} />
          </Match>
          <Match when={tab() === "pulls"}>
            <Pulls pulls={data().pulls} />
          </Match>
          <Match when={tab() === "ci"}>
            <Ci runs={data().runs} />
          </Match>
          <Match when={tab() === "tags"}>
            <Show when={data().tags.length > 0} fallback={<Empty>No tags yet.</Empty>}>
              <ul class="divide-y divide-border border-y border-border">
                <For each={data().tags}>
                  {(tag) => (
                    <li class="flex items-baseline gap-3 py-2.5 text-[13.5px]">
                      <a class="shrink-0 font-mono text-[12.5px] hover:underline" href={tagUrl(tag.name)} {...external}>
                        {tag.name}
                      </a>
                      <span class="min-w-0 flex-1 truncate text-muted-foreground">{tag.message}</span>
                      <span class="shrink-0 text-[12px] text-muted-foreground">{timeAgo(tag.createdAt)}</span>
                    </li>
                  )}
                </For>
              </ul>
            </Show>
          </Match>
        </Switch>
      </div>
    </>
  );
}

function Empty(props: { children: string }) {
  return <p class="text-[13px] text-muted-foreground">{props.children}</p>;
}

function Pill(props: { tone: "ok" | "bad" | "warn" | "muted"; children: string }) {
  return (
    <span
      class={cn(
        "inline-flex shrink-0 items-center rounded-full px-2 py-0.5 text-[11px] font-medium",
        props.tone === "ok" && "bg-ok/15 text-ok",
        props.tone === "bad" && "bg-destructive/20 text-destructive dark:text-destructive-foreground",
        props.tone === "warn" && "bg-warn/15 text-warn",
        props.tone === "muted" && "bg-muted text-muted-foreground",
      )}
    >
      {props.children}
    </span>
  );
}

function deployPill(deploy: Deploy): { tone: "ok" | "bad" | "warn" | "muted"; label: string } {
  if (deploy.status === "cancelled") return { tone: "muted", label: "Cancelled" };
  if (deploy.status === "failure") return { tone: "bad", label: "Failed" };
  if (deploy.healthOk === false) return { tone: "bad", label: "Unhealthy" };
  return { tone: "ok", label: "Live" };
}

function Production(props: { deploys: readonly Deploy[] }) {
  return (
    <Show when={props.deploys.length > 0} fallback={<Empty>No deploys recorded yet.</Empty>}>
      <div class="overflow-x-auto">
        <table class="w-full text-left text-[13.5px]">
          <thead class="text-[12px] text-muted-foreground">
            <tr class="border-b border-border">
              <th class="py-2 pr-4 font-normal">Component</th>
              <th class="py-2 pr-4 font-normal">Tag</th>
              <th class="py-2 pr-4 font-normal">Status</th>
              <th class="py-2 pr-4 font-normal">Deployed</th>
              <th class="py-2 font-normal">By</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <For each={props.deploys}>
              {(deploy) => {
                const pill = deployPill(deploy);
                return (
                  <tr>
                    <td class="py-2.5 pr-4 font-medium">{deploy.component}</td>
                    <td class="py-2.5 pr-4 font-mono text-[12.5px]">
                      <a class="hover:underline" href={commitUrl(deploy.sha)} {...external}>
                        {deploy.tag.slice(deploy.component.length + 1) || deploy.tag}
                      </a>
                    </td>
                    <td class="py-2.5 pr-4">
                      <Show when={deploy.runId} fallback={<Pill tone={pill.tone}>{pill.label}</Pill>}>
                        {(id) => (
                          <a href={runUrl(id())} {...external}>
                            <Pill tone={pill.tone}>{pill.label}</Pill>
                          </a>
                        )}
                      </Show>
                    </td>
                    <td class="whitespace-nowrap py-2.5 pr-4 text-muted-foreground">{timeAgo(deploy.at)}</td>
                    <td class="py-2.5 text-muted-foreground">{deploy.actor}</td>
                  </tr>
                );
              }}
            </For>
          </tbody>
        </table>
      </div>
    </Show>
  );
}

function reviewPill(pr: PullRequest): { tone: "ok" | "bad" | "warn" | "muted"; label: string } {
  if (pr.draft) return { tone: "muted", label: "Draft" };
  if (pr.reviewState === "approved") return { tone: "ok", label: "Approved" };
  if (pr.reviewState === "changes_requested") return { tone: "bad", label: "Changes requested" };
  return { tone: "warn", label: "Needs review" };
}

function Pulls(props: { pulls: readonly PullRequest[] }) {
  const open = () => props.pulls.filter((pr) => pr.state === "open");
  const merged = () => props.pulls.filter((pr) => pr.state === "merged");
  return (
    <>
      <Section
        title="Open"
        aside={
          <a href="https://github.com/TrieOH/TheTree/pulls" {...external}>
            All on GitHub
          </a>
        }
      >
        <PullList pulls={open()} empty="No open pull requests." pill={reviewPill} />
      </Section>
      <Section title="Merged" aside="last 30 days">
        <PullList
          pulls={merged()}
          empty="Nothing merged in the last 30 days."
          pill={(pr) => ({ tone: "muted", label: timeAgo(pr.mergedAt!) })}
        />
      </Section>
    </>
  );
}

function PullList(props: {
  pulls: readonly PullRequest[];
  empty: string;
  pill: (pr: PullRequest) => { tone: "ok" | "bad" | "warn" | "muted"; label: string };
}) {
  return (
    <Show when={props.pulls.length > 0} fallback={<Empty>{props.empty}</Empty>}>
      <ul class="divide-y divide-border border-y border-border">
        <For each={props.pulls}>
          {(pr) => {
            const pill = props.pill(pr);
            return (
              <li>
                <a
                  href={pullUrl(pr.number)}
                  {...external}
                  class="flex items-center gap-3 py-2.5 text-[13.5px] hover:bg-foreground/5"
                >
                  <span class="w-12 shrink-0 font-mono text-[12px] text-muted-foreground">#{pr.number}</span>
                  <span class="min-w-0 flex-1 truncate">{pr.title}</span>
                  <span class="hidden shrink-0 text-[12px] text-muted-foreground sm:inline">{pr.author}</span>
                  <Pill tone={pill.tone}>{pill.label}</Pill>
                </a>
              </li>
            );
          }}
        </For>
      </ul>
    </Show>
  );
}

function runTone(run: WorkflowRun): string {
  if (run.conclusion === null) return "bg-warn/70";
  if (run.conclusion === "success") return run.attempt > 1 ? "bg-ok/50" : "bg-ok";
  if (run.conclusion === "failure" || run.conclusion === "timed_out") return "bg-destructive-foreground";
  return "bg-muted-foreground/40";
}

function runLabel(run: WorkflowRun): string {
  const state = run.conclusion ?? "running";
  const retry = run.attempt > 1 ? `, attempt ${run.attempt}` : "";
  return `${run.branch}: ${state}${retry}`;
}

function Ci(props: { runs: readonly WorkflowRun[] }) {
  const workflows = createMemo(() => {
    const byName = new Map<string, WorkflowRun[]>();
    for (const run of props.runs) byName.set(run.workflow, [...(byName.get(run.workflow) ?? []), run]);
    return [...byName.entries()].map(([name, runs]) => ({ name, runs }));
  });
  return (
    <Show when={workflows().length > 0} fallback={<Empty>No workflow runs recorded yet.</Empty>}>
      <p class="mb-4 text-[13px] text-muted-foreground">
        Last 14 runs of each workflow, oldest first. CI runs on release tags; a lighter square passed on a retry.
      </p>
      <ul class="divide-y divide-border border-y border-border">
        <For each={workflows()}>
          {(workflow) => (
            <li class="flex flex-wrap items-center gap-x-4 gap-y-2 py-3">
              <span class="w-32 shrink-0 font-medium">{workflow.name}</span>
              <span class="flex flex-1 gap-1">
                <For each={workflow.runs}>
                  {(run) => (
                    <a
                      href={runUrl(run.id)}
                      {...external}
                      title={runLabel(run)}
                      aria-label={runLabel(run)}
                      class={cn("h-5 w-3 rounded-[3px]", runTone(run))}
                    />
                  )}
                </For>
              </span>
              <span class="shrink-0 text-[12px] tabular-nums text-muted-foreground">
                typically {formatDuration(typicalDuration(workflow.runs))}
              </span>
            </li>
          )}
        </For>
      </ul>
    </Show>
  );
}
