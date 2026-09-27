import type { JSX } from "@solidjs/web";
import QRCode from "qrcode";
import { createMemo } from "solid-js";

export interface QrRendererProps {
  value: string;
  foreground: string;
  background: string;
  style: "square" | "rounded" | "dots";
}

export function QrRenderer(props: QrRendererProps): JSX.Element {
  const matrix = createMemo(() => {
    try {
      return QRCode.create(props.value || "https://univents.app/check-in").modules;
    } catch {
      return null;
    }
  });

  const margin = 2;
  const size = createMemo(() => (matrix() ? matrix()!.size + margin * 2 : 33));

  const pathData = createMemo(() => {
    const m = matrix();
    if (!m) return "";
    const style = props.style;
    const parts: string[] = [];
    const r = style === "rounded" ? 0.28 : 0;
    for (let row = 0; row < m.size; row++) {
      for (let column = 0; column < m.size; column++) {
        if (!m.get(row, column)) continue;
        const x = column + margin;
        const y = row + margin;
        if (style === "dots") {
          const cx = x + 0.5;
          const cy = y + 0.5;
          const cr = 0.42;
          parts.push(`M ${cx - cr} ${cy} a ${cr} ${cr} 0 1 0 ${cr * 2} 0 a ${cr} ${cr} 0 1 0 -${cr * 2} 0`);
        } else if (style === "rounded") {
          parts.push(
            `M ${x + r} ${y} h ${1 - 2 * r} a ${r} ${r} 0 0 1 ${r} ${r} v ${1 - 2 * r} a ${r} ${r} 0 0 1 -${r} ${r} h -${1 - 2 * r} a ${r} ${r} 0 0 1 -${r} -${r} v -${1 - 2 * r} a ${r} ${r} 0 0 1 ${r} -${r} Z`,
          );
        } else {
          parts.push(`M ${x} ${y} h 1 v 1 h -1 Z`);
        }
      }
    }
    return parts.join(" ");
  });

  return (
    <svg
      role="img"
      aria-label="QR Code de check-in"
      class="size-full"
      viewBox={`0 0 ${size()} ${size()}`}
      style={{ background: props.background }}
      shape-rendering={
        props.style === "square" ? "crispEdges" : "geometricPrecision"
      }
    >
      <path d={pathData()} fill={props.foreground} />
    </svg>
  );
}
