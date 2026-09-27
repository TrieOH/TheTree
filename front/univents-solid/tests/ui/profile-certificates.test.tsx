import type { JSX } from "@solidjs/web";
import { render, screen } from "@solidjs/testing-library";
import { untrack } from "solid-js";
import { describe, expect, it, vi } from "vitest";
import type { CertificationI } from "@/features/certifications/model";
import { UserCertificationsSection } from "@/features/certifications/ui/UserCertificationsSection";

vi.mock("@tanstack/solid-router", () => ({
  Link: (props: {
    to: string;
    params?: Record<string, string>;
    preload?: string;
    "aria-label"?: string;
    class?: string;
    children?: JSX.Element;
  }) => {
    let href = untrack(() => props.to);
    const params = untrack(() => props.params);
    if (params) {
      for (const [key, value] of Object.entries(params)) {
        href = href.replace(`$${key}`, value);
      }
    }
    return (
      <a
        href={href}
        data-preload={props.preload}
        aria-label={props["aria-label"]}
        class={props.class}
      >
        {props.children}
      </a>
    );
  },
  useNavigate: () => vi.fn(),
}));

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

    expect(screen.getByText("Você ainda não possui certificados")).toBeDefined();
    expect(
      screen.getByText(
        "Seus certificados aparecerão aqui quando forem emitidos.",
      ),
    ).toBeDefined();
  });

  it("renders certificate cards as links directly to /verify/$hash", () => {
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

    expect(certCards[0].getAttribute("href")).toBe("/verify/abcd-1234-ef56");

    expect(certCards[1].getAttribute("href")).toBe("/verify/zzzz-9999-xxxx");
  });
});
