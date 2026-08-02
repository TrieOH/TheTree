import { orvalData } from "@trieoh/api-client";
import { listDraftEditions, listPublicEditions } from "@trieoh/univents-api";
import {
  registerUploadAssociationHandler,
  UploadAssociationError,
} from "@/features/upload-queue";
import { getContext } from "@/integrations/tanstack-query/root-provider";
import type { EditionApiI, EditionI } from "../model";
import { patchEditionFn } from "./index";
import { editionKeys } from "./query-keys";

registerUploadAssociationHandler("edition-image", async (task, url) => {
  const field = task.association?.input?.field;
  if (field !== "logo_url" && field !== "banner_url") {
    throw new UploadAssociationError("Campo de imagem inválido.", {
      status: 400,
    });
  }
  const eventId = task.association?.input?.eventId;
  if (typeof eventId !== "string")
    throw new UploadAssociationError("Evento da edição não encontrado.", {
      status: 400,
    });
  const [publicEditions, draftEditions] = await Promise.all([
    listPublicEditions(eventId, { public: true }).then(
      orvalData<EditionApiI[]>,
    ),
    listDraftEditions(eventId).then(orvalData<EditionApiI[]>),
  ]);
  const editions = [...publicEditions, ...draftEditions];
  const edition = editions.find((item) => item.id === task.owner.id) as
    | EditionI
    | undefined;
  if (!edition)
    throw new UploadAssociationError("Edição não encontrada.", { status: 404 });
  const updated = await patchEditionFn(edition.event_id, edition.id, {
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
  });
  getContext().queryClient.setQueryData<EditionI[]>(
    editionKeys.adminListByEvent(updated.event_id),
    (old = []) => old.map((item) => (item.id === updated.id ? updated : item)),
  );
  getContext().queryClient.setQueryData<EditionI[]>(
    editionKeys.publicListByEvent(updated.event_id),
    (old = []) => old.map((item) => (item.id === updated.id ? updated : item)),
  );
});
