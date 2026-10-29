import { useQueryClient } from "@trieoh/front-core-solid";
import {
  checkInOccurrence,
  createProgram,
  createProgramOccurrence,
  deleteOccurrence,
  deleteProgram,
  deregisterOccurrence,
  markParticipationAttended,
  patchOccurrence,
  patchProgram,
  registerOccurrence,
} from "@trieoh/univents-api";
import { orvalData } from "@trieoh/api-client";
import type {
  CreateProgramOccurrenceRequest,
  CreateProgramRequest,
  PatchProgramOccurrenceRequest,
  PatchProgramRequest,
  ProgramParticipation,
} from "@trieoh/univents-api/schemas";
import { Effect } from "effect";
import { apiEffect, useEffectMutation } from "@/shared/lib/effect-query";
import { invalidateParticipationCache } from "./cache";
import { programKeys } from "./query-keys";
import type {
  OccurrenceCreateOutput,
  OccurrenceI,
  ProgramCreateInput,
  ProgramCreateOutput,
  ProgramI,
} from "../model";

export const registerOccurrenceEffect = (occurrenceId: string) =>
  apiEffect(() =>
    registerOccurrence(occurrenceId).then(orvalData<ProgramParticipation>),
  );

export const deregisterOccurrenceEffect = (occurrenceId: string) =>
  apiEffect(() =>
    deregisterOccurrence(occurrenceId).then(orvalData<ProgramParticipation>),
  );

export const createProgramEffect = (
  editionId: string,
  data: ProgramCreateInput | ProgramCreateOutput,
) =>
  apiEffect(() =>
    createProgram(editionId, data as CreateProgramRequest).then(
      orvalData<ProgramI>,
    ),
  );

export const patchProgramEffect = (
  programId: string,
  data: ProgramCreateInput | ProgramCreateOutput,
) =>
  apiEffect(() =>
    patchProgram(programId, data as PatchProgramRequest).then(
      orvalData<ProgramI>,
    ),
  );

export const deleteProgramEffect = (programId: string) =>
  apiEffect(() => deleteProgram(programId).then(orvalData<ProgramI>));

export const createOccurrenceEffect = (
  programId: string,
  data: OccurrenceCreateOutput,
) =>
  apiEffect(() =>
    createProgramOccurrence(
      programId,
      data as CreateProgramOccurrenceRequest,
    ).then(orvalData<OccurrenceI>),
  );

export const patchOccurrenceEffect = (
  occurrenceId: string,
  data: OccurrenceCreateOutput,
) =>
  apiEffect(() =>
    patchOccurrence(
      occurrenceId,
      data as PatchProgramOccurrenceRequest,
    ).then(orvalData<OccurrenceI>),
  );

export const deleteOccurrenceEffect = (occurrenceId: string) =>
  apiEffect(() => deleteOccurrence(occurrenceId).then(orvalData<OccurrenceI>));

export const markParticipationAttendedEffect = (participationId: string) =>
  apiEffect(() =>
    markParticipationAttended(participationId).then(
      orvalData<ProgramParticipation>,
    ),
  );

export const checkInOccurrenceEffect = (
  occurrenceId: string,
  attendeeId: string,
) =>
  apiEffect(() =>
    checkInOccurrence(occurrenceId, { attendee_id: attendeeId }).then(
      orvalData<ProgramParticipation>,
    ),
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
      createProgramEffect(editionId(), data).pipe(
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
      patchProgramEffect(id, data).pipe(
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
      deleteProgramEffect(id).pipe(
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
      createOccurrenceEffect(programId, data).pipe(
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
      patchOccurrenceEffect(id, data).pipe(
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
      deleteOccurrenceEffect(id).pipe(
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
      markParticipationAttendedEffect(participationId).pipe(
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
      checkInOccurrenceEffect(occurrenceId(), attendeeId).pipe(
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
