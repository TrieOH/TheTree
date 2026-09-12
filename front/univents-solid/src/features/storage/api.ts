import { Data, Effect, Schedule } from "effect";

export const ALLOWED_IMAGE_TYPES = [
  "image/png",
  "image/jpeg",
  "image/webp",
] as const;

export const MAX_IMAGE_SIZE_BYTES = 10 * 1024 * 1024; // 10 MB

export class InvalidFileTypeError extends Data.TaggedError("InvalidFileTypeError")<{
  readonly fileType: string;
  readonly allowedTypes: readonly string[];
}> {
  readonly message = "Use uma imagem PNG, JPG ou WebP.";
}

export class FileSizeExceededError extends Data.TaggedError("FileSizeExceededError")<{
  readonly size: number;
  readonly maxSize: number;
}> {
  readonly message = "A imagem deve ter no máximo 10 MB.";
}

export class StorageUploadNetworkError extends Data.TaggedError("StorageUploadNetworkError")<{
  readonly cause: unknown;
  readonly message: string;
}> {}

export class StorageModerationError extends Data.TaggedError("StorageModerationError")<{
  readonly reason: string;
  readonly message: string;
}> {}

export type StorageError =
  | InvalidFileTypeError
  | FileSizeExceededError
  | StorageUploadNetworkError
  | StorageModerationError;

/**
 * Default retry schedule for transient network failures:
 * Exponential backoff (starting at 150ms) with jitter, limited to 2 retries (3 attempts total).
 */
export const uploadRetrySchedule = Schedule.exponential("150 millis").pipe(
  Schedule.jittered,
  Schedule.upTo({ times: 2 }),
);

export const validateImageFile = (
  file: File,
): Effect.Effect<File, InvalidFileTypeError | FileSizeExceededError> =>
  Effect.gen(function* () {
    if (
      !ALLOWED_IMAGE_TYPES.includes(
        file.type as (typeof ALLOWED_IMAGE_TYPES)[number],
      )
    ) {
      return yield* Effect.fail(
        new InvalidFileTypeError({
          fileType: file.type,
          allowedTypes: ALLOWED_IMAGE_TYPES,
        }),
      );
    }
    if (file.size > MAX_IMAGE_SIZE_BYTES) {
      return yield* Effect.fail(
        new FileSizeExceededError({
          size: file.size,
          maxSize: MAX_IMAGE_SIZE_BYTES,
        }),
      );
    }
    return file;
  });

export const validateImageFileSync = (
  file: File,
): { ok: true; file: File } | { ok: false; error: string } =>
  Effect.runSync(
    Effect.match(validateImageFile(file), {
      onFailure: (err) => ({ ok: false, error: err.message }),
      onSuccess: (validFile) => ({ ok: true, file: validFile }),
    }),
  );

export interface UploadOptions {
  readonly retrySchedule?: Schedule.Schedule<unknown, unknown>;
}

export const uploadProfileImageEffect = (
  file: File,
  field: "pfpUrl" | "bannerUrl",
  options?: UploadOptions,
): Effect.Effect<string, StorageError> =>
  Effect.gen(function* () {
    yield* validateImageFile(file);

    const uploadAttempt = Effect.gen(function* () {
      const data = new FormData();
      data.append("file", file);
      data.append("path", "profiles/images");
      data.append("idempotencyKey", crypto.randomUUID());
      data.append("field", field);

      const response = yield* Effect.tryPromise({
        try: () =>
          fetch("/storage/image/preprocess", {
            method: "POST",
            body: data,
          }),
        catch: (cause) =>
          new StorageUploadNetworkError({
            cause,
            message: "Falha na conexão ao enviar a imagem.",
          }),
      });

      const result = yield* Effect.tryPromise({
        try: () =>
          response.json().catch(() => ({})) as Promise<{
            publicUrl?: string;
            error?: string;
            approved?: boolean;
          }>,
        catch: (cause) =>
          new StorageUploadNetworkError({
            cause,
            message: "Falha ao processar a resposta do servidor de imagens.",
          }),
      });

      if (!response.ok || !result.publicUrl) {
        return yield* Effect.fail(
          new StorageModerationError({
            reason: result.error ?? "A imagem não foi aprovada.",
            message: result.error ?? "A imagem não foi aprovada.",
          }),
        );
      }

      return result.publicUrl;
    });

    const schedule = options?.retrySchedule ?? uploadRetrySchedule;

    return yield* Effect.retry(uploadAttempt, {
      schedule,
      while: (err) => err._tag === "StorageUploadNetworkError",
    });
  });

export interface UploadBatchResult {
  readonly uploaded: Partial<Record<"pfpUrl" | "bannerUrl", string>>;
  readonly failed: Array<"foto" | "banner">;
}

export const uploadProfileImagesBatchEffect = (
  uploads: ReadonlyArray<readonly ["pfpUrl" | "bannerUrl", File]>,
  options?: UploadOptions,
): Effect.Effect<UploadBatchResult> =>
  Effect.gen(function* () {
    const tasks = uploads.map(([field, file]) =>
      uploadProfileImageEffect(file, field, options).pipe(
        Effect.match({
          onFailure: () => ({ success: false as const, field }),
          onSuccess: (url) => ({ success: true as const, field, url }),
        }),
      ),
    );

    const outcomes = yield* Effect.all(tasks, { concurrency: "unbounded" });

    const uploaded: Partial<Record<"pfpUrl" | "bannerUrl", string>> = {};
    const failed: Array<"foto" | "banner"> = [];

    for (const outcome of outcomes) {
      if (outcome.success) {
        uploaded[outcome.field] = outcome.url;
      } else {
        failed.push(outcome.field === "pfpUrl" ? "foto" : "banner");
      }
    }

    return { uploaded, failed };
  });

export const uploadProfileImagesBatch = (
  uploads: ReadonlyArray<readonly ["pfpUrl" | "bannerUrl", File]>,
  options?: UploadOptions,
): Promise<UploadBatchResult> =>
  Effect.runPromise(uploadProfileImagesBatchEffect(uploads, options));

export async function uploadProfileImage(
  file: File,
  field: "pfpUrl" | "bannerUrl",
  options?: UploadOptions,
): Promise<string> {
  return Effect.runPromise(uploadProfileImageEffect(file, field, options));
}
