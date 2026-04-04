import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'
import type { ClassValue } from 'clsx'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * Returns an object containing only the fields from 'data' that are different 
 * from their counterparts in 'original'.
 */
export function getDirtyFields<T extends Record<string, unknown>>(
  data: T,
  original: T,
  fieldsToCompare?: (keyof T)[]
): Partial<T> {
  const changes: Partial<T> = {}
  const keys = fieldsToCompare ?? (Object.keys(data) as (keyof T)[])

  keys.forEach((key) => {
    const newVal = data[key]
    const oldVal = original[key]

    // Normalizing null/undefined to compare empty states correctly
    const normalizedNew = newVal === "" ? null : (newVal ?? null)
    const normalizedOld = oldVal === "" ? null : (oldVal ?? null)

    const isObject = (val: unknown): val is Record<string, unknown> =>
      typeof val === 'object' && val !== null && !Array.isArray(val)

    if (isObject(normalizedNew) && isObject(normalizedOld)) {
      if (JSON.stringify(normalizedNew) !== JSON.stringify(normalizedOld)) changes[key] = newVal
    } else if (normalizedNew !== normalizedOld) changes[key] = newVal
  })

  return changes
}
