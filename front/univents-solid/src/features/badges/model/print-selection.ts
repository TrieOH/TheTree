import type { BadgeEditionEmission, BadgePrintItem } from "./index";

export function selectBadgePrintItems(
  printItems: BadgePrintItem[],
  emissions: BadgeEditionEmission[],
  printedAfter?: string,
): BadgePrintItem[] {
  const emissionMap = new Map(emissions.map((e) => [e.id, e]));
  return printItems.filter((item) => {
    const emission = emissionMap.get(item.emission_id);
    if (!emission || emission.status !== "active") return false;
    if (printedAfter) {
      const afterDate = new Date(printedAfter).getTime();
      const emissionDate = new Date(emission.emitted_at).getTime();
      if (Number.isFinite(afterDate) && Number.isFinite(emissionDate) && emissionDate < afterDate) {
        return false;
      }
    }
    return true;
  });
}
