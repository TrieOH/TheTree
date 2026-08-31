import { describe, expect, it } from "vitest";
import { drawSequence, drawTimeline } from "../lib/draw-sequence";

describe("drawSequence", () => {
  it("shows every eligible participant and fills short draws", () => {
    const participants = ["Ana", "Bia", "Caio"];
    const sequence = drawSequence(participants, 10);

    expect(sequence).toHaveLength(10);
    expect(new Set(sequence)).toEqual(new Set(participants));
    expect(sequence.length).toBeGreaterThan(new Set(sequence).size);
    expect(
      sequence.slice(1).every((name, index) => name !== sequence[index]),
    ).toBe(true);
  });

  it("slows the name changes near the end", () => {
    const { delays, durationMs } = drawTimeline(24);
    const firstGap = (delays[1] ?? 0) - (delays[0] ?? 0);
    const lastGap = (delays.at(-1) ?? 0) - (delays.at(-2) ?? 0);

    expect(lastGap).toBeGreaterThan(firstGap * 3);
    expect(durationMs).toBeCloseTo(4800);
  });
});
