import { useEffect, useMemo, useState } from 'react'
import { Combobox, type ComboboxOption } from '../../toolbar-combobox'
import type { CustomRange, RangeKey } from '../model/types'
import { computeRangeWindow, PRESETS } from '../hooks/use-chart-filters'
import { addMonths, buildCalendarWeeks, isSameDay, isWithinRange } from '../utils/date-utils'
import { ChevronLeftIcon, ChevronRightIcon, CloseIcon } from './Icons'

const WEEKDAY_LABELS = ['D', 'S', 'T', 'Q', 'Q', 'S', 'S']

const MONTH_LABELS = Array.from({ length: 12 }, (_, i) =>
  new Intl.DateTimeFormat('pt-BR', { month: 'long' }).format(new Date(2000, i, 1)),
)
const MONTH_OPTIONS: ComboboxOption[] = MONTH_LABELS.map((label, value) => ({
  value: String(value),
  label: label.charAt(0).toLocaleUpperCase('pt-BR') + label.slice(1),
}))

const defaultLabelFormatter = (date: Date) =>
  new Intl.DateTimeFormat('pt-BR', { day: '2-digit', month: 'short', year: 'numeric' }).format(date)

interface DateRangeModalProps {
  open: boolean
  onClose: () => void
  range: RangeKey
  customRange: CustomRange
  onApply: (range: RangeKey, customRange: CustomRange) => void
  dateFormatter?: (date: Date) => string
}

/**
 * Opens as a centered modal (not inline) with a presets column on the left
 * and a single-month calendar on the right. Selecting two days highlights
 * the range between them as a continuous pill, matching common
 * product date pickers, rather than isolated circles. Nothing commits to
 * the parent until "Aplicar" — "Cancelar" or the backdrop discards the
 * draft.
 */
