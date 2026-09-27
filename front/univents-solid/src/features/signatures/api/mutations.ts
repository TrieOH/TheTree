import { orvalData } from "@trieoh/api-client";
import { useQueryClient } from "@trieoh/front-core-solid";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import {
  cancelSignatureRequest,
  createSignature,
  createSignatureRequest,
  deleteSignature,
  denySignatureRequest,
  fulfillSignatureRequest,
  revokeSignature,
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

export interface FulfillSignatureRequestInput {
  token: string;
  imageUrl: string;
  requestId?: string;
  editionId?: string;
}

export interface DenySignatureRequestInput {
  token: string;
  reason?: string;
  requestId?: string;
  editionId?: string;
}

export interface RevokeSignatureInput {
  token: string;
  signatureId?: string;
  editionId?: string;
}

export const fulfillSignatureRequestEffect = (
  token: string,
  imageUrl: string,
) =>
  apiEffect(() =>
    withSpan("action:signature-request-fulfill", () =>
      fulfillSignatureRequest(
        { image_url: imageUrl },
        { token },
      ).then(orvalData<SignatureI>),
    ),
  );

export const denySignatureRequestEffect = (
  token: string,
  reason?: string,
) =>
  apiEffect(() =>
    withSpan("action:signature-request-deny", () =>
      denySignatureRequest(
        reason ? { reason } : {},
        { token },
      ).then(() => null),
    ),
  );

export const revokeSignatureEffect = (token: string) =>
  apiEffect(() =>
    withSpan("action:signature-revoke", () =>
      revokeSignature({ token }).then(() => null),
    ),
  );

export const useFulfillSignatureRequestMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      token,
      imageUrl,
      requestId,
      editionId,
    }: FulfillSignatureRequestInput) =>
      fulfillSignatureRequestEffect(token, imageUrl).pipe(
        Effect.tap((signature) =>
          Effect.sync(() => {
            if (requestId) {
              void queryClient.invalidateQueries({
                queryKey: signatureKeys.requestById(requestId),
              });
            }
            const targetEditionId = editionId ?? signature?.edition_id;
            if (targetEditionId) {
              void queryClient.invalidateQueries({
                queryKey: signatureKeys.requestsByEdition(targetEditionId),
              });
              void queryClient.invalidateQueries({
                queryKey: signatureKeys.byEdition(targetEditionId),
              });
            }
          }),
        ),
      ),
  });
};

export const useDenySignatureRequestMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      token,
      reason,
      requestId,
      editionId,
    }: DenySignatureRequestInput) =>
      denySignatureRequestEffect(token, reason).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            if (requestId) {
              void queryClient.invalidateQueries({
                queryKey: signatureKeys.requestById(requestId),
              });
            }
            if (editionId) {
              void queryClient.invalidateQueries({
                queryKey: signatureKeys.requestsByEdition(editionId),
              });
            }
          }),
        ),
      ),
  });
};

export const useRevokeSignatureMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      token,
      signatureId,
      editionId,
    }: RevokeSignatureInput) =>
      revokeSignatureEffect(token).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            if (signatureId) {
              void queryClient.invalidateQueries({
                queryKey: signatureKeys.byId(signatureId),
              });
            }
            if (editionId) {
              void queryClient.invalidateQueries({
                queryKey: signatureKeys.byEdition(editionId),
              });
            }
          }),
        ),
      ),
  });
};
