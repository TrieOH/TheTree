import { useId, useMemo } from 'react'
import { defineChart, lineY, areaY, barY, dot, group } from '@tanstack/charts'
import { scaleLinear } from '@tanstack/charts/scales/linear'
import { scaleBand } from '@tanstack/charts/scales/band'
import { scaleOrdinal } from '@tanstack/charts/scales/ordinal'
import { d3Curve } from '@tanstack/charts/d3/shape'
import { tooltip } from '@tanstack/charts/tooltip'
import { portal } from '@tanstack/charts/tooltip/portal'
import { Chart } from '@tanstack/charts/react/tooltip'
import { scaleUtc } from 'd3-scale'
import { curveMonotoneX } from 'd3-shape'
import type { ChartDatum, ChartType, CurveStyle, SeriesMeta } from '../model/types'

const smoothCurve = d3Curve(curveMonotoneX)

const xAccessor = (d: ChartDatum): Date => d.date
const yAccessor = (d: ChartDatum): number => d.value
const zAccessor = (d: ChartDatum): string => d.series
const colorAccessor = (d: ChartDatum): string => d.series

interface GenericChartProps {
  data: ChartDatum[]
  type: ChartType
  series: SeriesMeta[]
  ariaLabel: string
  height?: number
  showPoints?: boolean
  curveStyle?: CurveStyle
  valueFormatter?: (value: number) => string
  dateFormatter?: (date: Date) => string
  tooltipDetails?: (datum: ChartDatum) => { label: string; value: string }[]
  emptyMessage?: string
}

const defaultValueFormatter = (value: number) =>
  new Intl.NumberFormat('pt-BR', { notation: 'compact', maximumFractionDigits: 1 }).format(value)

const defaultDateFormatter = (date: Date) =>
  new Intl.DateTimeFormat('pt-BR', { day: '2-digit', month: 'short' }).format(date)

/**
 * Renders one of four chart types (line / area / bar / scatter) from the
 * same tidy dataset. Swapping `type` swaps the mark, not the data shape —
 * that's what makes this "generic". Width is left unset so the chart
 * measures its container and stays responsive on mobile.
 */
