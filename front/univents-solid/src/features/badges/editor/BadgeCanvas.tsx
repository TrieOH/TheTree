import { resolveStorageUrl } from "@/shared/lib/storage-url";
import type { JSX } from "@solidjs/web";
import {
  For,
  Show,
  createEffect,
  createMemo,
  createSignal,
  onCleanup,
} from "solid-js";
import { CanvasControls } from "@trieoh/ui-solid";
import {
  EditableText,
  EditorElementFrame,
  QrRenderer,
  StaticText,
  type ElementBounds,
  type TextElementAdapter,
} from "@/features/editor";
import type { BadgeDesign, BadgeElement } from "../model";

export interface BadgeCanvasProps {
  design: BadgeDesign;
  selectedId: string | null;
  onSelect: (id: string | null) => void;
  onChangeElement: (id: string, element: Partial<BadgeElement>) => void;
  onDeleteElement?: (id: string) => void;
  textAdapter?: TextElementAdapter;
  previewValues?: Record<string, string>;
}

const MIN_ZOOM = 0.25;
const MAX_ZOOM = 4;

export function BadgeCanvas(props: BadgeCanvasProps): JSX.Element {
  let containerRef: HTMLDivElement | undefined;
  let ro: ResizeObserver | undefined;
  let wheelCleanup: (() => void) | undefined;
  const [zoom, setZoom] = createSignal(1);
  const [fitScale, setFitScale] = createSignal(1);
  const [pan, setPan] = createSignal({ x: 0, y: 0 });
  const [panning, setPanning] = createSignal(false);
  let panStart: { pointerX: number; pointerY: number; x: number; y: number } | null = null;
  const [editingId, setEditingId] = createSignal<string | null>(null);

  const elementIds = createMemo(() => props.design.elements.map((e) => e.id));

  const updateFitScale = () => {
    if (!containerRef) return;
    const { clientWidth, clientHeight } = containerRef;
    if (clientWidth <= 0 || clientHeight <= 0) return;
    const padding = 64;
    const availW = Math.max(100, clientWidth - padding);
    const availH = Math.max(100, clientHeight - padding);
    const scaleW = availW / props.design.canvas.width;
    const scaleH = availH / props.design.canvas.height;
    const nextScale = Math.min(1.5, Math.max(0.2, Math.min(scaleW, scaleH)));
    if (Number.isFinite(nextScale) && nextScale > 0) {
      setFitScale(nextScale);
    }
  };

  createEffect(
    () => ({
      width: props.design.canvas.width,
      height: props.design.canvas.height,
    }),
    () => {
      updateFitScale();
    },
  );

  onCleanup(() => {
    ro?.disconnect();
    wheelCleanup?.();
  });

  const effectiveScale = () => fitScale() * zoom();

  function startPan(event: PointerEvent) {
    const target = event.target as HTMLElement;
    if (target.closest("[data-editor-element], [data-canvas-controls], button, input")) {
      return;
    }
    props.onSelect(null);
    setEditingId(null);
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
        containerRef = el;
        wheelCleanup?.();
        if (!el) return;
        updateFitScale();
        ro?.disconnect();
        ro = new ResizeObserver(() => updateFitScale());
        ro.observe(el);

        const handleWheel = (event: WheelEvent) => {
          event.preventDefault();
          const factor = event.deltaY < 0 ? 1.1 : 0.9;
          setZoom((value) =>
            Math.min(MAX_ZOOM, Math.max(MIN_ZOOM, Number((value * factor).toFixed(3)))),
          );
        };
        el.addEventListener("wheel", handleWheel, { passive: false });
        wheelCleanup = () => {
          el.removeEventListener("wheel", handleWheel);
        };
      }}
      class={
        "relative flex min-h-0 min-w-0 flex-1 touch-none items-center justify-center overflow-hidden p-8 select-none " +
        (panning() ? "cursor-grabbing" : "cursor-grab")
      }
      style={{
        background:
          "radial-gradient(circle, var(--color-muted) 1px, transparent 1px)",
        "background-size": "16px 16px",
      }}
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
        scale={effectiveScale()}
        onZoomIn={() =>
          setZoom((v) => Math.min(MAX_ZOOM, Number((v * 1.1).toFixed(3))))
        }
        onZoomOut={() =>
          setZoom((v) => Math.max(MIN_ZOOM, Number((v * 0.9).toFixed(3))))
        }
        onReset={() => {
          setZoom(1);
          setPan({ x: 0, y: 0 });
        }}
      />

      <div
        class="relative shrink-0 shadow-2xl transition-all duration-150 m-auto"
        style={{
          width: `${props.design.canvas.width * effectiveScale()}px`,
          height: `${props.design.canvas.height * effectiveScale()}px`,
          transform: `translate(${pan().x}px, ${pan().y}px)`,
        }}
      >
        <div
          class="relative overflow-hidden"
          style={{
            width: `${props.design.canvas.width}px`,
            height: `${props.design.canvas.height}px`,
            transform: `scale(${effectiveScale()})`,
            "transform-origin": "top left",
            "background-color": props.design.backgroundColor,
            "background-image": props.design.background
              ? `url(${resolveStorageUrl(props.design.background)})`
              : undefined,
            "background-size": "cover",
            "background-position": "center",
          }}
        >
          <For each={elementIds()}>
            {(id) => (
              <CanvasElementItem
                id={id}
                design={() => props.design}
                scale={effectiveScale}
                selectedId={() => props.selectedId}
                editingId={editingId}
                onSelect={(selectedId) => {
                  props.onSelect(selectedId);
                  if (editingId() && editingId() !== selectedId) {
                    setEditingId(null);
                  }
                }}
                onStartEditing={(editId) => {
                  props.onSelect(editId);
                  setEditingId(editId);
                }}
                onStopEditing={() => setEditingId(null)}
                onChangeBounds={(bounds) => props.onChangeElement(id, bounds)}
                onDeleteElement={
                  props.onDeleteElement ? () => props.onDeleteElement?.(id) : undefined
                }
                previewValues={props.previewValues}
                textAdapter={props.textAdapter}
              />
            )}
          </For>
        </div>
      </div>
    </div>
  );
}

