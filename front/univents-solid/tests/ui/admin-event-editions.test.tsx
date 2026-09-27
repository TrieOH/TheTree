import type { JSX } from "@solidjs/web";
import { cleanup, fireEvent, render, screen, waitFor } from "@solidjs/testing-library";
import { afterEach, describe, expect, it, vi } from "vitest";

import type { EditionI } from "@/features/editions/model";
import { AdminCreateEditionCard } from "@/features/editions/ui/AdminCreateEditionCard";
import { ManageEditionDialog } from "@/features/editions/ui/ManageEditionDialog";

vi.mock("@tanstack/solid-router", () => ({
  Link: (props: { to: string; "aria-label"?: string; class?: string; children?: JSX.Element }) => (
    <a href={props.to} aria-label={props["aria-label"]} class={props.class}>
      {props.children}
    </a>
  ),
  useNavigate: () => vi.fn(),
  useLocation: () => () => ({ pathname: "/", href: "/" }),
}));

afterEach(cleanup);

const mockEdition: EditionI = {
  id: "ed-1",
  event_id: "event-1",
  name: "Edição 2026",
  slug: "edicao-2026",
  tagline: "O maior evento do ano",
  description: "Descrição da edição",
  is_draft: false,
  registration_opens_at: null,
  starts_at: "2026-10-01T10:00:00.000Z",
  ends_at: "2026-10-05T18:00:00.000Z",
  location_name: "Centro de Convenções",
  location_description: null,
  logo_url: null,
  banner_url: null,
  contact_email: "contato@edicao.com",
  created_by: "user-1",
  created_at: "2026-01-01T12:00:00.000Z",
  updated_at: null,
  deleted_at: null,
  status: "future",
};

const submit = () => {
  const form = document.getElementById("manage-edition-form") as HTMLFormElement;
  form.requestSubmit();
};

describe("AdminCreateEditionCard", () => {
  it("renders create edition card and handles click", () => {
    const onCreate = vi.fn();
    render(() => <AdminCreateEditionCard onCreate={onCreate} animate={false} />);

    const btn = screen.getByRole("button", { name: /Nova edição/i });
    expect(btn).toBeInTheDocument();
    fireEvent.click(btn);
    expect(onCreate).toHaveBeenCalledTimes(1);
  });
});

describe("AdminEditionCard", async () => {
  const { AdminEditionCard } = await import("@/features/editions/ui/AdminEditionCard");

  it("renders edition details and edit button", () => {
    const onEdit = vi.fn();
    render(() => (
      <AdminEditionCard
        edition={mockEdition}
        eventId="event-1"
        onEdit={onEdit}
        animate={false}
      />
    ));

    expect(screen.getByText("Edição 2026")).toBeInTheDocument();
    expect(screen.getByText("Centro de Convenções")).toBeInTheDocument();
    expect(screen.getByText("/edicao-2026")).toBeInTheDocument();

    const editBtn = screen.getByLabelText("Editar Edição 2026");
    fireEvent.click(editBtn);
    expect(onEdit).toHaveBeenCalledWith(mockEdition);
  });
});

describe("ManageEditionDialog", () => {
  it("opens an existing edition with its values and submits them", async () => {
    const onSubmit = vi.fn().mockResolvedValue(true);
    const onOpenChange = vi.fn();

    render(() => (
      <ManageEditionDialog
        open={true}
        onOpenChange={onOpenChange}
        edition={mockEdition}
        onSubmit={onSubmit}
      />
    ));

    expect(screen.getByRole("heading", { name: "Editar edição" })).toBeInTheDocument();

    submit();

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledTimes(1);
      expect(onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          name: "Edição 2026",
          slug: "edicao-2026",
          location_name: "Centro de Convenções",
        }),
      );
    });
  });

  it("refuses to submit an invalid form when creating", async () => {
    const onSubmit = vi.fn().mockResolvedValue(true);
    const onOpenChange = vi.fn();

    render(() => (
      <ManageEditionDialog
        open={true}
        onOpenChange={onOpenChange}
        edition={null}
        onSubmit={onSubmit}
      />
    ));

    // Tentativa de submissão no passo 0 (Edição) sem preenchimento
    submit();

    await waitFor(() => {
      expect(onSubmit).not.toHaveBeenCalled();
      expect(screen.getAllByText("Informe ao menos 2 caracteres.")).toHaveLength(2);
    });

    // Preenche o passo 0
    const nameInput = screen.getByPlaceholderText("Ex: Edição 2026");
    fireEvent.input(nameInput, { target: { value: "Edição 2027" } });

    const continueBtn = screen.getByRole("button", { name: "Continuar" });
    fireEvent.click(continueBtn);

    // Passo 1 (Cronograma): tenta avançar sem data de início
    await waitFor(() => {
      expect(screen.getByText("Cronograma")).toBeInTheDocument();
    });

    submit();

    await waitFor(() => {
      expect(onSubmit).not.toHaveBeenCalled();
      expect(screen.getByText("Informe a data de início.")).toBeInTheDocument();
    });
  });
});
