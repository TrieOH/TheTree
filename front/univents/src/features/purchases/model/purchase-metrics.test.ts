import type { EditionPurchase } from "@trieoh/univents-api/schemas";
import { describe, expect, it } from "vitest";
import { buildPurchaseMetrics } from "./purchase-metrics";

const purchase = (
  status: EditionPurchase["status"],
  overrides: Partial<EditionPurchase> = {},
): EditionPurchase => ({
  purchase_id: `${status}-purchase`,
  edition_id: "edition-1",
  status,
  total_cents: 1_000,
  currency: "BRL",
  created_at: "2026-01-01T12:00:00Z",
  items: [],
  attendees: [],
  ...overrides,
});

describe("purchase metrics", () => {
  it("counts statuses and excludes non-approved revenue", () => {
    const result = buildPurchaseMetrics(
      [
        purchase("approved"),
        purchase("approved", { purchase_id: "approved-2" }),
        purchase("pending", { total_cents: 5_000 }),
        purchase("refunded", { total_cents: 5_000 }),
      ],
      "revenue",
    );

    expect(result.revenue).toBe(2_000);
    expect(result.refundedPurchaseCount).toBe(1);
    expect(result.statusCounts).toMatchObject({
      approved: 2,
      pending: 1,
      refunded: 1,
    });
    expect(result.profitData).toHaveLength(1);
    expect(result.profitData[0]).toMatchObject({
      value: 20,
      series: "revenue",
      purchases: 2,
    });
  });
});
