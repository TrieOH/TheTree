import { cleanup, fireEvent, render, screen, waitFor } from "@solidjs/testing-library";
import { afterEach, describe, expect, it, vi } from "vitest";

import type { OccurrenceI, ProgramI } from "@/features/programs/model";
import { AdminCreateProgramCard } from "@/features/programs/ui/AdminCreateProgramCard";
import { AdminProgramCard } from "@/features/programs/ui/AdminProgramCard";
import { ManageProgramDialog } from "@/features/programs/ui/ManageProgramDialog";

afterEach(() => {
  cleanup();
});

const mockProgram: ProgramI = {
  id: "prog-1",
  edition_id: "ed-1",
  kind: "activity",
  name: "Keynote de Abertura",
  description: "Apresentação de abertura com os fundadores e convidados especiais.",
  min_access_level: 1,
  staff_only: false,
  banner_url: "https://example.com/keynote.png",
  price: 5000,
  created_at: "2026-09-01T10:00:00Z",
  updated_at: "2026-09-01T10:00:00Z",
  deleted_at: null,
};

const mockCheckpointProgram: ProgramI = {
  id: "prog-2",
  edition_id: "ed-1",
  kind: "checkpoint",
  name: "Checkpoint Entrada Principal",
  description: "Validação de crachás no hall principal.",
  min_access_level: 0,
  staff_only: true,
  banner_url: null,
  price: 0,
  created_at: "2026-09-01T10:00:00Z",
  updated_at: "2026-09-01T10:00:00Z",
  deleted_at: null,
};

const mockOccurrences: OccurrenceI[] = [
  {
    id: "occ-1",
    program_id: "prog-1",
    edition_id: "ed-1",
    starts_at: "2026-09-10T09:00:00Z",
    ends_at: "2026-09-10T10:30:00Z",
    max_capacity: 100,
    created_at: "2026-09-01T10:00:00Z",
    updated_at: "2026-09-01T10:00:00Z",
    deleted_at: null,
  },
  {
    id: "occ-2",
    program_id: "prog-1",
    edition_id: "ed-1",
    starts_at: "2026-09-10T14:00:00Z",
    ends_at: "2026-09-10T15:30:00Z",
    max_capacity: 100,
    created_at: "2026-09-01T10:00:00Z",
    updated_at: "2026-09-01T10:00:00Z",
    deleted_at: null,
  },
];

describe("AdminCreateProgramCard", () => {
  it("renders with title, description and triggers onCreate callback", () => {
    const onCreate = vi.fn();
    render(() => <AdminCreateProgramCard onCreate={onCreate} animate={false} />);

    expect(screen.getByText("Novo programa")).toBeInTheDocument();
    expect(
      screen.getByText(
        "Crie atividades, palestras ou checkpoints com horários e controle de acesso.",
      ),
    ).toBeInTheDocument();

    const button = screen.getByRole("button", { name: /Novo programa/i });
    fireEvent.click(button);
    expect(onCreate).toHaveBeenCalledTimes(1);
  });
});

