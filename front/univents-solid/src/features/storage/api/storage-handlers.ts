import { AwsClient } from "aws4fetch";
import { Data, Effect } from "effect";
import { privateJsonResponse } from "@/shared/lib/http-cache";
import type { StorageUploadRequest } from "../model";

type StorageOptionalEnvKeys =
  | "STORAGE_IMAGE_ALLOWED_TYPES"
  | "STORAGE_IMAGE_MAX_SIZE_BYTES"
  | "STORAGE_IMAGE_UPLOAD_EXPIRES_SECONDS"
  | "STORAGE_IMAGE_MODERATION_MODEL"
  | "STORAGE_IMAGE_MODERATION_PROMPT";

type StorageRuntimeEnv = Omit<Env, StorageOptionalEnvKeys> &
  Partial<Pick<Env, StorageOptionalEnvKeys>>;

const DEFAULT_ALLOWED_TYPES = ["image/png", "image/jpeg", "image/webp"];
const DEFAULT_MAX_SIZE_BYTES = 10 * 1024 * 1024; // 10MB
const DEFAULT_EXPIRES_SECONDS = 300;
const DEFAULT_MODERATION_MODEL = "@cf/llava-hf/llava-1.5-7b-hf";
const DEFAULT_MODERATION_PROMPT =
  "Does this image contain any explicit, violent, or inappropriate content? Reply with only 'safe' or 'unsafe'.";

export class StorageValidationError extends Data.TaggedError(
  "StorageValidationError",
)<{
  readonly message: string;
}> { }

export class StorageConfigError extends Data.TaggedError("StorageConfigError")<{
  readonly message: string;
}> { }

export class StorageUploadError extends Data.TaggedError("StorageUploadError")<{
  readonly message: string;
  readonly cause?: unknown;
}> { }

export function getStorageConfig(env: StorageRuntimeEnv) {
  const endpoint = (env.S3_ENDPOINT ?? "").trim();
  const bucket = (env.S3_BUCKET ?? "").trim();
  const accessKeyId = (env.S3_ACCESS_KEY_ID ?? "").trim();
  const secretAccessKey = (env.S3_SECRET_ACCESS_KEY ?? "").trim();
  const publicBaseUrl = (env.VITE_STORAGE_URL ?? "").trim().replace(/\/+$/, "");

  return {
    endpoint,
    bucket,
    accessKeyId,
    secretAccessKey,
    publicBaseUrl,
  };
}

