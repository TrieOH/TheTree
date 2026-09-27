import type { QueryClient } from "@trieoh/front-core-solid";
import type { EditionI } from "../model";
import { editionKeys } from "./query-keys";

function upsert(editions: EditionI[] | undefined, edition: EditionI) {
  if (!editions) return [edition];
  const index = editions.findIndex((item) => item.id === edition.id);
  return index < 0
    ? [...editions, edition]
    : editions.map((item) => (item.id === edition.id ? edition : item));
}

export function syncEditionCaches(queryClient: QueryClient, edition: EditionI) {
  queryClient.setQueryData(
    editionKeys.adminListByEvent(edition.event_id),
    (old: EditionI[] | undefined) => upsert(old, edition),
  );

  queryClient.setQueryData(
    editionKeys.publicListByEvent(edition.event_id),
    (old: EditionI[] | undefined) =>
      edition.is_draft
        ? old?.filter((item) => item.id !== edition.id)
        : upsert(old, edition),
  );

  if (!edition.is_draft) {
    void queryClient.invalidateQueries({
      queryKey: editionKeys.activeByEvent(edition.event_id),
    });
    void queryClient.invalidateQueries({
      queryKey: editionKeys.pastByEvent(edition.event_id),
    });
    void queryClient.invalidateQueries({
      queryKey: editionKeys.upcomingByEvent(edition.event_id),
    });
  }
}
