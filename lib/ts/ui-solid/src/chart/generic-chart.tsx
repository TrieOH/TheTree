import { defineChart, lineY, areaY, barY, dot, group } from "@tanstack/charts";
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
import { buildSeriesMeta } from "./theme";
import { cn } from "../lib/cn";
import type { ChartDatum, ChartType, CurveStyle, SeriesMeta } from "./types";

const smoothCurve = /* @__PURE__ */ d3Curve(curveMonotoneX);

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

type HostOptions = Parameters<
  typeof createChartRendererAdapter<ChartDatum, Date, number>
>[0];
type TooltipTarget = Parameters<
  NonNullable<HostOptions["onTooltipBodyChange"]>
>[0];

export interface GenericChartProps {
  data: readonly ChartDatum[];
  type?: ChartType;
  series?: readonly SeriesMeta[];
  ariaLabel: string;
  height?: number;
  showPoints?: boolean;
  curveStyle?: CurveStyle;
  valueFormatter?: (value: number) => string;
  dateFormatter?: (date: Date) => string;
  tooltipDetails?: (datum: ChartDatum) => { label: string; value: string }[];
  emptyMessage?: string;
  showAxes?: boolean;
  showTooltip?: boolean;
  class?: string;
}

export function GenericChart(props: GenericChartProps): JSX.Element {
  const valueFmt = () => props.valueFormatter ?? defaultValueFormatter;
  const dateFmt = () => props.dateFormatter ?? defaultDateFormatter;

  const [tooltipState, setTooltipState] = createSignal<TooltipTarget>(null);
  const [surface, setSurface] = createSignal<HTMLDivElement>();
  let adapterInstance: ReturnType<
    typeof createChartRendererAdapter<ChartDatum, Date, number>
  > | null = null;
  let mountedElement: HTMLDivElement | undefined;
  const data = createMemo(() =>
    props.data
      .filter(
        (datum) =>
          datum.date instanceof Date &&
          Number.isFinite(datum.date.getTime()) &&
          Number.isFinite(datum.value),
      )
      .slice()
      .sort((a, b) => a.date.getTime() - b.date.getTime()),
  );
  const series = createMemo(
    () =>
      props.series ??
      buildSeriesMeta([...new Set(data().map((datum) => datum.series))]),
  );
  const seriesKeys = createMemo(() => series().map((s) => s.key));
  const colorRange = createMemo(() => series().map((s) => s.color));
  const multiSeries = createMemo(() => seriesKeys().length > 1);

  const definition = createMemo(() => {
    const d = data();
    const t = props.type ?? "line";
    const showPts = props.showPoints ?? d.length === 1;
    const c =
      (props.curveStyle ?? "smooth") === "smooth" ? smoothCurve : undefined;
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
      margin: props.showAxes === false ? 4 : undefined,
      x:
        t === "bar"
          ? {
              scale: () => scaleBand().padding(0.28),
              axis:
                props.showAxes === false ? false : { ticks: { format: dFmt } },
            }
          : {
              scale: scaleUtc,
              nice: false,
              axis:
                props.showAxes === false ? false : { ticks: { format: dFmt } },
            },
      y: {
        scale: scaleLinear,
        nice: true,
        grid: props.showAxes !== false,
        axis: props.showAxes === false ? false : { ticks: { format: vFmt } },
      },
      color: {
        scale: scaleOrdinal(keys, colors),
      },
      focus: "group-x",
      tooltip:
        props.showTooltip === false
          ? undefined
          : {
              use: tooltip,
              portal,
              className:
                "[--ts-chart-tooltip-background:var(--popover)] [--ts-chart-tooltip-color:var(--popover-foreground)] [--ts-chart-tooltip-border:1px_solid_var(--border)] [--ts-chart-tooltip-border-radius:0.625rem]",
              anchor: "group-center",
              placement: ["top", "right", "bottom", "left"],
            },
    });
  });

  const renderer = createSvgChartRenderer<ChartDatum, Date, number>(
    renderChartSvg,
  );
  const idPrefix = `ts-chart-${Math.random().toString(36).slice(2, 8)}`;
  const destroy = () => {
    adapterInstance?.destroy();
    adapterInstance = null;
    mountedElement = undefined;
  };
  createEffect(
    () => ({
      element: surface(),
      definition: definition(),
      height: props.height ?? 320,
      label: props.ariaLabel,
      populated: data().length > 0,
    }),
    ({ element, definition, height, label, populated }) => {
      if (!populated || !element) {
        destroy();
        setTooltipState(null);
        return;
      }
      const options: HostOptions = {
        renderer,
        definition,
        ariaLabel: label,
        height: Math.max(24, height),
        initialWidth: 320,
        idPrefix,
        onTooltipBodyChange: setTooltipState,
      };
      if (mountedElement !== element) {
        destroy();
        adapterInstance = createChartRendererAdapter(options);
        adapterInstance.mount(element);
        mountedElement = element;
      } else adapterInstance?.update(options);
    },
  );
  onCleanup(destroy);

  const showTotal = () => series().length > 1;

  return (
    <Show
      when={data().length > 0}
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
        class={cn("ts-chart-host min-w-0 text-muted-foreground", props.class)}
        style={{
          position: "relative",
          width: "100%",
          height: `${props.height ?? 320}px`,
        }}
      >
        <div
          ref={setSurface}
          class="ts-chart-surface [&_svg:focus:not(:focus-visible)]:outline-none [&_svg:focus-visible]:outline-2 [&_svg:focus-visible]:outline-offset-2 [&_svg:focus-visible]:outline-ring"
          style={{ width: "100%", height: "100%" }}
        />

        <Show when={tooltipState()}>
          {(target) => {
            const points = () => target().points ?? [];
            const rows = () =>
              points().map((p) => {
                const meta = series().find((s) => s.key === p.datum.series);
                return {
                  key: p.datum.series,
                  label: meta?.label ?? p.datum.series,
                  color: meta?.color ?? "#94a3b8",
                  value: p.yValue,
                  synthetic: Boolean(p.datum.synthetic),
                };
              });
            const total = () => rows().reduce((sum, r) => sum + r.value, 0);

            return (
              <Portal mount={target().element}>
                <div class="min-w-36 px-1 py-0.5">
                  <p class="mb-1.5 text-[11px] font-medium text-muted-foreground">
                    Período: {points()[0] ? dateFmt()(points()[0].xValue) : ""}
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
