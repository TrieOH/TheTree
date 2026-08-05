import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import {
  createEdition,
  getActiveEdition,
  listDraftEditions,
  listPastEditions,
  listPublicEditions,
  listUpcomingEditions,
  patchEdition,
  publishEdition,
} from "@trieoh/univents-api";
import {
  type EditionApiI,
  type EditionCreateOutputI,
  type EditionPatchOutputI,
  normalizeEdition,
} from "../model";
import { editionKeys } from "./query-keys";

export const createEditionFn = createClientOnlyFn(
  async (editionData: EditionCreateOutputI, eventId: string) => {
    const edition = await createEdition(eventId, editionData).then(
      orvalData<EditionApiI>,
    );
    return normalizeEdition(edition);
  },
);

export const patchEditionFn = createClientOnlyFn(
  async (
    eventId: string,
    editionId: string,
    editionData: EditionPatchOutputI,
  ) => {
    const edition = await patchEdition(eventId, editionId, editionData).then(
      orvalData<EditionApiI>,
    );
    return normalizeEdition(edition);
  },
);

export const publishEditionFn = createClientOnlyFn(
  (eventId: string, editionId: string) =>
    publishEdition(eventId, editionId).then(orvalData<null>),
);

export const getAllPublicEditionsFn = async (eventId: string) => {
  const editions = await listPublicEditions(eventId, { public: true }).then(
    orvalData<EditionApiI[]>,
  );
  return editions.filter((edition) => !edition.is_draft).map(normalizeEdition);
};

const getActiveEditionPublic = async (eventId: string) => {
  try {
    const edition = await getActiveEdition(eventId, { public: true }).then(
      orvalData<EditionApiI>,
    );
    return normalizeEdition(edition);
  } catch {
    return null;
  }
};

const getPublicEditions = async (
  eventId: string,
  kind: "past" | "upcoming",
) => {
  const fn = kind === "past" ? listPastEditions : listUpcomingEditions;
  const editions = await fn(eventId, { public: true }).then(
    orvalData<EditionApiI[]>,
  );
  return editions.map(normalizeEdition);
};

export const activeEditionQueryOptions = (eventId: string) =>
  queryOptions({
    queryKey: editionKeys.activeByEvent(eventId),
    queryFn: () => getActiveEditionPublic(eventId),
  });

export const pastEditionsQueryOptions = (eventId: string) =>
  queryOptions({
    queryKey: editionKeys.pastByEvent(eventId),
    queryFn: () => getPublicEditions(eventId, "past"),
  });

export const upcomingEditionsQueryOptions = (eventId: string) =>
  queryOptions({
    queryKey: editionKeys.upcomingByEvent(eventId),
    queryFn: () => getPublicEditions(eventId, "upcoming"),
  });

export const allPublicEditionsQueryOptions = (eventId: string) => {
  return queryOptions({
    queryKey: editionKeys.publicListByEvent(eventId),
    queryFn: () => getAllPublicEditionsFn(eventId),
  });
};

export const getDraftEditionsFn = createClientOnlyFn(
  async (eventId: string) => {
    const editions = await listDraftEditions(eventId).then(
      orvalData<EditionApiI[]>,
    );
    return editions.map(normalizeEdition);
  },
);

export const getAllAdminEditionsFn = createClientOnlyFn(
  async (eventId: string) => {
    const [publicEditions, draftEditions] = await Promise.all([
      getAllPublicEditionsFn(eventId),
      getDraftEditionsFn(eventId),
    ]);

    return [...publicEditions, ...draftEditions].filter(
      (edition, index, editions) =>
        editions.findIndex((candidate) => candidate.id === edition.id) ===
        index,
    );
  },
);

export const allAdminEditionsQueryOptions = (eventId: string) => {
  return queryOptions({
    queryKey: editionKeys.adminListByEvent(eventId),
    queryFn: () => getAllAdminEditionsFn(eventId),
  });
};
