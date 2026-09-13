import { useQueryClient, type QueryClient } from "@trieoh/front-core-solid";
import type { CreateEventRequest, PatchEventRequest } from "@trieoh/univents-api/schemas";
import { orvalData } from "@trieoh/api-client";
import { createEvent, discontinueEvent, patchEvent, publishEvent } from "@trieoh/univents-api";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import { Effect } from "effect";
import { apiEffect, useEffectMutation, toMutationFn } from "@/shared/lib/effect-query";
import type { EventI } from "../model";
import { eventKeys } from "./query-keys";
import { syncEventCaches } from "./cache";

const invalidateEvents = (queryClient: QueryClient) =>
  Promise.all([
    queryClient.invalidateQueries({ queryKey: eventKeys.lists() }),
    queryClient.invalidateQueries({ queryKey: eventKeys.details() }),
  ]);

export const createEventEffect = (data: CreateEventRequest) =>
  apiEffect(() =>
    withSpan("action:event-create", () =>
      createEvent(data).then(orvalData<EventI>),
    ),
  );

export const patchEventEffect = (eventId: string, data: PatchEventRequest) =>
  apiEffect(() =>
    withSpan("action:event-patch", () =>
      patchEvent(eventId, data).then(orvalData<EventI>),
    ),
  );

export const publishEventEffect = (eventId: string) =>
  apiEffect(() =>
    withSpan("action:event-publish", () =>
      publishEvent(eventId).then(orvalData<null>),
    ),
  );

export const discontinueEventEffect = (eventId: string) =>
  apiEffect(() =>
    withSpan("action:event-discontinue", () =>
      discontinueEvent(eventId).then(orvalData<null>),
    ),
  );

export const useCreateEventMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: (data: CreateEventRequest) =>
      createEventEffect(data).pipe(
        Effect.tap((event) =>
          Effect.sync(() => {
            syncEventCaches(queryClient, event);
          }),
        ),
      ),
  });
};

export const usePatchEventMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      eventId,
      data,
    }: {
      eventId: string;
      data: PatchEventRequest;
    }) =>
      patchEventEffect(eventId, data).pipe(
        Effect.tap((event) =>
          Effect.sync(() => {
            syncEventCaches(queryClient, event);
          }),
        ),
      ),
  });
};

export const publishEventMutation = (queryClient: QueryClient) =>
  toMutationFn(
    (eventId: string) =>
      publishEventEffect(eventId).pipe(
        Effect.tap(() =>
          Effect.tryPromise({
            try: () => invalidateEvents(queryClient),
            catch: () => undefined,
          }),
        ),
      ),
    { retryTransient: true },
  );

export const discontinueEventMutation = (queryClient: QueryClient) =>
  toMutationFn(
    (eventId: string) =>
      discontinueEventEffect(eventId).pipe(
        Effect.tap(() =>
          Effect.tryPromise({
            try: () => invalidateEvents(queryClient),
            catch: () => undefined,
          }),
        ),
      ),
    { retryTransient: true },
  );
