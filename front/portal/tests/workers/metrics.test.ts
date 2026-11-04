import { env } from "cloudflare:workers";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { statusOf, type MetricsPayload } from "../../src/features/metrics/model";
import { syncMetrics } from "../../src/features/metrics/server/sync";
import worker from "../../src/server";

const NOW = new Date("2026-11-04T18:40:00Z");

async function readApi(): Promise<MetricsPayload> {
  const response = await worker.fetch(new Request("http://localhost/api/metrics"), env);
  expect(response.status).toBe(200);
  return response.json();
}

/** Answers Cloudflare's GraphQL and Neon's API the way they answer; records every call. */
function stubApis(options: { coldStarts?: "ok" | "missing" } = {}) {
  const calls: string[] = [];
  vi.spyOn(globalThis, "fetch").mockImplementation(async (input, init) => {
    const url = String(input instanceof Request ? input.url : input);
    calls.push(url);
    if (url.includes("/analytics_engine/sql")) {
      if (options.coldStarts === "missing") return new Response("unknown table container_cold_starts", { status: 422 });
      // The SQL API returns numbers as strings in places; the sync coerces.
      return Response.json({
        data: [
          { service: "identityx", day: "2026-11-03 00:00:00", count: "6", failed: "2", p50: "4200.5", p95: "9800" },
          { service: "identityx", day: "2026-11-04 00:00:00", count: 3, failed: 0, p50: 3900, p95: 6100 },
          { service: "not-a-service", day: "2026-11-04 00:00:00", count: 9, failed: 0, p50: 1, p95: 1 },
        ],
      });
    }
    if (url.includes("/graphql")) {
      const { query, variables } = JSON.parse(String(init?.body));
      if (query.includes("httpRequestsAdaptiveGroups")) {
        const trieoh = variables.zone === "6fa6890c8818fbdec24d20546ab06f08";
        const groups = trieoh
          ? [
              { count: 40, dimensions: { datetimeHour: "2026-11-04T17:00:00Z", clientRequestHTTPHost: "api.payssage.trieoh.com", edgeResponseStatus: 200 } },
              { count: 10, dimensions: { datetimeHour: "2026-11-04T17:00:00Z", clientRequestHTTPHost: "api.payssage.trieoh.com", edgeResponseStatus: 502 } },
              { count: 20, dimensions: { datetimeHour: "2026-11-04T16:00:00Z", clientRequestHTTPHost: "api.identityx.trieoh.com", edgeResponseStatus: 200 } },
              { count: 1, dimensions: { datetimeHour: "2026-11-04T16:00:00Z", clientRequestHTTPHost: "api.identityx.trieoh.com", edgeResponseStatus: 504 } },
            ]
          : [];
        // Like the real API, only the hours inside the queried window.
        const inWindow = groups.filter(
          (g) => g.dimensions.datetimeHour >= variables.since.replace(".000Z", "Z") && g.dimensions.datetimeHour < variables.until,
        );
        return Response.json({ data: { viewer: { zones: [{ httpRequestsAdaptiveGroups: inWindow }] } } });
      }
      return Response.json({
        data: {
          viewer: {
            accounts: [
              {
                workersInvocationsAdaptive: [
                  { sum: { requests: 5000, errors: 2, cpuTimeUs: 4_000_000 }, quantiles: { cpuTimeP50: 600, cpuTimeP99: 9000, wallTimeP50: 20_000, wallTimeP99: 8_000_000 }, dimensions: { scriptName: "payssage-api" } },
                  { sum: { requests: 3000, errors: 9, cpuTimeUs: 30_000_000 }, quantiles: { cpuTimeP50: 3300, cpuTimeP99: 40_000, wallTimeP50: 50_000, wallTimeP99: 150_000 }, dimensions: { scriptName: "payssage" } },
                  { sum: { requests: 340, errors: 0, cpuTimeUs: 400_000 }, quantiles: { cpuTimeP50: 1000, cpuTimeP99: 2000, wallTimeP50: 1, wallTimeP99: 2 }, dimensions: { scriptName: "__unknown__" } },
                ],
                d1AnalyticsAdaptiveGroups: [{ sum: { rowsRead: 700_000, rowsWritten: 4600 } }],
                kvOperationsAdaptiveGroups: [
                  { sum: { requests: 5000 }, dimensions: { actionType: "read" } },
                  { sum: { requests: 150 }, dimensions: { actionType: "write" } },
                  { sum: { requests: 20 }, dimensions: { actionType: "delete" } },
                ],
              },
            ],
          },
        },
      });
    }
    if (url.includes("console.neon.tech")) {
      if (url.endsWith("/branches")) return Response.json({ branches: [{ default: true, logical_size: 30_000_000 }] });
      if (url.includes("/operations")) {
        return Response.json({
          operations: [
            { action: "suspend_compute", status: "finished", created_at: "2026-11-04T18:26:16Z", total_duration_ms: 431 },
            { action: "start_compute", status: "finished", created_at: "2026-11-04T18:20:50Z", total_duration_ms: 1583 },
            { action: "start_compute", status: "finished", created_at: "2026-11-03T06:11:29Z", total_duration_ms: 368 },
            { action: "start_compute", status: "finished", created_at: "2026-10-30T06:11:29Z", total_duration_ms: 999 },
          ],
          pagination: {},
        });
      }
      return Response.json({
        project: {
          active_time_seconds: 3600,
          compute_last_active_at: "2026-11-04T18:21:00Z",
          consumption_period_start: "2026-11-03T14:58:42Z",
          consumption_period_end: "2026-12-01T00:00:00Z",
        },
      });
    }
    throw new Error(`unexpected fetch ${url}`);
  });
  return calls;
}

