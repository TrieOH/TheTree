import { useMemo, useState } from 'react'
import type { ChartDatum, ChartType, CurveStyle, RangeKey } from '../model/types'
import { useChartFilters } from '../hooks/use-chart-filters'
import { buildSeriesMeta } from '../model/theme'
import { fillContinuousSeries } from '../model/data-continuity'
import { ChartTypeSwitch, DateRangeControl, SearchToggle, SeriesFilterRow } from './chart-filter-bar'
import { GenericChart } from './generic-chart'

interface ChartCardProps {
  title: string
  subtitle?: string
  data: ChartDatum[]
  seriesLabels?: Record<string, string>
  seriesColors?: Record<string, string>

  initialType?: ChartType
  /**
   * Restrict which chart types are selectable. Pass a single type
   * (e.g. `['bar']`) to lock the chart — the type switch is hidden and the
   * chart always renders that type. Defaults to all four types.
   */
  allowedTypes?: ChartType[]

  initialRange?: RangeKey
  showRangeFilter?: boolean
  showSeriesFilter?: boolean
  showSearchFilter?: boolean
  showPointsToggle?: boolean

  /**
   * For line/area charts, carries the last known value forward through
   * gaps and pads the leading edge with 0 up to the first real data point,
   * so the chart reads as continuous instead of dropping to zero, breaking,
   * or stopping short of the selected range's end. Has no effect on
   * bar/scatter (a missing bar or point there is a normal, expected gap).
   * Defaults to true.
   */
  continuity?: boolean

  curveStyle?: CurveStyle
  height?: number
  valueFormatter?: (value: number) => string
  dateFormatter?: (date: Date) => string
  tooltipDetails?: (datum: ChartDatum) => { label: string; value: string }[]
}

/**
 * Drop-in chart card: title + controls + chart, all wired together. Every
 * active control (points toggle, type switch, search, date range) is packed
 * into the same header row as the title — when only one is active it just
 * sits there instead of a mostly-empty row underneath, and on wide screens
 * nothing is left stranded far from the title with empty space in between.
 * Series chips get their own row below since they can wrap to multiple
 * lines with many series.
 */
export function ChartCard({
  title,
  subtitle,
  data,
  seriesLabels,
  seriesColors,
  initialType = 'line',
  allowedTypes,
  initialRange = 'all',
  showRangeFilter = true,
  showSeriesFilter = true,
  showSearchFilter = true,
  showPointsToggle = true,
  continuity = true,
  curveStyle = 'smooth',
  height = 320,
  valueFormatter,
  dateFormatter,
  tooltipDetails,
}: ChartCardProps) {
  const seriesKeys = useMemo(() => Array.from(new Set(data.map(d => d.series))), [data])
  const series = useMemo(
    () => buildSeriesMeta(seriesKeys, seriesLabels, seriesColors),
    [seriesKeys, seriesLabels, seriesColors],
  )
  const [showPoints, setShowPoints] = useState(false)

  const lockedType = allowedTypes && allowedTypes.length === 1 ? allowedTypes[0] : undefined

  const {
    type,
    setType,
    range,
    setRange,
    customRange,
    setCustomRange,
    visible,
    toggleSeries,
    showAllSeries,
    query,
    setQuery,
    windowStart,
    windowEnd,
    seriesFilteredData,
    filteredData,
  } = useChartFilters({
    data,
    seriesKeys,
    initialType: lockedType ?? initialType,
    initialRange,
  })

  const effectiveType = lockedType ?? type
  const visibleSeries = series.filter(s => visible.has(s.key))
  const canTogglePoints = showPointsToggle && (effectiveType === 'line' || effectiveType === 'area')
  const isContinuousType = effectiveType === 'line' || effectiveType === 'area'

  const chartData = useMemo(() => {
    if (!isContinuousType) return filteredData
    return fillContinuousSeries({
      data: seriesFilteredData,
      seriesKeys: visibleSeries.map(s => s.key),
      windowStart,
      windowEnd,
      mode: continuity ? 'continuous' : 'none',
    })
  }, [isContinuousType, filteredData, seriesFilteredData, visibleSeries, windowStart, windowEnd, continuity])

  return (
    <div className="w-full rounded-2xl border border-border bg-card p-4 text-card-foreground shadow-sm sm:p-5">
      <div className="flex flex-col gap-3 sm:flex-row! sm:items-start! sm:justify-between!">
        <div className="min-w-0 flex-1">
          <h3 className="text-sm font-semibold">{title}</h3>
          {subtitle && <p className="mt-0.5 text-xs text-muted-foreground">{subtitle}</p>}
        </div>

        <div className="flex w-full flex-wrap items-center justify-start gap-2 sm:w-auto! sm:justify-end!">
          {canTogglePoints && (
            <button
              type="button"
              onClick={() => setShowPoints(v => !v)}
              aria-pressed={showPoints}
              className={`shrink-0 rounded-full border px-2.5 py-1 text-[11px] font-medium transition-colors ${showPoints
                ? 'border-primary bg-primary text-primary-foreground'
                : 'border-border text-muted-foreground hover:bg-muted hover:text-foreground'
                }`}
            >
              Pontos
            </button>
          )}
          <ChartTypeSwitch type={effectiveType} onTypeChange={setType} allowedTypes={allowedTypes} />
          {showSearchFilter && <SearchToggle query={query} onQueryChange={setQuery} />}
          {showRangeFilter && (
            <DateRangeControl
              range={range}
              onRangeChange={setRange}
              customRange={customRange}
              onCustomRangeChange={setCustomRange}
              dateFormatter={dateFormatter}
            />
          )}
        </div>
      </div>

      {showSeriesFilter && (
        <div className="mt-3">
          <SeriesFilterRow
            series={series}
            visible={visible}
            onToggleSeries={toggleSeries}
            onShowAllSeries={showAllSeries}
          />
        </div>
      )}

      <div className="mt-4">
        <GenericChart
          data={chartData}
          type={effectiveType}
          series={visibleSeries}
          ariaLabel={title}
          height={height}
          showPoints={showPoints}
          curveStyle={curveStyle}
          valueFormatter={valueFormatter}
          dateFormatter={dateFormatter}
          tooltipDetails={tooltipDetails}
        />
      </div>
    </div>
  )
}
