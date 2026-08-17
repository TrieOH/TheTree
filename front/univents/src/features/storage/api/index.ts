import type {
  StoragePreprocessResponse,
  StorageUploadRequest,
  StorageUploadResponse,
} from "../model";

const MODERATION_MAX_EDGE = 448;
const MODERATION_WEBP_QUALITY = 0.82;
const ALLOWED_IMAGE_TYPES = new Set(["image/png", "image/jpeg", "image/webp"]);

export class StorageImageError extends Error {
  constructor(
    message: string,
    readonly code: string,
    readonly status?: number,
  ) {
    super(message);
    this.name = "StorageImageError";
  }
}

async function resizeImageForModeration(file: File): Promise<File> {
  if (!file.type.startsWith("image/")) return file;

  if (typeof createImageBitmap !== "function") return file;
  if (typeof document === "undefined") return file;

  try {
    const bitmap = await createImageBitmap(file);
    const scale = Math.min(
      1,
      MODERATION_MAX_EDGE / Math.max(bitmap.width, bitmap.height),
    );
    const targetWidth = Math.max(1, Math.round(bitmap.width * scale));
    const targetHeight = Math.max(1, Math.round(bitmap.height * scale));

    if (
      targetWidth === bitmap.width &&
      targetHeight === bitmap.height &&
      file.type === "image/webp"
    ) {
      bitmap.close();
      return file;
    }

    const canvas = document.createElement("canvas");
    canvas.width = targetWidth;
    canvas.height = targetHeight;

    const context = canvas.getContext("2d");
    if (!context) {
      bitmap.close();
      return file;
    }

    context.drawImage(bitmap, 0, 0, targetWidth, targetHeight);
    bitmap.close();

    const blob = await new Promise<Blob | null>((resolve) => {
      canvas.toBlob(
        (nextBlob) => resolve(nextBlob),
        "image/webp",
        MODERATION_WEBP_QUALITY,
      );
    });

    if (!blob || blob.size >= file.size) return file;

    const baseName = file.name.replace(/\.[^.]+$/, "") || "image";
    return new File([blob], `${baseName}.webp`, {
      type: "image/webp",
      lastModified: file.lastModified,
    });
  } catch {
    return file;
  }
}

export async function preprocessImageUpload(
  file: File,
  path?: string,
  idempotencyKey?: string,
): Promise<string> {
  if (!ALLOWED_IMAGE_TYPES.has(file.type)) {
    throw new StorageImageError(
      "Formato de imagem não suportado. Use PNG, JPG ou WebP.",
      "UNSUPPORTED_IMAGE_TYPE",
      400,
    );
  }

  const moderationFile = await resizeImageForModeration(file);
  const formData = new FormData();
  formData.append("file", file);
  if (moderationFile !== file) {
    formData.append("moderationFile", moderationFile);
  }
  if (path) formData.append("path", path);
  if (idempotencyKey) formData.append("idempotencyKey", idempotencyKey);

  let res: Response;
  try {
    res = await fetch("/storage/image/preprocess", {
      method: "POST",
      body: formData,
    });
  } catch {
    throw new StorageImageError(
      "Não foi possível conectar ao servidor.",
      "NETWORK_ERROR",
    );
  }

  if (!res.ok) {
    const errorData = (await res.json().catch(() => ({}))) as {
      error?: string;
    };
    throw new StorageImageError(
      errorData.error ?? "Não foi possível processar a imagem.",
      `STORAGE_HTTP_${res.status}`,
      res.status,
    );
  }

  const data: StoragePreprocessResponse = await res.json();
  if (!data.approved || !data.publicUrl) {
    throw new StorageImageError(
      "A imagem não foi aprovada pela moderação.",
      "MODERATION_REJECTED",
      422,
    );
  }

  return data.publicUrl;
}

export const uploadFile = async (
  file: File,
  path?: string,
): Promise<string> => {
  const filename = path
    ? `${path}/${Date.now()}-${file.name}`
    : `${Date.now()}-${file.name}`;

  const uploadPayload: StorageUploadRequest = {
    filename,
    contentType: file.type,
    size: file.size,
  };

  const uploadRes = await fetch("/storage/upload", {
    method: "POST",
    body: JSON.stringify(uploadPayload),
  });

  if (!uploadRes.ok) {
    const errorData: { error?: string } = await uploadRes.json();
    throw new Error(errorData.error ?? "Failed to get upload URL");
  }
  const { uploadUrl, publicUrl }: StorageUploadResponse =
    await uploadRes.json();

  const putRes = await fetch(uploadUrl, {
    method: "PUT",
    body: file,
    headers: { "Content-Type": file.type },
  });

  if (!putRes.ok) throw new Error("Failed to upload file");

  return publicUrl;
};
