import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import {
  cancelSignatureRequest,
  createSignature,
  createSignatureRequest,
  deleteSignature,
  denySignatureRequest,
  fulfillSignatureRequest,
  getSignature,
  getSignatureRequest,
  listEditionSignatureRequests,
  listEditionSignatures,
  revokeSignature,
} from "@trieoh/univents-api";
import type {
  AddSignatureRequest,
  CancelSignatureRequestBody,
  DenySignatureRequestBody,
  FulfillSignatureRequestBody,
  RevokeSignatureParams,
} from "@trieoh/univents-api/schemas";
import type {
  SignatureCreateOutputI,
  SignatureI,
  SignatureRequestCreateI,
  SignatureRequestI,
} from "@/features/signatures/model";
import { signatureKeys } from "./query-keys";

export const createSignatureFn = createClientOnlyFn(
  (editionId: string, payload: SignatureCreateOutputI) => {
    return createSignature(editionId, payload as AddSignatureRequest).then(
      orvalData<SignatureI>,
    );
  },
);

export const getAllSignaturesFn = createClientOnlyFn(
  async (editionId: string) => {
    return listEditionSignatures(editionId, { public: true }).then(
      orvalData<SignatureI[]>,
    );
  },
);

export const allSignaturesQueryOptions = (eventId: string, editionId: string) =>
  queryOptions({
    queryKey: signatureKeys.byEdition(eventId, editionId),
    queryFn: () => getAllSignaturesFn(editionId),
  });

export const getSignatureFn = createClientOnlyFn((sigId: string) => {
  return getSignature(sigId, { public: true }).then(orvalData<SignatureI>);
});

export const signatureQueryOptions = (
  eventId: string,
  editionId: string,
  sigId: string,
) =>
  queryOptions({
    queryKey: signatureKeys.byId(eventId, editionId, sigId),
    queryFn: () => getSignatureFn(sigId),
  });

export const removeSignatureFn = createClientOnlyFn((sigId: string) => {
  return deleteSignature(sigId).then(orvalData<null>);
});

export const getAllSignatureRequestsFn = createClientOnlyFn(
  (editionId: string) =>
    listEditionSignatureRequests(editionId, { public: true }).then(
      orvalData<SignatureRequestI[]>,
    ),
);

export const allSignatureRequestsQueryOptions = (editionId: string) =>
  queryOptions({
    queryKey: [...signatureKeys.all, "requests", editionId],
    queryFn: () => getAllSignatureRequestsFn(editionId),
  });

export const createSignatureRequestFn = createClientOnlyFn(
  (editionId: string, data: SignatureRequestCreateI) =>
    createSignatureRequest(editionId, data).then(orvalData<SignatureRequestI>),
);

export const getSignatureRequestFn = createClientOnlyFn((requestId: string) =>
  getSignatureRequest(requestId, { public: true }).then(
    orvalData<SignatureRequestI>,
  ),
);

export const signatureRequestQueryOptions = (requestId: string) =>
  queryOptions({
    queryKey: [...signatureKeys.all, "request", requestId],
    queryFn: () => getSignatureRequestFn(requestId),
  });

export const cancelSignatureRequestFn = createClientOnlyFn(
  (requestId: string, reason?: string) =>
    cancelSignatureRequest(requestId, {
      reason: reason ?? "",
    } satisfies CancelSignatureRequestBody).then(orvalData<null>),
);

export const fulfillSignatureRequestFn = createClientOnlyFn(
  (token: string, imageUrl: string) =>
    fulfillSignatureRequest(
      { image_url: imageUrl } satisfies FulfillSignatureRequestBody,
      { token },
      { public: true },
    ).then(orvalData<SignatureI>),
);

export const denySignatureRequestFn = createClientOnlyFn(
  (token: string, reason?: string) =>
    denySignatureRequest(
      { reason: reason ?? "" } satisfies DenySignatureRequestBody,
      { token },
      { public: true },
    ).then(orvalData<null>),
);

export const revokeSignatureFn = createClientOnlyFn((token: string) =>
  revokeSignature({ token } satisfies RevokeSignatureParams, {
    public: true,
  }).then(orvalData<null>),
);
