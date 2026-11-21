import { createFileRoute } from "@tanstack/solid-router";
import { Tabs, cn, PageHeader, PageSection, DashboardStatCard, Badge, Sparkline, Table } from "@trieoh/ui-solid";
import { For, Match, Show, Switch, createMemo, createSignal } from "solid-js";
import { loadDelivery } from "@/features/delivery/api";
import {
  commitUrl,
  pullUrl,
  runUrl,
  tagUrl,
  type CiPayload,
  type CiRun,
  type CiStepStat,
  type Deploy,
  type PullRequest,
  type WorkflowRun,
} from "@/features/delivery/model";
import {
  changeFailureRate,
  deploysPerWeek,
  formatDuration,
  formatMetric,
  formatStepTime,
  prLeadTime,
  timeAgo,
  timeToRestore,
  typicalDuration,
} from "@/features/delivery/stats";
import { KpiRow } from "@/widgets/shell/Kpi";

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
        <DashboardStatCard label="Deploys / week" hint="last 4 weeks" value={deploysPerWeek(data().deploys).toFixed(1)} />
        <DashboardStatCard label="PR lead time" hint="median, 30 days" value={formatDuration(prLeadTime(data().pulls))} />
        <DashboardStatCard
          label="Change failure rate"
          hint="30 days"
          value={cfr() === null ? "—" : `${Math.round(cfr()! * 100)}%`}
        />
        <DashboardStatCard label="Time to restore" hint="median, 30 days" value={formatDuration(timeToRestore(data().deploys))} />
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
            <Pipelines ci={data().ci} />
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

type Tone = "success" | "destructive" | "warning" | "muted";

function deployPill(deploy: Deploy): { tone: Tone; label: string } {
  if (deploy.status === "cancelled") return { tone: "muted", label: "Cancelled" };
  if (deploy.status === "failure") return { tone: "destructive", label: "Failed" };
  return { tone: "success", label: "Deployed" };
}

function Production(props: { deploys: readonly Deploy[] }) {
  return (
    <Show when={props.deploys.length > 0} fallback={<Empty>No deploys recorded yet.</Empty>}>
      <Table
        caption="Production deploys"
        class="text-[13.5px]"
        items={props.deploys}
        columns={[
          { id: "component", header: "Component", class: "font-medium", cell: (deploy) => deploy.component },
          {
            id: "tag",
            header: "Tag",
            class: "font-mono text-[12.5px]",
            cell: (deploy) => (
              <a class="hover:underline" href={commitUrl(deploy.sha)} {...external}>
                {deploy.tag.slice(deploy.component.length + 1) || deploy.tag}
              </a>
            ),
          },
          {
            id: "status",
            header: "Status",
            cell: (deploy) => {
              const pill = deployPill(deploy);
              const badge = (
                <Badge shape="pill" variant={pill.tone}>
                  {pill.label}
                </Badge>
              );
              return (
                <Show when={deploy.runId} fallback={badge}>
                  {(id) => (
                    <a href={runUrl(id())} {...external}>
                      {badge}
                    </a>
                  )}
                </Show>
              );
            },
          },
          {
            id: "deployed",
            header: "Deployed",
            class: "whitespace-nowrap text-muted-foreground",
            cell: (deploy) => timeAgo(deploy.at),
          },
          { id: "by", header: "By", class: "text-muted-foreground", cell: (deploy) => deploy.actor },
        ]}
      />
    </Show>
  );
}

function reviewPill(pr: PullRequest): { tone: Tone; label: string } {
  if (pr.draft) return { tone: "muted", label: "Draft" };
  if (pr.reviewState === "approved") return { tone: "success", label: "Approved" };
  if (pr.reviewState === "changes_requested") return { tone: "destructive", label: "Changes requested" };
  return { tone: "warning", label: "Needs review" };
}

function Pulls(props: { pulls: readonly PullRequest[] }) {
  const open = () => props.pulls.filter((pr) => pr.state === "open");
  const merged = () => props.pulls.filter((pr) => pr.state === "merged");
  return (
    <>
      <PageSection
        title="Open"
        aside={
          <a href="https://github.com/TrieOH/TheTree/pulls" {...external}>
            All on GitHub
          </a>
        }
      >
        <PullList pulls={open()} empty="No open pull requests." pill={reviewPill} />
      </PageSection>
      <PageSection title="Merged" aside="last 30 days">
        <PullList
          pulls={merged()}
          empty="Nothing merged in the last 30 days."
          pill={(pr) => ({ tone: "muted", label: timeAgo(pr.mergedAt!) })}
        />
      </PageSection>
    </>
  );
}

