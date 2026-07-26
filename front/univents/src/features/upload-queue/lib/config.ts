import { env } from "@/env";

export const uploadQueueConfig = {
  maxRetries: env.VITE_UPLOAD_MAX_RETRIES,
  baseDelayMs: env.VITE_UPLOAD_RETRY_BASE_DELAY_MS,
  maxDelayMs: env.VITE_UPLOAD_RETRY_MAX_DELAY_MS,
};

export function getRetryDelay(retryCount: number) {
  const exponentialDelay =
    uploadQueueConfig.baseDelayMs * 2 ** Math.max(retryCount - 1, 0);
  const cappedDelay = Math.min(exponentialDelay, uploadQueueConfig.maxDelayMs);
  const jitter = Math.floor(
    Math.random() * Math.max(Math.round(cappedDelay * 0.2), 1),
  );
  return cappedDelay + jitter;
}
