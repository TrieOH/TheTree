import type {
  CertificationTemplateCreateI,
  CertificationTemplateElement,
} from "../model";

export interface CertificateCanvasSize {
  width: number;
  height: number;
}

export type CertificateElementType = CertificationTemplateElement["type"];
export type HashCertificateElement = Extract<
  CertificationTemplateElement,
  { type: "hash" }
>;
export type TextCertificateElement = Extract<
  CertificationTemplateElement,
  { type: "text" }
>;
export type ImageCertificateElement = Extract<
  CertificationTemplateElement,
  { type: "image" }
>;
export type SignatureCertificateElement = Extract<
  CertificationTemplateElement,
  { type: "signature" }
>;

export type CertificationTemplateDraft = CertificationTemplateCreateI;
