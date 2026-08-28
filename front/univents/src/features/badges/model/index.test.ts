import { describe, expect, it } from "vitest";
import {
  badgeMmToPx,
  badgePxToMm,
  MIN_BADGE_CANVAS_SIZE_MM,
  MIN_BADGE_CANVAS_SIZE_PX,
} from ".";

describe("badge canvas size conversion", () => {
  it("keeps the minimum value usable by a number input", () => {
    expect(badgePxToMm(MIN_BADGE_CANVAS_SIZE_PX)).toBe(
      MIN_BADGE_CANVAS_SIZE_MM,
    );
    expect(badgeMmToPx(MIN_BADGE_CANVAS_SIZE_MM + 1)).toBeGreaterThan(
      MIN_BADGE_CANVAS_SIZE_PX,
    );
    expect(badgePxToMm(badgeMmToPx(26.5))).toBe(26.5);
  });
});
