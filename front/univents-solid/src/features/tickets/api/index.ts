import { orvalData } from "@trieoh/api-client";
import { getEditionMyTicket, listTicketTypes } from "@trieoh/univents-api";
import type { MyTicket, TicketType } from "@trieoh/univents-api/schemas";
import { ticketKeys } from "./query-keys";

export const ticketsQueryOptions = (editionId: string) => ({
  queryKey: ticketKeys.byEdition(editionId),
  queryFn: () =>
    listTicketTypes(editionId, { public: true }).then(orvalData<TicketType[]>),
});

export const myTicketQueryOptions = (editionId: string) => ({
  queryKey: ticketKeys.mine(editionId),
  queryFn: () =>
    getEditionMyTicket(editionId)
      .then((response) => orvalData<MyTicket | null>(response) ?? null)
      .catch(() => null),
});
