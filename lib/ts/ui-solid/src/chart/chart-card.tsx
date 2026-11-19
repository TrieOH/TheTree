import type { JSX } from "@solidjs/web";
import {
  For,
  Show,
  createEffect,
  createMemo,
  createSignal,
  untrack,
} from "solid-js";
import EyeIcon from "~icons/lucide/eye";
import { Button } from "../button";
import { Combobox } from "../combobox";
import { SearchInput } from "../search-input";
import { Loading } from "../loading";
import { cn } from "../lib/cn";
import { shapeClass, type Shape } from "../lib/shape";
import { DateRangeControl } from "./date-range-control";
import { computeRangeWindow } from "./range";
import { fillContinuousSeries } from "./data-continuity";
import { GenericChart } from "../generic-chart";
import { buildSeriesMeta } from "./theme";
import type {
  ChartDatum,
  ChartType,
  CurveStyle,
  RangeKey,
  CustomRange,
} from "./types";

export interface ChartCardProps {
  title: string;
  subtitle?: string;
  data: readonly ChartDatum[];
  seriesLabels?: Record<string, string>;
  seriesColors?: Record<string, string>;
  initialType?: ChartType;
  /** Controlled chart type, paired with onTypeChange. */
  type?: ChartType;
  onTypeChange?: (type: ChartType) => void;
  allowedTypes?: readonly ChartType[];
  /** Show the chart type control. Defaults to true. */
  showTypeSelector?: boolean;
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
  loading?: boolean;
  error?: string;
  onRetry?: () => void;
  emptyMessage?: string;
  shape?: Shape;
  class?: string;
}
const typeLabels: Record<ChartType, string> = {
  line: "Linha",
  area: "Área",
  bar: "Barras",
  scatter: "Dispersão",
};
export function ChartCard(props: ChartCardProps): JSX.Element {
  const [internalType, setType] = createSignal<ChartType>(
    untrack(() => props.initialType ?? props.allowedTypes?.[0] ?? "line"),
  );
  const [range, setRange] = createSignal<RangeKey>(
    untrack(() => props.initialRange ?? "all"),
  );
  const [customRange, setCustomRange] = createSignal<CustomRange>({
    from: null,
    to: null,
  });
  const [excluded, setExcluded] = createSignal<ReadonlySet<string>>(new Set());
  const [query, setQuery] = createSignal("");
  const [showPoints, setShowPoints] = createSignal(false);
  const types = () =>
    props.allowedTypes?.length
      ? props.allowedTypes
      : (["line", "area", "bar", "scatter"] as const);
  const type = () => {
    const value = props.type ?? internalType();
    return types().includes(value) ? value : types()[0];
  };
  const validData = createMemo(() =>
    props.data.filter(
      (datum) =>
        datum.date instanceof Date &&
        Number.isFinite(datum.date.getTime()) &&
        Number.isFinite(datum.value),
    ),
  );
  const series = createMemo(() =>
    buildSeriesMeta(
      [...new Set(validData().map((datum) => datum.series))],
      props.seriesLabels,
      props.seriesColors,
    ),
  );
  const visibleSeries = createMemo(() =>
    series().filter(
      (item) =>
        !excluded().has(item.key) &&
        item.label
          .toLocaleLowerCase()
          .includes(query().trim().toLocaleLowerCase()),
    ),
  );
  const window = createMemo(() => computeRangeWindow(range(), customRange()));
  const selectedData = createMemo(() => {
    const keys = new Set(visibleSeries().map((item) => item.key));
    return validData().filter((datum) => keys.has(datum.series));
  });
  const chartData = createMemo(() => {
    const { windowStart, windowEnd } = window();
    const continuous =
      props.continuity !== false && (type() === "line" || type() === "area");
    if (continuous)
      return fillContinuousSeries({
        data: [...selectedData()],
        seriesKeys: visibleSeries().map((item) => item.key),
        windowStart,
        windowEnd,
        mode: "continuous",
      });
    return selectedData().filter(
      (datum) =>
        (!windowStart || datum.date >= windowStart) &&
        (!windowEnd || datum.date <= windowEnd),
    );
  });
  // Discard exclusions belonging to removed series; newly arriving series stay visible.
  createEffect(
    () => series().map((item) => item.key),
    (keys) => {
      setExcluded(
        (previous) =>
          new Set([...previous].filter((key) => keys.includes(key))),
      );
    },
  );
  const toggle = (key: string) =>
    setExcluded((previous) => {
      const next = new Set(previous);
      if (next.has(key)) next.delete(key);
      else next.add(key);
      return next;
    });
  const height = () => Math.max(80, props.height ?? 320);
  return (
    <section
      aria-label={props.title}
      aria-busy={props.loading ? "true" : "false"}
      class={cn(
        "w-full min-w-0 rounded-xl border border-border bg-card p-4 text-card-foreground shadow-xs",
        shapeClass(props.shape, "surface"),
        props.class,
      )}
    >
      <header class="flex min-w-0 flex-wrap items-start justify-between gap-3">
        <div class="min-w-0 flex-1 basis-40">
          <h3 class="text-sm font-semibold wrap-anywhere">{props.title}</h3>
          <Show when={props.subtitle}>
            <p class="mt-1 text-xs leading-relaxed text-muted-foreground wrap-anywhere">
              {props.subtitle}
            </p>
          </Show>
        </div>
        <div class="flex max-w-full flex-wrap items-center gap-2">
          <Show when={props.showTypeSelector !== false && types().length > 1}>
            <Combobox
              class="w-32 max-w-full"
              size="sm"
              shape="soft"
              searchable={false}
              aria-label="Tipo do gráfico"
              placeholder="Tipo"
              value={type()}
              options={types().map((value) => ({
                value,
                label: typeLabels[value],
              }))}
              disabled={props.loading}
              onChange={(value) => {
                setType(value as ChartType);
                props.onTypeChange?.(value as ChartType);
              }}
            />
          </Show>
          <Show when={props.showRangeFilter !== false}>
            <DateRangeControl
              range={range()}
              customRange={customRange()}
              onRangeChange={setRange}
              onCustomRangeChange={setCustomRange}
              dateFormatter={props.dateFormatter}
              disabled={props.loading}
            />
          </Show>
          <Show when={props.showPointsToggle !== false}>
            <Button
              size="sm"
              variant="outline"
              aria-pressed={showPoints() ? "true" : "false"}
              title="Pontos em gráficos de linha ou área"
              disabled={
                props.loading || (type() !== "line" && type() !== "area")
              }
              onClick={() => setShowPoints((value) => !value)}
            >
              Pontos
            </Button>
          </Show>
        </div>
      </header>
      <Show when={props.showSearchFilter}>
        <SearchInput
          class="mt-3"
          aria-label="Buscar série"
          placeholder="Buscar série…"
          value={query()}
          onInput={(event) => setQuery(event.currentTarget.value)}
          disabled={props.loading}
        />
      </Show>
      <Show when={props.showSeriesFilter !== false && series().length > 1}>
        <div
          role="group"
          aria-label="Séries do gráfico"
          class="mt-3 flex flex-wrap gap-2"
        >
          <For each={series()}>
            {(item) => (
              <Button
                size="sm"
                variant="ghost"
                aria-label={`Série ${item.label}`}
                aria-pressed={!excluded().has(item.key) ? "true" : "false"}
                disabled={props.loading}
                onClick={() => toggle(item.key)}
                class={
                  excluded().has(item.key)
                    ? "text-muted-foreground opacity-50"
                    : "text-foreground"
                }
              >
                <span
                  class="size-2 shrink-0 rounded-full"
                  style={{ "background-color": item.color }}
                  aria-hidden="true"
                />
                {item.label}
              </Button>
            )}
          </For>
        </div>
      </Show>
      <div class="mt-4 min-w-0" style={{ height: `${height()}px` }}>
        <Show
          when={!props.error}
          fallback={
            <div
              class="flex flex-col items-center justify-center gap-3 text-center"
              style={{ height: `${height()}px` }}
            >
              <p role="alert" class="text-sm text-muted-foreground">
                {props.error}
              </p>
              <Show when={props.onRetry}>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => props.onRetry?.()}
                >
                  Tentar novamente
                </Button>
              </Show>
            </div>
          }
        >
          <Show
            when={!props.loading}
            fallback={<Loading size="sm" tone="muted" class="h-full" />}
          >
            <Show
              when={!props.isMasked}
              fallback={
                <div
                  class="flex flex-col items-center justify-center gap-3 rounded-lg bg-muted/30 text-center"
                  style={{ height: `${height()}px` }}
                >
                  <p class="text-sm text-muted-foreground">
                    {props.maskedPlaceholder ?? "Valores ocultos"}
                  </p>
                  <Show when={props.onToggleMask}>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => props.onToggleMask?.()}
                    >
                      <EyeIcon class="size-4" />
                      Mostrar valores
                    </Button>
                  </Show>
                </div>
              }
            >
              <GenericChart
                data={chartData()}
                type={type()}
                series={visibleSeries()}
                ariaLabel={props.title}
                height={height()}
                showPoints={showPoints()}
                curveStyle={props.curveStyle}
                valueFormatter={props.valueFormatter}
                dateFormatter={props.dateFormatter}
                tooltipDetails={props.tooltipDetails}
                emptyMessage={props.emptyMessage}
              />
            </Show>
          </Show>
        </Show>
      </div>
    </section>
  );
}
