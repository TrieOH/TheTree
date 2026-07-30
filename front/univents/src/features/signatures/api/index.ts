import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import type {
  SignatureCreateOutputI,
  SignatureI,
  SignatureRequestCreateI,
  SignatureRequestI,
} from "@/features/signatures/model";
import {
  authFetcher,
  publicFetcher,
  publicQueryFetcher,
} from "@/shared/lib/api/fetch";
import { signatureKeys } from "./query-keys";

export const createSignatureFn = createClientOnlyFn(
  (editionId: string, payload: SignatureCreateOutputI) => {
    return authFetcher.post<SignatureI>(
      `/editions/${editionId}/signatures`,
      payload,
    );
  },
);

export const getAllSignaturesFn = createClientOnlyFn(
  async (editionId: string) => {
    return publicQueryFetcher<SignatureI[]>(
      `/editions/${editionId}/signatures`,
    );
  },
);

export const allSignaturesQueryOptions = (eventId: string, editionId: string) =>
  queryOptions({
    queryKey: signatureKeys.byEdition(eventId, editionId),
    queryFn: () => getAllSignaturesFn(editionId),
  });

export const getSignatureFn = createClientOnlyFn((sigId: string) => {
  return publicQueryFetcher<SignatureI>(`/signatures/${sigId}`);
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
  return authFetcher.delete<null>(`/signatures/${sigId}`);
});

export const getAllSignatureRequestsFn = createClientOnlyFn(
  (editionId: string) =>
    publicQueryFetcher<SignatureRequestI[]>(
      `/editions/${editionId}/signature-requests`,
    ),
);

export const allSignatureRequestsQueryOptions = (editionId: string) =>
  queryOptions({
    queryKey: [...signatureKeys.all, "requests", editionId],
    queryFn: () => getAllSignatureRequestsFn(editionId),
  });

export const createSignatureRequestFn = createClientOnlyFn(
  (editionId: string, data: SignatureRequestCreateI) =>
    authFetcher.post<SignatureRequestI>(
      `/editions/${editionId}/signature-requests`,
      data,
    ),
);

export const getSignatureRequestFn = createClientOnlyFn((requestId: string) =>
  publicQueryFetcher<SignatureRequestI>(`/signature-requests/${requestId}`),
);

export const signatureRequestQueryOptions = (requestId: string) =>
  queryOptions({
    queryKey: [...signatureKeys.all, "request", requestId],
    queryFn: () => getSignatureRequestFn(requestId),
  });

export const cancelSignatureRequestFn = createClientOnlyFn(
  (requestId: string, reason?: string) =>
    authFetcher.post<null>(`/signature-requests/${requestId}/cancel`, {
      reason,
    }),
);

export const fulfillSignatureRequestFn = createClientOnlyFn(
  (token: string, imageUrl: string) =>
    publicFetcher.post<SignatureI>(
      `/signature-requests/fulfill?token=${encodeURIComponent(token)}`,
      { image_url: imageUrl },
    ),
);

export const denySignatureRequestFn = createClientOnlyFn(
  (token: string, reason?: string) =>
    publicFetcher.post<null>(
      `/signature-requests/deny?token=${encodeURIComponent(token)}`,
      { reason },
    ),
);

export const revokeSignatureFn = createClientOnlyFn((token: string) =>
  publicFetcher.post<null>(
    `/signatures/revoke?token=${encodeURIComponent(token)}`,
  ),
);
