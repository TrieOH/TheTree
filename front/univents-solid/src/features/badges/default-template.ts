import type { BadgeTemplateCreate } from "./model";

export const DEFAULT_BADGE_TEMPLATE: BadgeTemplateCreate = {
  name: "Crachá padrão",
  ticket_type_id: null,
  origin: null,
  design_data: {
    canvas: { width: 321, height: 204 },
    backgroundColor: "#ffffff",
    background: null,
    elements: [
      {
        id: "default-event",
        type: "text",
        x: 32,
        y: 18,
        width: 257,
        height: 20,
        paragraphs: [
          {
            align: "center",
            lineHeight: 1.25,
            runs: [
              {
                text: "{{event_name}}",
                fontSize: 9,
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
        x: 27,
        y: 78,
        width: 267,
        height: 30,
        paragraphs: [
          {
            align: "center",
            lineHeight: 1.25,
            runs: [
              {
                text: "{{participant_name}}",
                fontSize: 18,
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
        x: 42,
        y: 112,
        width: 237,
        height: 16,
        paragraphs: [
          {
            align: "center",
            lineHeight: 1.25,
            runs: [
              {
                text: "{{ticket_name}}",
                fontSize: 10,
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
        x: 140,
        y: 148,
        width: 44,
        height: 44,
        value: "{{checkin_url}}",
        foreground: "#0f172a",
        background: "#ffffff",
        style: "square",
      },
    ],
  },
};
