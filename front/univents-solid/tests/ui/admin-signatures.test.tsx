import { render, screen, fireEvent, cleanup, waitFor } from "@solidjs/testing-library";
import { describe, it, expect, vi, afterEach, beforeEach } from "vitest";
import {
  AdminCreateSignatureCard,
  AdminCreateSignatureRequestCard,
  AdminSignatureCard,
  AdminSignatureRequestCard,
  SignatureSectionTabs,
} from "@/features/signatures/ui";
import { CreateSignatureModal } from "@/features/signatures/ui/CreateSignatureModal";
import { CreateSignatureRequestModal } from "@/features/signatures/ui/CreateSignatureRequestModal";
import type {
  SignatureI,
  SignatureRequestI,
} from "@/features/signatures/model";

const mockCreateSignature = vi.fn().mockResolvedValue({ id: "sig-new" });
const mockCreateSignatureRequest = vi.fn().mockResolvedValue({ id: "req-new" });
const mockUploadFile = vi.fn().mockResolvedValue("https://storage.example.com/signatures/test.png");

vi.mock("@/features/storage/api/index", () => ({
  uploadFile: (...args: unknown[]) => mockUploadFile(...args),
}));

vi.mock("@/features/signatures/api/mutations", () => ({
  useCreateSignatureMutation: () => ({
    mutateAsync: mockCreateSignature,
  }),
  useCreateSignatureRequestMutation: () => ({
    mutateAsync: mockCreateSignatureRequest,
  }),
}));

