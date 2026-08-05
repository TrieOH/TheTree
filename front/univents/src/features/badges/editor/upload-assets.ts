import { uploadFile } from "@/features/storage/api";
import type { BadgeTemplateCreate } from "../model";

function isDataImage(value: string | null): value is string {
  return Boolean(value?.startsWith("data:image/"));
}

function dataImageToFile(dataUrl: string, name: string): File {
  const [metadata, encoded = ""] = dataUrl.split(",", 2);
  const mimeType = metadata.match(/^data:([^;]+)/)?.[1] ?? "image/png";
  const bytes = Uint8Array.from(atob(encoded), (character) =>
    character.charCodeAt(0),
  );
  const extension = mimeType.split("/")[1]?.replace("jpeg", "jpg") ?? "png";
  return new File([bytes], `${name}.${extension}`, { type: mimeType });
}

export async function uploadBadgeAssets(
  draft: BadgeTemplateCreate,
  eventId: string,
  editionId: string,
): Promise<BadgeTemplateCreate> {
  const path = `events/${eventId}/editions/${editionId}/badges`;
  const uploads = new Map<string, Promise<string>>();
  const upload = (source: string, name: string) => {
    const pending = uploads.get(source);
    if (pending) return pending;
    const next = uploadFile(dataImageToFile(source, name), path);
    uploads.set(source, next);
    return next;
  };
  const background = isDataImage(draft.design_data.background)
    ? await upload(draft.design_data.background, "background")
    : draft.design_data.background;
  const elements = await Promise.all(
    draft.design_data.elements.map(async (element) =>
      element.type === "image" && isDataImage(element.src)
        ? { ...element, src: await upload(element.src, `image-${element.id}`) }
        : element,
    ),
  );
  return {
    ...draft,
    design_data: { ...draft.design_data, background, elements },
  };
}
