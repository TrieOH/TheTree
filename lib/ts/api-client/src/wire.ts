/**
 * Readers for the Backends' wire format, as envoy-fetch's `toSuccess` /
 * `toFailure`: the payload itself on success (nothing on a 204), lists as
 * `{ items }`, and on failure `{ code, message, fields }`, or the RFC 9457
 * problem errors are moving to (ADR-0016). They fill envoy-fetch's
 * `DefaultFetchResult`, whose `code` is the HTTP status on a failure and
 * whose `error_id` is the error's code.
 */
import {
  defaultToFailure,
  defaultToSuccess,
  type DefaultFailureEnvelope,
  type DefaultSuccessEnvelope,
} from "@trieoh/envoy-fetch-ts"

/** A field the request got wrong, as a failure names it. */
export interface FieldError {
  field: string
  message: string
}

/** What a failure body says, whichever error shape it came in. */
export interface WireFailure {
  /** Machine-readable error code (`NOT_FOUND`, `VALIDATION_ERROR`, …). */
  code?: string
  message?: string
  fields?: FieldError[]
}

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === "object" && value !== null && !Array.isArray(value)

/** Reads a failure body: the error object, or an RFC 9457 problem. */
export function readFailure(body: unknown): WireFailure {
  if (!isRecord(body)) return {}
  if (typeof body.code === "string" && typeof body.message === "string") {
    return body as WireFailure
  }
  if (typeof body.title === "string") {
    const message = typeof body.detail === "string" && body.detail ? body.detail : body.title
    return {
      ...(typeof body.code === "string" ? { code: body.code } : {}),
      message,
      ...(Array.isArray(body.errors) ? { fields: body.errors as FieldError[] } : {}),
    }
  }
  return {}
}

/** `toSuccess` for envoy-fetch: the body is the payload. */
export const toSuccess = <T>(raw: unknown): DefaultSuccessEnvelope & { data: T } =>
  defaultToSuccess<T>({ data: raw })

/** `toFailure` for envoy-fetch: the HTTP status, the error's code and message. */
export function toFailure(raw: unknown, status: number): DefaultFailureEnvelope {
  const failure = readFailure(raw)
  return defaultToFailure(
    {
      code: status,
      ...(failure.message ? { message: failure.message } : {}),
      ...(failure.code ? { error_id: failure.code } : {}),
    },
    status,
  )
}

/** The items of a list payload (`{ items }`). */
export const itemsOf = <T>(data: unknown): T[] =>
  (data as { items?: T[] } | undefined)?.items ?? []
