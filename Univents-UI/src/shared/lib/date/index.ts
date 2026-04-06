/**
 * Adjusts a Date object to the local timezone and formats it for a datetime-local input.
 * @param date The Date object to format.
 * @returns A string in `YYYY-MM-DDTHH:mm` format representing the local time.
 */
export function formatDateForDatetimeLocal(date: Date): string {
  const offset = date.getTimezoneOffset();
  const adjustedDate = new Date(date.getTime() - offset * 60 * 1000);
  return adjustedDate.toISOString().slice(0, 16);
}

/**
 * Parses a string from a datetime-local input into a Date object.
 * The input string is assumed to be in the user's local timezone.
 * @param localDateTimeString The `YYYY-MM-DDTHH:mm` string from the input.
 * @returns A Date object.
 */
export function parseDatetimeLocal(localDateTimeString: string): Date {
  return new Date(localDateTimeString);
}

/**
 * Formats a date range between two ISO date strings into a human-readable format.
 * If both dates are in the same month and year, only the day range is shown before the month/year.
 * @param starts The ISO date string for the start of the range.
 * @param ends The ISO date string for the end of the range.
 * @returns A formatted string representing the date range (e.g., "10 – 15 de set. de 2026").
 */
export function formatDateRange(starts: string, ends: string) {
  const s = new Date(starts)
  const e = new Date(ends)
  const opts: Intl.DateTimeFormatOptions = { day: 'numeric', month: 'short' }
  if (s.getMonth() === e.getMonth() && s.getFullYear() === e.getFullYear()) {
    return `${s.getDate()} – ${e.toLocaleDateString('pt-BR', { ...opts, year: 'numeric' })}`
  }
  return `${s.toLocaleDateString('pt-BR', opts)} – ${e.toLocaleDateString('pt-BR', { ...opts, year: 'numeric' })}`
}