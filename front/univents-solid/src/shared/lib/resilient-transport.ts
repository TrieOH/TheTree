import type { ApiTransport } from "@trieoh/api-client";
import type {
  BffTransport,
  ServerProxyRequest,
  ServerProxyResult,
} from "@trieoh/front-core/auth/bff/client";
import {
  classifyError,
  unavailableEnvelope,
  type ServiceFailure,
  type ServiceName,
} from "./service-error";
import { CircuitOpenError, serviceHealth, type ServiceHealth } from "./service-health";

/**
 * How long a read may take before the page stops waiting for it. These are
 * deadlines on the caller, not request cancellations: the request keeps going,
 * so a session cookie the server rotates mid-flight (single-use refresh token)
 * still lands instead of being dropped with an aborted response.
 */
export const READ_DEADLINE_MS = 8_000;
/**
 * Session operations gate the whole app (`AuthProvider` holds the first paint
 * until the session settles), so a hung IdentityX must not hold it for long.
 */
export const SESSION_DEADLINE_MS = 8_000;

class DeadlineError extends Error {
  readonly envelope = unavailableEnvelope("CLIENT_TIMEOUT");

  constructor(ms: number) {
    super(`Request timed out after ${ms}ms`);
    this.name = "TimeoutError";
  }
}

export function withDeadline<T>(promise: Promise<T>, ms: number): Promise<T> {
  let timer: ReturnType<typeof setTimeout>;
  const deadline = new Promise<never>((_, reject) => {
    timer = setTimeout(() => reject(new DeadlineError(ms)), ms);
  });
  // The request outlives its deadline; swallow its late failure instead of
  // leaving an unhandled rejection behind.
  promise.catch(() => undefined);
  return Promise.race([promise, deadline]).finally(() => clearTimeout(timer));
}

const isRead = (method: string | undefined) =>
  method === undefined || method === "GET" || method === "HEAD";

type ApiRequestOptions = Parameters<ApiTransport>[1];

/**
 * Wraps the generated API clients' transport with a read deadline, the
 * service's circuit breaker and health reporting. Writes are never cut short or
 * short-circuited: they are explicit user actions and may have already landed.
 */
export function withApiResilience(
  transport: ApiTransport,
  service: ServiceName,
  health: ServiceHealth = serviceHealth,
): ApiTransport {
  return (async (url: string, options: ApiRequestOptions) => {
    const read = isRead(options.method);
    if (read) health.guard(service);

    try {
      const response = await (read
        ? withDeadline(transport(url, options), READ_DEADLINE_MS)
        : transport(url, options));
      health.success(service);
      return response;
    } catch (error) {
      if (error instanceof CircuitOpenError) throw error;
      report(health, service, classifyError(error));
      throw error;
    }
  }) as ApiTransport;
}

function report(health: ServiceHealth, service: ServiceName, failure: ServiceFailure) {
  // A 404 or 422 is the service answering: it is alive, so it counts as healthy.
  if (failure.transient) health.failure(service, failure);
  else health.success(service);
}

/** What went wrong in a result envelope, or `null` if it succeeded. */
function failureOf(result: unknown): ServiceFailure | null {
  if (typeof result !== "object" || result === null) return null;
  if ((result as { success?: unknown }).success !== false) return null;
  return classifyError(result);
}

const SESSION_OPS = new Set(["restore", "refresh", "introspect", "isSetupDone"]);

const bffServiceFor = (op: string, input: unknown): ServiceName =>
  op === "request" && (input as ServerProxyRequest | undefined)?.target !== "identityx"
    ? "api"
    : "identityx";

function bffDeadlineFor(op: string, input: unknown): number | undefined {
  if (SESSION_OPS.has(op)) return SESSION_DEADLINE_MS;
  if (op === "request" && isRead((input as ServerProxyRequest | undefined)?.method)) {
    return READ_DEADLINE_MS;
  }
  return undefined;
}

function failedRequest(errorId: string): ServerProxyResult {
  const { success, code, error_id, message } = unavailableEnvelope(errorId);
  return { success, code, error_id, message };
}

/**
 * The BFF counterpart of {@link withApiResilience}. Authenticated traffic goes
 * through `call("request", …)`, so this is where authenticated reads get their
 * deadline and where IdentityX health is observed.
 *
 * A `request` that cannot be completed resolves to a failure result instead of
 * rejecting, which is the shape the proxy adapter already turns into a proper
 * `ApiError` envelope. Session operations still reject: the SDK handles those.
 */
export function withBffResilience(
  transport: BffTransport,
  health: ServiceHealth = serviceHealth,
): BffTransport {
  return {
    ...(transport.has ? { has: (op: string) => transport.has?.(op) ?? true } : {}),

    async call<T>(op: string, input?: unknown): Promise<T> {
      const isRequest = op === "request";
      const service = bffServiceFor(op, input);
      const deadline = bffDeadlineFor(op, input);

      if (isRequest && deadline !== undefined && health.isOpen(service)) {
        return failedRequest("CIRCUIT_OPEN") as T;
      }

      let result: T;
      try {
        const call = transport.call<T>(op, input);
        result = deadline === undefined ? await call : await withDeadline(call, deadline);
      } catch (error) {
        report(health, service, classifyError(error));
        if (!isRequest) throw error;
        return failedRequest(
          classifyError(error).kind === "timeout" ? "CLIENT_TIMEOUT" : "CLIENT_REQUEST_FAILED",
        ) as T;
      }

      const failure = failureOf(result);
      // A refresh that could not reach IdentityX fails an API request without
      // the API being at fault.
      const culprit =
        failure?.errorId === "AUTH_SERVICE_UNAVAILABLE" ? "identityx" : service;
      report(health, culprit, failure ?? { kind: "unknown", transient: false, message: "" });
      return result;
    },
  };
}
