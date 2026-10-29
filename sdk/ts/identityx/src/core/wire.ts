import {
  defaultToFailure,
  defaultToSuccess,
  type DefaultFailureEnvelope,
  type DefaultSuccessEnvelope,
} from "@trieoh/envoy-fetch-ts";

// IdentityX is moving off the `fun.Response` envelope: the payload itself on
// success, `{ code, message, fields }` on failure, lists as `{ items }`. These
// read both shapes, so the SDK works on either side of the move.

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === "object" && value !== null && !Array.isArray(value);

/** A `fun.Response` envelope: a numeric `code` next to `data` or `module`. */
export const isEnvelope = (raw: unknown): raw is { code: number; data?: unknown; message?: string } =>
  isRecord(raw) && typeof raw.code === "number" && ("data" in raw || "module" in raw);

/** `toSuccess` for the fetch client: the envelope's data, or the body itself. */
export function toSuccess<T>(raw: unknown): DefaultSuccessEnvelope & { data: T } {
  return defaultToSuccess<T>(isEnvelope(raw) ? raw : { module: "IdentityX", data: raw });
}

/** `toFailure` for the fetch client: the envelope, or a bare error object. */
export function toFailure(raw: unknown, status: number): DefaultFailureEnvelope {
  if (!isEnvelope(raw) && isRecord(raw) && typeof raw.code === "string" && typeof raw.message === "string") {
    return defaultToFailure({ module: "IdentityX", code: status, message: raw.message, error_id: raw.code }, status);
  }
  return defaultToFailure(raw, status);
}

/** The items of a list payload: `{ items }`, or the bare array of the envelope days. */
export const itemsOf = <T>(data: unknown): T[] =>
  Array.isArray(data) ? (data as T[]) : ((data as { items?: T[] } | undefined)?.items ?? []);
