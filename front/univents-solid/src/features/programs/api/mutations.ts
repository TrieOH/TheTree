import { useQueryClient } from "@trieoh/front-core-solid";
import { deregisterOccurrence, registerOccurrence } from "@trieoh/univents-api";
import { orvalData } from "@trieoh/api-client";
import type { ProgramParticipation } from "@trieoh/univents-api/schemas";
import { Effect } from "effect";
import { apiEffect, useEffectMutation } from "@/shared/lib/effect-query";
import { invalidateParticipationCache } from "./cache";

export const registerOccurrenceEffect = (occurrenceId: string) =>
  apiEffect(() =>
    registerOccurrence(occurrenceId).then(orvalData<ProgramParticipation>),
  );

export const deregisterOccurrenceEffect = (occurrenceId: string) =>
  apiEffect(() =>
    deregisterOccurrence(occurrenceId).then(orvalData<ProgramParticipation>),
  );

export const useRegisterOccurrenceMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      editionId,
      occurrenceId,
    }: {
      editionId: string;
      occurrenceId: string;
    }) =>
      registerOccurrenceEffect(occurrenceId).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            invalidateParticipationCache(queryClient, editionId);
          }),
        ),
      ),
  });
};

export const useDeregisterOccurrenceMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      editionId,
      occurrenceId,
    }: {
      editionId: string;
      occurrenceId: string;
    }) =>
      deregisterOccurrenceEffect(occurrenceId).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            invalidateParticipationCache(queryClient, editionId);
          }),
        ),
      ),
  });
};
