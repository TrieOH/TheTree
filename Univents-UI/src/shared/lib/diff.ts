function normalize(value: unknown): unknown {
  if (value === "" || value === undefined) return null
  if (Array.isArray(value)) {
    if (value.length === 0) return null
    return value.map(normalize)
  }
  if (value && typeof value === "object") {
    return Object.fromEntries(
      Object.entries(value).map(([k, v]) => [k, normalize(v)])
    )
  }
  return value ?? null
}

function deepEqual(a: unknown, b: unknown): boolean {
  if (a === b) return true

  if (Array.isArray(a) && Array.isArray(b)) {
    if (a.length !== b.length) return false
    return a.every((v, i) => deepEqual(v, b[i]))
  }

  if (a && b && typeof a === "object" && typeof b === "object") {
    const keysA = Object.keys(a)
    const keysB = Object.keys(b)

    if (keysA.length !== keysB.length) return false

    return keysA.every((key) =>
      deepEqual(
        (a as Record<string, unknown>)[key],
        (b as Record<string, unknown>)[key]
      )
    )
  }

  return false
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
    const newVal = normalize(data[key])
    const oldVal = normalize(original[key])

    if (!deepEqual(newVal, oldVal)) changes[key] = data[key]
  })

  return changes
}