export function DateRangeModal({ open, onClose, range, customRange, onApply, dateFormatter }: DateRangeModalProps) {
  const formatLabel = dateFormatter ?? defaultLabelFormatter

  const [draftPreset, setDraftPreset] = useState<RangeKey>(range)
  const [draftFrom, setDraftFrom] = useState<Date | null>(customRange.from)
  const [draftTo, setDraftTo] = useState<Date | null>(customRange.to)
  const [visibleMonth, setVisibleMonth] = useState<Date>(customRange.to ?? customRange.from ?? new Date())
  const [hoverDate, setHoverDate] = useState<Date | null>(null)
  const [mounted, setMounted] = useState(false)

  // Reset the draft to the committed values every time the modal opens, and
  // drive a small mount transition (fade + scale in) rather than popping in.
  useEffect(() => {
    if (!open) {
      setMounted(false)
      return
    }
    const { windowStart, windowEnd } = computeRangeWindow(range, customRange)
    setDraftPreset(range)
    setDraftFrom(windowStart)
    setDraftTo(windowEnd)
    setVisibleMonth(windowEnd ?? windowStart ?? new Date())
    const raf = requestAnimationFrame(() => setMounted(true))
    return () => cancelAnimationFrame(raf)
  }, [open, range, customRange])

  useEffect(() => {
    if (!open) return
    function onKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [open, onClose])

  const yearOptions = useMemo(() => {
    const visibleYear = visibleMonth.getFullYear()
    return Array.from({ length: 9 }, (_, i) => {
      const year = String(visibleYear - 7 + i)
      return { value: year, label: year }
    })
  }, [visibleMonth])

  if (!open) return null

  const today = new Date()
  const weeks = buildCalendarWeeks(visibleMonth)

  const rangeFrom = draftFrom
  const rangeTo = draftTo ?? (draftFrom && hoverDate ? hoverDate : null)
  const isPreviewOnly = draftFrom !== null && draftTo === null

  function handleDayClick(day: Date) {
    setDraftPreset('custom')
    if (!draftFrom || draftTo) {
      setDraftFrom(day)
      setDraftTo(null)
      return
    }
    if (day.getTime() < draftFrom.getTime()) {
      setDraftTo(draftFrom)
      setDraftFrom(day)
    } else {
      setDraftTo(day)
    }
  }

  function handlePresetClick(key: RangeKey) {
    const { windowStart, windowEnd } = computeRangeWindow(key, { from: null, to: null })
    setDraftPreset(key)
    setDraftFrom(windowStart)
    setDraftTo(windowEnd)
    setVisibleMonth(windowEnd ?? windowStart ?? new Date())
  }

  function handleApply() {
    if (draftPreset === 'custom') {
      if (!draftFrom) return
      onApply('custom', { from: draftFrom, to: draftTo ?? draftFrom })
    } else {
      onApply(draftPreset, { from: null, to: null })
    }
    onClose()
  }

  const summaryLabel =
    draftPreset === 'custom'
      ? draftFrom
        ? `${formatLabel(draftFrom)}${draftTo ? ` – ${formatLabel(draftTo)}` : ' – selecione o fim'}`
        : 'Selecione a data inicial'
      : draftFrom && draftTo
        ? `${formatLabel(draftFrom)} – ${formatLabel(draftTo)}`
        : (PRESETS.find(p => p.key === draftPreset)?.label ?? '')

  return (
    <div
      className={`fixed inset-0 z-60 flex items-center justify-center bg-foreground/40 p-2 backdrop-blur-sm transition-opacity duration-150 sm:p-4 ${mounted ? 'opacity-100' : 'opacity-0'
        }`}
      onClick={onClose}
      role="presentation"
    >
      <div
        onClick={event => event.stopPropagation()}
        role="dialog"
        aria-modal="true"
        aria-label="Selecionar período"
        style={{ maxHeight: 'calc(100vh - 2rem)' }}
        className={`flex min-w-[300px] max-w-sm flex-1 flex-col overflow-hidden rounded-lg border border-border bg-card text-card-foreground shadow-2xl transition-all duration-150 sm:max-w-2xl ${mounted ? 'scale-100 opacity-100' : 'scale-95 opacity-0'
          }`}
      >
        <div className="flex shrink-0 items-center justify-between border-b border-border px-4 py-3">
          <div>
            <p className="text-sm font-semibold">Período personalizado</p>
            <p className="hidden text-xs text-muted-foreground sm:block">Escolha o intervalo que deseja visualizar.</p>
          </div>
          <button
            type="button"
            onClick={onClose}
            aria-label="Fechar"
            className="rounded-lg p-1 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
          >
            <CloseIcon className="h-4 w-4" />
          </button>
        </div>

        <div className="flex min-h-0 flex-1 flex-col overflow-y-auto">
          {/* Presets */}
          <div className="grid shrink-0 grid-cols-2 gap-1.5 border-b border-border p-3 sm:grid-cols-6">
            {PRESETS.map(p => (
              <button
                key={p.key}
                type="button"
                onClick={() => handlePresetClick(p.key)}
                className={`shrink-0 rounded-lg px-3 py-2 text-center text-xs font-medium transition-colors ${draftPreset === p.key
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                  }`}
              >
                {p.label}
              </button>
            ))}
            <button
              type="button"
              onClick={() => handlePresetClick('custom')}
              className={`shrink-0 rounded-lg px-3 py-2 text-center text-xs font-medium transition-colors ${draftPreset === 'custom'
                ? 'bg-primary text-primary-foreground'
                : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                }`}
            >
              Personalizado
            </button>
          </div>

          {/* Calendar */}
          <div className="flex-1 p-4">
            <div className="mb-3 flex items-center justify-between gap-2">
              <button
                type="button"
                onClick={() => setVisibleMonth(m => addMonths(m, -1))}
                aria-label="Mês anterior"
                className="rounded-lg p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
              >
                <ChevronLeftIcon className="h-4 w-4" />
              </button>

              <div className="flex items-center gap-1.5">
                <Combobox
                  value={String(visibleMonth.getMonth())}
                  options={MONTH_OPTIONS}
                  placeholder="Mês"
                  searchPlaceholder="Buscar mês…"
                  className="w-28 sm:w-48"
                  triggerClassName="h-9 font-semibold"
                  dropdownClassName="w-60 max-w-none"
                  onChange={value => setVisibleMonth(new Date(visibleMonth.getFullYear(), Number(value), 1))}
                />
                <Combobox
                  value={String(visibleMonth.getFullYear())}
                  options={yearOptions}
                  placeholder="Ano"
                  searchPlaceholder="Buscar ano…"
                  className="w-20 sm:w-32"
                  triggerClassName="h-9 font-semibold"
                  dropdownClassName="w-40 max-w-none"
                  onChange={value => setVisibleMonth(new Date(Number(value), visibleMonth.getMonth(), 1))}
                />
              </div>

              <button
                type="button"
                onClick={() => setVisibleMonth(m => addMonths(m, 1))}
                aria-label="Próximo mês"
                className="rounded-lg p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
              >
                <ChevronRightIcon className="h-4 w-4" />
              </button>
            </div>

            <div className="grid grid-cols-7 gap-y-0.5">
              {WEEKDAY_LABELS.map((label, i) => (
                <span
                  key={i}
                  className="pb-1 text-center text-[10px] font-medium uppercase tracking-wide text-muted-foreground"
                >
                  {label}
                </span>
              ))}
              {weeks.flatMap((week, weekIndex) =>
                week.map((day, dayIndex) => {
                  if (!day) return <div key={`${weekIndex}-${dayIndex}`} />

                  const isStartCell = isSameDay(day, rangeFrom)
                  const isEndCell = rangeTo ? isSameDay(day, rangeTo) : false
                  const isEdgeCell = isStartCell || isEndCell
                  const inRange = Boolean(rangeFrom && rangeTo && isWithinRange(day, rangeFrom, rangeTo))
                  const isToday = isSameDay(day, today)
                  const roundLeft = dayIndex === 0 || isStartCell
                  const roundRight = dayIndex === 6 || isEndCell

                  return (
                    <div key={`${weekIndex}-${dayIndex}`} className="relative">
                      {inRange && (
                        <span
                          aria-hidden
                          className={`absolute inset-y-0.5 inset-x-0 ${isPreviewOnly ? 'bg-muted/50' : 'bg-primary/10'
                            } ${roundLeft ? 'rounded-l-full' : ''} ${roundRight ? 'rounded-r-full' : ''}`}
                        />
                      )}
                      <button
                        type="button"
                        onClick={() => handleDayClick(day)}
                        onMouseEnter={() => setHoverDate(day)}
                        className={`relative z-10 mx-auto flex h-8 w-8 items-center justify-center rounded-full text-xs transition-colors ${isEdgeCell
                          ? isPreviewOnly && isEndCell
                            ? 'font-semibold text-primary ring-2 ring-primary'
                            : 'bg-primary font-semibold text-primary-foreground'
                          : 'text-foreground hover:bg-muted'
                          } ${isToday && !isEdgeCell ? 'ring-1 ring-inset ring-border' : ''}`}
                      >
                        {day.getDate()}
                      </button>
                    </div>
                  )
                }),
              )}
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="flex shrink-0 items-center justify-between gap-3 border-t border-border px-4 py-3">
          <p className="min-w-0 truncate text-xs text-muted-foreground">{summaryLabel}</p>
          <div className="flex shrink-0 gap-2">
            <button
              type="button"
              onClick={onClose}
              className="rounded-lg px-3 py-1.5 text-xs font-medium text-muted-foreground hover:bg-muted hover:text-foreground"
            >
              Cancelar
            </button>
            <button
              type="button"
              onClick={handleApply}
              disabled={draftPreset === 'custom' && !draftFrom}
              className="rounded-lg bg-primary px-3.5 py-1.5 text-xs font-semibold text-primary-foreground transition-opacity hover:opacity-90 disabled:opacity-40"
            >
              Aplicar
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
