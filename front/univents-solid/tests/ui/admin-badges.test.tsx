import { fireEvent, render, screen } from "@solidjs/testing-library";
import { describe, expect, it, vi } from "vitest";

import { AdminBadgeCard, AdminCreateBadgeCard } from "@/features/badges/ui";
import type { BadgePrintItem, BadgeTemplate } from "@/features/badges/model";

const mockTemplate: BadgeTemplate = {
  id: "tmpl-1",
  edition_id: "ed-1",
  ticket_type_id: "tt-1",
  origin: null,
  created_at: "2026-01-01T00:00:00Z",
  updated_at: "2026-01-01T00:00:00Z",
  name: "Crachá VIP 2026",
  design_data: {
    canvas: { width: 321, height: 204 },
    backgroundColor: "#ffffff",
    background: null,
    elements: [
      {
        id: "txt-1",
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
                bold: true,
                italic: false,
                underline: false,
                color: "#64748b",
                fontFamily: "Inter, sans-serif",
                fontSize: 9,
              },
            ],
          },
        ],
      },
      {
        id: "txt-2",
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
                bold: true,
                italic: false,
                underline: false,
                color: "#0f172a",
                fontFamily: "Inter, sans-serif",
                fontSize: 18,
              },
            ],
          },
        ],
      },
      {
        id: "txt-3",
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
                bold: false,
                italic: false,
                underline: false,
                color: "#475569",
                fontFamily: "Inter, sans-serif",
                fontSize: 10,
              },
            ],
          },
        ],
      },
      {
        id: "qr-1",
        type: "qr",
        x: 140,
        y: 148,
        width: 44,
        height: 44,
        value: "https://univents.app/check-in/test",
        foreground: "#000000",
        background: "#ffffff",
        style: "square",
      },
    ],
  },
};

const mockPrintItem: BadgePrintItem = {
  emission_id: "em-1",
  user_id: "u-1",
  origin: "participant",
  action_url: "https://univents.app/check-in/test",
  event_name: "TrieOH Conf",
  edition_name: "Edição 2026",
  ticket_name: "VIP",
  design_data: mockTemplate.design_data,
};

describe("AdminBadgeCard", () => {
  it("renders template details and triggers callbacks", () => {
    const onEdit = vi.fn();
    const onDelete = vi.fn();

    render(() => (
      <AdminBadgeCard
        item={mockTemplate}
        kind="template"
        ticketName="Ingresso VIP"
        onEdit={onEdit}
        onDelete={onDelete}
      />
    ));

    expect(screen.getByText("Crachá VIP 2026")).toBeDefined();
    expect(screen.getAllByText("Ingresso VIP").length).toBeGreaterThan(0);

    // Click card invokes onEdit
    const card = screen.getByRole("button", { name: "Visualizar crachá Crachá VIP 2026" });
    fireEvent.click(card);
    expect(onEdit).toHaveBeenCalledTimes(1);
  });

  it("renders emission card with participant information", () => {
    render(() => (
      <AdminBadgeCard
        item={mockPrintItem}
        kind="emission"
        participantName="Carlos Silva"
      />
    ));

    expect(screen.getAllByText("TrieOH Conf").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Carlos Silva").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Crachá").length).toBeGreaterThan(0);
  });
});

describe("AdminCreateBadgeCard", () => {
  it("renders create card and fires onCreate", () => {
    const onCreate = vi.fn();
    render(() => (
      <AdminCreateBadgeCard
        index={0}
        onCreate={onCreate}
      />
    ));

    const card = screen.getByRole("button", { name: "Novo template" });
    fireEvent.click(card);
    expect(onCreate).toHaveBeenCalledTimes(1);
  });
});
