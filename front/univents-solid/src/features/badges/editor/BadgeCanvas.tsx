import type { JSX } from "@solidjs/web";
import {
  For,
  Show,
  createEffect,
  createMemo,
  createSignal,
} from "solid-js";
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

export function BadgeCanvas(props: BadgeCanvasProps): JSX.Element {
  let containerRef: HTMLDivElement | undefined;
  const [scale, setScale] = createSignal(1, { ownedWrite: true });
  const [editingId, setEditingId] = createSignal<string | null>(null, {
    ownedWrite: true,
  });

  const elementIds = createMemo(() => props.design.elements.map((e) => e.id));

  createEffect(
    () => ({
      width: props.design.canvas.width,
      height: props.design.canvas.height,
    }),
    ({ width, height }) => {
      if (!containerRef) return;
      const { clientWidth, clientHeight } = containerRef;
      const padding = 64;
      const availW = Math.max(100, clientWidth - padding);
      const availH = Math.max(100, clientHeight - padding);
      const scaleW = availW / width;
      const scaleH = availH / height;
      setScale(Math.min(1.5, Math.max(0.2, Math.min(scaleW, scaleH))));
    },
  );

  return (
    <div
      ref={(el) => (containerRef = el)}
      class="relative flex flex-1 items-center justify-center overflow-auto p-8 select-none"
      style={{
        background:
          "radial-gradient(circle, var(--color-muted) 1px, transparent 1px)",
        "background-size": "16px 16px",
      }}
      onPointerDown={(event) => {
        if (event.target === event.currentTarget) {
          props.onSelect(null);
          setEditingId(null);
        }
      }}
    >
      <div
        class="relative shadow-2xl transition-all duration-150"
        style={{
          width: `${props.design.canvas.width}px`,
          height: `${props.design.canvas.height}px`,
          transform: `scale(${scale()})`,
          "transform-origin": "center center",
          "background-color": props.design.backgroundColor,
          "background-image": props.design.background
            ? `url(${props.design.background})`
            : undefined,
          "background-size": "cover",
          "background-position": "center",
        }}
        onPointerDown={(event) => {
          if (event.target === event.currentTarget) {
            props.onSelect(null);
            setEditingId(null);
          }
        }}
      >
        <For each={elementIds()}>
          {(id) => (
            <CanvasElementItem
              id={id}
              design={() => props.design}
              scale={scale}
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
          onDelete={el().type !== "qr" ? props.onDeleteElement : undefined}
        >
          <Show
            when={el().type === "text"}
            fallback={
              <Show
                when={el().type === "image"}
                fallback={
                  <div class="h-full w-full pointer-events-none">
                    <QrRenderer
                      value="https://univents.app/check-in"
                      foreground={
                        el().type === "qr"
                          ? (el() as Extract<BadgeElement, { type: "qr" }>).foreground
                          : "#000000"
                      }
                      background={
                        el().type === "qr"
                          ? (el() as Extract<BadgeElement, { type: "qr" }>).background
                          : "#ffffff"
                      }
                      style={
                        el().type === "qr"
                          ? (el() as Extract<BadgeElement, { type: "qr" }>).style
                          : "square"
                      }
                    />
                  </div>
                }
              >
                {(() => {
                  const img = () =>
                    el() as Extract<BadgeElement, { type: "image" }>;
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
