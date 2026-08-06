import { ProfileQrCode } from "../../profile/ui/profile-qr-code";
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
  actionUrl: actionUrlOverride,
}: {
  badge: BadgeProfileBadge | BadgePrintItem | BadgeTemplate;
  className?: string;
  contain?: boolean;
  framed?: boolean;
  eventName?: string;
  editionName?: string;
  ticketName?: string;
  participantName?: string;
  actionUrl?: string;
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
    checkin_url: actionUrl,
  };

  function textFor(element: Extract<BadgeElement, { type: "text" }>) {
    return element.paragraphs
      .flatMap((paragraph) => paragraph.runs.map((run) => run.text))
      .join("\n")
      .replace(/\{\{([^}]+)\}\}/g, (_, key: string) => values[key] ?? "");
  }

  return (
    <div
      className={`${className} shrink-0 overflow-hidden rounded-md${framed ? " border shadow-sm" : ""}`}
      style={{
        aspectRatio: `${design.canvas.width} / ${design.canvas.height}`,
        backgroundColor: design.backgroundColor,
        backgroundImage: design.background
          ? `url(${design.background})`
          : undefined,
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
        backgroundSize: "cover",
        ...(contain
          ? design.canvas.width >= design.canvas.height
            ? { width: "100%", height: "auto" }
            : { width: "auto", height: "100%" }
          : {}),
      }}
    >
      {design.elements.map((element) => (
        <BadgeElementPreview
          key={element.id}
          element={element}
          design={design}
          actionUrl={actionUrl}
          text={element.type === "text" ? textFor(element) : undefined}
        />
      ))}
    </div>
  );
}

function BadgeElementPreview({
  element,
  design,
  actionUrl,
  text,
}: {
  element: BadgeElement;
  design: typeof DEFAULT_BADGE_TEMPLATE.design_data;
  actionUrl: string;
  text?: string;
}) {
  const style = {
    left: `${(element.x / design.canvas.width) * 100}%`,
    top: `${(element.y / design.canvas.height) * 100}%`,
    width: `${(element.width / design.canvas.width) * 100}%`,
    height: `${(element.height / design.canvas.height) * 100}%`,
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
        <ProfileQrCode value={actionUrl} size={160} className="size-full" />
      </div>
    );
  }
  return (
    <div
      className="absolute overflow-hidden whitespace-pre-wrap"
      style={{
        ...style,
        color: element.paragraphs[0]?.runs[0]?.color,
        fontWeight: element.paragraphs[0]?.runs[0]?.bold ? "bold" : "normal",
        fontSize: Math.max(
          10,
          (element.paragraphs[0]?.runs[0]?.fontSize ?? 20) * 0.5,
        ),
        textAlign: element.paragraphs[0]?.align,
      }}
    >
      {text}
    </div>
  );
}
