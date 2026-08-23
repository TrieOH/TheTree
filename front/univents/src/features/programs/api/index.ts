import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import {
  createProgram,
  createProgramOccurrence,
  deleteOccurrence,
  deleteProgram,
  deregisterOccurrence,
  listEditionOccurrences,
  listEditionPrograms,
  listMyParticipations,
  patchOccurrence,
  patchProgram,
  registerOccurrence,
} from "@trieoh/univents-api";
import type {
  CreateProgramRequest,
  MyParticipation,
  PatchProgramRequest,
  ProgramParticipation,
} from "@trieoh/univents-api/schemas";
import type {
  OccurrenceCreateOutput,
  OccurrenceI,
  ProgramCreateInput,
  ProgramCreateOutput,
  ProgramI,
} from "../model";
import { programKeys } from "./query-keys";

export const listProgramsFn = createClientOnlyFn((editionId: string) =>
  listEditionPrograms(editionId, { public: true }).then(orvalData<ProgramI[]>),
);

export const listOccurrencesFn = createClientOnlyFn((editionId: string) =>
  listEditionOccurrences(editionId, { public: true }).then(
    orvalData<OccurrenceI[]>,
  ),
);

export const createProgramFn = createClientOnlyFn(
  (editionId: string, data: ProgramCreateInput | ProgramCreateOutput) =>
    withSpan("action:program-create", () =>
      createProgram(editionId, data as CreateProgramRequest).then(
        orvalData<ProgramI>,
      ),
    ),
);

export const patchProgramFn = createClientOnlyFn(
  (id: string, data: ProgramCreateInput | ProgramCreateOutput) =>
    withSpan("action:program-patch", () =>
      patchProgram(id, data as PatchProgramRequest).then(orvalData<ProgramI>),
    ),
);

export const deleteProgramFn = createClientOnlyFn((id: string) =>
  withSpan("action:program-delete", () =>
    deleteProgram(id).then(orvalData<ProgramI>),
  ),
);

export const createOccurrenceFn = createClientOnlyFn(
  (programId: string, data: OccurrenceCreateOutput) =>
    withSpan("action:occurrence-create", () =>
      createProgramOccurrence(programId, data).then(orvalData<OccurrenceI>),
    ),
);

export const patchOccurrenceFn = createClientOnlyFn(
  (id: string, data: OccurrenceCreateOutput) =>
    withSpan("action:occurrence-patch", () =>
      patchOccurrence(id, data).then(orvalData<OccurrenceI>),
    ),
);

export const deleteOccurrenceFn = createClientOnlyFn((id: string) =>
  withSpan("action:occurrence-delete", () =>
    deleteOccurrence(id).then(orvalData<OccurrenceI>),
  ),
);

export const listMyParticipationsFn = createClientOnlyFn((editionId: string) =>
  listMyParticipations(editionId).then(orvalData<MyParticipation[] | null>),
);

export const registerOccurrenceFn = createClientOnlyFn((occurrenceId: string) =>
  registerOccurrence(occurrenceId).then(orvalData<ProgramParticipation>),
);

export const deregisterOccurrenceFn = createClientOnlyFn(
  (occurrenceId: string) =>
    deregisterOccurrence(occurrenceId).then(orvalData<ProgramParticipation>),
);

export const programsQueryOptions = (editionId: string) =>
  queryOptions({
    queryKey: programKeys.byEdition(editionId),
    queryFn: () => listProgramsFn(editionId),
  });

export const occurrencesQueryOptions = (editionId: string) =>
  queryOptions({
    queryKey: programKeys.occurrences(editionId),
    queryFn: () => listOccurrencesFn(editionId),
  });

export const myParticipationsQueryOptions = (
  editionId: string,
  enabled = true,
) =>
  queryOptions({
    queryKey: programKeys.myParticipations(editionId),
    queryFn: () => listMyParticipationsFn(editionId),
    enabled: Boolean(editionId) && enabled,
  });
