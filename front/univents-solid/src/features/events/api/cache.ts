import type { QueryClient } from "@trieoh/front-core-solid";
import type { EventI } from "../model";
import { eventKeys } from "./query-keys";

function upsert(events: EventI[] | undefined, event: EventI) {
  if (!events) return events;
  const index = events.findIndex((item) => item.id === event.id);
  return index < 0
    ? [...events, event]
    : events.map((item) => (item.id === event.id ? event : item));
}

export function syncEventCaches(queryClient: QueryClient, event: EventI) {
  queryClient.setQueryData(eventKeys.ownedLists(), (old: EventI[] | undefined) => upsert(old, event));
  queryClient.setQueryData(eventKeys.joinedLists(), (old: EventI[] | undefined) =>
    old?.some((item) => item.id === event.id) ? upsert(old, event) : old,
  );
  queryClient.setQueryData(eventKeys.publicLists(), (old: EventI[] | undefined) =>
    event.status === "draft" ? old?.filter((item) => item.id !== event.id) : upsert(old, event),
  );
  queryClient.setQueryData(eventKeys.detail.publicBySlug(event.slug), event);
  void queryClient.invalidateQueries({ queryKey: eventKeys.catalog("", false).slice(0, 2) });
  void queryClient.invalidateQueries({ queryKey: eventKeys.store("", false).slice(0, 2) });
}
