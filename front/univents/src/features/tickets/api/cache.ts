import type { QueryClient } from "@tanstack/react-query";
import type { TicketI } from "../model";
import { ticketKeys } from "./query-keys";

function upsertById(tickets: TicketI[] | undefined, ticket: TicketI) {
  if (!tickets) return tickets;
  const index = tickets.findIndex((item) => item.id === ticket.id);

  if (index === -1) return [...tickets, ticket];

  const next = [...tickets];
  next[index] = ticket;
  return next;
}

export function syncTicketCaches(queryClient: QueryClient, ticket: TicketI) {
  queryClient.setQueryData<TicketI[]>(
    ticketKeys.listByEdition(ticket.edition_id),
    (old) => upsertById(old, ticket),
  );
  queryClient.setQueryData<TicketI>(ticketKeys.detail.byId(ticket.id), (old) =>
    old ? ticket : old,
  );
}
