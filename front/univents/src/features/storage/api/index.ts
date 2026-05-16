import type {
  StorageUploadRequest,
  StorageUploadResponse,
  StorageModerateRequest,
  StorageModerateResponse
} from "../model";

export const uploadAndModerateFile = async (file: File, path?: string): Promise<string> => {
  // 1. Get signed URL
  const filename = path ? `${path}/${Date.now()}-${file.name}` : `${Date.now()}-${file.name}`;

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
  const { uploadUrl, key, publicUrl }: StorageUploadResponse = (await uploadRes.json());

  // 2. Upload to MinIO
  const putRes = await fetch(uploadUrl, {
    method: "PUT",
    body: file,
    headers: { "Content-Type": file.type },
  });

  if (!putRes.ok) throw new Error("Failed to upload file");

  // 3. Moderate
  const moderatePayload: StorageModerateRequest = { key };
  const modRes = await fetch("/storage/moderate", {
    method: "POST",
    body: JSON.stringify(moderatePayload),
  });

  if (!modRes.ok) throw new Error("Failed to moderate file");

  const { approved }: StorageModerateResponse = (await modRes.json());

  if (!approved) throw new Error("Image not approved by moderation");

  return publicUrl;
};
