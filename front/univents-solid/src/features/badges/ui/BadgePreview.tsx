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
  event_name: "Nome do Evento",
  edition_name: "Nome da Edição",
  ticket_name: "Nome do Ingresso",
  ticket_type: "Nome do Ingresso",
  ticket: "Nome do Ingresso",
  participant_name: "Nome do Participante",
  name: "Nome do Participante",
  activity_name: "Nome da Atividade",
  program_name: "Nome da Programação",
  participation_type: "Participação",
  location: "Local da Edição",
  location_name: "Local da Edição",
  workload_hours: "10 horas",
  participation_date: "19/09/2026",
  certified_at: "19/09/2026",
  issue_date: "19/09/2026",
  cert_hash: "UNIV-2026-CERT-SAMPLE",
  verify_url: "https://univents.app/verify/sample",
  checkin_url: "https://univents.app/check-in",
};

function previewText(
  element: Extract<BadgeElement, { type: "text" }>,
  v: Record<string, string>,
  show: boolean = true,
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
          if (key === "ticket_name" || key === "ticket_type" || key === "ticket") {
            const ticketVal = v.ticket_name ?? v.ticket_type ?? v.ticket;
            if (ticketVal) return ticketVal;
          }
          if (key === "participant_name" || key === "name") {
            const nameVal = v.participant_name ?? v.name;
            if (nameVal) return nameVal;
          }
          if (key === "activity_name" || key === "program_name") {
            const actVal = v.activity_name ?? v.program_name;
            if (actVal) return actVal;
          }
          if (key === "location" || key === "location_name") {
            const locVal = v.location ?? v.location_name;
            if (locVal) return locVal;
          }
          if (!show) return "";
          return FALLBACK_LABELS[key] ?? `{{${key}}}`;
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
              showVariables={() => props.showVariables ?? true}
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
  const el = () => props.element;

  return (
    <div
      class="absolute overflow-hidden"
      style={{
        left: `${el().x}px`,
        top: `${el().y}px`,
        width: `${el().width}px`,
        height: `${el().height}px`,
      }}
    >
      <Show when={el().type === "text"}>
        <StaticText
          paragraphs={
            previewText(
              el() as Extract<BadgeElement, { type: "text" }>,
              props.values(),
              props.showVariables() ?? true,
            ).paragraphs
          }
          values={props.values()}
          showVariables={props.showVariables() ?? true}
        />
      </Show>

      <Show when={el().type === "image"}>
        <img
          src={(el() as Extract<BadgeElement, { type: "image" }>).src}
          alt="Elemento"
          class="h-full w-full object-contain pointer-events-none"
        />
      </Show>

      <Show when={el().type === "qr"}>
        <div class="flex h-full w-full items-center justify-center pointer-events-none">
          <QrRenderer
            value={props.actionUrl()}
            foreground={
              (el() as Extract<BadgeElement, { type: "qr" }>).foreground ??
              "#000000"
            }
            background={
              (el() as Extract<BadgeElement, { type: "qr" }>).background ??
              "#ffffff"
            }
            style={
              (el() as Extract<BadgeElement, { type: "qr" }>).style ?? "square"
            }
          />
        </div>
      </Show>
    </div>
  );
}
