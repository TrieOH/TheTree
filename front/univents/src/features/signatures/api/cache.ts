import type { QueryClient } from "@tanstack/react-query";
import type { SignatureI, SignatureRequestI } from "../model";
import { signatureKeys } from "./query-keys";

function upsertById<T extends { id: string }>(list: T[] | undefined, item: T) {
  if (!list) return list;
  const index = list.findIndex((candidate) => candidate.id === item.id);
  if (index === -1) return [...list, item];
  const next = [...list];
  next[index] = item;
  return next;
}

export function syncSignatureCache(
  queryClient: QueryClient,
  signature: SignatureI,
) {
  queryClient.setQueryData<SignatureI[]>(
    signatureKeys.byEdition(signature.edition_id),
    (old) => upsertById(old, signature),
  );
  queryClient.setQueryData<SignatureI>(
    signatureKeys.byId(signature.id),
    (old) => (old ? signature : old),
  );
}

export function removeSignatureCache(
  queryClient: QueryClient,
  editionId: string,
  signatureId: string,
) {
  queryClient.setQueryData<SignatureI[]>(
    signatureKeys.byEdition(editionId),
    (old) => old?.filter((signature) => signature.id !== signatureId),
  );
  queryClient.removeQueries({ queryKey: signatureKeys.byId(signatureId) });
}

export function syncSignatureRequestCache(
  queryClient: QueryClient,
  request: SignatureRequestI,
) {
  queryClient.setQueryData<SignatureRequestI[]>(
    signatureKeys.requestsByEdition(request.edition_id),
    (old) => upsertById(old, request),
  );
  queryClient.setQueryData<SignatureRequestI>(
    signatureKeys.requestById(request.id),
    (old) => (old ? request : old),
  );
}

export function invalidateSignatureRequestCache(
  queryClient: QueryClient,
  editionId: string,
  requestId: string,
) {
  void queryClient.invalidateQueries({
    queryKey: signatureKeys.requestsByEdition(editionId),
  });
  void queryClient.invalidateQueries({
    queryKey: signatureKeys.requestById(requestId),
  });
}
