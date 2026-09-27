import { useCallback, useEffect, useMemo, useState } from 'react'
import type { ChartDatum, ChartType, CustomRange, RangeKey, RangeOption } from '../model/types'

export const PRESETS: RangeOption[] = [
  { key: '7d', label: 'Últimos 7 dias', days: 7 },
  { key: '30d', label: 'Últimos 30 dias', days: 30 },
  { key: '90d', label: 'Últimos 90 dias', days: 90 },
  { key: 'all', label: 'Tudo', days: null },
]

interface UseChartFiltersOptions {
  data: ChartDatum[]
  seriesKeys: string[]
  initialType?: ChartType
  initialRange?: RangeKey
}

interface Window {
  windowStart: Date | null
  windowEnd: Date | null
}

export function computeRangeWindow(range: RangeKey, customRange: CustomRange, now = new Date()): Window {
  if (range === 'custom') {
    return { windowStart: customRange.from, windowEnd: customRange.to }
  }
  if (range === 'all') {
    return { windowStart: null, windowEnd: null }
  }
  const preset = PRESETS.find(p => p.key === range)
  if (!preset?.days) return { windowStart: null, windowEnd: null }
  return { windowStart: new Date(now.getTime() - preset.days * 86_400_000), windowEnd: now }
}

export function useChartFilters({
  data,
  seriesKeys,
  initialType = 'line',
  initialRange = 'all',
}: UseChartFiltersOptions) {
  const [type, setType] = useState<ChartType>(initialType)
  const [range, setRange] = useState<RangeKey>(initialRange)
  const [customRange, setCustomRange] = useState<CustomRange>({ from: null, to: null })
  const [visible, setVisible] = useState<Set<string>>(() => new Set(seriesKeys))
  const [query, setQuery] = useState('')

  useEffect(() => {
    setVisible(current => (current.size === 0 && seriesKeys.length > 0 ? new Set(seriesKeys) : current))
  }, [seriesKeys])

  const toggleSeries = useCallback((key: string) => {
    setVisible(prev => {
      if (prev.has(key)) {
        if (prev.size === 1) return prev // always keep at least one series visible
        const next = new Set(prev)
        next.delete(key)
        return next
      }
      return new Set(prev).add(key)
    })
  }, [])

  const isolateSeries = useCallback((key: string) => {
    setVisible(new Set([key]))
  }, [])

  const showAllSeries = useCallback(() => {
    setVisible(new Set(seriesKeys))
  }, [seriesKeys])

  const { windowStart, windowEnd } = useMemo(() => computeRangeWindow(range, customRange), [range, customRange])

  // Series/search filtered, but NOT date-bounded — this is what
  // fillContinuousSeries needs to see points before windowStart in order to
  // seed the carried-forward baseline correctly.
  const seriesFilteredData = useMemo(() => {
    const q = query.trim().toLowerCase()
    return data.filter(d => visible.has(d.series) && (!q || d.series.toLowerCase().includes(q)))
  }, [data, visible, query])

  // Series/search + date bounded — what non-continuous chart types (bar,
  // scatter) render directly, with no gap-filling.
  const filteredData = useMemo(() => {
    return seriesFilteredData.filter(d => {
      const t = d.date.getTime()
      if (windowStart !== null && t < windowStart.getTime()) return false
      if (windowEnd !== null && t > windowEnd.getTime()) return false
      return true
    })
  }, [seriesFilteredData, windowStart, windowEnd])

  return {
    type,
    setType,
    range,
    setRange,
    customRange,
    setCustomRange,
    visible,
    toggleSeries,
    isolateSeries,
    showAllSeries,
    query,
    setQuery,
    windowStart,
    windowEnd,
    seriesFilteredData,
    filteredData,
  }
}
