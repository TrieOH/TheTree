import type { QueryClient } from "@tanstack/react-query";
import { upsertById } from "../../../shared/lib/query-cache";
import type { EventI } from "../model";
import { eventKeys } from "./query-keys";

function removeById(events: EventI[] | undefined, eventId: string) {
  return events?.filter((item) => item.id !== eventId);
}

function shouldBePublic(event: Pick<EventI, "status">) {
  return event.status !== "draft";
}

export function syncEventCaches(queryClient: QueryClient, event: EventI) {
  queryClient.setQueryData<EventI[]>(eventKeys.ownLists(), (old) =>
    upsertById(old, event),
  );

  queryClient.setQueryData<EventI[]>(eventKeys.publicLists(), (old) => {
    if (shouldBePublic(event)) return upsertById(old, event);
    return removeById(old, event.id);
  });

  queryClient.setQueryData<EventI[]>(eventKeys.joinedLists(), (old) =>
    old?.some((item) => item.id === event.id) ? upsertById(old, event) : old,
  );

  queryClient.setQueryData<EventI>(
    eventKeys.detail.publicBySlug(event.slug),
    (old) => (old ? event : old),
  );
}
