import { cleanup, render } from "@solidjs/testing-library";
import { afterEach, describe, expect, it } from "vitest";
import { formatBytes, formatMs, formatPercent } from "../../src/features/metrics/format";
import { cycleBounds } from "../../src/features/metrics/model";
import { UsageBar } from "../../src/features/metrics/ui/UsageBar";

afterEach(cleanup);

describe("metrics formatting", () => {
  it("formats durations, sizes and shares", () => {
    expect(formatMs(0.6)).toBe("600 µs");
    expect(formatMs(8000)).toBe("8.0 s");
    expect(formatBytes(30_000_000)).toBe("29 MB");
    expect(formatBytes(3_400_000)).toBe("3.2 MB");
    expect(formatPercent(0.8)).toBe("80.00%");
    expect(formatPercent(null)).toBe("—");
  });

  it("turns cycles over on the 24th like the container guard", () => {
    expect(cycleBounds(new Date("2026-11-04T12:00:00Z")).start.toISOString()).toBe("2026-10-24T00:00:00.000Z");
    expect(cycleBounds(new Date("2026-11-24T00:00:00Z")).end.toISOString()).toBe("2026-12-24T00:00:00.000Z");
  });
});

describe("<UsageBar>", () => {
  it("warns when usage runs ahead of an even pace, and marks the guard's cut-off", () => {
    const { container } = render(() => (
      <UsageBar usage={{ resource: "containers_memory", used: 10, included: 25, stopsAt: 22.5 }} pace={0.3} />
    ));

    expect(container.textContent).toContain("10.0 GiB-h of 25.0 GiB-h");
    expect(container.querySelector(".bg-warn")).not.toBeNull();
    expect(container.querySelector("[title='The container guard stops containers here']")).not.toBeNull();
    expect(container.textContent).toContain("stops every API container at 22.5 GiB-h");
  });

  it("stays calm behind the pace", () => {
    const { container } = render(() => (
      <UsageBar usage={{ resource: "workers_requests", used: 50_000, included: 10_000_000 }} pace={0.3} />
    ));

    expect(container.querySelector(".bg-warn")).toBeNull();
    expect(container.querySelector(".bg-ok")).not.toBeNull();
  });
});

describe("cold start summary", () => {
  it("weights the typical time by each day's starts and keeps the worst p95", async () => {
    const { summarize } = await import("../../src/features/metrics/cold-starts");
    const rows = [
      { service: "identityx" as const, day: "2026-11-03", count: 6, failed: 2, p50Ms: 4000, p95Ms: 9800 },
      { service: "identityx" as const, day: "2026-11-04", count: 2, failed: 0, p50Ms: 2000, p95Ms: 6100 },
      { service: "payssage" as const, day: "2026-11-04", count: 5, failed: 0, p50Ms: 1, p95Ms: 1 },
      { service: "identityx" as const, day: "2026-10-01", count: 99, failed: 9, p50Ms: 1, p95Ms: 99_999 },
    ];

    const week = summarize(rows, "identityx", 7, new Date("2026-11-04T19:00:00Z"));
    expect(week).toMatchObject({ count: 8, failed: 2, typicalMs: 3500, worstP95Ms: 9800 });
    expect(week.daily).toEqual([0, 0, 0, 0, 0, 6, 2]);
    expect(summarize(rows, "informd", 7, new Date("2026-11-04T19:00:00Z"))).toMatchObject({ count: 0, typicalMs: null });
  });
});
