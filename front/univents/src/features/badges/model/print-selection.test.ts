import { describe, expect, it } from "vitest";
import type { BadgeEditionEmission, BadgePrintItem } from ".";
import { selectBadgePrintItems } from "./print-selection";

const item = (emissionId: string): BadgePrintItem => ({
  emission_id: emissionId,
  user_id: `user-${emissionId}`,
  origin: "participant",
  event_name: "Event",
  edition_name: "Edition",
  action_url: `https://example.com/${emissionId}`,
});

const emission = (id: string, emittedAt: string): BadgeEditionEmission => ({
  id,
  user_id: `user-${id}`,
  origin: "participant",
  status: "active",
  emitted_at: emittedAt,
});

describe("badge print selection", () => {
  const items = [item("before"), item("equal"), item("after")];
  const emissions = [
    emission("before", "2026-01-01T11:59:59Z"),
    emission("equal", "2026-01-01T12:00:00Z"),
    emission("after", "2026-01-01T12:30:00Z"),
  ];

  it("includes emissions at or after the selected instant", () => {
    expect(
      selectBadgePrintItems(items, emissions, "2026-01-01T12:00:00Z").map(
        (badge) => badge.emission_id,
      ),
    ).toEqual(["equal", "after"]);
  });

  it("keeps every item without a date and rejects invalid dates", () => {
    expect(selectBadgePrintItems(items, emissions, "")).toBe(items);
    expect(selectBadgePrintItems(items, emissions, "invalid")).toEqual([]);
  });
});
