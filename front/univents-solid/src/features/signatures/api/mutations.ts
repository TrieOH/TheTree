import { orvalData } from "@trieoh/api-client";
import { useQueryClient } from "@trieoh/front-core-solid";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import {
  cancelSignatureRequest,
  createSignature,
  createSignatureRequest,
  deleteSignature,
} from "@trieoh/univents-api";
import type {
  AddSignatureRequest,
  CancelSignatureRequestBody,
  CreateSignatureRequestRequest,
} from "@trieoh/univents-api/schemas";
import { Effect } from "effect";
import { apiEffect, useEffectMutation } from "@/shared/lib/effect-query";
import type {
  SignatureCreateOutputI,
  SignatureI,
  SignatureRequestCreateI,
  SignatureRequestI,
} from "../model";
import { signatureKeys } from "./query-keys";

export interface CreateSignatureInput {
  editionId: string;
  data: SignatureCreateOutputI;
}

export interface DeleteSignatureInput {
  editionId: string;
  signatureId: string;
}

export interface CreateSignatureRequestInput {
  editionId: string;
  data: SignatureRequestCreateI;
}

export interface CancelSignatureRequestInput {
  editionId: string;
  requestId: string;
  reason?: string;
}

export const createSignatureEffect = (
  editionId: string,
  data: SignatureCreateOutputI,
) =>
  apiEffect(() =>
    withSpan("action:signature-create", () =>
      createSignature(
        editionId,
        data as unknown as AddSignatureRequest,
      ).then(orvalData<SignatureI>),
    ),
  );

export const deleteSignatureEffect = (signatureId: string) =>
  apiEffect(() =>
    withSpan("action:signature-delete", () =>
      deleteSignature(signatureId).then(orvalData<null>),
    ),
  );

export const createSignatureRequestEffect = (
  editionId: string,
  data: SignatureRequestCreateI,
) =>
  apiEffect(() =>
    withSpan("action:signature-request-create", () =>
      createSignatureRequest(
        editionId,
        data as unknown as CreateSignatureRequestRequest,
      ).then(orvalData<SignatureRequestI>),
    ),
  );

export const cancelSignatureRequestEffect = (
  requestId: string,
  reason?: string,
) =>
  apiEffect(() =>
    withSpan("action:signature-request-cancel", () =>
      cancelSignatureRequest(
        requestId,
        (reason ? { reason } : {}) as CancelSignatureRequestBody,
      ).then(orvalData<null>),
    ),
  );

export const useCreateSignatureMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ editionId, data }: CreateSignatureInput) =>
      createSignatureEffect(editionId, data).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: signatureKeys.byEdition(editionId),
            });
          }),
        ),
      ),
  });
};

export const useDeleteSignatureMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ editionId, signatureId }: DeleteSignatureInput) =>
      deleteSignatureEffect(signatureId).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: signatureKeys.byEdition(editionId),
            });
            void queryClient.invalidateQueries({
              queryKey: signatureKeys.byId(signatureId),
            });
          }),
        ),
      ),
  });
};

export const useCreateSignatureRequestMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ editionId, data }: CreateSignatureRequestInput) =>
      createSignatureRequestEffect(editionId, data).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: signatureKeys.requestsByEdition(editionId),
            });
          }),
        ),
      ),
  });
};

export const useCancelSignatureRequestMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ editionId, requestId, reason }: CancelSignatureRequestInput) =>
      cancelSignatureRequestEffect(requestId, reason).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: signatureKeys.requestsByEdition(editionId),
            });
            void queryClient.invalidateQueries({
              queryKey: signatureKeys.requestById(requestId),
            });
          }),
        ),
      ),
  });
};
