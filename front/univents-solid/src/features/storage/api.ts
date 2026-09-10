const allowed = new Set(["image/png", "image/jpeg", "image/webp"]);

export async function uploadProfileImage(file: File, field: "pfpUrl" | "bannerUrl") {
  if (!allowed.has(file.type)) throw new Error("Use uma imagem PNG, JPG ou WebP.");
  if (file.size > 10 * 1024 * 1024) throw new Error("A imagem deve ter no máximo 10 MB.");
  const data = new FormData();
  data.append("file", file);
  data.append("path", "profiles/images");
  data.append("idempotencyKey", crypto.randomUUID());
  data.append("field", field);
  const response = await fetch("/storage/image/preprocess", { method: "POST", body: data });
  const result = (await response.json().catch(() => ({}))) as { publicUrl?: string; error?: string };
  if (!response.ok || !result.publicUrl) throw new Error(result.error ?? "A imagem não foi aprovada.");
  return result.publicUrl;
}
