import { createMemo } from "solid-js";
import Wallet from "~icons/lucide/wallet";
import ChartNoAxesCombined from "~icons/lucide/chart-no-axes-combined";
import {
  DashboardStatCard,
  DashboardPanel,
  DashboardBarList,
  MetricDelta,
  Sparkline,
  ChartCard,
  GenericChart,
  Button,
  type ChartDatum,
  type ChartType,
  type CurveStyle,
  type Shape,
} from "@trieoh/ui-solid";
import type { Values } from "../catalog";
type Props = {
  values: Values;
  update: (key: string, value: string | boolean) => void;
};
const history = [32, 27, 40, 37, 54, 51, 68, 64, 85, 83, 102, 118];
const bars = [
  { id: "online", label: "Inscrições online", value: 680 },
  { id: "partners", label: "Parceiros", value: 240 },
  { id: "venue", label: "No local", value: 80 },
];
const currency = (value: number) =>
  new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
    maximumFractionDigits: 0,
  }).format(value);
export function demoChartData(): ChartDatum[] {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  return Array.from({ length: 21 }, (_, index) => {
    const date = new Date(today);
    date.setDate(date.getDate() - 20 + index);
    return [
      {
        date,
        value: 80 + index * 13 + Math.round(Math.sin(index) * 20),
        series: "receita",
      },
      {
        date,
        value: 45 + index * 6 + Math.round(Math.cos(index) * 10),
        series: "despesas",
      },
    ];
  }).flat();
}
export function StatCardExample(props: Props) {
  return (
    <div class="w-full max-w-100">
      <DashboardStatCard
        label={String(props.values.label)}
        value={currency(Number(props.values.value))}
        hint={String(props.values.hint)}
        icon={Wallet}
        loading={!!props.values.loading}
        shape={props.values.shape as Shape}
        comparison={
          props.values.comparison
            ? {
                value: Number(props.values.value),
                previousValue: Number(props.values.previousValue),
                label: "vs. mês anterior",
              }
            : undefined
        }
        sparkline={
          props.values.sparkline
            ? { data: history, label: "Evolução da receita", type: "area" }
            : undefined
        }
      />
    </div>
  );
}
export function MetricDeltaExample(props: Props) {
  return (
    <MetricDelta
      value={Number(props.values.value)}
      previousValue={Number(props.values.previousValue)}
      improvement={props.values.improvement as "increase" | "decrease"}
      label={String(props.values.label)}
    />
  );
}
export function SparklineExample(props: Props) {
  return (
    <div data-wide class="w-full max-w-120">
      <Sparkline
        data={
          props.values.empty
            ? []
            : props.values.flat
              ? history.map(() => 50)
              : history
        }
        label="Evolução das inscrições"
        type={props.values.type as "line" | "area"}
        height={Number(props.values.height)}
        showTooltip={props.values.showTooltip !== false}
      />
    </div>
  );
}
export function BarListExample(props: Props) {
  return (
    <div class="w-full">
      <DashboardBarList
        items={props.values.empty ? [] : bars}
        total={1000}
        showPercentage={!!props.values.showPercentage}
        loading={!!props.values.loading}
      />
    </div>
  );
}
export function PanelExample(props: Props) {
  return (
    <div data-wide class="w-full">
      <DashboardPanel
        title={String(props.values.title)}
        description={String(props.values.description)}
        icon={ChartNoAxesCombined}
        loading={!!props.values.loading}
        error={
          props.values.error
            ? "Não foi possível carregar as métricas."
            : undefined
        }
        onRetry={() => props.update("error", false)}
        action={
          <Button
            size="sm"
            variant="outline"
            onClick={() => props.update("loading", !props.values.loading)}
          >
            Atualizar
          </Button>
        }
      >
        <div class="grid min-w-0 grid-cols-[repeat(auto-fit,minmax(min(100%,12rem),1fr))] gap-3">
          <DashboardStatCard
            label="Inscrições"
            value="1.000"
            hint="Inscrições confirmadas"
            icon={Wallet}
            comparison={{ value: 1000, previousValue: 820 }}
          />
          <DashboardStatCard
            label="Receita"
            value="R$ 24.850"
            hint="Receita aprovada"
            icon={Wallet}
            sparkline={{ data: history, height: 40 }}
          />
        </div>
        <div class="mt-4 rounded-xl border border-border bg-card p-4">
          <DashboardBarList items={bars} total={1000} showPercentage />
        </div>
      </DashboardPanel>
    </div>
  );
}
export function ChartExample(props: Props & { bare?: boolean }) {
  const initial = demoChartData();
  const data = createMemo(() =>
    props.values.empty
      ? []
      : props.values.multiSeries === false
        ? initial.filter((row) => row.series === "receita")
        : initial,
  );
  return (
    <div data-wide class="w-full min-w-0">
      {props.bare ? (
        <GenericChart
          data={data()}
          type={props.values.type as ChartType}
          ariaLabel="Receita e despesas por dia"
          height={Number(props.values.height)}
          showPoints={!!props.values.showPoints}
          curveStyle={props.values.curveStyle as CurveStyle}
          valueFormatter={currency}
        />
      ) : (
        <ChartCard
          title={String(props.values.title)}
          subtitle="Receita e despesas ao longo do período."
          data={data()}
          type={props.values.type as ChartType}
          onTypeChange={(value) => props.update("type", value)}
          seriesLabels={{ receita: "Receita", despesas: "Despesas" }}
          valueFormatter={currency}
          height={Number(props.values.height)}
          shape={props.values.shape as Shape}
          showTypeSelector={props.values.showTypeSelector !== false}
          showRangeFilter={!!props.values.showRangeFilter}
          showSeriesFilter={!!props.values.showSeriesFilter}
          showSearchFilter={!!props.values.showSearchFilter}
          showPointsToggle={!!props.values.showPointsToggle}
          continuity={false}
          loading={!!props.values.loading}
          error={
            props.values.error
              ? "Não foi possível carregar o gráfico."
              : undefined
          }
          onRetry={() => props.update("error", false)}
          isMasked={!!props.values.masked}
          onToggleMask={() => props.update("masked", false)}
        />
      )}
    </div>
  );
}
