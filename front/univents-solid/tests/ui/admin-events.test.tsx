import { cleanup, fireEvent, render, screen, waitFor } from "@solidjs/testing-library";
import { afterEach, describe, expect, it, vi } from "vitest";

import { PaginatedContainer } from "@trieoh/ui-solid";

import { ManageEventDialog, type ManageEventValues } from "@/features/events/ui/ManageEventDialog";
import type { EventI } from "@/features/events/model";

// The dialog renders through a portal into `document.body`, and this library
// does not unmount for us, so each case would otherwise stack on the previous.
afterEach(cleanup);

/** Async handler defined at module scope: inline `async` arrows in JSX trip a diagnostic. */
const acceptSubmit = async (): Promise<boolean> => true;
const rejectSubmit = async (): Promise<boolean> => false;

const submitSpy = () => vi.fn<(values: ManageEventValues) => Promise<boolean>>(acceptSubmit);

const event = {
  id: "event-1",
  full_name: "Semana de Tecnologia",
  slug: "semana-de-tecnologia",
  acronym: "ST",
  description: "Uma semana de palestras.",
  contact_email: "contato@example.test",
  status: "draft",
  created_at: "2026-01-01T00:00:00.000Z",
} as EventI;

/** Labels carry the required marker, hence a regex instead of an exact string. */
const field = (name: string) => screen.getByLabelText(new RegExp(`^${name}`));

const typeInto = (name: string, value: string) => {
  const input = field(name);
  fireEvent.input(input, { target: { value } });
  return input;
};

const submit = () => {
  const form = document.getElementById("manage-event-form") as HTMLFormElement;
  form.requestSubmit();
};

describe("ManageEventDialog", () => {
  it("derives slug and acronym from the name", async () => {
    render(() => (
      <ManageEventDialog open event={null} onOpenChange={() => undefined} onSubmit={acceptSubmit} />
    ));

    const name = typeInto("Nome", "Semana de Tecnologia");

    expect(name).toHaveValue("Semana de Tecnologia");
    await waitFor(() => {
      expect(field("Slug")).toHaveValue("semana-de-tecnologia");
      expect(field("Sigla")).toHaveValue("ST");
    });
  });

  it("stops deriving once the slug is edited by hand", async () => {
    render(() => (
      <ManageEventDialog open event={null} onOpenChange={() => undefined} onSubmit={acceptSubmit} />
    ));

    typeInto("Nome", "Semana de Tecnologia");
    typeInto("Slug", "st-2026");
    typeInto("Nome", "Semana de Tecnologia e Arte");

    expect(field("Slug")).toHaveValue("st-2026");
  });

  it("refuses to submit an invalid form", async () => {
    const onSubmit = submitSpy();
    render(() => (
      <ManageEventDialog open event={null} onOpenChange={() => undefined} onSubmit={onSubmit} />
    ));

    submit();

    await waitFor(() => expect(screen.getAllByRole("alert").length).toBeGreaterThan(0));
    expect(onSubmit).not.toHaveBeenCalled();
  });

  /**
   * Typing over a pre-filled field is covered by the Playwright e2e: in jsdom the
   * delegated `input` event does not reach the portaled dialog's subtree when the
   * field starts with a value, which says more about the harness than the app.
   */
  it("opens an existing event with its values and submits them", async () => {
    const onSubmit = submitSpy();
    render(() => (
      <ManageEventDialog open event={event} onOpenChange={() => undefined} onSubmit={onSubmit} />
    ));

    expect(field("Nome")).toHaveValue("Semana de Tecnologia");
    expect(field("Slug")).toHaveValue("semana-de-tecnologia");
    expect(field("Sigla")).toHaveValue("ST");

    submit();

    await waitFor(() => expect(onSubmit).toHaveBeenCalledTimes(1));
    expect(onSubmit.mock.calls[0]?.[0]).toMatchObject({
      full_name: "Semana de Tecnologia",
      slug: "semana-de-tecnologia",
      acronym: "ST",
      contact_email: "contato@example.test",
    });
  });

  it("surfaces the API error by keeping the dialog open", async () => {
    const onOpenChange = vi.fn();
    render(() => (
      <ManageEventDialog open event={event} onOpenChange={onOpenChange} onSubmit={rejectSubmit} />
    ));

    submit();

    await waitFor(() => expect(onOpenChange).not.toHaveBeenCalled());
  });
});

describe("PaginatedContainer", () => {
  const items = [{ id: "a" }, { id: "b" }, { id: "c" }];

  it("renders one page and reports the range", () => {
    render(() => (
      <PaginatedContainer
        items={items}
        pageSize={2}
        itemLabel="eventos"
        renderItems={(slice) => <ul>{slice.map((item) => <li>{item.id}</li>)}</ul>}
      />
    ));

    expect(screen.getAllByRole("listitem")).toHaveLength(2);
    expect(screen.getByText(/Mostrando 1–2 de 3 eventos/)).toBeInTheDocument();
  });

  it("moves to the next page and disables it at the end", async () => {
    render(() => (
      <PaginatedContainer
        items={items}
        pageSize={2}
        itemLabel="eventos"
        renderItems={(slice) => <ul>{slice.map((item) => <li>{item.id}</li>)}</ul>}
      />
    ));

    fireEvent.click(screen.getByRole("button", { name: "Próxima página" }));

    await waitFor(() => expect(screen.getAllByRole("listitem")).toHaveLength(1));
    expect(screen.getByText("c")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Próxima página" })).toBeDisabled();
  });

  it("shows the empty state instead of an empty grid", () => {
    render(() => (
      <PaginatedContainer
        items={[]}
        renderItems={() => <ul />}
        emptyState={<p>Nada por aqui</p>}
      />
    ));

    expect(screen.getByText("Nada por aqui")).toBeInTheDocument();
  });
});
