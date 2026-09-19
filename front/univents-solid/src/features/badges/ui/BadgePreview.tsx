import type { JSX } from "@solidjs/web";
import { type Accessor, For, Show, createMemo, createSignal, onCleanup } from "solid-js";
import { QrRenderer } from "@/features/editor/qr-renderer";
import { StaticText } from "@/features/editor/static-text";
import { DEFAULT_BADGE_TEMPLATE } from "../default-template";
import {
  type BadgeDesign,
  type BadgeElement,
  type BadgePrintItem,
  type BadgeProfileBadge,
  type BadgeTemplate,
  badgeDesignSchema,
} from "../model";

export interface BadgePreviewProps {
  badge: BadgeProfileBadge | BadgePrintItem | BadgeTemplate;
  class?: string;
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
  style?: JSX.CSSProperties;
}

const FALLBACK_LABELS: Record<string, string> = {
  event_name: "Nome do evento",
  edition_name: "Nome da edição",
  ticket_name: "Nome do ingresso",
  ticket_type: "Nome do ingresso",
  ticket: "Nome do ingresso",
  participant_name: "Nome do participante",
  name: "Nome do participante",
  location: "Local da edição",
  checkin_url: "Link de check-in",
};

function previewText(
  element: Extract<BadgeElement, { type: "text" }>,
  v: Record<string, string>,
  show?: boolean,
) {
  return {
    ...element,
    paragraphs: element.paragraphs.map((paragraph) => ({
      ...paragraph,
      runs: paragraph.runs.map((run) => ({
        ...run,
        text: run.text.replace(/\{\{([^}]+)\}\}/g, (_, rawKey: string) => {
          const key = rawKey.trim();
          if (v[key] !== undefined && v[key] !== "") return v[key];
          if (!show) return "";
          return FALLBACK_LABELS[key] ?? key;
        }),
      })),
    })),
  };
}

export function BadgePreview(props: BadgePreviewProps): JSX.Element {
  const [containerWidth, setContainerWidth] = createSignal(0);
  let observer: ResizeObserver | undefined;

  onCleanup(() => {
    observer?.disconnect();
  });

  const setupContainer = (el: HTMLDivElement) => {
    const rect = el.getBoundingClientRect();
    if (rect.width > 0) {
      setContainerWidth(rect.width);
    }

    if (typeof ResizeObserver !== "undefined") {
      observer = new ResizeObserver(([entry]) => {
        if (entry && entry.contentRect.width > 0) {
          setContainerWidth(entry.contentRect.width);
        }
      });
      observer.observe(el);
    }
  };

  const design = createMemo<BadgeDesign>(() => {
    const raw = props.badge.design_data;
    const parsed = badgeDesignSchema.safeParse(raw);
    return parsed.success ? parsed.data : DEFAULT_BADGE_TEMPLATE.design_data;
  });

  const qrElement = createMemo(() =>
    design().elements.find((el) => el.type === "qr"),
  );

  const actionUrl = createMemo(() => {
    if (props.actionUrl) return props.actionUrl;
    if ("action_url" in props.badge && props.badge.action_url) {
      return props.badge.action_url;
    }
    const qr = qrElement();
    if (qr && qr.type === "qr" && qr.value) return qr.value;
    return "https://univents.app/check-in";
  });

  const values = createMemo<Record<string, string>>(() => {
    const b = props.badge;
    const ticketName =
      props.ticketName ??
      ("ticket_name" in b && b.ticket_name ? b.ticket_name : "");
    const eventName =
      props.eventName ??
      ("event_name" in b && b.event_name ? b.event_name : "");
    const editionName =
      props.editionName ??
      ("edition_name" in b && b.edition_name ? b.edition_name : "");
    const participantName = props.participantName ?? "";
    const location = props.location ?? "";
    const action = actionUrl();

    return {
      event_name: eventName,
      edition_name: editionName,
      ticket_name: ticketName,
      ticket_type: ticketName,
      ticket: ticketName,
      participant_name: participantName,
      name: participantName,
      location: location,
      checkin_url: action,
    };
  });

  const scale = createMemo(() => {
    const width = containerWidth();
    const canvasWidth = design().canvas.width;
    if (!width || !canvasWidth) return 1;
    return width / canvasWidth;
  });

  const framed = () => props.framed !== false;
  const cls = () => props.class ?? props.className ?? "relative h-40 w-auto";

  return (
    <div
      ref={setupContainer}
      class={`${cls()} shrink-0 overflow-hidden rounded-md${framed() ? " border shadow-xs" : ""}`}
      style={{
        "aspect-ratio": `${design().canvas.width} / ${design().canvas.height}`,
        position: "relative",
        ...(props.contain
          ? design().canvas.width > design().canvas.height
            ? { width: "100%", height: "auto" }
            : { width: "auto", height: "100%" }
          : {}),
        ...props.style,
      }}
    >
      <div
        class="absolute left-0 top-0 overflow-hidden select-none"
        style={{
          width: `${design().canvas.width}px`,
          height: `${design().canvas.height}px`,
          transform: `scale(${scale()})`,
          "transform-origin": "top left",
          "background-color": design().backgroundColor,
          "background-image": design().background
            ? `url(${design().background})`
            : undefined,
          "background-position": "center",
          "background-repeat": "no-repeat",
          "background-size": "cover",
        }}
      >
        <For each={design().elements}>
          {(element) => (
            <BadgeElementPreview
              element={element}
              values={values}
              showVariables={() => props.showVariables}
              actionUrl={actionUrl}
            />
          )}
        </For>
      </div>
    </div>
  );
}

