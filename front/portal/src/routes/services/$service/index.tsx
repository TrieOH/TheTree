import type { JSX } from "@solidjs/web";
import { Link, createFileRoute, notFound, useNavigate } from "@tanstack/solid-router";
import { DashboardStatCard, EmptyState, PageHeader, PageSection, StatusDot as Dot, buttonVariants, cn } from "@trieoh/ui-solid";
import { For, Show, createMemo } from "solid-js";
import { REPO_URL, findService } from "@/features/catalog/services";
import { coldStartDays, serviceAttention, type OverviewSource } from "@/features/catalog/overview";
import { ColdStartStrip } from "@/features/catalog/ui/ColdStartStrip";
import { loadDelivery } from "@/features/delivery/api";
import { runUrl, tagUrl, type Deploy } from "@/features/delivery/model";
import { timeAgo } from "@/features/delivery/stats";
import { loadIssues } from "@/features/issues/api";
import { isOpen } from "@/features/issues/model";
import { IssueList } from "@/features/issues/ui/IssueList";
import { loadMetrics } from "@/features/metrics/api";
import { summarize } from "@/features/metrics/cold-starts";
import { formatBytes, formatCount, formatMs, formatPercent } from "@/features/metrics/format";
import { DEGRADED_P95_MS } from "@/features/metrics/model";
import { EndpointTable, ProblemTable } from "@/features/metrics/ui/LatencyTables";
import { Attention } from "@/shared/ui/attention";
import { KpiRow } from "@/widgets/shell/Kpi";
import { StatusDot, statusLabel } from "@/widgets/shell/StatusDot";

export const Route = createFileRoute("/services/$service/")({
  loader: async ({ params }) => {
    const service = findService(params.service);
    if (!service) throw notFound();
    const [{ issues }, delivery, metrics] = await Promise.all([
      loadIssues(),
      loadDelivery().catch(() => null),
      loadMetrics().catch(() => null),
    ]);
    const deploys = (component: string) => delivery?.deploys.filter((d) => d.component === component) ?? [];
    return {
      service,
      attention: serviceAttention(service, metrics, delivery),
      metrics: metrics?.services.find((s) => s.id === service.id) ?? null,
      syncedAt: metrics?.syncedAt ?? null,
      coldStarts: metrics ? summarize(metrics.coldStarts, service.id, 7) : null,
      coldStartDays: metrics ? coldStartDays(metrics.coldStarts, service.id) : [],
      endpoints: metrics?.endpoints.filter((e) => e.service === service.id).slice(0, 5) ?? [],
      problems: metrics?.problems.filter((p) => p.service === service.id).slice(0, 5) ?? [],
      workers: metrics?.workers.filter((w) => w.service === service.id) ?? [],
      neon: metrics?.neon.find((n) => n.service === service.id) ?? null,
      api: deploys(service.id),
      front: deploys(`${service.id}-ui`),
      issues: issues.filter((issue) => issue.service === service.id && isOpen(issue)),
    };
  },
  head: ({ loaderData }) => ({ meta: [{ title: `${loaderData?.service.name ?? "Service"} · TrieOH dev` }] }),
  component: ServicePage,
});

const SOURCE_LABEL: Record<OverviewSource, string> = { metrics: "Metrics", delivery: "Delivery" };
const external = { target: "_blank", rel: "noopener" } as const;
const action = buttonVariants({ variant: "outline", size: "sm" });
const note = "text-[13px] text-muted-foreground";

