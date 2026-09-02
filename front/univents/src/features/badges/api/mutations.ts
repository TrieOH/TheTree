import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { getErrorMessage } from "@/shared/lib/errors";
import type { BadgeTemplateCreate } from "../model";
import {
  createBadgeTemplateFn,
  deleteBadgeTemplateFn,
  updateBadgeTemplateFn,
} from ".";
import { removeBadgeTemplateCache, syncBadgeTemplateCache } from "./cache";

export function useCreateBadgeTemplateMutation() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: ({
      editionId,
      data,
    }: {
      editionId: string;
      data: BadgeTemplateCreate;
    }) => createBadgeTemplateFn(editionId, data),
    onSuccess: (template) => {
      syncBadgeTemplateCache(client, template);
      toast.success("Template de crachá criado");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Não foi possível criar o template")),
  });
}

export function useUpdateBadgeTemplateMutation() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: ({
      templateId,
      data,
    }: {
      templateId: string;
      data: BadgeTemplateCreate;
    }) => updateBadgeTemplateFn(templateId, data),
    onSuccess: (template) => {
      syncBadgeTemplateCache(client, template);
      toast.success("Template de crachá atualizado");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível atualizar o template"),
      ),
  });
}

export function useDeleteBadgeTemplateMutation() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: ({ templateId }: { editionId: string; templateId: string }) =>
      deleteBadgeTemplateFn(templateId),
    onSuccess: (_, { editionId, templateId }) => {
      removeBadgeTemplateCache(client, editionId, templateId);
      toast.success("Template excluído");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível excluir o template"),
      ),
  });
}