describe("AdminProgramCard", () => {
  it("renders activity program info, thumbnail with upload trigger, metadata, occurrences count and link", () => {
    const onEdit = vi.fn();
    const onDelete = vi.fn();
    const onManageOccurrences = vi.fn();

    render(() => (
      <AdminProgramCard
        program={mockProgram}
        eventId="event-1"
        occurrences={mockOccurrences}
        occurrencesHref="/admin/events/1/editions/1/programs/prog-1/occurrences"
        animate={false}
        onEdit={onEdit}
        onDelete={onDelete}
        onManageOccurrences={onManageOccurrences}
      />
    ));

    expect(screen.getByText("Keynote de Abertura")).toBeInTheDocument();
    expect(
      screen.getByText("Apresentação de abertura com os fundadores e convidados especiais."),
    ).toBeInTheDocument();
    expect(screen.getByText("Atividade")).toBeInTheDocument();
    expect(screen.getByText("Nível 1")).toBeInTheDocument();
    expect(screen.getByText("2 ocorrências")).toBeInTheDocument();

    const img = screen.getByRole("img", { name: "Keynote de Abertura" });
    expect(img).toHaveAttribute("src", "https://example.com/keynote.png");

    // Image card interactive upload trigger and hidden file input
    const imageCardTrigger = screen.getByRole("button", {
      name: "Alterar imagem de Keynote de Abertura",
    });
    expect(imageCardTrigger).toBeInTheDocument();

    const fileInput = document.getElementById("program-prog-1-banner-upload");
    expect(fileInput).toBeInTheDocument();
    expect(fileInput).toHaveAttribute("type", "file");
    expect(fileInput).toHaveAttribute("accept", "image/*");

    const editBtn = screen.getByLabelText("Editar Keynote de Abertura");
    fireEvent.click(editBtn);
    expect(onEdit).toHaveBeenCalledWith(mockProgram);

    const deleteBtn = screen.getByLabelText("Excluir Keynote de Abertura");
    fireEvent.click(deleteBtn);
    expect(onDelete).toHaveBeenCalledWith(mockProgram);

    const occLink = screen.getByRole("link", { name: /Ocorrências/i });
    expect(occLink).toHaveAttribute("href", "/admin/events/1/editions/1/programs/prog-1/occurrences");
    fireEvent.click(occLink);
    expect(onManageOccurrences).toHaveBeenCalledWith(mockProgram);
  });

  it("renders checkpoint program with staff only badge and fallback icon", () => {
    const onEdit = vi.fn();
    const onDelete = vi.fn();

    render(() => (
      <AdminProgramCard
        program={mockCheckpointProgram}
        occurrences={[]}
        animate={false}
        onEdit={onEdit}
        onDelete={onDelete}
      />
    ));

    expect(screen.getByText("Checkpoint Entrada Principal")).toBeInTheDocument();
    expect(screen.getByText("Checkpoint")).toBeInTheDocument();
    expect(screen.getByText("Staff")).toBeInTheDocument();
    expect(screen.getByText("0 ocorrências")).toBeInTheDocument();
  });
});

