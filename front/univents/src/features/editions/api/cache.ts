import type { QueryClient } from "@tanstack/react-query";
import { upsertById } from "../../../shared/lib/query-cache";
import type { EditionI } from "../model";
import { editionKeys } from "./query-keys";

function invalidateDerivedEditions(queryClient: QueryClient, eventId: string) {
  void queryClient.invalidateQueries({
    queryKey: editionKeys.activeByEvent(eventId),
  });
  void queryClient.invalidateQueries({
    queryKey: editionKeys.pastByEvent(eventId),
  });
  void queryClient.invalidateQueries({
    queryKey: editionKeys.upcomingByEvent(eventId),
  });
}

export function syncEditionCaches(queryClient: QueryClient, edition: EditionI) {
  queryClient.setQueryData<EditionI[]>(
    editionKeys.adminListByEvent(edition.event_id),
    (old) => upsertById(old, edition),
  );

  queryClient.setQueryData<EditionI[]>(
    editionKeys.publicListByEvent(edition.event_id),
    (old) =>
      edition.is_draft
        ? old?.filter((item) => item.id !== edition.id)
        : upsertById(old, edition),
  );

  if (!edition.is_draft) {
    invalidateDerivedEditions(queryClient, edition.event_id);
  }
}
