import {
  DEFAULT_EDITOR_FONT,
  DEFAULT_EDITOR_LINE_HEIGHT,
  DEFAULT_EDITOR_TEXT_COLOR,
  type RichParagraph,
  type RichRun,
} from "./types";

const convertedColorCache = new Map<string, string>();

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}

function runToStyle(run: RichRun): string {
  return [
    `color:${run.color}`,
    `font-size:${run.fontSize}px`,
    `font-family:${run.fontFamily}`,
    `font-weight:${run.bold ? 700 : 400}`,
    `font-style:${run.italic ? "italic" : "normal"}`,
    `text-decoration:${run.underline ? "underline" : "none"}`,
  ].join(";");
}

export function paragraphsToHtml(
  paragraphs: readonly RichParagraph[],
): string {
  if (!paragraphs.length) return "<p style=\"margin:0;min-height:1.25em\"><br></p>";
  return paragraphs
    .map((paragraph, paragraphIndex) => {
      const nearestRun =
        paragraph.runs[0] ??
        paragraphs
          .slice(0, paragraphIndex)
          .reverse()
          .find((candidate) => candidate.runs.length > 0)
          ?.runs.at(-1) ??
        paragraphs
          .slice(paragraphIndex + 1)
          .find((candidate) => candidate.runs.length > 0)?.runs[0];
      const content =
        paragraph.runs
          .map(
            (run) =>
              `<span style="${runToStyle(run)}">${escapeHtml(run.text) || "\u200b"}</span>`,
          )
          .join("") ||
        (nearestRun
          ? `<span style="${runToStyle(nearestRun)}">\u200b</span>`
          : "\u200b");
      return `<p style="text-align:${paragraph.align};line-height:${paragraph.lineHeight};margin:0;min-height:1.25em">${content}</p>`;
    })
    .join("");
}

function channelsToHex(channels: number[]): string {
  return `#${channels
    .map((channel) =>
      Math.max(0, Math.min(255, Math.round(channel)))
        .toString(16)
        .padStart(2, "0"),
    )
    .join("")}`;
}

export function computedColorToHex(color: string): string {
  if (!color) return DEFAULT_EDITOR_TEXT_COLOR;
  if (color.startsWith("#")) return color;
  const cached = convertedColorCache.get(color);
  if (cached) return cached;

  const rgb = color.match(
    /^rgba?\(\s*([\d.]+)(?:\s*,\s*|\s+)([\d.]+)(?:\s*,\s*|\s+)([\d.]+)/i,
  );
  if (rgb) {
    const converted = channelsToHex([
      Number(rgb[1]),
      Number(rgb[2]),
      Number(rgb[3]),
    ]);
    convertedColorCache.set(color, converted);
    return converted;
  }

  if (typeof document === "undefined") return DEFAULT_EDITOR_TEXT_COLOR;

  try {
    const canvas = document.createElement("canvas");
    canvas.width = 1;
    canvas.height = 1;
    const context = canvas.getContext("2d");
    if (!context) return DEFAULT_EDITOR_TEXT_COLOR;
    context.clearRect(0, 0, 1, 1);
    context.fillStyle = color;
    context.fillRect(0, 0, 1, 1);
    const pixel = context.getImageData(0, 0, 1, 1).data;
    const converted = channelsToHex([pixel[0] ?? 0, pixel[1] ?? 0, pixel[2] ?? 0]);
    convertedColorCache.set(color, converted);
    return converted;
  } catch {
    return DEFAULT_EDITOR_TEXT_COLOR;
  }
}

function readRunStyle(element: Element): Omit<RichRun, "text"> {
  if (typeof window === "undefined") {
    return {
      bold: false,
      italic: false,
      underline: false,
      color: DEFAULT_EDITOR_TEXT_COLOR,
      fontSize: 16,
      fontFamily: DEFAULT_EDITOR_FONT,
    };
  }

  const style = window.getComputedStyle(element);
  const numericWeight = Number(style.fontWeight);
  return {
    bold:
      style.fontWeight === "bold" ||
      (!Number.isNaN(numericWeight) && numericWeight >= 600),
    italic: style.fontStyle === "italic",
    underline: (style.textDecorationLine || style.textDecoration || "").includes("underline"),
    color: computedColorToHex(style.color),
    fontSize: Math.round(Number.parseFloat(style.fontSize) || 16),
    fontFamily: style.fontFamily || DEFAULT_EDITOR_FONT,
  };
}

function collectRuns(root: Node): RichRun[] {
  const runs: RichRun[] = [];

  function visit(node: Node) {
    if (node.nodeType === Node.TEXT_NODE) {
      const text = node.textContent?.replace(/\u200b/g, "") ?? "";
      if (!text) return;
      const style = node.parentElement
        ? readRunStyle(node.parentElement)
        : {
            bold: false,
            italic: false,
            underline: false,
            color: DEFAULT_EDITOR_TEXT_COLOR,
            fontSize: 16,
            fontFamily: DEFAULT_EDITOR_FONT,
          };
      const previous = runs[runs.length - 1];
      if (
        previous &&
        previous.bold === style.bold &&
        previous.italic === style.italic &&
        previous.underline === style.underline &&
        previous.color === style.color &&
        previous.fontSize === style.fontSize &&
        previous.fontFamily === style.fontFamily
      ) {
        previous.text += text;
      } else {
        runs.push({ ...style, text });
      }
      return;
    }

    node.childNodes.forEach(visit);
  }

  visit(root);
  return runs;
}

function paragraphAlign(element: Element): RichParagraph["align"] {
  if (typeof window === "undefined") return "left";
  const align = window.getComputedStyle(element).textAlign;
  return align === "center" || align === "right" || align === "justify"
    ? align
    : "left";
}

export function domToParagraphs(container: HTMLElement): RichParagraph[] {
  const blocks = Array.from(container.children).filter(
    (child) => child.tagName === "DIV" || child.tagName === "P",
  );
  if (blocks.length === 0) {
    const runs = collectRuns(container);
    return [
      {
        align: "left",
        lineHeight: DEFAULT_EDITOR_LINE_HEIGHT,
        runs: runs.length
          ? runs
          : [
              {
                text: "",
                bold: false,
                italic: false,
                underline: false,
                color: DEFAULT_EDITOR_TEXT_COLOR,
                fontSize: 16,
                fontFamily: DEFAULT_EDITOR_FONT,
              },
            ],
      },
    ];
  }

  return blocks.map((block) => {
    const runs = collectRuns(block);
    const style = window.getComputedStyle(block);
    const fontSize = Number.parseFloat(style.fontSize) || 16;
    const lineHeightPx = Number.parseFloat(style.lineHeight);
    const lineHeight = Number.isNaN(lineHeightPx)
      ? DEFAULT_EDITOR_LINE_HEIGHT
      : lineHeightPx / fontSize || DEFAULT_EDITOR_LINE_HEIGHT;

    return {
      align: paragraphAlign(block),
      lineHeight,
      runs: runs.length
        ? runs
        : [
            {
              text: "",
              bold: false,
              italic: false,
              underline: false,
              color: DEFAULT_EDITOR_TEXT_COLOR,
              fontSize: 16,
              fontFamily: DEFAULT_EDITOR_FONT,
            },
          ],
    };
  });
}
