import { Link, createFileRoute } from "@tanstack/solid-router";
import { Tabs, cn } from "@trieoh/ui-solid";
import { For, Match, Show, Switch, createSignal } from "solid-js";
import { findService } from "@/features/catalog/services";
import { loadMetrics } from "@/features/metrics/api";
import { formatBytes, formatCount, formatHours, formatMs, formatPercent } from "@/features/metrics/format";
import { DEGRADED_P95_MS, type MetricsPayload, type WorkerStats } from "@/features/metrics/model";
import { summarize } from "@/features/metrics/cold-starts";
import { Sparkline } from "@/features/metrics/ui/Sparkline";
import { EndpointTable, ProblemTable } from "@/features/metrics/ui/LatencyTables";
import { UsageBar } from "@/features/metrics/ui/UsageBar";
import { timeAgo } from "@/features/delivery/stats";
import { formatDay } from "@/features/issues/stats";
import { NotConnected, PageHeader } from "@/widgets/shell/PageHeader";
import { StatusDot, statusLabel } from "@/widgets/shell/StatusDot";

export const Route = createFileRoute("/metrics")({
  loader: () => loadMetrics(),
  head: () => ({ meta: [{ title: "Metrics · TrieOH dev" }] }),
  component: MetricsPage,
});

const TABS = [
  { id: "services", label: "Services" },
  { id: "latency", label: "Latency & errors" },
  { id: "usage", label: "Plan usage" },
  { id: "database", label: "Database" },
  { id: "fronts", label: "Frontends" },
  { id: "cold-starts", label: "Cold starts" },
] as const;

type Tab = (typeof TABS)[number]["id"];

const th = "py-2 pr-4 font-normal";
const td = "py-2.5 pr-4 tabular-nums";

function MetricsPage() {
  const data = Route.useLoaderData();
  const [tab, setTab] = createSignal<Tab>("services", { name: "metrics-tab" });

  return (
    <>
      <PageHeader
        title="Metrics"
        description={
          <>
            From Cloudflare's analytics, Neon's API and Grafana Cloud, pulled every 30 minutes; nothing here reaches a
            service. Full dashboards and traces live in Grafana.
          </>
        }
      />
      <Tabs ariaLabel="Metric views" items={TABS} value={tab()} onChange={setTab} />
      <div class="mt-6">
        <Show when={data().syncedAt} fallback={<NotConnected source="Cloudflare and Neon" />}>
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
              <Database data={data()} />
            </Match>
            <Match when={tab() === "fronts"}>
              <Workers workers={data().workers.filter((w) => w.kind === "front" || w.kind === "portal")} />
            </Match>
            <Match when={tab() === "cold-starts"}>
              <ColdStarts data={data()} />
              <h3 class="mb-2 mt-10 font-heading text-[15px] font-semibold">API Workers</h3>
              <Workers workers={data().workers.filter((w) => w.kind === "api")} />
            </Match>
          </Switch>
          <p class="mt-8 text-[12px] text-muted-foreground">Last pulled {timeAgo(data().syncedAt!)}.</p>
        </Show>
      </div>
    </>
  );
}

