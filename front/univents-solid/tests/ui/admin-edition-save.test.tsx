import type { JSX } from "@solidjs/web";
import { cleanup, fireEvent, render, screen, waitFor } from "@solidjs/testing-library";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

const state = vi.hoisted(() => ({
  editions: [] as unknown[],
  create: vi.fn(), patch: vi.fn(), success: vi.fn(), error: vi.fn(),
}));

vi.mock("@tanstack/solid-router", async (original) => ({
  ...await original<typeof import("@tanstack/solid-router")>(),
  createFileRoute: () => (options: unknown) => ({ options, useParams: () => () => ({ eventId: "event-1" }) }),
  Link: (props: { children?: JSX.Element }) => <a>{props.children}</a>,
  useNavigate: () => vi.fn(),
  useLocation: () => () => ({ pathname: "/", href: "/" }),
}));
vi.mock("@trieoh/front-core-solid", async (original) => ({
  ...await original<typeof import("@trieoh/front-core-solid")>(),
  useQuery: () => () => ({ data: state.editions, isLoading: false }),
}));
vi.mock("@/features/editions/api", () => ({ allAdminEditionsQueryOptions: () => ({}) }));
vi.mock("@/features/editions/api/mutations", () => ({
  useCreateEditionMutation: () => ({ mutateAsync: state.create }),
  usePatchEditionMutation: () => ({ mutateAsync: state.patch }),
}));
vi.mock("@/shared/ui/toast", () => ({ toast: { success: state.success, error: state.error } }));

import { Route } from "@/routes/admin/events/$eventId/editions/index";
const Page = Route.options.component as () => JSX.Element;
const edition = {
  id: "ed-1", event_id: "event-1", name: "Edição 2026", slug: "edicao-2026",
  starts_at: "2026-10-01T10:00:00.000Z", ends_at: "2026-10-05T18:00:00.000Z",
  created_at: "2026-01-01T12:00:00.000Z", status: "future", is_draft: false,
};

afterEach(cleanup);
beforeEach(async () => {
  vi.clearAllMocks();
  state.editions = [];
  state.create.mockReset().mockResolvedValue(edition);
  state.patch.mockReset().mockResolvedValue(edition);
  await (Page as typeof Page & { preload?: () => Promise<void> }).preload?.();
});

async function openCreateAndFill() {
  fireEvent.click(screen.getAllByRole("button", { name: /Nova edição/i })[0]!);
  await fireEvent.input(await screen.findByPlaceholderText("Ex: Edição 2026"), { target: { value: "Edição 2026" } });
  await waitFor(() => expect(screen.getByPlaceholderText("edicao-2026")).toHaveValue("edicao-2026"));
  fireEvent.click(screen.getByRole("button", { name: "Continuar" }));
  await waitFor(() => expect(screen.getByLabelText(/Início/)).toBeInTheDocument());
  await fireEvent.input(screen.getByLabelText(/Início/), { target: { value: "2026-10-01T10:00" } });
  await fireEvent.input(screen.getByLabelText(/Término/), { target: { value: "2026-10-05T18:00" } });
}

function submit() { fireEvent.submit(document.getElementById("manage-edition-form")!); }

describe("edition save notifications and writes", () => {
  it("creates with exactly one mutation and one success notification", async () => {
    render(() => <Page />);
    await openCreateAndFill();
    submit();
    await waitFor(() => expect(state.success).toHaveBeenCalled());
    expect(state.create).toHaveBeenCalledTimes(1);
    expect(state.patch).not.toHaveBeenCalled();
    expect(state.success).toHaveBeenCalledExactlyOnceWith("Edição criada com sucesso!");
  });

  it("updates with exactly one mutation and one success notification", async () => {
    state.editions = [edition];
    render(() => <Page />);
    fireEvent.click(screen.getByLabelText("Editar Edição 2026"));
    await screen.findByRole("heading", { name: "Editar edição" });
    submit();
    await waitFor(() => expect(state.success).toHaveBeenCalled());
    expect(state.patch).toHaveBeenCalledTimes(1);
    expect(state.create).not.toHaveBeenCalled();
    expect(state.success).toHaveBeenCalledExactlyOnceWith("Edição atualizada com sucesso!");
  });

  it("reports a failed update once and keeps the editor open", async () => {
    state.editions = [edition];
    state.patch.mockRejectedValue(new Error("Falha ao gravar"));
    render(() => <Page />);
    fireEvent.click(screen.getByLabelText("Editar Edição 2026"));
    await screen.findByRole("heading", { name: "Editar edição" });
    submit();
    await waitFor(() => expect(state.error).toHaveBeenCalled());
    expect(state.patch).toHaveBeenCalledTimes(1);
    expect(state.error).toHaveBeenCalledTimes(1);
    expect(state.success).not.toHaveBeenCalled();
    expect(screen.getByRole("dialog")).toBeInTheDocument();
  });

  it("sends one update when the form is submitted repeatedly before completion", async () => {
    state.editions = [edition];
    let finish!: (value: typeof edition) => void;
    state.patch.mockImplementation(() => new Promise((resolve) => { finish = resolve; }));
    render(() => <Page />);
    fireEvent.click(screen.getByLabelText("Editar Edição 2026"));
    await screen.findByRole("heading", { name: "Editar edição" });
    submit();
    submit();
    try {
      expect(state.patch).toHaveBeenCalledTimes(1);
    } finally {
      finish(edition);
    }
    await waitFor(() => expect(state.success).toHaveBeenCalledTimes(1));
    expect(state.patch).toHaveBeenCalledTimes(1);
  });

  it("allows a new attempt after a failed update", async () => {
    state.editions = [edition];
    state.patch.mockRejectedValueOnce(new Error("Falha temporária"));
    render(() => <Page />);
    fireEvent.click(screen.getByLabelText("Editar Edição 2026"));
    await screen.findByRole("heading", { name: "Editar edição" });
    submit();
    await waitFor(() => expect(state.error).toHaveBeenCalledTimes(1));
    await waitFor(() => expect(screen.getByRole("button", { name: "Continuar" })).toBeEnabled());
    submit();
    await waitFor(() => expect(state.success).toHaveBeenCalledTimes(1));
    expect(state.patch).toHaveBeenCalledTimes(2);
    expect(state.error).toHaveBeenCalledTimes(1);
  });

  it("reports overlapping dates once without writing", async () => {
    state.editions = [edition];
    render(() => <Page />);
    await openCreateAndFill();
    submit();
    await waitFor(() => expect(state.error).toHaveBeenCalled());
    expect(state.error).toHaveBeenCalledExactlyOnceWith("Já existe uma edição cadastrada nesse período.");
    expect(state.create).not.toHaveBeenCalled();
    expect(state.patch).not.toHaveBeenCalled();
  });
});
