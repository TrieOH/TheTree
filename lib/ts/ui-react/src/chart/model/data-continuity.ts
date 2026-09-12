import type { ChartDatum } from './types'
import { addDays, isSameDay, startOfDay } from '../utils/date-utils'

export type ContinuityMode = 'continuous' | 'none'

interface FillContinuousSeriesOptions {
  /** Rows for the series in question, already filtered to the visible series/search — but NOT date-bounded, so points before the window can seed the carried-forward baseline. */
  data: ChartDatum[]
  seriesKeys: string[]
  /** Inclusive lower bound of the window to render. `null` = no fixed window — each series gets a single 0-value anchor the day before its own first real point (a visible "rise from zero"), instead of a full day-by-day pad back to some arbitrary date. */
  windowStart: Date | null
  /** Inclusive upper bound of the window to render. `null` = end at each series' own last point (no trailing pad). */
  windowEnd: Date | null
  mode: ContinuityMode
}

/**
 * Produces exactly one row per day per series across [windowStart, windowEnd]:
 *
 * - A day with real data keeps that row untouched.
 * - A day with NO data, but with an earlier real point somewhere before it,
 *   carries that last known value forward (`synthetic: true`) — so a metric
 *   that simply hasn't changed reads as a flat continuation, not a drop to
 *   zero or a broken line.
 * - With an explicit `windowStart` (e.g. the "últimos 90 dias" preset), every
 *   day from `windowStart` up to the series' first real point is padded with
 *   0 (`synthetic: true`) — so the chart shows exactly when the metric began
 *   relative to the selected window, instead of silently starting mid-air.
 * - With `windowStart: null` (e.g. "Tudo", where there's no fixed lower
 *   bound to pad back to), each series still gets a single 0-value anchor
 *   one day before its own first real point — enough for the line to
 *   visibly rise from zero instead of appearing to start already elevated,
 *   without fabricating a potentially huge run of daily zeros back to an
 *   arbitrary date.
 * - The trailing edge is carried forward the same way up to `windowEnd`
 *   (e.g. "últimos 7 dias" with data only up to 2 days ago still reaches
 *   today on the chart).
 *
 * Only meaningful for continuous marks (line/area) — bar and scatter should
 * pass `mode: 'none'` (a missing bar/point is a normal, expected gap there).
 */
export function fillContinuousSeries({
  data,
  seriesKeys,
  windowStart,
  windowEnd,
  mode,
}: FillContinuousSeriesOptions): ChartDatum[] {
  if (mode === 'none' || data.length === 0) return data

  const bySeries = new Map<string, ChartDatum[]>()
  for (const key of seriesKeys) bySeries.set(key, [])
  for (const datum of data) {
    bySeries.get(datum.series)?.push(datum)
  }

  const result: ChartDatum[] = []

  for (const key of seriesKeys) {
    const points = (bySeries.get(key) ?? []).slice().sort((a, b) => a.date.getTime() - b.date.getTime())
    if (points.length === 0) continue

    // With no fixed window, don't pad every day back to the dawn of time —
    // just add one anchor day before the series actually starts, so the
    // line has somewhere to rise from.
    const from = windowStart ?? addDays(points[0].date, -1)
    const to = windowEnd ?? points[points.length - 1].date
    if (from.getTime() > to.getTime()) continue

    let pointIndex = 0
    let lastKnownValue: number | null = null
    const cursor = startOfDay(from)
    const end = startOfDay(to)

    // Consume any points strictly before the window to seed the
    // carried-forward baseline right from windowStart, instead of treating
    // an earlier-starting series as if it had no data yet.
    while (pointIndex < points.length && startOfDay(points[pointIndex].date).getTime() < cursor.getTime()) {
      lastKnownValue = points[pointIndex].value
      pointIndex++
    }

    while (cursor.getTime() <= end.getTime()) {
      let matchedToday = false
      while (pointIndex < points.length && isSameDay(points[pointIndex].date, cursor)) {
        lastKnownValue = points[pointIndex].value
        result.push(points[pointIndex])
        pointIndex++
        matchedToday = true
      }
      if (!matchedToday) {
        result.push({
          date: new Date(cursor),
          value: lastKnownValue ?? 0,
          series: key,
          synthetic: true,
        })
      }
      cursor.setDate(cursor.getDate() + 1)
    }
  }

  return result.sort((a, b) => a.date.getTime() - b.date.getTime())
}