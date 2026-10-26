import { resolveStorageUrl } from "@/shared/lib/storage-url";
import type { JSX } from "@solidjs/web";
import {
  For,
  Show,
  createMemo,
  createSignal,
  onCleanup,
} from "solid-js";
import { CanvasControls } from "@trieoh/ui-solid";
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
                      src={resolveStorageUrl(img().src)}
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
  let wheelCleanup: (() => void) | undefined;

  const fitScale = createMemo(() => {
    const vp = viewport();
    const c = canvas();
    if (vp.width <= 0 || vp.height <= 0) {
      return lastFitScale;
    }

    const availableWidth = Math.max(120, vp.width - 24);
    const availableHeight = Math.max(120, vp.height - 24);
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

  onCleanup(() => {
    wheelCleanup?.();
  });

  function startPan(event: PointerEvent) {
    const target = event.target as HTMLElement;
    if (target.closest("[data-editor-element], [data-canvas-controls], button, input")) {
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
        stageRef(el);
        wheelCleanup?.();
        if (!el) return;

        const handleWheel = (event: WheelEvent) => {
          event.preventDefault();
          const delta =
            event.deltaMode === 1
              ? event.deltaY * 20
              : event.deltaMode === 2
                ? event.deltaY * 100
                : event.deltaY;
          const factor = Math.exp(-delta * 0.0015);
          const clampedFactor = Math.min(1.2, Math.max(0.8, factor));
          setZoom((value) => {
            const next = Math.min(MAX_ZOOM, Math.max(MIN_ZOOM, value * clampedFactor));
            return Number(next.toFixed(3));
          });
        };

        el.addEventListener("wheel", handleWheel, { passive: false });
        wheelCleanup = () => {
          el.removeEventListener("wheel", handleWheel);
        };
      }}
      class={
        "relative flex min-h-0 min-w-0 flex-1 touch-none items-center justify-center overflow-hidden p-8 sm:p-12 " +
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
      <CanvasControls
        scale={scale()}
        onZoomIn={() =>
          setZoom((value) =>
            Math.min(MAX_ZOOM, Number((value * 1.1).toFixed(3))),
          )
        }
        onZoomOut={() =>
          setZoom((value) =>
            Math.max(MIN_ZOOM, Number((value * 0.9).toFixed(3))),
          )
        }
        onReset={() => {
          setZoom(1);
          setPan({ x: 0, y: 0 });
        }}
      />

      <div
        class="relative shrink-0 shadow-2xl ring-1 ring-black/5 m-auto"
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
                  "background-image": `url(${resolveStorageUrl(backgroundUrl())})`,
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
