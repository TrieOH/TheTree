import { Table, type TableColumn } from "@trieoh/ui-solid";
import { findService, type ServiceId } from "@/features/catalog/services";
import { formatCount, formatMs } from "../format";
import { DEGRADED_P95_MS, type EndpointStats, type ProblemStats } from "../model";

// Slowest endpoints and top error types, on Metrics and on each service page.
// Route templates and problem slugs only: never paths, messages or tenants.

const empty = (text: string) => <p class="text-[13px] text-muted-foreground">{text}</p>;

function serviceColumn<T extends { service: ServiceId }>(): TableColumn<T> {
  return { id: "service", header: "Service", class: "font-medium", cell: (row) => findService(row.service)?.name };
}

export function EndpointTable(props: { endpoints: EndpointStats[]; showService?: boolean }) {
  const columns = (): TableColumn<EndpointStats>[] => [
    ...(props.showService ? [serviceColumn<EndpointStats>()] : []),
    {
      id: "endpoint",
      header: "Endpoint",
      class: "font-mono text-[12.5px]",
      cell: (endpoint) => (
        <>
          <span class="text-muted-foreground">{endpoint.method}</span> {endpoint.route}
        </>
      ),
    },
    {
      id: "p95",
      header: "p95",
      class: (endpoint) => ((endpoint.p95Ms ?? 0) > DEGRADED_P95_MS ? "tabular-nums text-warn" : "tabular-nums"),
      cell: (endpoint) => formatMs(endpoint.p95Ms),
    },
    { id: "requests", header: "Requests", class: "tabular-nums", cell: (endpoint) => formatCount(endpoint.requests) },
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

export function ProblemTable(props: { problems: ProblemStats[]; showService?: boolean }) {
  const columns = (): TableColumn<ProblemStats>[] => [
    ...(props.showService ? [serviceColumn<ProblemStats>()] : []),
    { id: "problem", header: "Problem", class: "font-mono text-[12.5px]", cell: (problem) => problem.problem },
    {
      id: "status",
      header: "Status",
      class: (problem) => (problem.status >= 500 ? "tabular-nums text-warn" : "tabular-nums"),
      cell: (problem) => problem.status || "—",
    },
    { id: "answers", header: "Answers", class: "tabular-nums", cell: (problem) => formatCount(problem.count) },
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
