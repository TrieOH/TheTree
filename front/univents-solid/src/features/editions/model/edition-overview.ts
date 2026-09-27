import type { EditionPurchase } from "@trieoh/univents-api/schemas";
import type { EditionI } from "./index";
import {
  type PurchaseProfitDatum,
  buildPurchaseMetrics,
} from "@/features/purchases/model/purchase-metrics";

export interface EditionOverviewMetrics {
  revenue: number;
  refundedPurchaseCount: number;
  purchasesCount: number;
  attendeeCount: number;
  ticketCount: number;
  productCount: number;
  programCount: number;
  occurrenceCount: number;
  profitData: PurchaseProfitDatum[];
  statusCounts: Record<string, number>;
}

export function buildEditionOverviewMetrics(params: {
  edition: EditionI;
  purchases: EditionPurchase[];
  attendeeCount: number;
  ticketCount: number;
  productCount: number;
  programCount: number;
  occurrenceCount: number;
}): EditionOverviewMetrics {
  const { revenue, refundedPurchaseCount, statusCounts, profitData } =
    buildPurchaseMetrics(params.purchases, "revenue");

  return {
    revenue,
    refundedPurchaseCount,
    purchasesCount: params.purchases.length,
    attendeeCount: params.attendeeCount,
    ticketCount: params.ticketCount,
    productCount: params.productCount,
    programCount: params.programCount,
    occurrenceCount: params.occurrenceCount,
    profitData,
    statusCounts,
  };
}
