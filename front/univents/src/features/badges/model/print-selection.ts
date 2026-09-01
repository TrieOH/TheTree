import type { BadgeEditionEmission, BadgePrintItem } from ".";

export function selectBadgePrintItems(
  items: BadgePrintItem[],
  emissions: BadgeEditionEmission[],
  emittedAfter: string,
) {
  const timestamp = emittedAfter ? new Date(emittedAfter).getTime() : null;
  if (timestamp !== null && Number.isNaN(timestamp)) return [];

  const emissionIds = new Set(
    emissions
      .filter((emission) => {
        if (emission.status !== "active") return false;
        if (timestamp === null) return true;
        const emittedAt = new Date(emission.emitted_at).getTime();
        return !Number.isNaN(emittedAt) && emittedAt >= timestamp;
      })
      .map((emission) => emission.id),
  );
  return items.filter((item) => emissionIds.has(item.emission_id));
}
