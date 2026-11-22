import { describe, expect, it } from "vitest";
import { docsIn, findDoc } from "../../src/features/docs/catalog";
import { durationOf, formatMoment, formatSpan, phasesOf, serviceName, severityLabel } from "../../src/features/docs/incident";

describe("incident reports", () => {
  it("lists the migrated reports newest first, preamble read", () => {
    const reports = docsIn("incident-reports");
    expect(reports.length).toBeGreaterThanOrEqual(2);
    expect(reports[0].meta.date! >= reports[1].meta.date!).toBe(true);
    expect(findDoc("incident-reports", "2026-08-26-host-saturation")?.meta).toMatchObject({
      severity: "high",
      status: "resolved",
      started: "2026-08-26T17:04Z",
    });
  });

  it("orders phases and measures the spans between them", () => {
    const phases = phasesOf({ started: "2026-10-19", detected: "2026-11-02T11:27Z", mitigated: "2026-11-02T11:40Z", ended: "2026-11-02T17:50Z" });
    expect(phases.map((p) => p.label)).toEqual(["Started", "Noticed", "Contained", "Over"]);
    expect(formatSpan(phases[0], phases[1])).toBe("about 14 days");
    expect(formatSpan(phases[1], phases[2])).toBe("13 min");
    expect(formatSpan(phases[2], phases[3])).toBe("6 h 10 min");
    expect(formatMoment(phases[1])).toBe("2 Nov, 11:27 UTC");
    expect(formatMoment(phases[0], true)).toBe("19 Oct 2026");
  });

  it("gives a duration only when the end is known", () => {
    expect(durationOf({ started: "2026-08-26T17:04Z", ended: "2026-08-26T17:34Z" })).toBe("30 min");
    expect(durationOf({ started: "2026-08-26T17:04Z" })).toBeUndefined();
  });

  it("names services and severities for people", () => {
    expect(serviceName("identityx")).toBe("IdentityX");
    expect(serviceName("all hosted services")).toBe("all hosted services");
    expect(severityLabel("medium")).toBe("Medium severity");
  });
});
