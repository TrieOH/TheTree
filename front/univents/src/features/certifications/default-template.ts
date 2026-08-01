import type { CertificationTemplateI } from "./model";

/**
 * Default certificate used when an emission has no linked template.
 *
 * Replace only `data` (or this entire object) when the final certificate JSON
 * is provided. Keeping this shape stable lets the viewer and exporter work
 * with the default certificate exactly like a user-created template.
 */
export const DEFAULT_CERTIFICATION_TEMPLATE: CertificationTemplateI = {
  id: "default-certification-template",
  edition_id: "",
  kind: "edition_attendance",
  name: "Certificado",
  description: "Certificado padrão",
  created_at: "",
  design_data: {
    canvas: { width: 1123, height: 794 },
    background: null,
    elements: [],
  },
};

export function getCertificationTemplateOrDefault(
  template: CertificationTemplateI | null | undefined,
): CertificationTemplateI {
  return template ?? DEFAULT_CERTIFICATION_TEMPLATE;
}
