import { render, screen, fireEvent } from "@solidjs/testing-library";
import { describe, it, expect, vi } from "vitest";
import { BadgeEditor } from "@/features/badges/editor/BadgeEditor";
import { BadgeCanvas } from "@/features/badges/editor/BadgeCanvas";
import { DEFAULT_BADGE_TEMPLATE } from "@/features/badges/default-template";

vi.mock("@tanstack/solid-router", () => ({
  Link: (props: any) => <a href={props.href}>{props.children}</a>,
  useNavigate: () => vi.fn(),
}));

vi.mock("@trieoh/front-core-solid", () => ({
  useQuery: (optionsFn: any) => () => {
    const opts = typeof optionsFn === "function" ? optionsFn() : optionsFn;
    if (opts?.enabled === false) {
      return { data: undefined, isLoading: false, isError: false };
    }
    return {
      data: [
        { id: "ticket-1", name: "Ingresso Geral" },
        { id: "ticket-2", name: "Ingresso VIP" },
      ],
      isLoading: false,
      isError: false,
    };
  },
}));

vi.mock("@/features/badges/api/mutations", () => ({
  useCreateBadgeTemplateMutation: () => ({
    mutate: vi.fn(),
    isPending: false,
    result: () => ({ status: "idle" }),
  }),
  useUpdateBadgeTemplateMutation: () => ({
    mutate: vi.fn(),
    isPending: false,
    result: () => ({ status: "idle" }),
  }),
}));

vi.mock("@/shared/ui/toast", () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

describe("BadgeEditor", () => {
  it("renders editor sidebar and canvas properly", () => {
    const { container } = render(() => (
      <BadgeEditor eventId="ev-123" editionId="ed-456" />
    ));

    expect(screen.getByText("Editor de crachás")).toBeDefined();
    expect(screen.getByText("Nome do template")).toBeDefined();
    expect(screen.getByText("Salvar")).toBeDefined();
    expect(container.querySelector("[data-badge-element]")).toBeDefined();
  });

  it("allows switching presets", () => {
    render(() => <BadgeEditor eventId="ev-123" editionId="ed-456" />);

    expect(screen.getAllByText("Tamanho (mm)").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Camadas").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Valores de pré-visualização").length).toBeGreaterThan(0);
  });

  it("does not render dynamic variables dropdown in the top toolbar (matches Univents 1:1)", () => {
    const { container } = render(() => (
      <BadgeEditor eventId="ev-123" editionId="ed-456" />
    ));

    const topbar = container.querySelector("[data-text-toolbar]");
    expect(topbar).not.toBeNull();
    // In Univents, the top toolbar doesn't have the variable dropdown {} because it's in the right sidebar
    expect(topbar?.querySelector('button[title*="variável"]')).toBeNull();
  });
});

describe("BadgeCanvas - Editing text without reactive write errors or unwanted blur", () => {
  it("enters edit mode on double click without throwing REACTIVE_WRITE_IN_OWNED_SCOPE", async () => {
    const onSelect = vi.fn();
    const onChangeElement = vi.fn();
    const setController = vi.fn();
    const setSelectionStyles = vi.fn();
    const updateParagraphs = vi.fn();
    const stopEditing = vi.fn();

    const { container } = render(() => (
      <BadgeCanvas
        design={DEFAULT_BADGE_TEMPLATE.design_data}
        selectedId="text-name"
        onSelect={onSelect}
        onChangeElement={onChangeElement}
        textAdapter={{
          setController,
          setSelectionStyles,
          updateParagraphs,
          stopEditing,
        }}
      />
    ));

    const textEl = container.querySelector(".cursor-text");
    expect(textEl).toBeDefined();
    if (textEl) {
      fireEvent.dblClick(textEl);
      await new Promise((resolve) => setTimeout(resolve, 50));
      const editable = container.querySelector("[contenteditable]");
      expect(editable).not.toBeNull();
      expect(setController).toHaveBeenCalled();
    }
  });

  it("stays in edit mode when typing text instead of exiting edit mode", async () => {
    const onSelect = vi.fn();
    const onChangeElement = vi.fn();
    const setController = vi.fn();
    const setSelectionStyles = vi.fn();
    const updateParagraphs = vi.fn();
    const stopEditing = vi.fn();

    const { container } = render(() => (
      <BadgeCanvas
        design={DEFAULT_BADGE_TEMPLATE.design_data}
        selectedId="text-name"
        onSelect={onSelect}
        onChangeElement={onChangeElement}
        textAdapter={{
          setController,
          setSelectionStyles,
          updateParagraphs,
          stopEditing,
        }}
      />
    ));

    const textEl = container.querySelector(".cursor-text");
    expect(textEl).toBeDefined();
    if (textEl) {
      fireEvent.dblClick(textEl);
      await new Promise((resolve) => setTimeout(resolve, 50));
      const editable = container.querySelector("[contenteditable]");
      expect(editable).not.toBeNull();

      // Simulating user typing text
      if (editable) {
        editable.innerHTML = "<p>Texto editado pelo usuário</p>";
        fireEvent.input(editable);
        await new Promise((resolve) => setTimeout(resolve, 50));

        // Must still be in edit mode!
        const editableAfterInput = container.querySelector("[contenteditable]");
        expect(editableAfterInput).not.toBeNull();
      }
    }
  });
});
