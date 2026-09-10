import { orvalData } from "@trieoh/api-client";
import {
  deregisterOccurrence,
  getOccurrence,
  listEditionOccurrences,
  listEditionPrograms,
  listMyParticipations,
  registerOccurrence,
} from "@trieoh/univents-api";
import type {
  MyParticipation,
  Program,
  ProgramOccurrence,
  ProgramParticipation,
} from "@trieoh/univents-api/schemas";
import { programKeys } from "./query-keys";

export const programsQueryOptions = (editionId: string) => ({
  queryKey: programKeys.byEdition(editionId),
  queryFn: () =>
    listEditionPrograms(editionId, { public: true }).then(orvalData<Program[]>),
});

export const occurrencesQueryOptions = (editionId: string) => ({
  queryKey: programKeys.occurrences(editionId),
  queryFn: () =>
    listEditionOccurrences(editionId, { public: true }).then(
      orvalData<ProgramOccurrence[]>,
    ),
});

export const occurrenceQueryOptions = (occurrenceId: string) => ({
  queryKey: programKeys.occurrence(occurrenceId),
  queryFn: () =>
    getOccurrence(occurrenceId, { public: true }).then(
      orvalData<ProgramOccurrence>,
    ),
});

export const myParticipationsQueryOptions = (editionId: string) => ({
  queryKey: programKeys.mine(editionId),
  queryFn: () =>
    listMyParticipations(editionId).then(
      (response) => orvalData<MyParticipation[] | null>(response) ?? [],
    ),
});

export const registerOccurrenceFn = (occurrenceId: string) =>
  registerOccurrence(occurrenceId).then(orvalData<ProgramParticipation>);

export const deregisterOccurrenceFn = (occurrenceId: string) =>
  deregisterOccurrence(occurrenceId).then(orvalData<ProgramParticipation>);
