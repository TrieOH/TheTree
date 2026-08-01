import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { authFetcher, publicQueryFetcher } from "@/shared/lib/api/fetch";
import type { TicketCreateOutputI, TicketI } from "../model";
import { ticketKeys } from "./query-keys";

/**
 * Creates a new Ticket on the server. (Needs Auth)
 * @param ticketData - The data for the new ticket.
 * @param editionId - the edition id
 * @returns A promise that resolves to the API response containing the newly created ticket.
 */
export const createTicketFn = createClientOnlyFn(
  (ticketData: TicketCreateOutputI, editionId: string) => {
    return authFetcher.post<TicketI>(
      `/editions/${editionId}/ticket-types`,
      ticketData,
    );
  },
);

/**
 * Patches a Ticket on the server. (Needs Auth)
 * @param ticketData - The data for the updated ticket.
 * @param ticketId - the ticket id
 * @returns A promise that resolves to the API response containing the updated ticket.
 */
export const patchTicketFn = createClientOnlyFn(
  (ticketData: TicketCreateOutputI, ticketId: string) => {
    return authFetcher.patch<TicketI>(`/ticket-types/${ticketId}`, ticketData);
  },
);

/**
 * Fetches all tickets for a specific edition from the server.
 * @param editionId - the edition id
 * @returns A promise that resolves to an array of Ticket objects.
 */
const getAllTicketsFn = async (editionId: string) => {
  return publicQueryFetcher<TicketI[]>(`/editions/${editionId}/ticket-types`);
};

/**
 * Query options for fetching all tickets for a specific edition, using TanStack Query.
 * @param editionId - The edition id
 * @returns An object containing the query key and query function for fetching all tickets for a specific edition.
 */
export const allTicketsQueryOptions = (editionId: string) => {
  return queryOptions({
    queryKey: ticketKeys.listByEdition(editionId),
    queryFn: () => getAllTicketsFn(editionId),
    enabled: Boolean(editionId),
  });
};

/**
 * Fetches a ticket by ID from the server.
 * @param ticketId - the ticket id
 * @returns A promise that resolves to a Ticket object.
 */
const getTicketByIdFn = async (ticketId: string) => {
  return publicQueryFetcher<TicketI | null>(`/ticket-types/${ticketId}`).catch(
    () => null,
  );
};

/**
 * Query options for fetching a ticket by ID, using TanStack Query.
 * @param ticketId - the ticket id
 * @returns An object containing the query key and query function for fetching the ticket.
 */
export const ticketByIdQueryOptions = (ticketId: string) => {
  return queryOptions({
    queryKey: ticketKeys.detail.byId(ticketId),
    queryFn: () => getTicketByIdFn(ticketId),
  });
};
