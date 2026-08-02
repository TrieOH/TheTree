import type { QueryClient } from "@tanstack/react-query";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import type {
  EditionCreateOutputI,
  EditionI,
  EditionPatchOutputI,
} from "../model";
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

function upsertById(editions: EditionI[] | undefined, edition: EditionI) {
  const list = editions ?? [];
  const index = list.findIndex((item) => item.id === edition.id);

  if (index === -1) return [...list, edition];

  const next = [...list];
  next[index] = edition;
  return next;
}

function syncEditionCaches(queryClient: QueryClient, edition: EditionI) {
  queryClient.setQueryData<EditionI[]>(
    editionKeys.adminListByEvent(edition.event_id),
    (old) => upsertById(old, edition),
  );

  queryClient.setQueryData<EditionI[]>(
    editionKeys.publicListByEvent(edition.event_id),
    (old) =>
      edition.is_draft
        ? (old ?? []).filter((item) => item.id !== edition.id)
        : upsertById(old, edition),
  );
}

export function useCreateEditionMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ eventId, data }: CreateEditionInput) =>
      createEditionFn(data, eventId),
    onSuccess: (edition) => {
      syncEditionCaches(queryClient, edition);
      toast.success("Edição criada com sucesso!");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
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
    onError: () => toast.error("Erro ao conectar com o servidor"),
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
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}
