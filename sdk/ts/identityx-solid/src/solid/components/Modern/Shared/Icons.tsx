import { omit } from "solid-js";
import type { JSX } from "@solidjs/web";

export function createIcon(
  iconNode: [string, Record<string, string>][],
) {
  return function Icon(props: JSX.SvgSVGAttributes<SVGSVGElement> & { size?: number | string }) {
    const rest = omit(props, "size", "class");
    const local = { size: props.size, class: props.class };
    return (
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width={local.size ?? 24}
        height={local.size ?? 24}
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
        class={local.class}
        {...rest}
      >
        {iconNode.map(([tag, attrs]) => {
          if (tag === "circle") return <circle {...attrs} />;
          if (tag === "path") return <path {...attrs} />;
          if (tag === "line") return <line {...attrs} />;
          if (tag === "polyline") return <polyline {...attrs} />;
          if (tag === "rect") return <rect {...attrs} />;
          return null;
        })}
      </svg>
    );
  };
}

export const ArrowRight = createIcon([["path", { d: "M5 12h14" }], ["path", { d: "m12 5 7 7-7 7" }]]);
export const Loader2 = createIcon([["path", { d: "M21 12a9 9 0 1 1-6.219-8.56" }]]);
export const CheckCircle2 = createIcon([["circle", { cx: "12", cy: "12", r: "10" }], ["path", { d: "m9 12 2 2 4-4" }]]);
export const XCircle = createIcon([["circle", { cx: "12", cy: "12", r: "10" }], ["path", { d: "m15 9-6 6" }], ["path", { d: "m9 9 6 6" }]]);
export const Sparkles = createIcon([["path", { d: "m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z" }], ["path", { d: "M5 3v4" }], ["path", { d: "M19 17v4" }], ["path", { d: "M3 5h4" }], ["path", { d: "M17 19h4" }]]);
export const Copy = createIcon([["rect", { width: "14", height: "14", x: "8", y: "8", rx: "2", ry: "2" }], ["path", { d: "M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2" }]]);
export const Check = createIcon([["path", { d: "M20 6 9 17l-5-5" }]]);
export const User = createIcon([["path", { d: "M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2" }], ["circle", { cx: "12", cy: "7", r: "4" }]]);
export const FolderOpen = createIcon([["path", { d: "m6 14 1.45-2.9A2 2 0 0 1 9.24 10H20a2 2 0 0 1 1.94 2.5l-1.55 6a2 2 0 0 1-1.94 1.5H4a2 2 0 0 1-2-2V5c0-1.1.9-2 2-2h3.93a2 2 0 0 1 1.66.9l.82 1.2a2 2 0 0 0 1.66.9H18a2 2 0 0 1 2 2v2" }]]);
export const Key = createIcon([["path", { d: "m15.5 7.5 2.3 2.3a1 1 0 0 0 1.4 0l2.1-2.1a1 1 0 0 0 0-1.4L19 4" }], ["path", { d: "m21 2-9.6 9.6" }], ["circle", { cx: "7.5", cy: "15.5", r: "5.5" }]]);
export const Mail = createIcon([["rect", { width: "20", height: "16", x: "2", y: "4", rx: "2" }], ["path", { d: "m2 7 10 7 10-7" }]]);
export const Eye = createIcon([["path", { d: "M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z" }], ["circle", { cx: "12", cy: "12", r: "3" }]]);
export const EyeOff = createIcon([["path", { d: "M9.88 9.88a3 3 0 1 0 4.24 4.24" }], ["path", { d: "M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68" }], ["path", { d: "M6.61 6.61A13.526 13.526 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61" }], ["line", { x1: "2", y1: "2", x2: "22", y2: "22" }]]);
