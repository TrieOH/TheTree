import { uploadFile } from "@/features/storage/api";
import type { CertificationTemplateCreateI } from "../model";

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

export async function uploadCertificateAssets(
  draft: CertificationTemplateCreateI,
  eventId: string,
  editionId: string,
): Promise<CertificationTemplateCreateI> {
  const path = `events/${eventId}/editions/${editionId}/certificates`;
  const uploads = new Map<string, Promise<string>>();

  function uploadDataImage(dataUrl: string, name: string) {
    const existing = uploads.get(dataUrl);
    if (existing) return existing;
    const upload = uploadFile(dataImageToFile(dataUrl, name), path);
    uploads.set(dataUrl, upload);
    return upload;
  }

  const backgroundSource = draft.url ?? draft.data.background;
  const backgroundUrl = isDataImage(backgroundSource)
    ? await uploadDataImage(backgroundSource, "background")
    : backgroundSource;

  const elements = await Promise.all(
    draft.data.elements.map(async (element) => {
      if (
        (element.type === "image" || element.type === "signature") &&
        isDataImage(element.src)
      ) {
        return {
          ...element,
          src: await uploadDataImage(
            element.src,
            `${element.type}-${element.id}`,
          ),
        };
      }
      return element;
    }),
  );

  return {
    ...draft,
    url: backgroundUrl,
    data: {
      ...draft.data,
      background: backgroundUrl,
      elements,
    },
  };
}
