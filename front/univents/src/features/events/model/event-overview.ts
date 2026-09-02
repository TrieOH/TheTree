import type { EditionPurchase } from "@trieoh/univents-api/schemas";
import { buildPurchaseMetrics } from "../../purchases/model/purchase-metrics";

export interface EventOverviewEdition {
  id: string;
  name: string;
}

export function buildEventOverviewMetrics({
  editions,
  purchasesByEdition,
  attendeeCounts,
  ticketCounts,
  productCounts,
  programCounts,
  occurrenceCounts,
}: {
  editions: EventOverviewEdition[];
  purchasesByEdition: EditionPurchase[][];
  attendeeCounts: number[];
  ticketCounts: number[];
  productCounts: number[];
  programCounts: number[];
  occurrenceCounts: number[];
}) {
  const purchases = purchasesByEdition.flat();
  const purchaseMetrics = buildPurchaseMetrics(purchases, "Lucro");

  const editionSales = editions.map((edition, index) => {
    const editionPurchases = purchasesByEdition[index] ?? [];
    return {
      name: edition.name,
      purchases: editionPurchases.length,
      revenue: editionPurchases
        .filter((purchase) => purchase.status === "approved")
        .reduce((total, purchase) => total + purchase.total_cents, 0),
    };
  });

  return {
    purchases,
    revenue: purchaseMetrics.revenue,
    refundedPurchaseCount: purchaseMetrics.refundedPurchaseCount,
    participantCount: sum(attendeeCounts),
    ticketCount: sum(ticketCounts),
    productCount: sum(productCounts),
    programCount: sum(programCounts),
    occurrenceCount: sum(occurrenceCounts),
    profitData: purchaseMetrics.profitData,
    editionSales,
    maxEditionRevenue: Math.max(
      ...editionSales.map((edition) => edition.revenue),
      1,
    ),
  };
}

export type EventOverviewMetrics = ReturnType<typeof buildEventOverviewMetrics>;

function sum(values: number[]) {
  return values.reduce((total, value) => total + value, 0);
}
