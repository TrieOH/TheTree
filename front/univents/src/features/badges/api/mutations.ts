import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { getErrorMessage } from "@/shared/lib/errors";
import type { BadgeTemplateCreate } from "../model";
import { createBadgeTemplateFn, deleteBadgeTemplateFn } from ".";
import { badgeKeys } from "./query-keys";

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
    onSuccess: (_, { editionId }) => {
      void client.invalidateQueries({
        queryKey: badgeKeys.byEdition(editionId),
      });
      toast.success("Template de crachá criado");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Não foi possível criar o template")),
  });
}

export function useDeleteBadgeTemplateMutation() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: ({ templateId }: { templateId: string }) =>
      deleteBadgeTemplateFn(templateId),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: badgeKeys.all });
      toast.success("Template excluído");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível excluir o template"),
      ),
  });
}
