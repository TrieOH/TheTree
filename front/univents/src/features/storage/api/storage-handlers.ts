import { AwsClient } from "aws4fetch";
import type { StorageUploadRequest } from "../model";

type StorageOptionalEnvKeys =
  | "STORAGE_IMAGE_ALLOWED_TYPES"
  | "STORAGE_IMAGE_MAX_SIZE_BYTES"
  | "STORAGE_IMAGE_UPLOAD_EXPIRES_SECONDS"
  | "STORAGE_IMAGE_MODERATION_MODEL"
  | "STORAGE_IMAGE_MODERATION_PROMPT";

type StorageRuntimeEnv =
  Omit<Env, StorageOptionalEnvKeys>
  & Partial<Pick<Env, StorageOptionalEnvKeys>>;

const DEFAULT_ALLOWED_TYPES = ["image/png", "image/jpeg", "image/webp"];
const DEFAULT_MAX_SIZE_BYTES = 10 * 1024 * 1024; // 10MB
const DEFAULT_EXPIRES_SECONDS = 300;
const DEFAULT_MODERATION_MODEL = "@cf/llava-hf/llava-1.5-7b-hf";
const DEFAULT_MODERATION_PROMPT =
  "Does this image contain any explicit, violent, or inappropriate content? Reply with only 'safe' or 'unsafe'.";

