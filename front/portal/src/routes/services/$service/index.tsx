import { Link, createFileRoute, notFound } from "@tanstack/solid-router";
import { For, Show } from "solid-js";
import { REPO_URL, findService } from "@/features/catalog/services";
import { loadDelivery } from "@/features/delivery/api";
import { loadMetrics } from "@/features/metrics/api";
import { summarize } from "@/features/metrics/cold-starts";
import { formatCount, formatMs, formatPercent } from "@/features/metrics/format";
import { tagUrl } from "@/features/delivery/model";
import { loadIssues } from "@/features/issues/api";
import { isOpen } from "@/features/issues/model";
import { IssueList } from "@/features/issues/ui/IssueList";
import { EndpointTable, ProblemTable } from "@/features/metrics/ui/LatencyTables";
import { Kpi, KpiRow } from "@/widgets/shell/Kpi";
import { PageHeader, Section } from "@/widgets/shell/PageHeader";
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
    // What's running: the newest successful deploy of the service.
    const current = delivery?.deploys.find((deploy) => deploy.component === service.id && deploy.status === "success");
    return {
      service,
      metrics: metrics?.services.find((s) => s.id === service.id) ?? null,
      coldStarts: metrics ? summarize(metrics.coldStarts, service.id, 7) : null,
      endpoints: metrics?.endpoints.filter((e) => e.service === service.id).slice(0, 5) ?? [],
      problems: metrics?.problems.filter((p) => p.service === service.id).slice(0, 5) ?? [],
      tag: current?.tag ?? null,
      issues: issues.filter((issue) => issue.service === service.id && isOpen(issue)),
    };
  },
  head: ({ loaderData }) => ({ meta: [{ title: `${loaderData?.service.name ?? "Service"} · TrieOH dev` }] }),
  component: ServicePage,
});

const button =
  "inline-flex h-8 items-center rounded-lg border border-border px-3 text-[13px] transition-colors hover:bg-foreground/5";

function ServicePage() {
  const data = Route.useLoaderData();
  const service = () => data().service;

  const links = () => [
    { label: "Code", href: `${REPO_URL}/tree/main/${service().dir}` },
    { label: "Front", href: service().front },
  ];

  return (
    <>
      <PageHeader title={service().name} description={service().description}>
        <p class="flex flex-wrap gap-x-3 font-mono text-[12px] text-muted-foreground">
          <span>{service().apiUrl.replace("https://", "")}</span>
          <Show when={data().tag}>
            {(tag) => (
              <a class="hover:underline" href={tagUrl(tag())} target="_blank" rel="noopener">
                {tag()}
              </a>
            )}
          </Show>
        </p>
        <div class="flex flex-wrap gap-2 pt-3">
          <Link to="/services/$service/api" params={{ service: service().id }} class={button}>
            API reference
          </Link>
          <Link to="/problems/$set" params={{ set: service().id }} class={button}>
            Problems
          </Link>
          <For each={links()}>
            {(link) => (
              <a class={button} href={link.href} target="_blank" rel="noopener">
                {link.label}
              </a>
            )}
          </For>
        </div>
      </PageHeader>

      <Section title="Health">
        <KpiRow>
          <Kpi label="Status">
            <span class="inline-flex items-center gap-2">
              <StatusDot status={data().metrics?.status ?? "unknown"} />
              {statusLabel(data().metrics?.status ?? "unknown")}
            </span>
          </Kpi>
          <Kpi label="Success rate 30d">{formatPercent(data().metrics?.successRate30d ?? null)}</Kpi>
          <Kpi label="Requests 24h">
            {data().metrics?.requests24h == null ? "—" : formatCount(data().metrics!.requests24h!)}
          </Kpi>
          <Kpi label="p95 24h" note="inside the container">
            {formatMs(data().metrics?.p95Ms24h ?? null)}
          </Kpi>
        </KpiRow>
        <Show when={data().coldStarts?.count}>
          <p class="mt-3 text-[13px] text-muted-foreground">
            {data().coldStarts!.count} cold starts in the last 7 days, typically {formatMs(data().coldStarts!.typicalMs)}
            {data().coldStarts!.failed ? `; ${data().coldStarts!.failed} failed` : ""}.{" "}
            <Link to="/metrics" class="underline underline-offset-4">
              Metrics
            </Link>
          </p>
        </Show>
      </Section>

      <Show when={data().endpoints.length > 0 || data().problems.length > 0}>
        <Section title="Latency and errors" aside={<Link to="/metrics">Metrics</Link>}>
          <p class="mb-3 text-[13px] text-muted-foreground">Slowest endpoints and error types, last 7 days.</p>
          <EndpointTable endpoints={data().endpoints} />
          <div class="mt-6">
            <ProblemTable problems={data().problems} />
          </div>
        </Section>
      </Show>

      <Section title={`Open issues (${data().issues.length})`} aside={<Link to="/issues">All issues</Link>}>
        <IssueList issues={data().issues} empty={`No open issues for ${service().name}.`} />
      </Section>
    </>
  );
}
