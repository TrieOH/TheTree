import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { getErrorMessage } from "@/shared/lib/errors";
import type { CertificationTemplateCreateI } from "../model";
import {
  removeCertificationTemplateCache,
  syncCertificationTemplateCache,
} from "./cache";
import {
  createCertificationTemplateFn,
  deleteCertificationTemplateFn,
  emitProgramCertificationsFn,
  invalidateCertificationFn,
  linkCertificationTemplateFn,
  unlinkCertificationTemplateFn,
  updateCertificationTemplateFn,
} from "./index";
import { certificationKeys } from "./query-keys";

interface CreateTemplateInput {
  editionId: string;
  data: CertificationTemplateCreateI;
}

export function useCreateCertificationTemplateMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ editionId, data }: CreateTemplateInput) =>
      createCertificationTemplateFn(editionId, data),
    onSuccess: (template) => {
      syncCertificationTemplateCache(queryClient, template);
      toast.success("Template de certificado criado");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Não foi possível criar o template")),
  });
}

export function useLinkCertificationTemplateMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      templateId,
      programId,
    }: {
      templateId: string;
      programId: string;
    }) => linkCertificationTemplateFn(templateId, programId),
    onSuccess: (_, variables) => {
      void queryClient.invalidateQueries({
        queryKey: certificationKeys.templateLinks(variables.templateId),
      });
      toast.success("Template vinculado à atividade");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível vincular o template"),
      ),
  });
}

export function useUpdateCertificationTemplateMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      templateId,
      data,
    }: {
      templateId: string;
      data: CertificationTemplateCreateI;
    }) => updateCertificationTemplateFn(templateId, data),
    onSuccess: (template) => {
      syncCertificationTemplateCache(queryClient, template);
      toast.success("Template atualizado");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível atualizar o template"),
      ),
  });
}

export function useDeleteCertificationTemplateMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ templateId }: { editionId: string; templateId: string }) =>
      deleteCertificationTemplateFn(templateId),
    onSuccess: (_, { editionId, templateId }) => {
      removeCertificationTemplateCache(queryClient, editionId, templateId);
      toast.success("Template excluído");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível excluir o template"),
      ),
  });
}

export function useInvalidateCertificationMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      certificationId,
      reason,
    }: {
      certificationId: string;
      reason: string;
      verificationHash: string;
    }) => invalidateCertificationFn({ certificationId, reason }),
    onSuccess: (_, { verificationHash }) => {
      toast.success("Certificado invalidado");
      void queryClient.invalidateQueries({
        queryKey: certificationKeys.issued(),
      });
      void queryClient.invalidateQueries({
        queryKey: certificationKeys.verification(verificationHash),
      });
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível invalidar o certificado"),
      ),
  });
}

export function useUnlinkCertificationTemplateMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      templateId,
      programId,
    }: {
      templateId: string;
      programId: string;
    }) => unlinkCertificationTemplateFn(templateId, programId),
    onSuccess: (_, variables) => {
      void queryClient.invalidateQueries({
        queryKey: certificationKeys.templateLinks(variables.templateId),
      });
      toast.success("Vínculos removidos");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível remover os vínculos"),
      ),
  });
}

export function useEmitProgramCertificationsMutation() {
  return useMutation({
    mutationFn: emitProgramCertificationsFn,
    onSuccess: () => toast.success("Emissão de certificados iniciada"),
    onError: (error) =>
      toast.error(getErrorMessage(error, "Não foi possível iniciar a emissão")),
  });
}
