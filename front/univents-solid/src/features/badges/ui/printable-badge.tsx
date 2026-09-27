import type { JSX } from "@solidjs/web";
import QRCode from "qrcode";
import { Show, createMemo } from "solid-js";
import { badgeDesignSchema, type BadgePrintItem } from "../model";
import { BadgePreview } from "./BadgePreview";

export function PrintableBadge(props: {
  badge: BadgePrintItem;
  participantName: string;
  location: string;
}): JSX.Element {
  const parsed = createMemo(() => badgeDesignSchema.safeParse(props.badge.design_data));
  const design = () => (parsed().success ? parsed().data : null);
  const widthMm = () => (design() ? design()!.canvas.width / (96 / 25.4) : 85);
  const heightMm = () => (design() ? design()!.canvas.height / (96 / 25.4) : 54);

  return (
    <article
      class="overflow-hidden bg-transparent text-black shadow-xs print:break-inside-avoid print:shadow-none"
      style={{
        width: `${widthMm()}mm`,
        height: `${heightMm()}mm`,
      }}
    >
      <BadgePreview
        badge={props.badge}
        participantName={props.participantName}
        location={props.location}
        class="relative size-full rounded-none border-0 shadow-none"
        style={{
          width: "100%",
          height: "100%",
          "aspect-ratio": "auto",
        }}
      />
    </article>
  );
}

export function PrintableQr(props: {
  badge: BadgePrintItem;
  size: number;
  participant: string;
}): JSX.Element {
  const qr = createMemo(() => {
    try {
      return QRCode.create(props.badge.action_url).modules;
    } catch {
      return null;
    }
  });
  const margin = 2;
  const viewSize = () => (qr() ? qr()!.size + margin * 2 : 30);

  const pathData = createMemo(() => {
    const matrix = qr();
    if (!matrix) return "";
    let d = "";
    for (let row = 0; row < matrix.size; row++) {
      for (let col = 0; col < matrix.size; col++) {
        if (matrix.get(row, col)) {
          d += `M${col + margin},${row + margin}h1v1h-1z `;
        }
      }
    }
    return d;
  });

  return (
    <article
      class="flex break-inside-avoid flex-col items-center gap-2 text-center text-black"
      style={{ width: `${props.size}px` }}
    >
      <Show when={qr()}>
        <svg
          role="img"
          aria-label={`QR Code de ${props.participant}`}
          viewBox={`0 0 ${viewSize()} ${viewSize()}`}
          style={{ width: `${props.size}px`, height: `${props.size}px` }}
          shape-rendering="crispEdges"
        >
          <rect width="100%" height="100%" fill="white" />
          <path d={pathData()} fill="black" />
        </svg>
      </Show>
      <strong class="max-w-full truncate text-sm">{props.participant}</strong>
    </article>
  );
}