describe("ManageProgramDialog", () => {
  it("maintains focus during typing and guides through creation of an activity without banner field in modal", async () => {
    const onOpenChange = vi.fn();
    const onSubmit = vi.fn().mockResolvedValue(true);

    render(() => (
      <ManageProgramDialog
        open={true}
        onOpenChange={onOpenChange}
        onSubmit={onSubmit}
      />
    ));

    expect(screen.getByText("Nova programação")).toBeInTheDocument();

    // Step 1: Validation failure on empty name
    const continueBtn = screen.getByRole("button", { name: /Continuar/i });
    fireEvent.click(continueBtn);

    await waitFor(() => {
      expect(
        screen.getByText("Nome deve ter pelo menos 2 caracteres."),
      ).toBeInTheDocument();
    });

    // Fill Step 1 fields and check focus retention
    const nameInput = screen.getByLabelText(/Nome da programação/i) as HTMLInputElement;
    nameInput.focus();
    expect(document.activeElement).toBe(nameInput);

    fireEvent.input(nameInput, { target: { value: "Workshop de Rust" } });
    expect(document.activeElement).toBe(nameInput);

    // Step error should be cleared when user types
    await waitFor(() => {
      expect(
        screen.queryByText("Nome deve ter pelo menos 2 caracteres."),
      ).not.toBeInTheDocument();
    });

    const descInput = screen.getByLabelText(/Descrição \(opcional\)/i);
    fireEvent.input(descInput, { target: { value: "Do básico ao avançado em sistemas." } });

    // Banner URL input is not present in modal steps
    expect(screen.queryByPlaceholderText("https://exemplo.com/imagem.png")).not.toBeInTheDocument();

    // Advance to Step 2: Acesso e Valores
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByLabelText(/Valor adicional/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/Nível de acesso mínimo/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/Apenas para equipe/i)).toBeInTheDocument();
    });

    // Fill price and min_access_level
    const priceInput = screen.getByLabelText(/Valor adicional/i);
    fireEvent.input(priceInput, { target: { value: "15000" } }); // R$ 150,00

    const levelInput = screen.getByLabelText(/Nível de acesso mínimo/i);
    fireEvent.input(levelInput, { target: { value: "2" } });

    // Advance to Step 3: Resumo
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByText("Dados da programação")).toBeInTheDocument();
      expect(screen.getByText("Pronto para cadastrar")).toBeInTheDocument();
      expect(screen.getAllByText("Workshop de Rust")[0]).toBeInTheDocument();
      expect(screen.getAllByText("Atividade / Palestra")[0]).toBeInTheDocument();
      expect(screen.getAllByText("Nível 2")[0]).toBeInTheDocument();
      expect(screen.getAllByText("R$ 150,00")[0]).toBeInTheDocument();
    });

    // Submit form
    const submitBtn = screen.getByRole("button", { name: /Criar programação/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        kind: "activity",
        name: "Workshop de Rust",
        description: "Do básico ao avançado em sistemas.",
        min_access_level: 2,
        staff_only: false,
        banner_url: null,
        price: 15000,
      });
      expect(onOpenChange).toHaveBeenCalledWith(false);
    });
  });

  it("creates a checkpoint with full-width staff restriction and no price field", async () => {
    const onOpenChange = vi.fn();
    const onSubmit = vi.fn().mockResolvedValue(true);

    render(() => (
      <ManageProgramDialog
        open={true}
        onOpenChange={onOpenChange}
        onSubmit={onSubmit}
      />
    ));

    // Select checkpoint kind
    const checkpointBtn = screen.getByRole("button", { name: /Checkpoint de Presença/i });
    fireEvent.click(checkpointBtn);

    const nameInput = screen.getByLabelText(/Nome da programação/i);
    fireEvent.input(nameInput, { target: { value: "Credenciamento Portaria A" } });

    // Advance to Step 2
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      // Checkpoint does not have price field
      expect(screen.queryByLabelText(/Valor adicional/i)).not.toBeInTheDocument();
      expect(screen.getByLabelText(/Nível de acesso mínimo/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/Apenas para equipe/i)).toBeInTheDocument();
    });

    // Toggle staff only
    const staffCheckbox = screen.getByLabelText(/Apenas para equipe/i);
    fireEvent.click(staffCheckbox);

    // Advance to Step 3: Resumo
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getAllByText("Credenciamento Portaria A")[0]).toBeInTheDocument();
      expect(screen.getAllByText("Checkpoint de Presença")[0]).toBeInTheDocument();
      expect(screen.getByText("Não aplicável (Checkpoint)")).toBeInTheDocument();
      expect(screen.getByText("Exclusivo para equipe (Staff only)")).toBeInTheDocument();
    });

    // Submit form
    const submitBtn = screen.getByRole("button", { name: /Criar programação/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        kind: "checkpoint",
        name: "Credenciamento Portaria A",
        description: undefined,
        min_access_level: 0,
        staff_only: true,
        banner_url: null,
        price: undefined,
      });
    });
  });

  it("populates existing program data when editing and preserves existing banner_url without clearing it", async () => {
    const onOpenChange = vi.fn();
    const onSubmit = vi.fn().mockResolvedValue(true);

    render(() => (
      <ManageProgramDialog
        open={true}
        program={mockProgram}
        onOpenChange={onOpenChange}
        onSubmit={onSubmit}
      />
    ));

    expect(screen.getByText("Editar programação")).toBeInTheDocument();

    const nameInput = screen.getByLabelText(/Nome da programação/i) as HTMLInputElement;
    expect(nameInput.value).toBe("Keynote de Abertura");

    fireEvent.input(nameInput, { target: { value: "Keynote Principal de Abertura" } });

    // Advance to Step 2
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByLabelText(/Valor adicional/i)).toBeInTheDocument();
    });

    // Advance to Step 3
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByText("Pronto para salvar alterações")).toBeInTheDocument();
      expect(screen.getAllByText("Keynote Principal de Abertura")[0]).toBeInTheDocument();
    });

    const submitBtn = screen.getByRole("button", { name: /Salvar alterações/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        kind: "activity",
        name: "Keynote Principal de Abertura",
        description: "Apresentação de abertura com os fundadores e convidados especiais.",
        min_access_level: 1,
        staff_only: false,
        banner_url: "https://example.com/keynote.png",
        price: 5000,
      });
    });
  });
});
