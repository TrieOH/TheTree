import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import type { SignatureCreateOutputI } from "@/features/signatures/model";
import { getErrorMessage } from "@/shared/lib/errors";
import type { SignatureRequestCreateOutputI } from "../model";
import { findActorIdByEmailServerFn } from "./actor-lookup";
import {
  invalidateSignatureRequestCache,
  removeSignatureCache,
  syncSignatureCache,
  syncSignatureRequestCache,
} from "./cache";
import {
  cancelSignatureRequestFn,
  createSignatureFn,
  createSignatureRequestFn,
  removeSignatureFn,
} from "./index";

type CreateSignatureInput = {
  editionId: string;
  data: SignatureCreateOutputI;
};

type RemoveSignatureInput = {
  editionId: string;
  signatureId: string;
};

export function useCreateSignatureMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ editionId, data }: CreateSignatureInput) =>
      createSignatureFn(editionId, data),
    onSuccess: (signature) => {
      syncSignatureCache(queryClient, signature);
      toast.success("Assinatura criada com sucesso");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível criar a assinatura"),
      ),
  });
}

export function useRemoveSignatureMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ signatureId }: RemoveSignatureInput) =>
      removeSignatureFn(signatureId),
    onSuccess: (_, { editionId, signatureId }) => {
      removeSignatureCache(queryClient, editionId, signatureId);
      toast.success("Assinatura removida");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível remover a assinatura"),
      ),
  });
}

export function useCreateSignatureRequestMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      editionId,
      data,
    }: {
      editionId: string;
      data: SignatureRequestCreateOutputI;
    }) =>
      createSignatureRequestFn(editionId, {
        ...data,
        signatory_title: data.signatory_title || undefined,
        signatory_user_id: await findActorIdByEmailServerFn({
          data: { email: data.signatory_email },
        }),
      }),
    onSuccess: (request) => {
      syncSignatureRequestCache(queryClient, request);
      toast.success("Convite enviado");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Não foi possível enviar o convite")),
  });
}

export function useCancelSignatureRequestMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ requestId }: { editionId: string; requestId: string }) =>
      cancelSignatureRequestFn(requestId),
    onSuccess: (_, { editionId, requestId }) => {
      invalidateSignatureRequestCache(queryClient, editionId, requestId);
      toast.success("Convite cancelado");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível cancelar o convite"),
      ),
  });
}
