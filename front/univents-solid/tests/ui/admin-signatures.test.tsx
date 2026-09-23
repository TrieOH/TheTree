import { render, screen, fireEvent, cleanup, waitFor } from "@solidjs/testing-library";
import { describe, it, expect, vi, afterEach } from "vitest";
import {
  AdminCreateSignatureCard,
  AdminCreateSignatureRequestCard,
  AdminSignatureCard,
  AdminSignatureRequestCard,
  SignatureSectionTabs,
} from "@/features/signatures/ui";
import type {
  SignatureI,
  SignatureRequestI,
} from "@/features/signatures/model";

describe("Admin Signatures UI", () => {
  afterEach(() => {
    cleanup();
  });

  const mockSignature: SignatureI = {
    id: "sig-1",
    edition_id: "ed-100",
    created_by: "usr-admin",
    signatory_name: "Dra. Carolina Mendes",
    signatory_title: "Coordenadora Científica",
    signatory_email: "carolina.mendes@evento.com",
    signatory_user_id: null,
    image_url: "https://example.com/signatures/carolina.png",
    created_at: "2026-03-15T10:00:00Z",
    updated_at: null,
    deleted_at: null,
  };

  const mockRequest: SignatureRequestI = {
    id: "req-1",
    edition_id: "ed-100",
    created_by: "usr-admin",
    signatory_name: "Prof. Marcos Silveira",
    signatory_title: "Diretor Geral",
    signatory_email: "marcos.silveira@universidade.edu",
    signatory_user_id: null,
    idempotency_key: "idem-1",
    status: "pending",
    status_reason: null,
    expires_at: "2026-10-01T00:00:00Z",
    signature_id: null,
    created_at: "2026-03-20T10:00:00Z",
    updated_at: null,
  };

  it("renders AdminCreateSignatureCard and handles click", () => {
    const onCreate = vi.fn();
    render(() => <AdminCreateSignatureCard onCreate={onCreate} />);

    expect(screen.getByText("Nova assinatura")).toBeInTheDocument();
    fireEvent.click(screen.getByText("Nova assinatura"));
    expect(onCreate).toHaveBeenCalled();
  });

  it("renders AdminCreateSignatureRequestCard and handles click", () => {
    const onCreate = vi.fn();
    render(() => <AdminCreateSignatureRequestCard onCreate={onCreate} />);

    expect(screen.getByText("Novo convite")).toBeInTheDocument();
    fireEvent.click(screen.getByText("Novo convite"));
    expect(onCreate).toHaveBeenCalled();
  });

  it("renders AdminSignatureCard with correct signatory details", () => {
    render(() => <AdminSignatureCard signature={mockSignature} />);

    expect(screen.getByText("Dra. Carolina Mendes")).toBeInTheDocument();
    expect(screen.getByText("Coordenadora Científica")).toBeInTheDocument();
    expect(screen.getByText("carolina.mendes@evento.com")).toBeInTheDocument();
  });

  it("opens delete modal in AdminSignatureCard when delete button is clicked", async () => {
    const onDelete = vi.fn();
    render(() => (
      <AdminSignatureCard signature={mockSignature} onDelete={onDelete} />
    ));

    const deleteBtn = screen.getByTitle("Excluir assinatura");
    expect(deleteBtn).toBeInTheDocument();
    fireEvent.click(deleteBtn);

    await waitFor(() => {
      expect(
        screen.getByText("Remover assinatura?"),
      ).toBeInTheDocument();
    });
  });

  it("renders AdminSignatureRequestCard with status and signatory info", () => {
    render(() => <AdminSignatureRequestCard request={mockRequest} />);

    expect(screen.getByText("Prof. Marcos Silveira")).toBeInTheDocument();
    expect(screen.getByText("Diretor Geral")).toBeInTheDocument();
    expect(
      screen.getByText("marcos.silveira@universidade.edu"),
    ).toBeInTheDocument();
    expect(screen.getByText("Pendente")).toBeInTheDocument();
  });

  it("renders SignatureSectionTabs and handles tab switching", () => {
    const onChange = vi.fn();
    render(() => (
      <SignatureSectionTabs
        active="signatures"
        onChange={onChange}
        signaturesCount={5}
        invitesCount={2}
      />
    ));

    expect(screen.getByText("Assinaturas")).toBeInTheDocument();
    expect(screen.getByText("Convites de assinatura")).toBeInTheDocument();
    expect(screen.getByText("5")).toBeInTheDocument();
    expect(screen.getByText("2")).toBeInTheDocument();

    const invitesTab = screen.getByText("Convites de assinatura");
    fireEvent.click(invitesTab);
    expect(onChange).toHaveBeenCalledWith("invites");
  });
});
