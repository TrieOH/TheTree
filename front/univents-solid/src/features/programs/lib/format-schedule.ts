/**
 * Formats occurrence start and end dates into a readable schedule string in Brazilian Portuguese.
 * Example: "Sex, 10 de set • 09:00 – 10:30"
 */
export function formatOccurrenceSchedule(
  startsAt?: string | null,
  endsAt?: string | null,
): string {
  if (!startsAt) return "";
  try {
    const sDate = new Date(startsAt);
    if (Number.isNaN(sDate.getTime())) return "";

    const datePart = new Intl.DateTimeFormat("pt-BR", {
      weekday: "short",
      day: "2-digit",
      month: "short",
    })
      .format(sDate)
      .replace(".", "");

    const capitalizedDate =
      datePart.charAt(0).toUpperCase() + datePart.slice(1);
    const sTime = sDate.toLocaleTimeString("pt-BR", {
      hour: "2-digit",
      minute: "2-digit",
    });

    if (!endsAt) return `${capitalizedDate} • ${sTime}`;

    const eDate = new Date(endsAt);
    if (Number.isNaN(eDate.getTime())) return `${capitalizedDate} • ${sTime}`;

    const eTime = eDate.toLocaleTimeString("pt-BR", {
      hour: "2-digit",
      minute: "2-digit",
    });
    return `${capitalizedDate} • ${sTime} – ${eTime}`;
  } catch {
    return startsAt;
  }
}
