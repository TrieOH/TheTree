import { describe, expect, it } from "vitest";
import type { EditionPurchase } from "@trieoh/univents-api/schemas";
import type { EditionI } from "@/features/editions/model";
import { buildEventOverviewMetrics } from "@/features/events/model/event-overview";
import { buildPurchaseMetrics } from "@/features/purchases/model/purchase-metrics";

describe("Event Overview & Purchase Metrics", () => {
  it("computes purchase metrics correctly for approved and refunded purchases", () => {
    const purchases: EditionPurchase[] = [
      {
        purchase_id: "p1",
        edition_id: "ed1",
        status: "approved",
        total_cents: 5000,
        currency: "BRL",
        created_at: "2026-03-01T10:00:00Z",
        items: [],
        attendees: [],
      },
      {
        purchase_id: "p2",
        edition_id: "ed1",
        status: "approved",
        total_cents: 7500,
        currency: "BRL",
        created_at: "2026-03-02T12:00:00Z",
        items: [],
        attendees: [],
      },
      {
        purchase_id: "p3",
        edition_id: "ed1",
        status: "refunded",
        total_cents: 3000,
        currency: "BRL",
        created_at: "2026-03-03T14:00:00Z",
        items: [],
        attendees: [],
      },
      {
        purchase_id: "p4",
        edition_id: "ed1",
        status: "pending",
        total_cents: 2000,
        currency: "BRL",
        created_at: "2026-03-04T15:00:00Z",
        items: [],
        attendees: [],
      },
    ];

    const metrics = buildPurchaseMetrics(purchases, "Edição 2026");

    expect(metrics.revenue).toBe(12500);
    expect(metrics.statusCounts.approved).toBe(2);
    expect(metrics.refundedPurchaseCount).toBe(1);
    expect(metrics.profitData).toHaveLength(2);
    expect(metrics.profitData[0].value).toBe(50);
    expect(metrics.profitData[1].value).toBe(125);
  });

  it("aggregates multiple editions into event overview metrics", () => {
    const mockEdition1: EditionI = {
      id: "ed1",
      event_id: "ev1",
      name: "Edição 1",
      slug: "ed-1",
      is_draft: false,
      starts_at: "2026-04-01T00:00:00Z",
      ends_at: "2026-04-03T00:00:00Z",
      created_by: "u1",
      created_at: "2026-01-01T00:00:00Z",
      status: "future",
    };

    const mockEdition2: EditionI = {
      id: "ed2",
      event_id: "ev1",
      name: "Edição 2",
      slug: "ed-2",
      is_draft: false,
      starts_at: "2026-05-01T00:00:00Z",
      ends_at: "2026-05-03T00:00:00Z",
      created_by: "u1",
      created_at: "2026-01-02T00:00:00Z",
      status: "future",
    };

    const purchasesEd1: EditionPurchase[] = [
      {
        purchase_id: "p1",
        edition_id: "ed1",
        status: "approved",
        total_cents: 10000,
        currency: "BRL",
        created_at: "2026-03-01T10:00:00Z",
        items: [],
        attendees: [],
      },
    ];

    const purchasesEd2: EditionPurchase[] = [
      {
        purchase_id: "p2",
        edition_id: "ed2",
        status: "approved",
        total_cents: 20000,
        currency: "BRL",
        created_at: "2026-03-02T10:00:00Z",
        items: [],
        attendees: [],
      },
    ];

    const overview = buildEventOverviewMetrics({
      editions: [mockEdition1, mockEdition2],
      purchasesByEdition: [purchasesEd1, purchasesEd2],
      attendeeCounts: [15, 30],
      ticketCounts: [2, 4],
      productCounts: [3, 5],
      programCounts: [1, 2],
      occurrenceCounts: [4, 8],
    });

    expect(overview.revenue).toBe(30000);
    expect(overview.participantCount).toBe(45);
    expect(overview.ticketCount).toBe(6);
    expect(overview.productCount).toBe(8);
    expect(overview.programCount).toBe(3);
    expect(overview.occurrenceCount).toBe(12);
    expect(overview.editionSales).toEqual([
      { name: "Edição 1", revenue: 10000, purchases: 1 },
      { name: "Edição 2", revenue: 20000, purchases: 1 },
    ]);
    expect(overview.maxEditionRevenue).toBe(20000);
  });

  it("handles empty editions gracefully", () => {
    const overview = buildEventOverviewMetrics({
      editions: [],
      purchasesByEdition: [],
      attendeeCounts: [],
      ticketCounts: [],
      productCounts: [],
      programCounts: [],
      occurrenceCounts: [],
    });

    expect(overview.revenue).toBe(0);
    expect(overview.participantCount).toBe(0);
    expect(overview.editionSales).toEqual([]);
    expect(overview.maxEditionRevenue).toBe(1);
  });
});
