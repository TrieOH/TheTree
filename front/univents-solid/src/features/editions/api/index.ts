import { listPublicEditions } from "@trieoh/univents-api";
import { orvalData } from "@trieoh/api-client";
import type { EditionI } from "../model";
import { editionKeys } from "./query-keys";

const getPublicEditionsFn = (eventId: string) =>
  listPublicEditions(eventId, { public: true }).then(orvalData<EditionI[]>);

export const allPublicEditionsQueryOptions = (eventId: string) => ({
  queryKey: editionKeys.publicListByEvent(eventId),
  queryFn: () => getPublicEditionsFn(eventId),
});