function Services(props: { data: MetricsPayload }) {
  return (
    <div class="grid gap-3 sm:grid-cols-2">
      <For each={props.data.services}>
        {(service) => (
          <Link
            to="/services/$service"
            params={{ service: service.id }}
            class="rounded-xl border border-border bg-card p-4 transition-colors hover:bg-foreground/5"
          >
            <div class="flex items-center justify-between gap-3">
              <span class="font-medium">{findService(service.id)?.name}</span>
              <span class="inline-flex items-center gap-1.5 text-[12px] text-muted-foreground">
                <StatusDot status={service.status} />
                {statusLabel(service.status)}
              </span>
            </div>
            <div class="mt-3">
              <Sparkline values={service.hourly} label={`Requests per hour, last 24 hours`} />
            </div>
            <dl class="mt-3 grid grid-cols-4 gap-2 text-[12px]">
              <div>
                <dt class="text-muted-foreground">Requests 24h</dt>
                <dd class="tabular-nums">{service.requests24h === null ? "—" : formatCount(service.requests24h)}</dd>
              </div>
              <div>
                <dt class="text-muted-foreground">5xx 24h</dt>
                <dd class={cn("tabular-nums", (service.errors24h ?? 0) > 0 && "text-warn")}>{service.errors24h ?? "—"}</dd>
              </div>
              <div>
                <dt class="text-muted-foreground">p95 24h</dt>
                <dd class={cn("tabular-nums", (service.p95Ms24h ?? 0) > DEGRADED_P95_MS && "text-warn")}>
                  {formatMs(service.p95Ms24h)}
                </dd>
              </div>
              <div>
                <dt class="text-muted-foreground">Success 30d</dt>
                <dd class="tabular-nums">{formatPercent(service.successRate30d)}</dd>
              </div>
            </dl>
            <Show when={service.lastRequestAt}>
              {(at) => <p class="mt-2 text-[11.5px] text-muted-foreground">Last request {timeAgo(at())}</p>}
            </Show>
          </Link>
        )}
      </For>
    </div>
  );
}

function LatencyAndErrors(props: { data: MetricsPayload }) {
  return (
    <div class="space-y-4">
      <p class="text-[13px] text-muted-foreground">
        Measured inside each API's container and read from Grafana Cloud, so a cold start's wait is not included (see
        Cold starts). Last 7 days; an endpoint's p95 is its worst day.
      </p>
      <Show
        when={props.data.endpoints.length > 0 || props.data.problems.length > 0}
        fallback={<NotConnected source="Grafana Cloud">No latency or error data pulled yet.</NotConnected>}
      >
        <h3 class="font-heading text-[15px] font-semibold">Slowest endpoints</h3>
        <EndpointTable endpoints={props.data.endpoints.slice(0, 15)} showService />
        <h3 class="mt-10 font-heading text-[15px] font-semibold">Top error types</h3>
        <ProblemTable problems={props.data.problems.slice(0, 15)} showService />
      </Show>
    </div>
  );
}

function ColdStarts(props: { data: MetricsPayload }) {
  return (
    <div class="space-y-4">
      <p class="text-[13px] text-muted-foreground">
        A request that finds its API's container asleep waits for it to start. Each API Worker records that wait (start
        plus first answer) as it happens; nothing is started to measure it. Last 7 days; bars are starts per day over 14.
      </p>
      <Show
        when={props.data.coldStarts.length > 0}
        fallback={
          <NotConnected source="the API Workers">
            No cold starts recorded yet. They start arriving once each API is released with the cold-start recording.
          </NotConnected>
        }
      >
        <div class="overflow-x-auto">
          <table class="w-full text-left text-[13px]">
            <thead class="text-[12px] text-muted-foreground">
              <tr class="border-b border-border">
                <th class={th}>Service</th>
                <th class={th}>Starts</th>
                <th class={th}>Failed</th>
                <th class={th}>Typical</th>
                <th class={th}>Worst day p95</th>
                <th class="py-2 font-normal">Per day</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-border">
              <For each={props.data.services}>
                {(service) => {
                  const week = summarize(props.data.coldStarts, service.id, 7);
                  const fortnight = summarize(props.data.coldStarts, service.id, 14);
                  return (
                    <tr>
                      <td class="py-2.5 pr-4 font-medium">{findService(service.id)?.name}</td>
                      <td class={td}>{week.count}</td>
                      <td class={cn(td, week.failed > 0 && "text-warn")}>{week.failed}</td>
                      <td class={td}>{formatMs(week.typicalMs)}</td>
                      <td class={td}>{formatMs(week.worstP95Ms)}</td>
                      <td class="py-2.5">
                        <Sparkline values={fortnight.daily} label="Cold starts per day, last 14 days" />
                      </td>
                    </tr>
                  );
                }}
              </For>
            </tbody>
          </table>
        </div>
      </Show>
    </div>
  );
}

