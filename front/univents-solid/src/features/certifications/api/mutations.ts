import { orvalData } from "@trieoh/api-client";
import { useQueryClient } from "@trieoh/front-core-solid";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import {
  createCertificationTemplate,
  deleteCertificationTemplate,
  emitProgramCertifications,
  invalidateCertification,
  linkCertificationTemplate,
  unlinkCertificationTemplate,
  updateCertificationTemplate,
} from "@trieoh/univents-api";
import type {
  CreateCertificationTemplateRequest,
  UpdateCertificationTemplateRequest,
} from "@trieoh/univents-api/schemas";
import { Effect } from "effect";
import { apiEffect, useEffectMutation } from "@/shared/lib/effect-query";
import type {
  CertificationI,
  CertificationTemplateCreateI,
  CertificationTemplateI,
} from "../model";
import { certificationKeys } from "./query-keys";

export interface CreateCertificationTemplateInput {
  editionId: string;
  data: CertificationTemplateCreateI;
}

export interface UpdateCertificationTemplateInput {
  templateId: string;
  editionId: string;
  data: CertificationTemplateCreateI;
}

export interface DeleteCertificationTemplateInput {
  templateId: string;
  editionId: string;
}

export interface InvalidateCertificationInput {
  editionId: string;
  certificationId: string;
  reason: string;
}

export interface EmitProgramCertificationsInput {
  editionId: string;
  programId: string;
}

export const createCertificationTemplateEffect = (
  editionId: string,
  data: CertificationTemplateCreateI,
) =>
  apiEffect(() =>
    withSpan("action:certification-template-create", () =>
      createCertificationTemplate(
        editionId,
        data as unknown as CreateCertificationTemplateRequest,
      ).then(orvalData<CertificationTemplateI>),
    ),
  );

export const updateCertificationTemplateEffect = (
  templateId: string,
  data: CertificationTemplateCreateI,
) =>
  apiEffect(() =>
    withSpan("action:certification-template-update", () =>
      updateCertificationTemplate(
        templateId,
        data as unknown as UpdateCertificationTemplateRequest,
      ).then(orvalData<CertificationTemplateI>),
    ),
  );

export const deleteCertificationTemplateEffect = (templateId: string) =>
  apiEffect(() =>
    withSpan("action:certification-template-delete", () =>
      deleteCertificationTemplate(templateId).then(orvalData<null>),
    ),
  );

export const invalidateCertificationEffect = (
  certificationId: string,
  reason: string,
) =>
  apiEffect(() =>
    withSpan("action:certification-invalidate", () =>
      invalidateCertification(certificationId, { reason }).then(
        orvalData<CertificationI>,
      ),
    ),
  );

export const emitProgramCertificationsEffect = (programId: string) =>
  apiEffect(() =>
    withSpan("action:certification-emit-program", () =>
      emitProgramCertifications(programId).then(orvalData<null>),
    ),
  );

export const useCreateCertificationTemplateMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ editionId, data }: CreateCertificationTemplateInput) =>
      createCertificationTemplateEffect(editionId, data).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: certificationKeys.templatesByEdition(editionId),
            });
          }),
        ),
      ),
  });
};

export const useUpdateCertificationTemplateMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      templateId,
      editionId,
      data,
    }: UpdateCertificationTemplateInput) =>
      updateCertificationTemplateEffect(templateId, data).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: certificationKeys.templateById(templateId),
            });
            void queryClient.invalidateQueries({
              queryKey: certificationKeys.templatesByEdition(editionId),
            });
          }),
        ),
      ),
  });
};

export const useDeleteCertificationTemplateMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      templateId,
      editionId,
    }: DeleteCertificationTemplateInput) =>
      deleteCertificationTemplateEffect(templateId).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: certificationKeys.templatesByEdition(editionId),
            });
            void queryClient.invalidateQueries({
              queryKey: certificationKeys.templateById(templateId),
            });
          }),
        ),
      ),
  });
};

export const useInvalidateCertificationMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      editionId,
      certificationId,
      reason,
    }: InvalidateCertificationInput) =>
      invalidateCertificationEffect(certificationId, reason).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: certificationKeys.issuedByEdition(editionId),
            });
            void queryClient.invalidateQueries({
              queryKey: certificationKeys.issuedById(certificationId),
            });
          }),
        ),
      ),
  });
};

export const useEmitProgramCertificationsMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      editionId,
      programId,
    }: EmitProgramCertificationsInput) =>
      emitProgramCertificationsEffect(programId).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: certificationKeys.issuedByEdition(editionId),
            });
            void queryClient.invalidateQueries({
              queryKey: certificationKeys.emissionErrorsByEdition(editionId),
            });
          }),
        ),
      ),
  });
};

export const linkCertificationTemplateEffect = (
  templateId: string,
  programId: string,
) =>
  apiEffect(() =>
    withSpan("action:certification-template-link", () =>
      linkCertificationTemplate(templateId, {
        program_id: programId,
      }).then(orvalData<null>),
    ),
  );

export const unlinkCertificationTemplateEffect = (
  templateId: string,
  programId: string,
) =>
  apiEffect(() =>
    withSpan("action:certification-template-unlink", () =>
      unlinkCertificationTemplate(templateId, {
        program_id: programId,
      }).then(orvalData<null>),
    ),
  );

export const useLinkCertificationTemplateMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      templateId,
      programId,
    }: {
      templateId: string;
      programId: string;
    }) =>
      linkCertificationTemplateEffect(templateId, programId).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: certificationKeys.templateLinks(templateId),
            });
            void queryClient.invalidateQueries({
              queryKey: certificationKeys.allEditionProgramLinks(),
            });
          }),
        ),
      ),
  });
};

export const useUnlinkCertificationTemplateMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({
      templateId,
      programId,
    }: {
      templateId: string;
      programId: string;
    }) =>
      unlinkCertificationTemplateEffect(templateId, programId).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: certificationKeys.templateLinks(templateId),
            });
            void queryClient.invalidateQueries({
              queryKey: certificationKeys.allEditionProgramLinks(),
            });
          }),
        ),
      ),
  });
};
