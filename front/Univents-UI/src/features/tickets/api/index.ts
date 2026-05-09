import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type { TicketCreateI, TicketI } from "../model";
import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";

/**
 * Creates a new Ticket on the server.
 * @param ticketData - The data for the new ticket.
 * @returns A promise that resolves to the API response containing the newly created ticket.
 */
export const createTicketFn = createClientOnlyFn((
  ticketData: TicketCreateI, eventId: string, editionId: string
) => {
  return authFetcher.post<TicketI>(
    `/events/${eventId}/editions/${editionId}/tickets`,
    ticketData
  );
});

/**
 * Fetches all tickets for a specific edition from the server.
 * @returns A promise that resolves to an array of Ticket objects.
 */
export const getAllTicketsFn = createClientOnlyFn(async (eventId: string, editionId: string) => {
  try {
    return await tanstackQueryFetcher<TicketI[]>(`/events/${eventId}/editions/${editionId}/tickets`);
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all tickets for a specific edition, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all tickets for a specific edition.
 */
export const allTicketsQueryOptions = (eventId: string, editionId: string) => {
  return queryOptions({
    queryKey: ['tickets', 'public', eventId, editionId],
    queryFn: () => getAllTicketsFn(eventId, editionId),
  })
}
