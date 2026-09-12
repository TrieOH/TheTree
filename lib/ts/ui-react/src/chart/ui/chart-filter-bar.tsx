import { useState, type ReactElement, type SVGProps } from 'react'
import type { ChartType, CustomRange, RangeKey, SeriesMeta } from '../model/types'
import { PRESETS } from '../hooks/use-chart-filters'
import { isSameDay } from '../utils/date-utils'
import { AreaIcon, BarIcon, CalendarIcon, ChevronDownIcon, CloseIcon, LineIcon, ScatterIcon, SearchIcon } from './Icons'
import { DateRangeModal } from './date-range-modal'

const ALL_TYPE_OPTIONS: { key: ChartType; label: string; Icon: (props: SVGProps<SVGSVGElement>) => ReactElement }[] = [
  { key: 'line', label: 'Linha', Icon: LineIcon },
  { key: 'area', label: 'Área', Icon: AreaIcon },
  { key: 'bar', label: 'Barra', Icon: BarIcon },
  { key: 'scatter', label: 'Dispersão', Icon: ScatterIcon },
]

const defaultTriggerFormatter = (date: Date) =>
  new Intl.DateTimeFormat('pt-BR', { day: '2-digit', month: 'short' }).format(date)

/**
 * Segmented pill switch. The active option is highlighted by a single
 * absolutely-positioned pill that slides via CSS transform — the detail
 * that makes the control feel alive instead of just a row of buttons.
 */
export function SegmentedSwitch<T extends string>({
  options,
  value,
  onChange,
}: {
  options: { key: T; label: string; Icon?: (props: SVGProps<SVGSVGElement>) => ReactElement }[]
  value: T
  onChange: (key: T) => void
}) {
  const activeIndex = Math.max(
    0,
    options.findIndex(o => o.key === value),
  )

  return (
    <div className="relative inline-grid grid-flow-col auto-cols-fr rounded-full bg-muted p-1">
      <div
        aria-hidden
        className="absolute inset-y-1 rounded-full bg-card shadow-sm ring-1 ring-border transition-transform duration-300 ease-out"
        style={{
          width: `calc(${100 / options.length}% - 4px)`,
          transform: `translateX(calc(${activeIndex} * (100% + 4px) + 2px))`,
        }}
      />
      {options.map(({ key, label, Icon }) => (
        <button
          key={key}
          type="button"
          onClick={() => onChange(key)}
          aria-pressed={value === key}
          className={`relative z-10 flex items-center justify-center gap-1.5 rounded-full px-2.5 py-1.5 text-xs font-medium transition-colors sm:px-3.5 ${value === key
            ? 'text-foreground'
            : 'text-muted-foreground hover:text-foreground'
            }`}
        >
          {Icon && <Icon className="h-3.5 w-3.5 shrink-0" />}
          <span className={Icon ? 'hidden sm:inline' : ''}>{label}</span>
        </button>
      ))}
    </div>
  )
}

/** Chart-type switch. Renders nothing when 0-1 types are allowed (locked chart). */
export function ChartTypeSwitch({
  type,
  onTypeChange,
  allowedTypes,
}: {
  type: ChartType
  onTypeChange: (type: ChartType) => void
  allowedTypes?: ChartType[]
}) {
  const options = ALL_TYPE_OPTIONS.filter(o => !allowedTypes || allowedTypes.includes(o.key))
  if (options.length <= 1) return null
  return <SegmentedSwitch options={options} value={type} onChange={onTypeChange} />
}

/** Collapsible search input for filtering series by name. */
export function SearchToggle({ query, onQueryChange }: { query: string; onQueryChange: (query: string) => void }) {
  const [open, setOpen] = useState(false)
  return (
    <div
      className={`flex h-8 items-center overflow-hidden rounded-full border border-border transition-all duration-300 ${open ? 'w-40 px-1' : 'w-8'
        }`}
    >
      <button
        type="button"
        onClick={() => {
          if (open) {
            onQueryChange('')
            setOpen(false)
          } else {
            setOpen(true)
          }
        }}
        aria-label={open ? 'Fechar busca' : 'Buscar série'}
        className="flex h-8 w-8 shrink-0 items-center justify-center text-muted-foreground hover:text-foreground"
      >
        {open ? <CloseIcon className="h-3.5 w-3.5" /> : <SearchIcon className="h-3.5 w-3.5" />}
      </button>
      {open && (
        <input
          autoFocus
          value={query}
          onChange={e => onQueryChange(e.target.value)}
          placeholder="Filtrar série..."
          className="h-8 w-full border-0 bg-transparent text-xs text-foreground outline-none placeholder:text-muted-foreground"
        />
      )}
    </div>
  )
}

