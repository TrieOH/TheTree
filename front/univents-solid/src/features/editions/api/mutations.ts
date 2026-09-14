import { useQueryClient } from "@trieoh/front-core-solid";
import { orvalData } from "@trieoh/api-client";
import { createEdition, patchEdition, publishEdition } from "@trieoh/univents-api";
import type { CreateEditionRequest, PatchEditionRequest } from "@trieoh/univents-api/schemas";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import { Effect } from "effect";
import { apiEffect, useEffectMutation } from "@/shared/lib/effect-query";
import { normalizeEdition, type EditionI } from "../model";
import { syncEditionCaches } from "./cache";

export interface CreateEditionInput {
  eventId: string;
  data: CreateEditionRequest;
}

export interface PatchEditionInput {
  eventId: string;
  editionId: string;
  data: PatchEditionRequest;
}

export const createEditionEffect = (eventId: string, data: CreateEditionRequest) =>
  apiEffect(() =>
    withSpan("action:edition-create", () =>
      createEdition(eventId, data)
        .then(orvalData<EditionI>)
        .then(normalizeEdition),
    ),
  );

export const patchEditionEffect = (
  eventId: string,
  editionId: string,
  data: PatchEditionRequest,
) =>
  apiEffect(() =>
    withSpan("action:edition-patch", () =>
      patchEdition(eventId, editionId, data)
        .then(orvalData<EditionI>)
        .then(normalizeEdition),
    ),
  );

export const publishEditionEffect = (eventId: string, editionId: string) =>
  apiEffect(() =>
    withSpan("action:edition-publish", () =>
      publishEdition(eventId, editionId).then(orvalData<null>),
    ),
  );

export const useCreateEditionMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ eventId, data }: CreateEditionInput) =>
      createEditionEffect(eventId, data).pipe(
        Effect.tap((edition) =>
          Effect.sync(() => {
            syncEditionCaches(queryClient, edition);
          }),
        ),
      ),
  });
};

export const usePatchEditionMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ eventId, editionId, data }: PatchEditionInput) =>
      patchEditionEffect(eventId, editionId, data).pipe(
        Effect.tap((edition) =>
          Effect.sync(() => {
            syncEditionCaches(queryClient, edition);
          }),
        ),
      ),
  });
};
