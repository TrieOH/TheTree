import type { QueryClient } from "@tanstack/react-query";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { getErrorMessage } from "@/shared/lib/errors";
import type { EventI } from "../model";
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

export function syncEventCaches(queryClient: QueryClient, event: EventI) {
  queryClient.setQueryData<EventI[]>(eventKeys.ownLists(), (old) =>
    upsertById(old, event),
  );

  queryClient.setQueryData<EventI[]>(eventKeys.publicLists(), (old) => {
    if (shouldBePublic(event)) return upsertById(old, event);
    return removeById(old, event.id);
  });

  queryClient.setQueryData<EventI[]>(eventKeys.joinedLists(), (old) =>
    old?.some((item) => item.id === event.id) ? upsertById(old, event) : old,
  );
}

function syncEventStatusInCaches(
  queryClient: QueryClient,
  eventId: string,
  status: EventI["status"],
  updatedAt: string,
) {
  const ownEvent = queryClient
    .getQueryData<EventI[]>(eventKeys.ownLists())
    ?.find((event) => event.id === eventId);
  const joinedEvent = queryClient
    .getQueryData<EventI[]>(eventKeys.joinedLists())
    ?.find((event) => event.id === eventId);
  const publicEvent = queryClient
    .getQueryData<EventI[]>(eventKeys.publicLists())
    ?.find((event) => event.id === eventId);
  const cachedEvent = ownEvent ?? joinedEvent ?? publicEvent;

  if (!cachedEvent) return;

  const nextEvent = { ...cachedEvent, status, updated_at: updatedAt };

  queryClient.setQueryData<EventI[]>(eventKeys.ownLists(), (old) =>
    old?.some((event) => event.id === eventId)
      ? upsertById(old, nextEvent)
      : old,
  );

  queryClient.setQueryData<EventI[]>(eventKeys.publicLists(), (old) =>
    upsertById(old, nextEvent),
  );

  queryClient.setQueryData<EventI[]>(eventKeys.joinedLists(), (old) =>
    old?.some((event) => event.id === eventId)
      ? upsertById(old, nextEvent)
      : old,
  );
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
    onSuccess: (_, eventId) => {
      syncEventStatusInCaches(
        queryClient,
        eventId,
        "active",
        new Date().toISOString(),
      );
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
    onSuccess: (_, eventId) => {
      syncEventStatusInCaches(
        queryClient,
        eventId,
        "discontinued",
        new Date().toISOString(),
      );
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
        (old = []) => [...old.filter((member) => member.id !== res.id), res],
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
        (old = []) => old.filter((member) => member.user_id !== input.userId),
      );
      toast.success("Membro removido com sucesso!");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Erro ao remover membro")),
  });
}
