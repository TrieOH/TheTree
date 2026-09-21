import type { JSX } from "@solidjs/web";
import {
  For,
  Show,
  createEffect,
  createMemo,
  createSignal,
} from "solid-js";
import MaximizeIcon from "~icons/lucide/maximize";
import MinusIcon from "~icons/lucide/minus";
import PlusIcon from "~icons/lucide/plus";
import {
  EditorElementFrame,
  EditableText,
  StaticText,
} from "@/features/editor";
import type { CertificationTemplateElement } from "../../model";
import type {
  HashCertificateElement,
  ImageCertificateElement,
  SignatureCertificateElement,
  TextCertificateElement,
} from "../types";
import { useElementSize } from "../hooks/use-element-size";
import { certificateEditorActions, certificateEditorStore } from "../store";
import { HashElementView } from "./elements/hash-element-view";
import { SignatureElementView } from "./elements/signature-element-view";

const Maximize = MaximizeIcon as unknown as (props: { class?: string }) => JSX.Element;
const Minus = MinusIcon as unknown as (props: { class?: string }) => JSX.Element;
const Plus = PlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export const CERTIFICATE_CANVAS_DOM_ID = "certificate-canvas-root";

const MIN_ZOOM = 0.25;
const MAX_ZOOM = 4;

const CHECKERBOARD_STYLE = {
  "background-color": "#e5e7eb",
  "background-image":
    "linear-gradient(45deg, #d1d5db 25%, transparent 25%), linear-gradient(-45deg, #d1d5db 25%, transparent 25%), linear-gradient(45deg, transparent 75%, #d1d5db 75%), linear-gradient(-45deg, transparent 75%, #d1d5db 75%)",
  "background-size": "16px 16px",
  "background-position": "0 0, 0 8px, 8px -8px, -8px 0",
};

interface CertificateCanvasElementItemProps {
  id: string;
  scale: () => number;
}

function CertificateCanvasElementItem(props: CertificateCanvasElementItemProps): JSX.Element {
  const elements = () =>
    certificateEditorStore.draft().design_data.elements ?? [];
  const element = createMemo<CertificationTemplateElement | undefined>(() =>
    elements().find((e) => e.id === props.id),
  );
  const elementIndex = createMemo(() =>
    elements().findIndex((e) => e.id === props.id),
  );
  const canvas = () => certificateEditorStore.canvas();
  const isSelected = createMemo(
    () => certificateEditorStore.selectedElementId() === props.id,
  );
  const isEditing = createMemo(
    () => certificateEditorStore.editingElementId() === props.id,
  );

  return (
    <Show when={element()}>
      {(el) => (
        <EditorElementFrame
          bounds={el()}
          scale={props.scale()}
          canvas={canvas()}
          zIndex={elementIndex() + 1}
          selected={isSelected()}
          editing={isEditing()}
          onSelect={() => {
            const currentEditing = certificateEditorStore.editingElementId();
            if (currentEditing && currentEditing !== props.id) {
              certificateEditorActions.stopEditing();
            }
            certificateEditorActions.selectElement(props.id);
          }}
          onDoubleClick={() => {
            if (el().type === "text") {
              certificateEditorActions.startEditing(props.id);
            }
          }}
          onChangeBounds={(bounds) =>
            certificateEditorActions.updateElementBounds(props.id, bounds)
          }
          onDelete={
            el().type !== "hash"
              ? () => certificateEditorActions.removeElement(props.id)
              : undefined
          }
        >
          <Show
            when={el().type === "text"}
            fallback={
              <Show
                when={el().type === "image"}
                fallback={
                  <Show
                    when={el().type === "signature"}
                    fallback={
                      <HashElementView element={el() as HashCertificateElement} />
                    }
                  >
                    <SignatureElementView
                      element={el() as SignatureCertificateElement}
                    />
                  </Show>
                }
              >
                {(() => {
                  const img = () => el() as ImageCertificateElement;
                  return (
                    <img
                      src={img().src}
                      alt=""
                      class="h-full w-full pointer-events-none select-none"
                      style={{
                        "object-fit": img().fit,
                        opacity: img().opacity,
                        "border-radius": `${img().radius}px`,
                      }}
                      draggable={false}
                    />
                  );
                })()}
              </Show>
            }
          >
            {(() => {
              const txt = () => el() as TextCertificateElement;

              return (
                <div
                  class="relative h-full w-full cursor-text"
                  onDblClick={(e) => {
                    e.stopPropagation();
                    certificateEditorActions.startEditing(props.id);
                  }}
                >
                  <Show
                    when={isEditing()}
                    fallback={
                      <StaticText
                        paragraphs={txt().paragraphs}
                        showVariables={false}
                      />
                    }
                  >
                    <EditableText
                      element={txt()}
                      adapter={certificateEditorActions.textAdapter}
                      onStopEditing={() => certificateEditorActions.stopEditing()}
                    />
                  </Show>
                </div>
              );
            })()}
          </Show>
        </EditorElementFrame>
      )}
    </Show>
  );
}

