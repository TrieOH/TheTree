import { useQueryClient } from "@trieoh/front-core-solid";
import { deregisterOccurrence, registerOccurrence } from "@trieoh/univents-api";
import { orvalData } from "@trieoh/api-client";
import type { ProgramParticipation } from "@trieoh/univents-api/schemas";
import { Effect } from "effect";
import { apiEffect, useEffectMutation } from "@/shared/lib/effect-query";
import { invalidateParticipationCache } from "./cache";
import {
  checkInOccurrenceFn,
  createOccurrenceFn,
  createProgramFn,
  deleteOccurrenceFn,
  deleteProgramFn,
  markParticipationAttendedFn,
  patchOccurrenceFn,
  patchProgramFn,
} from "./index";
import { programKeys } from "./query-keys";
import type {
  OccurrenceCreateOutput,
  ProgramCreateInput,
  ProgramCreateOutput,
} from "../model";

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

export const useCreateProgramMutation = (editionId: () => string) => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: (data: ProgramCreateInput | ProgramCreateOutput) =>
      apiEffect(() => createProgramFn(editionId(), data)).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: programKeys.byEdition(editionId()),
            });
          }),
        ),
      ),
  });
};

export const useUpdateProgramMutation = (editionId: () => string) => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      id,
      data,
    }: {
      id: string;
      data: ProgramCreateInput | ProgramCreateOutput;
    }) =>
      apiEffect(() => patchProgramFn(id, data)).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: programKeys.byEdition(editionId()),
            });
          }),
        ),
      ),
  });
};

export const useDeleteProgramMutation = (editionId: () => string) => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: (id: string) =>
      apiEffect(() => deleteProgramFn(id)).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: programKeys.byEdition(editionId()),
            });
            void queryClient.invalidateQueries({
              queryKey: programKeys.occurrences(editionId()),
            });
          }),
        ),
      ),
  });
};

export const useCreateOccurrenceMutation = (editionId: () => string) => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      programId,
      data,
    }: {
      programId: string;
      data: OccurrenceCreateOutput;
    }) =>
      apiEffect(() => createOccurrenceFn(programId, data)).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: programKeys.occurrences(editionId()),
            });
          }),
        ),
      ),
  });
};

export const useUpdateOccurrenceMutation = (editionId: () => string) => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      id,
      data,
    }: {
      id: string;
      data: OccurrenceCreateOutput;
    }) =>
      apiEffect(() => patchOccurrenceFn(id, data)).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: programKeys.occurrences(editionId()),
            });
          }),
        ),
      ),
  });
};

export const useDeleteOccurrenceMutation = (editionId: () => string) => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: (id: string) =>
      apiEffect(() => deleteOccurrenceFn(id)).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: programKeys.occurrences(editionId()),
            });
          }),
        ),
      ),
  });
};

export const useMarkParticipationAttendedMutation = (occurrenceId: () => string) => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: (participationId: string) =>
      apiEffect(() => markParticipationAttendedFn(participationId)).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: programKeys.participants(occurrenceId()),
            });
          }),
        ),
      ),
  });
};

export const useCheckpointCheckInMutation = (occurrenceId: () => string) => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: (attendeeId: string) =>
      apiEffect(() => checkInOccurrenceFn(occurrenceId(), attendeeId)).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: programKeys.participants(occurrenceId()),
            });
          }),
        ),
      ),
  });
};
