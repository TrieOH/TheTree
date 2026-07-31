import { Check, ShoppingCart, Star } from "lucide-react";
import { useCart } from "@/features/products/hooks/use-cart";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import type { TicketI } from "../model";

function formatPrice(cents: number): string {
  if (cents === 0) return "Gratuito";
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
    minimumFractionDigits: 0,
  }).format(cents / 100);
}

interface TicketCardProps {
  ticket: TicketI;
  isFeatured?: boolean;
  editionId?: string;
}

export function TicketCard({ ticket, isFeatured, editionId }: TicketCardProps) {
  const { addItem, items } = useCart(editionId ?? "");
  const inCart = items.find(
    (item) => item.id === ticket.id && item.type === "ticket",
  );
  const isFree = ticket.price_cents === 0;

  return (
    <div
      className={cn(
        "relative flex flex-col w-72 h-56 rounded-2xl border overflow-hidden transition-all duration-300 p-5",
        isFeatured
          ? "bg-card border-border/60 lg:scale-105 lg:bg-primary/4 lg:border-primary/25 lg:shadow-lg lg:shadow-primary/5 lg:z-10 hover:border-border hover:shadow-md"
          : "bg-card border-border/60 hover:border-border hover:shadow-md",
      )}
    >
      {/* Badge */}
      <div className="absolute top-0 right-0">
        {isFree ? (
          <span className="inline-flex items-center gap-1 px-3 py-1 rounded-bl-xl rounded-tr-2xl bg-emerald-500 text-white text-[10px] font-bold tracking-wide uppercase">
            <Check className="w-3 h-3" />
            Gratuito
          </span>
        ) : (
          <span className="inline-flex items-center gap-1 px-3 py-1 rounded-bl-xl rounded-tr-2xl bg-primary text-primary-foreground text-[10px] font-bold tracking-wide uppercase">
            <Star className="w-3 h-3 fill-current" />
            Pago
          </span>
        )}
      </div>

      {/* Name */}
      <h3 className="font-semibold text-foreground leading-tight pr-16 text-base">
        {ticket.name}
      </h3>

      {/* Price */}
      <div className="mt-2">
        <span
          className={cn(
            "text-2xl font-bold",
            isFree
              ? "text-emerald-600 dark:text-emerald-400"
              : isFeatured
                ? "text-foreground lg:text-primary"
                : "text-foreground",
          )}
        >
          {formatPrice(ticket.price_cents)}
        </span>
      </div>

      {/* Description */}
      {ticket.description && (
        <p className="mt-3 text-sm text-muted-foreground leading-relaxed whitespace-pre-line line-clamp-4">
          {ticket.description}
        </p>
      )}
      {editionId && (
        <Button
          size="sm"
          variant={inCart ? "secondary" : "default"}
          className="mt-5 h-9 w-full gap-2 text-xs font-semibold shadow-sm"
          onClick={() =>
            addItem(
              {
                id: ticket.id,
                type: "ticket",
                name: ticket.name,
                price_cents: ticket.price_cents,
                inventory_remaining: ticket.max_quantity ?? 999,
                has_inventory: ticket.max_quantity !== null,
              },
              1,
            )
          }
        >
          {inCart ? (
            <Check className="h-4 w-4" />
          ) : (
            <ShoppingCart className="h-4 w-4" />
          )}
          {inCart ? `Adicionado (${inCart.quantity})` : "Adicionar ao carrinho"}
        </Button>
      )}
    </div>
  );
}
