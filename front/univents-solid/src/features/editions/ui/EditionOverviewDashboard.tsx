import type { JSX } from "@solidjs/web";
import { createMemo, For, Show } from "solid-js";
import CalendarRangeIcon from "~icons/lucide/calendar-range";
import CircleAlertIcon from "~icons/lucide/circle-alert";
import EyeIcon from "~icons/lucide/eye";
import EyeOffIcon from "~icons/lucide/eye-off";
import Layers3Icon from "~icons/lucide/layers-3";
import PackageIcon from "~icons/lucide/package";
import ShoppingBagIcon from "~icons/lucide/shopping-bag";
import TicketIcon from "~icons/lucide/ticket";
import WalletIcon from "~icons/lucide/wallet";
import { ChartCard } from "@/widgets/ui/ChartCard";
import { DashboardBarList } from "@/widgets/ui/DashboardBarList";
import { DashboardPanel } from "@/widgets/ui/DashboardPanel";
import { DashboardStatCard } from "@/widgets/ui/DashboardStatCard";
import { useMonetaryVisibility } from "@/features/events/lib/use-monetary-visibility";
import type { EditionOverviewMetrics } from "../model/edition-overview";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideCalendarRange = CalendarRangeIcon as unknown as IconComp;
const LucideCircleAlert = CircleAlertIcon as unknown as IconComp;
const LucideEye = EyeIcon as unknown as IconComp;
const LucideEyeOff = EyeOffIcon as unknown as IconComp;
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

export function EditionOverviewDashboard(props: {
  metrics: EditionOverviewMetrics;
}): JSX.Element {
  const m = () => props.metrics;
  const { showMonetary, toggleMonetary } = useMonetaryVisibility();

  const profitValueFormatter = createMemo(() => {
    const show = showMonetary();
    return (value: number) => (show ? wholeCurrency.format(value) : "••••");
  });

  const statusBars = () =>
    purchaseStatuses.map((item) => ({
      id: item.status,
      label: item.label,
      value: m().statusCounts[item.status] ?? 0,
      color: item.color,
    }));

  const summaryMetrics = () => [
    {
      label: "Receita aprovada",
      value: (
        <span class="inline-flex items-center gap-2">
          <span>{showMonetary() ? currency.format(m().revenue / 100) : "••••"}</span>
        </span>
      ),
      hint: showMonetary() ? "Compras aprovadas" : "Valores monetários ocultos",
      icon: (p: { class?: string }) => <LucideWallet class={p.class ?? "size-4"} />,
      action: (
        <button
          type="button"
          onClick={toggleMonetary}
          title={showMonetary() ? "Ocultar valores monetários" : "Exibir valores monetários"}
          aria-label={showMonetary() ? "Ocultar valores monetários" : "Exibir valores monetários"}
          class="rounded p-1 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
        >
          <Show when={showMonetary()} fallback={<LucideEyeOff class="size-3.5" />}>
            <LucideEye class="size-3.5" />
          </Show>
        </button>
      ),
    },
    {
      label: "Participantes",
      value: m().attendeeCount,
      hint: "Ingressos comprados",
      icon: (p: { class?: string }) => <LucideTicket class={p.class ?? "size-4"} />,
    },
    {
      label: "Produtos Vendidos",
      value: m().purchasesCount,
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
      hint: "Horários cadastrados",
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
        description="Ingressos, produtos e programação desta edição."
        icon={(p) => <LucidePackage class={p.class} />}
        class="order-3"
      >
        <div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
          <For each={catalogMetrics()}>
            {(metric) => <DashboardStatCard {...metric} />}
          </For>
        </div>
      </DashboardPanel>

      <section class="order-4 grid gap-4 xl:grid-cols-[1fr_1.8fr]">
        <DashboardPanel
          title="Status das compras"
          description={`${m().purchasesCount} compra${m().purchasesCount === 1 ? "" : "s"} registrada${m().purchasesCount === 1 ? "" : "s"}.`}
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

        <ChartCard
          title="Lucro"
          subtitle="Crescimento acumulado das compras da edição."
          data={m().profitData}
          allowedTypes={["line"]}
          initialRange="30d"
          showSeriesFilter={false}
          showSearchFilter={false}
          showPointsToggle={false}
          seriesLabels={{ revenue: "Lucro acumulado" }}
          seriesColors={{ revenue: "#10b981" }}
          tooltipDetails={(datum) => [
            { label: "Compras", value: String(datum.purchases ?? 0) },
          ]}
          valueFormatter={profitValueFormatter()}
          isMasked={!showMonetary()}
          onToggleMask={toggleMonetary}
          maskedPlaceholder="Valores financeiros ocultos por padrão"
        />
      </section>
    </>
  );
}
