import QRCode from "qrcode";
import type { BadgePrintItem } from "../model";
import { badgeDesignSchema } from "../model";
import { BadgePreview } from "./badge-preview";

export function PrintableBadge({
  badge,
  participantName,
  location,
}: {
  badge: BadgePrintItem;
  participantName: string;
  location: string;
}) {
  const parsed = badgeDesignSchema.safeParse(badge.design_data);
  const design = parsed.success ? parsed.data : null;
  const widthMm = design ? design.canvas.width / (96 / 25.4) : 85;
  const heightMm = design ? design.canvas.height / (96 / 25.4) : 54;

  return (
    <article
      className="overflow-hidden bg-transparent text-black shadow print:break-inside-avoid print:shadow-none"
      style={{ width: `${widthMm}mm`, height: `${heightMm}mm` }}
    >
      <BadgePreview
        badge={badge}
        participantName={participantName}
        location={location}
        className="relative h-full w-full rounded-none border-0 shadow-none"
        style={{ width: "100%", height: "100%", aspectRatio: "auto" }}
      />
    </article>
  );
}

export function PrintableQr({
  badge,
  size,
  participant,
}: {
  badge: BadgePrintItem;
  size: number;
  participant: string;
}) {
  const matrix = QRCode.create(badge.action_url).modules;
  const margin = 2;
  const viewSize = matrix.size + margin * 2;
  return (
    <article
      className="flex break-inside-avoid flex-col items-center gap-2 text-center text-black"
      style={{ width: size }}
    >
      <svg
        role="img"
        aria-label={`QR Code de ${participant}`}
        viewBox={`0 0 ${viewSize} ${viewSize}`}
        style={{ width: size, height: size }}
        shapeRendering="crispEdges"
      >
        <rect width="100%" height="100%" fill="white" />
        {Array.from({ length: matrix.size * matrix.size }, (_, index) => {
          const row = Math.floor(index / matrix.size);
          const column = index % matrix.size;
          return matrix.get(row, column) ? (
            <rect
              key={`${row}-${column}`}
              x={column + margin}
              y={row + margin}
              width="1"
              height="1"
            />
          ) : null;
        })}
      </svg>
      <strong className="max-w-full truncate text-sm">{participant}</strong>
    </article>
  );
}