function ServicePage() {
  const data = Route.useLoaderData();
  const navigate = useNavigate();
  const service = () => data().service;
  const m = () => data().metrics;
  const peak = () => Math.max(0, ...(m()?.hourly ?? []));

  return (
    <>
      <PageHeader title={service().name} description={service().description}>
        <div class="flex flex-wrap gap-2 pt-3">
          <Link to="/services/$service/api" params={{ service: service().id }} class={action}>
            API reference
          </Link>
          <Link to="/problems/$set" params={{ set: service().id }} class={action}>
            Problems
          </Link>
          <a class={action} href={`${REPO_URL}/tree/main/${service().dir}`} {...external}>
            Code
          </a>
          <a class={action} href={service().front} {...external}>
            Open the front
          </a>
        </div>
      </PageHeader>

      <Attention
        items={data().attention}
        tabLabel={(source) => SOURCE_LABEL[source]}
        onOpen={(source) => navigate({ to: source === "metrics" ? "/metrics" : "/delivery" })}
      />

      <KpiRow>
        <DashboardStatCard
          label="Status"
          hint={m()?.lastRequestAt ? `called ${timeAgo(m()!.lastRequestAt!)}` : "no calls yet"}
          value={
            <span class="inline-flex items-center gap-2">
              <StatusDot status={m()?.status ?? "unknown"} />
              {statusLabel(m()?.status ?? "unknown")}
            </span>
          }
        />
        <DashboardStatCard
          label="Requests 24h"
          hint={peak() > 0 ? `busiest hour ${formatCount(peak())}` : "idle"}
          value={m()?.requests24h == null ? "—" : formatCount(m()!.requests24h!)}
          sparkline={
            m()?.hourly.length
              ? {
                  type: "bar",
                  fill: true,
                  data: m()!.hourly,
                  height: 22,
                  color: "var(--accent)",
                  label: "Requests per hour, last 24 hours",
                }
              : undefined
          }
        />
        <DashboardStatCard
          label="p95 24h"
          hint={`${m()?.errors24h ?? 0} server ${(m()?.errors24h ?? 0) === 1 ? "error" : "errors"}`}
          value={
            <span class={(m()?.p95Ms24h ?? 0) > DEGRADED_P95_MS ? "text-warn" : undefined}>
              {formatMs(m()?.p95Ms24h ?? null)}
            </span>
          }
        />
        <DashboardStatCard label="Success" hint="last 30 days" value={formatPercent(m()?.successRate30d ?? null)} />
      </KpiRow>

      <div class="mt-10 grid gap-x-12 gap-y-2 lg:grid-cols-[minmax(0,1fr)_300px]">
        <div class="min-w-0">
          <PageSection title="Cold starts" aside="last 30 days">
            <Show
              when={data().coldStartDays.some((d) => d.count > 0)}
              fallback={<EmptyState variant="inline" title="No cold starts recorded in the last 30 days." />}
            >
              <ColdStartStrip bars={data().coldStartDays} />
              <Show when={data().coldStarts}>
                {(week) => (
                  <p class={cn(note, "mt-3")}>
                    Last 7 days: {week().count} {week().count === 1 ? "start" : "starts"}, usually{" "}
                    {formatMs(week().typicalMs)} to the first answer
                    <Show when={week().failed > 0}>
                      , <span class="font-medium text-foreground">{week().failed} failed</span>
                    </Show>
                    . Failed starts are the red part of each bar.
                  </p>
                )}
              </Show>
            </Show>
          </PageSection>

          <PageSection title="Slowest endpoints" aside="last 7 days">
            <Show
              when={data().endpoints.length > 0}
              fallback={<EmptyState variant="inline" title="No endpoint timings in the last 7 days." />}
            >
              <EndpointTable endpoints={data().endpoints} />
            </Show>
          </PageSection>

          <PageSection title="Error answers" aside="last 7 days">
            <Show
              when={data().problems.length > 0}
              fallback={<EmptyState variant="inline" title="No problem answers in the last 7 days." />}
            >
              <ProblemTable problems={data().problems} />
            </Show>
          </PageSection>

          <PageSection title={`Open issues (${data().issues.length})`} aside={<Link to="/issues">All issues</Link>}>
            <IssueList issues={data().issues} empty={`No open issues for ${service().name}.`} />
          </PageSection>
        </div>

        <aside class="min-w-0 space-y-10 lg:border-l lg:border-border lg:pl-8">
          <RailSection title="Running now">
            <Running label="API" host={service().apiUrl} deploys={data().api} />
            <Running label="Front" host={service().front} deploys={data().front} />
          </RailSection>

          <Show when={data().api.length > 0}>
            <RailSection title="Recent API releases">
              <ul class="space-y-1.5 text-[13px]">
                <For each={data().api.slice(0, 5)}>
                  {(deploy) => (
                    <li class="flex items-center gap-2">
                      <Dot
                        tone={deploy.status === "success" ? "success" : deploy.status === "failure" ? "destructive" : "muted"}
                        label={deploy.status === "success" ? "Shipped" : deploy.status === "failure" ? "Failed" : "Cancelled"}
                      />
                      <a class="font-mono text-[12.5px] hover:underline" href={tagUrl(deploy.tag)} {...external}>
                        {deploy.tag.slice(deploy.component.length + 1)}
                      </a>
                      <span class="ml-auto text-[12px] text-muted-foreground">{timeAgo(deploy.at)}</span>
                    </li>
                  )}
                </For>
              </ul>
              <p class={cn(note, "mt-2 text-[12px]")}>
                {data().api.filter((d) => d.status === "success" && Date.now() - Date.parse(d.at) < 30 * 86_400_000).length}{" "}
                releases in the last 30 days. <Link to="/delivery" class="underline underline-offset-4">Delivery</Link>
              </p>
            </RailSection>
          </Show>

          <RailSection title="Database">
            <Show when={data().neon} fallback={<p class={note}>Not pulled from Neon yet.</p>}>
              {(db) => (
                <Facts
                  rows={[
                    ["Size", formatBytes(db().storageBytes)],
                    ["Asleep", db().suspendedShare === null ? "—" : `${formatPercent(db().suspendedShare, 0)} of the time`],
                    ["Wake-ups", `${formatCount(db().resumes)}, usually ${formatMs(db().resumeP50Ms)}`],
                    ["Slowest wake", formatMs(db().resumeMaxMs)],
                    ["Last active", db().lastActiveAt ? timeAgo(db().lastActiveAt!) : "—"],
                  ]}
                />
              )}
            </Show>
          </RailSection>

          <RailSection title="Workers" aside="this billing cycle">
            <Show when={data().workers.length > 0} fallback={<p class={note}>Not pulled from Cloudflare yet.</p>}>
              <div class="space-y-4">
                <For each={[...data().workers].sort((a) => (a.kind === "api" ? -1 : 1))}>
                  {(w) => (
                    <div>
                      <p class="mb-1 text-[13px] font-medium">
                        {w.kind === "api" ? "API Worker" : "Front Worker"}{" "}
                        <span class="font-mono text-[11.5px] font-normal text-muted-foreground">{w.script}</span>
                      </p>
                      <Facts
                        rows={[
                          ["Requests", formatCount(w.requests)],
                          [
                            "Errors",
                            w.errors === 0 ? "none" : `${formatCount(w.errors)} (${formatPercent(w.errors / w.requests, 1)})`,
                            w.requests > 0 && w.errors / w.requests > 0.01,
                          ],
                          ["Wall time p99", formatMs(w.wallP99Ms)],
                          ["CPU p99", formatMs(w.cpuP99Ms)],
                        ]}
                      />
                    </div>
                  )}
                </For>
              </div>
            </Show>
          </RailSection>

          <Show when={data().syncedAt}>
            {(at) => <p class="text-[12px] text-muted-foreground">Metrics pulled {timeAgo(at())}.</p>}
          </Show>
        </aside>
      </div>
    </>
  );
}

