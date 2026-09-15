import { cleanup, fireEvent, render, screen, waitFor } from "@solidjs/testing-library";
import { afterEach, describe, expect, it, vi } from "vitest";

import type { EditionPurchase } from "@trieoh/univents-api/schemas";
import type { TicketI } from "@/features/tickets/model";
import { AdminCreateTicketCard } from "@/features/tickets/ui/AdminCreateTicketCard";
import { AdminTicketCard } from "@/features/tickets/ui/AdminTicketCard";
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
