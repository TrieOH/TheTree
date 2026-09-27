import { uploadFile } from "@/features/storage/api/index";
import type { CertificationTemplateCreateI } from "../model";

function isDataImage(value: string | null): value is string {
  return Boolean(value?.startsWith("data:image/"));
}

function dataImageToFile(dataUrl: string, name: string): File {
  const [metadata, encoded = ""] = dataUrl.split(",", 2);
  const mimeType = metadata?.match(/^data:([^;]+)/)?.[1] ?? "image/png";
  const binary = atob(encoded);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i);
  }
  const extension = mimeType.split("/")[1]?.replace("jpeg", "jpg") ?? "png";
  return new File([bytes], `${name}.${extension}`, { type: mimeType });
}

export async function uploadCertificateAssets(
  draft: CertificationTemplateCreateI,
  eventId: string,
  editionId: string,
): Promise<CertificationTemplateCreateI> {
  const path = `events/${eventId}/editions/${editionId}/certificates`;
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
    draft.design_data.elements.map(async (element) => {
      if (
        (element.type === "image" || element.type === "signature") &&
        isDataImage(element.src)
      ) {
        return {
          ...element,
          src: await upload(element.src, `${element.type}-${element.id}`),
        };
      }
      return element;
    }),
  );

  return {
    ...draft,
    design_data: { ...draft.design_data, background, elements },
  };
}
