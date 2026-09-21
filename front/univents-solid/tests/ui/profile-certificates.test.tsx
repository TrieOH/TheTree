import { fireEvent, render, screen } from "@solidjs/testing-library";
import { describe, expect, it, vi } from "vitest";
import type { CertificationI } from "@/features/certifications/model";
import { UserCertificationsSection } from "@/features/certifications/ui/UserCertificationsSection";

vi.mock("@trieoh/front-core-solid", async (importOriginal) => {
  const actual =
    await importOriginal<typeof import("@trieoh/front-core-solid")>();
  return {
    ...actual,
    useQuery: () => () => ({ data: [], isLoading: false, isError: false }),
  };
});

const mockValidCert: CertificationI = {
  id: "cert-1",
  edition_id: "ed-1",
  template_id: null,
  registration_id: "reg-1",
  user_id: "user-1",
  program_id: null,
  verification_hash: "abcd-1234-ef56",
  valid: true,
  invalid_reason: null,
  email_sent: true,
  issued_at: "2026-03-15T12:00:00Z",
  created_at: "2026-03-15T12:00:00Z",
  updated_at: null,
};

const mockInvalidCert: CertificationI = {
  id: "cert-2",
  edition_id: "ed-1",
  template_id: null,
  registration_id: "reg-2",
  user_id: "user-1",
  program_id: null,
  verification_hash: "zzzz-9999-xxxx",
  valid: false,
  invalid_reason: "Cancelado",
  email_sent: false,
  issued_at: "2026-01-10T08:00:00Z",
  created_at: "2026-01-10T08:00:00Z",
  updated_at: null,
};

describe("UserCertificationsSection", () => {
  it("renders empty state when there are no certificates", () => {
    render(() => (
      <UserCertificationsSection
        certifications={[]}
        participantName="Carlos Silva"
      />
    ));

    expect(screen.getByText("Nenhum certificado emitido")).toBeDefined();
    expect(
      screen.getByText(
        "Seus certificados aparecerão aqui quando forem liberados.",
      ),
    ).toBeDefined();
  });

  it("renders certificate preview cards and opens viewer on click", async () => {
    render(() => (
      <UserCertificationsSection
        certifications={[mockValidCert, mockInvalidCert]}
        participantName="Carlos Silva"
      />
    ));

    const certCards = screen.getAllByRole("link", {
      name: /Abrir certificado de/i,
    });
    expect(certCards.length).toBe(2);

    // Clicking a card opens the viewer modal with download actions
    fireEvent.click(certCards[0]);
    expect(await screen.findByText("PNG")).toBeDefined();
    expect(screen.getByText("PDF")).toBeDefined();
    expect(screen.getByText("Fechar")).toBeDefined();
  });
});
