import { createFileRoute } from "@tanstack/react-router";
import { ShoppingBag } from "lucide-react";
import { useState } from "react";
import type { ProductI } from "@/features/products/model";
import { ProductList } from "@/features/products/ui/ProductList";
import { Cart } from "@/features/products/ui/Cart";
import { useCart } from "@/features/products/hooks/use-cart";
import { Button } from "@/shared/ui/shadcn/button";
import { cn } from "@/shared/lib/utils";

export const Route = createFileRoute("/events/$eventId/editions/$editionId/products")({
  component: ProductsPage,
});

function ProductsPage() {
  const { eventId, editionId } = Route.useParams();
  const { totalCents } = useCart(editionId);
  const [isCartOpen, setIsCartOpen] = useState(false);

  const totalFormatted = new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(totalCents / 100);

  // Mock data as requested
  const mockProducts: ProductI[] = [
    {
      id: "1",
      scope_id: "scope-1",
      edition_id: editionId,
      name: "Ingresso VIP - Lote Antecipado",
      description: "Acesso exclusivo à área VIP, lounge com buffet e vista privilegiada do palco principal.",
      type: "ticket",
      ticket_id: "t1",
      price_cents: 45000,
      status: "available",
      available_from: new Date().toISOString(),
      available_until: new Date(Date.now() + 86400000 * 7).toISOString(),
      has_inventory: true,
      inventory_quantity: 100,
      inventory_remaining: 4,
      created_by: "admin",
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      deleted_at: null,
    },
    {
      id: "2",
      scope_id: "scope-1",
      edition_id: editionId,
      name: "Combo Camiseta + Boné Oficial",
      description: "Kit exclusivo da edição comemorativa. Tecido 100% algodão e boné snapback bordado.",
      type: "bundle",
      ticket_id: null,
      price_cents: 12990,
      status: "available",
      available_from: new Date().toISOString(),
      available_until: null,
      has_inventory: true,
      inventory_quantity: 200,
      inventory_remaining: 156,
      created_by: "admin",
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      deleted_at: null,
    },
    {
      id: "3",
      scope_id: "scope-1",
      edition_id: editionId,
      name: "Token de Consumo - R$ 100",
      description: "Crédito antecipado para consumo de bebidas e alimentos durante o evento. Evite filas!",
      type: "token",
      ticket_id: null,
      price_cents: 10000,
      status: "available",
      available_from: new Date().toISOString(),
      available_until: null,
      has_inventory: false,
      inventory_quantity: 0,
      inventory_remaining: 0,
      created_by: "admin",
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      deleted_at: null,
    },
    {
      id: "4",
      scope_id: "scope-1",
      edition_id: editionId,
      name: "Ingresso Pista - Inteira",
      description: "Acesso geral ao evento e todas as áreas comuns de entretenimento.",
      type: "ticket",
      ticket_id: "t2",
      price_cents: 15000,
      status: "sold_out",
      available_from: new Date().toISOString(),
      available_until: null,
      has_inventory: true,
      inventory_quantity: 1000,
      inventory_remaining: 0,
      created_by: "admin",
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      deleted_at: null,
    },
  ];

  return (
    <div className="min-h-screen pb-20">
      <div className="max-w-7xl mx-auto px-4 pt-8 pb-4 flex items-center justify-end">
        <Button
          onClick={() => { setIsCartOpen(true); }}
          size="lg"
          className={cn(
            "rounded-full bg-accent hover:bg-accent/90 text-white shadow-xl",
            "shadow-accent/20 px-6 font-bold gap-3"
          )}
        >
          <ShoppingBag className="h-5 w-5" />
          <span>Carrinho</span>
          {totalCents > 0 && (
            <span className="border-l border-white/30 pl-3 font-mono text-sm">
              {totalFormatted}
            </span>
          )}
        </Button>
      </div>

      <ProductList products={mockProducts} />

      {/* Cart Drawer */}
      <Cart
        isOpen={isCartOpen}
        eventId={eventId}
        editionId={editionId}
        onClose={() => { setIsCartOpen(false); }}
      />
    </div>
  );
}