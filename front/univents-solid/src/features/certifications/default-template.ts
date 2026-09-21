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
    canvas: { width: 1484, height: 1060 },
    background: "/certificate-template.webp",
    elements: [
      {
        id: "hash_0655c435-ae03-4b92-b46f-edd6f96890a8",
        x: 77.30357356680216,
        y: 806.8932934927003,
        width: 848.0607536197364,
        height: 129.14380903918297,
        type: "hash",
        hashLabel: "Código de verificação",
        hash: "{{cert_hash}}",
        linkLabel: "Link de Verificação",
        url: "{{verify_url}}",
        fontSize: 35,
        color: "#4B5563",
        align: "left",
      },
      {
        id: "text_5a9d948a-6a88-4518-afe9-a5e5bdb21e24",
        x: 553.3272519954389,
        y: 65.15291459771808,
        width: 408.1707828611893,
        height: 54.484472562035315,
        type: "text",
        paragraphs: [
          {
            align: "center",
            lineHeight: 1.25,
            runs: [
              {
                text: "Univents",
                bold: false,
                italic: false,
                underline: false,
                color: "#0037b0",
                fontSize: 45,
                fontFamily: "Inter, ui-sans-serif, system-ui, sans-serif",
              },
            ],
          },
        ],
      },
      {
        id: "text_e383b3f8-f866-4890-a031-e11e73660732",
        x: 131.72878917537273,
        y: 152.93733517056515,
        width: 1241.8137415637495,
        height: 97.31634419366043,
        type: "text",
        paragraphs: [
          {
            align: "center",
            lineHeight: 1.25,
            runs: [
              {
                text: "Certificado de Participação",
                bold: true,
                italic: false,
                underline: false,
                color: "#0037b0",
                fontSize: 80,
                fontFamily: "Arial, Helvetica, sans-serif",
              },
            ],
          },
        ],
      },
      {
        id: "text_a20880ed-e786-4fb2-b4d1-c9b4e29cb951",
        x: 0,
        y: 313.6307800507964,
        width: 1484,
        height: 44.44711877707854,
        type: "text",
        paragraphs: [
          {
            align: "center",
            lineHeight: 1.25,
            runs: [
              {
                text: "Certificamos, para os devidos fins, que",
                bold: false,
                italic: false,
                underline: false,
                color: "#4b5563",
                fontSize: 34,
                fontFamily: "Arial, Helvetica, sans-serif",
              },
            ],
          },
        ],
      },
      {
        id: "text_486afa2a-25e3-4a7a-a9e7-46a40ebfe535",
        x: 87.80213088716067,
        y: 730.2095262934866,
        width: 831.5993473931637,
        height: 64.45155355387776,
        type: "text",
        paragraphs: [
          {
            align: "left",
            lineHeight: 1.25,
            runs: [
              {
                text: "Data de Emissão: {{certified_at}}",
                bold: false,
                italic: false,
                underline: false,
                color: "#4b5563",
                fontSize: 35,
                fontFamily: "Arial, Helvetica, sans-serif",
              },
            ],
          },
        ],
      },
      {
        id: "text_a642ce8f-9a13-4a9c-b54e-cd1f32080916",
        x: 0,
        y: 412.4984066111217,
        width: 1484,
        height: 58.41100517297852,
        type: "text",
        paragraphs: [
          {
            align: "center",
            lineHeight: 1.25,
            runs: [
              {
                text: "{{participant_name}}",
                bold: false,
                italic: false,
                underline: false,
                color: "#000000",
                fontSize: 46,
                fontFamily: "Arial, Helvetica, sans-serif",
              },
            ],
          },
        ],
      },
      {
        id: "text_e840b4b3-2724-403e-a22d-cd845f2da0d0",
        x: 232.0979081044035,
        y: 503.03829658962906,
        width: 951.9448654280517,
        height: 72.72305421442013,
        type: "text",
        paragraphs: [
          {
            align: "center",
            lineHeight: 1.25,
            runs: [
              {
                text: "concluiu com êxito {{participation_type}} ",
                bold: false,
                italic: false,
                underline: false,
                color: "#4b5563",
                fontSize: 24,
                fontFamily: "Arial, Helvetica, sans-serif",
              },
              {
                text: "{{activity_name}}",
                bold: true,
                italic: false,
                underline: false,
                color: "#000000",
                fontSize: 24,
                fontFamily: "Arial, Helvetica, sans-serif",
              },
              {
                text: " durante o evento ",
                bold: false,
                italic: false,
                underline: false,
                color: "#4b5563",
                fontSize: 24,
                fontFamily: "Arial, Helvetica, sans-serif",
              },
              {
                text: "{{event_name}}",
                bold: true,
                italic: false,
                underline: false,
                color: "#000000",
                fontSize: 24,
                fontFamily: "Arial, Helvetica, sans-serif",
              },
              {
                text: ", realizada em {{location}}.",
                bold: false,
                italic: false,
                underline: false,
                color: "#4b5563",
                fontSize: 24,
                fontFamily: "Arial, Helvetica, sans-serif",
              },
            ],
          },
        ],
      },
    ],
  },
};

export function getCertificationTemplateOrDefault(
  template: CertificationTemplateI | null | undefined,
): CertificationTemplateI {
  return template ?? DEFAULT_CERTIFICATION_TEMPLATE;
}
