import type { SeriesMeta } from './types'

/**
 * Categorical palette used across all chart types and the filter/legend chips.
 * Chosen for contrast against a white/slate-950 surface and to stay distinct
 * from the default Tailwind indigo/blue that shows up in every dashboard.
 */
export const CHART_PALETTE = [
  '#6C5CE7', // signal violet
  '#FF6B4A', // coral
  '#14B8A6', // teal
  '#F5A623', // amber
  '#3B82F6', // sky
  '#EC4899', // pink
] as const

export function colorForIndex(index: number): string {
  return CHART_PALETTE[index % CHART_PALETTE.length]
}

export function buildSeriesMeta(
  seriesKeys: string[],
  labels?: Record<string, string>,
  colors?: Record<string, string>,
): SeriesMeta[] {
  return seriesKeys.map((key, index) => ({
    key,
    label: labels?.[key] ?? key,
    color: colors?.[key] ?? colorForIndex(index),
  }))
}