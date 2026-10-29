import { cleanup, fireEvent, render, screen } from "@solidjs/testing-library";
import { createSignal, flush } from "solid-js";
import { afterEach, describe, expect, it, vi } from "vitest";

import { Tabs } from "@trieoh/ui-solid";

afterEach(cleanup);

const items = [
  { id: "overview", label: "Visão geral", content: <p>painel visão geral</p> },
  { id: "activity", label: "Atividade", count: 8, content: <p>painel atividade</p> },
  { id: "billing", label: "Cobrança", disabled: true, content: <p>painel cobrança</p> },
  { id: "settings", label: "Ajustes", content: <p>painel ajustes</p> },
];

const tab = (name: string) => screen.getByRole("tab", { name: new RegExp(name) });
// Solid 2 batches writes, so each interaction is flushed before it is asserted.
const click = (element: HTMLElement) => {
  fireEvent.click(element);
  flush();
};
const key = (name: string, k: string) => {
  fireEvent.keyDown(tab(name), { key: k });
  flush();
};

describe("Tabs", () => {
  it("wires tabs to their panel with ARIA", () => {
    render(() => <Tabs items={items} ariaLabel="Seções" />);

    expect(screen.getByRole("tablist", { name: "Seções" })).toHaveAttribute(
      "aria-orientation",
      "horizontal",
    );
    expect(tab("Visão geral")).toHaveAttribute("aria-selected", "true");
    expect(tab("Atividade")).toHaveAttribute("aria-selected", "false");

    const panel = screen.getByRole("tabpanel");
    expect(panel).toHaveTextContent("painel visão geral");
    expect(panel.getAttribute("aria-labelledby")).toBe(tab("Visão geral").id);
    expect(tab("Visão geral").getAttribute("aria-controls")).toBe(panel.id);
    // Unselected panels are not mounted, so nothing points at them.
    expect(tab("Atividade")).not.toHaveAttribute("aria-controls");
  });

  it("selects on click and shows the matching panel", () => {
    const onChange = vi.fn();
    render(() => <Tabs items={items} onChange={onChange} />);

    click(tab("Atividade"));

    expect(onChange).toHaveBeenCalledWith("activity");
    expect(tab("Atividade")).toHaveAttribute("aria-selected", "true");
    expect(screen.getByRole("tabpanel")).toHaveTextContent("painel atividade");
    expect(screen.queryByText("painel visão geral")).toBeNull();
  });

  it("starts on defaultValue", () => {
    render(() => <Tabs items={items} defaultValue="settings" />);

    expect(tab("Ajustes")).toHaveAttribute("aria-selected", "true");
  });

  it("is controlled when given a value", () => {
    const onChange = vi.fn();
    const [value, setValue] = createSignal("overview");
    render(() => (
      <Tabs
        items={items}
        value={value()}
        onChange={(id) => {
          onChange(id);
          setValue(id);
        }}
      />
    ));

    click(tab("Ajustes"));
    expect(onChange).toHaveBeenCalledWith("settings");
    expect(tab("Ajustes")).toHaveAttribute("aria-selected", "true");

    // A parent that ignores the change keeps the tab where it is.
    cleanup();
    render(() => <Tabs items={items} value="overview" onChange={() => {}} />);
    click(tab("Atividade"));
    expect(tab("Visão geral")).toHaveAttribute("aria-selected", "true");
  });

  it("gives only the selected tab a tab stop (roving tabindex)", () => {
    render(() => <Tabs items={items} />);

    expect(tab("Visão geral")).toHaveAttribute("tabindex", "0");
    expect(tab("Atividade")).toHaveAttribute("tabindex", "-1");
    expect(tab("Ajustes")).toHaveAttribute("tabindex", "-1");
  });

  describe("keyboard", () => {
    it("moves with the arrows, wrapping and skipping disabled tabs", () => {
      render(() => <Tabs items={items} />);

      key("Visão geral", "ArrowRight");
      expect(tab("Atividade")).toHaveAttribute("aria-selected", "true");
      expect(tab("Atividade")).toHaveFocus();

      // "Cobrança" is disabled, so the next stop is "Ajustes".
      key("Atividade", "ArrowRight");
      expect(tab("Ajustes")).toHaveAttribute("aria-selected", "true");

      key("Ajustes", "ArrowRight");
      expect(tab("Visão geral")).toHaveAttribute("aria-selected", "true");

      key("Visão geral", "ArrowLeft");
      expect(tab("Ajustes")).toHaveAttribute("aria-selected", "true");
    });

    it("jumps to the first and last enabled tab", () => {
      render(() => <Tabs items={items} defaultValue="activity" />);

      key("Atividade", "End");
      expect(tab("Ajustes")).toHaveAttribute("aria-selected", "true");

      key("Ajustes", "Home");
      expect(tab("Visão geral")).toHaveAttribute("aria-selected", "true");
    });

    it("follows the vertical axis when vertical", () => {
      render(() => <Tabs items={items} orientation="vertical" />);

      expect(screen.getByRole("tablist")).toHaveAttribute(
        "aria-orientation",
        "vertical",
      );
      key("Visão geral", "ArrowRight");
      expect(tab("Visão geral")).toHaveAttribute("aria-selected", "true");

      key("Visão geral", "ArrowDown");
      expect(tab("Atividade")).toHaveAttribute("aria-selected", "true");
    });

    it("only moves focus with manual activation", () => {
      render(() => <Tabs items={items} activation="manual" />);

      key("Visão geral", "ArrowRight");
      expect(tab("Atividade")).toHaveFocus();
      expect(tab("Visão geral")).toHaveAttribute("aria-selected", "true");

      click(tab("Atividade"));
      expect(tab("Atividade")).toHaveAttribute("aria-selected", "true");
    });

    it("ignores other keys", () => {
      render(() => <Tabs items={items} />);

      key("Visão geral", "a");
      expect(tab("Visão geral")).toHaveAttribute("aria-selected", "true");
    });
  });

  it("never selects a disabled tab", () => {
    const onChange = vi.fn();
    render(() => <Tabs items={items} onChange={onChange} />);

    click(tab("Cobrança"));

    expect(onChange).not.toHaveBeenCalled();
    expect(tab("Cobrança")).toBeDisabled();
  });

  it("keeps every panel mounted when asked, hiding the inactive ones", () => {
    render(() => <Tabs items={items} keepMounted />);

    const panels = screen.getAllByRole("tabpanel", { hidden: true });
    expect(panels).toHaveLength(items.length);
    expect(panels.filter((panel) => !panel.hidden)).toHaveLength(1);
    expect(tab("Atividade").getAttribute("aria-controls")).toBe(panels[1].id);
  });

  it("shows counts and badges on the triggers", () => {
    render(() => <Tabs items={items} />);

    expect(tab("Atividade")).toHaveTextContent("8");
  });

  it.each([
    ["underline", "border-b-2"],
    ["pills", "rounded-full"],
    ["segmented", "rounded-md"],
  ] as const)("renders the %s variant", (variant, expected) => {
    render(() => <Tabs items={items} variant={variant} />);

    expect(tab("Visão geral").className).toContain(expected);
  });

  describe("without panels (a bare switch for content you render)", () => {
    const bare = items.map(({ content: _content, ...item }) => item);

    it("mounts no panel and points no tab at one", () => {
      render(() => <Tabs items={bare} ariaLabel="Seções" />);

      expect(screen.queryByRole("tabpanel", { hidden: true })).toBeNull();
      expect(tab("Visão geral")).not.toHaveAttribute("aria-controls");
      expect(tab("Visão geral")).toHaveAttribute("aria-selected", "true");
    });

    it("adds no gap for the missing panel area and puts class on the wrapper", () => {
      const { container } = render(() => <Tabs items={bare} class="mb-4" />);
      const root = container.firstElementChild!;

      expect(root.className).toContain("mb-4");
      expect(root.className).not.toContain("gap-4");
      expect(root.children).toHaveLength(1);
    });

    it("still drives the page through value and onChange, with the keyboard", () => {
      const onChange = vi.fn();
      const [value, setValue] = createSignal("overview");
      render(() => (
        <Tabs
          items={bare}
          value={value()}
          onChange={(id) => {
            onChange(id);
            setValue(id);
          }}
        />
      ));

      click(tab("Atividade"));
      expect(onChange).toHaveBeenLastCalledWith("activity");

      key("Atividade", "ArrowRight");
      expect(onChange).toHaveBeenLastCalledWith("settings");
      expect(tab("Ajustes")).toHaveAttribute("aria-selected", "true");
    });
  });

  it("styles the segmented list as a track and stretches with fullWidth", () => {
    render(() => <Tabs items={items} variant="segmented" fullWidth />);

    expect(screen.getByRole("tablist").className).toContain("bg-muted");
    expect(tab("Visão geral").className).toContain("flex-1");
  });
});
