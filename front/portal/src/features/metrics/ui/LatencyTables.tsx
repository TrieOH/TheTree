import { Link } from "@tanstack/solid-router";
import { Meter, Table, type TableColumn } from "@trieoh/ui-solid";
import { Show } from "solid-js";
import { findService, type ServiceId } from "@/features/catalog/services";
import { findProblem } from "@/features/problems/model";
import { StatusPill } from "@/features/problems/ui/ProblemList";
import { formatCount, formatMs } from "../format";
import { DEGRADED_P95_MS, type EndpointStats, type ProblemStats } from "../model";

// Slowest endpoints and top error types, on Metrics and on each service page.
// Route templates and problem slugs only: never paths, messages or tenants.

const empty = (text: string) => <p class="text-[13px] text-muted-foreground">{text}</p>;

function serviceColumn<T extends { service: ServiceId }>(): TableColumn<T> {
  return {
    id: "service",
    header: "Service",
    class: "text-muted-foreground",
    cell: (row) => findService(row.service)?.name,
  };
}

/** p95 as a bar against the degraded line, so a slow route stands out before you read the number. */
function P95(props: { ms: number | null; scale: number }) {
  return (
    <div class="flex items-center gap-3">
      <Meter
        class="w-28 shrink-0 [&>div]:h-1.5"
        ariaLabel="p95"
        value={props.ms ?? 0}
        max={props.scale}
        tone={(props.ms ?? 0) > DEGRADED_P95_MS ? "warning" : "accent"}
        markers={[{ value: DEGRADED_P95_MS, label: `Degraded above ${formatMs(DEGRADED_P95_MS)}` }]}
      />
      <span class={(props.ms ?? 0) > DEGRADED_P95_MS ? "text-warn" : undefined}>{formatMs(props.ms)}</span>
    </div>
  );
}

export function EndpointTable(props: { endpoints: EndpointStats[]; showService?: boolean }) {
  // The degraded line sits at the same spot on every row; slower routes stretch the scale.
  const scale = () => Math.max(DEGRADED_P95_MS * 1.25, ...props.endpoints.map((e) => e.p95Ms ?? 0));
  const columns = (): TableColumn<EndpointStats>[] => [
    {
      id: "endpoint",
      header: "Endpoint",
      cell: (endpoint) => (
        <span class="font-mono text-[12.5px]">
          <span class="text-muted-foreground">{endpoint.method}</span> {endpoint.route}
        </span>
      ),
    },
    ...(props.showService ? [serviceColumn<EndpointStats>()] : []),
    { id: "p95", header: "p95, worst day", class: "tabular-nums", cell: (e) => <P95 ms={e.p95Ms} scale={scale()} /> },
    {
      id: "requests",
      header: "Requests",
      align: "right",
      class: "tabular-nums",
      cell: (endpoint) => formatCount(endpoint.requests),
    },
  ];
  return (
    <Table
      caption="Slowest endpoints"
      items={props.endpoints}
      columns={columns()}
      emptyState={empty("No requests recorded.")}
    />
  );
}

/** The problem's doc, under its service or the shared set. */
function problemSet(problem: ProblemStats): string | null {
  if (findProblem(problem.service, problem.problem)) return problem.service;
  if (findProblem("shared", problem.problem)) return "shared";
  return null;
}

export function ProblemTable(props: { problems: ProblemStats[]; showService?: boolean }) {
  const columns = (): TableColumn<ProblemStats>[] => [
    {
      id: "problem",
      header: "Problem",
      class: "font-mono text-[12.5px]",
      cell: (problem) => (
        <Show when={problemSet(problem)} fallback={problem.problem}>
          {(set) => (
            <Link
              to="/problems/$set/$slug"
              params={{ set: set(), slug: problem.problem }}
              class="underline decoration-border underline-offset-4 hover:decoration-foreground"
            >
              {problem.problem}
            </Link>
          )}
        </Show>
      ),
    },
    ...(props.showService ? [serviceColumn<ProblemStats>()] : []),
    {
      id: "status",
      header: "Status",
      cell: (problem) => (problem.status ? <StatusPill status={problem.status} /> : "—"),
    },
    { id: "answers", header: "Answers", align: "right", class: "tabular-nums", cell: (p) => formatCount(p.count) },
  ];
  return (
    <Table
      caption="Top problem answers"
      items={props.problems}
      columns={columns()}
      emptyState={empty("No problem answers recorded.")}
    />
  );
}
