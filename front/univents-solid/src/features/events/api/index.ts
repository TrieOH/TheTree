import {
  getEventBySlug,
  listJoinedEvents,
  listOwnedEvents,
  listPublicEvents,
} from "@trieoh/univents-api";
import { orvalData, orvalItems } from "@trieoh/api-client";
import type { EventI } from "../model";
import { eventKeys } from "./query-keys";
import "./upload-association";

const getPublicEventsFn = () =>
  listPublicEvents({ public: true }).then(orvalItems<EventI>);

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


const getOwnedEventsFn = async () =>
  listOwnedEvents({ public: false }).then(orvalItems<EventI>);

const getJoinedEventsFn = async () =>
  listJoinedEvents({ public: false }).then(orvalItems<EventI>);

export const allOwnEventsQueryOptions = () => ({
  queryKey: eventKeys.ownedLists(),
  queryFn: getOwnedEventsFn,
});

export const allJoinedEventsQueryOptions = () => ({
  queryKey: eventKeys.joinedLists(),
  queryFn: getJoinedEventsFn,
});
