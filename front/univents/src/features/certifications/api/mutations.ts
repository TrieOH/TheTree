import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import type {
  CertificationTargetType,
  CertificationTemplateCreateI,
  CertifyUserRequestI,
  SetCertificationTemplateRequestI,
} from "../model";
import {
  certifyUserFn,
  createCertificationTemplateFn,
  setActivityCertificationTemplateFn,
  setEditionCertificationTemplateFn,
} from "./index";
import { certificationKeys } from "./query-keys";

interface EditionInput {
  eventId: string;
  editionId: string;
}

interface CreateTemplateInput extends EditionInput {
  data: CertificationTemplateCreateI;
}

interface SetEditionTemplateInput extends EditionInput {
  data: SetCertificationTemplateRequestI;
}

interface SetActivityTemplateInput extends SetEditionTemplateInput {
  activityId: string;
}

interface CertifyUserInput extends EditionInput {
  userId: string;
  data: CertifyUserRequestI;
}

export function useCreateCertificationTemplateMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ eventId, editionId, data }: CreateTemplateInput) =>
      createCertificationTemplateFn(eventId, editionId, data),
    onSuccess: (response, variables) => {
      if (!response.success) {
        toast.error(response.message || "Não foi possível salvar o template");
        return;
      }
      void queryClient.invalidateQueries({
        queryKey: certificationKeys.templatesByEdition(
          variables.eventId,
          variables.editionId,
        ),
      });
      toast.success("Template de certificado criado");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}

export function useSetEditionCertificationTemplateMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ eventId, editionId, data }: SetEditionTemplateInput) =>
      setEditionCertificationTemplateFn(eventId, editionId, data),
    onSuccess: (response, variables) => {
      if (!response.success) {
        toast.error(response.message || "Erro ao definir template da edição");
        return;
      }
      void queryClient.invalidateQueries({
        queryKey: certificationKeys.templatesByEdition(
          variables.eventId,
          variables.editionId,
        ),
      });
      toast.success("Template definido para a edição");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}

export function useSetActivityCertificationTemplateMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      eventId,
      editionId,
      activityId,
      data,
    }: SetActivityTemplateInput) =>
      setActivityCertificationTemplateFn(eventId, editionId, activityId, data),
    onSuccess: (response, variables) => {
      if (!response.success) {
        toast.error(
          response.message || "Erro ao definir template da atividade",
        );
        return;
      }
      void queryClient.invalidateQueries({
        queryKey: certificationKeys.templatesByEdition(
          variables.eventId,
          variables.editionId,
        ),
      });
      toast.success("Template definido para a atividade");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}

export function useCertifyUserMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ eventId, editionId, userId, data }: CertifyUserInput) =>
      certifyUserFn(eventId, editionId, userId, data),
    onSuccess: (response, variables) => {
      if (!response.success) {
        toast.error(response.message || "Erro ao emitir certificado");
        return;
      }
      const targetType = variables.data.target_type as CertificationTargetType;
      void Promise.all([
        queryClient.invalidateQueries({
          queryKey: certificationKeys.issuedByUser(variables.userId),
        }),
        queryClient.invalidateQueries({
          queryKey: certificationKeys.issuedByTarget(
            targetType,
            variables.data.target_id,
          ),
        }),
      ]);
      toast.success("Certificado emitido com sucesso");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}
