import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import {
  createTicketType,
  getEditionAttendeeCount,
  getEditionMyTicket,
  getTicketType,
  listTicketTypes,
  patchTicketType,
} from "@trieoh/univents-api";
import type {
  AttendeeCount,
  MyTicket,
  PatchTicketTypeRequest,
} from "@trieoh/univents-api/schemas";
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
    return createTicketType(editionId, ticketData).then(orvalData<TicketI>);
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
    return patchTicketType(ticketId, ticketData as PatchTicketTypeRequest).then(
      orvalData<TicketI>,
    );
  },
);

/**
 * Fetches all tickets for a specific edition from the server. (Public)
 * @param editionId - the edition id
 * @returns A promise that resolves to an array of Ticket objects.
 */
const getAllTicketsFn = async (editionId: string) => {
  return listTicketTypes(editionId, { public: true }).then(
    orvalData<TicketI[]>,
  );
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
 * Fetches a ticket by ID from the server. (Public)
 * @param ticketId - the ticket id
 * @returns A promise that resolves to a Ticket object.
 */
const getTicketByIdFn = async (ticketId: string) => {
  return getTicketType(ticketId, { public: true })
    .then(orvalData<TicketI | null>)
    .catch(() => null);
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

/**
 * Fetches the authenticated user's held ticket for an edition. (Needs
 * Auth) Returns `null` when the caller holds no ticket — the front falls
 * back to the normal buy flow. Powers the upgrade options shown on the
 * more expensive ticket types (the upgrade action itself is a follow-up).
 */
const getMyTicketFn = async (editionId: string) => {
  const response = await getEditionMyTicket(editionId);
  return (orvalData<MyTicket | null>(response) ?? null) as MyTicket | null;
};

/**
 * Query options for the caller's held ticket for an edition, using
 * TanStack Query. Disable when unauthenticated (a public visitor has no
 * ticket to check).
 * @param editionId - The edition id
 * @param enabled - Whether the query may run (e.g. only when logged in)
 */
export const myTicketQueryOptions = (editionId: string, enabled = true) => {
  return queryOptions({
    queryKey: ticketKeys.myTicket(editionId),
    queryFn: () => getMyTicketFn(editionId),
    enabled: Boolean(editionId) && enabled,
  });
};

export const attendeeCountQueryOptions = (editionId: string) =>
  queryOptions({
    queryKey: ticketKeys.attendeeCount(editionId),
    queryFn: () =>
      getEditionAttendeeCount(editionId, { public: true }).then(
        orvalData<AttendeeCount>,
      ),
  });
