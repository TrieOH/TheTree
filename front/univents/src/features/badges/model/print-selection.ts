import type { BadgeEditionEmission, BadgePrintItem } from ".";

export function selectBadgePrintItems(
  items: BadgePrintItem[],
  emissions: BadgeEditionEmission[],
  emittedAfter: string,
) {
  if (!emittedAfter) return items;
  const timestamp = new Date(emittedAfter).getTime();
  if (Number.isNaN(timestamp)) return [];

  const emissionIds = new Set(
    emissions
      .filter((emission) => {
        const emittedAt = new Date(emission.emitted_at).getTime();
        return !Number.isNaN(emittedAt) && emittedAt >= timestamp;
      })
      .map((emission) => emission.id),
  );
  return items.filter((item) => emissionIds.has(item.emission_id));
}