function PlanUsage(props: { data: MetricsPayload }) {
  const pace = () => {
    const start = Date.parse(props.data.cycle.start);
    const end = Date.parse(props.data.cycle.end);
    return Math.min(1, Math.max(0, (Date.now() - start) / (end - start)));
  };
  return (
    <div class="space-y-6">
      <p class="text-[13px] text-muted-foreground">
        Billing cycle {formatDay(props.data.cycle.start)} – {formatDay(props.data.cycle.end)}. The tick marks an even
        pace through it; a bar past its tick is on course to run out.
      </p>
      <For each={props.data.usage}>{(usage) => <UsageBar usage={usage} pace={pace()} />}</For>
    </div>
  );
}

function Database(props: { data: MetricsPayload }) {
  return (
    <div class="space-y-4">
      <p class="text-[13px] text-muted-foreground">
        One Neon project per service, read from Neon's API; nothing connects to Postgres. Neon's month runs to its own
        quota reset.
      </p>
      <div class="overflow-x-auto">
        <table class="w-full text-left text-[13px]">
          <thead class="text-[12px] text-muted-foreground">
            <tr class="border-b border-border">
              <th class={th}>Service</th>
              <th class={th}>Storage</th>
              <th class={th}>Awake</th>
              <th class={th}>Suspended</th>
              <th class={th}>Resumes</th>
              <th class={th}>Resume p50 / max</th>
              <th class="py-2 font-normal">Last active</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <For each={props.data.neon}>
              {(project) => (
                <tr>
                  <td class="py-2.5 pr-4 font-medium">{findService(project.service)?.name}</td>
                  <td class={td}>{formatBytes(project.storageBytes)}</td>
                  <td class={td}>{formatHours(project.activeSeconds)}</td>
                  <td class={td}>{formatPercent(project.suspendedShare, 1)}</td>
                  <td class={td}>{project.resumes}</td>
                  <td class={td}>
                    {formatMs(project.resumeP50Ms)} / {formatMs(project.resumeMaxMs)}
                  </td>
                  <td class="py-2.5 text-muted-foreground">{project.lastActiveAt ? timeAgo(project.lastActiveAt) : "—"}</td>
                </tr>
              )}
            </For>
          </tbody>
        </table>
      </div>
    </div>
  );
}

function Workers(props: { workers: WorkerStats[] }) {
  return (
    <div class="overflow-x-auto">
      <table class="w-full text-left text-[13px]">
        <thead class="text-[12px] text-muted-foreground">
          <tr class="border-b border-border">
            <th class={th}>Worker</th>
            <th class={th}>Requests</th>
            <th class={th}>Errors</th>
            <th class={th}>CPU p50 / p99</th>
            <th class="py-2 font-normal">Wall p50 / p99</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border">
          <For each={props.workers}>
            {(worker) => (
              <tr>
                <td class="py-2.5 pr-4 font-mono text-[12.5px]">{worker.script}</td>
                <td class={td}>{formatCount(worker.requests)}</td>
                <td class={cn(td, worker.errors > 0 && "text-warn")}>{worker.errors}</td>
                <td class={td}>
                  {formatMs(worker.cpuP50Ms)} / {formatMs(worker.cpuP99Ms)}
                </td>
                <td class="py-2.5 tabular-nums">
                  {formatMs(worker.wallP50Ms)} / {formatMs(worker.wallP99Ms)}
                </td>
              </tr>
            )}
          </For>
        </tbody>
      </table>
      <p class="mt-3 text-[12px] text-muted-foreground">This billing cycle. Errors are uncaught exceptions in the Worker.</p>
    </div>
  );
}
