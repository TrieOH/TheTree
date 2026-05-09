import { AwsClient } from "aws4fetch";
import type {
  StorageUploadRequest,
  StorageModerateRequest
} from "../model";

const ALLOWED_TYPES = ["image/png", "image/jpeg", "image/webp"];
const MAX_SIZE = 10 * 1024 * 1024; // 10MB

/**
 * Utility to construct a sanitized URL for S3/MinIO
 */
function getS3Url(key: string, env: Env): URL {
  const endpoint = env.MINIO_ENDPOINT.trim();
  if (!endpoint || !/^https?:\/\//.test(endpoint)) {
    throw new Error("Invalid or missing MINIO_ENDPOINT protocol (http/https)");
  }

  const baseUrl = endpoint.replace(/\/+$/, "");
  const cleanKey = key.replace(/^\/+/, "");
  return new URL(`${baseUrl}/${env.BUCKET_NAME}/${cleanKey}`);
}

/**
 * Validates that all required environment variables are present
 */
function validateEnv(env: Env) {
  const keys: (keyof Env)[] = ["MINIO_ENDPOINT", "BUCKET_NAME", "MINIO_ACCESS_KEY", "MINIO_SECRET_KEY"];
  for (const key of keys) {
    if (!env[key]) throw new Error(`Missing environment variable: ${key}`);
  }
}

const getAwsClient = (env: Env) =>
  new AwsClient({
    accessKeyId: env.MINIO_ACCESS_KEY,
    secretAccessKey: env.MINIO_SECRET_KEY,
    service: "s3",
    region: "auto",
  });

export async function handleStorageUpload(request: Request, env: Env): Promise<Response> {
  try {
    validateEnv(env);
    const { filename, contentType, size } = await request.json<StorageUploadRequest>();

    if (!ALLOWED_TYPES.includes(contentType)) {
      return Response.json({ error: "Only PNG, JPEG, and WebP are allowed" }, { status: 400 });
    }

    if (size > MAX_SIZE) {
      return Response.json({ error: "File exceeds 10MB limit" }, { status: 400 });
    }

    const aws = getAwsClient(env);
    const uploadUrl = getS3Url(filename, env);
    uploadUrl.searchParams.set("X-Amz-Expires", "300");

    const signed = await aws.sign(
      new Request(uploadUrl, {
        method: "PUT",
        headers: { "Content-Type": contentType },
      }),
      { aws: { signQuery: true } }
    );

    const publicUrl = `${getS3Url("", env).toString()}${filename}`;

    return Response.json({
      uploadUrl: signed.url,
      key: filename,
      publicUrl,
    });
  } catch (error) {
    return Response.json(
      { error: error instanceof Error ? error.message : "Upload failed" },
      { status: 500 }
    );
  }
}

export async function handleStorageModerate(request: Request, env: Env): Promise<Response> {
  try {
    validateEnv(env);
    const { key } = await request.json<StorageModerateRequest>();

    const isSafe = await moderateFile(key, env);
    if (!isSafe) {
      await deleteFile(key, env);
      return Response.json({ approved: false });
    }

    return Response.json({ approved: true });
  } catch (error) {
    return Response.json(
      { error: error instanceof Error ? error.message : "Moderation failed" },
      { status: 500 }
    );
  }
}

async function moderateFile(key: string, env: Env): Promise<boolean> {
  const aws = getAwsClient(env);
  const downloadUrl = getS3Url(key, env);
  downloadUrl.searchParams.set("X-Amz-Expires", "60");

  const signed = await aws.sign(new Request(downloadUrl, { method: "GET" }), {
    aws: { signQuery: true },
  });

  const res = await fetch(signed.url);
  if (!res.ok) throw new Error("Failed to fetch file for moderation");

  const buffer = await res.arrayBuffer();
  const response = await env.AI.run("@cf/llava-hf/llava-1.5-7b-hf", {
    prompt: "Does this image contain any explicit, violent, or inappropriate content? Reply with only 'safe' or 'unsafe'.",
    image: Array.from(new Uint8Array(buffer)),
    max_tokens: 5,
  });

  return (response).description.trim().toLowerCase().startsWith("safe");
}

async function deleteFile(key: string, env: Env): Promise<void> {
  const aws = getAwsClient(env);
  const deleteUrl = getS3Url(key, env);

  const signed = await aws.sign(new Request(deleteUrl, { method: "DELETE" }), {
    aws: { signQuery: true },
  });

  const res = await fetch(signed.url, { method: "DELETE" });
  if (!res.ok) throw new Error("Failed to delete file");
}
