import {
  ChartCard,
  DashboardBarList,
  DashboardStatCard,
} from "@trieoh/ui-base";
import {
  CalendarRange,
  CircleAlert,
  Layers3,
  Package,
  ShoppingBag,
  Ticket,
  Wallet,
} from "lucide-react";
import { DashboardPanel } from "@/widgets/ui/dashboard-panel";
import type { EventOverviewMetrics } from "../model/event-overview";

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

export function EventOverviewDashboard({
  metrics,
}: {
  metrics: EventOverviewMetrics;
}) {
  const revenueBars = metrics.editionSales.slice(0, 6).map((edition) => ({
    id: edition.name,
    label: edition.name,
    value: edition.revenue,
    detail: currency.format(edition.revenue / 100),
  }));
  const statusBars = purchaseStatuses.map((item) => ({
    id: item.status,
    label: item.label,
    value: metrics.purchases.filter(
      (purchase) => purchase.status === item.status,
    ).length,
    color: item.color,
  }));
  const summaryMetrics = [
    {
      label: "Receita aprovada",
      value: currency.format(metrics.revenue / 100),
      hint: "Somente compras aprovadas",
      icon: Wallet,
    },
    {
      label: "Participantes",
      value: metrics.participantCount,
      hint: "Ingressos comprados",
      icon: Ticket,
    },
    {
      label: "Produtos Vendidos",
      value: metrics.purchases.length,
      hint: "Compras realizadas",
      icon: ShoppingBag,
    },
    {
      label: "Reembolsadas",
      value: metrics.refundedPurchaseCount,
      hint: "Compras reembolsadas",
      icon: CircleAlert,
    },
  ];
  const catalogMetrics = [
    {
      label: "Ingressos",
      value: metrics.ticketCount,
      hint: "Tipos cadastrados",
      icon: Ticket,
    },
    {
      label: "Produtos",
      value: metrics.productCount,
      hint: "Produtos cadastrados",
      icon: Package,
    },
    {
      label: "Programas",
      value: metrics.programCount,
      hint: "Atividades cadastradas",
      icon: CalendarRange,
    },
    {
      label: "Ocorrências",
      value: metrics.occurrenceCount,
      hint: "Horários da programação",
      icon: Layers3,
    },
  ];

  return (
    <>
      <section className="order-2 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        {summaryMetrics.map((metric) => (
          <DashboardStatCard key={metric.label} {...metric} />
        ))}
      </section>

      <DashboardPanel
        title="Catálogo"
        description="Conteúdo publicado nas edições deste evento."
        icon={Package}
        className="order-4"
      >
        <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
          {catalogMetrics.map((metric) => (
            <DashboardStatCard key={metric.label} {...metric} />
          ))}
        </div>
      </DashboardPanel>

      <section className="order-5 grid gap-4 xl:grid-cols-[1fr_1.35fr]">
        <DashboardPanel
          title="Status das compras"
          description={`${metrics.purchases.length} compra${metrics.purchases.length === 1 ? "" : "s"} no total.`}
          icon={ShoppingBag}
          className="rounded-lg bg-card p-5 ring-1 ring-foreground/10"
        >
          <div className="mt-2">
            <DashboardBarList
              items={statusBars}
              emptyMessage="Nenhuma compra registrada."
            />
          </div>
        </DashboardPanel>
        <DashboardPanel
          title="Receita por edição"
          description="Receita aprovada comparada entre as edições."
          icon={Wallet}
          className="rounded-lg bg-card p-5 ring-1 ring-foreground/10"
        >
          <div className="mt-2">
            <DashboardBarList
              items={revenueBars}
              maxValue={metrics.maxEditionRevenue}
              emptyMessage="Ainda não há edições para comparar."
            />
          </div>
        </DashboardPanel>
      </section>

      <section className="order-6">
        <ChartCard
          title="Lucro"
          subtitle="Crescimento acumulado no período selecionado."
          data={metrics.profitData}
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
