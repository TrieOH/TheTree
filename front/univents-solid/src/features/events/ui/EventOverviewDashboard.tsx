import type { JSX } from "@solidjs/web";
import { For } from "solid-js";
import CalendarRangeIcon from "~icons/lucide/calendar-range";
import CircleAlertIcon from "~icons/lucide/circle-alert";
import Layers3Icon from "~icons/lucide/layers-3";
import PackageIcon from "~icons/lucide/package";
import ShoppingBagIcon from "~icons/lucide/shopping-bag";
import TicketIcon from "~icons/lucide/ticket";
import WalletIcon from "~icons/lucide/wallet";
import { ChartCard } from "@/widgets/ui/ChartCard";
import { DashboardBarList } from "@/widgets/ui/DashboardBarList";
import { DashboardPanel } from "@/widgets/ui/DashboardPanel";
import { DashboardStatCard } from "@/widgets/ui/DashboardStatCard";
import type { EventOverviewMetrics } from "../model/event-overview";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideCalendarRange = CalendarRangeIcon as unknown as IconComp;
const LucideCircleAlert = CircleAlertIcon as unknown as IconComp;
const LucideLayers3 = Layers3Icon as unknown as IconComp;
const LucidePackage = PackageIcon as unknown as IconComp;
const LucideShoppingBag = ShoppingBagIcon as unknown as IconComp;
const LucideTicket = TicketIcon as unknown as IconComp;
const LucideWallet = WalletIcon as unknown as IconComp;

const currency = new Intl.NumberFormat("pt-BR", {
  style: "currency",
  currency: "BRL",
});

const wholeCurrency = new Intl.NumberFormat("pt-BR", {
  style: "currency",
  currency: "BRL",
  maximumFractionDigits: 0,
});

const purchaseStatuses = [
  { label: "Aprovadas", status: "approved", color: "bg-emerald-500" },
  { label: "Pendentes", status: "pending", color: "bg-amber-500" },
  { label: "Reembolsadas", status: "refunded", color: "bg-sky-500" },
  { label: "Expiradas", status: "expired", color: "bg-slate-400" },
  { label: "Canceladas", status: "cancelled", color: "bg-rose-500" },
] as const;

export function EventOverviewDashboard(props: {
  metrics: EventOverviewMetrics;
}): JSX.Element {
  const m = () => props.metrics;

  const revenueBars = () =>
    m()
      .editionSales.slice(0, 6)
      .map((edition) => ({
        id: edition.name,
        label: edition.name,
        value: edition.revenue,
        detail: currency.format(edition.revenue / 100),
      }));

  const statusBars = () =>
    purchaseStatuses.map((item) => ({
      id: item.status,
      label: item.label,
      value: m().purchases.filter(
        (purchase) => purchase.status === item.status,
      ).length,
      color: item.color,
    }));

  const summaryMetrics = () => [
    {
      label: "Receita aprovada",
      value: currency.format(m().revenue / 100),
      hint: "Somente compras aprovadas",
      icon: (p: { class?: string }) => <LucideWallet class={p.class ?? "size-4"} />,
    },
    {
      label: "Participantes",
      value: m().participantCount,
      hint: "Ingressos comprados",
      icon: (p: { class?: string }) => <LucideTicket class={p.class ?? "size-4"} />,
    },
    {
      label: "Produtos Vendidos",
      value: m().purchases.length,
      hint: "Compras realizadas",
      icon: (p: { class?: string }) => <LucideShoppingBag class={p.class ?? "size-4"} />,
    },
    {
      label: "Reembolsadas",
      value: m().refundedPurchaseCount,
      hint: "Compras reembolsadas",
      icon: (p: { class?: string }) => <LucideCircleAlert class={p.class ?? "size-4"} />,
    },
  ];

  const catalogMetrics = () => [
    {
      label: "Ingressos",
      value: m().ticketCount,
      hint: "Tipos cadastrados",
      icon: (p: { class?: string }) => <LucideTicket class={p.class ?? "size-4"} />,
    },
    {
      label: "Produtos",
      value: m().productCount,
      hint: "Produtos cadastrados",
      icon: (p: { class?: string }) => <LucidePackage class={p.class ?? "size-4"} />,
    },
    {
      label: "Programas",
      value: m().programCount,
      hint: "Atividades cadastradas",
      icon: (p: { class?: string }) => <LucideCalendarRange class={p.class ?? "size-4"} />,
    },
    {
      label: "Ocorrências",
      value: m().occurrenceCount,
      hint: "Horários da programação",
      icon: (p: { class?: string }) => <LucideLayers3 class={p.class ?? "size-4"} />,
    },
  ];

  return (
    <>
      <section class="order-2 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        <For each={summaryMetrics()}>
          {(metric) => <DashboardStatCard {...metric} />}
        </For>
      </section>

      <DashboardPanel
        title="Catálogo"
        description="Conteúdo publicado nas edições deste evento."
        icon={(p) => <LucidePackage class={p.class} />}
        class="order-4"
      >
        <div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
          <For each={catalogMetrics()}>
            {(metric) => <DashboardStatCard {...metric} />}
          </For>
        </div>
      </DashboardPanel>

      <section class="order-5 grid gap-4 xl:grid-cols-[1fr_1.35fr]">
        <DashboardPanel
          title="Status das compras"
          description={`${m().purchases.length} compra${m().purchases.length === 1 ? "" : "s"} no total.`}
          icon={(p) => <LucideShoppingBag class={p.class} />}
          class="rounded-xl border border-border bg-card p-5 shadow-xs"
        >
          <div class="mt-2">
            <DashboardBarList
              items={statusBars()}
              emptyMessage="Nenhuma compra registrada."
            />
          </div>
        </DashboardPanel>

        <DashboardPanel
          title="Receita por edição"
          description="Receita aprovada comparada entre as edições."
          icon={(p) => <LucideWallet class={p.class} />}
          class="rounded-xl border border-border bg-card p-5 shadow-xs"
        >
          <div class="mt-2">
            <DashboardBarList
              items={revenueBars()}
              maxValue={m().maxEditionRevenue}
              emptyMessage="Ainda não há edições para comparar."
            />
          </div>
        </DashboardPanel>
      </section>

      {/* Gráfico de lucro acumulado */}
      <section class="order-6">
        <ChartCard
          title="Lucro"
          subtitle="Crescimento acumulado no período selecionado."
          data={m().profitData}
          allowedTypes={["line"]}
          initialRange="30d"
          showSeriesFilter={false}
          showSearchFilter={false}
          continuity
          showPointsToggle={false}
          seriesLabels={{ Lucro: "Lucro acumulado" }}
          tooltipDetails={(datum) => [
            {
              label: "Compras",
              value: String(datum.purchases ?? 0),
            },
          ]}
          valueFormatter={(value) => wholeCurrency.format(value)}
        />
      </section>
    </>
  );
}
