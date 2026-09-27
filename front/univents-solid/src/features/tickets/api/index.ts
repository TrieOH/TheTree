import { orvalData } from "@trieoh/api-client";
import {
  createTicketType,
  getEditionAttendeeCount,
  getEditionMyTicket,
  listTicketTypes,
  patchTicketType,
} from "@trieoh/univents-api";
import type {
  AttendeeCount,
  CreateTicketTypeRequest,
  MyTicket,
  PatchTicketTypeRequest,
  TicketType,
} from "@trieoh/univents-api/schemas";
import { ticketKeys } from "./query-keys";

export const ticketsQueryOptions = (editionId: string) => ({
  queryKey: ticketKeys.byEdition(editionId),
  queryFn: () =>
    listTicketTypes(editionId, { public: true }).then(orvalData<TicketType[]>),
});

export const allTicketsQueryOptions = ticketsQueryOptions;

export const myTicketQueryOptions = (editionId: string) => ({
  queryKey: ticketKeys.mine(editionId),
  queryFn: () =>
    getEditionMyTicket(editionId)
      .then((response) => orvalData<MyTicket | null>(response) ?? null)
      .catch(() => null),
});

export const attendeeCountQueryOptions = (editionId: string) => ({
  queryKey: ticketKeys.attendeeCount(editionId),
  queryFn: () =>
    getEditionAttendeeCount(editionId, { public: true }).then(
      orvalData<AttendeeCount>,
    ),
});

export const createTicketFn = (
  data: CreateTicketTypeRequest,
  editionId: string,
) =>
  createTicketType(editionId, data).then(orvalData<TicketType>);

export const patchTicketFn = (
  data: PatchTicketTypeRequest,
  ticketId: string,
) =>
  patchTicketType(ticketId, data).then(orvalData<TicketType>);
