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
import type { CreateEventRequest, PatchEventRequest } from "@trieoh/univents-api/schemas";
import { orvalData } from "@trieoh/api-client";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import type { EventI } from "../model";
import { eventKeys } from "./query-keys";

const getPublicEventsFn = () =>
  listPublicEvents({ public: true }).then(orvalData<EventI[]>);

export const allPublicEventsQueryOptions = () => ({
  queryKey: eventKeys.publicLists(),
  queryFn: getPublicEventsFn,
});

const getPublicEventBySlugFn = async (slug: string) =>
  getEventBySlug(slug, { public: true }).then(orvalData<EventI | null>);

export const publicEventBySlugQueryOptions = (slug: string) => ({
  queryKey: eventKeys.detail.publicBySlug(slug),
  queryFn: () => getPublicEventBySlugFn(slug),
});

// --- admin -------------------------------------------------------------

const getOwnedEventsFn = async () =>
  listOwnedEvents({ public: false }).then(orvalData<EventI[]>);

const getJoinedEventsFn = async () =>
  listJoinedEvents({ public: false }).then(orvalData<EventI[]>);

export const allOwnEventsQueryOptions = () => ({
  queryKey: eventKeys.ownedLists(),
  queryFn: getOwnedEventsFn,
});

export const allJoinedEventsQueryOptions = () => ({
  queryKey: eventKeys.joinedLists(),
  queryFn: getJoinedEventsFn,
});

export const createEventFn = (data: CreateEventRequest) =>
  withSpan("action:event-create", () => createEvent(data).then(orvalData<EventI>));

export const patchEventFn = (eventId: string, data: PatchEventRequest) =>
  withSpan("action:event-patch", () =>
    patchEvent(eventId, data).then(orvalData<EventI>),
  );

export const publishEventFn = (eventId: string) =>
  withSpan("action:event-publish", () =>
    publishEvent(eventId).then(orvalData<null>),
  );

export const discontinueEventFn = (eventId: string) =>
  withSpan("action:event-discontinue", () =>
    discontinueEvent(eventId).then(orvalData<null>),
  );
