import type { QueryClient } from "@tanstack/react-query";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { getErrorMessage } from "@/shared/lib/errors";
import { syncEventCaches } from "./cache";
import {
  createEventFn,
  discontinueEventFn,
  patchEventFn,
  publishEventFn,
} from "./index";
import {
  addEventMemberFn,
  type EventMemberWithEmailI,
  removeEventMemberFn,
} from "./members";
import { eventKeys } from "./query-keys";

function invalidateEventReads(queryClient: QueryClient) {
  void queryClient.invalidateQueries({ queryKey: eventKeys.lists() });
  void queryClient.invalidateQueries({ queryKey: eventKeys.details() });
}

export function useCreateEventMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createEventFn,
    onSuccess: (res) => {
      syncEventCaches(queryClient, res);
      toast.success("Evento criado com sucesso!");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Erro ao criar evento")),
  });
}

export function usePublishEventMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (eventId: string) => publishEventFn(eventId),
    onSuccess: () => {
      invalidateEventReads(queryClient);
      toast.success("Evento publicado com sucesso!");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Erro ao publicar evento")),
  });
}

export function usePatchEventMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      eventId,
      data,
    }: {
      eventId: string;
      data: Parameters<typeof patchEventFn>[1];
    }) => patchEventFn(eventId, data),
    onSuccess: (res) => {
      syncEventCaches(queryClient, res);
      void queryClient.invalidateQueries({ queryKey: eventKeys.details() });
      toast.success("Evento atualizado com sucesso!");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Erro ao editar evento")),
  });
}

export function useDiscontinueEventMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (eventId: string) => discontinueEventFn(eventId),
    onSuccess: () => {
      invalidateEventReads(queryClient);
      toast.success("Evento descontinuado com sucesso!");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Erro ao descontinuar evento")),
  });
}

export function useAddEventMemberMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: addEventMemberFn,
    onSuccess: (res, input) => {
      queryClient.setQueryData<EventMemberWithEmailI[]>(
        eventKeys.members(input.eventId),
        (old) =>
          old ? [...old.filter((member) => member.id !== res.id), res] : old,
      );
      toast.success("Membro adicionado com sucesso!");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Erro ao adicionar membro")),
  });
}

export function useRemoveEventMemberMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: removeEventMemberFn,
    onSuccess: (_, input) => {
      queryClient.setQueryData<EventMemberWithEmailI[]>(
        eventKeys.members(input.eventId),
        (old) =>
          old ? old.filter((member) => member.user_id !== input.userId) : old,
      );
      toast.success("Membro removido com sucesso!");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Erro ao remover membro")),
  });
}