export function CertificateCanvas(): JSX.Element {
  const elements = () =>
    certificateEditorStore.draft().design_data.elements ?? [];
  const elementIds = createMemo(() => elements().map((e) => e.id));
  const backgroundUrl = () =>
    certificateEditorStore.draft().design_data.background;
  const canvas = () => certificateEditorStore.canvas();

  const { ref: stageRef, size: viewport } = useElementSize<HTMLDivElement>();
  const [zoom, setZoom] = createSignal(1);
  const [pan, setPan] = createSignal({ x: 0, y: 0 });
  const [panning, setPanning] = createSignal(false);
  let panStart: { pointerX: number; pointerY: number; x: number; y: number } | null = null;
  let lastFitScale = 1;

  const fitScale = createMemo(() => {
    const vp = viewport();
    const c = canvas();
    if (vp.width <= 0 || vp.height <= 0) {
      return lastFitScale;
    }

    const padding = 96;
    const availableWidth = Math.max(120, vp.width - padding);
    const availableHeight = Math.max(120, vp.height - padding);
    const nextScale = Math.min(
      availableWidth / c.width,
      availableHeight / c.height,
      2,
    );

    if (Number.isFinite(nextScale) && nextScale > 0) {
      lastFitScale = nextScale;
    }
    return lastFitScale;
  });

  const scale = () => fitScale() * zoom();

  let stageEl: HTMLDivElement | undefined;

  createEffect(
    () => stageEl,
    (stage) => {
      if (!stage) return;

      const handleWheel = (event: WheelEvent) => {
        event.preventDefault();
        const factor = event.deltaY < 0 ? 1.1 : 0.9;
        setZoom((value) =>
          Math.min(MAX_ZOOM, Math.max(MIN_ZOOM, value * factor)),
        );
      };

      stage.addEventListener("wheel", handleWheel, { passive: false });
      return () => {
        stage.removeEventListener("wheel", handleWheel);
      };
    },
  );

  function startPan(event: PointerEvent) {
    const target = event.target as HTMLElement;
    if (target.closest("[data-editor-element], [data-canvas-controls]")) {
      return;
    }

    certificateEditorActions.selectElement(null);
    certificateEditorActions.stopEditing();
    panStart = {
      pointerX: event.clientX,
      pointerY: event.clientY,
      ...pan(),
    };
    setPanning(true);
    (event.currentTarget as HTMLElement).setPointerCapture(event.pointerId);
  }

  function finishPan(event: PointerEvent) {
    panStart = null;
    setPanning(false);
    const el = event.currentTarget as HTMLElement;
    if (el.hasPointerCapture(event.pointerId)) {
      el.releasePointerCapture(event.pointerId);
    }
  }

  return (
    <div
      ref={(el) => {
        stageEl = el;
        stageRef(el);
      }}
      class={
        "relative flex min-h-0 min-w-0 flex-1 touch-none items-center justify-center overflow-hidden p-12 " +
        (panning() ? "cursor-grabbing" : "cursor-grab")
      }
      style={CHECKERBOARD_STYLE}
      onPointerDown={startPan}
      onPointerMove={(event) => {
        if (!panStart) return;
        setPan({
          x: panStart.x + event.clientX - panStart.pointerX,
          y: panStart.y + event.clientY - panStart.pointerY,
        });
      }}
      onPointerUp={finishPan}
      onPointerCancel={finishPan}
    >
      <div
        data-canvas-controls="true"
        class="absolute bottom-4 left-4 z-20 flex items-center gap-1 rounded-lg border border-border bg-popover/90 p-1 shadow-md backdrop-blur-xs"
      >
        <button
          type="button"
          aria-label="Diminuir zoom"
          class="rounded p-1.5 hover:bg-muted"
          onClick={() => setZoom((value) => Math.max(MIN_ZOOM, value - 0.1))}
        >
          <Minus class="size-4" />
        </button>
        <span class="min-w-12 text-center text-xs font-medium">
          {Math.round(scale() * 100)}%
        </span>
        <button
          type="button"
          aria-label="Aumentar zoom"
          class="rounded p-1.5 hover:bg-muted"
          onClick={() => setZoom((value) => Math.min(MAX_ZOOM, value + 0.1))}
        >
          <Plus class="size-4" />
        </button>
        <div class="mx-1 h-4 w-px bg-border" />
        <button
          type="button"
          aria-label="Ajustar zoom na tela"
          class="rounded p-1.5 hover:bg-muted"
          onClick={() => {
            setZoom(1);
            setPan({ x: 0, y: 0 });
          }}
        >
          <Maximize class="size-4" />
        </button>
      </div>

      <div
        class="relative shrink-0 shadow-2xl ring-1 ring-black/5"
        style={{
          width: `${canvas().width * scale()}px`,
          height: `${canvas().height * scale()}px`,
          transform: `translate(${pan().x}px, ${pan().y}px)`,
        }}
      >
        <div
          id={CERTIFICATE_CANVAS_DOM_ID}
          class="relative overflow-hidden"
          style={{
            width: `${canvas().width}px`,
            height: `${canvas().height}px`,
            transform: `scale(${scale()})`,
            "transform-origin": "top left",
            "background-color": "#ffffff",
            ...(backgroundUrl()
              ? {
                  "background-image": `url(${backgroundUrl()})`,
                  "background-size": "cover",
                  "background-position": "center",
                  "background-repeat": "no-repeat",
                }
              : {}),
          }}
        >
          <For each={elementIds()}>
            {(id) => (
              <CertificateCanvasElementItem
                id={id}
                scale={scale}
              />
            )}
          </For>
        </div>
      </div>
    </div>
  );
}
