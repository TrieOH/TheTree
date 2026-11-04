import { createFileRoute } from "@tanstack/solid-router";
import { Tabs } from "@trieoh/ui-solid";
import { NotConnected, PageHeader } from "@/widgets/shell/PageHeader";

export const Route = createFileRoute("/metrics")({
  head: () => ({ meta: [{ title: "Metrics · TrieOH dev" }] }),
  component: MetricsPage,
});

const TABS = [
  { id: "services", label: "Services" },
  { id: "cold-starts", label: "Cold starts" },
  { id: "usage", label: "Plan usage" },
  { id: "database", label: "Database" },
  { id: "fronts", label: "Frontends" },
] as const;

function MetricsPage() {
  return (
    <>
      <PageHeader
        title="Metrics"
        description="Last 30 days. Full dashboards and traces live in Grafana."
      />
      <Tabs
        ariaLabel="Metric views"
        items={TABS.map((tab) => ({
          ...tab,
          content: <NotConnected source="Cloudflare, Grafana and Neon" />,
        }))}
        panelClass="pt-6"
      />
    </>
  );
}
