import { readFailure } from "@trieoh/api-client/wire";
import type { AuthTokens, TokenClaims, TokenSubject } from "@trieoh/identityx-sdk-ts";
import type {
  BffIntrospectResponse,
  ServerAuthResult,
  ServerOperationResult,
  ServerProxyRequest,
  ServerProxyResult,
  ServerSessionSnapshot,
  SerializableValue,
  IdentityXTransportLogEvent,
} from "./types";

/**
 * Server-managed IdentityX session, resolved per request. The BFF only knows
 * how to read and write tokens through this, so the runtime stays swap-
 * pable: TanStack Start's `useSession` today, a signed cookie on a Worker.
 */
export interface BffSession {
  read(): { tokens?: AuthTokens } | null;
  write(tokens: AuthTokens): void | Promise<void>;
  clear(): void | Promise<void>;
}

export interface IdentityXBffConfig {
  identityX: {
    baseURL: string;
    projectId?: string;
  };
  apiBaseURL: string;
  /**
   * How long the server waits for an upstream before giving up on a read, in
   * ms (default {@link DEFAULT_READ_TIMEOUT_MS}). Only idempotent reads are
   * cut short: a refresh or a write that is aborted may already have taken
   * effect upstream (the refresh token is single-use), so those run to
   * completion. Pass `0` to disable.
   */
  readTimeoutMs?: number;
  /**
   * Resolved per operation, never cached: each request must get its own
   * session, otherwise tokens leak across callers.
   */
  session: () => BffSession | Promise<BffSession>;
  /** Also per call: the proxy checks the origin of the incoming request. */
  request: () => Request;
  observability?: {
    log?: (event: IdentityXTransportLogEvent) => void;
    logSuccesses?: boolean;
  };
}

type TokenResolution =
  | { success: true; tokens: AuthTokens }
  | { success: false; error: ServerAuthResult; sessionInvalid: boolean };

/** What a Backend answered: the payload on success, else the error's code and message. */
interface Answer<T> {
  data?: T;
  message?: string;
  error_id?: string;
}

/**
 * Generous on purpose: the browser gives up on a read after a few seconds on
 * its own, so this only exists to stop a hung upstream from pinning the server
 * request (and its isolate) open indefinitely.
 */
export const DEFAULT_READ_TIMEOUT_MS = 20_000;

const isSafeMethod = (method: string) => method === "GET" || method === "HEAD";

const isTimeout = (error: unknown) =>
  error instanceof Error && (error.name === "TimeoutError" || error.name === "AbortError");

function joinURL(base: string, path: string): string {
  return `${base.replace(/\/$/, "")}/${path.replace(/^\//, "")}`;
}

function defaultBffLogger(event: IdentityXTransportLogEvent): void {
  const method = event.success ? "info" : "error";
  console[method]("[identityx-bff]", event);
}

function decodeClaims(token: string): TokenClaims | null {
  try {
    const payload = token.split(".")[1];
    if (!payload) return null;
    const base64 = payload.replace(/-/g, "+").replace(/_/g, "/");
    const padded = base64.padEnd(Math.ceil(base64.length / 4) * 4, "=");
    return JSON.parse(atob(padded)) as TokenClaims;
  } catch {
    return null;
  }
}

function getProfile(
  tokens?: AuthTokens,
  options: { allowExpired?: boolean } = {},
): TokenSubject | null {
  if (!tokens) return null;
  const claims = decodeClaims(tokens.access_token);
  if (
    !claims ||
    (!options.allowExpired && claims.exp * 1000 <= Date.now())
  )
    return null;
  return claims.subject;
}

function normalize<T>(response: Response, answer: Answer<T>): ServerProxyResult<T> {
  return {
    success: response.ok,
    code: response.status,
    ...(answer.data === undefined ? {} : { data: answer.data }),
    ...(answer.message ? { message: answer.message } : {}),
    ...(answer.error_id ? { error_id: answer.error_id } : {}),
  };
}

async function readAnswer<T>(response: Response): Promise<Answer<T>> {
  if (response.status === 204 || response.body === null) return {};

  try {
    const body: unknown = await response.json();
    if (response.ok) return { data: body as T };
    const failure = readFailure(body);
    return {
      ...(failure.message ? { message: failure.message } : {}),
      ...(failure.code ? { error_id: failure.code } : {}),
    };
  } catch {
    return { message: response.statusText };
  }
}

