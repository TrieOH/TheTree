export function formatDate(dateString: string): string {
  const date = new Date(dateString)

  const formatted = date.toLocaleDateString("en-US", {
    month: "long",
    day: "2-digit",
    year: "numeric",
    timeZone: "UTC"
  })

  return formatted
}
