import { describe, expect, it, vi, afterEach } from "vitest";
import { render, screen, fireEvent, cleanup } from "@solidjs/testing-library";
import { StepChecklist } from "@/widgets/ui/StepChecklist";

describe("StepChecklist", () => {
  afterEach(() => {
    cleanup();
  });

  const items = [
    {
      id: "step1",
      title: "Passo 1",
      description: "Descrição 1",
      completed: true,
    },
    {
      id: "step2",
      title: "Passo 2",
      description: "Descrição 2",
      completed: false,
      action: {
        label: "Fazer passo 2",
        onClick: vi.fn(),
      },
    },
    {
      id: "step3",
      title: "Passo 3",
      description: "Descrição 3",
      completed: false,
    },
  ];

  it("renders trigger button at the bottom and shows completion count", () => {
    render(() => <StepChecklist title="Checklist" items={items} />);

    const button = screen.getByRole("button", {
      name: /Abrir checklist de configuração/i,
    });
    expect(button).toBeDefined();
    expect(button.textContent).toContain("Configuração (1/3)");
  });

  it("renders discreet trigger button when all items are completed", () => {
    const allCompletedItems = items.map((item) => ({ ...item, completed: true }));
    render(() => <StepChecklist title="Checklist" items={allCompletedItems} />);

    const button = screen.getByRole("button", {
      name: /Abrir checklist de configuração/i,
    });
    expect(button).toBeDefined();
    expect(button.textContent).toContain("Configurado (3/3)");
    expect(button.className).toContain("opacity-60");
  });

  it("opens panel with mobile backdrop when clicked and renders all steps", async () => {
    render(() => <StepChecklist title="Checklist de Teste" items={items} />);

    const trigger = screen.getByRole("button", {
      name: /Abrir checklist de configuração/i,
    });
    await fireEvent.click(trigger);

    expect(screen.getByText("Checklist de Teste")).toBeDefined();
    expect(screen.getByText("1 de 3 concluídos (33%)")).toBeDefined();
    expect(screen.getByText("Passo 1")).toBeDefined();
    expect(screen.getByText("Passo 2")).toBeDefined();
    expect(screen.getByText("Passo 3")).toBeDefined();

    const closeButton = screen.getByRole("button", {
      name: /Fechar checklist/i,
    });
    await fireEvent.click(closeButton);

    expect(
      screen.getByRole("button", {
        name: /Abrir checklist de configuração/i,
      }),
    ).toBeDefined();
  });

  it("allows clicking mobile backdrop to dismiss drawer", async () => {
    render(() => <StepChecklist title="Checklist de Teste" items={items} />);

    const trigger = screen.getByRole("button", {
      name: /Abrir checklist de configuração/i,
    });
    await fireEvent.click(trigger);

    expect(screen.getByText("Checklist de Teste")).toBeDefined();

    const backdrop = screen.getByTestId("backdrop");
    await fireEvent.click(backdrop);

    expect(
      screen.getByRole("button", {
        name: /Abrir checklist de configuração/i,
      }),
    ).toBeDefined();
  });

  it("allows dragging down on mobile handle to dismiss drawer", async () => {
    render(() => <StepChecklist title="Checklist de Teste" items={items} />);

    const trigger = screen.getByRole("button", {
      name: /Abrir checklist de configuração/i,
    });
    await fireEvent.click(trigger);

    expect(screen.getByText("Checklist de Teste")).toBeDefined();

    const handle = screen.getByTestId("drag-handle");
    expect(handle).toBeDefined();

    await fireEvent.mouseDown(handle, { clientY: 100 });
    await fireEvent.mouseMove(window, { clientY: 250 });
    await fireEvent.mouseUp(window);

    expect(
      screen.getByRole("button", {
        name: /Abrir checklist de configuração/i,
      }),
    ).toBeDefined();
  });
});
