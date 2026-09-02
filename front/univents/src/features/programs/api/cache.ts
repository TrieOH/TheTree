import type { QueryClient } from "@tanstack/react-query";
import { upsertById } from "../../../shared/lib/query-cache";
import { productKeys } from "../../products/api/query-keys";
import type { OccurrenceI, ProgramI } from "../model";
import { programKeys } from "./query-keys";

export function syncProgramCache(queryClient: QueryClient, program: ProgramI) {
  queryClient.setQueryData<ProgramI[]>(
    programKeys.byEdition(program.edition_id),
    (old) => upsertById(old, program),
  );
}

export function removeProgramCaches(
  queryClient: QueryClient,
  program: ProgramI,
) {
  const occurrences = queryClient.getQueryData<OccurrenceI[]>(
    programKeys.occurrences(program.edition_id),
  );

  for (const occurrence of occurrences ?? []) {
    if (occurrence.program_id === program.id) {
      queryClient.removeQueries({
        queryKey: programKeys.participants(occurrence.id),
      });
    }
  }

  queryClient.setQueryData<ProgramI[]>(
    programKeys.byEdition(program.edition_id),
    (old) => old?.filter((item) => item.id !== program.id),
  );
  queryClient.setQueryData<OccurrenceI[]>(
    programKeys.occurrences(program.edition_id),
    (old) => old?.filter((item) => item.program_id !== program.id),
  );
  void queryClient.invalidateQueries({
    queryKey: productKeys.storeStock(program.edition_id),
  });
}

export function syncOccurrenceCache(
  queryClient: QueryClient,
  occurrence: OccurrenceI,
) {
  queryClient.setQueryData<OccurrenceI[]>(
    programKeys.occurrences(occurrence.edition_id),
    (old) => upsertById(old, occurrence),
  );
  void queryClient.invalidateQueries({
    queryKey: productKeys.storeStock(occurrence.edition_id),
  });
}

export function removeOccurrenceCaches(
  queryClient: QueryClient,
  occurrence: OccurrenceI,
) {
  queryClient.setQueryData<OccurrenceI[]>(
    programKeys.occurrences(occurrence.edition_id),
    (old) => old?.filter((item) => item.id !== occurrence.id),
  );
  queryClient.removeQueries({
    queryKey: programKeys.participants(occurrence.id),
  });
  void queryClient.invalidateQueries({
    queryKey: productKeys.storeStock(occurrence.edition_id),
  });
}
