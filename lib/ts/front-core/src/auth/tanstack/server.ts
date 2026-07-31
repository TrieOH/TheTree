import { useSession, getRequest } from "@tanstack/react-start/server";
import type { AuthTokens, TokenClaims, TokenSubject } from "@trieoh/identityx-sdk-ts";
import type {
  ServerAuthResult,
  ServerOperationResult,
  ServerProxyRequest,
  ServerProxyResult,
  ServerSessionSnapshot,
  SerializableValue,
} from "./types";

interface IdentityXSessionData {
  tokens: AuthTokens;
}

interface ApiEnvelope<T> {
  code?: number;
  data?: T;
  message?: string;
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
}

function joinURL(base: string, path: string): string {
  return `${base.replace(/\/$/, "")}/${path.replace(/^\//, "")}`;
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

function getProfile(tokens?: AuthTokens): TokenSubject | null {
  if (!tokens) return null;
  const claims = decodeClaims(tokens.access_token);
  if (!claims || claims.exp * 1000 <= Date.now()) return null;
  return claims.subject;
}

function normalize<T>(response: Response, envelope: ApiEnvelope<T>): ServerProxyResult<T> {
  const code = envelope.code ?? response.status;
  return {
    success: response.ok && code >= 200 && code < 300,
    code,
    ...(envelope.data === undefined ? {} : { data: envelope.data }),
    ...(envelope.message ? { message: envelope.message } : {}),
    ...(envelope.trace ? { trace: envelope.trace } : {}),
  };
}

async function readEnvelope<T>(response: Response): Promise<ApiEnvelope<T>> {
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

  const projectQuery = () => config.identityX.projectId
    ? `?project_id=${encodeURIComponent(config.identityX.projectId)}`
    : "";

  async function authenticate(
    path: string,
    init: RequestInit,
  ): Promise<ServerAuthResult> {
    const response = await fetch(joinURL(config.identityX.baseURL, path), init);
    const envelope = await readEnvelope<AuthTokens>(response);
    const normalized = normalize(response, envelope);
    if (!normalized.success || !envelope.data) return normalized;

    const current = await session();
    await current.update({ tokens: envelope.data });
    return { ...normalized, profile: getProfile(envelope.data) };
  }

  async function refreshTokens(): Promise<AuthTokens | null> {
    const current = await session();
    const refreshToken = current.data.tokens?.refresh_token;
    if (!refreshToken) return null;

    const response = await fetch(joinURL(config.identityX.baseURL, "/auth/refresh"), {
      method: "POST",
      headers: { "Refresh-Token": refreshToken },
    });
    const envelope = await readEnvelope<AuthTokens>(response);
    if (!response.ok || !envelope.data?.access_token) {
      if (response.status >= 400 && response.status < 500) await current.clear();
      return null;
    }

    await current.update({ tokens: envelope.data });
    return envelope.data;
  }

  async function validTokens(): Promise<AuthTokens | null> {
    const current = await session();
    const tokens = current.data.tokens;
    if (!tokens) return null;

    const accessExpiry = new Date(tokens.access_expires_at).getTime();
    if (Number.isFinite(accessExpiry) && accessExpiry > Date.now() + 30_000) {
      return tokens;
    }
    return refreshTokens();
  }

  return {
    async isSetupDone(): Promise<ServerAuthResult> {
      const response = await fetch(joinURL(config.identityX.baseURL, "/auth/setup"));
      return normalize(response, await readEnvelope(response));
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
      const response = await fetch(
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
      const response = await fetch(joinURL(
        config.identityX.baseURL,
        `/auth/${provider}/connect${projectQuery()}`,
      ));
      return normalize(response, await readEnvelope<{ url: string }>(response));
    },

    async completeProviderLogin(
      provider: "github" | "google",
      code: string,
    ): Promise<ServerAuthResult> {
      return authenticate(
        `/auth/${provider}/callback?code=${encodeURIComponent(code)}`,
        { method: "GET" },
      );
    },

    async logout(): Promise<ServerAuthResult> {
      const current = await session();
      const accessToken = (await validTokens())?.access_token;
      let result: ServerAuthResult = { success: true, code: 204 };

      if (accessToken) {
        const response = await fetch(joinURL(config.identityX.baseURL, "/auth/logout"), {
          method: "POST",
          headers: { Authorization: `Bearer ${accessToken}` },
        });
        result = normalize(response, await readEnvelope(response));
      }

      await current.clear();
      return result;
    },

    async refresh(): Promise<ServerAuthResult> {
      const tokens = await refreshTokens();
      return tokens
        ? { success: true, code: 200, profile: getProfile(tokens) }
        : { success: false, code: 401, message: "Session expired" };
    },

    async restore(): Promise<ServerSessionSnapshot> {
      const tokens = await validTokens();
      const profile = getProfile(tokens ?? undefined);
      return { isAuthenticated: !!tokens && !!profile, profile };
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

      const tokens = await validTokens();
      if (!tokens) return { success: false, code: 401, message: "Session expired" };

      const headers = new Headers(input.headers);
      headers.set("Authorization", `Bearer ${tokens.access_token}`);
      if (input.body !== undefined && !headers.has("Content-Type")) {
        headers.set("Content-Type", "application/json");
      }

      const response = await fetch(joinURL(config.apiBaseURL, input.path), {
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
