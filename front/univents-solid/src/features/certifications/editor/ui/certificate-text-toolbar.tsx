import type { JSX } from "@solidjs/web";
import { createMemo } from "solid-js";
import {
  DEFAULT_EDITOR_FONT,
  DEFAULT_EDITOR_TEXT_COLOR,
  RichTextToolbar,
} from "@/features/editor";
import type { RichTextController, TextSelectionStyles } from "@/features/editor/types";
import { certificateEditorActions, certificateEditorStore } from "../store";
import type { TextCertificateElement } from "../types";

export function CertificateTextToolbar(): JSX.Element {
  const selectedTextElement = createMemo<TextCertificateElement | null>(() => {
    const selectedId = certificateEditorStore.selectedElementId();
    if (!selectedId) return null;
    const el = certificateEditorStore
      .draft()
      .design_data.elements.find((item) => item.id === selectedId);
    return el && el.type === "text" ? (el as TextCertificateElement) : null;
  });

  const fallbackController = createMemo<RichTextController | null>(() => {
    const item = selectedTextElement();
    if (!item) return null;
    const elementId = item.id;

    return {
      elementId,
      commit: () => {},
      toggleBold: () => {
        const current = selectedTextElement();
        if (!current) return;
        const anyBold = current.paragraphs.some((p) =>
          p.runs.some((r) => r.bold),
        );
        certificateEditorActions.updateElement(current.id, (el) => {
          if (el.type !== "text") return el;
          return {
            ...el,
            paragraphs: el.paragraphs.map((p) => ({
              ...p,
              runs: p.runs.map((r) => ({ ...r, bold: !anyBold })),
            })),
          };
        });
      },
      toggleItalic: () => {
        const current = selectedTextElement();
        if (!current) return;
        const anyItalic = current.paragraphs.some((p) =>
          p.runs.some((r) => r.italic),
        );
        certificateEditorActions.updateElement(current.id, (el) => {
          if (el.type !== "text") return el;
          return {
            ...el,
            paragraphs: el.paragraphs.map((p) => ({
              ...p,
              runs: p.runs.map((r) => ({ ...r, italic: !anyItalic })),
            })),
          };
        });
      },
      toggleUnderline: () => {
        const current = selectedTextElement();
        if (!current) return;
        const anyUnderline = current.paragraphs.some((p) =>
          p.runs.some((r) => r.underline),
        );
        certificateEditorActions.updateElement(current.id, (el) => {
          if (el.type !== "text") return el;
          return {
            ...el,
            paragraphs: el.paragraphs.map((p) => ({
              ...p,
              runs: p.runs.map((r) => ({ ...r, underline: !anyUnderline })),
            })),
          };
        });
      },
      setAlign: (align) => {
        const current = selectedTextElement();
        if (!current) return;
        certificateEditorActions.updateElement(current.id, (el) => {
          if (el.type !== "text") return el;
          return {
            ...el,
            paragraphs: el.paragraphs.map((p) => ({ ...p, align })),
          };
        });
      },
      setLineHeight: (lineHeight) => {
        const current = selectedTextElement();
        if (!current) return;
        certificateEditorActions.updateElement(current.id, (el) => {
          if (el.type !== "text") return el;
          return {
            ...el,
            paragraphs: el.paragraphs.map((p) => ({ ...p, lineHeight })),
          };
        });
      },
      setColor: (color) => {
        const current = selectedTextElement();
        if (!current) return;
        certificateEditorActions.updateElement(current.id, (el) => {
          if (el.type !== "text") return el;
          return {
            ...el,
            paragraphs: el.paragraphs.map((p) => ({
              ...p,
              runs: p.runs.map((r) => ({ ...r, color })),
            })),
          };
        });
      },
      setFontSize: (fontSize) => {
        const current = selectedTextElement();
        if (!current) return;
        certificateEditorActions.updateElement(current.id, (el) => {
          if (el.type !== "text") return el;
          return {
            ...el,
            paragraphs: el.paragraphs.map((p) => ({
              ...p,
              runs: p.runs.map((r) => ({ ...r, fontSize })),
            })),
          };
        });
      },
      setFontFamily: (fontFamily) => {
        const current = selectedTextElement();
        if (!current) return;
        certificateEditorActions.updateElement(current.id, (el) => {
          if (el.type !== "text") return el;
          return {
            ...el,
            paragraphs: el.paragraphs.map((p) => ({
              ...p,
              runs: p.runs.map((r) => ({ ...r, fontFamily })),
            })),
          };
        });
      },
      insertText: (text) => {
        const current = selectedTextElement();
        if (!current) return;
        const lastIdx = current.paragraphs.length - 1;
        certificateEditorActions.updateElement(current.id, (el) => {
          if (el.type !== "text") return el;
          return {
            ...el,
            paragraphs: el.paragraphs.map((p, idx) =>
              idx === lastIdx
                ? {
                    ...p,
                    runs: [
                      ...p.runs,
                      {
                        text: ` ${text} `,
                        fontSize: p.runs[0]?.fontSize ?? 18,
                        fontFamily: p.runs[0]?.fontFamily ?? DEFAULT_EDITOR_FONT,
                        color: p.runs[0]?.color ?? DEFAULT_EDITOR_TEXT_COLOR,
                        bold: p.runs[0]?.bold ?? false,
                        italic: p.runs[0]?.italic ?? false,
                        underline: p.runs[0]?.underline ?? false,
                      },
                    ],
                  }
                : p,
            ),
          };
        });
      },
    };
  });

  const fallbackStyles = createMemo<TextSelectionStyles | null>(() => {
    const item = selectedTextElement();
    if (!item) return null;
    const firstRun = item.paragraphs[0]?.runs[0];
    return {
      bold: Boolean(firstRun?.bold),
      italic: Boolean(firstRun?.italic),
      underline: Boolean(firstRun?.underline),
      align: item.paragraphs[0]?.align ?? "left",
      lineHeight: item.paragraphs[0]?.lineHeight ?? 1.25,
      color: firstRun?.color ?? DEFAULT_EDITOR_TEXT_COLOR,
      fontSize: firstRun?.fontSize ?? 18,
      fontFamily: firstRun?.fontFamily ?? DEFAULT_EDITOR_FONT,
    };
  });

  const effectiveController = () =>
    certificateEditorStore.richTextController() ?? fallbackController();
  const effectiveStyles = () =>
    certificateEditorStore.textSelectionStyles() ?? fallbackStyles();

  return (
    <RichTextToolbar
      controller={effectiveController()}
      selectionStyles={effectiveStyles()}
    />
  );
}
