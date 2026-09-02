import type { BadgeElement, BadgeTemplateCreate } from "../model";

export function resizeBadgeDesign(
  design: BadgeTemplateCreate["design_data"],
  canvas: { width: number; height: number },
): BadgeTemplateCreate["design_data"] {
  const scaleX = canvas.width / design.canvas.width;
  const scaleY = canvas.height / design.canvas.height;
  const fontScale = Math.sqrt(scaleX * scaleY);

  return {
    ...design,
    canvas,
    elements: design.elements.map((element) => ({
      ...element,
      x: element.x * scaleX,
      y: element.y * scaleY,
      width: element.width * scaleX,
      height: element.height * scaleY,
      ...(element.type === "text"
        ? {
            paragraphs: element.paragraphs.map((paragraph) => ({
              ...paragraph,
              runs: paragraph.runs.map((run) => ({
                ...run,
                fontSize: Math.max(6, Math.round(run.fontSize * fontScale)),
              })),
            })),
          }
        : {}),
    })) as BadgeElement[],
  };
}