function CanvasElementItem(props: {
  id: string;
  design: () => BadgeDesign;
  scale: () => number;
  selectedId: () => string | null;
  editingId: () => string | null;
  onSelect: (id: string) => void;
  onStartEditing: (id: string) => void;
  onStopEditing: () => void;
  onChangeBounds: (bounds: ElementBounds) => void;
  onDeleteElement?: () => void;
  previewValues?: Record<string, string>;
  textAdapter?: TextElementAdapter;
}): JSX.Element {
  const element = createMemo(() =>
    props.design().elements.find((e) => e.id === props.id),
  );
  const index = createMemo(() =>
    props.design().elements.findIndex((e) => e.id === props.id),
  );
  const isSelected = createMemo(() => props.selectedId() === props.id);
  const isEditing = createMemo(() => props.editingId() === props.id);

  return (
    <Show when={element()}>
      {(el) => (
        <EditorElementFrame
          bounds={el()}
          scale={props.scale()}
          canvas={props.design().canvas}
          zIndex={index() + 1}
          selected={isSelected()}
          editing={isEditing()}
          onSelect={() => props.onSelect(props.id)}
          onDoubleClick={() => {
            if (el().type === "text") {
              props.onStartEditing(props.id);
            }
          }}
          onChangeBounds={props.onChangeBounds}
          onDelete={props.onDeleteElement}
        >
          <Show
            when={el().type === "text"}
            fallback={
              <Show
                when={el().type === "image"}
                fallback={
                  <Show when={el().type === "qr"}>
                    <div class="h-full w-full pointer-events-none">
                      <QrRenderer
                        value="https://univents.app/check-in"
                        foreground={
                          (el() as Extract<BadgeElement, { type: "qr" }>).foreground ?? "#000000"
                        }
                        background={
                          (el() as Extract<BadgeElement, { type: "qr" }>).background ?? "#ffffff"
                        }
                        style={
                          (el() as Extract<BadgeElement, { type: "qr" }>).style ?? "square"
                        }
                      />
                    </div>
                  </Show>
                }
              >
                {(() => {
                  const img = () =>
                    el() as Extract<BadgeElement, { type: "image" }>;
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
              const txt = () =>
                el() as Extract<BadgeElement, { type: "text" }>;
              return (
                <div
                  class="relative h-full w-full cursor-text"
                  onDblClick={(e) => {
                    e.stopPropagation();
                    props.onStartEditing(props.id);
                  }}
                >
                  <Show
                    when={isEditing()}
                    fallback={
                      <StaticText
                        paragraphs={txt().paragraphs}
                        values={props.previewValues}
                        showVariables={false}
                      />
                    }
                  >
                    <EditableText
                      element={txt()}
                      adapter={props.textAdapter}
                      onStopEditing={props.onStopEditing}
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
