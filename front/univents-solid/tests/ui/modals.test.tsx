import { cleanup, fireEvent, render, screen, waitFor } from "@solidjs/testing-library";
import { afterEach, describe, expect, it, vi } from "vitest";
import { createSignal } from "solid-js";
import { AlertModal, Dialog, Modal } from "@trieoh/ui-solid";

afterEach(cleanup);

describe("Dialog and Modal", () => {
  it("renders when open and unmounts when closed", async () => {
    function TestWrapper() {
      const [open, setOpen] = createSignal(true);
      return (
        <Dialog
          open={open()}
          onOpenChange={setOpen}
          title="Título de Teste"
          description="Descrição do modal"
        >
          <div data-testid="dialog-body">Conteúdo do modal</div>
        </Dialog>
      );
    }

    render(() => <TestWrapper />);

    expect(screen.getByRole("heading", { name: "Título de Teste" })).toBeInTheDocument();
    expect(screen.getByText("Descrição do modal")).toBeInTheDocument();
    expect(screen.getByTestId("dialog-body")).toBeInTheDocument();

    const closeBtn = screen.getByRole("button", { name: "Fechar" });
    fireEvent.click(closeBtn);

    await waitFor(() => {
      expect(screen.queryByRole("heading", { name: "Título de Teste" })).toBeNull();
    });
  });

  it("supports Modal as an alias for Dialog with custom size", () => {
    render(() => (
      <Modal
        open={true}
        onOpenChange={() => {}}
        title="Modal Alias Test"
        size="lg"
      >
        <p>Conteúdo</p>
      </Modal>
    ));

    expect(screen.getByRole("heading", { name: "Modal Alias Test" })).toBeInTheDocument();
    const dialogElement = screen.getByRole("dialog");
    expect(dialogElement.className).toContain("max-w-2xl");
  });

  it("closes on Escape key press when preventClose is false", async () => {
    const onOpenChange = vi.fn();
    render(() => (
      <Dialog
        open={true}
        onOpenChange={onOpenChange}
        title="Esc Key Test"
      >
        <p>Press ESC</p>
      </Dialog>
    ));

    fireEvent.keyDown(window, { key: "Escape" });
    expect(onOpenChange).toHaveBeenCalledWith(false);
  });

  it("does not close on Escape key press when preventClose is true", async () => {
    const onOpenChange = vi.fn();
    render(() => (
      <Dialog
        open={true}
        onOpenChange={onOpenChange}
        title="Prevent Close Test"
        preventClose={true}
      >
        <p>Prevented</p>
      </Dialog>
    ));

    fireEvent.keyDown(window, { key: "Escape" });
    expect(onOpenChange).not.toHaveBeenCalled();
  });

  it("ensures scrollable content container has min-h-0 and overflow-y-auto", () => {
    render(() => (
      <Dialog
        open={true}
        onOpenChange={() => {}}
        title="Scroll Container Test"
        footer={<button type="button">Action</button>}
      >
        <div data-testid="scroll-content">Big long content</div>
      </Dialog>
    ));

    const contentWrapper = screen.getByTestId("scroll-content").parentElement;
    expect(contentWrapper).not.toBeNull();
    expect(contentWrapper?.className).toContain("min-h-0");
    expect(contentWrapper?.className).toContain("overflow-y-auto");
  });
});

describe("AlertModal", () => {
  it("renders title, description and triggers onConfirm and onCancel", () => {
    const onConfirm = vi.fn();
    const onCancel = vi.fn();
    const onOpenChange = vi.fn();

    render(() => (
      <AlertModal
        open={true}
        onOpenChange={onOpenChange}
        title="Remover assinatura?"
        description="Esta ação não poderá ser desfeita."
        confirmLabel="Sim, remover"
        cancelLabel="Voltar"
        variant="destructive"
        onConfirm={onConfirm}
        onCancel={onCancel}
      />
    ));

    expect(screen.getByRole("alertdialog")).toBeInTheDocument();
    expect(screen.getByText("Remover assinatura?")).toBeInTheDocument();
    expect(screen.getByText("Esta ação não poderá ser desfeita.")).toBeInTheDocument();

    const cancelBtn = screen.getByRole("button", { name: "Voltar" });
    fireEvent.click(cancelBtn);
    expect(onCancel).toHaveBeenCalledTimes(1);
    expect(onOpenChange).toHaveBeenCalledWith(false);

    const confirmBtn = screen.getByRole("button", { name: "Sim, remover" });
    fireEvent.click(confirmBtn);
    expect(onConfirm).toHaveBeenCalledTimes(1);
  });

  it("disables buttons and shows spinner when loading is true", () => {
    const onConfirm = vi.fn();
    render(() => (
      <AlertModal
        open={true}
        onOpenChange={() => {}}
        title="Excluir item"
        variant="destructive"
        loading={true}
        onConfirm={onConfirm}
      />
    ));

    const cancelBtn = screen.getByRole("button", { name: "Cancelar" });
    const confirmBtn = screen.getByRole("button", { name: /Processando/i });

    expect(cancelBtn).toBeDisabled();
    expect(confirmBtn).toBeDisabled();
  });
});
