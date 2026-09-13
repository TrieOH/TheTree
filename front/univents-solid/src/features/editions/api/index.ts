import { listDraftEditions, listPublicEditions } from "@trieoh/univents-api";
import { orvalData } from "@trieoh/api-client";
import type { Edition } from "@trieoh/univents-api/schemas";
import { normalizeEdition } from "../model";
import { editionKeys } from "./query-keys";

export const getPublicEditionsFn = (eventId: string) =>
  listPublicEditions(eventId, { public: true })
    .then(orvalData<Edition[]>)
    .then((list) => (list ?? []).map(normalizeEdition));

export const getDraftEditionsFn = (eventId: string) =>
  listDraftEditions(eventId)
    .then(orvalData<Edition[]>)
    .then((list) => (list ?? []).map(normalizeEdition));

export const getAllAdminEditionsFn = async (eventId: string) => {
  const [publicEditions, draftEditions] = await Promise.all([
    getPublicEditionsFn(eventId).catch(() => []),
    getDraftEditionsFn(eventId).catch(() => []),
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
