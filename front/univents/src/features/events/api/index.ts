import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import {
  createEvent,
  discontinueEvent,
  getEventBySlug,
  listJoinedEvents,
  listOwnedEvents,
  listPublicEvents,
  patchEvent,
  publishEvent,
} from "@trieoh/univents-api";
import type { EventCreateOutputI, EventI } from "../model";
import { eventKeys } from "./query-keys";

/**
 * Creates a new Event on the server.
 * @param eventData - The data for the new event.
 * @returns A promise that resolves to the API response containing the newly created event.
 */
export const createEventFn = createClientOnlyFn(
  (eventData: EventCreateOutputI) =>
    withSpan("action:event-create", () =>
      createEvent(eventData).then(orvalData<EventI>),
    ),
);

export const patchEventFn = createClientOnlyFn(
  (eventId: string, eventData: EventCreateOutputI) =>
    withSpan("action:event-patch", () =>
      patchEvent(eventId, eventData).then(orvalData<EventI>),
    ),
);

/**
 * Publish a Event on the server.
 * @param eventId - The event id
 * @returns A promise that resolves to the API null response.
 */
export const publishEventFn = createClientOnlyFn((eventId: string) => {
  return withSpan("action:event-publish", () =>
    publishEvent(eventId).then(orvalData<null>),
  );
});

export const discontinueEventFn = createClientOnlyFn((eventId: string) => {
  return withSpan("action:event-discontinue", () =>
    discontinueEvent(eventId).then(orvalData<null>),
  );
});

/**
 * Fetches all public events from the server.
 * @returns A promise that resolves to an array of Event objects.
 */
const getPublicEventsFn = async () => {
  return listPublicEvents({ public: true }).then(orvalData<EventI[]>);
};

/**
 * Query options for fetching events, using TanStack Query.
 * @returns An object containing the query key and query function for fetching events.
 */
export const allPublicEventsQueryOptions = () => {
  return queryOptions({
    queryKey: eventKeys.publicLists(),
    queryFn: getPublicEventsFn,
  });
};

/**
 * Fetches all own events from the server.
 * @returns A promise that resolves to an array of Event objects.
 */
const getOwnEventsFn = createClientOnlyFn(async () => {
  return listOwnedEvents().then(orvalData<EventI[]>);
});

/**
 * Query options for fetching own events, using TanStack Query.
 * @returns An object containing the query key and query function for fetching own events.
 */
export const allOwnEventsQueryOptions = () => {
  return queryOptions({
    queryKey: eventKeys.ownLists(),
    queryFn: getOwnEventsFn,
  });
};

/**
 * Fetches all joined events from the server.
 * @returns A promise that resolves to an array of Event objects.
 */
const getJoinedEventsFn = createClientOnlyFn(async () => {
  return listJoinedEvents().then(orvalData<EventI[]>);
});

/**
 * Query options for fetching joined events, using TanStack Query.
 * @returns An object containing the query key and query function for fetching joined events.
 */
export const allJoinedEventsQueryOptions = () => {
  return queryOptions({
    queryKey: eventKeys.joinedLists(),
    queryFn: getJoinedEventsFn,
  });
};

/**
 * Fetches a public events from the server.
 * @param slug - the event slug
 * @returns A promise that resolves to an Event object.
 */
const getPublicEventBySlugFn = async (slug: string) => {
  return getEventBySlug(slug, { public: true })
    .then(orvalData<EventI | null>)
    .catch(() => null);
};

/**
 * Query options for fetching event, using TanStack Query.
 * @param slug - the event slug
 * @returns An object containing the query key and query function for fetching the event.
 */
export const publicEventBySlugQueryOptions = (slug: string) => {
  return queryOptions({
    queryKey: eventKeys.detail.publicBySlug(slug),
    queryFn: () => getPublicEventBySlugFn(slug),
  });
};
