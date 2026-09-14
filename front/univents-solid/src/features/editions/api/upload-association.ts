import { orvalData } from "@trieoh/api-client";
import {
  listDraftEditions,
  listPublicEditions,
  patchEdition,
} from "@trieoh/univents-api";
import type { Edition } from "@trieoh/univents-api/schemas";
import {
  registerUploadAssociationHandler,
  UploadAssociationError,
} from "@/features/upload-queue";
import { normalizeEdition, type EditionI } from "../model";

type ImageField = "logo_url" | "banner_url";

export function patchEditionData(
  edition: EditionI,
  field: ImageField,
  url: string | null,
) {
  return {
    name: edition.name,
    slug: edition.slug,
    starts_at: edition.starts_at,
    ends_at: edition.ends_at,
    tagline: edition.tagline,
    description: edition.description,
    registration_opens_at: edition.registration_opens_at,
    location_name: edition.location_name,
    location_description: edition.location_description,
    contact_email: edition.contact_email,
    logo_url: field === "logo_url" ? url : edition.logo_url,
    banner_url: field === "banner_url" ? url : edition.banner_url,
  };
}

async function associateEditionImage(
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

  const eventId = task.association?.input?.eventId;
  if (typeof eventId !== "string") {
    throw new UploadAssociationError("Evento da edição não encontrado.", {
      status: 400,
    });
  }

  const [publicEditions, draftEditions] = await Promise.all([
    listPublicEditions(eventId, { public: true })
      .then(orvalData<Edition[]>)
      .then((list) => (list ?? []).map(normalizeEdition)),
    listDraftEditions(eventId)
      .then(orvalData<Edition[]>)
      .then((list) => (list ?? []).map(normalizeEdition)),
  ]);

  const editions = [...publicEditions, ...draftEditions];
  const edition = editions.find((item) => item.id === task.owner.id);
  if (!edition) {
    throw new UploadAssociationError("Edição não encontrada.", { status: 404 });
  }

  try {
    const response = await patchEdition(
      eventId,
      edition.id,
      patchEditionData(edition, field, uploadedUrl),
    );
    if (!response) {
      throw new UploadAssociationError("Não foi possível associar a imagem.", {
        status: 500,
      });
    }
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

registerUploadAssociationHandler("edition-image", associateEditionImage);
