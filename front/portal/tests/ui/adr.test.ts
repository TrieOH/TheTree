import { describe, expect, it } from "vitest";
import { areasOf, byMonth, formatDate, isSuperseded, matches, relationsOf } from "../../src/features/docs/adr";
import { decisionByNumber, docsIn } from "../../src/features/docs/catalog";

describe("ADR helpers", () => {
  it("formats dates without shifting the day", () => {
    expect(formatDate("2026-11-05")).toBe("5 Nov 2026");
    expect(formatDate("2026-11-05", false)).toBe("5 Nov");
    expect(formatDate("soon")).toBe("soon");
  });

  it("sorts scope paths into areas", () => {
    const areas = areasOf({ scope: ["api/*", "lib/go/httpserver", "lib/ts/container-budget", "sdk/go/IdentityX", "front", ".husky"] });
    expect([...areas].sort()).toEqual(["backends", "fronts", "go-libs", "sdks", "tooling", "ts-libs"]);
  });

  it("lists relations with the replacement first", () => {
    expect(relationsOf({ superseded_by: "ADR-0011", amends: ["ADR-0003"] })).toEqual([
      ["superseded_by", ["ADR-0011"]],
      ["amends", ["ADR-0003"]],
    ]);
    expect(isSuperseded({ superseded_by: "ADR-0011" })).toBe(true);
    expect(isSuperseded({ status: "accepted" })).toBe(false);
  });

  it("groups the real decisions by month, newest first, and finds them by number", () => {
    const months = byMonth(docsIn("decisions"));
    expect(months.reduce((n, m) => n + m.docs.length, 0)).toBe(docsIn("decisions").length);
    expect(new Set(months.map((m) => m.month)).size).toBe(months.length);
    expect(decisionByNumber("ADR-0011")?.slug).toBe("adr-0011");
    expect(matches(decisionByNumber("ADR-0020")!, "edge")).toBe(true);
  });
});
