import { useMutation, useQueryClient } from "@trieoh/front-core-solid";
import { deregisterOccurrence, registerOccurrence } from "@trieoh/univents-api";
import { orvalData } from "@trieoh/api-client";
import type { ProgramParticipation } from "@trieoh/univents-api/schemas";
import { invalidateParticipationCache } from "./cache";

export const useRegisterOccurrenceMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: async ({ editionId, occurrenceId }: { editionId: string; occurrenceId: string }) => {
    const participation = await registerOccurrence(occurrenceId).then(orvalData<ProgramParticipation>);
    invalidateParticipationCache(queryClient, editionId);
    return participation;
  }});
};

export const useDeregisterOccurrenceMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: async ({ editionId, occurrenceId }: { editionId: string; occurrenceId: string }) => {
    const participation = await deregisterOccurrence(occurrenceId).then(orvalData<ProgramParticipation>);
    invalidateParticipationCache(queryClient, editionId);
    return participation;
  }});
};
