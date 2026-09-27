import { cleanup, fireEvent, render, screen, waitFor } from "@solidjs/testing-library";
import {
  MultiStepCustomField,
  MultiStepDialog,
  MultiStepTextField,
  MultiStepTextareaField,
  type MultiStepItem,
} from "@trieoh/ui-solid";
import { afterEach, describe, expect, it, vi } from "vitest";

afterEach(cleanup);

const STEPS: MultiStepItem[] = [
  { id: "step1", title: "Primeiro Passo", description: "Dados iniciais" },
  { id: "step2", title: "Segundo Passo", description: "Confirmação" },
];

describe("MultiStepDialog", () => {
  it("renders the first step with mobile progressbar and navigation", () => {
    render(() => (
      <MultiStepDialog
        open
        onOpenChange={() => undefined}
        title="Assistente de Teste"
        steps={STEPS}
      >
        {(ctx) => (
          <div>
            <p>Conteúdo atual: {ctx.step.title}</p>
          </div>
        )}
      </MultiStepDialog>
    ));

    expect(screen.getByText("Assistente de Teste")).toBeInTheDocument();
    expect(screen.getByText("Conteúdo atual: Primeiro Passo")).toBeInTheDocument();

    const progress = screen.getByRole("progressbar");
    expect(progress).toHaveAttribute("aria-valuenow", "1");
    expect(progress).toHaveAttribute("aria-valuemax", "2");

    // First step has "Cancelar" and "Continuar"
    expect(screen.getByRole("button", { name: /cancelar/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /continuar/i })).toBeInTheDocument();
  });

  it("applies fixedHeight style to dialog panel by default to avoid jumping", () => {
    render(() => (
      <MultiStepDialog
        open
        onOpenChange={() => undefined}
        title="Assistente de Teste"
        steps={STEPS}
      >
        {() => <div>Conteúdo</div>}
      </MultiStepDialog>
    ));

    const dialogPanel = screen.getByRole("dialog");
    expect(dialogPanel).toHaveStyle({ height: "calc(100dvh - 2rem)" });
  });

  it("navigates forward and backward through steps", async () => {
    render(() => (
      <MultiStepDialog
        open
        onOpenChange={() => undefined}
        title="Assistente de Teste"
        steps={STEPS}
      >
        {(ctx) => (
          <div>
            <p>Passo ativo: {ctx.currentStep}</p>
          </div>
        )}
      </MultiStepDialog>
    ));

    expect(screen.getByText("Passo ativo: 0")).toBeInTheDocument();

    // Advance to step 1
    const nextBtn = screen.getByRole("button", { name: /continuar/i });
    fireEvent.click(nextBtn);

    await waitFor(() => {
      expect(screen.getByText("Passo ativo: 1")).toBeInTheDocument();
    });

    // Step 1 has "Voltar" and "Salvar"
    const backBtn = screen.getByRole("button", { name: /voltar/i });
    expect(backBtn).toBeInTheDocument();
    const saveBtn = screen.getByRole("button", { name: /salvar/i });
    expect(saveBtn).toBeInTheDocument();

    // Click back to step 0
    fireEvent.click(backBtn);

    await waitFor(() => {
      expect(screen.getByText("Passo ativo: 0")).toBeInTheDocument();
    });
  });

  it("blocks advance when onBeforeNext returns false", async () => {
    let allowed = false;
    const guard = vi.fn(() => allowed);

    render(() => (
      <MultiStepDialog
        open
        onOpenChange={() => undefined}
        title="Assistente de Teste"
        steps={STEPS}
        onBeforeNext={guard}
      >
        {(ctx) => (
          <div>
            <p>Passo ativo: {ctx.currentStep}</p>
          </div>
        )}
      </MultiStepDialog>
    ));

    const nextBtn = screen.getByRole("button", { name: /continuar/i });
    fireEvent.click(nextBtn);

    await waitFor(() => {
      expect(guard).toHaveBeenCalledTimes(1);
    });
    expect(screen.getByText("Passo ativo: 0")).toBeInTheDocument();

    // Now permit
    allowed = true;
    fireEvent.click(nextBtn);

    await waitFor(() => {
      expect(guard).toHaveBeenCalledTimes(2);
      expect(screen.getByText("Passo ativo: 1")).toBeInTheDocument();
    });
  });

  it("triggers onSubmit on the final step and closes upon completion", async () => {
    const onSubmit = vi.fn().mockResolvedValue(true);
    const onOpenChange = vi.fn();

    render(() => (
      <MultiStepDialog
        open
        onOpenChange={onOpenChange}
        title="Assistente de Teste"
        steps={STEPS}
        onSubmit={onSubmit}
      >
        {(ctx) => (
          <div>
            <p>Passo ativo: {ctx.currentStep}</p>
          </div>
        )}
      </MultiStepDialog>
    ));

    // Advance to step 1
    fireEvent.click(screen.getByRole("button", { name: /continuar/i }));

    await waitFor(() => {
      expect(screen.getByText("Passo ativo: 1")).toBeInTheDocument();
    });

    // Submit
    fireEvent.click(screen.getByRole("button", { name: /salvar/i }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledTimes(1);
      expect(onOpenChange).toHaveBeenCalledWith(false);
    });
  });

  it("renders declarative fields and summary cards without manual children boilerplate", async () => {
    interface TestValues {
      title: string;
      notes: string;
    }

    const onChange = vi.fn();
    const onSubmit = vi.fn().mockResolvedValue(true);
    const onOpenChange = vi.fn();

    const declarativeSteps: MultiStepItem<TestValues>[] = [
      {
        id: "info",
        title: "Informações",
        fields: [
          {
            name: "title",
            label: "Título Principal",
            placeholder: "Digite o título...",
            required: true,
          },
          {
            kind: "textarea",
            name: "notes",
            label: "Anotações",
            placeholder: "Notas extras...",
            rows: 3,
          },
        ],
      },
      {
        id: "confirm",
        title: "Confirmação",
        summary: {
          title: "Revisão Final",
          badge: "Pronto",
          items: [
            { label: "Título", value: "Meu Teste" },
            { label: "Anotações", value: "Nota de exemplo", fullWidth: true },
          ],
        },
      },
    ];

    render(() => (
      <MultiStepDialog<TestValues>
        open
        onOpenChange={onOpenChange}
        title="Cadastro Declarativo"
        steps={declarativeSteps}
        values={{ title: "Meu Teste", notes: "Nota de exemplo" }}
        onChange={onChange}
        onSubmit={onSubmit}
      />
    ));

    // Renders automatically generated labels and inputs
    expect(screen.getByText("Título Principal")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("Digite o título...")).toHaveValue("Meu Teste");
    expect(screen.getByPlaceholderText("Notas extras...")).toHaveValue("Nota de exemplo");

    // Advance to summary step
    fireEvent.click(screen.getByRole("button", { name: /continuar/i }));

    await waitFor(() => {
      expect(screen.getByText("Revisão Final")).toBeInTheDocument();
      expect(screen.getByText("Pronto")).toBeInTheDocument();
      expect(screen.getByText("Nota de exemplo")).toBeInTheDocument();
    });

    // Submit from summary step
    fireEvent.click(screen.getByRole("button", { name: /salvar/i }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledTimes(1);
      expect(onOpenChange).toHaveBeenCalledWith(false);
    });
  });

  it("renders modular input-specific subcomponents directly", () => {
    const onTextChange = vi.fn();
    const onTextareaChange = vi.fn();

    render(() => (
      <div>
        <MultiStepTextField
          field={{ name: "name", label: "Nome Teste", placeholder: "Insira o nome" }}
          value="Valor Inicial"
          onChange={onTextChange}
        />
        <MultiStepTextareaField
          field={{ kind: "textarea", name: "bio", label: "Bio Teste", placeholder: "Insira a bio" }}
          value="Bio inicial"
          onChange={onTextareaChange}
        />
        <MultiStepCustomField
          field={{
            kind: "custom",
            name: "custom",
            label: "Custom Label",
            render: (args) => <span id={args.ids.id}>Render Customizado</span>,
          }}
        />
      </div>
    ));

    expect(screen.getByLabelText("Nome Teste")).toHaveValue("Valor Inicial");
    expect(screen.getByLabelText("Bio Teste")).toHaveValue("Bio inicial");
    expect(screen.getByText("Render Customizado")).toBeInTheDocument();
  });
});
