import type { EditionPurchase } from "@trieoh/univents-api/schemas";

export interface PurchaseProfitDatum {
  [key: string]: unknown;
  date: Date;
  value: number;
  series: string;
  purchases: number;
}

export function buildPurchaseMetrics(
  purchases: EditionPurchase[],
  series: string,
) {
  const approved = purchases.filter(
    (purchase) => purchase.status === "approved",
  );
  const statusCounts: Record<string, number> = {};
  for (const purchase of purchases) {
    statusCounts[purchase.status] = (statusCounts[purchase.status] ?? 0) + 1;
  }

  const revenueByDay = new Map<
    number,
    { revenueCents: number; purchases: number }
  >();
  for (const purchase of approved) {
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
  const profitData: PurchaseProfitDatum[] = [...revenueByDay.entries()]
    .sort(([left], [right]) => left - right)
    .map(([timestamp, day]) => {
      accumulatedRevenue += day.revenueCents;
      return {
        date: new Date(timestamp),
        value: accumulatedRevenue / 100,
        series,
        purchases: day.purchases,
      };
    });

  return {
    revenue: approved.reduce(
      (total, purchase) => total + purchase.total_cents,
      0,
    ),
    refundedPurchaseCount: statusCounts.refunded ?? 0,
    statusCounts,
    profitData,
  };
}
