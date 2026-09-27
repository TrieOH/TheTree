import type { JSX } from "@solidjs/web";
import { For, Show, createMemo, createSignal, onCleanup } from "solid-js";
import { StaticText } from "@/features/editor/static-text";
import { DEFAULT_CERTIFICATION_TEMPLATE } from "../default-template";
import type {
  CertificationTemplateElement,
  CertificationTemplateI,
} from "../model";

export interface CertificatePreviewProps {
  template: CertificationTemplateI;
  class?: string;
  className?: string;
  contain?: boolean;
  framed?: boolean;
  showVariables?: boolean;
  variables?: Record<string, string | number>;
  style?: JSX.CSSProperties;
}

export function CertificatePreview(props: CertificatePreviewProps): JSX.Element {
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

  const canvasWidth = () =>
    props.template.design_data?.canvas?.width ??
    DEFAULT_CERTIFICATION_TEMPLATE.design_data.canvas?.width ??
    1484;

  const canvasHeight = () =>
    props.template.design_data?.canvas?.height ??
    DEFAULT_CERTIFICATION_TEMPLATE.design_data.canvas?.height ??
    1060;

  const scale = createMemo(() => {
    const w = containerWidth();
    if (w <= 0) return 1;
    return w / canvasWidth();
  });

  const backgroundUrl = () => props.template.design_data?.background;
  const elements = () => props.template.design_data?.elements ?? [];

  const values = createMemo<Record<string, string>>(() => {
    const base: Record<string, string> = {
      participant_name: "Nome do Participante",
      event_name: "Nome do Evento",
      edition_name: "Nome da Edição",
      activity_name: "Nome da Atividade",
      participation_type: "participante",
      location: "Local do evento",
      workload_hours: "8 horas",
      participation_date: "19/09/2026",
      certified_at: "19/09/2026",
      cert_hash: "UNIV-2026-CERT-SAMPLE",
      verify_url: "https://univents.app/verify/sample",
    };
    if (props.variables) {
      for (const [k, v] of Object.entries(props.variables)) {
        if (v !== undefined && v !== "") {
          base[k] = String(v);
        }
      }
    }
    return base;
  });

  const framed = () => props.framed !== false;
  const cls = () => props.class ?? props.className ?? "relative h-40 w-auto";

  return (
    <div
      ref={setupContainer}
      class={`${cls()} shrink-0 overflow-hidden rounded-md${framed() ? " border border-border/40 shadow-xs" : ""}`}
      style={{
        "aspect-ratio": `${canvasWidth()} / ${canvasHeight()}`,
        position: "relative",
        ...(props.contain
          ? canvasWidth() > canvasHeight()
            ? { width: "100%", height: "auto" }
            : { width: "auto", height: "100%" }
          : {}),
        "max-width": "100%",
        "max-height": "100%",
        ...props.style,
      }}
    >
      <div
        class="absolute left-0 top-0 select-none overflow-hidden"
        style={{
          width: `${canvasWidth()}px`,
          height: `${canvasHeight()}px`,
          transform: `scale(${scale()})`,
          "transform-origin": "top left",
          "background-color": "#ffffff",
          "background-image": backgroundUrl() ? `url(${JSON.stringify(backgroundUrl())})` : undefined,
          "background-size": "cover",
          "background-position": "center",
          "background-repeat": "no-repeat",
        }}
      >
        <For each={elements()}>
          {(element: CertificationTemplateElement) => (
            <div
              class="absolute overflow-hidden"
              style={{
                left: `${element.x}px`,
                top: `${element.y}px`,
                width: `${element.width}px`,
                height: `${element.height}px`,
              }}
            >
              <Show when={element.type === "text"}>
                <StaticText
                  paragraphs={element.type === "text" ? element.paragraphs : []}
                  values={values()}
                  showVariables={props.showVariables ?? true}
                />
              </Show>
              <Show when={element.type === "image"}>
                <img
                  src={element.type === "image" ? element.src : ""}
                  alt="Elemento"
                  class="h-full w-full object-contain pointer-events-none"
                />
              </Show>
              <Show when={element.type === "signature"}>
                <div class="flex h-full w-full flex-col items-center justify-end pointer-events-none">
                  <Show when={element.type === "signature" && (element.src || (element as { imageUrl?: string }).imageUrl)}>
                    <img
                      src={element.type === "signature" ? ((element as { imageUrl?: string }).imageUrl ?? element.src) : ""}
                      alt={element.type === "signature" ? element.name : ""}
                      class="max-h-full max-w-full object-contain"
                    />
                  </Show>
                </div>
              </Show>
            </div>
          )}
        </For>
      </div>
    </div>
  );
}
