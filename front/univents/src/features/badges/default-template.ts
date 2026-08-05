import type { BadgeTemplateCreate } from "./model";

export const DEFAULT_BADGE_TEMPLATE: BadgeTemplateCreate = {
  name: "Crachá padrão",
  ticket_type_id: null,
  design_data: {
    canvas: { width: 638, height: 1011 },
    backgroundColor: "#ffffff",
    background: null,
    elements: [
      {
        id: "default-event",
        type: "text",
        x: 64,
        y: 90,
        width: 510,
        height: 80,
        paragraphs: [
          {
            align: "center",
            lineHeight: 1.25,
            runs: [
              {
                text: "{{event_name}}",
                fontSize: 28,
                fontFamily: "Inter, sans-serif",
                color: "#64748b",
                bold: true,
                italic: false,
                underline: false,
              },
            ],
          },
        ],
      },
      {
        id: "default-name",
        type: "text",
        x: 54,
        y: 390,
        width: 530,
        height: 130,
        paragraphs: [
          {
            align: "center",
            lineHeight: 1.25,
            runs: [
              {
                text: "{{participant_name}}",
                fontSize: 48,
                fontFamily: "Inter, sans-serif",
                color: "#0f172a",
                bold: true,
                italic: false,
                underline: false,
              },
            ],
          },
        ],
      },
      {
        id: "default-ticket",
        type: "text",
        x: 84,
        y: 535,
        width: 470,
        height: 60,
        paragraphs: [
          {
            align: "center",
            lineHeight: 1.25,
            runs: [
              {
                text: "{{ticket_name}}",
                fontSize: 24,
                fontFamily: "Inter, sans-serif",
                color: "#475569",
                bold: false,
                italic: false,
                underline: false,
              },
            ],
          },
        ],
      },
      {
        id: "default-qr",
        type: "qr",
        x: 219,
        y: 720,
        width: 200,
        height: 200,
        value: "{{checkin_url}}",
        foreground: "#0f172a",
        background: "#ffffff",
        style: "square",
      },
    ],
  },
};