/** Date-range trigger button — shows the current period as text and opens DateRangeModal. */
export function DateRangeControl({
  range,
  onRangeChange,
  customRange,
  onCustomRangeChange,
  dateFormatter,
}: {
  range: RangeKey
  onRangeChange: (range: RangeKey) => void
  customRange: CustomRange
  onCustomRangeChange: (customRange: CustomRange) => void
  dateFormatter?: (date: Date) => string
}) {
  const [open, setOpen] = useState(false)
  const formatDate = dateFormatter ?? defaultTriggerFormatter

  const label =
    range === 'custom' && customRange.from
      ? customRange.to && !isSameDay(customRange.from, customRange.to)
        ? `${formatDate(customRange.from)} – ${formatDate(customRange.to)}`
        : formatDate(customRange.from)
      : (PRESETS.find(p => p.key === range)?.label ?? 'Período')

  return (
    <>
      <button
        type="button"
        onClick={() => setOpen(true)}
        className="flex h-8 items-center gap-1.5 rounded-full border border-border bg-background px-3 text-xs font-medium text-muted-foreground shadow-sm transition-colors hover:bg-muted hover:text-foreground"
      >
        <CalendarIcon className="h-3.5 w-3.5 text-muted-foreground" />
        <span className="max-w-[9rem] truncate sm:max-w-none">{label}</span>
        <ChevronDownIcon className="h-3 w-3 text-muted-foreground" />
      </button>
      <DateRangeModal
        open={open}
        onClose={() => setOpen(false)}
        range={range}
        customRange={customRange}
        onApply={(nextRange, nextCustomRange) => {
          onRangeChange(nextRange)
          onCustomRangeChange(nextCustomRange)
        }}
        dateFormatter={dateFormatter}
      />
    </>
  )
}

/** Series chips — doubles as legend and visibility filter. Renders nothing for a single series. */
export function SeriesFilterRow({
  series,
  visible,
  onToggleSeries,
  onShowAllSeries,
}: {
  series: SeriesMeta[]
  visible: Set<string>
  onToggleSeries: (key: string) => void
  onShowAllSeries: () => void
}) {
  if (series.length <= 1) return null
  const someHidden = visible.size < series.length

  return (
    <div className="flex items-center gap-2">
      <div className="flex flex-1 items-center gap-1.5 overflow-x-auto [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
        {series.map(s => {
          const active = visible.has(s.key)
          return (
            <button
              key={s.key}
              type="button"
              onClick={() => onToggleSeries(s.key)}
              aria-pressed={active}
              className={`flex shrink-0 items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-medium transition-all duration-200 ${active
                ? 'border-border bg-background text-foreground shadow-sm'
                : 'border-transparent bg-muted text-muted-foreground opacity-60'
                }`}
            >
              <span
                className="h-2 w-2 shrink-0 rounded-full transition-transform duration-200"
                style={{ backgroundColor: s.color, transform: active ? 'scale(1)' : 'scale(0.65)' }}
              />
              {s.label}
            </button>
          )
        })}
      </div>
      {someHidden && (
        <button
          type="button"
          onClick={onShowAllSeries}
          className="shrink-0 text-xs font-medium text-muted-foreground underline-offset-2 hover:text-foreground hover:underline"
        >
          Mostrar tudo
        </button>
      )}
    </div>
  )
}

interface ChartFilterBarProps {
  type: ChartType
  onTypeChange: (type: ChartType) => void
  allowedTypes?: ChartType[]
  range: RangeKey
  onRangeChange: (range: RangeKey) => void
  customRange: CustomRange
  onCustomRangeChange: (customRange: CustomRange) => void
  showRangeFilter?: boolean
  dateFormatter?: (date: Date) => string
  series: SeriesMeta[]
  visible: Set<string>
  onToggleSeries: (key: string) => void
  onShowAllSeries: () => void
  showSeriesFilter?: boolean
  query: string
  onQueryChange: (query: string) => void
  showSearch?: boolean
}

/**
 * Convenience composition of all controls in one dedicated block (type +
 * range + search on top, series chips below). `ChartCard` does NOT use this
 * by default — it packs the individual pieces (`ChartTypeSwitch`,
 * `DateRangeControl`, `SearchToggle`, `SeriesFilterRow`) directly into its
 * header row instead, so a single active filter doesn't get a whole row to
 * itself. Use this component if you want that dedicated-row look, or when
 * composing a custom layout that isn't `ChartCard`.
 */
export function ChartFilterBar({
  type,
  onTypeChange,
  allowedTypes,
  range,
  onRangeChange,
  customRange,
  onCustomRangeChange,
  showRangeFilter = true,
  dateFormatter,
  series,
  visible,
  onToggleSeries,
  onShowAllSeries,
  showSeriesFilter = true,
  query,
  onQueryChange,
  showSearch = true,
}: ChartFilterBarProps) {
  return (
    <div className="flex flex-col gap-3">
      <div className="flex flex-wrap items-center justify-between gap-2.5">
        <ChartTypeSwitch type={type} onTypeChange={onTypeChange} allowedTypes={allowedTypes} />
        <div className="flex flex-wrap items-center gap-2">
          {showSearch && <SearchToggle query={query} onQueryChange={onQueryChange} />}
          {showRangeFilter && (
            <DateRangeControl
              range={range}
              onRangeChange={onRangeChange}
              customRange={customRange}
              onCustomRangeChange={onCustomRangeChange}
              dateFormatter={dateFormatter}
            />
          )}
        </div>
      </div>
      {showSeriesFilter && (
        <SeriesFilterRow
          series={series}
          visible={visible}
          onToggleSeries={onToggleSeries}
          onShowAllSeries={onShowAllSeries}
        />
      )}
    </div>
  )
}
