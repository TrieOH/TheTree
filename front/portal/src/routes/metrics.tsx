import { createFileRoute } from "@tanstack/solid-router";
import {
  Badge,
  ChartCard,
  DashboardStatCard,
  Meter,
  PageHeader,
  PageSection,
  Table,
  Tabs,
  ToggleGroup,
  type ChartDatum,
} from "@trieoh/ui-solid";
import { For, Match, Show, Switch, createMemo, createSignal } from "solid-js";
import { NotConnected } from "@/widgets/shell/NotConnected";
import { KpiRow } from "@/widgets/shell/Kpi";
import { SERVICES, findService } from "@/features/catalog/services";
import { loadMetrics } from "@/features/metrics/api";
import { cyclePace, needsAttention, projectedShare, type MetricsTab } from "@/features/metrics/attention";
import { summarize } from "@/features/metrics/cold-starts";
import { formatBytes, formatCount, formatHours, formatMs, formatPercent } from "@/features/metrics/format";
import { PLAN, type MetricsPayload, type WorkerStats } from "@/features/metrics/model";
import { Attention } from "@/shared/ui/attention";
import { EndpointTable, ProblemTable } from "@/features/metrics/ui/LatencyTables";
import { ServiceLanes } from "@/features/metrics/ui/ServiceLanes";
import { UsageBar } from "@/features/metrics/ui/UsageBar";
import { timeAgo } from "@/features/delivery/stats";

export const Route = createFileRoute("/metrics")({
  loader: () => loadMetrics(),
  head: () => ({ meta: [{ title: "Metrics · TrieOH dev" }] }),
  component: MetricsPage,
});

const TABS: { id: MetricsTab; label: string }[] = [
  { id: "services", label: "APIs" },
  { id: "latency", label: "Latency & errors" },
  { id: "usage", label: "Plan usage" },
  { id: "database", label: "Databases" },
  { id: "workers", label: "Workers" },
  { id: "cold-starts", label: "Cold starts" },
];
const tabLabel = (id: MetricsTab) => TABS.find((tab) => tab.id === id)!.label;

const note = "text-[13px] text-muted-foreground max-w-[68ch]";

function MetricsPage() {
  const data = Route.useLoaderData();
  const [tab, setTab] = createSignal<MetricsTab>("services", { name: "metrics-tab" });
  const attention = createMemo(() => (data().syncedAt ? needsAttention(data()) : []), { name: "metrics-attention" });
  let tabs: HTMLDivElement | undefined;
  const open = (next: MetricsTab) => {
    setTab(next);
    tabs?.scrollIntoView({ behavior: "smooth", block: "start" });
  };

  return (
    <>
      <PageHeader
        title="Metrics"
        description={
          <>
            The four APIs, their Workers and databases, pulled from Cloudflare, Neon and Grafana every 30 minutes.
            <Show when={data().syncedAt}>{(at) => <> Last pulled {timeAgo(at())}.</>}</Show>
          </>
        }
      />
      <Show when={data().syncedAt} fallback={<NotConnected source="Cloudflare and Neon" />}>
        <Attention items={attention()} tabLabel={tabLabel} onOpen={open} />
        <div ref={(el) => (tabs = el)} class="scroll-mt-6">
          <Tabs ariaLabel="Metric views" items={TABS} value={tab()} onChange={setTab} />
        </div>
        <div class="mt-6">
          <Switch>
            <Match when={tab() === "services"}>
              <Services data={data()} />
            </Match>
            <Match when={tab() === "latency"}>
              <LatencyAndErrors data={data()} />
            </Match>
            <Match when={tab() === "usage"}>
              <PlanUsage data={data()} />
            </Match>
            <Match when={tab() === "database"}>
              <Databases data={data()} />
            </Match>
            <Match when={tab() === "workers"}>
              <Workers workers={data().workers} />
            </Match>
            <Match when={tab() === "cold-starts"}>
              <ColdStarts data={data()} />
            </Match>
          </Switch>
        </div>
      </Show>
    </>
  );
}

