import { describe, expect, it } from "vitest";
import {
  getAdminBackLink,
  getAdminRouteContext,
  getAdminShellLabel,
  getAdminSidebarSections,
} from "@/widgets/sidebar/sidebar-menu";

describe("Admin Sidebar Menu Context", () => {
  it("extracts route context correctly", () => {
    expect(getAdminRouteContext("/admin/events")).toEqual({});
    expect(getAdminRouteContext("/admin/events/ev_123")).toEqual({
      eventId: "ev_123",
      editionId: undefined,
    });
    expect(
      getAdminRouteContext("/admin/events/ev_123/editions/ed_456"),
    ).toEqual({
      eventId: "ev_123",
      editionId: "ed_456",
    });
  });

  it("returns event-specific sections when on an event page", () => {
    const sections = getAdminSidebarSections("/admin/events/ev_123");
    expect(sections).toHaveLength(1);
    expect(sections[0]?.title).toBe("Evento");
    expect(sections[0]?.items.map((item) => item.id)).toEqual([
      "event-overview",
      "event-editions",
      "event-members",
    ]);
    expect(sections[0]?.items[0]?.params).toEqual({ eventId: "ev_123" });
  });

  it("returns edition-specific sections when on an edition page", () => {
    const sections = getAdminSidebarSections(
      "/admin/events/ev_123/editions/ed_456",
    );
    expect(sections).toHaveLength(1);
    expect(sections[0]?.title).toBe("Edição");
    expect(sections[0]?.items.map((item) => item.id)).toEqual([
      "edition-overview",
      "edition-programs",
      "edition-products",
      "edition-purchases",
      "edition-badges",
      "edition-certifications",
      "edition-signatures",
      "edition-tickets",
    ]);
  });

  it("returns default admin sections when on root admin pages", () => {
    const sections = getAdminSidebarSections("/admin/events");
    expect(sections).toHaveLength(1);
    expect(sections[0]?.title).toBe("Admin");
    expect(sections[0]?.items.map((item) => item.id)).toEqual([
      "events",
      "uploads",
    ]);
  });

  it("returns proper shell labels and back links", () => {
    expect(getAdminShellLabel("/admin/events/ev_123").title).toBe("Evento");
    expect(getAdminBackLink("/admin/events/ev_123")).toEqual({
      to: "/admin/events",
      params: undefined,
    });

    expect(getAdminShellLabel("/admin/events").title).toBe("Eventos");
    expect(getAdminBackLink("/admin/events")).toBeNull();
  });
});
