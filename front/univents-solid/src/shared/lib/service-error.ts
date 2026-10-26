/**
 * Turns whatever a failed request threw into something the UI can act on.
 *
 * Two shapes reach the query layer: the fetch client's `FetchClientError` and
 * the app's own `QueryError`. Both carry an `envelope` with the HTTP-ish `code`
 * (0 for "never reached the server", 408 for a client-side timeout) and an
 * `error_id`. Everything else (a `TypeError` from `fetch`, a bare `Error`) is
 * treated as a network failure.
 */

export type ServiceName = "api" | "identityx";

export type FailureKind =
  | "offline"
  | "timeout"
  | "unavailable"
  | "unauthorized"
  | "forbidden"
  | "not-found"
  | "client"
  | "unknown";

export interface ServiceFailure {
  kind: FailureKind;
  status?: number;
  errorId?: string;
  /** Worth retrying: the server never answered, or answered with a 5xx. */
  transient: boolean;
  /** Short, user-facing reason in pt-BR. */
  message: string;
}

const MESSAGES: Record<FailureKind, string> = {
  offline: "Você parece estar sem conexão com a internet.",
  timeout: "O servidor demorou demais para responder.",
  unavailable: "O serviço está temporariamente indisponível.",
  unauthorized: "Sua sessão expirou. Entre novamente para continuar.",
  forbidden: "Você não tem permissão para acessar isto.",
  "not-found": "Não encontramos o que você procurou.",
  client: "A solicitação não pôde ser concluída.",
  unknown: "Ocorreu um erro inesperado.",
};

/** Raised locally, without touching the network, while a service's circuit is open. */
export const CIRCUIT_OPEN_ERROR_ID = "CIRCUIT_OPEN";

interface ErrorLike {
  name?: unknown;
  message?: unknown;
  code?: unknown;
  error_id?: unknown;
  envelope?: { code?: unknown; error_id?: unknown; message?: unknown };
}

/** What Chromium, Firefox and Safari put in the TypeError `fetch` rejects with. */
const NETWORK_FAILURE_MESSAGE =
  /failed to fetch|networkerror|network request failed|load failed|fetch failed/i;

const isNumber = (value: unknown): value is number =>
  typeof value === "number" && Number.isFinite(value);

const isString = (value: unknown): value is string =>
  typeof value === "string" && value.length > 0;

function kindFor(status: number | undefined, errorId: string | undefined): FailureKind {
  if (errorId === "CLIENT_OFFLINE") return "offline";
  if (errorId === "CLIENT_TIMEOUT") return "timeout";
  if (errorId === CIRCUIT_OPEN_ERROR_ID) return "unavailable";
  if (errorId?.startsWith("CLIENT_")) return "unavailable";
  if (errorId === "AUTH_SERVICE_UNAVAILABLE") return "unavailable";

  if (status === undefined) return "unknown";
  if (status === 0 || status === 502 || status === 503 || status === 504) return "unavailable";
  if (status === 408) return "timeout";
  if (status === 401) return "unauthorized";
  if (status === 403) return "forbidden";
  if (status === 404) return "not-found";
  if (status === 429 || status >= 500) return "unavailable";
  if (status >= 400) return "client";
  return "unknown";
}

export function classifyError(error: unknown): ServiceFailure {
  if (error === null || error === undefined) {
    return { kind: "unknown", transient: false, message: MESSAGES.unknown };
  }

  const value = (typeof error === "object" ? error : { message: String(error) }) as ErrorLike;
  const status = [value.envelope?.code, value.code].find(isNumber);
  const errorId = [value.envelope?.error_id, value.error_id].find(isString);

  // `fetch` rejects with a bare TypeError when the request never completes. A
  // TypeError from our own code ("x is not a function") is a bug, not an outage,
  // so it only counts when it carries the browsers' network-failure wording.
  const hasNoAnswer = status === undefined && errorId === undefined;
  const isNetworkTypeError =
    value.name === "TypeError" &&
    typeof value.message === "string" &&
    NETWORK_FAILURE_MESSAGE.test(value.message);
  const isTimeoutLike = value.name === "AbortError" || value.name === "TimeoutError";

  const kind = hasNoAnswer && isNetworkTypeError
    ? "unavailable"
    : hasNoAnswer && isTimeoutLike
      ? "timeout"
      : kindFor(status, errorId);

  return {
    kind,
    ...(status === undefined ? {} : { status }),
    ...(errorId === undefined ? {} : { errorId }),
    transient:
      kind === "offline" ||
      kind === "timeout" ||
      kind === "unavailable",
    message: MESSAGES[kind],
  };
}

/**
 * The code to report when no server answered, in the envelope shape the fetch
 * client already uses, so errors raised locally classify like remote ones.
 */
export function unavailableEnvelope(errorId: string, message?: string) {
  return {
    success: false as const,
    code: errorId === "CLIENT_TIMEOUT" ? 408 : 503,
    error_id: errorId,
    message: message ?? MESSAGES[errorId === "CLIENT_TIMEOUT" ? "timeout" : "unavailable"],
  };
}
