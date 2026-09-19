import { render, screen } from "@solidjs/testing-library";
import { describe, expect, it, vi } from "vitest";
import { allProfileBadges } from "@/features/badges/model";
import { ProfileBadges } from "@/features/badges/ui/ProfileBadges";

vi.mock("@tanstack/solid-router", () => ({
  Link: (props: any) => (
    <a href={props.to} aria-label={props["aria-label"]} class={props.class}>
      {props.children}
    </a>
  ),
  useNavigate: () => vi.fn(),
  useLocation: () => () => ({ pathname: "/", href: "/" }),
}));

const mockBadge = {
  emission_id: "em-1",
  edition_id: "ed-1",
  event_name: "TrieOH Conf",
  edition_name: "Edição 2026",
  ticket_name: "VIP",
  template_name: "Crachá VIP 2026",
  action_url: "https://univents.app/check-in/em-1",
  role: "attendee",
  origin: "participant" as const,
  design_data: {
    canvas: { width: 321, height: 204 },
    backgroundColor: "#ffffff",
    elements: [
      {
        id: "txt-1",
        type: "text" as const,
        x: 32,
        y: 18,
        width: 257,
        height: 20,
        paragraphs: [
          {
            align: "center" as const,
            lineHeight: 1.25,
            runs: [
              {
                text: "{{event_name}}",
                bold: true,
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
        type: "text" as const,
        x: 27,
        y: 78,
        width: 267,
        height: 30,
        paragraphs: [
          {
            align: "center" as const,
            lineHeight: 1.25,
            runs: [
              {
                text: "{{participant_name}}",
                bold: true,
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
        type: "text" as const,
        x: 42,
        y: 114,
        width: 237,
        height: 16,
        paragraphs: [
          {
            align: "center" as const,
            lineHeight: 1.25,
            runs: [
              {
                text: "{{ticket_name}}",
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
        type: "qr" as const,
        x: 140,
        y: 148,
        width: 44,
        height: 44,
        value: "https://univents.app/check-in/em-1",
      },
    ],
  },
};

const mockStaffBadge = {
  ...mockBadge,
  emission_id: "em-2",
  template_name: "Crachá Staff",
  role: "staff",
  origin: "staff" as const,
};

describe("allProfileBadges helper", () => {
  it("combines attendant and staff badges across current and past", () => {
    const groups = {
      attendant: {
        current: [
          {
            ...mockBadge,
            emission_id: "em-1",
          },
        ],
        past: [
          {
            ...mockBadge,
            emission_id: "em-2",
          },
        ],
      },
      staff: {
        current: [
          {
            ...mockBadge,
            emission_id: "em-3",
          },
        ],
        past: [],
      },
    };

    const combined = allProfileBadges(groups);
    expect(combined).toHaveLength(3);
    expect(combined.map((b) => b.emission_id)).toEqual([
      "em-1",
      "em-2",
      "em-3",
    ]);
  });

  it("handles null or undefined input gracefully", () => {
    expect(allProfileBadges(null as any)).toEqual([]);
    expect(allProfileBadges(undefined as any)).toEqual([]);
  });
});

describe("ProfileBadges component", () => {
  it("renders nothing when badge list is empty", () => {
    const { container } = render(() => (
      <ProfileBadges
        badges={[]}
        profileIdentifier="user-1"
        participantName="Maria da Silva"
      />
    ));

    expect(container.firstChild).toBeNull();
  });

  it("renders badge cards with preview and accessible link", () => {
    render(() => (
      <ProfileBadges
        badges={[mockBadge, mockStaffBadge]}
        profileIdentifier="user-1"
        participantName="Maria da Silva"
      />
    ));

    expect(
      screen.getByLabelText("Abrir badge Crachá VIP 2026"),
    ).toBeInTheDocument();
    expect(
      screen.getByLabelText("Abrir badge Crachá Staff"),
    ).toBeInTheDocument();
    expect(screen.getAllByText("TrieOH Conf")).toHaveLength(2);
    expect(screen.getAllByText("Maria da Silva")).toHaveLength(2);
    expect(screen.getAllByText("VIP")).toHaveLength(2);
  });
});
