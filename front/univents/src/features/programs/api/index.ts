import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import {
  createProgram,
  createProgramOccurrence,
  deleteOccurrence,
  deleteProgram,
  listEditionOccurrences,
  listEditionPrograms,
  patchOccurrence,
  patchProgram,
} from "@trieoh/univents-api";
import type {
  CreateProgramRequest,
  PatchProgramRequest,
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
    createProgram(editionId, data as CreateProgramRequest).then(
      orvalData<ProgramI>,
    ),
);

export const patchProgramFn = createClientOnlyFn(
  (id: string, data: ProgramCreateInput | ProgramCreateOutput) =>
    patchProgram(id, data as PatchProgramRequest).then(orvalData<ProgramI>),
);

export const deleteProgramFn = createClientOnlyFn((id: string) =>
  deleteProgram(id).then(orvalData<ProgramI>),
);

export const createOccurrenceFn = createClientOnlyFn(
  (programId: string, data: OccurrenceCreateOutput) =>
    createProgramOccurrence(programId, data).then(orvalData<OccurrenceI>),
);

export const patchOccurrenceFn = createClientOnlyFn(
  (id: string, data: OccurrenceCreateOutput) =>
    patchOccurrence(id, data).then(orvalData<OccurrenceI>),
);

export const deleteOccurrenceFn = createClientOnlyFn((id: string) =>
  deleteOccurrence(id).then(orvalData<OccurrenceI>),
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