/**
 * In-flight refreshes, keyed by refresh token and shared across requests.
 *
 * The refresh token is single-use, so two concurrent requests arriving with the
 * same expired session must not both call `/auth/refresh`: the loser would be
 * logged out. Two requests that land in different isolates are separate
 * processes and still race — the next request then carries the winner's cookie.
 */
const refreshes = new Map<string, Promise<TokenResolution>>();

/**
 * The IdentityX server-managed session (BFF), with no framework imports. It owns
 * the IdentityX contract — endpoints, responses, refresh and the proxy guard —
 * and stays out of the storage and transport decisions, which the caller
 * injects.
 */
export function createIdentityXBff(config: IdentityXBffConfig) {
  const session = () => config.session();
  const readTimeoutMs = config.readTimeoutMs ?? DEFAULT_READ_TIMEOUT_MS;

  const observedFetch = async (
    operation: string,
    url: string,
    init?: RequestInit,
    { timeout = false }: { timeout?: boolean } = {},
  ): Promise<Response> => {
    const started = performance.now();
    const method = (init?.method ?? "GET").toUpperCase();
    const path = new URL(url).pathname;
    // Carry the caller's trace across the boundary. Without this the span
    // stops at the BFF and the IdentityX/API calls that follow — refresh
    // included — are unattributable.
    const traceparent = config.request().headers.get("traceparent");
    try {
      const headers = new Headers(init?.headers);
      if (traceparent && !headers.has("traceparent")) {
        headers.set("traceparent", traceparent);
      }

      const response = await fetch(url, {
        ...init,
        headers,
        ...(timeout && readTimeoutMs > 0 ? { signal: AbortSignal.timeout(readTimeoutMs) } : {}),
      });
      const failure = response.ok
        ? undefined
        : await readAnswer(response.clone());
      const event: IdentityXTransportLogEvent = {
        layer: "bff-server",
        operation,
        method,
        path,
        duration_ms: Math.round(performance.now() - started),
        success: response.ok,
        status: response.status,
        ...(traceparent ? { traceparent } : {}),
        ...(failure?.error_id ? { error_id: failure.error_id } : {}),
        ...(!response.ok ? { message: failure?.message ?? response.statusText } : {}),
      };
      if (config.observability?.logSuccesses || !response.ok) {
        (config.observability?.log ?? defaultBffLogger)(event);
      }
      return response;
    } catch (error) {
      const event: IdentityXTransportLogEvent = {
        layer: "bff-server",
        operation,
        method,
        path,
        duration_ms: Math.round(performance.now() - started),
        success: false,
        message: error instanceof Error ? error.message : "Network request failed",
        ...(traceparent ? { traceparent } : {}),
      };
      (config.observability?.log ?? defaultBffLogger)(event);
      throw error;
    }
  };

  const projectQuery = () => config.identityX.projectId
    ? `?project_id=${encodeURIComponent(config.identityX.projectId)}`
    : "";

  async function authenticate(
    path: string,
    init: RequestInit,
  ): Promise<ServerAuthResult> {
    const response = await observedFetch("authenticate", joinURL(config.identityX.baseURL, path), init);
    const answer = await readAnswer<AuthTokens>(response);
    const normalized = normalize(response, answer);
    if (!normalized.success || !answer.data) return normalized;

    const current = await session();
    await current.write(answer.data);
    return { ...normalized, profile: getProfile(answer.data) };
  }

  async function refreshTokens(): Promise<TokenResolution> {
    const current = await session();
    const refreshToken = current.read()?.tokens?.refresh_token;
    if (!refreshToken) {
      return {
        success: false,
        sessionInvalid: true,
        error: {
          success: false,
          code: 401,
          error_id: "REFRESH_TOKEN_MISSING",
          message: "Session expired",
        },
      };
    }
    if (refreshToken.split(".").length !== 3) {
      await current.clear();
      return {
        success: false,
        sessionInvalid: true,
        error: {
          success: false,
          code: 401,
          error_id: "REFRESH_TOKEN_MALFORMED",
          message: "Session expired",
        },
      };
    }

    let refresh = refreshes.get(refreshToken);
    if (!refresh) {
      refresh = (async (): Promise<TokenResolution> => {
        let response: Response;
        try {
          response = await observedFetch(
            "refresh",
            joinURL(config.identityX.baseURL, "/auth/refresh"),
            {
              method: "POST",
              headers: { "Refresh-Token": refreshToken },
            },
          );
        } catch {
          return {
            success: false,
            sessionInvalid: false,
            error: {
              success: false,
              code: 503,
              error_id: "AUTH_SERVICE_UNAVAILABLE",
              message: "Authentication service unavailable",
            },
          };
        }

        const answer = await readAnswer<AuthTokens>(response);
        if (!response.ok || !answer.data?.access_token) {
          const normalized = normalize(response, answer);
          const sessionInvalid = response.status === 401 || response.status === 403;
          return {
            success: false,
            sessionInvalid,
            error: response.ok
              ? {
                success: false,
                code: 502,
                error_id: "INVALID_REFRESH_RESPONSE",
                message: "Authentication service returned invalid tokens",
              }
              : normalized,
          };
        }

        return { success: true, tokens: answer.data };
      })();
      refreshes.set(refreshToken, refresh);
      void refresh.finally(() => refreshes.delete(refreshToken));
    }

    const resolution = await refresh;
    if (resolution.success) await current.write(resolution.tokens);
    return resolution;
  }

  async function validTokens(): Promise<TokenResolution> {
    const current = await session();
    const tokens = current.read()?.tokens;
    if (!tokens) {
      return {
        success: false,
        sessionInvalid: true,
        error: {
          success: false,
          code: 401,
          error_id: "SESSION_MISSING",
          message: "Session expired",
        },
      };
    }

    const accessExpiry = new Date(tokens.access_expires_at).getTime();
    if (Number.isFinite(accessExpiry) && accessExpiry > Date.now() + 30_000) {
      return { success: true, tokens };
    }
    return refreshTokens();
  }

  return {
    async isSetupDone(): Promise<ServerAuthResult> {
      const response = await observedFetch(
        "isSetupDone",
        joinURL(config.identityX.baseURL, "/auth/setup"),
        undefined,
        { timeout: true },
      );
      return normalize(response, await readAnswer(response));
    },

    async introspect(
      apiKey?: string,
    ): Promise<ServerProxyResult<BffIntrospectResponse>> {
      if (apiKey) {
        const response = await observedFetch(
          "introspectApiKey",
          joinURL(config.identityX.baseURL, "/auth/introspect"),
          { headers: { "X-API-KEY": apiKey } },
          { timeout: true },
        );
        return normalize(
          response,
          await readAnswer<BffIntrospectResponse>(response),
        );
      }

      const resolution = await validTokens();
      if (!resolution.success) return resolution.error;

      const response = await observedFetch(
        "introspect",
        joinURL(config.identityX.baseURL, "/auth/introspect"),
        {
          headers: {
            Authorization: `Bearer ${resolution.tokens.access_token}`,
          },
        },
        { timeout: true },
      );
      return normalize(
        response,
        await readAnswer<BffIntrospectResponse>(response),
      );
    },

    async setup(email: string, password: string): Promise<ServerAuthResult> {
      return authenticate(`/auth/setup${projectQuery()}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });
    },

    async login(email: string, password: string): Promise<ServerAuthResult> {
      return authenticate(`/auth/login${projectQuery()}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });
    },

    async register(email: string, password: string): Promise<ServerAuthResult> {
      const response = await observedFetch(
        "register",
        joinURL(config.identityX.baseURL, `/auth/register${projectQuery()}`),
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email, password }),
        },
      );
      return normalize(response, await readAnswer(response));
    },

    async loginWithProvider(
      provider: "github" | "google",
    ): Promise<ServerOperationResult<{ url: string }>> {
      const response = await observedFetch("loginWithProvider", joinURL(
        config.identityX.baseURL,
        `/auth/${provider}/connect${projectQuery()}`,
      ), undefined, { timeout: true });
      return normalize(response, await readAnswer<{ url: string }>(response));
    },

    async completeProviderLogin(
      provider: "github" | "google",
      code: string,
      state: string,
    ): Promise<ServerAuthResult> {
      const query = new URLSearchParams({ code, state });
      return authenticate(
        `/auth/${provider}/callback?${query}`,
        { method: "GET" },
      );
    },

    async logout(): Promise<ServerAuthResult> {
      const current = await session();
      const resolution = await validTokens();
      let result: ServerAuthResult = { success: true, code: 204 };

      if (resolution.success) {
        const response = await observedFetch("logout", joinURL(config.identityX.baseURL, "/auth/logout"), {
          method: "POST",
          headers: {
            Authorization: `Bearer ${resolution.tokens.access_token}`,
            "Refresh-Token": resolution.tokens.refresh_token,
          },
        });
        result = normalize(response, await readAnswer(response));
      }

      await current.clear();
      return result;
    },

    async refresh(): Promise<ServerAuthResult> {
      const resolution = await refreshTokens();
      return resolution.success
        ? {
          success: true,
          code: 200,
          profile: getProfile(resolution.tokens),
        }
        : resolution.error;
    },

    async restore(): Promise<ServerSessionSnapshot> {
      const current = await session();
      const tokens = current.read()?.tokens;
      const resolution = await validTokens();
      if (resolution.success) {
        const profile = getProfile(resolution.tokens);
        return { isAuthenticated: !!profile, profile };
      }
      if (!resolution.sessionInvalid) {
        const profile = getProfile(tokens, { allowExpired: true });
        return { isAuthenticated: !!profile, profile };
      }
      return { isAuthenticated: false, profile: null };
    },

    async request<T extends SerializableValue = SerializableValue>(
      input: ServerProxyRequest,
    ): Promise<ServerProxyResult<T>> {
      if (!input.path.startsWith("/") || input.path.startsWith("//")) {
        return { success: false, code: 400, message: "Invalid service path" };
      }

      const request = config.request();
      const method = (input.method ?? "GET").toUpperCase();
      if (!["GET", "HEAD", "OPTIONS"].includes(method)) {
        const origin = request.headers.get("Origin");
        if (origin && origin !== new URL(request.url).origin) {
          return { success: false, code: 403, message: "Invalid request origin" };
        }
      }

      const anonymousAuthRequest =
        input.target === "identityx" &&
        method === "POST" &&
        [
          "/auth/verify-email",
          "/auth/resend-verification",
          "/auth/forgot-password",
          "/auth/reset-password",
        ].includes(input.path);
      const resolution = await validTokens();
      if (
        !resolution.success &&
        !anonymousAuthRequest &&
        (method !== "GET" || !resolution.sessionInvalid)
      ) {
        return resolution.error;
      }

      const headers = new Headers(input.headers);
      if (resolution.success) {
        headers.set("Authorization", `Bearer ${resolution.tokens.access_token}`);
      }
      if (input.body !== undefined && !headers.has("Content-Type")) {
        headers.set("Content-Type", "application/json");
      }

      const targetBaseURL = input.target === "identityx"
        ? config.identityX.baseURL
        : config.apiBaseURL;
      const operation = input.target === "identityx"
        ? "identityXRequest"
        : "apiRequest";
      const safe = isSafeMethod(method);
      const send = () =>
        observedFetch(operation, joinURL(targetBaseURL, input.path), {
          method,
          headers,
          body: input.body === undefined ? undefined : JSON.stringify(input.body),
        }, { timeout: safe });

      try {
        let response = await send();
        if (response.status === 401 && resolution.success) {
          const refreshed = await refreshTokens();
          if (!refreshed.success) return refreshed.error;
          headers.set("Authorization", `Bearer ${refreshed.tokens.access_token}`);
          response = await send();
        }
        return normalize(response, await readAnswer<T>(response));
      } catch (error) {
        // A read that could not be completed is reported as an upstream
        // failure the client knows how to show, rather than a server error.
        // Writes still throw: whether they landed is not ours to guess.
        if (!safe) throw error;
        return isTimeout(error)
          ? { success: false, code: 504, error_id: "UPSTREAM_TIMEOUT", message: "Upstream service timed out" }
          : { success: false, code: 503, error_id: "UPSTREAM_UNAVAILABLE", message: "Upstream service unavailable" };
      }
    },
  };
}

export type IdentityXBff = ReturnType<typeof createIdentityXBff>;
export type {
  IdentityXOAuthProvider,
  ProxyHttpMethod,
  ServerAuthResult,
  ServerOperationResult,
  ServerProxyRequest,
  ServerProxyResult,
  ServerSessionSnapshot,
} from "./types";
