import type { JSX } from "@solidjs/web";
import { onCleanup, untrack } from "solid-js";
import { domToParagraphs, paragraphsToHtml } from "./dom-serializer";
import { normalizeHexColor } from "./rich-text-toolbar";
import { DEFAULT_EDITOR_FONT } from "./types";
import type {
  RichParagraph,
  TextElementAdapter,
  TextSelectionStyles,
} from "./types";

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
    span.style.setProperty(
      key.replace(/([A-Z])/g, "-$1").toLowerCase(),
      value,
    );
  }
  try {
    range.surroundContents(span);
  } catch {
    const contents = range.extractContents();
    span.appendChild(contents);
    range.insertNode(span);
  }
}

export interface EditableTextProps {
  element: { id: string; paragraphs: readonly RichParagraph[] | RichParagraph[] };
  adapter?: TextElementAdapter;
  onStopEditing?: () => void;
}

export function EditableText(props: EditableTextProps): JSX.Element {
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
      color: normalizeHexColor(computed?.color),
      fontSize: computed?.fontSize
        ? Number.parseFloat(computed.fontSize)
        : 18,
      fontFamily:
        computed?.fontFamily?.split(",")[0]?.trim().replace(/^['"]|['"]$/g, "") ||
        DEFAULT_EDITOR_FONT,
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
      onPointerDown={(event) => event.stopPropagation()}
      onInput={() => {
        saveSelection();
        sync();
      }}
      onKeyUp={() => {
        saveSelection();
        const adapter = untrack(() => props.adapter);
        adapter?.setSelectionStyles(readStyles());
      }}
      onMouseUp={() => {
        saveSelection();
        const adapter = untrack(() => props.adapter);
        adapter?.setSelectionStyles(readStyles());
      }}
      onBlur={() => {
        saveSelection();
        sync();
      }}
      onKeyDown={(event) => {
        if (event.key === "Escape") {
          event.preventDefault();
          sync();
          props.onStopEditing?.();
        }
      }}
    />
  );
}
