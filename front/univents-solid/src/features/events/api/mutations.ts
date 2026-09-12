import { useMutation, useQueryClient, type QueryClient } from "@trieoh/front-core-solid";
import type { CreateEventRequest, PatchEventRequest } from "@trieoh/univents-api/schemas";
import { orvalData } from "@trieoh/api-client";
import { createEvent, discontinueEvent, patchEvent, publishEvent } from "@trieoh/univents-api";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import type { EventI } from "../model";
import { eventKeys } from "./query-keys";
import { syncEventCaches } from "./cache";

const invalidateEvents = (queryClient: QueryClient) => Promise.all([
  queryClient.invalidateQueries({ queryKey: eventKeys.lists() }),
  queryClient.invalidateQueries({ queryKey: eventKeys.details() }),
]);

export const useCreateEventMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: (data: CreateEventRequest) => withSpan("action:event-create", () => createEvent(data).then(orvalData<EventI>)), onSuccess: (event) => syncEventCaches(queryClient, event) });
};

export const usePatchEventMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: ({ eventId, data }: { eventId: string; data: PatchEventRequest }) => withSpan("action:event-patch", () => patchEvent(eventId, data).then(orvalData<EventI>)), onSuccess: (event) => syncEventCaches(queryClient, event) });
};

export const publishEventMutation = (queryClient: QueryClient) => async (eventId: string) => {
  await withSpan("action:event-publish", () => publishEvent(eventId).then(orvalData<null>));
  await invalidateEvents(queryClient);
};

export const discontinueEventMutation = (queryClient: QueryClient) => async (eventId: string) => {
  await withSpan("action:event-discontinue", () => discontinueEvent(eventId).then(orvalData<null>));
  await invalidateEvents(queryClient);
};
