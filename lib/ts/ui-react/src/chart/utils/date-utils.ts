export function startOfDay(date: Date): Date {
  const d = new Date(date)
  d.setHours(0, 0, 0, 0)
  return d
}

export function endOfDay(date: Date): Date {
  const d = new Date(date)
  d.setHours(23, 59, 59, 999)
  return d
}

/** Formats a Date as 'YYYY-MM-DD' for use as an <input type="date"> value. */
export function toDateInputValue(date: Date | null): string {
  if (!date) return ''
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

/** Parses an <input type="date"> value ('YYYY-MM-DD') as a local Date. */
export function fromDateInputValue(value: string): Date | null {
  if (!value) return null
  const [year, month, day] = value.split('-').map(Number)
  if (!year || !month || !day) return null
  return new Date(year, month - 1, day)
}

export function isSameDay(a: Date | null, b: Date | null): boolean {
  if (!a || !b) return false
  return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate()
}

export function isWithinRange(day: Date, start: Date, end: Date): boolean {
  const t = startOfDay(day).getTime()
  const [lo, hi] = start.getTime() <= end.getTime() ? [start, end] : [end, start]
  return t >= startOfDay(lo).getTime() && t <= startOfDay(hi).getTime()
}

export function addMonths(date: Date, delta: number): Date {
  return new Date(date.getFullYear(), date.getMonth() + delta, 1)
}

export function addDays(date: Date, delta: number): Date {
  const result = new Date(date)
  result.setDate(result.getDate() + delta)
  return result
}

/**
 * Builds a 6x7 (or fewer trailing rows trimmed) calendar grid for the given
 * month. `null` cells are padding before day 1 / after the last day.
 */
export function buildCalendarWeeks(monthDate: Date): (Date | null)[][] {
  const year = monthDate.getFullYear()
  const month = monthDate.getMonth()
  const firstWeekday = new Date(year, month, 1).getDay() // 0 = Sunday
  const daysInMonth = new Date(year, month + 1, 0).getDate()

  const cells: (Date | null)[] = []
  for (let i = 0; i < firstWeekday; i++) cells.push(null)
  for (let day = 1; day <= daysInMonth; day++) cells.push(new Date(year, month, day))
  while (cells.length % 7 !== 0) cells.push(null)

  const weeks: (Date | null)[][] = []
  for (let i = 0; i < cells.length; i += 7) weeks.push(cells.slice(i, i + 7))
  return weeks
}