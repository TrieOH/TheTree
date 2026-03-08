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
