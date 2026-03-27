import { AwsClient } from "aws4fetch";
import type {
  StorageUploadRequest,
  StorageUploadResponse,
  StorageModerateRequest,
  StorageModerateResponse,
  StorageErrorResponse
} from "../model";

const ALLOWED_TYPES = ["image/png", "image/jpeg", "image/webp"];
const MAX_SIZE = 10 * 1024 * 1024; // 10MB

// Initialize AWS client for MinIO/S3-compatible storage
const getAwsClient = (env: Env) =>
  new AwsClient({
    accessKeyId: env.MINIO_ACCESS_KEY,
    secretAccessKey: env.MINIO_SECRET_KEY,
    service: "s3",
    region: "auto",
  });

async function json<T>(req: Request): Promise<T> {
  return req.json();
}
export async function handleStorageUpload(
  request: Request,
  env: Env
): Promise<Response> {
  const { filename, contentType, size } = await json<StorageUploadRequest>(request);

  if (!ALLOWED_TYPES.includes(contentType)) {
    return Response.json(
      { error: "Only PNG, JPEG, and WebP are allowed" } as StorageErrorResponse,
      { status: 400 }
    );
  }

  if (size > MAX_SIZE) {
    return Response.json(
      { error: "File exceeds 10MB limit" } as StorageErrorResponse,
      { status: 400 }
    );
  }

  const uploadUrl = await getUploadUrl(filename, contentType, env);
  const publicUrl = `${env.MINIO_ENDPOINT}/${env.BUCKET_NAME}/${filename}`;

  return Response.json({ uploadUrl, key: filename, publicUrl } as StorageUploadResponse);
}

export async function handleStorageModerate(
  request: Request,
  env: Env
): Promise<Response> {
  try {
    const { key } = await json<StorageModerateRequest>(request);
    const isSafe = await moderateFile(key, env);

    if (!isSafe) {
      await deleteFile(key, env);
      return Response.json({ approved: false } as StorageModerateResponse);
    }

    return Response.json({ approved: true } as StorageModerateResponse);
  } catch (error) {
    return Response.json(
      { error: error instanceof Error ? error.message : "Moderation failed" } as StorageErrorResponse,
      { status: 500 }
    );
  }
}

async function getUploadUrl(key: string, contentType: string, env: Env): Promise<string> {
  const aws = getAwsClient(env);
  const uploadUrl = new URL(`${env.MINIO_ENDPOINT}/${env.BUCKET_NAME}/${key}`);
  uploadUrl.searchParams.set("X-Amz-Expires", "300"); // 5 minutes

  const signed = await aws.sign(
    new Request(uploadUrl, {
      method: "PUT",
      headers: { "Content-Type": contentType },
    }),
    { aws: { signQuery: true } }
  );

  return signed.url;
}

async function moderateFile(key: string, env: Env): Promise<boolean> {
  const aws = getAwsClient(env);
  const downloadUrl = new URL(`${env.MINIO_ENDPOINT}/${env.BUCKET_NAME}/${key}`);
  downloadUrl.searchParams.set("X-Amz-Expires", "60"); // 1 minute

  const signed = await aws.sign(
    new Request(downloadUrl, { method: "GET" }),
    { aws: { signQuery: true } }
  );

  const fileRes = await fetch(signed.url);
  if (!fileRes.ok) throw new Error("Failed to fetch file for moderation");

  const buffer = await fileRes.arrayBuffer();

  const response = await env.AI.run("@cf/llava-hf/llava-1.5-7b-hf", {
    prompt: "Does this image contain any explicit, violent, or inappropriate content? Reply with only 'safe' or 'unsafe'.",
    image: [...new Uint8Array(buffer)],
    max_tokens: 5,
  });

  const result = (response as { description: string }).description.trim().toLowerCase();
  return result.startsWith("safe");
}

async function deleteFile(key: string, env: Env): Promise<void> {
  const aws = getAwsClient(env);
  const deleteUrl = new URL(`${env.MINIO_ENDPOINT}/${env.BUCKET_NAME}/${key}`);

  const signed = await aws.sign(
    new Request(deleteUrl, { method: "DELETE" }),
    { aws: { signQuery: true } }
  );

  const res = await fetch(signed.url, { method: "DELETE" });
  if (!res.ok) throw new Error(`Failed to delete file: ${res.statusText}`);
}
