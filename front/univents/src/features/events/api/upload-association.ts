import { orvalData } from "@trieoh/api-client";
import { listJoinedEvents, listOwnedEvents } from "@trieoh/univents-api";
import {
  registerUploadAssociationHandler,
  UploadAssociationError,
} from "@/features/upload-queue";
import { getContext } from "@/integrations/tanstack-query/root-provider";
import type { EventI } from "../model";
import { patchEventFn } from "./index";
import { syncEventCaches } from "./mutations";

type ImageField = "logo_url" | "banner_url";

function patchData(event: EventI, field: ImageField, url: string | null) {
  return {
    full_name: event.full_name,
    slug: event.slug,
    acronym: event.acronym,
    description: event.description,
    contact_email: event.contact_email,
    logo_url: field === "logo_url" ? url : event.logo_url,
    banner_url: field === "banner_url" ? url : event.banner_url,
  };
}

async function associateEventImage(
  task: {
    owner: { id: string };
    association?: { input?: Record<string, unknown> };
  },
  uploadedUrl: string,
) {
  const field = task.association?.input?.field;
  if (field !== "logo_url" && field !== "banner_url") {
    throw new UploadAssociationError("Campo de imagem inválido.", {
      status: 400,
    });
  }

  const [owned, joined] = await Promise.all([
    listOwnedEvents().then(orvalData<EventI[]>),
    listJoinedEvents().then(orvalData<EventI[]>),
  ]);
  const event = [...owned, ...joined].find((item) => item.id === task.owner.id);
  if (!event) {
    throw new UploadAssociationError("Evento não encontrado.", { status: 404 });
  }

  try {
    const response = await patchEventFn(
      event.id,
      patchData(event, field, uploadedUrl),
    );
    if (!response) {
      throw new UploadAssociationError("Não foi possível associar a imagem.", {
        status: 500,
      });
    }
    syncEventCaches(getContext().queryClient, response);
  } catch (err: unknown) {
    const error = err instanceof Error ? err : undefined;
    const details =
      err && typeof err === "object"
        ? (err as { status?: number; message?: string })
        : undefined;
    throw new UploadAssociationError(
      error?.message ||
        details?.message ||
        "Não foi possível associar a imagem.",
      { status: details?.status || 500 },
    );
  }
}

registerUploadAssociationHandler("event-image", associateEventImage);

export { patchData };
