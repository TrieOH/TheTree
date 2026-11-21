import { createFileRoute } from "@tanstack/solid-router";
import { DashboardStatCard, PageHeader, Tabs } from "@trieoh/ui-solid";
import { Match, Switch, createMemo, createSignal } from "solid-js";
import { loadDelivery } from "@/features/delivery/api";
import { needsAttention, type DeliveryTab } from "@/features/delivery/attention";
import { groupReleases } from "@/features/delivery/releases";
import { changeFailureRate, deploysPerWeek, formatDuration, prLeadTime, timeToRestore } from "@/features/delivery/stats";
import { Pipelines } from "@/features/delivery/ui/Pipelines";
import { Production } from "@/features/delivery/ui/Production";
import { Pulls } from "@/features/delivery/ui/Pulls";
import { Releases } from "@/features/delivery/ui/Releases";
import { Attention } from "@/shared/ui/attention";
import { KpiRow } from "@/widgets/shell/Kpi";

export const Route = createFileRoute("/delivery")({
  loader: () => loadDelivery(),
  head: () => ({ meta: [{ title: "Delivery · TrieOH dev" }] }),
  component: DeliveryPage,
});

const TABS: { id: DeliveryTab; label: string }[] = [
  { id: "production", label: "Production" },
  { id: "releases", label: "Releases" },
  { id: "pipelines", label: "Pipelines" },
  { id: "pulls", label: "Pull requests" },
];
const tabLabel = (id: DeliveryTab) => TABS.find((tab) => tab.id === id)!.label;

function DeliveryPage() {
  const data = Route.useLoaderData();
  const [tab, setTab] = createSignal<DeliveryTab>("production", { name: "delivery-tab" });
  const attention = createMemo(() => needsAttention(data()), { name: "delivery-attention" });
  const releases = createMemo(() => groupReleases(data()), { name: "delivery-releases" });
  const cfr = () => changeFailureRate(data().deploys);
  let tabs: HTMLDivElement | undefined;
  const open = (next: DeliveryTab) => {
    setTab(next);
    tabs?.scrollIntoView({ behavior: "smooth", block: "start" });
  };

  return (
    <>
      <PageHeader
        title="Delivery"
        description="What runs in production, how each change shipped and where release time goes. From GitHub and the deploy workflow."
      />

      <Attention items={attention()} tabLabel={tabLabel} onOpen={open} />

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

      <div ref={(el) => (tabs = el)} class="mt-10 scroll-mt-6">
        <Tabs ariaLabel="Delivery views" items={TABS} value={tab()} onChange={setTab} />
      </div>

      <div class="mt-6">
        <Switch>
          <Match when={tab() === "production"}>
            <Production production={data().production} deploys={data().deploys} />
          </Match>
          <Match when={tab() === "releases"}>
            <Releases releases={releases()} />
          </Match>
          <Match when={tab() === "pipelines"}>
            <Pipelines ci={data().ci} />
          </Match>
          <Match when={tab() === "pulls"}>
            <Pulls pulls={data().pulls} />
          </Match>
        </Switch>
      </div>
    </>
  );
}
