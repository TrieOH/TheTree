import type { CertificateCanvasSize, CertificateElementType } from "./types";

export const DEFAULT_CERTIFICATE_CANVAS: CertificateCanvasSize = {
  width: 1123,
  height: 794,
};

export interface CertificateCanvasPreset {
  id: string;
  label: string;
  size: CertificateCanvasSize;
}

export const CERTIFICATE_CANVAS_PRESETS: CertificateCanvasPreset[] = [
  {
    id: "a4-landscape",
    label: "A4 Paisagem",
    size: { width: 1123, height: 794 },
  },
  {
    id: "a4-portrait",
    label: "A4 Retrato",
    size: { width: 794, height: 1123 },
  },
  {
    id: "letter-landscape",
    label: "Carta Paisagem",
    size: { width: 1056, height: 816 },
  },
  {
    id: "letter-portrait",
    label: "Carta Retrato",
    size: { width: 816, height: 1056 },
  },
  {
    id: "square",
    label: "Quadrado",
    size: { width: 1000, height: 1000 },
  },
];

export const CERTIFICATE_CANVAS_MAX_WIDTH = 980;
export const CERTIFICATE_CANVAS_MAX_HEIGHT = 640;

export const CERTIFICATE_ELEMENT_OVERFLOW: Record<
  CertificateElementType,
  number
> = {
  text: 0,
  hash: 0,
  image: 0.5,
  signature: 0.5,
};

export const CERTIFICATE_VARIABLES = [
  {
    key: "participant_name",
    token: "{{participant_name}}",
    label: "Nome completo do participante",
    description: "Nome legal informado no perfil do participante",
  },
  {
    key: "event_name",
    token: "{{event_name}}",
    label: "Nome do evento",
    description: "Nome do evento",
  },
  {
    key: "edition_name",
    token: "{{edition_name}}",
    label: "Nome da edição",
    description: "Nome da edição do evento",
  },
  {
    key: "activity_name",
    token: "{{activity_name}}",
    label: "Nome da atividade/edição",
    description: "Nome da atividade ou edição",
  },
  {
    key: "participation_type",
    token: "{{participation_type}}",
    label: "Tipo de participação",
    description: "Edição ou atividade",
  },
  {
    key: "location",
    token: "{{location}}",
    label: "Local",
    description: "Local informado na edição",
  },
  {
    key: "certified_at",
    token: "{{certified_at}}",
    label: "Data de certificação",
    description: "Data em que o certificado foi emitido",
  },
  {
    key: "cert_hash",
    token: "{{cert_hash}}",
    label: "Hash do certificado",
    description: "Identificador único para validação do certificado",
  },
  {
    key: "verify_url",
    token: "{{verify_url}}",
    label: "URL de verificação",
    description: "Endereço público para validar o certificado",
  },
] as const;

export type CertificateVariableKey =
  (typeof CERTIFICATE_VARIABLES)[number]["key"];

export const CERTIFICATE_FONT_FAMILIES = [
  { label: "Arial", value: "Arial, Helvetica, sans-serif" },
  { label: "Inter", value: "Inter, ui-sans-serif, system-ui, sans-serif" },
  { label: "Georgia", value: "Georgia, 'Times New Roman', serif" },
  { label: "Times New Roman", value: "'Times New Roman', Times, serif" },
  { label: "Playfair Display", value: "'Playfair Display', Georgia, serif" },
  { label: "Courier New", value: "'Courier New', Courier, monospace" },
] as const;

export const DEFAULT_CERTIFICATE_FONT = CERTIFICATE_FONT_FAMILIES[0].value;
export const DEFAULT_CERTIFICATE_TEXT_COLOR = "#000000";
export const DEFAULT_CERTIFICATE_SUBTLE_COLOR = "#4B5563";
export const MIN_CERTIFICATE_ELEMENT_SIZE = { width: 24, height: 24 };

export const CERTIFICATE_IMAGE_ACCEPT = "image/png,image/jpeg,image/webp";
export const CERTIFICATE_IMAGE_MIME_TYPES = new Set([
  "image/png",
  "image/jpeg",
  "image/webp",
]);
