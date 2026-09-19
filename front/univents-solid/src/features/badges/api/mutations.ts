import { orvalData } from "@trieoh/api-client";
import { useQueryClient } from "@trieoh/front-core-solid";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import {
  createBadgeTemplate,
  deleteBadgeTemplate,
  updateBadgeTemplate,
} from "@trieoh/univents-api";
import type {
  CreateBadgeTemplateRequest,
  UpdateBadgeTemplateRequest,
} from "@trieoh/univents-api/schemas";
import { Effect } from "effect";
import { apiEffect, useEffectMutation } from "@/shared/lib/effect-query";
import type { BadgeTemplate } from "../model";
import { badgeKeys } from "./query-keys";

export interface CreateBadgeTemplateInput {
  editionId: string;
  data: CreateBadgeTemplateRequest;
}

export interface UpdateBadgeTemplateInput {
  templateId: string;
  editionId: string;
  data: UpdateBadgeTemplateRequest;
}

export interface DeleteBadgeTemplateInput {
  templateId: string;
  editionId: string;
}

export const createBadgeTemplateEffect = (
  editionId: string,
  data: CreateBadgeTemplateRequest,
) =>
  apiEffect(() =>
    withSpan("action:badge-template-create", () =>
      createBadgeTemplate(editionId, {
        ...data,
        origin: data.origin ?? null,
      }).then(orvalData<BadgeTemplate>),
    ),
  );

export const updateBadgeTemplateEffect = (
  templateId: string,
  data: UpdateBadgeTemplateRequest,
) =>
  apiEffect(() =>
    withSpan("action:badge-template-update", () =>
      updateBadgeTemplate(templateId, {
        name: data.name,
        design_data: data.design_data,
      }).then(orvalData<BadgeTemplate>),
    ),
  );

export const deleteBadgeTemplateEffect = (templateId: string) =>
  apiEffect(() =>
    withSpan("action:badge-template-delete", () =>
      deleteBadgeTemplate(templateId).then(orvalData<null>),
    ),
  );

export const useCreateBadgeTemplateMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ editionId, data }: CreateBadgeTemplateInput) =>
      createBadgeTemplateEffect(editionId, data).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: badgeKeys.byEdition(editionId),
            });
          }),
        ),
      ),
  });
};

export const useUpdateBadgeTemplateMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ templateId, editionId, data }: UpdateBadgeTemplateInput) =>
      updateBadgeTemplateEffect(templateId, data).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: badgeKeys.byEdition(editionId),
            });
            void queryClient.invalidateQueries({
              queryKey: badgeKeys.detail(templateId),
            });
          }),
        ),
      ),
  });
};

export const useDeleteBadgeTemplateMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ templateId, editionId }: DeleteBadgeTemplateInput) =>
      deleteBadgeTemplateEffect(templateId).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: badgeKeys.byEdition(editionId),
            });
            void queryClient.invalidateQueries({
              queryKey: badgeKeys.detail(templateId),
            });
          }),
        ),
      ),
  });
};
