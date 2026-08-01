import {
  registerUploadAssociationHandler,
  UploadAssociationError,
  uploadAssociationErrorFromResponse,
} from "@/features/upload-queue";
import { getContext } from "@/integrations/tanstack-query/root-provider";
import { authQueryFetcher } from "@/shared/lib/api/fetch";
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
    authQueryFetcher<EditionApiI[]>(`/events/${eventId}/editions`),
    authQueryFetcher<EditionApiI[]>(`/events/${eventId}/editions/draft`),
  ]);
  const editions = [...publicEditions, ...draftEditions];
  const edition = editions.find((item) => item.id === task.owner.id) as
    | EditionI
    | undefined;
  if (!edition)
    throw new UploadAssociationError("Edição não encontrada.", { status: 404 });
  const response = await patchEditionFn(edition.event_id, edition.id, {
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
  if (!response.success)
    throw uploadAssociationErrorFromResponse(
      response,
      "Não foi possível associar a imagem.",
    );
  getContext().queryClient.setQueryData<EditionI[]>(
    editionKeys.adminListByEvent(response.data.event_id),
    (old = []) =>
      old.map((item) => (item.id === response.data.id ? response.data : item)),
  );
  getContext().queryClient.setQueryData<EditionI[]>(
    editionKeys.publicListByEvent(response.data.event_id),
    (old = []) =>
      old.map((item) => (item.id === response.data.id ? response.data : item)),
  );
});
