/**
 * Format a date string to a human-readable format.
 * Example: "January 15, 2024"
 */
export function formatDate(dateString: string): string {
  const date = new Date(dateString)

  const formatted = date.toLocaleDateString("en-US", {
    month: "long",
    day: "2-digit",
    year: "numeric",
    timeZone: "UTC",
  })

  return formatted
}

/**
 * Format a date string to a relative time (e.g., "3 hours ago").
 */
export function timeAgo(dateString: string): string {
  const rtf = new Intl.RelativeTimeFormat("en", { numeric: "auto" })

  const date = new Date(dateString)
  const diffInSeconds = (date.getTime() - Date.now()) / 1000

  const divisions = [
    { amount: 60, unit: "second" as const },
    { amount: 60, unit: "minute" as const },
    { amount: 24, unit: "hour" as const },
    { amount: 7, unit: "day" as const },
    { amount: 4.34524, unit: "week" as const },
    { amount: 12, unit: "month" as const },
    { amount: Number.POSITIVE_INFINITY, unit: "year" as const },
  ]

  let duration = diffInSeconds

  for (const division of divisions) {
    if (Math.abs(duration) < division.amount)
      return rtf.format(Math.round(duration), division.unit)

    duration /= division.amount
  }

  return "just now"
}
