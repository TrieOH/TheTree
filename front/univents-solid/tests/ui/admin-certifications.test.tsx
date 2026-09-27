import type { JSX } from "@solidjs/web";
import { render, screen, fireEvent, cleanup } from "@solidjs/testing-library";
import { describe, it, expect, vi, afterEach } from "vitest";
import {
  CertificationList,
  CertificationEmissionErrorsList,
} from "@/features/certifications/ui/CertificationLists";
import type {
  CertificationI,
  CertificationEmissionErrorI,
} from "@/features/certifications/model";

vi.mock("@tanstack/solid-router", () => ({
  Link: (props: { href: string; children: JSX.Element }) => (
    <a href={props.href}>{props.children}</a>
  ),
  useNavigate: () => vi.fn(),
}));

vi.mock("@trieoh/identityx-sdk-ts-solid", () => ({
  useAuth: () => ({
    auth: {
      getActorProfile: vi.fn().mockResolvedValue({
        success: true,
        data: {
          profile: {
            display_name: "Carolina Santos",
          },
        },
      }),
    },
  }),
}));

const testCerts: CertificationI[] = [
  {
    id: "cert-1",
    edition_id: "ed-456",
    template_id: "tmpl-1",
    registration_id: "reg-1",
    user_id: "usr_carolina_santos",
    program_id: null,
    verification_hash: "UNIV-2026-A8F2-7C3B",
    valid: true,
    invalid_reason: null,
    email_sent: true,
    issued_at: new Date().toISOString(),
    created_at: new Date().toISOString(),
    updated_at: null,
  },
  {
    id: "cert-2",
    edition_id: "ed-456",
    template_id: "tmpl-1",
    registration_id: "reg-2",
    user_id: "usr_rodrigo_melo",
    program_id: null,
    verification_hash: "UNIV-2026-F9D1-4E2A",
    valid: true,
    invalid_reason: null,
    email_sent: true,
    issued_at: new Date().toISOString(),
    created_at: new Date().toISOString(),
    updated_at: null,
  },
];

const testErrors: CertificationEmissionErrorI[] = [
  {
    id: "err-1",
    edition_id: "ed-456",
    user_id: "usr_carolina_santos",
    template_id: "tmpl-1",
    program_id: "prog-1",
    error_message: "Presença mínima de 75% não atingida pelo participante.",
    created_at: new Date().toISOString(),
  },
  {
    id: "err-2",
    edition_id: "ed-456",
    user_id: "usr_rodrigo_melo",
    template_id: "tmpl-1",
    program_id: null,
    error_message: "Assinatura digital do coordenador não configurada.",
    created_at: new Date().toISOString(),
  },
  {
    id: "err-3",
    edition_id: "ed-456",
    user_id: "usr_beatriz_lima",
    template_id: "tmpl-1",
    program_id: null,
    error_message: "Reembolso solicitado antes da conclusão do evento.",
    created_at: new Date().toISOString(),
  },
];

vi.mock("@trieoh/front-core-solid", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@trieoh/front-core-solid")>();
  return {
    ...actual,
    useQuery: (optionsFn: () => { queryKey?: unknown[] }) => () => {
      const opts = typeof optionsFn === "function" ? optionsFn() : optionsFn;
      const key = (opts as { queryKey?: unknown[] })?.queryKey;
      if (Array.isArray(key)) {
        if (key[0] === "certifications" && key[1] === "issued") {
          return { data: testCerts, isLoading: false, isError: false };
        }
        if (key[0] === "certifications" && key[1] === "emission-errors") {
          return { data: testErrors, isLoading: false, isError: false };
        }
      }
      return { data: [], isLoading: false, isError: false };
    },
  };
});

vi.mock("@/features/certifications/api/mutations", () => ({
  useInvalidateCertificationMutation: () => ({
    mutateAsync: vi.fn(),
  }),
}));

describe("CertificationLists", () => {
  afterEach(() => {
    cleanup();
  });

  it("renders CertificationList with sort options and certificates", () => {
    const { container } = render(() => (
      <CertificationList eventId="ev-123" editionId="ed-456" />
    ));

    expect(
      screen.getByPlaceholderText("Buscar por participante, hash ou atividade..."),
    ).toBeDefined();

    // Verify sort button is present in PaginatedContainer
    const sortButton = container.querySelector('button[aria-label="Ordenar"]');
    expect(sortButton).not.toBeNull();

    // Verify certificates are rendered
    expect(screen.getByText("UNIV-2026-A8F2-7C3B")).toBeDefined();
    expect(screen.getByText("UNIV-2026-F9D1-4E2A")).toBeDefined();
  });

  it("renders CertificationEmissionErrorsList with PaginatedContainer and errors", () => {
    const { container } = render(() => (
      <CertificationEmissionErrorsList editionId="ed-456" />
    ));

    expect(
      screen.getByPlaceholderText("Buscar por participante, erro ou ID..."),
    ).toBeDefined();

    // Verify sort button is present in PaginatedContainer
    const sortButton = container.querySelector('button[aria-label="Ordenar"]');
    expect(sortButton).not.toBeNull();

    // Verify errors are rendered
    expect(screen.getByText(/Presença mínima de 75%/i)).toBeDefined();
    expect(screen.getByText(/Assinatura digital do coordenador/i)).toBeDefined();
    expect(screen.getAllByText("Falha de emissão").length).toBeGreaterThan(0);
  });

  it("filters emission errors according to search query", async () => {
    render(() => <CertificationEmissionErrorsList editionId="ed-456" />);

    const searchInput = screen.getByPlaceholderText(
      "Buscar por participante, erro ou ID...",
    ) as HTMLInputElement;

    fireEvent.input(searchInput, { target: { value: "reembolso" } });
    await new Promise((resolve) => setTimeout(resolve, 50));

    expect(screen.getByText(/reembolso solicitado/i)).toBeDefined();
    expect(screen.queryByText(/Presença mínima de 75%/i)).toBeNull();
  });
});
