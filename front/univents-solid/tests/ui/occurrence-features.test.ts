import { describe, expect, it } from "vitest";
import {
  drawSequence,
  drawTimeline,
  randomItem,
} from "../../src/features/programs/lib/draw-sequence";
import { actorIdFromQr } from "../../src/features/programs/lib/actor-id-from-qr";
import { formatOccurrenceSchedule } from "../../src/features/programs/lib/format-schedule";

describe("draw-sequence", () => {
  it("randomItem selects an item from non-empty array", () => {
    const items = ["Alice", "Bob", "Charlie"];
    const chosen = randomItem(items);
    expect(chosen).toBeDefined();
    expect(items).toContain(chosen);
  });

  it("randomItem returns undefined for empty array", () => {
    expect(randomItem([])).toBeUndefined();
  });

  it("drawSequence creates a sequence with minimum length", () => {
    const items = ["A", "B", "C", "D"];
    const seq = drawSequence(items, 20);
    expect(seq.length).toBeGreaterThanOrEqual(20);
    seq.forEach((item) => {
      expect(items).toContain(item);
    });
  });

  it("drawSequence returns empty array for empty inputs", () => {
    expect(drawSequence([])).toEqual([]);
  });

  it("drawTimeline creates delays with progressive deceleration", () => {
    const timeline = drawTimeline(10, 3000);
    expect(timeline.delays.length).toBe(10);
    expect(timeline.delays[0]).toBe(0);
    expect(timeline.durationMs).toBeCloseTo(3000, 0);

    // Gaps between consecutive frames should increase (easing slowdown)
    const gaps = [];
    for (let i = 1; i < timeline.delays.length; i++) {
      gaps.push(timeline.delays[i]! - timeline.delays[i - 1]!);
    }
    expect(gaps[gaps.length - 1]).toBeGreaterThan(gaps[0]!);
  });
});

describe("actorIdFromQr", () => {
  it("extracts ID from plain string", () => {
    expect(actorIdFromQr("user-123")).toBe("user-123");
  });

  it("extracts ID from URL path", () => {
    expect(actorIdFromQr("https://univents.app/badges/user-456/")).toBe("user-456");
  });

  it("handles URL encoded characters", () => {
    expect(actorIdFromQr("https%3A%2F%2Funivents.app%2Fbadges%2Fuser-789")).toBe("user-789");
  });
});

describe("formatOccurrenceSchedule", () => {
  it("returns empty string when startsAt is undefined or null", () => {
    expect(formatOccurrenceSchedule(undefined)).toBe("");
    expect(formatOccurrenceSchedule(null)).toBe("");
  });

  it("formats start and end time correctly", () => {
    const start = "2026-09-17T14:00:00.000Z";
    const end = "2026-09-17T15:30:00.000Z";
    const result = formatOccurrenceSchedule(start, end);
    expect(result).toContain("17");
    expect(result).toContain("–");
  });

  it("formats single date when endsAt is not provided", () => {
    const start = "2026-09-17T14:00:00.000Z";
    const result = formatOccurrenceSchedule(start);
    expect(result).toContain("17");
  });
});
