/** The failure an ApiError carries: the HTTP status and the error's message. */
type ErrorEnvelope = {
  code?: unknown;
  message?: unknown;
};

/** Returns the most useful user-facing message exposed by local or API errors. */
export function getErrorMessage(error: unknown, fallback: string): string {
  if (typeof error === "string" && error.trim()) return error;
  if (!error || typeof error !== "object") return fallback;

  const candidate = error as {
    message?: unknown;
    envelope?: ErrorEnvelope;
  };
  for (const message of [candidate.envelope?.message, candidate.message]) {
    if (typeof message === "string" && message.trim()) return message;
  }

  return fallback;
}

export function isVerifiedEmailRequiredError(error: unknown): boolean {
  if (!error || typeof error !== "object") return false;

  const candidate = error as {
    message?: unknown;
    envelope?: ErrorEnvelope;
  };
  const message = getErrorMessage(error, "").toLowerCase();

  return (
    (candidate.envelope?.code === 401 || candidate.envelope?.code === 403) &&
    message.includes("verified email required")
  );
}
