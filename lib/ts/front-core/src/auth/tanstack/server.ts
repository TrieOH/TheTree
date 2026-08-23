import { useSession, getRequest } from "@tanstack/react-start/server";
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

interface IdentityXSessionData {
  tokens: AuthTokens;
}

type TokenResolution =
  | { success: true; tokens: AuthTokens }
  | { success: false; error: ServerAuthResult; sessionInvalid: boolean };

interface ApiEnvelope<T> {
  code?: number;
  data?: T;
  message?: string;
  error_id?: string;
  error?: {
    code?: string;
    message?: string;
  };
  trace?: string[];
}

export interface TanStackIdentityXBffConfig {
  identityX: {
    baseURL: string;
    projectId?: string;
  };
  session: {
    password: string;
    name?: string;
    maxAge?: number;
    secure?: boolean;
  };
  apiBaseURL: string;
  observability?: {
    log?: (event: IdentityXTransportLogEvent) => void;
    logSuccesses?: boolean;
  };
}

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

function normalize<T>(response: Response, envelope: ApiEnvelope<T>): ServerProxyResult<T> {
  const code = envelope.code ?? response.status;
  const message = envelope.message ?? envelope.error?.message;
  const errorId = envelope.error_id ?? envelope.error?.code;
  return {
    success: response.ok && code >= 200 && code < 300,
    code,
    ...(envelope.data === undefined ? {} : { data: envelope.data }),
    ...(message ? { message } : {}),
    ...(errorId ? { error_id: errorId } : {}),
    ...(envelope.trace ? { trace: envelope.trace } : {}),
  };
}

async function readEnvelope<T>(response: Response): Promise<ApiEnvelope<T>> {
  if (response.status === 204 || response.body === null) {
    return { code: response.status };
  }

  try {
    return await response.json() as ApiEnvelope<T>;
  } catch {
    return { code: response.status, message: response.statusText };
  }
}

export function createTanStackIdentityXBff(config: TanStackIdentityXBffConfig) {
  if (config.session.password.length < 32) {
    throw new Error("IdentityX session password must contain at least 32 characters");
  }

  const sessionConfig = {
    name: config.session.name ?? "trieoh-auth",
    password: config.session.password,
    maxAge: config.session.maxAge ?? 60 * 60 * 24 * 30,
    cookie: {
      httpOnly: true,
      secure: config.session.secure ?? true,
      sameSite: "lax" as const,
      path: "/",
    },
  };

  const session = () => useSession<IdentityXSessionData>(sessionConfig);
  const refreshes = new Map<string, Promise<TokenResolution>>();

  const observedFetch = async (
    operation: string,
    url: string,
    init?: RequestInit,
  ): Promise<Response> => {
    const started = performance.now();
    const method = (init?.method ?? "GET").toUpperCase();
    const path = new URL(url).pathname;
    try {
      const response = await fetch(url, init);
      const failureEnvelope = response.ok
        ? undefined
        : await readEnvelope(response.clone());
      const event: IdentityXTransportLogEvent = {
        layer: "bff-server",
        operation,
        method,
        path,
        duration_ms: Math.round(performance.now() - started),
        success: response.ok,
        status: response.status,
        ...(failureEnvelope?.error_id || failureEnvelope?.error?.code
          ? { error_id: failureEnvelope.error_id ?? failureEnvelope.error?.code }
          : {}),
        ...(!response.ok
          ? { message: failureEnvelope?.message ?? failureEnvelope?.error?.message ?? response.statusText }
          : {}),
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
    const envelope = await readEnvelope<AuthTokens>(response);
    const normalized = normalize(response, envelope);
    if (!normalized.success || !envelope.data) return normalized;

    const current = await session();
    await current.update({ tokens: envelope.data });
    return { ...normalized, profile: getProfile(envelope.data) };
  }

  async function refreshTokens(): Promise<TokenResolution> {
    const current = await session();
    const refreshToken = current.data.tokens?.refresh_token;
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

        const envelope = await readEnvelope<AuthTokens>(response);
        if (!response.ok || !envelope.data?.access_token) {
          const normalized = normalize(response, envelope);
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

        return { success: true, tokens: envelope.data };
      })();
      refreshes.set(refreshToken, refresh);
      void refresh.finally(() => refreshes.delete(refreshToken));
    }

    const resolution = await refresh;
    if (resolution.success) await current.update({ tokens: resolution.tokens });
    else if (resolution.sessionInvalid) await current.clear();
    return resolution;
  }

  async function validTokens(): Promise<TokenResolution> {
    const current = await session();
    const tokens = current.data.tokens;
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
      const response = await observedFetch("isSetupDone", joinURL(config.identityX.baseURL, "/auth/setup"));
      return normalize(response, await readEnvelope(response));
    },

    async introspect(
      apiKey?: string,
    ): Promise<ServerProxyResult<BffIntrospectResponse>> {
      if (apiKey) {
        const response = await observedFetch(
          "introspectApiKey",
          joinURL(config.identityX.baseURL, "/auth/introspect"),
          { headers: { "X-API-KEY": apiKey } },
        );
        return normalize(
          response,
          await readEnvelope<BffIntrospectResponse>(response),
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
      );
      return normalize(
        response,
        await readEnvelope<BffIntrospectResponse>(response),
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
      return normalize(response, await readEnvelope(response));
    },

    async loginWithProvider(
      provider: "github" | "google",
    ): Promise<ServerOperationResult<{ url: string }>> {
      const response = await observedFetch("loginWithProvider", joinURL(
        config.identityX.baseURL,
        `/auth/${provider}/connect${projectQuery()}`,
      ));
      return normalize(response, await readEnvelope<{ url: string }>(response));
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
        result = normalize(response, await readEnvelope(response));
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
      const resolution = await validTokens();
      if (resolution.success) {
        const profile = getProfile(resolution.tokens);
        return { isAuthenticated: !!profile, profile };
      }
      if (!resolution.sessionInvalid) {
        const profile = getProfile(current.data.tokens, { allowExpired: true });
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

      const request = getRequest();
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
      const response = await observedFetch(
        input.target === "identityx" ? "identityXRequest" : "apiRequest",
        joinURL(targetBaseURL, input.path), {
        method,
        headers,
        body: input.body === undefined ? undefined : JSON.stringify(input.body),
      });
      return normalize(response, await readEnvelope<T>(response));
    },
  };
}

export type TanStackIdentityXBff = ReturnType<typeof createTanStackIdentityXBff>;
export type {
  IdentityXOAuthProvider,
  ProxyHttpMethod,
  ServerAuthResult,
  ServerOperationResult,
  ServerProxyRequest,
  ServerProxyResult,
  ServerSessionSnapshot,
} from "./types";
