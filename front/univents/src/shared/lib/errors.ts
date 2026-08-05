type ErrorEnvelope = {
  message?: unknown;
  error?: { message?: unknown } | unknown;
};

/** Returns the most useful user-facing message exposed by local or API errors. */
export function getErrorMessage(error: unknown, fallback: string): string {
  if (typeof error === "string" && error.trim()) return error;
  if (!error || typeof error !== "object") return fallback;

  const candidate = error as {
    message?: unknown;
    envelope?: ErrorEnvelope;
  };
  const nestedMessage =
    candidate.envelope?.error &&
    typeof candidate.envelope.error === "object" &&
    "message" in candidate.envelope.error
      ? candidate.envelope.error.message
      : undefined;

  for (const message of [
    candidate.envelope?.message,
    nestedMessage,
    candidate.message,
  ]) {
    if (typeof message === "string" && message.trim()) return message;
  }

  return fallback;
}
