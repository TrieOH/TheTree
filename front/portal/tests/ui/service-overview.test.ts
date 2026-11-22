import { describe, expect, it } from "vitest";
import { coldStartDays, serviceAttention } from "../../src/features/catalog/overview";
import { findService } from "../../src/features/catalog/services";
import type { MetricsPayload } from "../../src/features/metrics/model";

const now = new Date("2026-11-21T23:00:00Z");
const metrics: MetricsPayload = {
  cycle: { start: "2026-10-24T00:00:00Z", end: "2026-11-24T00:00:00Z" },
  services: (["identityx", "payssage"] as const).map((id) => ({
    id,
    status: "operational" as const,
    requests24h: 10,
    errors24h: 0,
    successRate30d: 1,
    hourly: [],
    lastRequestAt: null,
    p95Ms24h: 100,
    p95Hourly: [],
  })),
  usage: [],
  neon: [],
  endpoints: [],
  problems: [],
  syncedAt: "2026-11-21T23:00:00Z",
  coldStarts: [
    { service: "identityx", day: "2026-11-19", count: 52, failed: 45, p50Ms: 2860, p95Ms: 7470 },
    { service: "payssage", day: "2026-11-19", count: 10, failed: 9, p50Ms: 2000, p95Ms: 4000 },
  ],
  workers: [
    { script: "identityx", service: "identityx", kind: "front", requests: 1000, errors: 50, cpuP50Ms: 1, cpuP99Ms: 1, wallP50Ms: 1, wallP99Ms: 1 },
    { script: "payssage", service: "payssage", kind: "front", requests: 1000, errors: 50, cpuP50Ms: 1, cpuP99Ms: 1, wallP50Ms: 1, wallP99Ms: 1 },
  ],
};

describe("service overview", () => {
  it("keeps only the items about this service, its front and its Workers", () => {
    const items = serviceAttention(findService("identityx")!, metrics, null, now);
    expect(items.map((i) => [i.tab, i.text])).toEqual([
      ["metrics", "IdentityX had 45 of 52 cold starts fail in the last 7 days."],
      ["metrics", "The identityx Worker threw on 5.0% of its requests (50 of 1K) this cycle."],
    ]);
  });

  it("lays cold starts out one bar per day, zero where none were recorded", () => {
    const bars = coldStartDays(metrics.coldStarts, "identityx", 3, now);
    expect(bars).toEqual([
      { day: "2026-11-19", count: 52, failed: 45 },
      { day: "2026-11-20", count: 0, failed: 0 },
      { day: "2026-11-21", count: 0, failed: 0 },
    ]);
  });
});
