import { cleanup, fireEvent, render, screen } from "@solidjs/testing-library";
import { afterEach, describe, expect, it, vi } from "vitest";

import type { OccurrenceI, ProgramI } from "@/features/programs/model";
import { AdminCreateProgramCard } from "@/features/programs/ui/AdminCreateProgramCard";
import { AdminProgramCard } from "@/features/programs/ui/AdminProgramCard";

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
  price: 0,
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
  it("renders activity program info, thumbnail, metadata, occurrences count and link", () => {
    const onEdit = vi.fn();
    const onDelete = vi.fn();
    const onManageOccurrences = vi.fn();

    render(() => (
      <AdminProgramCard
        program={mockProgram}
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
