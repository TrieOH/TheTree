/**
 * The Backends are moving off the `fun.Response` envelope, one at a time:
 * plain JSON payloads on success, and on failure the envelope's `error`
 * object on its own (`{ code, message, fields }`), later RFC 9457
 * `application/problem+json`. Until every Backend has moved, the fetch stack keeps reading the
 * envelope (`DefaultFetchResult`, `ApiError.envelope`), and these helpers
 * turn whichever shape a Backend sent into it, so callers never see which.
 * Drop them, and read the plain shapes directly, once no Backend sends the
 * envelope.
 */

/** The envelope fields the fetch stack and the BFF read. */
export interface WireEnvelope<T = unknown> {
  code: number
  message?: string
  module?: string
  timestamp?: string
  data?: T
  error_id?: string
  error?: {
    code?: string
    message?: string
    fields?: { field: string; message: string }[]
  }
  trace?: string[]
}

/** An RFC 9457 problem, with the extension members the Backends send. */
interface Problem {
  type?: string
  status: number
  title: string
  detail?: string
  instance?: string
  /** Machine-readable error code (`NOT_FOUND`, `VALIDATION_ERROR`, …). */
  code?: string
  /** The request fields at fault, on validation failures. */
  errors?: { field: string; message: string }[]
}

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === "object" && value !== null && !Array.isArray(value)

/**
 * True for a `fun.Response` envelope: a numeric `code` next to `data` or
 * `module`. No payload in any spec has those fields.
 */
export const isEnvelope = (body: unknown): body is WireEnvelope =>
  isRecord(body) && typeof body.code === "number" && ("data" in body || "module" in body)

/** The envelope's `error` object sent on its own: a string `code` and a `message`. */
interface BareError {
  code: string
  message: string
  fields?: { field: string; message: string }[]
}

const isBareError = (body: unknown): body is BareError =>
  isRecord(body) && typeof body.code === "string" && typeof body.message === "string"

const isProblem = (body: unknown, contentType?: string | null): body is Problem =>
  isRecord(body) &&
  (contentType?.includes("application/problem+json") ||
    (typeof body.status === "number" && typeof body.title === "string"))

/**
 * The envelope for a response body of either shape. `body` is the parsed
 * JSON, or `undefined` when there was none (a 204, say).
 */
export function toEnvelope<T = unknown>(
  body: unknown,
  status: number,
  contentType?: string | null,
): WireEnvelope<T> {
  if (isEnvelope(body)) return body as WireEnvelope<T>
  const ok = status >= 200 && status < 300
  if (!ok && isBareError(body)) {
    return { code: status, message: body.message, error_id: body.code, error: body }
  }
  if (!ok && isProblem(body, contentType)) {
    const message = body.detail || body.title
    return {
      code: status,
      message,
      ...(body.code ? { error_id: body.code } : {}),
      error: {
        ...(body.code ? { code: body.code } : {}),
        message,
        ...(body.errors ? { fields: body.errors } : {}),
      },
    }
  }
  if (!ok) return { code: status }
  return { code: status, ...(body === undefined ? {} : { data: body as T }) }
}

// Response cannot carry a body on these statuses; the envelope keeps the real one in `code`.
const NULL_BODY_STATUSES = new Set([101, 204, 205, 304])

/**
 * Rewrites a response of either shape into an envelope response, for the
 * fetch clients that only read the envelope. A response already enveloped
 * is returned as is.
 */
export async function toEnvelopeResponse(response: Response): Promise<Response> {
  const text = await response.clone().text()
  let body: unknown
  try {
    body = text ? JSON.parse(text) : undefined
  } catch {
    return response
  }
  if (isEnvelope(body)) return response

  const envelope = toEnvelope(body, response.status, response.headers.get("content-type"))
  const headers = new Headers(response.headers)
  headers.set("content-type", "application/json")
  headers.delete("content-length")
  return new Response(JSON.stringify(envelope), {
    status: NULL_BODY_STATUSES.has(response.status) ? 200 : response.status,
    statusText: response.statusText,
    headers,
  })
}

/**
 * A `fetch` that answers with envelope responses whichever shape the
 * Backend sent. Pass it as the fetch clients' transport
 * (`clientConfig.fetch`, or an envoy-fetch `adapter`).
 */
export const envelopeFetch =
  // The global is looked up per call, so tracing that patches it later still applies.
  (fetchImpl: typeof fetch = (input, init) => fetch(input, init)): typeof fetch =>
  async (input, init) =>
    toEnvelopeResponse(await fetchImpl(input, init))

/**
 * The items of a list payload: `{ items }` from a Backend that moved off the
 * envelope, or the bare array one still on it sends.
 */
export const itemsOf = <T>(data: unknown): T[] =>
  Array.isArray(data) ? (data as T[]) : ((data as { items?: T[] } | undefined)?.items ?? [])
