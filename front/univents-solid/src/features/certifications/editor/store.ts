import { createSignal } from "solid-js";
import type {
  CertificationTemplateCreateI,
  CertificationTemplateElement,
  CertificationTemplateI,
} from "../model";
import {
  DEFAULT_CERTIFICATE_CANVAS,
  DEFAULT_CERTIFICATE_DRAFT,
  MIN_CERTIFICATE_ELEMENT_SIZE,
} from "./constants";
import type { CanvasSize, CertificateSignatory } from "./types";
import type {
  RichParagraph,
  RichTextController,
  TextElementAdapter,
  TextSelectionStyles,
} from "@/features/editor/types";
import { scaleCertificateElements } from "./utils";

export type CertificateTextSelectionStyles = TextSelectionStyles;
export type CertificateRichTextController = RichTextController;

export interface CertificateEditorState {
  draft: CertificationTemplateCreateI;
  canvas: CanvasSize;
  availableSignatures: CertificateSignatory[];
  selectedElementId: string | null;
  editingElementId: string | null;
  richTextController: CertificateRichTextController | null;
  textSelectionStyles: CertificateTextSelectionStyles | null;
}

const [draft, setDraft] = createSignal<CertificationTemplateCreateI>(
  JSON.parse(JSON.stringify(DEFAULT_CERTIFICATE_DRAFT)),
);
const [canvas, setCanvas] = createSignal<CanvasSize>(
  { ...DEFAULT_CERTIFICATE_CANVAS },
);
const [availableSignatures, setAvailableSignatures] = createSignal<
  CertificateSignatory[]
>([]);
const [selectedElementId, setSelectedElementId] = createSignal<string | null>(
  null,
);
const [editingElementId, setEditingElementId] = createSignal<string | null>(
  null,
);
const [richTextController, setRichTextController] =
  createSignal<CertificateRichTextController | null>(null);
const [textSelectionStyles, setTextSelectionStyles] =
  createSignal<CertificateTextSelectionStyles | null>(null);

export const certificateEditorStore = {
  get state() {
    return {
      draft: draft(),
      canvas: canvas(),
      availableSignatures: availableSignatures(),
      selectedElementId: selectedElementId(),
      editingElementId: editingElementId(),
      richTextController: richTextController(),
      textSelectionStyles: textSelectionStyles(),
    };
  },
  draft,
  canvas,
  availableSignatures,
  selectedElementId,
  editingElementId,
  richTextController,
  textSelectionStyles,
};

