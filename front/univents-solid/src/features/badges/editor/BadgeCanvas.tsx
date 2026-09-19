import type { JSX } from "@solidjs/web";
import {
  For,
  Show,
  createEffect,
  createMemo,
  createSignal,
  onCleanup,
  untrack,
} from "solid-js";
import {
  EditorElementFrame,
  QrRenderer,
  StaticText,
  domToParagraphs,
  paragraphsToHtml,
  type ElementBounds,
  type TextElementAdapter,
  type TextSelectionStyles,
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
  const [scale, setScale] = createSignal(1, { ownedWrite: true } as any);
  const [editingId, setEditingId] = createSignal<string | null>(null, {
    ownedWrite: true,
  } as any);

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

function EditableText(props: {
  element: Extract<BadgeElement, { type: "text" }>;
  adapter?: TextElementAdapter;
  onStopEditing: () => void;
}): JSX.Element {
  let editorRef: HTMLDivElement | undefined;
  let isMounted = true;
  let savedRange: Range | null = null;

  const saveSelection = () => {
    const sel = window.getSelection();
    if (sel && sel.rangeCount > 0 && editorRef && rangeWithin(sel, editorRef)) {
      savedRange = sel.getRangeAt(0).cloneRange();
    }
  };

  const restoreSelection = () => {
    if (!editorRef) return;
    editorRef.focus();
    if (savedRange) {
      const sel = window.getSelection();
      if (sel) {
        sel.removeAllRanges();
        sel.addRange(savedRange);
      }
    }
  };

  const readStyles = (): TextSelectionStyles => {
    const sel = window.getSelection();
    let targetNode: Node | null = sel?.anchorNode ?? null;
    if (targetNode?.nodeType === Node.TEXT_NODE) {
      targetNode = targetNode.parentElement;
    }
    const computed = targetNode
      ? window.getComputedStyle(targetNode as Element)
      : null;

    return {
      bold:
        computed?.fontWeight === "bold" ||
        Number(computed?.fontWeight ?? 400) >= 700,
      italic: computed?.fontStyle === "italic",
      underline: computed?.textDecorationLine?.includes("underline") ?? false,
      align:
        (computed?.textAlign as
          | "left"
          | "center"
          | "right"
          | "justify"
          | undefined) ?? "left",
      lineHeight: computed?.lineHeight
        ? Number.parseFloat(computed.lineHeight) /
          Number.parseFloat(computed.fontSize || "16")
        : 1.25,
      color: computed?.color ?? "#0f172a",
      fontSize: computed?.fontSize
        ? Number.parseFloat(computed.fontSize)
        : 18,
      fontFamily: computed?.fontFamily ?? "Inter, sans-serif",
    };
  };

  const sync = () => {
    if (!editorRef) return;
    const paragraphs = domToParagraphs(editorRef);
    const adapter = untrack(() => props.adapter);
    const elementId = untrack(() => props.element.id);
    adapter?.updateParagraphs(elementId, paragraphs);
    adapter?.setSelectionStyles(readStyles());
  };

  const initEditor = (el: HTMLDivElement) => {
    editorRef = el;
    const initialParagraphs = untrack(() => props.element.paragraphs);
    el.innerHTML = paragraphsToHtml(initialParagraphs);
    el.focus();
    try {
      const range = document.createRange();
      range.selectNodeContents(el);
      range.collapse(false);
      const sel = window.getSelection();
      sel?.removeAllRanges();
      sel?.addRange(range);
      savedRange = range.cloneRange();
    } catch {
      // ignore
    }

    const adapter = untrack(() => props.adapter);
    const elementId = untrack(() => props.element.id);

    queueMicrotask(() => {
      if (!isMounted) return;
      adapter?.setController({
        elementId,
        commit: sync,
        toggleBold: () => {
          restoreSelection();
          document.execCommand("bold");
          sync();
          saveSelection();
        },
        toggleItalic: () => {
          restoreSelection();
          document.execCommand("italic");
          sync();
          saveSelection();
        },
        toggleUnderline: () => {
          restoreSelection();
          document.execCommand("underline");
          sync();
          saveSelection();
        },
        setAlign: (align) => {
          restoreSelection();
          if (align === "left") document.execCommand("justifyLeft");
          else if (align === "center") document.execCommand("justifyCenter");
          else if (align === "right") document.execCommand("justifyRight");
          else if (align === "justify") document.execCommand("justifyFull");
          sync();
          saveSelection();
        },
        setLineHeight: (lineHeight) => {
          restoreSelection();
          const sel = window.getSelection();
          if (sel && sel.anchorNode && editorRef) {
            let node: Node | null = sel.anchorNode;
            if (node.nodeType === Node.TEXT_NODE) node = node.parentElement;
            const block = (node as HTMLElement)?.closest("p, div");
            if (block && editorRef.contains(block)) {
              (block as HTMLElement).style.lineHeight = `${lineHeight}`;
            }
          }
          sync();
          saveSelection();
        },
        setColor: (color) => {
          restoreSelection();
          document.execCommand("foreColor", false, color);
          sync();
          saveSelection();
        },
        setFontSize: (fontSize) => {
          restoreSelection();
          const sel = window.getSelection();
          if (sel && !sel.isCollapsed && rangeWithin(sel, editorRef!)) {
            applyInlineStyleToSelection({ fontSize: `${fontSize}px` }, editorRef!);
          } else if (editorRef) {
            editorRef.style.fontSize = `${fontSize}px`;
          }
          sync();
          saveSelection();
        },
        setFontFamily: (fontFamily) => {
          restoreSelection();
          const sel = window.getSelection();
          if (sel && !sel.isCollapsed && rangeWithin(sel, editorRef!)) {
            applyInlineStyleToSelection({ fontFamily }, editorRef!);
          } else if (editorRef) {
            editorRef.style.fontFamily = fontFamily;
          }
          sync();
          saveSelection();
        },
        insertText: (text) => {
          restoreSelection();
          document.execCommand("insertText", false, text);
          sync();
          saveSelection();
        },
      });

      adapter?.setSelectionStyles(readStyles());
    });
  };

  onCleanup(() => {
    isMounted = false;
    const adapter = untrack(() => props.adapter);
    queueMicrotask(() => {
      adapter?.setController(null);
      adapter?.setSelectionStyles(null);
    });
  });

  return (
    <div
      ref={initEditor}
      contenteditable
      class="h-full w-full outline-none focus:outline-none"
      style={{ "overflow-wrap": "anywhere", "white-space": "pre-wrap" }}
      onInput={() => {
        saveSelection();
        sync();
      }}
      onBlur={() => {
        saveSelection();
        sync();
      }}
      onKeyDown={(event) => {
        if (event.key === "Escape") {
          event.preventDefault();
          sync();
          props.onStopEditing();
        }
      }}
    />
  );
}

function rangeWithin(selection: Selection, container: HTMLElement): boolean {
  if (selection.rangeCount === 0) return false;
  const range = selection.getRangeAt(0);
  return (
    container.contains(range.commonAncestorContainer) ||
    container === range.commonAncestorContainer
  );
}

function applyInlineStyleToSelection(
  styles: Record<string, string>,
  container: HTMLElement,
): void {
  const sel = window.getSelection();
  if (!sel || sel.rangeCount === 0) return;
  const range = sel.getRangeAt(0);
  if (range.collapsed || !rangeWithin(sel, container)) return;

  const span = document.createElement("span");
  for (const [key, value] of Object.entries(styles)) {
    (span.style as any)[key] = value;
  }
  try {
    range.surroundContents(span);
  } catch {
    const contents = range.extractContents();
    span.appendChild(contents);
    range.insertNode(span);
  }
}
