import type { JSX } from "@solidjs/web";
import {
  Show,
  createMemo,
  createSignal,
  untrack,
} from "solid-js";
import EyeIcon from "~icons/lucide/eye";
import { DateRangeControl } from "./chart-filter-bar";
import { fillContinuousSeries } from "./data-continuity";
import { GenericChart } from "./generic-chart";
import { buildSeriesMeta } from "./theme";
import type { ChartDatum, ChartType, CurveStyle, RangeKey } from "./types";
import { useChartFilters } from "./use-chart-filters";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideEye = EyeIcon as unknown as IconComp;

export interface ChartCardProps {
  title: string;
  subtitle?: string;
  data: ChartDatum[];
  seriesLabels?: Record<string, string>;
  seriesColors?: Record<string, string>;

  initialType?: ChartType;
  allowedTypes?: ChartType[];

  initialRange?: RangeKey;
  showRangeFilter?: boolean;
  showSeriesFilter?: boolean;
  showSearchFilter?: boolean;
  showPointsToggle?: boolean;

  continuity?: boolean;
  curveStyle?: CurveStyle;
  height?: number;
  valueFormatter?: (value: number) => string;
  dateFormatter?: (date: Date) => string;
  tooltipDetails?: (datum: ChartDatum) => { label: string; value: string }[];
  isMasked?: boolean;
  onToggleMask?: () => void;
  maskedPlaceholder?: string;
}

export function ChartCard(props: ChartCardProps): JSX.Element {
  const seriesKeys = createMemo(() =>
    Array.from(new Set(props.data.map((d) => d.series))),
  );

  const series = createMemo(() =>
    buildSeriesMeta(seriesKeys(), props.seriesLabels, props.seriesColors),
  );

  const [showPoints, setShowPoints] = createSignal(false);

  const lockedType = () =>
    props.allowedTypes && props.allowedTypes.length === 1
      ? props.allowedTypes[0]
      : undefined;

  const {
    type,
    range,
    setRange,
    customRange,
    setCustomRange,
    visible,
    windowStart,
    windowEnd,
    seriesFilteredData,
    filteredData,
  } = useChartFilters({
    data: () => props.data,
    seriesKeys,
    initialType: untrack(
      () => props.allowedTypes?.[0] ?? props.initialType ?? "line",
    ),
    initialRange: untrack(() => props.initialRange ?? "all"),
  });

  const effectiveType = () => lockedType() ?? type();

  const visibleSeries = createMemo(() =>
    series().filter((s) => visible().has(s.key)),
  );

  const isContinuousType = () => {
    const t = effectiveType();
    return t === "line" || t === "area";
  };

  const chartData = createMemo(() => {
    if (!isContinuousType()) return filteredData();
    return fillContinuousSeries({
      data: seriesFilteredData(),
      seriesKeys: visibleSeries().map((s) => s.key),
      windowStart: windowStart(),
      windowEnd: windowEnd(),
      mode: props.continuity !== false ? "continuous" : "none",
    });
  });

  const canTogglePoints = () =>
    (props.showPointsToggle ?? true) &&
    (effectiveType() === "line" || effectiveType() === "area");

  return (
    <div class="w-full rounded-2xl border border-border bg-card p-4 text-card-foreground shadow-xs sm:p-5">
      <div class="flex flex-col gap-3 sm:flex-row! sm:items-start! sm:justify-between!">
        <div class="min-w-0 flex-1">
          <h3 class="text-sm font-semibold">{props.title}</h3>
          <Show when={props.subtitle}>
            <p class="mt-0.5 text-xs text-muted-foreground">
              {props.subtitle}
            </p>
          </Show>
        </div>

        <div class="flex w-full flex-wrap items-center justify-start gap-2 sm:w-auto! sm:justify-end!">
          <Show when={canTogglePoints()}>
            <button
              type="button"
              onClick={() => setShowPoints((v) => !v)}
              aria-pressed={showPoints() ? "true" : "false"}
              class={`shrink-0 rounded-full border px-2.5 py-1 text-[11px] font-medium transition-colors ${showPoints()
                  ? "border-primary bg-primary text-primary-foreground"
                  : "border-border text-muted-foreground hover:bg-muted hover:text-foreground"
                }`}
            >
              Pontos
            </button>
          </Show>

          <Show when={props.showRangeFilter !== false}>
            <DateRangeControl
              range={range()}
              onRangeChange={setRange}
              customRange={customRange()}
              onCustomRangeChange={setCustomRange}
              dateFormatter={props.dateFormatter}
            />
          </Show>
        </div>
      </div>

      <div class="relative mt-4">
        <div
          class={props.isMasked ? "pointer-events-none select-none filter blur-md transition-all" : "transition-all"}
        >
          <GenericChart
            data={chartData()}
            type={effectiveType()}
            series={visibleSeries()}
            ariaLabel={props.title}
            height={props.height ?? 320}
            showPoints={showPoints()}
            curveStyle={props.curveStyle}
            valueFormatter={props.valueFormatter}
            dateFormatter={props.dateFormatter}
            tooltipDetails={props.tooltipDetails}
          />
        </div>

        <Show when={props.isMasked}>
          <div class="absolute inset-0 z-10 flex flex-col items-center justify-center gap-2 rounded-xl bg-background/60 p-4 text-center backdrop-blur-xs">
            <p class="text-sm font-medium text-foreground">
              {props.maskedPlaceholder ?? "Valores ocultos"}
            </p>
            <Show when={props.onToggleMask}>
              <button
                type="button"
                onClick={() => props.onToggleMask?.()}
                class="inline-flex items-center gap-1.5 rounded-md border border-border bg-background px-3 py-1.5 text-xs font-medium text-foreground shadow-xs transition-colors hover:bg-muted"
              >
                <LucideEye class="size-3.5" />
                <span>Mostrar valores</span>
              </button>
            </Show>
          </div>
        </Show>
      </div>
    </div>
  );
}
