import type { EditionPurchase } from "@trieoh/univents-api/schemas";

export interface EventOverviewEdition {
  id: string;
  name: string;
}

export interface EventProfitDatum {
  [key: string]: unknown;
  date: Date;
  value: number;
  series: string;
  purchases: number;
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
  const approvedPurchases = purchases.filter(
    (purchase) => purchase.status === "approved",
  );
  const revenueByDay = new Map<
    number,
    { revenueCents: number; purchases: number }
  >();

  for (const purchase of approvedPurchases) {
    if (!purchase.created_at) continue;
    const date = new Date(purchase.created_at);
    if (Number.isNaN(date.getTime())) continue;
    date.setHours(0, 0, 0, 0);
    const timestamp = date.getTime();
    const day = revenueByDay.get(timestamp) ?? {
      revenueCents: 0,
      purchases: 0,
    };
    day.revenueCents += purchase.total_cents;
    day.purchases += 1;
    revenueByDay.set(timestamp, day);
  }

  let accumulatedRevenue = 0;
  const profitData: EventProfitDatum[] = [...revenueByDay.entries()]
    .sort(([left], [right]) => left - right)
    .map(([timestamp, day]) => {
      accumulatedRevenue += day.revenueCents;
      return {
        date: new Date(timestamp),
        value: accumulatedRevenue / 100,
        series: "Lucro",
        purchases: day.purchases,
      };
    });

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
    revenue: approvedPurchases.reduce(
      (total, purchase) => total + purchase.total_cents,
      0,
    ),
    refundedPurchaseCount: purchases.filter(
      (purchase) => purchase.status === "refunded",
    ).length,
    participantCount: sum(attendeeCounts),
    ticketCount: sum(ticketCounts),
    productCount: sum(productCounts),
    programCount: sum(programCounts),
    occurrenceCount: sum(occurrenceCounts),
    profitData,
    editionSales,
    maxEditionRevenue: Math.max(
      ...editionSales.map((edition) => edition.revenue),
      1,
    ),
  };
}

function sum(values: number[]) {
  return values.reduce((total, value) => total + value, 0);
}
