import {
  defineChart,
  lineY,
  areaY,
  barY,
  dot,
  group,
} from "@tanstack/charts";
import { createChartRendererAdapter } from "@tanstack/charts/adapter/renderer";
import { d3Curve } from "@tanstack/charts/d3/shape";
import { scaleBand } from "@tanstack/charts/scales/band";
import { scaleLinear } from "@tanstack/charts/scales/linear";
import { scaleOrdinal } from "@tanstack/charts/scales/ordinal";
import { renderChartSvg } from "@tanstack/charts/svg";
import { createSvgChartRenderer } from "@tanstack/charts/svg/renderer";
import { tooltip } from "@tanstack/charts/tooltip";
import { portal } from "@tanstack/charts/tooltip/portal";
import { scaleUtc } from "d3-scale";
import { curveMonotoneX } from "d3-shape";
import {
  For,
  Show,
  createEffect,
  createMemo,
  createSignal,
  onCleanup,
} from "solid-js";
import { Portal, type JSX } from "@solidjs/web";
import type { ChartDatum, ChartType, CurveStyle, SeriesMeta } from "./types";

const smoothCurve = d3Curve(curveMonotoneX);

const xAccessor = (d: ChartDatum): Date => d.date;
const yAccessor = (d: ChartDatum): number => d.value;
const zAccessor = (d: ChartDatum): string => d.series;
const colorAccessor = (d: ChartDatum): string => d.series;

const defaultValueFormatter = (value: number) =>
  new Intl.NumberFormat("pt-BR", {
    notation: "compact",
    maximumFractionDigits: 1,
  }).format(value);

const defaultDateFormatter = (date: Date) =>
  new Intl.DateTimeFormat("pt-BR", {
    day: "2-digit",
    month: "short",
  }).format(date);

interface TooltipPoint {
  xValue: Date;
  yValue: number;
  datum: ChartDatum;
}

interface TooltipTarget {
  element: HTMLElement;
  points: TooltipPoint[];
}

export interface GenericChartProps {
  data: ChartDatum[];
  type: ChartType;
  series: SeriesMeta[];
  ariaLabel: string;
  height?: number;
  showPoints?: boolean;
  curveStyle?: CurveStyle;
  valueFormatter?: (value: number) => string;
  dateFormatter?: (date: Date) => string;
  tooltipDetails?: (datum: ChartDatum) => { label: string; value: string }[];
  emptyMessage?: string;
}

