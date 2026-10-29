export const shapes = ["subtle", "sharp", "soft", "round", "pill"] as const;
export type Shape = (typeof shapes)[number];

const classes: Record<Shape, string> = {
  subtle: "rounded-sm",
  sharp: "rounded-none",
  soft: "rounded-md",
  round: "rounded-xl",
  pill: "rounded-full",
};

/** Omitted shape preserves each component's existing radius. */
export function shapeClass(
  shape?: Shape,
  geometry: "control" | "surface" = "control",
) {
  if (shape === "pill" && geometry === "surface") return "rounded-3xl";
  return shape ? classes[shape] : undefined;
}