function getS3Url(key: string, env: StorageRuntimeEnv): URL {
  const config = getStorageConfig(env);
  if (!config.endpoint || !/^https?:\/\//.test(config.endpoint))
    throw new Error("Invalid or missing S3_ENDPOINT protocol (http/https)");

  const baseUrl = config.endpoint.replace(/\/+$/, "");
  const cleanKey = key.replace(/^\/+/, "");
  return new URL(`${baseUrl}/${config.bucket}/${cleanKey}`);
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
  return Number.isFinite(parsed) && parsed > 0
    ? parsed
    : DEFAULT_MAX_SIZE_BYTES;
}

function getUploadExpiresSeconds(env: StorageRuntimeEnv) {
  const raw = env.STORAGE_IMAGE_UPLOAD_EXPIRES_SECONDS?.trim();
  if (!raw) return DEFAULT_EXPIRES_SECONDS;

  const parsed = Number(raw);
  return Number.isFinite(parsed) && parsed > 0
    ? parsed
    : DEFAULT_EXPIRES_SECONDS;
}

function getModerationModel(env: StorageRuntimeEnv) {
  return env.STORAGE_IMAGE_MODERATION_MODEL?.trim() || DEFAULT_MODERATION_MODEL;
}

function getModerationPrompt(env: StorageRuntimeEnv) {
  return (
    env.STORAGE_IMAGE_MODERATION_PROMPT?.trim() || DEFAULT_MODERATION_PROMPT
  );
}

function buildAllowedTypesErrorMessage(types: string[]) {
  return `Only ${types.join(", ")} are allowed`;
}

/**
 * Validates that all required S3 environment variables are present
 */
export const validateEnvEffect = (
  env: StorageRuntimeEnv,
): Effect.Effect<void, StorageConfigError> =>
  Effect.gen(function* () {
    const keys: (keyof StorageRuntimeEnv)[] = [
      "S3_ENDPOINT",
      "S3_BUCKET",
      "S3_ACCESS_KEY_ID",
      "S3_SECRET_ACCESS_KEY",
    ];
    for (const key of keys) {
      if (!env[key]) {
        yield* Effect.fail(
          new StorageConfigError({
            message: `Missing environment variable: ${key}`,
          }),
        );
      }
    }
  });

const getAwsClient = (env: StorageRuntimeEnv) => {
  const config = getStorageConfig(env);
  return new AwsClient({
    accessKeyId: config.accessKeyId,
    secretAccessKey: config.secretAccessKey,
    service: "s3",
    region: "auto",
  });
};

function sanitizePath(path?: string | null) {
  const value = (path ?? "").trim();
  if (!value) return "";
  return value.replace(/^\/+|\/+$/g, "").replace(/\.\./g, "");
}

function buildStorageKey(
  filename: string,
  path?: string | null,
  idempotencyKey?: string | null,
) {
  const cleanPath = sanitizePath(path);
  const cleanFilename = filename.replace(/\s+/g, "-").replace(/[^\w.-]/g, "");
  const cleanIdempotencyKey = (idempotencyKey ?? "").replace(/[^\w-]/g, "");
  const baseName = `${cleanIdempotencyKey || Date.now()}-${cleanFilename}`;
  return cleanPath ? `${cleanPath}/${baseName}` : baseName;
}

/**
 * The 5-token cap keeps the vision model's reply terse, so the prompt's "reply
 * with only 'safe' or 'unsafe'" actually holds. Anything that does not start
 * with "safe" — "unsafe", empty, a refusal — counts as a rejection.
 */
export function isSafeModerationOutput(output: unknown): boolean {
  const description = (output as { description?: unknown } | null)?.description;
  return String(description ?? "")
    .trim()
    .toLowerCase()
    .startsWith("safe");
}

export const moderateFileBytesEffect = (
  file: File,
  env: StorageRuntimeEnv,
): Effect.Effect<boolean, StorageUploadError> =>
  Effect.tryPromise({
    try: async () => {
      const buffer = await file.arrayBuffer();
      const response = await env.AI.run(getModerationModel(env), {
        prompt: getModerationPrompt(env),
        image: Array.from(new Uint8Array(buffer)),
        max_tokens: 5,
        temperature: 0.0,
      });
      return isSafeModerationOutput(response);
    },
    catch: (cause) =>
      new StorageUploadError({
        message:
          cause instanceof Error ? cause.message : "Falha na moderação de imagem",
        cause,
      }),
  });

export const putFileToStorageEffect = (
  file: File,
  key: string,
  env: StorageRuntimeEnv,
): Effect.Effect<string, StorageUploadError> =>
  Effect.tryPromise({
    try: async () => {
      const aws = getAwsClient(env);
      const uploadUrl = getS3Url(key, env);
      const res = await aws.fetch(uploadUrl.toString(), {
        method: "PUT",
        body: file,
        headers: {
          "Content-Type": file.type,
          "Cache-Control": "public, max-age=31536000, immutable",
        },
      });

      if (!res.ok) {
        const responseBody = await res.text().catch(() => "");
        const storageMessage = responseBody.match(
          /<Message>(.*?)<\/Message>/s,
        )?.[1];
        throw new Error(
          `Falha ao enviar imagem para o storage (${res.status})${storageMessage ? `: ${storageMessage}` : ""
          }`,
        );
      }

      const config = getStorageConfig(env);
      const cleanKey = key.replace(/^\/+/, "");
      if (config.publicBaseUrl) {
        return `${config.publicBaseUrl}/${cleanKey}`;
      }
      return `${getS3Url("", env).toString()}${cleanKey}`;
    },
    catch: (cause) =>
      new StorageUploadError({
        message: cause instanceof Error ? cause.message : "Falha ao enviar imagem",
        cause,
      }),
  });

const readImageFormDataEffect = (
  request: Request,
): Effect.Effect<
  {
    file: File;
    moderationFile: File;
    path: string;
    idempotencyKey: string;
  },
  StorageValidationError
> =>
  Effect.tryPromise({
    try: async () => {
      const formData = await request.formData();
      const file = formData.get("file");
      const moderationFile = formData.get("moderationFile");
      const path = formData.get("path");
      const idempotencyKey = formData.get("idempotencyKey");

      if (!(file instanceof File)) throw new Error("Missing file");

      return {
        file,
        moderationFile: moderationFile instanceof File ? moderationFile : file,
        path: typeof path === "string" ? path : "",
        idempotencyKey: typeof idempotencyKey === "string" ? idempotencyKey : "",
      };
    },
    catch: (cause) =>
      new StorageValidationError({
        message: cause instanceof Error ? cause.message : "Missing file",
      }),
  });

export const handleStorageUploadEffect = (
  request: Request,
  env: StorageRuntimeEnv,
): Effect.Effect<Response, never> =>
  Effect.gen(function* () {
    yield* validateEnvEffect(env);

    const { filename, contentType, size } = yield* Effect.tryPromise({
      try: () => request.json() as Promise<StorageUploadRequest>,
      catch: () =>
        new StorageValidationError({
          message: "Invalid JSON payload",
        }),
    });

    const allowedTypes = getAllowedTypes(env);
    const maxSize = getMaxSize(env);

    if (!allowedTypes.includes(contentType)) {
      yield* Effect.fail(
        new StorageValidationError({
          message: buildAllowedTypesErrorMessage(allowedTypes),
        }),
      );
    }

    if (size > maxSize) {
      yield* Effect.fail(
        new StorageValidationError({
          message: "File exceeds 10MB limit",
        }),
      );
    }

    const config = getStorageConfig(env);
    const aws = getAwsClient(env);
    const uploadUrl = getS3Url(filename, env);
    uploadUrl.searchParams.set(
      "X-Amz-Expires",
      String(getUploadExpiresSeconds(env)),
    );

    const signed = yield* Effect.tryPromise({
      try: () =>
        aws.sign(
          new Request(uploadUrl, {
            method: "PUT",
            headers: { "Content-Type": contentType },
          }),
          { aws: { signQuery: true } },
        ),
      catch: (cause) =>
        new StorageUploadError({
          message: "Upload failed",
          cause,
        }),
    });

    const cleanFilename = filename.replace(/^\/+/, "");
    const publicUrl = config.publicBaseUrl
      ? `${config.publicBaseUrl}/${cleanFilename}`
      : `${getS3Url("", env).toString()}${cleanFilename}`;

    return privateJsonResponse({
      uploadUrl: signed.url,
      key: filename,
      publicUrl,
    });
  }).pipe(
    Effect.catchTags({
      StorageValidationError: (err) =>
        Effect.succeed(
          privateJsonResponse({ error: err.message }, { status: 400 }),
        ),
      StorageConfigError: (err) =>
        Effect.succeed(
          privateJsonResponse({ error: err.message }, { status: 500 }),
        ),
      StorageUploadError: (err) =>
        Effect.succeed(
          privateJsonResponse({ error: err.message }, { status: 500 }),
        ),
    }),
    Effect.catch((err: unknown) =>
      Effect.succeed(
        privateJsonResponse(
          { error: err instanceof Error ? err.message : "Upload failed" },
          { status: 500 },
        ),
      ),
    ),
  );

export async function handleStorageUpload(
  request: Request,
  env: StorageRuntimeEnv,
): Promise<Response> {
  return Effect.runPromise(handleStorageUploadEffect(request, env));
}

export const handleStorageImagePreprocessEffect = (
  request: Request,
  env: StorageRuntimeEnv,
): Effect.Effect<Response, never> =>
  Effect.gen(function* () {
    yield* validateEnvEffect(env);

    const { file, moderationFile, path, idempotencyKey } =
      yield* readImageFormDataEffect(request);

    const allowedTypes = getAllowedTypes(env);
    const maxSize = getMaxSize(env);

    if (!allowedTypes.includes(file.type)) {
      yield* Effect.fail(
        new StorageValidationError({
          message: buildAllowedTypesErrorMessage(allowedTypes),
        }),
      );
    }

    if (file.size > maxSize) {
      yield* Effect.fail(
        new StorageValidationError({
          message: "File exceeds 10MB limit",
        }),
      );
    }

    if (!allowedTypes.includes(moderationFile.type)) {
      yield* Effect.fail(
        new StorageValidationError({
          message: buildAllowedTypesErrorMessage(allowedTypes),
        }),
      );
    }

    if (moderationFile.size > maxSize) {
      yield* Effect.fail(
        new StorageValidationError({
          message: "File exceeds 10MB limit",
        }),
      );
    }

    const approved = yield* moderateFileBytesEffect(moderationFile, env);
    if (!approved) {
      return privateJsonResponse({ approved: false });
    }

    const key = buildStorageKey(file.name, path, idempotencyKey);
    const publicUrl = yield* putFileToStorageEffect(file, key, env);
    return privateJsonResponse({ approved: true, key, publicUrl });
  }).pipe(
    Effect.catchTags({
      StorageValidationError: (err) =>
        Effect.succeed(
          privateJsonResponse({ error: err.message }, { status: 400 }),
        ),
      StorageConfigError: (err) =>
        Effect.succeed(
          privateJsonResponse({ error: err.message }, { status: 500 }),
        ),
      StorageUploadError: (err) =>
        Effect.succeed(
          privateJsonResponse({ error: err.message }, { status: 500 }),
        ),
    }),
    Effect.catch((err: unknown) =>
      Effect.succeed(
        privateJsonResponse(
          {
            error:
              err instanceof Error ? err.message : "Preprocessing failed",
          },
          { status: 500 },
        ),
      ),
    ),
  );

export async function handleStorageImagePreprocess(
  request: Request,
  env: StorageRuntimeEnv,
): Promise<Response> {
  return Effect.runPromise(handleStorageImagePreprocessEffect(request, env));
}
