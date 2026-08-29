import QRCode from "qrcode";
import type { CSSProperties } from "react";
import { useElementSize } from "../../certifications/editor/hooks/use-element-size";
import { StaticTextElement } from "../../certifications/editor/ui/elements/text-element-view";
import { DEFAULT_BADGE_TEMPLATE } from "../default-template";
import {
  type BadgeElement,
  type BadgePrintItem,
  type BadgeProfileBadge,
  type BadgeTemplate,
  badgeDesignSchema,
} from "../model";

export function BadgePreview({
  badge,
  className = "relative h-40 w-auto",
  contain = false,
  framed = true,
  eventName,
  editionName,
  ticketName,
  participantName,
  location,
  actionUrl: actionUrlOverride,
  showVariables = false,
  style,
}: {
  badge: BadgeProfileBadge | BadgePrintItem | BadgeTemplate;
  className?: string;
  contain?: boolean;
  framed?: boolean;
  eventName?: string;
  editionName?: string;
  ticketName?: string;
  participantName?: string;
  location?: string;
  actionUrl?: string;
  showVariables?: boolean;
  style?: CSSProperties;
}) {
  const design = badgeDesignSchema.safeParse(badge.design_data).success
    ? badgeDesignSchema.parse(badge.design_data)
    : DEFAULT_BADGE_TEMPLATE.design_data;

  const qrElement = design.elements.find((element) => element.type === "qr");

  const qrValue =
    qrElement && typeof qrElement.value === "string"
      ? qrElement.value
      : undefined;

  const actionUrl =
    actionUrlOverride ??
    ("action_url" in badge ? badge.action_url : (qrValue ?? ""));

  const values: Record<string, string> = {
    event_name: eventName ?? ("event_name" in badge ? badge.event_name : ""),
    edition_name:
      editionName ?? ("edition_name" in badge ? badge.edition_name : ""),
    ticket_name:
      ticketName ?? ("ticket_name" in badge ? (badge.ticket_name ?? "") : ""),
    participant_name: participantName ?? "",
    location: location ?? "",
    checkin_url: actionUrl,
  };

  function previewText(element: Extract<BadgeElement, { type: "text" }>) {
    return {
      ...element,
      paragraphs: element.paragraphs.map((paragraph) => ({
        ...paragraph,
        runs: paragraph.runs.map((run) => ({
          ...run,
          text: run.text.replace(/\{\{([^}]+)\}\}/g, (_, key: string) => {
            if (values[key]) return values[key];
            if (!showVariables) return "";
            return (
              {
                event_name: "Nome do evento",
                edition_name: "Nome da edição",
                ticket_name: "Nome do ingresso",
                participant_name: "Nome do participante",
                location: "Local da edição",
                checkin_url: "Link de check-in",
              }[key] ?? key
            );
          }),
        })),
      })),
    };
  }

  const { ref, size } = useElementSize<HTMLDivElement>();
  const scale = size.width ? size.width / design.canvas.width : 1;

  return (
    <div
      ref={ref}
      className={`${className} shrink-0 overflow-hidden rounded-md${framed ? " border shadow-sm" : ""}`}
      style={{
        aspectRatio: `${design.canvas.width} / ${design.canvas.height}`,
        position: "relative",
        ...(contain
          ? design.canvas.width > design.canvas.height
            ? { width: "100%", height: "auto" }
            : { width: "auto", height: "100%" }
          : {}),
        ...style,
      }}
    >
      <div
        className="absolute left-0 top-0 overflow-hidden"
        style={{
          width: design.canvas.width,
          height: design.canvas.height,
          transform: `scale(${scale})`,
          transformOrigin: "top left",
          backgroundColor: design.backgroundColor,
          backgroundImage: design.background
            ? `url(${design.background})`
            : undefined,
          backgroundPosition: "center",
          backgroundRepeat: "no-repeat",
          backgroundSize: "cover",
        }}
      >
        {design.elements.map((element) => (
          <BadgeElementPreview
            key={element.id}
            element={element.type === "text" ? previewText(element) : element}
            actionUrl={actionUrl}
          />
        ))}
      </div>
    </div>
  );
}

function BadgeElementPreview({
  element,
  actionUrl,
}: {
  element: BadgeElement;
  actionUrl: string;
}) {
  const style = {
    left: element.x,
    top: element.y,
    width: element.width,
    height: element.height,
  };

  if (element.type === "image") {
    return (
      <img
        src={element.src}
        alt=""
        className="absolute transition-transform duration-700 ease-out"
        loading="lazy"
        decoding="async"
        style={{
          ...style,
          objectFit: element.fit,
          opacity: element.opacity,
          borderRadius: element.radius,
        }}
      />
    );
  }
  if (element.type === "qr") {
    return (
      <div className="absolute" style={style}>
        <PreviewQrCode element={element} value={actionUrl} />
      </div>
    );
  }
  return (
    <div className="absolute overflow-hidden" style={style}>
      <StaticTextElement element={element} />
    </div>
  );
}

function PreviewQrCode({
  element,
  value,
}: {
  element: Extract<BadgeElement, { type: "qr" }>;
  value: string;
}) {
  const matrix = QRCode.create(
    value || "https://univents.app/check-in",
  ).modules;
  const margin = 2;
  const size = matrix.size + margin * 2;
  const modules = [];

  for (let row = 0; row < matrix.size; row++) {
    for (let column = 0; column < matrix.size; column++) {
      if (!matrix.get(row, column)) continue;
      const x = column + margin;
      const y = row + margin;
      modules.push(
        element.style === "dots" ? (
          <circle
            key={`${row}-${column}`}
            cx={x + 0.5}
            cy={y + 0.5}
            r="0.42"
            fill={element.foreground}
          />
        ) : (
          <rect
            key={`${row}-${column}`}
            x={x}
            y={y}
            width="1"
            height="1"
            rx={element.style === "rounded" ? 0.28 : 0}
            fill={element.foreground}
          />
        ),
      );
    }
  }

  return (
    <svg
      role="img"
      aria-label="QR Code de check-in"
      className="size-full"
      viewBox={`0 0 ${size} ${size}`}
      style={{ background: element.background }}
      shapeRendering={
        element.style === "square" ? "crispEdges" : "geometricPrecision"
      }
    >
      {modules}
    </svg>
  );
}
