import type { QueryClient } from "@tanstack/react-query";
import { upsertById } from "../../../shared/lib/query-cache";
import type { TicketI } from "../model";
import { ticketKeys } from "./query-keys";

export function syncTicketCaches(queryClient: QueryClient, ticket: TicketI) {
  queryClient.setQueryData<TicketI[]>(
    ticketKeys.listByEdition(ticket.edition_id),
    (old) => upsertById(old, ticket),
  );
  queryClient.setQueryData<TicketI>(ticketKeys.detail.byId(ticket.id), (old) =>
    old ? ticket : old,
  );
}
