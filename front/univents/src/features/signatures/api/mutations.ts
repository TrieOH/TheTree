import {
  type QueryClient,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { toast } from "sonner";
import type {
  SignatureCreateOutputI,
  SignatureI,
} from "@/features/signatures/model";
import { createSignatureFn, removeSignatureFn } from "./index";
import { signatureKeys } from "./query-keys";

type CreateSignatureInput = {
  eventId: string;
  editionId: string;
  data: SignatureCreateOutputI;
};

type RemoveSignatureInput = {
  eventId: string;
  editionId: string;
  signatureId: string;
};

function upsertById(
  signatures: SignatureI[] | undefined,
  signature: SignatureI,
) {
  const list = signatures ?? [];
  const index = list.findIndex((item) => item.id === signature.id);

  if (index === -1) return [...list, signature];

  const next = [...list];
  next[index] = signature;
  return next;
}

function removeById(signatures: SignatureI[] | undefined, signatureId: string) {
  return (signatures ?? []).filter((signature) => signature.id !== signatureId);
}

function syncSignatureCaches(
  queryClient: QueryClient,
  eventId: string,
  editionId: string,
  signature: SignatureI,
) {
  queryClient.setQueryData<SignatureI[]>(
    signatureKeys.byEdition(eventId, editionId),
    (old) => upsertById(old, signature),
  );
}

function syncSignatureRemoval(
  queryClient: QueryClient,
  eventId: string,
  editionId: string,
  signatureId: string,
) {
  queryClient.setQueryData<SignatureI[]>(
    signatureKeys.byEdition(eventId, editionId),
    (old) => removeById(old, signatureId),
  );
}

export function useCreateSignatureMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ editionId, data }: CreateSignatureInput) =>
      createSignatureFn(editionId, data),
    onSuccess: (signature, variables) => {
      syncSignatureCaches(
        queryClient,
        variables.eventId,
        variables.editionId,
        signature,
      );
      toast.success("Assinatura criada com sucesso");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}

export function useRemoveSignatureMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ signatureId }: RemoveSignatureInput) =>
      removeSignatureFn(signatureId),
    onSuccess: (_res, variables) => {
      syncSignatureRemoval(
        queryClient,
        variables.eventId,
        variables.editionId,
        variables.signatureId,
      );
      toast.success("Assinatura removida");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}
