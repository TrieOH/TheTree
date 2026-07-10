import { useMutation, useQueryClient, type QueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import type { EventI } from "../model";
import { createEventFn, patchEventFn } from "./index";
import { eventKeys } from "./query-keys";

type PatchEventInput = {
  id: string;
  data: Partial<EventI>;
};

function upsertById(events: EventI[] | undefined, event: EventI) {
  const list = events ?? [];
  const index = list.findIndex((item) => item.id === event.id);

  if (index === -1) return [...list, event];

  const next = [...list];
  next[index] = event;
  return next;
}

function removeById(events: EventI[] | undefined, eventId: string) {
  return (events ?? []).filter((item) => item.id !== eventId);
}

function shouldBePublic(event: Pick<EventI, "status">) {
  return event.status !== "draft";
}

function syncEventCaches(queryClient: QueryClient, event: EventI) {
  queryClient.setQueryData<EventI[]>(eventKeys.ownLists(), (old) => upsertById(old, event));

  queryClient.setQueryData<EventI[]>(eventKeys.publicLists(), (old) => {
    if (shouldBePublic(event)) return upsertById(old, event);
    return removeById(old, event.id);
  });
}

export function useCreateEventMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createEventFn,
    onSuccess: (res) => {
      if (!res.success) {
        toast.error(res.message || "Erro ao criar evento");
        return;
      }

      syncEventCaches(queryClient, res.data);
      toast.success("Evento criado com sucesso!");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}

export function usePatchEventMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: PatchEventInput) => patchEventFn(id, data),
    onSuccess: (res) => {
      if (!res.success) {
        toast.error(res.message || "Erro ao atualizar evento");
        return;
      }

      syncEventCaches(queryClient, res.data);
      toast.success("Evento atualizado com sucesso!");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}
