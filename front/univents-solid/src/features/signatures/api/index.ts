import { orvalData } from "@trieoh/api-client";
import {
  getSignature,
  getSignatureRequest,
  listEditionSignatureRequests,
  listEditionSignatures,
} from "@trieoh/univents-api";
import type { SignatureI, SignatureRequestI } from "../model";
import { signatureKeys } from "./query-keys";

export const allSignaturesQueryOptions = (editionId: string) => ({
  queryKey: signatureKeys.byEdition(editionId),
  queryFn: async () => {
    return listEditionSignatures(editionId, { public: true }).then(
      orvalData<SignatureI[]>,
    );
  },
  enabled: Boolean(editionId),
});

export const signatureQueryOptions = (sigId: string) => ({
  queryKey: signatureKeys.byId(sigId),
  queryFn: async () => {
    return getSignature(sigId, { public: true }).then(orvalData<SignatureI>);
  },
  enabled: Boolean(sigId),
});

export const allSignatureRequestsQueryOptions = (editionId: string) => ({
  queryKey: signatureKeys.requestsByEdition(editionId),
  queryFn: async () => {
    return listEditionSignatureRequests(editionId, { public: true }).then(
      orvalData<SignatureRequestI[]>,
    );
  },
  enabled: Boolean(editionId),
});

export const signatureRequestQueryOptions = (requestId: string) => ({
  queryKey: signatureKeys.requestById(requestId),
  queryFn: async () => {
    return getSignatureRequest(requestId, { public: true }).then(
      orvalData<SignatureRequestI>,
    );
  },
  enabled: Boolean(requestId),
});

export * from "./query-keys";
