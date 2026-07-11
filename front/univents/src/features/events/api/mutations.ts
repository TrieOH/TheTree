import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { QueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import type { EventCreateOutputI, EventI } from "../model";
import {
  addImageToTheEventGalleryFn,
  createEventFn,
  patchEventFn,
  publishEventFn,
  removeImageFromTheEventGalleryFn,
  setEventBannerFn,
  setEventLogoFn,
  unsetEventBannerFn,
  unsetEventLogoFn,
} from "./index";
import { eventKeys } from "./query-keys";
import type { ImageFieldChange } from "@/widgets/multi-step-form/model/types";

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

function syncEventStatusInCaches(queryClient: QueryClient, eventId: string, status: EventI["status"]) {
  const ownEvent = queryClient.getQueryData<EventI[]>(eventKeys.ownLists())?.find((event) => event.id === eventId);

  if (!ownEvent) return;

  const nextEvent = { ...ownEvent, status };

  queryClient.setQueryData<EventI[]>(eventKeys.ownLists(), (old) => upsertById(old, nextEvent));

  queryClient.setQueryData<EventI[]>(eventKeys.publicLists(), (old) => upsertById(old, nextEvent));
}

function syncEventMediaInCaches(queryClient: QueryClient, event: EventI) {
  syncEventCaches(queryClient, event);
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

export function usePublishEventMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (eventId: string) => publishEventFn(eventId),
    onSuccess: (res, eventId) => {
      if (!res.success) {
        toast.error(res.message || "Erro ao publicar evento");
        return;
      }

      syncEventStatusInCaches(queryClient, eventId, "active");
      toast.success("Evento publicado com sucesso!");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}

export interface SyncEventMediaInput {
  eventId: string;
  values: EventCreateOutputI;
  logoChanges: ImageFieldChange;
  bannerChanges: ImageFieldChange;
  galleryChanges: ImageFieldChange;
}

export function useSyncEventMediaMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ eventId, values, logoChanges, bannerChanges, galleryChanges }: SyncEventMediaInput) => {
      let latestEvent: EventI | null = null;

      if (galleryChanges.removed.length > 0) {
        for (const image of galleryChanges.removed) {
          const res = await removeImageFromTheEventGalleryFn(eventId, { url: image.url });
          if (res.success) {
            latestEvent = res.data;
            syncEventMediaInCaches(queryClient, res.data);
          }
        }
      }

      if (galleryChanges.added.length > 0) {
        for (const image of galleryChanges.added) {
          const res = await addImageToTheEventGalleryFn(eventId, { url: image.url });
          if (res.success) {
            latestEvent = res.data;
            syncEventMediaInCaches(queryClient, res.data);
          }
        }
      }

      if (bannerChanges.added.length > 0) {
        const res = await setEventBannerFn(eventId, { url: bannerChanges.added.at(-1)?.url ?? values.banner_url ?? "" });
        if (res.success) {
          latestEvent = res.data;
          syncEventMediaInCaches(queryClient, res.data);
        }
      } else if (bannerChanges.removed.length > 0) {
        const res = await unsetEventBannerFn(eventId);
        if (res.success) {
          latestEvent = res.data;
          syncEventMediaInCaches(queryClient, res.data);
        }
      }

      if (logoChanges.added.length > 0) {
        const res = await setEventLogoFn(eventId, { url: logoChanges.added.at(-1)?.url ?? values.logo_url ?? "" });
        if (res.success) {
          latestEvent = res.data;
          syncEventMediaInCaches(queryClient, res.data);
        }
      } else if (logoChanges.removed.length > 0) {
        const res = await unsetEventLogoFn(eventId);
        if (res.success) {
          latestEvent = res.data;
          syncEventMediaInCaches(queryClient, res.data);
        }
      }

      return latestEvent;
    },
    onError: () => toast.error("Erro ao sincronizar imagens do evento"),
  });
}
