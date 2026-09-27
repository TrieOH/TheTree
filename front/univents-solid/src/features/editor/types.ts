export const DEFAULT_EDITOR_FONT = "Inter, sans-serif";
export const DEFAULT_EDITOR_LINE_HEIGHT = 1.25;
export const DEFAULT_EDITOR_TEXT_COLOR = "#000000";

export const FONT_FAMILIES = [
  { label: "Arial", value: "Arial, Helvetica, sans-serif" },
  { label: "Inter", value: "Inter, ui-sans-serif, system-ui, sans-serif" },
  { label: "Georgia", value: "Georgia, 'Times New Roman', serif" },
  { label: "Times New Roman", value: "'Times New Roman', Times, serif" },
  { label: "Playfair Display", value: "'Playfair Display', Georgia, serif" },
  { label: "Courier New", value: "'Courier New', Courier, monospace" },
] as const;

export const LINE_HEIGHT_OPTIONS = [
  { value: "1", label: "Simples" },
  { value: "1.15", label: "1,15" },
  { value: "1.25", label: "1,25" },
  { value: "1.5", label: "1,5" },
  { value: "2", label: "Duplo" },
] as const;

export interface ElementBounds {
  x: number;
  y: number;
  width: number;
  height: number;
}

export type ResizeHandle = "nw" | "ne" | "sw" | "se";

export interface RichRun {
  text: string;
  bold: boolean;
  italic: boolean;
  underline: boolean;
  color: string;
  fontSize: number;
  fontFamily: string;
}

export interface RichParagraph {
  align: "left" | "center" | "right" | "justify";
  lineHeight: number;
  runs: RichRun[];
}

export interface TextSelectionStyles {
  bold: boolean | null;
  italic: boolean | null;
  underline: boolean | null;
  align: "left" | "center" | "right" | "justify" | null;
  lineHeight: number | null;
  color: string | null;
  fontSize: number | null;
  fontFamily: string | null;
}

export interface RichTextController {
  elementId: string;
  commit: () => void;
  toggleBold: () => void;
  toggleItalic: () => void;
  toggleUnderline: () => void;
  setAlign: (align: "left" | "center" | "right" | "justify") => void;
  setLineHeight: (lineHeight: number) => void;
  setColor: (color: string) => void;
  setFontSize: (fontSize: number) => void;
  setFontFamily: (fontFamily: string) => void;
  insertText: (text: string) => void;
}

export interface TextElementAdapter {
  updateParagraphs: (elementId: string, paragraphs: RichParagraph[]) => void;
  setController: (controller: RichTextController | null) => void;
  setSelectionStyles: (styles: TextSelectionStyles | null) => void;
  stopEditing: () => void;
}

export interface EditorVariable {
  value: string;
  label: string;
  description?: string;
}

export interface CanvasDimensions {
  width: number;
  height: number;
}
