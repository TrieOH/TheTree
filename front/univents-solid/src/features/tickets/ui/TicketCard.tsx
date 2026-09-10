import type { MyTicket, TicketType } from "@trieoh/univents-api/schemas";
import type { JSX } from "@solidjs/web";
import { untrack } from "solid-js";
import CheckIcon from "~icons/lucide/check";
import StarIcon from "~icons/lucide/star";
import TicketIcon from "~icons/lucide/ticket-check";
import { useCart } from "@/features/products/hooks/use-cart";

const Check = CheckIcon as unknown as () => JSX.Element;
const Star = StarIcon as unknown as () => JSX.Element;
const Ticket = TicketIcon as unknown as () => JSX.Element;

const price = (cents: number) =>
  cents === 0
    ? "Grátis"
    : new Intl.NumberFormat("pt-BR", {
      style: "currency",
      currency: "BRL",
    }).format(cents / 100);

export function TicketCard(props: {
  ticket: TicketType;
  featured?: boolean;
  stock?: number | null;
  editionId: string;
  heldTicket?: MyTicket | null;
}) {
  const cart = useCart(untrack(() => props.editionId));
  const free = () => props.ticket.price_cents === 0;
  const soldOut = () =>
    props.stock === 0 ||
    (props.stock === undefined && props.ticket.max_quantity === 0);
  const held = () => props.heldTicket?.ticket_type.id === props.ticket.id;
  const inCart = () =>
    cart
      .items()
      .some((item) => item.type === "ticket" && item.id === props.ticket.id);
  const anotherTicket = () =>
    cart
      .items()
      .some((item) => item.type === "ticket" && item.id !== props.ticket.id);

  return (
    <article
      class={`relative flex h-64 w-72 flex-col overflow-hidden rounded-2xl border p-5 transition-all duration-300 hover:border-border hover:shadow-md ${props.featured
        ? "border-primary/25 bg-primary/4 shadow-lg shadow-primary/5 lg:scale-105"
        : "border-border/60 bg-card"
        }`}
    >
      <span
        class={`absolute right-0 top-0 inline-flex items-center gap-1 rounded-bl-xl rounded-tr-2xl px-3 py-1 text-[10px] font-bold uppercase tracking-wide ${held()
          ? "bg-amber-500 text-white"
          : free()
            ? "bg-emerald-500 text-white"
            : "bg-primary text-primary-foreground"
          }`}
      >
        <span class="size-3 [&>svg]:size-3 [&>svg]:fill-current">
          {held() || free() ? <Check /> : <Star />}
        </span>
        {held() ? "Seu ingresso" : free() ? "Gratuito" : "Pago"}
      </span>
      <h3 class="pr-16 text-base font-semibold leading-tight text-foreground">
        {props.ticket.name}
      </h3>
      <p
        class={`mt-2 text-2xl font-bold ${free() ? "text-emerald-600 dark:text-emerald-400" : "text-foreground"
          }`}
      >
        {price(props.ticket.price_cents)}
      </p>
      {props.stock !== undefined && (
        <p class="mt-1 text-xs text-muted-foreground">
          {props.stock === null
            ? "Estoque ilimitado"
            : soldOut()
              ? "Esgotado"
              : `${props.stock} disponível${props.stock === 1 ? "" : "is"}`}
        </p>
      )}
      {props.ticket.description && (
        <p class="mt-3 line-clamp-3 whitespace-pre-line text-sm leading-relaxed text-muted-foreground">
          {props.ticket.description}
        </p>
      )}
      <button
        type="button"
        disabled={soldOut() || held() || anotherTicket()}
        class={`mt-auto inline-flex h-9 w-full items-center justify-center gap-2 rounded-md px-3 text-xs font-semibold ${soldOut() || held() || anotherTicket()
          ? "bg-muted text-muted-foreground"
          : "bg-primary text-primary-foreground"
          }`}
        onClick={() =>
          untrack(() =>
            cart.add({
              id: props.ticket.id,
              type: "ticket",
              name: props.ticket.name,
              price_cents: props.ticket.price_cents,
              stock: props.stock ?? props.ticket.max_quantity ?? null,
            })
          )
        }
      >
        <span class="size-4">
          <Ticket />
        </span>
        {held()
          ? "Ingresso atual"
          : anotherTicket()
            ? "Limite de 1 ingresso"
            : soldOut()
              ? "Esgotado"
              : inCart()
                ? "Adicionado"
                : "Adicionar ao carrinho"}
      </button>
    </article>
  );
}
