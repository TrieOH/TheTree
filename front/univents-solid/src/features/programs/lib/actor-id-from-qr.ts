export function actorIdFromQr(value: string): string {
  const trimmed = decodeURIComponent(value.trim()).replace(/\/+$/, "");
  return trimmed.split("/").pop() ?? trimmed;
}