beforeEach(async () => {
  await env.DB.batch(
    ["metric_rollups", "usage_snapshots", "worker_stats", "neon_projects", "cold_starts", "sync_state"].map((t) => env.DB.prepare(`DELETE FROM ${t}`)),
  );
  // The container guard's pool for this cycle: 2 GiB-hours so far.
  await env.BUDGET.put("awake:2026-10-24:payssage", String(2 * 3600));
});

afterEach(() => vi.restoreAllMocks());

describe("statusOf", () => {
  it("calls a quiet service idle, never down", () => {
    expect(statusOf(null, null)).toBe("unknown");
    expect(statusOf(0, 0)).toBe("idle");
    expect(statusOf(100, 2)).toBe("operational");
    expect(statusOf(100, 5)).toBe("degraded");
    // A couple of errors on almost no traffic isn't degraded.
    expect(statusOf(2, 2)).toBe("operational");
  });
});

describe("syncMetrics", () => {
  it("reads only Cloudflare analytics, Neon's API and the budget KV", async () => {
    const calls = stubApis();
    vi.useFakeTimers({ now: NOW, toFake: ["Date"] });
    await syncMetrics(env, NOW);
    vi.useRealTimers();

    expect(
      calls.every(
        (url) =>
          url.startsWith("https://api.cloudflare.com/client/v4/graphql") ||
          url.startsWith("https://api.cloudflare.com/client/v4/accounts/fdc6590fe9b028c608005c515140ff74/analytics_engine/sql") ||
          url.startsWith("https://console.neon.tech/api/v2/"),
      ),
    ).toBe(true);
    expect(calls.some((url) => /api\.(identityx|payssage|univents|informd)\.|univents\.com\.br/.test(url) && !url.includes("graphql"))).toBe(false);
  });

  it("serves status, success rate, plan usage, Workers and Neon", async () => {
    stubApis();
    await syncMetrics(env, NOW);
    vi.useFakeTimers({ now: NOW, toFake: ["Date"] });
    const payload = await readApi();
    vi.useRealTimers();

    const payssage = payload.services.find((s) => s.id === "payssage")!;
    expect(payssage).toMatchObject({ status: "degraded", requests24h: 50, errors24h: 10, successRate30d: 0.8 });
    expect(payload.services.find((s) => s.id === "identityx")).toMatchObject({ status: "operational", requests24h: 21 });
    // Pulled, but nobody called it: idle, not down.
    expect(payload.services.find((s) => s.id === "univents")).toMatchObject({ status: "idle", requests24h: 0, successRate30d: null });

    const usage = Object.fromEntries(payload.usage.map((u) => [u.resource, u]));
    expect(usage.containers_memory).toMatchObject({ used: 2, included: 25, stopsAt: 22.5 });
    expect(usage.workers_requests.used).toBe(8340);
    expect(usage.workers_cpu.used).toBe(34_400);
    expect(usage.kv_writes.used).toBe(170);
    expect(payload.cycle.start).toBe("2026-10-24T00:00:00Z");

    expect(payload.workers.map((w) => [w.script, w.kind])).toEqual([
      ["payssage-api", "api"],
      ["payssage", "front"],
    ]);
    expect(payload.workers[0].wallP99Ms).toBe(8000);

    const neon = payload.neon.find((n) => n.service === "payssage")!;
    // Only resumes inside the project's consumption period (from Nov 3, 14:58).
    expect(neon).toMatchObject({ storageBytes: 30_000_000, activeSeconds: 3600, resumes: 1, resumeMaxMs: 1583 });
    expect(neon.periodStart).toBe("2026-11-03T14:58:42.000Z");
    const elapsed = (NOW.getTime() - Date.parse("2026-11-03T14:58:42Z")) / 1000;
    expect(neon.suspendedShare).toBeCloseTo(1 - 3600 / elapsed, 6);
    expect(payload.syncedAt).toBe(NOW.toISOString());
  });

  it("stores cold starts per service per day from Analytics Engine", async () => {
    stubApis();
    await syncMetrics(env, NOW);
    vi.useFakeTimers({ now: NOW, toFake: ["Date"] });
    const { coldStarts } = await readApi();
    vi.useRealTimers();

    expect(coldStarts).toEqual([
      { service: "identityx", day: "2026-11-03", count: 6, failed: 2, p50Ms: 4201, p95Ms: 9800 },
      { service: "identityx", day: "2026-11-04", count: 3, failed: 0, p50Ms: 3900, p95Ms: 6100 },
    ]);
  });

  it("still lands the rest of the pull while the cold-start dataset doesn't exist", async () => {
    stubApis({ coldStarts: "missing" });
    await syncMetrics(env, NOW);
    vi.useFakeTimers({ now: NOW, toFake: ["Date"] });
    const payload = await readApi();
    vi.useRealTimers();

    expect(payload.coldStarts).toEqual([]);
    expect(payload.syncedAt).toBe(NOW.toISOString());
    expect(payload.usage.length).toBeGreaterThan(0);
  });
});
