import type { BadgeDesign, BadgeElement } from "./index";

export function resizeBadgeDesign(
  design: BadgeDesign,
  canvas: { width: number; height: number },
): BadgeDesign {
  if (
    !design.canvas.width ||
    !design.canvas.height ||
    !canvas.width ||
    !canvas.height
  ) {
    return { ...design, canvas };
  }

  const scaleX = canvas.width / design.canvas.width;
  const scaleY = canvas.height / design.canvas.height;
  const fontScale = Math.sqrt(scaleX * scaleY);

  return {
    ...design,
    canvas,
    elements: design.elements.map((element) => {
      const nextX = Math.round(element.x * scaleX);
      const nextY = Math.round(element.y * scaleY);
      const nextWidth = Math.max(4, Math.round(element.width * scaleX));
      const nextHeight = Math.max(4, Math.round(element.height * scaleY));

      if (element.type === "text") {
        return {
          ...element,
          x: nextX,
          y: nextY,
          width: nextWidth,
          height: nextHeight,
          paragraphs: element.paragraphs.map((paragraph) => ({
            ...paragraph,
            runs: paragraph.runs.map((run) => ({
              ...run,
              fontSize: Math.max(6, Math.round(run.fontSize * fontScale)),
            })),
          })),
        };
      }

      return {
        ...element,
        x: nextX,
        y: nextY,
        width: nextWidth,
        height: nextHeight,
      };
    }) as BadgeElement[],
  };
}
