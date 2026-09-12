/**
 * Shared types for the generic chart component set.
 *
 * Data model: "tidy" / long format. One row per (date, series) pair.
 * This is what lets a single dataset drive line, area, bar, or scatter
 * without reshaping — every mark type reads the same three fields.
 */

export type ChartType = 'line' | 'area' | 'bar' | 'scatter'

export interface ChartDatum {
  date: Date
  value: number
  series: string
  /** Set by fillContinuousSeries on rows it synthesized (gap-fill / leading pad) rather than rows from your original data. */
  synthetic?: boolean
  /** Extra fields are allowed and ignored by the chart (useful for tooltips/analytics). */
  [key: string]: unknown
}

export interface SeriesMeta {
  key: string
  label: string
  color: string
}

export type RangeKey = '7d' | '30d' | '90d' | 'custom' | 'all'

export interface RangeOption {
  key: RangeKey
  label: string
  /** Lookback window in days from now. null = no lower bound (or n/a for 'custom'). */
  days: number | null
}

export interface CustomRange {
  from: Date | null
  to: Date | null
}

/** Which lines of the filter bar are shown. All default to true. */
export interface ChartFilterVisibility {
  typeSwitch?: boolean
  rangeFilter?: boolean
  seriesFilter?: boolean
  search?: boolean
  pointsToggle?: boolean
}

export type CurveStyle = 'smooth' | 'straight'