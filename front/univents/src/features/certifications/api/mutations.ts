import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import type { CertificationTemplateCreateI } from "../model";
import {
  createCertificationTemplateFn,
  deleteCertificationTemplateFn,
  invalidateCertificationFn,
  linkCertificationTemplateFn,
  unlinkCertificationTemplateFn,
  updateCertificationTemplateFn,
} from "./index";
import { certificationKeys } from "./query-keys";

interface EditionInput {
  eventId: string;
  editionId: string;
}

interface CreateTemplateInput extends EditionInput {
  data: CertificationTemplateCreateI;
}

export function useCreateCertificationTemplateMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ editionId, data }: CreateTemplateInput) =>
      createCertificationTemplateFn(editionId, data),
    onSuccess: (_template, variables) => {
      void queryClient.invalidateQueries({
        queryKey: certificationKeys.templatesByEdition(variables.editionId),
      });
      toast.success("Template de certificado criado");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
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
        queryKey: [
          ...certificationKeys.templates(),
          "links",
          variables.templateId,
        ],
      });
      toast.success("Template vinculado à atividade");
    },
    onError: () => toast.error("Não foi possível vincular o template"),
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
    onSuccess: (_, variables) => {
      void queryClient.invalidateQueries({
        queryKey: certificationKeys.templates(),
      });
      void queryClient.invalidateQueries({
        queryKey: [
          ...certificationKeys.templates(),
          "detail",
          variables.templateId,
        ],
      });
      toast.success("Template atualizado");
    },
    onError: () => toast.error("Não foi possível atualizar o template"),
  });
}

export function useDeleteCertificationTemplateMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ templateId }: { templateId: string }) =>
      deleteCertificationTemplateFn(templateId),
    onSuccess: (_, variables) => {
      void queryClient.invalidateQueries({
        queryKey: certificationKeys.templates(),
      });
      void queryClient.removeQueries({
        queryKey: [
          ...certificationKeys.templates(),
          "detail",
          variables.templateId,
        ],
      });
      toast.success("Template excluído");
    },
    onError: () => toast.error("Não foi possível excluir o template"),
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
    }) => invalidateCertificationFn({ certificationId, reason }),
    onSuccess: () => {
      toast.success("Certificado invalidado");
      void queryClient.invalidateQueries({
        queryKey: certificationKeys.issued(),
      });
    },
    onError: () => toast.error("Não foi possível invalidar o certificado"),
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
        queryKey: [
          ...certificationKeys.templates(),
          "links",
          variables.templateId,
        ],
      });
      toast.success("Vínculos removidos");
    },
    onError: () => toast.error("Não foi possível remover os vínculos"),
  });
}
