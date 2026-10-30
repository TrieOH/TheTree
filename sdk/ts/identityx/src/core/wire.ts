import {
  defaultToFailure,
  defaultToSuccess,
  type DefaultFailureEnvelope,
  type DefaultSuccessEnvelope,
} from "@trieoh/envoy-fetch-ts";

// IdentityX's wire: the payload itself on success (nothing on a 204), lists
// as `{ items }`, and an RFC 9457 problem on failure.

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === "object" && value !== null && !Array.isArray(value);

/** `toSuccess` for the fetch client: the body is the payload. */
export function toSuccess<T>(raw: unknown): DefaultSuccessEnvelope & { data: T } {
  return defaultToSuccess<T>({ module: "IdentityX", data: raw });
}

/** `toFailure` for the fetch client: the HTTP status and the problem's detail (else its title). */
export function toFailure(raw: unknown, status: number): DefaultFailureEnvelope {
  const message = isRecord(raw)
    ? [raw.detail, raw.title].find((m): m is string => typeof m === "string" && m !== "")
    : undefined;
  return defaultToFailure({ module: "IdentityX", code: status, ...(message ? { message } : {}) }, status);
}

/** The items of a list payload (`{ items }`). */
export const itemsOf = <T>(data: unknown): T[] =>
  (data as { items?: T[] } | undefined)?.items ?? [];
