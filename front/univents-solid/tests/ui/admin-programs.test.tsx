import { cleanup, fireEvent, render, screen } from "@solidjs/testing-library";
import { afterEach, describe, expect, it, vi } from "vitest";

import type { OccurrenceI, ProgramI } from "@/features/programs/model";
import { AdminCreateProgramCard } from "@/features/programs/ui/AdminCreateProgramCard";
import { AdminProgramCard } from "@/features/programs/ui/AdminProgramCard";

afterEach(cleanup);

const mockProgram: ProgramI = {
  id: "prog-1",
  edition_id: "ed-1",
  name: "Keynote de Abertura",
  description: "Apresentação de abertura com os fundadores e convidados especiais.",
  kind: "activity",
  min_access_level: 1,
  staff_only: false,
  price: 0,
  banner_url: "https://example.com/keynote.png",
  created_at: "2026-03-01T10:00:00Z",
  updated_at: null,
  deleted_at: null,
};

const mockCheckpointProgram: ProgramI = {
  id: "prog-2",
  edition_id: "ed-1",
  name: "Checkpoint Entrada Principal",
  description: "Validação de crachás no hall principal.",
  kind: "checkpoint",
  min_access_level: 0,
  staff_only: true,
  price: 0,
  banner_url: null,
  created_at: "2026-03-01T10:00:00Z",
  updated_at: null,
  deleted_at: null,
};

const mockOccurrences: OccurrenceI[] = [
  {
    id: "occ-1",
    program_id: "prog-1",
    edition_id: "ed-1",
    starts_at: "2026-03-20T09:00:00Z",
    ends_at: "2026-03-20T10:30:00Z",
    max_capacity: 150,
    created_at: "2026-03-01T10:00:00Z",
    updated_at: null,
    deleted_at: null,
  },
  {
    id: "occ-2",
    program_id: "prog-1",
    edition_id: "ed-1",
    starts_at: "2026-03-21T09:00:00Z",
    ends_at: "2026-03-21T10:30:00Z",
    max_capacity: 150,
    created_at: "2026-03-01T10:00:00Z",
    updated_at: null,
    deleted_at: null,
  },
];

describe("AdminCreateProgramCard", () => {
  it("renders with title, description and triggers onCreate callback", () => {
    const onCreate = vi.fn();
    render(() => <AdminCreateProgramCard onCreate={onCreate} animate={false} />);

    const btn = screen.getByRole("button", { name: /Novo programa/i });
    expect(btn).toBeInTheDocument();
    expect(
      screen.getByText("Crie atividades, palestras ou checkpoints com horários e controle de acesso."),
    ).toBeInTheDocument();

    fireEvent.click(btn);
    expect(onCreate).toHaveBeenCalledTimes(1);
  });
});

describe("AdminProgramCard", () => {
  it("renders activity program info, thumbnail, badges, occurrences count and actions", () => {
    const onEdit = vi.fn();
    const onDelete = vi.fn();
    const onOpenCalendar = vi.fn();

    render(() => (
      <AdminProgramCard
        program={mockProgram}
        occurrences={mockOccurrences}
        animate={false}
        onEdit={onEdit}
        onDelete={onDelete}
        onOpenCalendar={onOpenCalendar}
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

    const calBtn = screen.getByRole("button", { name: /Ver no calendário/i });
    fireEvent.click(calBtn);
    expect(onOpenCalendar).toHaveBeenCalledWith(mockProgram);
  });

  it("renders checkpoint program with staff only badge and fallback icon", () => {
    const onEdit = vi.fn();
    const onDelete = vi.fn();
    const onOpenCalendar = vi.fn();

    render(() => (
      <AdminProgramCard
        program={mockCheckpointProgram}
        occurrences={[]}
        animate={false}
        onEdit={onEdit}
        onDelete={onDelete}
        onOpenCalendar={onOpenCalendar}
      />
    ));

    expect(screen.getByText("Checkpoint Entrada Principal")).toBeInTheDocument();
    expect(screen.getByText("Checkpoint")).toBeInTheDocument();
    expect(screen.getByText("Staff")).toBeInTheDocument();
    expect(screen.getByText("0 ocorrências")).toBeInTheDocument();
  });
});
