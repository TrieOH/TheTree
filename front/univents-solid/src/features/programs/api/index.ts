import { orvalData } from "@trieoh/api-client";
import {
  checkInOccurrence,
  createProgram,
  createProgramOccurrence,
  deleteOccurrence,
  deleteProgram,
  deregisterOccurrence,
  getOccurrence,
  listEditionOccurrences,
  listEditionPrograms,
  listMyParticipations,
  listOccurrenceParticipants,
  markParticipationAttended,
  patchOccurrence,
  patchProgram,
  registerOccurrence,
} from "@trieoh/univents-api";
import type {
  CreateProgramOccurrenceRequest,
  CreateProgramRequest,
  MyParticipation,
  MyTicket,
  PatchProgramOccurrenceRequest,
  PatchProgramRequest,
  Program,
  ProgramOccurrence,
  ProgramParticipant,
  ProgramParticipation,
} from "@trieoh/univents-api/schemas";
import { programKeys } from "./query-keys";
import type { QueryClient } from "@tanstack/query-core";
import { publicEventBySlugQueryOptions } from "@/features/events/api";
import { allPublicEditionsQueryOptions } from "@/features/editions/api";
import type { EventI } from "@/features/events/model";
import type { EditionI } from "@/features/editions/model";
import { myTicketQueryOptions } from "@/features/tickets/api";
import type {
  OccurrenceCreateOutput,
  OccurrenceI,
  ProgramCreateInput,
  ProgramCreateOutput,
  ProgramI,
} from "../model";

export type ProgramPageData = {
  event: EventI;
  edition: EditionI | null;
  programs: Program[];
  occurrences: ProgramOccurrence[];
  heldTicket: MyTicket | null;
  participations: MyParticipation[];
  authenticated: boolean;
};

export const programPageQueryOptions = (
  slug: string,
  authenticated: boolean,
  client: QueryClient,
) => ({
  queryKey: programKeys.page(slug, authenticated),
  queryFn: async (): Promise<ProgramPageData | null> => {
    const event = await publicEventBySlugQueryOptions(slug).queryFn();
    if (!event) return null;
    const editions = await allPublicEditionsQueryOptions(event.id).queryFn();
    const edition = editions.find((item) => new Date(item.ends_at).getTime() >= Date.now()) ?? editions.at(-1);
    if (!edition) return { event, edition: null, programs: [], occurrences: [], heldTicket: null, participations: [], authenticated };
    const [programs, occurrences, heldTicket, participations] = await Promise.all([
      client.fetchQuery(programsQueryOptions(edition.id)),
      client.fetchQuery(occurrencesQueryOptions(edition.id)),
      authenticated ? client.fetchQuery(myTicketQueryOptions(edition.id)) : null,
      authenticated ? client.fetchQuery(myParticipationsQueryOptions(edition.id)) : [],
    ]);
    return { event, edition, programs, occurrences, heldTicket, participations, authenticated };
  },
});

export const programsQueryOptions = (editionId: string) => ({
  queryKey: programKeys.byEdition(editionId),
  queryFn: () =>
    listEditionPrograms(editionId).then(orvalData<ProgramI[]>),
});

export const occurrencesQueryOptions = (editionId: string) => ({
  queryKey: programKeys.occurrences(editionId),
  queryFn: () =>
    listEditionOccurrences(editionId).then(
      orvalData<OccurrenceI[]>,
    ),
});

export const occurrenceQueryOptions = (occurrenceId: string) => ({
  queryKey: programKeys.occurrence(occurrenceId),
  queryFn: () =>
    getOccurrence(occurrenceId, { public: true }).then(
      orvalData<ProgramOccurrence>,
    ),
});

export const occurrenceParticipantsQueryOptions = (occurrenceId: string) => ({
  queryKey: programKeys.participants(occurrenceId),
  queryFn: () => listOccurrenceParticipantsFn(occurrenceId),
  enabled: Boolean(occurrenceId),
});

export const myParticipationsQueryOptions = (editionId: string) => ({
  queryKey: programKeys.mine(editionId),
  queryFn: () =>
    listMyParticipations(editionId).then(
      (response) => orvalData<MyParticipation[] | null>(response) ?? [],
    ),
});

export const createProgramFn = (
  editionId: string,
  data: ProgramCreateInput | ProgramCreateOutput,
) =>
  createProgram(editionId, data as CreateProgramRequest).then(
    orvalData<ProgramI>,
  );

export const patchProgramFn = (
  programId: string,
  data: ProgramCreateInput | ProgramCreateOutput,
) =>
  patchProgram(programId, data as PatchProgramRequest).then(
    orvalData<ProgramI>,
  );

export const deleteProgramFn = (programId: string) =>
  deleteProgram(programId).then(orvalData<ProgramI>);

export const createOccurrenceFn = (
  programId: string,
  data: OccurrenceCreateOutput,
) =>
  createProgramOccurrence(
    programId,
    data as CreateProgramOccurrenceRequest,
  ).then(orvalData<OccurrenceI>);

export const patchOccurrenceFn = (
  occurrenceId: string,
  data: OccurrenceCreateOutput,
) =>
  patchOccurrence(
    occurrenceId,
    data as PatchProgramOccurrenceRequest,
  ).then(orvalData<OccurrenceI>);

export const deleteOccurrenceFn = (occurrenceId: string) =>
  deleteOccurrence(occurrenceId).then(orvalData<OccurrenceI>);

export const registerOccurrenceFn = (occurrenceId: string) =>
  registerOccurrence(occurrenceId).then(orvalData<ProgramParticipation>);

export const deregisterOccurrenceFn = (occurrenceId: string) =>
  deregisterOccurrence(occurrenceId).then(orvalData<ProgramParticipation>);

export const listOccurrenceParticipantsFn = (occurrenceId: string) =>
  listOccurrenceParticipants(occurrenceId).then(
    orvalData<ProgramParticipant[]>,
  );

export const markParticipationAttendedFn = (participationId: string) =>
  markParticipationAttended(participationId).then(
    orvalData<ProgramParticipation>,
  );

export const checkInOccurrenceFn = (occurrenceId: string, attendeeId: string) =>
  checkInOccurrence(occurrenceId, { attendee_id: attendeeId }).then(
    orvalData<ProgramParticipation>,
  );
