export const uploadQueueConfig = {
  maxRetries: Number(import.meta.env.VITE_UPLOAD_MAX_RETRIES ?? 5),
  baseDelayMs: Number(import.meta.env.VITE_UPLOAD_RETRY_BASE_DELAY_MS ?? 1000),
  maxDelayMs: Number(import.meta.env.VITE_UPLOAD_RETRY_MAX_DELAY_MS ?? 30000),
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