function Services(props: { data: MetricsPayload }) {
  return (
    <div class="space-y-4">
      <ServiceLanes services={props.data.services} />
      <p class={note}>
        Requests, 5xx and p95 cover the last 24 hours. p95 is measured inside the container, so a cold start's wait is
        not in it. Open an API for its endpoints, problems and releases.
      </p>
    </div>
  );
}

function LatencyAndErrors(props: { data: MetricsPayload }) {
  return (
    <>
      <p class={`${note} mb-6`}>
        Last 7 days, measured inside each API's container. An endpoint's p95 is its worst day; past the tick it counts
        as degraded.
      </p>
      <Show
        when={props.data.endpoints.length > 0 || props.data.problems.length > 0}
        fallback={<NotConnected source="Grafana Cloud">No latency or error data pulled yet.</NotConnected>}
      >
        <PageSection title="Slowest endpoints">
          <EndpointTable endpoints={props.data.endpoints} showService />
        </PageSection>
        <PageSection title="Most common problems" aside="answers sent, by type">
          <ProblemTable problems={props.data.problems} showService />
        </PageSection>
      </Show>
    </>
  );
}

/** Cycles turn over at midnight UTC, so their days are UTC days; the end is exclusive. */
const utcDay = (at: string | number) =>
  new Date(at).toLocaleDateString("en", { month: "short", day: "numeric", timeZone: "UTC" });

/** Below this, a resource is a rounding error and goes in the small table. */
const MINOR_SHARE = 0.01;

function PlanUsage(props: { data: MetricsPayload }) {
  const pace = () => cyclePace(props.data.cycle);
  const daysLeft = () => Math.max(0, Math.ceil((Date.parse(props.data.cycle.end) - Date.now()) / 86_400_000));
  const ranked = createMemo(() =>
    props.data.usage
      .map((usage) => ({ usage, end: projectedShare(usage, pace()) ?? usage.used / usage.included }))
      .sort((a, b) => b.end - a.end),
  );
  const major = () => ranked().filter((row) => row.end >= MINOR_SHARE);
  const minor = () => ranked().filter((row) => row.end < MINOR_SHARE);
  return (
    <div class="space-y-10">
      <Meter
        label="Billing cycle"
        valueLabel={`${utcDay(props.data.cycle.start)} to ${utcDay(Date.parse(props.data.cycle.end) - 1)}, ${daysLeft()} ${daysLeft() === 1 ? "day" : "days"} left`}
        value={pace()}
        tone="default"
      />
      <PageSection title="On the Workers Paid plan" aside="sorted by where each lands at the end">
        <div class="space-y-6">
          <For each={major()}>{(row) => <UsageBar usage={row.usage} pace={pace()} />}</For>
        </div>
      </PageSection>
      <Show when={minor().length > 0}>
        <PageSection title="Barely used" aside="under 1% by the end of the cycle">
          <Table
            caption="Barely used resources"
            density="compact"
            items={minor()}
            columns={[
              { id: "resource", header: "Resource", cell: (row) => PLAN[row.usage.resource].label },
              {
                id: "used",
                header: "Used",
                align: "right",
                class: "tabular-nums",
                cell: (row) => formatCount(row.usage.used),
              },
              {
                id: "included",
                header: "Included",
                align: "right",
                class: "tabular-nums text-muted-foreground",
                cell: (row) => formatCount(row.usage.included),
              },
              {
                id: "end",
                header: "By the end",
                align: "right",
                class: "tabular-nums text-muted-foreground",
                cell: (row) => formatPercent(row.end, 2),
              },
            ]}
          />
        </PageSection>
      </Show>
    </div>
  );
}