function BadgeElementPreview(props: {
  element: BadgeElement;
  values: Accessor<Record<string, string>>;
  showVariables: Accessor<boolean | undefined>;
  actionUrl: Accessor<string>;
}): JSX.Element {
  const transformedElement = createMemo(() => {
    if (props.element.type === "text") {
      return previewText(
        props.element as Extract<BadgeElement, { type: "text" }>,
        props.values(),
        props.showVariables(),
      );
    }
    return props.element;
  });

  const style = createMemo(() => ({
    position: "absolute" as const,
    left: `${props.element.x}px`,
    top: `${props.element.y}px`,
    width: `${props.element.width}px`,
    height: `${props.element.height}px`,
  }));

  return (
    <Show
      when={props.element.type === "text"}
      fallback={
        <Show
          when={props.element.type === "image"}
          fallback={
            <div style={style()}>
              <QrRenderer
                value={props.actionUrl()}
                foreground={
                  props.element.type === "qr"
                    ? (props.element as Extract<BadgeElement, { type: "qr" }>).foreground
                    : "#000000"
                }
                background={
                  props.element.type === "qr"
                    ? (props.element as Extract<BadgeElement, { type: "qr" }>).background
                    : "#ffffff"
                }
                style={
                  props.element.type === "qr"
                    ? (props.element as Extract<BadgeElement, { type: "qr" }>).style
                    : "square"
                }
              />
            </div>
          }
        >
          {(() => {
            const img = props.element as Extract<BadgeElement, { type: "image" }>;
            return (
              <img
                src={img.src}
                alt=""
                class="absolute transition-transform duration-700 ease-out pointer-events-none"
                loading="lazy"
                decoding="async"
                style={{
                  ...style(),
                  "object-fit": img.fit,
                  opacity: img.opacity,
                  "border-radius": `${img.radius}px`,
                }}
              />
            );
          })()}
        </Show>
      }
    >
      <div class="absolute overflow-hidden" style={style()}>
        <StaticText
          paragraphs={
            (transformedElement() as Extract<BadgeElement, { type: "text" }>)
              .paragraphs
          }
          values={props.values()}
          showVariables={props.showVariables()}
        />
      </div>
    </Show>
  );
}
