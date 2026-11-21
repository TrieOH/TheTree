import { describe, expect, it } from "vitest";
import { cyclePace, needsAttention, projectedShare } from "../../src/features/metrics/attention";
import type { MetricsPayload } from "../../src/features/metrics/model";

const now = new Date("2026-11-21T23:00:00Z");
const service = (id: "identityx" | "payssage", extra: Partial<MetricsPayload["services"][number]> = {}) => ({
  id,
  status: "operational" as const,
  requests24h: 40,
  errors24h: 0,
  successRate30d: 0.999,
  hourly: [],
  lastRequestAt: null,
  p95Ms24h: 200,
  p95Hourly: [],
  ...extra,
});
const base: MetricsPayload = {
  cycle: { start: "2026-10-24T00:00:00Z", end: "2026-11-24T00:00:00Z" },
  services: [service("identityx"), service("payssage")],
  usage: [{ resource: "workers_requests", used: 50_000, included: 10_000_000 }],
  workers: [],
  neon: [],
  coldStarts: [],
  endpoints: [],
  problems: [],
  syncedAt: "2026-11-21T23:00:00Z",
};

describe("needsAttention", () => {
  it("stays quiet when nothing is off", () => {
    expect(needsAttention(base, now)).toEqual([]);
  });

  it("flags server errors, failing cold starts, erroring Workers and slow routes, worst first", () => {
    const items = needsAttention(
      {
        ...base,
        services: [service("identityx", { errors24h: 1 }), service("payssage")],
        coldStarts: [{ service: "identityx", day: "2026-11-19", count: 52, failed: 45, p50Ms: 2860, p95Ms: 7470 }],
        workers: [
          {
            script: "portal",
            service: null,
            kind: "portal",
            requests: 4000,
            errors: 100,
            cpuP50Ms: 1,
            cpuP99Ms: 1,
            wallP50Ms: 1,
            wallP99Ms: 1,
          },
          {
            script: "guara",
            service: null,
            kind: "other",
            requests: 100,
            errors: 50,
            cpuP50Ms: 1,
            cpuP99Ms: 1,
            wallP50Ms: 1,
            wallP99Ms: 1,
          },
        ],
        endpoints: [{ service: "payssage", method: "POST", route: "/payments", requests: 9, p95Ms: 2400 }],
      },
      now,
    );
    expect(items.map((item) => item.id)).toEqual([
      "cold-identityx",
      "errors-identityx",
      "worker-portal",
      "slow-payssage-POST-/payments",
    ]);
    expect(items[0]).toMatchObject({ tone: "destructive", tab: "cold-starts" });
    expect(items[0].text).toBe("IdentityX had 45 of 52 cold starts fail in the last 7 days.");
    expect(items[2].text).toContain("2.5%");
  });

  it("ignores a handful of cold starts even when some fail", () => {
    const items = needsAttention(
      { ...base, coldStarts: [{ service: "payssage", day: "2026-11-20", count: 3, failed: 2, p50Ms: 1, p95Ms: 1 }] },
      now,
    );
    expect(items).toEqual([]);
  });

  it("warns as usage nears the guard's cut-off, and alarms once it is on course to pass it", () => {
    const at = (used: number) =>
      needsAttention({ ...base, usage: [{ resource: "containers_memory", used, included: 25, stopsAt: 22.5 }] }, now);
    expect(at(20)).toMatchObject([{ tone: "warning", tab: "usage" }]);
    expect(at(22)).toMatchObject([{ tone: "destructive", tab: "usage" }]);
    expect(at(5)).toEqual([]);
  });
});

describe("usage projection", () => {
  it("extends usage so far to the end of the cycle, once there is enough of it", () => {
    expect(cyclePace(base.cycle, new Date("2026-11-08T12:00:00Z"))).toBeCloseTo(0.5, 2);
    expect(projectedShare({ resource: "kv_reads", used: 25, included: 100 }, 0.5)).toBe(0.5);
    expect(projectedShare({ resource: "kv_reads", used: 25, included: 100 }, 0.05)).toBeNull();
  });
});