function Databases(props: { data: MetricsPayload }) {
  const neon = () => props.data.neon;
  const total = (pick: (p: MetricsPayload["neon"][number]) => number | null) =>
    neon().reduce((n, project) => n + (pick(project) ?? 0), 0);
  const asleep = () =>
    neon().length ? neon().reduce((n, project) => n + (project.suspendedShare ?? 0), 0) / neon().length : null;
  return (
    <div class="space-y-8">
      <p class={note}>
        One Neon project per API, read from Neon's API; nothing here connects to Postgres. Neon's month runs to its own
        quota reset, not the Cloudflare cycle.
      </p>
      <KpiRow>
        <DashboardStatCard label="Stored" value={formatBytes(total((p) => p.storageBytes))} hint="all four" />
        <DashboardStatCard label="Asleep" value={formatPercent(asleep(), 1)} hint="of the month, on average" />
        <DashboardStatCard label="Woke up" value={formatCount(total((p) => p.resumes))} hint="times this month" />
        <DashboardStatCard label="Awake" value={formatHours(total((p) => p.activeSeconds))} hint="compute, all four" />
      </KpiRow>
      <Table
        caption="Neon projects"
        items={neon()}
        columns={[
          { id: "service", header: "API", class: "font-medium", cell: (p) => findService(p.service)?.name },
          {
            id: "storage",
            header: "Stored",
            align: "right",
            class: "tabular-nums",
            cell: (p) => formatBytes(p.storageBytes),
          },
          {
            id: "awake",
            header: "Awake",
            align: "right",
            class: "tabular-nums",
            cell: (p) => formatHours(p.activeSeconds),
          },
          { id: "resumes", header: "Woke up", align: "right", class: "tabular-nums", cell: (p) => p.resumes },
          {
            id: "wake",
            header: "Wake time, typical",
            align: "right",
            class: "tabular-nums",
            cell: (p) => formatMs(p.resumeP50Ms),
          },
          {
            id: "wake-max",
            header: "slowest",
            align: "right",
            class: "tabular-nums text-muted-foreground",
            cell: (p) => formatMs(p.resumeMaxMs),
          },
          {
            id: "active",
            header: "Last awake",
            align: "right",
            class: "text-muted-foreground",
            cell: (p) => (p.lastActiveAt ? timeAgo(p.lastActiveAt) : "—"),
          },
        ]}
      />
    </div>
  );
}

type WorkerKind = "all" | "front" | "api" | "portal";
const KINDS: { value: WorkerKind; label: string }[] = [
  { value: "all", label: "All" },
  { value: "front", label: "Fronts" },
  { value: "api", label: "APIs" },
  { value: "portal", label: "Portal" },
];
const KIND_LABEL: Record<WorkerStats["kind"], string> = {
  front: "Front",
  api: "API",
  portal: "Portal",
  other: "Other",
};

function Workers(props: { workers: WorkerStats[] }) {
  const [kind, setKind] = createSignal<WorkerKind>("all", { name: "workers-kind" });
  // Workers outside the platform (personal projects on the same account) stay off the portal.
  const shown = () => props.workers.filter((w) => w.kind !== "other" && (kind() === "all" || w.kind === kind()));
  const errorShare = (w: WorkerStats) => (w.requests ? w.errors / w.requests : 0);
  return (
    <div class="space-y-4">
      <ToggleGroup ariaLabel="Filter Workers by kind" items={KINDS} value={kind()} onChange={setKind} />
      <Table
        caption="Workers this billing cycle"
        items={shown()}
        columns={[
          {
            id: "worker",
            header: "Worker",
            cell: (w) => (
              <span class="flex flex-col">
                <span class="font-medium">{w.script}</span>
                <span class="text-[12px] text-muted-foreground">
                  {KIND_LABEL[w.kind]}
                  {w.service ? ` of ${findService(w.service)?.name}` : ""}
                </span>
              </span>
            ),
          },
          {
            id: "requests",
            header: "Requests",
            align: "right",
            class: "tabular-nums",
            cell: (w) => formatCount(w.requests),
          },
          {
            id: "errors",
            header: "Errors",
            align: "right",
            class: "tabular-nums",
            cell: (w) => (
              <Show when={w.errors > 0} fallback={<span class="text-muted-foreground">none</span>}>
                <Badge variant={errorShare(w) > 0.01 ? "warning" : "muted"} shape="pill">
                  {formatPercent(errorShare(w), 1)}
                </Badge>
              </Show>
            ),
            title: (w) => `${formatCount(w.errors)} of ${formatCount(w.requests)} requests threw`,
          },
          { id: "cpu", header: "CPU p99", align: "right", class: "tabular-nums", cell: (w) => formatMs(w.cpuP99Ms) },
          { id: "wall", header: "Time p50", align: "right", class: "tabular-nums", cell: (w) => formatMs(w.wallP50Ms) },
          {
            id: "wall99",
            header: "Time p99",
            align: "right",
            class: "tabular-nums text-muted-foreground",
            cell: (w) => formatMs(w.wallP99Ms),
          },
        ]}
      />
      <p class={note}>
        This billing cycle. Errors are uncaught exceptions in the Worker. Time is wall-clock, so an API Worker's
        includes waiting on its container, cold starts too.
      </p>
    </div>
  );
}