function RailSection(props: { title: string; aside?: string; children: JSX.Element }) {
  return (
    <section>
      <h2 class="mb-3 flex items-baseline justify-between gap-2 font-heading text-[14px] font-semibold">
        {props.title}
        <Show when={props.aside}>
          <span class="text-[12px] font-normal text-muted-foreground">{props.aside}</span>
        </Show>
      </h2>
      {props.children}
    </section>
  );
}

/** Label and value pairs, one per line. A third item flags the value. */
function Facts(props: { rows: readonly (readonly [string, string, boolean?])[] }) {
  return (
    <dl class="grid grid-cols-[auto_minmax(0,1fr)] gap-x-4 gap-y-1.5 text-[13px]">
      <For each={props.rows}>
        {([label, value, warn]) => (
          <>
            <dt class="text-muted-foreground">{label}</dt>
            <dd class={cn("text-right tabular-nums", warn && "font-medium text-warn")}>{value}</dd>
          </>
        )}
      </For>
    </dl>
  );
}

/** What runs at one host: the newest successful deploy, with its version and age. */
function Running(props: { label: string; host: string; deploys: readonly Deploy[] }) {
  const current = createMemo(() => props.deploys.find((d) => d.status === "success"));
  return (
    <div class="mb-4 last:mb-0">
      <p class="flex items-baseline justify-between gap-2 text-[13px]">
        <span class="font-medium">{props.label}</span>
        <Show when={current()} fallback={<span class="text-muted-foreground">no deploy recorded</span>}>
          {(d) => (
            <a class="font-mono text-[12.5px] hover:underline" href={tagUrl(d().tag)} {...external}>
              {d().tag.slice(d().component.length + 1)}
            </a>
          )}
        </Show>
      </p>
      <p class="mt-0.5 flex items-baseline justify-between gap-2 text-[12px] text-muted-foreground">
        <a class="truncate hover:underline" href={props.host} {...external}>
          {props.host.replace("https://", "")}
        </a>
        <Show when={current()}>
          {(d) => (
            <Show when={d().runId} fallback={<span class="shrink-0">{timeAgo(d().at)}</span>}>
              {(id) => (
                <a class="shrink-0 hover:underline" href={runUrl(id())} {...external}>
                  {timeAgo(d().at)}
                </a>
              )}
            </Show>
          )}
        </Show>
      </p>
    </div>
  );
}