export function GenericChart(props: GenericChartProps): JSX.Element {
  const valueFmt = () => props.valueFormatter ?? defaultValueFormatter;
  const dateFmt = () => props.dateFormatter ?? defaultDateFormatter;

  const [tooltipState, setTooltipState] = createSignal<TooltipTarget | null>(
    null,
  );
  let containerRef: HTMLDivElement | undefined;
  let adapterInstance: ReturnType<typeof createChartRendererAdapter> | null = null;

  const seriesKeys = createMemo(() => props.series.map((s) => s.key));
  const colorRange = createMemo(() => props.series.map((s) => s.color));
  const multiSeries = createMemo(() => seriesKeys().length > 1);

  const definition = createMemo(() => {
    const d = props.data;
    const t = props.type;
    const showPts = props.showPoints ?? false;
    const c = (props.curveStyle ?? "smooth") === "smooth" ? smoothCurve : undefined;
    const dFmt = dateFmt();
    const vFmt = valueFmt();
    const keys = seriesKeys();
    const colors = colorRange();

    const marks = (() => {
      switch (t) {
        case "area":
          return [
            areaY(d, {
              x: xAccessor,
              y1: 0,
              y2: yAccessor,
              z: zAccessor,
              color: colorAccessor,
              fillOpacity: 0.14,
              curve: c,
            }),
            lineY(d, {
              x: xAccessor,
              y: yAccessor,
              z: zAccessor,
              color: colorAccessor,
              points: showPts,
              strokeWidth: 2,
              curve: c,
            }),
          ];
        case "bar":
          return [
            barY(d, {
              x: xAccessor,
              y: yAccessor,
              z: zAccessor,
              color: colorAccessor,
              radius: 4,
              layout: multiSeries() ? group({ padding: 0.24 }) : undefined,
            }),
          ];
        case "scatter":
          return [
            dot(d, {
              x: xAccessor,
              y: yAccessor,
              z: zAccessor,
              color: colorAccessor,
              r: 4,
              states: [
                {
                  when: { focus: "primary" },
                  style: { r: 7 },
                  transition: {
                    type: "tween",
                    duration: 140,
                    easing: "ease-out",
                  },
                },
              ],
            }),
          ];
        case "line":
        default:
          return [
            lineY(d, {
              x: xAccessor,
              y: yAccessor,
              z: zAccessor,
              color: colorAccessor,
              points: showPts,
              strokeWidth: 2.25,
              curve: c,
            }),
          ];
      }
    })();

    return defineChart({
      marks,
      x:
        t === "bar"
          ? {
            scale: () => scaleBand().padding(0.28),
            axis: { ticks: { format: dFmt } },
          }
          : {
            scale: scaleUtc,
            nice: false,
            axis: { ticks: { format: dFmt } },
          },
      y: {
        scale: scaleLinear,
        nice: true,
        grid: true,
        axis: { ticks: { format: vFmt } },
      },
      color: {
        scale: scaleOrdinal(keys, colors),
      },
      focus: "group-x",
      tooltip: {
        use: tooltip,
        portal,
        className:
          "[--ts-chart-tooltip-background:var(--popover)] [--ts-chart-tooltip-color:var(--popover-foreground)] [--ts-chart-tooltip-border:1px_solid_var(--border)] [--ts-chart-tooltip-border-radius:0.625rem]",
        anchor: "group-center",
        placement: ["top", "right", "bottom", "left"],
      },
    });
  });

  const renderer = createSvgChartRenderer(renderChartSvg);
  const idPrefix = `ts-chart-${Math.random().toString(36).slice(2, 8)}`;

  const buildHostOptions = (def: ReturnType<typeof definition>) => ({
    renderer,
    definition: def,
    ariaLabel: props.ariaLabel,
    height: props.height ?? 320,
    width: undefined,
    idPrefix,
    onTooltipBodyChange: (target: TooltipTarget | null) => {
      setTooltipState(target);
    },
  });

  createEffect(
    () => definition(),
    (def) => {
      const el = containerRef;
      if (!el) return;

      const opts = buildHostOptions(def) as unknown as Parameters<typeof createChartRendererAdapter>[0];
      if (!adapterInstance) {
        adapterInstance = createChartRendererAdapter(opts);
        adapterInstance.mount(el);
      }
      adapterInstance.update(opts);
    },
  );

  onCleanup(() => {
    if (adapterInstance) {
      adapterInstance.destroy();
      adapterInstance = null;
    }
  });

  const showTotal = () => props.series.length > 1;

  return (
    <Show
      when={props.data.length > 0}
      fallback={
        <div
          class="flex items-center justify-center rounded-xl border border-dashed border-border text-sm text-muted-foreground"
          style={{ height: `${props.height ?? 320}px` }}
        >
          {props.emptyMessage ?? "Nenhum dado para os filtros selecionados."}
        </div>
      }
    >
      <div
        class="ts-chart-host text-muted-foreground"
        style={{
          position: "relative",
          width: "100%",
          height: `${props.height ?? 320}px`,
        }}
      >
        <div
          ref={(el) => (containerRef = el)}
          class="ts-chart-surface"
          style={{ width: "100%", height: "100%" }}
        />

        <Show when={tooltipState()}>
          {(target) => {
            const points = () => target().points ?? [];
            const rows = () =>
              points().map((p) => {
                const meta = props.series.find(
                  (s) => s.key === p.datum.series,
                );
                return {
                  key: p.datum.series,
                  label: meta?.label ?? p.datum.series,
                  color: meta?.color ?? "#94a3b8",
                  value: p.yValue,
                  synthetic: Boolean(p.datum.synthetic),
                };
              });
            const total = () =>
              rows().reduce((sum, r) => sum + r.value, 0);

            return (
              <Portal mount={target().element}>
                <div class="min-w-36 px-1 py-0.5">
                  <p class="mb-1.5 text-[11px] font-medium text-muted-foreground">
                    Período:{" "}
                    {points()[0] ? dateFmt()(points()[0].xValue) : ""}
                  </p>
                  <ul class="space-y-1">
                    <For each={rows()}>
                      {(r) => (
                        <li class="flex items-center justify-between gap-4 text-xs">
                          <span class="flex items-center gap-1.5 text-muted-foreground">
                            <span
                              class="size-1.5 shrink-0 rounded-full"
                              style={{ "background-color": r.color }}
                            />
                            {r.label}
                            <Show when={r.synthetic}>
                              <span
                                class="text-muted-foreground/50"
                                title="Valor mantido — sem novo dado neste dia"
                              >
                                ⋯
                              </span>
                            </Show>
                          </span>
                          <span class="font-mono font-medium tabular-nums text-foreground">
                            {valueFmt()(r.value)}
                          </span>
                        </li>
                      )}
                    </For>
                  </ul>
                  <Show when={points()[0] && props.tooltipDetails}>
                    <For each={props.tooltipDetails!(points()[0].datum)}>
                      {(detail) => (
                        <div class="mt-1 flex items-center justify-between gap-4 text-xs">
                          <span class="text-muted-foreground">
                            {detail.label}
                          </span>
                          <span class="font-mono font-medium tabular-nums text-foreground">
                            {detail.value}
                          </span>
                        </div>
                      )}
                    </For>
                  </Show>
                  <Show when={showTotal()}>
                    <div class="mt-1.5 flex items-center justify-between gap-4 border-t border-border pt-1.5 text-xs">
                      <span class="font-medium text-muted-foreground">
                        Total
                      </span>
                      <span class="font-mono font-semibold tabular-nums text-foreground">
                        {valueFmt()(total())}
                      </span>
                    </div>
                  </Show>
                </div>
              </Portal>
            );
          }}
        </Show>
      </div>
    </Show>
  );
}