function PullList(props: {
  pulls: readonly PullRequest[];
  empty: string;
  pill: (pr: PullRequest) => { tone: Tone; label: string };
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
                  <Badge shape="pill" variant={pill.tone}>{pill.label}</Badge>
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

const PIPELINES: Record<string, string> = { service: "Services", front: "Fronts", sdk: "SDKs", pages: "Pages" };

function stagePill(run: CiRun): { tone: Tone; label: string } {
  if (run.status === "cancelled") return { tone: "muted", label: "Cancelled" };
  if (run.status === "failure") return run.stage === "publish" ? { tone: "destructive", label: "Failed" } : { tone: "warning", label: `Blocked at ${run.stage}` };
  return { tone: "success", label: run.type === "sdk" ? "Gate passed" : "Shipped" };
}


function Pipelines(props: { ci: CiPayload }) {
  const types = createMemo(() => Object.keys(PIPELINES).filter((type) => props.ci.runs.some((run) => run.type === type)));
  const [pipeline, setPipeline] = createSignal<string>("service", { name: "delivery-pipeline" });
  const current = () => (types().includes(pipeline()) ? pipeline() : (types()[0] ?? "service"));
  const steps = createMemo(() => props.ci.steps.filter((step) => step.type === current()));
  const metrics = createMemo(() => props.ci.metrics.filter((metric) => metric.type === current()));
  const stepLabel = (step: CiStepStat) => `${step.job} · ${step.step}`;

  return (
    <Show when={props.ci.runs.length > 0} fallback={<p class="mt-8 text-[13px] text-muted-foreground">No step timings reported yet.</p>}>
      <div class="mt-10">
        <Tabs
          ariaLabel="Pipelines"
          items={types().map((type) => ({ id: type, label: PIPELINES[type] }))}
          value={current()}
          onChange={setPipeline}
        />
      </div>

      <div class="mt-6">
        <PageSection title="Step times" aside={`last ${Math.max(0, ...steps().map((s) => s.runs))} runs`}>
          <Table
            caption="Step times"
            class="text-[13.5px]"
            density="compact"
            items={steps()}
            columns={[
              {
                id: "step",
                header: "Step",
                cell: (step) => (
                  <>
                    <span class="text-muted-foreground">{step.job} · </span>
                    {step.step}
                  </>
                ),
              },
              {
                id: "median",
                header: "Median",
                align: "right",
                class: "tabular-nums",
                cell: (step) => formatStepTime(step.medianMs),
              },
              {
                id: "p90",
                header: "p90",
                align: "right",
                class: "tabular-nums text-muted-foreground",
                cell: (step) => formatStepTime(step.p90Ms),
              },
              {
                id: "last",
                header: "Last",
                align: "right",
                class: "tabular-nums text-muted-foreground",
                cell: (step) => formatStepTime(step.lastMs),
              },
              {
                id: "recent",
                header: "Recent",
                title: stepLabel,
                cell: (step) => (
                  <Sparkline
                    type="bar"
                    data={step.recentMs}
                    label={`${stepLabel(step)}, recent runs`}
                    height={20}
                    color="color-mix(in oklab, var(--muted-foreground) 50%, transparent)"
                  />
                ),
              },
            ]}
          />
        </PageSection>

        <PageSection title="Measurements" aside="from .github/scripts/metric">
          <Show when={metrics().length > 0} fallback={<Empty>No measurements for this pipeline yet.</Empty>}>
            <Table
              caption="Measurements"
              class="text-[13.5px]"
              density="compact"
              items={metrics()}
              columns={[
                {
                  id: "name",
                  header: "Measurement",
                  cell: (metric) => (
                    <>
                      <span class="text-muted-foreground">
                        {metric.job} · {metric.source} ·{" "}
                      </span>
                      {metric.name}
                      <For each={Object.entries(metric.labels)}>
                        {([key, value]) => (
                          <span class="ml-2 rounded bg-muted px-1.5 py-0.5 font-mono text-[11px] text-muted-foreground">
                            {key}={value}
                          </span>
                        )}
                      </For>
                    </>
                  ),
                },
                {
                  id: "last",
                  header: "Last",
                  align: "right",
                  class: "tabular-nums",
                  cell: (metric) => formatMetric(metric.name, metric.last),
                },
                {
                  id: "median",
                  header: "Median",
                  align: "right",
                  class: "tabular-nums text-muted-foreground",
                  cell: (metric) => formatMetric(metric.name, metric.median),
                },
                {
                  id: "runs",
                  header: "Runs",
                  align: "right",
                  class: "tabular-nums text-muted-foreground",
                  cell: (metric) => metric.runs,
                },
              ]}
            />
          </Show>
        </PageSection>

        <PageSection title="Release runs" aside="newest first, all pipelines">
          <ul class="divide-y divide-border border-y border-border">
            <For each={props.ci.runs}>
              {(run) => {
                const pill = stagePill(run);
                const ms = run.startedAt && run.completedAt ? Date.parse(run.completedAt) - Date.parse(run.startedAt) : null;
                return (
                  <li>
                    <a href={runUrl(run.runId)} {...external} class="flex items-center gap-3 py-2.5 text-[13.5px] hover:bg-foreground/5">
                      <span class="min-w-0 flex-1 truncate font-mono text-[12.5px]">{run.tag}</span>
                      <Show when={run.lane === "express"}>
                        <Badge shape="pill" variant="muted">hotfix</Badge>
                      </Show>
                      <span class="hidden shrink-0 text-[12px] tabular-nums text-muted-foreground sm:inline">
                        {formatDuration(ms)}
                      </span>
                      <span class="hidden w-20 shrink-0 text-right text-[12px] text-muted-foreground sm:inline">
                        {run.startedAt ? timeAgo(run.startedAt) : ""}
                      </span>
                      <Badge shape="pill" variant={pill.tone}>{pill.label}</Badge>
                    </a>
                  </li>
                );
              }}
            </For>
          </ul>
        </PageSection>
      </div>
    </Show>
  );
}