export function GenericChart({
  data,
  type,
  series,
  ariaLabel,
  height = 320,
  showPoints = false,
  curveStyle = 'smooth',
  valueFormatter = defaultValueFormatter,
  dateFormatter = defaultDateFormatter,
  tooltipDetails,
  emptyMessage = 'Nenhum dado para os filtros selecionados.',
}: GenericChartProps) {
  const idPrefix = useId()
  const seriesKeys = useMemo(() => series.map(s => s.key), [series])
  const colorRange = useMemo(() => series.map(s => s.color), [series])
  const multiSeries = seriesKeys.length > 1
  const curve = curveStyle === 'smooth' ? smoothCurve : undefined

  const definition = useMemo(() => {
    const marks = (() => {
      switch (type) {
        case 'area':
          return [
            areaY(data, {
              x: xAccessor,
              y1: 0,
              y2: yAccessor,
              z: zAccessor,
              color: colorAccessor,
              fillOpacity: 0.14,
              curve,
            }),
            lineY(data, {
              x: xAccessor,
              y: yAccessor,
              z: zAccessor,
              color: colorAccessor,
              points: showPoints,
              strokeWidth: 2,
              curve,
            }),
          ]
        case 'bar':
          return [
            barY(data, {
              x: xAccessor,
              y: yAccessor,
              z: zAccessor,
              color: colorAccessor,
              radius: 4,
              layout: multiSeries ? group({ padding: 0.24 }) : undefined,
            }),
          ]
        case 'scatter':
          return [
            dot(data, {
              x: xAccessor,
              y: yAccessor,
              z: zAccessor,
              color: colorAccessor,
              r: 4,
              // Grows the hovered/focused point so pointer feedback reads
              // clearly before the tooltip text even lands.
              states: [
                {
                  when: { focus: 'primary' },
                  style: { r: 7 },
                  transition: { type: 'tween', duration: 140, easing: 'ease-out' },
                },
              ],
            }),
          ]
        case 'line':
        default:
          return [
            lineY(data, {
              x: xAccessor,
              y: yAccessor,
              z: zAccessor,
              color: colorAccessor,
              points: showPoints,
              strokeWidth: 2.25,
              curve,
            }),
          ]
      }
    })()

    return defineChart({
      marks,
      x:
        type === 'bar'
          ? {
            scale: () => scaleBand().padding(0.28),
            axis: { ticks: { format: dateFormatter } },
          }
          : {
            scale: scaleUtc,
            nice: false,
            axis: { ticks: { format: dateFormatter } },
          },
      y: {
        scale: scaleLinear,
        nice: true,
        grid: true,
        axis: { ticks: { format: valueFormatter } },
      },
      color: {
        scale: scaleOrdinal(seriesKeys, colorRange),
      },
      focus: 'group-x',
      tooltip: {
        use: tooltip,
        portal,
        className:
          '[--ts-chart-tooltip-background:var(--popover)] [--ts-chart-tooltip-color:var(--popover-foreground)] [--ts-chart-tooltip-border:1px_solid_var(--border)] [--ts-chart-tooltip-border-radius:0.625rem]',
        anchor: 'group-center',
        placement: ['top', 'right', 'bottom', 'left'],
      },
    })
  }, [data, type, seriesKeys, colorRange, multiSeries, showPoints, curve, dateFormatter, valueFormatter])

  if (data.length === 0) {
    return (
      <div
        className="flex items-center justify-center rounded-xl border border-dashed border-border text-sm text-muted-foreground"
        style={{ height }}
      >
        {emptyMessage}
      </div>
    )
  }

  const showTotal = series.length > 1

  return (
    <Chart
      idPrefix={idPrefix}
      definition={definition}
      ariaLabel={ariaLabel}
      height={height}
      className="text-muted-foreground"
      renderTooltipBody={({ points }) => {
        const rows = points.map(p => {
          const meta = series.find(s => s.key === p.datum.series)
          return {
            key: p.datum.series,
            label: meta?.label ?? p.datum.series,
            color: meta?.color ?? '#94a3b8',
            value: p.yValue,
            synthetic: Boolean(p.datum.synthetic),
          }
        })
        const total = rows.reduce((sum, r) => sum + r.value, 0)

        return (
          <div className="min-w-[9rem] px-1 py-0.5">
            <p className="mb-1.5 text-[11px] font-medium text-muted-foreground">
              Período: {points[0] ? dateFormatter(points[0].xValue) : ''}
            </p>
            <ul className="space-y-1">
              {rows.map(r => (
                <li key={r.key} className="flex items-center justify-between gap-4 text-xs">
                  <span className="flex items-center gap-1.5 text-muted-foreground">
                    <span className="h-1.5 w-1.5 shrink-0 rounded-full" style={{ backgroundColor: r.color }} />
                    {r.label}
                    {r.synthetic && (
                      <span className="text-muted-foreground/50" title="Valor mantido — sem novo dado neste dia">
                        ⋯
                      </span>
                    )}
                  </span>
                  <span className="font-mono font-medium tabular-nums text-foreground">
                    {valueFormatter(r.value)}
                  </span>
                </li>
              ))}
            </ul>
            {points[0] && tooltipDetails?.(points[0].datum).map(detail => (
              <div key={detail.label} className="mt-1 flex items-center justify-between gap-4 text-xs">
                <span className="text-muted-foreground">{detail.label}</span>
                <span className="font-mono font-medium tabular-nums text-foreground">{detail.value}</span>
              </div>
            ))}
            {showTotal && (
              <div className="mt-1.5 flex items-center justify-between gap-4 border-t border-border pt-1.5 text-xs">
                <span className="font-medium text-muted-foreground">Total</span>
                <span className="font-mono font-semibold tabular-nums text-foreground">
                  {valueFormatter(total)}
                </span>
              </div>
            )}
          </div>
        )
      }}
    />
  )
}
