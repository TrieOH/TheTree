import { cleanup, fireEvent, render, screen, waitFor } from "@solidjs/testing-library";
import { afterEach, describe, expect, it, vi } from "vitest";

import type { EditionPurchase } from "@trieoh/univents-api/schemas";
import type { TicketI } from "@/features/tickets/model";
import { AdminCreateTicketCard } from "@/features/tickets/ui/AdminCreateTicketCard";
import { AdminTicketCard } from "@/features/tickets/ui/AdminTicketCard";
import { ManageTicketDialog } from "@/features/tickets/ui/ManageTicketDialog";
import { AdminPurchaseCard } from "@/features/purchases/ui/AdminPurchaseCard";

afterEach(cleanup);

const mockTicket: TicketI = {
  id: "ticket-1",
  edition_id: "ed-1",
  name: "Ingresso VIP",
  description: "Acesso total à área VIP",
  price_cents: 15000,
  access_level: 2,
  max_quantity: 50,
  created_at: "2026-01-01T12:00:00Z",
  updated_at: null,
};

const mockPurchase: EditionPurchase = {
  purchase_id: "018f3a9a-bcde-7123-89ab-cdef01234567",
  edition_id: "ed-1",
  status: "approved",
  status_reason: "Pagamento aprovado via Pix",
  total_cents: 15000,
  currency: "BRL",
  payment_method: "pix",
  payer_email: "comprador@trieoh.com",
  created_at: "2026-03-01T14:30:00Z",
  items: [
    {
      item_type: "ticket",
      item_id: "ticket-1",
      quantity: 1,
      unit_price_cents: 15000,
    },
  ],
  attendees: [
    {
      name: "João Silva",
      email: "joao@trieoh.com",
    },
  ],
};

describe("AdminCreateTicketCard", () => {
  it("renders creation card with dashed style and responds to click", () => {
    const onCreate = vi.fn();
    render(() => <AdminCreateTicketCard onCreate={onCreate} animate={false} />);

    const btn = screen.getByRole("button", { name: /Novo ticket/i });
    expect(btn).toBeInTheDocument();
    expect(
      screen.getByText("Configure ingressos pagos ou gratuitos, lotes, limite de vagas e níveis de acesso."),
    ).toBeInTheDocument();

    fireEvent.click(btn);
    expect(onCreate).toHaveBeenCalledTimes(1);
  });
});

describe("AdminTicketCard", () => {
  it("renders ticket information and allows editing", () => {
    const onEdit = vi.fn();
    render(() => (
      <AdminTicketCard
        ticket={mockTicket}
        onEdit={onEdit}
        animate={false}
      />
    ));

    expect(screen.getByText("Ingresso VIP")).toBeInTheDocument();
    expect(screen.getByText("Acesso total à área VIP")).toBeInTheDocument();
    expect(screen.getByText("Acesso nível 2")).toBeInTheDocument();
    expect(screen.getByText("50 vagas")).toBeInTheDocument();
    expect(screen.getByText("R$ 150,00")).toBeInTheDocument();

    const editBtn = screen.getByRole("button", { name: /Editar Ingresso VIP/i });
    fireEvent.click(editBtn);
    expect(onEdit).toHaveBeenCalledWith(mockTicket);
  });

  it("handles free tickets and unlimited capacity correctly", () => {
    const freeTicket: TicketI = {
      ...mockTicket,
      id: "free-1",
      name: "Entrada Gratuita",
      price_cents: 0,
      max_quantity: null,
    };

    render(() => (
      <AdminTicketCard
        ticket={freeTicket}
        onEdit={vi.fn()}
        animate={false}
      />
    ));

    expect(screen.getByText("Entrada Gratuita")).toBeInTheDocument();
    expect(screen.getByText("Gratuito")).toBeInTheDocument();
    expect(screen.getByText("Ilimitado")).toBeInTheDocument();
  });
});

describe("AdminPurchaseCard", () => {
  it("renders purchase data, expands details, and triggers refund", async () => {
    const onRefund = vi.fn();
    render(() => (
      <AdminPurchaseCard
        purchase={mockPurchase}
        onRefund={onRefund}
        animate={false}
      />
    ));

    expect(screen.getByText("Aprovado")).toBeInTheDocument();
    expect(screen.getByText("Pagamento aprovado via Pix")).toBeInTheDocument();
    expect(screen.getByText("comprador@trieoh.com")).toBeInTheDocument();
    expect(screen.getByText("R$ 150,00")).toBeInTheDocument();

    // Toggle attendees details
    const attendeesBtn = screen.getByRole("button", { name: /1 participante/i });
    fireEvent.click(attendeesBtn);

    await waitFor(() => {
      expect(screen.getByText("João Silva")).toBeInTheDocument();
      expect(screen.getByText("joao@trieoh.com")).toBeInTheDocument();
    });

    // Trigger refund
    const refundBtn = screen.getByRole("button", { name: /Reembolsar/i });
    fireEvent.click(refundBtn);
    expect(onRefund).toHaveBeenCalledWith(mockPurchase);
  });
});

