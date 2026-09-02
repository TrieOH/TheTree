import { describe, expect, it } from "vitest";
import { DEFAULT_BADGE_TEMPLATE } from "../default-template";
import { resizeBadgeDesign } from "./resize-design";

describe("resizeBadgeDesign", () => {
  it("scales bounds and text proportionally without mutating the source", () => {
    const resized = resizeBadgeDesign(DEFAULT_BADGE_TEMPLATE.design_data, {
      width: 642,
      height: 408,
    });
    const text = resized.elements[0];

    expect(text).toMatchObject({ x: 64, y: 36, width: 514, height: 40 });
    expect(text.type === "text" && text.paragraphs[0].runs[0].fontSize).toBe(
      18,
    );
    expect(DEFAULT_BADGE_TEMPLATE.design_data.canvas.width).toBe(321);
  });

  it("keeps text at the minimum supported font size", () => {
    const resized = resizeBadgeDesign(DEFAULT_BADGE_TEMPLATE.design_data, {
      width: 160.5,
      height: 102,
    });
    const text = resized.elements[0];

    expect(text.type === "text" && text.paragraphs[0].runs[0].fontSize).toBe(6);
  });
});