function getS3Url(key: string, env: StorageRuntimeEnv): URL {
  const endpoint = env.MINIO_ENDPOINT.trim();
  if (!endpoint || !/^https?:\/\//.test(endpoint))
    throw new Error("Invalid or missing MINIO_ENDPOINT protocol (http/https)");

  const baseUrl = endpoint.replace(/\/+$/, "");
  const cleanKey = key.replace(/^\/+/, "");
  return new URL(`${baseUrl}/${env.BUCKET_NAME}/${cleanKey}`);
}

function getAllowedTypes(env: StorageRuntimeEnv) {
  const raw = env.STORAGE_IMAGE_ALLOWED_TYPES?.trim();
  if (!raw) return DEFAULT_ALLOWED_TYPES;

  const parsed = raw
    .split(",")
    .map((type) => type.trim())
    .filter(Boolean);

  return parsed.length > 0 ? parsed : DEFAULT_ALLOWED_TYPES;
}

function getMaxSize(env: StorageRuntimeEnv) {
  const raw = env.STORAGE_IMAGE_MAX_SIZE_BYTES?.trim();
  if (!raw) return DEFAULT_MAX_SIZE_BYTES;

  const parsed = Number(raw);
  return Number.isFinite(parsed) && parsed > 0 ? parsed : DEFAULT_MAX_SIZE_BYTES;
}

function getUploadExpiresSeconds(env: StorageRuntimeEnv) {
  const raw = env.STORAGE_IMAGE_UPLOAD_EXPIRES_SECONDS?.trim();
  if (!raw) return DEFAULT_EXPIRES_SECONDS;

  const parsed = Number(raw);
  return Number.isFinite(parsed) && parsed > 0 ? parsed : DEFAULT_EXPIRES_SECONDS;
}

function getModerationModel(env: StorageRuntimeEnv) {
  return env.STORAGE_IMAGE_MODERATION_MODEL?.trim() || DEFAULT_MODERATION_MODEL;
}

function getModerationPrompt(env: StorageRuntimeEnv) {
  return env.STORAGE_IMAGE_MODERATION_PROMPT?.trim() || DEFAULT_MODERATION_PROMPT;
}

function buildAllowedTypesErrorMessage(types: string[]) {
  return `Only ${types.join(", ")} are allowed`;
}

/**
 * Validates that all required environment variables are present
 */
function validateEnv(env: StorageRuntimeEnv) {
  const keys: (keyof Env)[] = ["MINIO_ENDPOINT", "BUCKET_NAME", "MINIO_ACCESS_KEY", "MINIO_SECRET_KEY"];
  for (const key of keys) {
    if (!env[key]) throw new Error(`Missing environment variable: ${key}`);
  }
}

const getAwsClient = (env: StorageRuntimeEnv) =>
  new AwsClient({
    accessKeyId: env.MINIO_ACCESS_KEY,
    secretAccessKey: env.MINIO_SECRET_KEY,
    service: "s3",
    region: "auto",
  });

function sanitizePath(path?: string | null) {
  const value = (path ?? "").trim();
  if (!value) return "";
  return value.replace(/^\/+|\/+$/g, "").replace(/\.\./g, "");
}

function buildStorageKey(filename: string, path?: string | null) {
  const cleanPath = sanitizePath(path);
  const cleanFilename = filename.replace(/\s+/g, "-").replace(/[^\w.-]/g, "");
  const baseName = `${Date.now()}-${cleanFilename}`;
  return cleanPath ? `${cleanPath}/${baseName}` : baseName;
}

async function moderateFileBytes(file: File, env: StorageRuntimeEnv): Promise<boolean> {
  const buffer = await file.arrayBuffer();
  const response = await env.AI.run(getModerationModel(env), {
    prompt: getModerationPrompt(env),
    image: Array.from(new Uint8Array(buffer)),
    max_tokens: 5,
    temperature: 0.0,
  });
  return String(response.description ?? "").trim().toLowerCase().startsWith("safe");
}

async function putFileToStorage(file: File, key: string, env: StorageRuntimeEnv): Promise<string> {
  const aws = getAwsClient(env);
  const uploadUrl = getS3Url(key, env);
  uploadUrl.searchParams.set("X-Amz-Expires", String(getUploadExpiresSeconds(env)));

  const signed = await aws.sign(
    new Request(uploadUrl, {
      method: "PUT",
      headers: { "Content-Type": file.type },
    }),
    { aws: { signQuery: true } }
  );

  const res = await fetch(signed.url, {
    method: "PUT",
    body: file,
    headers: { "Content-Type": file.type },
  });

  if (!res.ok) throw new Error("Failed to upload file");

  return `${getS3Url("", env).toString()}${key}`;
}

async function readImageFormData(request: Request) {
  const formData = await request.formData();
  const file = formData.get("file");
  const path = formData.get("path");

  if (!(file instanceof File)) throw new Error("Missing file");

  return {
    file,
    path: typeof path === "string" ? path : "",
  };
}

export async function handleStorageUpload(request: Request, env: StorageRuntimeEnv): Promise<Response> {
  try {
    validateEnv(env);
    const { filename, contentType, size } = await request.json<StorageUploadRequest>();
    const allowedTypes = getAllowedTypes(env);
    const maxSize = getMaxSize(env);

    if (!allowedTypes.includes(contentType)) {
      return Response.json({ error: buildAllowedTypesErrorMessage(allowedTypes) }, { status: 400 });
    }

    if (size > maxSize) {
      return Response.json({ error: "File exceeds 10MB limit" }, { status: 400 });
    }

    const aws = getAwsClient(env);
    const uploadUrl = getS3Url(filename, env);
    uploadUrl.searchParams.set("X-Amz-Expires", String(getUploadExpiresSeconds(env)));

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

export async function handleStorageImagePreprocess(request: Request, env: StorageRuntimeEnv): Promise<Response> {
  try {
    validateEnv(env);
    const { file, path } = await readImageFormData(request);
    const allowedTypes = getAllowedTypes(env);
    const maxSize = getMaxSize(env);

    if (!allowedTypes.includes(file.type)) {
      return Response.json({ error: buildAllowedTypesErrorMessage(allowedTypes) }, { status: 400 });
    }

    if (file.size > maxSize) {
      return Response.json({ error: "File exceeds 10MB limit" }, { status: 400 });
    }

    const approved = await moderateFileBytes(file, env);
    if (!approved) return Response.json({ approved: false });

    const key = buildStorageKey(file.name, path);
    const publicUrl = await putFileToStorage(file, key, env);
    return Response.json({ approved: true, publicUrl });
  } catch (error) {
    return Response.json(
      { error: error instanceof Error ? error.message : "Preprocessing failed" },
      { status: 500 }
    );
  }
}
