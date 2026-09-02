import type { MyTicket } from "@trieoh/univents-api/schemas";
import { ArrowUp, Check, ShoppingCart, Star, TicketCheck } from "lucide-react";
import { toast } from "sonner";
import { useCart } from "@/features/products/hooks/use-cart";
import { formatPrice } from "@/shared/lib/money";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import type { TicketI } from "../model";

interface TicketCardProps {
  ticket: TicketI;
  isFeatured?: boolean;
  editionId?: string;
  heldTicket?: MyTicket | null;
  stock?: number | null;
}

export function TicketCard({
  ticket,
  isFeatured,
  editionId,
  heldTicket = null,
  stock,
}: TicketCardProps) {
  const { addItem, items } = useCart(editionId ?? "");
  const inCart = items.find(
    (item) => item.id === ticket.id && item.type === "ticket",
  );
  const anotherTicketInCart = items.some(
    (item) => item.type === "ticket" && item.id !== ticket.id,
  );
  const isFree = ticket.price_cents === 0;
  const isOutOfStock =
    stock !== undefined
      ? stock === 0
      : ticket.max_quantity != null && ticket.max_quantity <= 0;

  const isHeld = heldTicket?.ticket_type.id === ticket.id;
  const hasHeldTicket = !!heldTicket;
  const isUpgrade =
    hasHeldTicket &&
    !isHeld &&
    ticket.access_level > (heldTicket?.ticket_type.access_level ?? -1);
  const upgradePrice = heldTicket
    ? ticket.price_cents - heldTicket.ticket_type.price_cents
    : 0;
  const buttonLabel = isHeld
    ? "Ingresso atual"
    : isUpgrade
      ? "Upgrade em breve"
      : hasHeldTicket
        ? "Você já possui um ingresso"
        : anotherTicketInCart
          ? "Limite de 1 ingresso"
          : isOutOfStock
            ? "Esgotado"
            : inCart
              ? "Adicionado"
              : "Adicionar ao carrinho";

  return (
    <div
      className={cn(
        "relative flex h-64 flex-col w-72 rounded-2xl border overflow-hidden transition-all duration-300 p-5",
        isFeatured
          ? "bg-card border-border/60 lg:scale-105 lg:bg-primary/4 lg:border-primary/25 lg:shadow-lg lg:shadow-primary/5 lg:z-10 hover:border-border hover:shadow-md"
          : "bg-card border-border/60 hover:border-border hover:shadow-md",
      )}
    >
      {/* Badge */}
      <div className="absolute top-0 right-0">
        {isHeld ? (
          <span className="inline-flex items-center gap-1 px-3 py-1 rounded-bl-xl rounded-tr-2xl bg-amber-500 text-white text-[10px] font-bold tracking-wide uppercase">
            <TicketCheck className="w-3 h-3" />
            Seu ingresso
          </span>
        ) : isUpgrade ? (
          <span className="inline-flex items-center gap-1 px-3 py-1 rounded-bl-xl rounded-tr-2xl bg-primary text-primary-foreground text-[10px] font-bold tracking-wide uppercase">
            <ArrowUp className="w-3 h-3" />
            Upgrade
          </span>
        ) : isFree ? (
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
        {isUpgrade && upgradePrice > 0 && (
          <span className="ml-2 text-xs text-muted-foreground">
            +{formatPrice(upgradePrice)}
          </span>
        )}
      </div>
      {stock !== undefined && (
        <p className="mt-1 text-xs text-muted-foreground">
          {stock === null
            ? "Estoque ilimitado"
            : isOutOfStock
              ? "Esgotado"
              : `${stock} disponível${stock === 1 ? "" : "is"}`}
        </p>
      )}

      {/* Description */}
      {ticket.description && (
        <p className="mt-3 line-clamp-3 text-sm leading-relaxed whitespace-pre-line text-muted-foreground">
          {ticket.description}
        </p>
      )}
      {editionId && (
        <Button
          size="sm"
          variant={inCart ? "secondary" : "default"}
          disabled={isOutOfStock || anotherTicketInCart}
          className="mt-auto h-9 w-full gap-2 text-xs font-semibold shadow-sm"
          onClick={() => {
            if (isUpgrade) {
              toast.info(
                "Ingresso adicionado. No checkout, marque-o como presente; upgrades estarão disponíveis em breve.",
              );
            } else if (hasHeldTicket) {
              toast.info(
                "Você já possui um ingresso. Este novo ingresso deverá ser atribuído como presente no checkout.",
              );
            }
            addItem(
              {
                id: ticket.id,
                type: "ticket",
                name: ticket.name,
                price_cents: ticket.price_cents,
                inventory_remaining: stock ?? ticket.max_quantity ?? 999,
                has_inventory: stock !== null && ticket.max_quantity !== null,
                is_upgrade: isUpgrade,
              },
              1,
            );
          }}
        >
          {isOutOfStock ? (
            <TicketCheck className="h-4 w-4" />
          ) : isHeld ? (
            <TicketCheck className="h-4 w-4" />
          ) : isUpgrade ? (
            <ArrowUp className="h-4 w-4" />
          ) : inCart ? (
            <Check className="h-4 w-4" />
          ) : (
            <ShoppingCart className="h-4 w-4" />
          )}
          {buttonLabel}
        </Button>
      )}
    </div>
  );
}
