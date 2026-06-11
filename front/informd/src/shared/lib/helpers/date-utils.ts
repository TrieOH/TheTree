export function timeAgo(dateString: string): string {
  const rtf = new Intl.RelativeTimeFormat("en", { numeric: "auto" });

  const date = new Date(dateString);
  const diffInSeconds = (date.getTime() - Date.now()) / 1000;

  const divisions = [
    { amount: 60, unit: "second" },
    { amount: 60, unit: "minute" },
    { amount: 24, unit: "hour" },
    { amount: 7, unit: "day" },
    { amount: 4.34524, unit: "week" },
    { amount: 12, unit: "month" },
    { amount: Number.POSITIVE_INFINITY, unit: "year" },
  ] as const;

  let duration = diffInSeconds;

  for (const division of divisions) {
    if (Math.abs(duration) < division.amount)
      return rtf.format(Math.round(duration), division.unit);

    duration /= division.amount;
  }

  return "just now";
}