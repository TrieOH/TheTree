import type { QueryClient } from "@tanstack/query-core";
import { listDraftEditions, listPublicEditions } from "@trieoh/univents-api";
import { orvalItems } from "@trieoh/api-client";
import type { Edition } from "@trieoh/univents-api/schemas";
import { userBadgesQueryOptions } from "@/features/badges/api";
import { allProfileBadges } from "@/features/badges/model/profile-badges";
import { allPublicEventsQueryOptions } from "@/features/events/api";
import { degradeUnlessTransient } from "@/shared/lib/degrade";
import { normalizeEdition } from "../model";
import { editionKeys } from "./query-keys";
import "./upload-association";

export const getPublicEditionsFn = (eventId: string) =>
  listPublicEditions(eventId, { public: true })
    .then(orvalItems<Edition>)
    .then((list) => (list ?? []).map(normalizeEdition));

export const getDraftEditionsFn = (eventId: string) =>
  listDraftEditions(eventId)
    .then(orvalItems<Edition>)
    .then((list) => (list ?? []).map(normalizeEdition));

export const getAllAdminEditionsFn = async (eventId: string) => {
  // A member without draft access gets a 403: that list is simply empty for
  // them. An unreachable API is not an empty list, though, so it propagates.
  const [publicEditions, draftEditions] = await Promise.all([
    degradeUnlessTransient(getPublicEditionsFn(eventId), []),
    degradeUnlessTransient(getDraftEditionsFn(eventId), []),
  ]);

  return [...publicEditions, ...draftEditions].filter(
    (edition, index, editions) =>
      editions.findIndex((candidate) => candidate.id === edition.id) === index,
  );
};

export const allPublicEditionsQueryOptions = (eventId: string) => ({
  queryKey: editionKeys.publicListByEvent(eventId),
  queryFn: () => getPublicEditionsFn(eventId),
});

export const allAdminEditionsQueryOptions = (eventId: string) => ({
  queryKey: editionKeys.adminListByEvent(eventId),
  queryFn: () => getAllAdminEditionsFn(eventId),
});

/**
 * The edition behind a certificate or badge, found among the public events'
 * editions. The per-event lists go through the query cache, so the lookups for
 * a page full of certificates share requests instead of repeating them.
 */
export const editionLocationQueryOptions = (
  editionId: string | undefined,
  events: Array<{ id: string }>,
) => ({
  queryKey: editionKeys.detail(editionId ?? ""),
  queryFn: async ({ client }: { client: QueryClient }) => {
    if (!editionId) return null;
    const lists = await Promise.all(
      events.map((event) =>
        degradeUnlessTransient(
          client.fetchQuery(allPublicEditionsQueryOptions(event.id)),
          [],
        ),
      ),
    );
    return lists.flat().find((item) => item.id === editionId) ?? null;
  },
  enabled: Boolean(editionId && events.length > 0),
});

/** How many events' edition lists are asked for at once when resolving locations. */
const LOCATION_LOOKUP_BATCH = 4;

/**
 * The location of every edition the user's badges were issued for, by edition
 * id. A badge names its edition but not its event, so the owning event has to
 * be found among the public ones. That is searched a few events at a time and
 * stops as soon as every edition is found, instead of asking for every event's
 * editions up front. (A badge carrying its event, or an edition lookup by id,
 * would make this one request; that needs the API.)
 */
export const badgeEditionLocationsQueryOptions = (actorId: string) => ({
  queryKey: [...editionKeys.all, "badge-locations", actorId] as const,
  queryFn: async ({ client }: { client: QueryClient }) => {
    const [groups, events] = await Promise.all([
      client.fetchQuery(userBadgesQueryOptions(actorId)),
      client.fetchQuery(allPublicEventsQueryOptions()),
    ]);
    const wanted = new Set(allProfileBadges(groups).map((badge) => badge.edition_id));
    const locations: Record<string, string> = {};

    for (
      let start = 0;
      start < events.length && wanted.size > 0;
      start += LOCATION_LOOKUP_BATCH
    ) {
      const lists = await Promise.all(
        events.slice(start, start + LOCATION_LOOKUP_BATCH).map((event) =>
          degradeUnlessTransient(
            client.fetchQuery(allPublicEditionsQueryOptions(event.id)),
            [],
          ),
        ),
      );
      for (const edition of lists.flat()) {
        if (!wanted.delete(edition.id)) continue;
        locations[edition.id] = edition.location_name ?? "";
      }
    }
    return locations;
  },
});