describe("ManageTicketDialog", () => {
  it("steps through creation flow, validates access level required, and submits valid ticket", async () => {
    const onSubmit = vi.fn().mockResolvedValue(true);
    const onOpenChange = vi.fn();

    render(() => (
      <ManageTicketDialog
        open={true}
        onOpenChange={onOpenChange}
        ticket={null}
        onSubmit={onSubmit}
      />
    ));

    expect(screen.getByText("Novo ticket")).toBeInTheDocument();

    // Try advancing with empty name
    const continueBtn = screen.getByRole("button", { name: /Continuar/i });
    fireEvent.click(continueBtn);

    await waitFor(() => {
      expect(screen.getByText("O nome deve ter pelo menos 2 caracteres.")).toBeInTheDocument();
    });

    // Fill name and description
    const nameInput = screen.getByLabelText(/Nome do ticket/i);
    fireEvent.input(nameInput, { target: { value: "Ingresso Regular" } });

    // Step error should be cleared when user types
    await waitFor(() => {
      expect(screen.queryByText("O nome deve ter pelo menos 2 caracteres.")).not.toBeInTheDocument();
    });

    // Advance to Step 2: Valores e Vagas
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByLabelText(/^Preço/i)).toBeInTheDocument();
    });

    // Test access level required validation: try advancing with unset access level
    const step2ContinueBtn = screen.getByRole("button", { name: /Continuar/i });
    fireEvent.click(step2ContinueBtn);

    await waitFor(() => {
      expect(screen.getByText("O nível de acesso é obrigatório.")).toBeInTheDocument();
    });

    // Fill price using monetary mask (typing 4990 gives R$ 49,90)
    const priceInput = screen.getByLabelText(/^Preço/i);
    fireEvent.input(priceInput, { target: { value: "4990" } });

    // Fill access level
    const accessLevelInput = screen.getByLabelText(/Nível de acesso/i);
    fireEvent.input(accessLevelInput, { target: { value: "0" } });

    // Fill max quantity
    const maxQtyInput = screen.getByLabelText(/Quantidade máxima \(vagas\)/i);
    fireEvent.input(maxQtyInput, { target: { value: "100" } });

    // Advance to Step 3: Resumo
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByText("Dados do ticket")).toBeInTheDocument();
      expect(screen.getByText("Pronto para criar")).toBeInTheDocument();
      expect(screen.getAllByText("100 vagas").length).toBeGreaterThan(0);
    });

    // Submit
    const submitBtn = screen.getByRole("button", { name: /Criar ticket/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        name: "Ingresso Regular",
        description: null,
        price_cents: 4990,
        access_level: 0,
        max_quantity: 100,
      });
      expect(onOpenChange).toHaveBeenCalledWith(false);
    });
  });

  it("populates existing ticket in edit mode and updates", async () => {
    const onSubmit = vi.fn().mockResolvedValue(true);
    const onOpenChange = vi.fn();

    render(() => (
      <ManageTicketDialog
        open={true}
        onOpenChange={onOpenChange}
        ticket={mockTicket}
        onSubmit={onSubmit}
      />
    ));

    expect(screen.getByText("Editar ticket")).toBeInTheDocument();
    const nameInput = screen.getByLabelText(/Nome do ticket/i) as HTMLInputElement;
    expect(nameInput.value).toBe("Ingresso VIP");

    // Advance to Step 2
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      const priceInput = screen.getByLabelText(/^Preço/i) as HTMLInputElement;
      expect(priceInput.value).toMatch(/150/);
    });

    // Advance to Step 3 (Resumo)
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByText("Pronto para atualizar")).toBeInTheDocument();
    });

    const submitBtn = screen.getByRole("button", { name: /Salvar alterações/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        name: "Ingresso VIP",
        description: "Acesso total à área VIP",
        price_cents: 15000,
        access_level: 2,
        max_quantity: 50,
      });
      expect(onOpenChange).toHaveBeenCalledWith(false);
    });
  });
});
