import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { getErrorMessage } from "@/shared/lib/errors";
import type { EditionCreateOutputI, EditionPatchOutputI } from "../model";
import { syncEditionCaches } from "./cache";
import { createEditionFn, patchEditionFn, publishEditionFn } from "./index";
import { editionKeys } from "./query-keys";

type CreateEditionInput = {
  eventId: string;
  data: EditionCreateOutputI;
};

type PatchEditionInput = {
  eventId: string;
  editionId: string;
  data: EditionPatchOutputI;
};

type PublishEditionInput = {
  eventId: string;
  editionId: string;
};

export function useCreateEditionMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ eventId, data }: CreateEditionInput) =>
      createEditionFn(data, eventId),
    onSuccess: (edition) => {
      syncEditionCaches(queryClient, edition);
      toast.success("Edição criada com sucesso!");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Não foi possível criar a edição")),
  });
}

export function usePatchEditionMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ eventId, editionId, data }: PatchEditionInput) =>
      patchEditionFn(eventId, editionId, data),
    onSuccess: (edition) => {
      syncEditionCaches(queryClient, edition);
      toast.success("Edição atualizada com sucesso!");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível atualizar a edição"),
      ),
  });
}

export function usePublishEditionMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ eventId, editionId }: PublishEditionInput) =>
      publishEditionFn(eventId, editionId),
    onSuccess: (_res, { eventId }) => {
      void queryClient.invalidateQueries({
        queryKey: editionKeys.adminListByEvent(eventId),
      });
      void queryClient.invalidateQueries({
        queryKey: editionKeys.publicListByEvent(eventId),
      });
      toast.success("Edição publicada com sucesso!");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Não foi possível publicar a edição")),
  });
}
