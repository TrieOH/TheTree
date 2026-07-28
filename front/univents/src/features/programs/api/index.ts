import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import type {
  OccurrenceCreateOutput,
  OccurrenceI,
  ProgramCreateInput,
  ProgramCreateOutput,
  ProgramI,
} from "../model";
import { programKeys } from "./query-keys";
export const listProgramsFn = createClientOnlyFn((editionId: string) =>
  tanstackQueryFetcher<ProgramI[]>(`/editions/${editionId}/programs`),
);
export const listOccurrencesFn = createClientOnlyFn((editionId: string) =>
  tanstackQueryFetcher<OccurrenceI[]>(`/editions/${editionId}/occurrences`),
);
export const createProgramFn = createClientOnlyFn(
  (editionId: string, data: ProgramCreateInput | ProgramCreateOutput) =>
    authFetcher.post<ProgramI>(`/editions/${editionId}/programs`, data),
);
export const patchProgramFn = createClientOnlyFn(
  (id: string, data: ProgramCreateInput | ProgramCreateOutput) =>
    authFetcher.patch<ProgramI>(`/programs/${id}`, data),
);
export const createOccurrenceFn = createClientOnlyFn(
  (programId: string, data: OccurrenceCreateOutput) =>
    authFetcher.post<OccurrenceI>(`/programs/${programId}/occurrences`, data),
);
export const patchOccurrenceFn = createClientOnlyFn(
  (id: string, data: OccurrenceCreateOutput) =>
    authFetcher.patch<OccurrenceI>(`/occurrences/${id}`, data),
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
