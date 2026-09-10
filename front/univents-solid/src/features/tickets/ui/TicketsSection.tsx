import type {
  MyTicket,
  StoreStockItem,
  TicketType,
} from "@trieoh/univents-api/schemas";
import { For, Show, createMemo } from "solid-js";
import { Link } from "@tanstack/solid-router";
import { TicketCard } from "./TicketCard";

export function TicketsSection(props: {
  tickets: TicketType[];
  stock: StoreStockItem[];
  eventSlug: string;
  editionId: string;
  heldTicket: MyTicket | null;
}) {
  const stock = createMemo(
    () => new Map(props.stock.map((item) => [item.id, item.stock])),
  );
  return (
    <Show when={props.tickets.length > 0}>
      <section class="w-full py-10">
        <div class="mb-8 text-center">
          <h2 class="text-3xl font-semibold tracking-tight text-foreground">
            Ingressos
          </h2>
          <p class="mx-auto mt-2 max-w-xl text-sm leading-relaxed text-muted-foreground">
            Escolha o tipo de ingresso que melhor se encaixa na sua experiência
            no evento.
          </p>
        </div>
        <div class="mx-auto flex max-w-5xl justify-center gap-4">
          <For each={props.tickets.slice(0, 3)}>
            {(ticket, index) => (
              <div
                class={
                  index() === 1
                    ? "hidden shrink-0 sm:block"
                    : index() === 2
                      ? "hidden shrink-0 lg:block"
                      : "shrink-0"
                }
              >
                <TicketCard
                  ticket={ticket}
                  editionId={props.editionId}
                  stock={stock().get(ticket.id)}
                  featured={props.tickets.length > 2 && index() === 1}
                  heldTicket={props.heldTicket}
                />
              </div>
            )}
          </For>
        </div>
        <div class="mt-6 text-center">
          <Link
            class="text-sm font-semibold text-primary"
            to="/events/$slug/store"
            params={{ slug: props.eventSlug }}
            search={{ tab: "tickets" }}
          >
            Ver todos os ingressos →
          </Link>
        </div>
      </section>
    </Show>
  );
}