export const certificateEditorActions = {
  reset() {
    setDraft(JSON.parse(JSON.stringify(DEFAULT_CERTIFICATE_DRAFT)));
    setCanvas({ ...DEFAULT_CERTIFICATE_CANVAS });
    setSelectedElementId(null);
    setEditingElementId(null);
    setRichTextController(null);
    setTextSelectionStyles(null);
    setAvailableSignatures([]);
  },

  loadDraft(template: CertificationTemplateI) {
    const nextCanvas =
      template.design_data.canvas ?? DEFAULT_CERTIFICATE_CANVAS;
    setDraft(JSON.parse(JSON.stringify(template)));
    setCanvas({ ...nextCanvas });
    setSelectedElementId(null);
    setEditingElementId(null);
    setRichTextController(null);
    setTextSelectionStyles(null);
  },

  getDraft(): CertificationTemplateCreateI {
    const currentDraft = draft();
    const currentCanvas = canvas();
    return {
      ...currentDraft,
      design_data: {
        ...currentDraft.design_data,
        canvas: { ...currentCanvas },
      },
    };
  },

  setName(name: string) {
    setDraft((prev) => ({ ...prev, name }));
  },

  setKind(kind: "edition_attendance" | "program_attendance") {
    setDraft((prev) => ({ ...prev, kind }));
  },

  setDescription(description: string) {
    setDraft((prev) => ({ ...prev, description }));
  },

  setBackgroundUrl(background: string | null) {
    setDraft((prev) => ({
      ...prev,
      design_data: {
        ...prev.design_data,
        background,
      },
    }));
  },

  setCanvasSize(nextSize: CanvasSize) {
    const current = canvas();
    if (current.width === nextSize.width && current.height === nextSize.height) {
      return;
    }

    const scaleX = nextSize.width / current.width;
    const scaleY = nextSize.height / current.height;
    const scale = Math.min(scaleX, scaleY);

    setCanvas({ ...nextSize });
    setDraft((prev) => ({
      ...prev,
      design_data: {
        ...prev.design_data,
        canvas: { ...nextSize },
        elements: scaleCertificateElements(
          prev.design_data.elements,
          scaleX,
          scaleY,
          scale,
        ),
      },
    }));
  },

  addElement(element: CertificationTemplateElement) {
    setDraft((prev) => ({
      ...prev,
      design_data: {
        ...prev.design_data,
        elements: [...prev.design_data.elements, element],
      },
    }));
    setSelectedElementId(element.id);
  },

  removeElement(elementId: string) {
    setDraft((prev) => ({
      ...prev,
      design_data: {
        ...prev.design_data,
        elements: prev.design_data.elements.filter(
          (item: CertificationTemplateElement) => item.id !== elementId,
        ),
      },
    }));
    if (selectedElementId() === elementId) {
      setSelectedElementId(null);
    }
    if (editingElementId() === elementId) {
      setEditingElementId(null);
      setRichTextController(null);
      setTextSelectionStyles(null);
    }
  },

  updateElement(
    elementId: string,
    updater: (
      element: CertificationTemplateElement,
    ) => CertificationTemplateElement,
  ) {
    setDraft((prev) => ({
      ...prev,
      design_data: {
        ...prev.design_data,
        elements: prev.design_data.elements.map(
          (el: CertificationTemplateElement) =>
            el.id === elementId ? updater(el) : el,
        ),
      },
    }));
  },

  updateElementBounds(
    elementId: string,
    bounds: Partial<{ x: number; y: number; width: number; height: number }>,
  ) {
    this.updateElement(elementId, (element: CertificationTemplateElement) => {
      const nextWidth = Math.max(
        MIN_CERTIFICATE_ELEMENT_SIZE.width,
        bounds.width ?? element.width,
      );
      const nextHeight = Math.max(
        MIN_CERTIFICATE_ELEMENT_SIZE.height,
        bounds.height ?? element.height,
      );
      return {
        ...element,
        x: bounds.x ?? element.x,
        y: bounds.y ?? element.y,
        width: nextWidth,
        height: nextHeight,
      };
    });
  },

  selectElement(elementId: string | null) {
    setSelectedElementId(elementId);
  },

  startEditing(elementId: string) {
    setSelectedElementId(elementId);
    setEditingElementId(elementId);
  },

  stopEditing() {
    setEditingElementId(null);
    setRichTextController(null);
    setTextSelectionStyles(null);
  },

  setAvailableSignatures(signatures: CertificateSignatory[]) {
    setAvailableSignatures(signatures);
  },

  bringForward(elementId: string) {
    const list = [...draft().design_data.elements];
    const index = list.findIndex(
      (element: CertificationTemplateElement) => element.id === elementId,
    );
    if (index === -1 || index >= list.length - 1) return;
    const [item] = list.splice(index, 1);
    list.splice(index + 1, 0, item);
    setDraft((prev) => ({
      ...prev,
      design_data: { ...prev.design_data, elements: list },
    }));
  },

  sendBackward(elementId: string) {
    const list = [...draft().design_data.elements];
    const index = list.findIndex(
      (element: CertificationTemplateElement) => element.id === elementId,
    );
    if (index <= 0) return;
    const [item] = list.splice(index, 1);
    list.splice(index - 1, 0, item);
    setDraft((prev) => ({
      ...prev,
      design_data: { ...prev.design_data, elements: list },
    }));
  },

  setRichTextController(controller: CertificateRichTextController | null) {
    setRichTextController(controller);
  },

  setTextSelectionStyles(styles: CertificateTextSelectionStyles | null) {
    setTextSelectionStyles(styles);
  },

  textAdapter: {
    updateParagraphs(elementId: string, paragraphs: RichParagraph[]) {
      certificateEditorActions.updateElement(elementId, (element) =>
        element.type === "text"
          ? ({ ...element, paragraphs } as CertificationTemplateElement)
          : element,
      );
    },
    setController(controller: RichTextController | null) {
      setRichTextController(controller);
    },
    setSelectionStyles(styles: TextSelectionStyles | null) {
      setTextSelectionStyles(styles);
    },
    stopEditing() {
      certificateEditorActions.stopEditing();
    },
  } as TextElementAdapter,
};
