import type { EditionPurchase } from "@trieoh/univents-api/schemas";
import { describe, expect, it } from "vitest";
import { buildEventOverviewMetrics } from "./event-overview";

const purchase = (
  overrides: Partial<EditionPurchase> = {},
): EditionPurchase => ({
  purchase_id: "purchase-1",
  edition_id: "edition-1",
  status: "approved",
  total_cents: 1_000,
  currency: "BRL",
  created_at: "2026-01-01T00:00:00Z",
  items: [],
  attendees: [],
  ...overrides,
});

describe("event overview metrics", () => {
  it("aggregates only approved revenue and builds chronological profit", () => {
    const result = buildEventOverviewMetrics({
      editions: [
        { id: "edition-1", name: "First" },
        { id: "edition-2", name: "Second" },
      ],
      purchasesByEdition: [
        [
          purchase({ purchase_id: "late", created_at: "2026-01-02T12:00:00Z" }),
          purchase({
            purchase_id: "refunded",
            status: "refunded",
            total_cents: 5_000,
          }),
        ],
        [
          purchase({
            purchase_id: "early",
            edition_id: "edition-2",
            total_cents: 2_000,
          }),
          purchase({ purchase_id: "invalid", created_at: "invalid" }),
        ],
      ],
      attendeeCounts: [2, 3],
      ticketCounts: [1, 2],
      productCounts: [3, 4],
      programCounts: [5, 6],
      occurrenceCounts: [7, 8],
    });

    expect(result.revenue).toBe(4_000);
    expect(result.refundedPurchaseCount).toBe(1);
    expect(result.participantCount).toBe(5);
    expect(result.ticketCount).toBe(3);
    expect(result.productCount).toBe(7);
    expect(result.programCount).toBe(11);
    expect(result.occurrenceCount).toBe(15);
    expect(result.editionSales.map((edition) => edition.revenue)).toEqual([
      1_000, 3_000,
    ]);
    expect(result.maxEditionRevenue).toBe(3_000);
    expect(
      result.profitData.map(({ value, purchases }) => ({
        value,
        purchases,
      })),
    ).toEqual([
      { value: 20, purchases: 1 },
      { value: 30, purchases: 1 },
    ]);
  });

  it("keeps chart scaling valid without editions or purchases", () => {
    const result = buildEventOverviewMetrics({
      editions: [],
      purchasesByEdition: [],
      attendeeCounts: [],
      ticketCounts: [],
      productCounts: [],
      programCounts: [],
      occurrenceCounts: [],
    });

    expect(result.maxEditionRevenue).toBe(1);
    expect(result.profitData).toEqual([]);
  });
});