describe("Admin Signatures UI", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockCreateSignature.mockResolvedValue({ id: "sig-new" });
    mockCreateSignatureRequest.mockResolvedValue({ id: "req-new" });
    mockUploadFile.mockResolvedValue("https://storage.example.com/signatures/test.png");
  });

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

  it("guides through CreateSignatureModal steps, validates fields, and submits successfully", async () => {
    const onOpenChange = vi.fn();
    const onSuccess = vi.fn();

    // Mock canvas methods for jsdom
    const originalGetContext = HTMLCanvasElement.prototype.getContext;
    HTMLCanvasElement.prototype.getContext = vi.fn().mockReturnValue({
      lineWidth: 3.5,
      lineCap: "round",
      lineJoin: "round",
      strokeStyle: "#000",
      clearRect: vi.fn(),
      beginPath: vi.fn(),
      moveTo: vi.fn(),
      lineTo: vi.fn(),
      stroke: vi.fn(),
    }) as any;
    HTMLCanvasElement.prototype.toDataURL = vi.fn().mockReturnValue("data:image/png;base64,mockcanvas");
    HTMLCanvasElement.prototype.toBlob = vi.fn().mockImplementation((cb) => {
      cb(new Blob(["mock-signature"], { type: "image/png" }));
    });

    render(() => (
      <CreateSignatureModal
        open={true}
        onOpenChange={onOpenChange}
        eventId="ev-1"
        editionId="ed-100"
        onSuccess={onSuccess}
      />
    ));

    expect(screen.getByText("Adicionar assinatura")).toBeInTheDocument();

    // Try advancing with empty name
    const continueBtn = screen.getByRole("button", { name: /Continuar/i });
    fireEvent.click(continueBtn);

    await waitFor(() => {
      expect(
        screen.getByText("O nome do signatário deve ter pelo menos 2 caracteres."),
      ).toBeInTheDocument();
    });

    // Fill signatory name and title
    const nameInput = screen.getByLabelText(/Nome do signatário/i);
    fireEvent.input(nameInput, { target: { value: "Prof. Dr. Ricardo Fonseca" } });

    const titleInput = screen.getByLabelText(/Cargo \/ Função/i);
    fireEvent.input(titleInput, { target: { value: "Reitor" } });

    // Step error should be cleared when user types
    await waitFor(() => {
      expect(
        screen.queryByText("O nome do signatário deve ter pelo menos 2 caracteres."),
      ).not.toBeInTheDocument();
    });

    // Advance to Step 2: Assinatura Digital
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByText("Origem da assinatura")).toBeInTheDocument();
    });

    // Switch to upload mode
    const uploadModeBtn = screen.getByRole("button", { name: /Importar imagem/i });
    fireEvent.click(uploadModeBtn);

    // Try to advance without file
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(
        screen.getByText("Selecione um arquivo de imagem com a assinatura."),
      ).toBeInTheDocument();
    });

    // Upload an image file via drop or input
    const file = new File(["dummy-content"], "signature.png", { type: "image/png" });
    const dropzone = document.getElementById("sig-file-upload")?.closest("label") as HTMLElement;
    fireEvent.drop(dropzone, {
      dataTransfer: {
        files: [file],
      },
    });

    await waitFor(() => {
      expect(screen.getByText("Clique para trocar de imagem")).toBeInTheDocument();
    });

    // Advance to Step 3: Resumo
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByText("Dados da assinatura")).toBeInTheDocument();
      expect(screen.getByText("Pronto para cadastrar")).toBeInTheDocument();
      expect(screen.getAllByText("Prof. Dr. Ricardo Fonseca")[0]).toBeInTheDocument();
      expect(screen.getAllByText("Reitor")[0]).toBeInTheDocument();
      expect(screen.getByText("Arquivo de imagem importado")).toBeInTheDocument();
    });

    // Submit form
    const submitBtn = screen.getByRole("button", { name: /Salvar assinatura/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(mockUploadFile).toHaveBeenCalled();
      expect(mockCreateSignature).toHaveBeenCalledWith({
        editionId: "ed-100",
        data: expect.objectContaining({
          signatory_name: "Prof. Dr. Ricardo Fonseca",
          signatory_title: "Reitor",
          image_url: "https://storage.example.com/signatures/test.png",
        }),
      });
      expect(onOpenChange).toHaveBeenCalledWith(false);
      expect(onSuccess).toHaveBeenCalled();
    });

    // Restore canvas
    HTMLCanvasElement.prototype.getContext = originalGetContext;
  });

  it("guides through CreateSignatureRequestModal steps and submits invite", async () => {
    const onOpenChange = vi.fn();
    const onSuccess = vi.fn();

    render(() => (
      <CreateSignatureRequestModal
        open={true}
        onOpenChange={onOpenChange}
        editionId="ed-100"
        onSuccess={onSuccess}
      />
    ));

    expect(screen.getByText("Enviar convite de assinatura")).toBeInTheDocument();

    // Try advancing with empty fields
    const continueBtn = screen.getByRole("button", { name: /Continuar/i });
    fireEvent.click(continueBtn);

    await waitFor(() => {
      expect(
        screen.getByText("O nome do signatário deve ter pelo menos 2 caracteres."),
      ).toBeInTheDocument();
      expect(
        screen.getByText("Informe um e-mail válido para envio do convite."),
      ).toBeInTheDocument();
    });

    // Fill name, email, title and expiration in the same Destinatário step
    const nameInput = screen.getByLabelText(/Nome do signatário/i);
    fireEvent.input(nameInput, { target: { value: "Dra. Beatriz Santos" } });

    const emailInput = screen.getByLabelText(/E-mail para envio/i);
    fireEvent.input(emailInput, { target: { value: "beatriz.santos@hospital.org" } });

    const titleInput = screen.getByLabelText(/Cargo \/ Função/i);
    fireEvent.input(titleInput, { target: { value: "Diretora Clínica" } });

    const expiresInput = screen.getByLabelText(/Validade do link \(dias\)/i);
    fireEvent.input(expiresInput, { target: { value: "14" } });

    // Advance directly to Step 2: Resumo
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByText("Dados do convite")).toBeInTheDocument();
      expect(screen.getByText("Pronto para enviar")).toBeInTheDocument();
      expect(screen.getAllByText("Dra. Beatriz Santos")[0]).toBeInTheDocument();
      expect(screen.getAllByText("beatriz.santos@hospital.org")[0]).toBeInTheDocument();
      expect(screen.getByText("14 dias")).toBeInTheDocument();
    });

    // Submit invite
    const submitBtn = screen.getByRole("button", { name: /Enviar convite/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(mockCreateSignatureRequest).toHaveBeenCalledWith({
        editionId: "ed-100",
        data: expect.objectContaining({
          signatory_name: "Dra. Beatriz Santos",
          signatory_email: "beatriz.santos@hospital.org",
          signatory_title: "Diretora Clínica",
          expires_in_days: 14,
        }),
      });
      expect(onOpenChange).toHaveBeenCalledWith(false);
      expect(onSuccess).toHaveBeenCalled();
    });
  });
});