const CHART_DAYS = 14;

function ColdStarts(props: { data: MetricsPayload }) {
  const days = () =>
    Array.from({ length: CHART_DAYS }, (_, i) => {
      const now = new Date();
      return new Date(Date.UTC(now.getUTCFullYear(), now.getUTCMonth(), now.getUTCDate() - (CHART_DAYS - 1 - i)));
    });
  const chart = createMemo<ChartDatum[]>(() =>
    SERVICES.flatMap((service) =>
      days().map((date) => ({
        date,
        series: service.name,
        value:
          props.data.coldStarts.find((row) => row.service === service.id && row.day === date.toISOString().slice(0, 10))
            ?.count ?? 0,
      })),
    ),
  );
  const rows = () =>
    props.data.services.map((service) => ({ id: service.id, week: summarize(props.data.coldStarts, service.id, 7) }));
  return (
    <Show
      when={props.data.coldStarts.length > 0}
      fallback={
        <NotConnected source="the API Workers">
          No cold starts recorded yet. They start arriving once each API is released with the cold-start recording.
        </NotConnected>
      }
    >
      <div class="space-y-10">
        <p class={note}>
          A request that finds its API's container asleep waits for it to start. Each API Worker records that wait,
          start plus first answer, as it happens; nothing is started to measure it.
        </p>
        <ChartCard
          title="Starts per day"
          subtitle={`Last ${CHART_DAYS} days. Pick an API to show or hide it.`}
          data={chart()}
          initialType="bar"
          initialRange="all"
          showTypeSelector={false}
          showRangeFilter={false}
          showSearchFilter={false}
          showPointsToggle={false}
          height={220}
          dateFormatter={(date) => date.toLocaleDateString("en", { month: "short", day: "numeric", timeZone: "UTC" })}
        />
        <PageSection title="Last 7 days">
          <Table
            caption="Cold starts by API, last 7 days"
            items={rows()}
            columns={[
              { id: "service", header: "API", class: "font-medium", cell: (row) => findService(row.id)?.name },
              { id: "starts", header: "Starts", align: "right", class: "tabular-nums", cell: (row) => row.week.count },
              {
                id: "failed",
                header: "Failed",
                align: "right",
                class: "tabular-nums",
                cell: (row) => (
                  <Show when={row.week.failed > 0} fallback={<span class="text-muted-foreground">none</span>}>
                    <Badge
                      shape="pill"
                      variant={
                        row.week.failed / row.week.count > 0.5
                          ? "destructive"
                          : row.week.failed / row.week.count > 0.1
                            ? "warning"
                            : "muted"
                      }
                    >
                      {row.week.failed} of {row.week.count}
                    </Badge>
                  </Show>
                ),
              },
              {
                id: "typical",
                header: "Typical wait",
                align: "right",
                class: "tabular-nums",
                cell: (row) => formatMs(row.week.typicalMs),
              },
              {
                id: "worst",
                header: "Worst day, p95",
                align: "right",
                class: "tabular-nums text-muted-foreground",
                cell: (row) => formatMs(row.week.worstP95Ms),
              },
            ]}
          />
        </PageSection>
      </div>
    </Show>
  );
}
