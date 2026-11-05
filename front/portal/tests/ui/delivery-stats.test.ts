import { describe, expect, it } from "vitest";
import type { Deploy, PullRequest } from "../../src/features/delivery/model";
import {
  changeFailureRate,
  deploysPerWeek,
  prLeadTime,
  timeToRestore,
} from "../../src/features/delivery/stats";

const now = new Date("2026-11-30T12:00:00Z");
const hoursAgo = (h: number) => new Date(now.getTime() - h * 3_600_000).toISOString();

const deploy = (component: string, h: number, overrides: Partial<Deploy> = {}): Deploy => ({
  component,
  tag: `${component}/v1.0.0`,
  sha: "abc",
  actor: "MintzyG",
  status: "success",
  at: hoursAgo(h),
  runId: null,
  ...overrides,
});

describe("DORA numbers", () => {
  it("counts successful deploys per week over 28 days", () => {
    const deploys = [deploy("a", 1), deploy("a", 24 * 27), deploy("a", 24 * 29), deploy("b", 2, { status: "failure" })];

    expect(deploysPerWeek(deploys, now)).toBe(0.5);
  });

  it("takes the median PR lead time over merged PRs", () => {
    const pull = (openedH: number, mergedH: number | null): PullRequest => ({
      number: openedH,
      title: "",
      author: "",
      state: mergedH === null ? "open" : "merged",
      draft: false,
      reviewState: null,
      openedAt: hoursAgo(openedH),
      mergedAt: mergedH === null ? null : hoursAgo(mergedH),
    });

    expect(prLeadTime([pull(10, 8), pull(10, 4), pull(10, 0), pull(5, null)], now)).toBe(6 * 3_600_000);
  });

  it("counts failed deploys against the change failure rate", () => {
    const deploys = [
      deploy("a", 1),
      deploy("a", 2, { status: "failure" }),
      deploy("b", 3, { status: "failure" }),
      deploy("b", 4, { status: "cancelled" }),
    ];

    expect(changeFailureRate(deploys, now)).toBeCloseTo(2 / 3);
    expect(changeFailureRate([], now)).toBeNull();
  });

  it("measures time to restore from the first failure to the next good deploy", () => {
    const deploys = [
      deploy("a", 10, { status: "failure" }),
      deploy("a", 8, { status: "failure" }),
      deploy("a", 6),
      deploy("b", 5, { status: "failure" }),
      deploy("b", 4),
      deploy("c", 3, { status: "failure" }),
    ];

    // a: 10h → 6h = 4h; b: 5h → 4h = 1h; c is not restored yet.
    expect(timeToRestore(deploys, now)).toBe(2.5 * 3_600_000);
  });
});
