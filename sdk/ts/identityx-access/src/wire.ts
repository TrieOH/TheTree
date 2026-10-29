import {
  defaultToFailure,
  defaultToSuccess,
  type DefaultFailureEnvelope,
  type DefaultSuccessEnvelope,
} from "@trieoh/envoy-fetch-ts";

// IdentityX's wire: the payload itself on success (nothing on a 204), lists
// as `{ items }`, and `{ code, message, fields }` on failure.

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === "object" && value !== null && !Array.isArray(value);

/** `toSuccess` for the fetch client: the body is the payload. */
export function toSuccess<T>(raw: unknown): DefaultSuccessEnvelope & { data: T } {
  return defaultToSuccess<T>({ module: "IdentityX", data: raw });
}

/** `toFailure` for the fetch client: the HTTP status, the error's code and message. */
export function toFailure(raw: unknown, status: number): DefaultFailureEnvelope {
  const error = isRecord(raw) && typeof raw.code === "string" && typeof raw.message === "string"
    ? { message: raw.message, error_id: raw.code }
    : {};
  return defaultToFailure({ module: "IdentityX", code: status, ...error }, status);
}

/** The items of a list payload (`{ items }`). */
export const itemsOf = <T>(data: unknown): T[] =>
  (data as { items?: T[] } | undefined)?.items ?? [];
