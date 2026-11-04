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
function stubApis() {
  const calls: string[] = [];
  vi.spyOn(globalThis, "fetch").mockImplementation(async (input, init) => {
    const url = String(input instanceof Request ? input.url : input);
    calls.push(url);
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
        project: { active_time: 3600, compute_last_active_at: "2026-11-04T18:21:00Z", quota_reset_at: "2026-12-01T00:00:00Z" },
      });
    }
    throw new Error(`unexpected fetch ${url}`);
  });
  return calls;
}

beforeEach(async () => {
  await env.DB.batch(
    ["metric_rollups", "usage_snapshots", "worker_stats", "neon_projects", "sync_state"].map((t) => env.DB.prepare(`DELETE FROM ${t}`)),
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

    expect(calls.every((url) => url.startsWith("https://api.cloudflare.com/client/v4/graphql") || url.startsWith("https://console.neon.tech/api/v2/"))).toBe(true);
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
    expect(neon).toMatchObject({ storageBytes: 30_000_000, activeSeconds: 3600, resumes: 2, resumeMaxMs: 1583 });
    // Neon's month began Nov 1; an hour awake by Nov 4 18:40 is ~98.9% asleep.
    const elapsed = (NOW.getTime() - Date.parse("2026-11-01T00:00:00Z")) / 1000;
    expect(neon.suspendedShare).toBeCloseTo(1 - 3600 / elapsed, 6);
    expect(payload.syncedAt).toBe(NOW.toISOString());
  });
});
