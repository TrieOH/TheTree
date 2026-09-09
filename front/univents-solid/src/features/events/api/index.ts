import { listPublicEvents } from "@trieoh/univents-api";
import { orvalData } from "@trieoh/api-client";
import type { EventI } from "../model";
import { eventKeys } from "./query-keys";

const getPublicEventsFn = () =>
  listPublicEvents({ public: true }).then(orvalData<EventI[]>);

export const allPublicEventsQueryOptions = () => ({
  queryKey: eventKeys.publicLists(),
  queryFn: getPublicEventsFn,
});
