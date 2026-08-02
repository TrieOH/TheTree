import { createDefaultFetchClient, type DefaultFetchResult } from "@trieoh/envoy-fetch-ts";
import type { AuthProviderAdapter, AuthService } from "@trieoh/identityx-sdk-ts/react";
import type { ApiResponse } from "@trieoh/identityx-sdk-ts";
import type { AuthTokens } from "@trieoh/identityx-sdk-ts";
import type {
  BffIntrospectResponse,
  SerializableValue,
  IdentityXOAuthProvider,
  ProxyHttpMethod,
  ServerAuthResult,
  ServerOperationResult,
  ServerProxyRequest,
  ServerProxyResult,
  ServerSessionSnapshot,
} from "./types";

function isSerializable(value: unknown): value is SerializableValue {
  if (value === null) return true;
  if (["string", "number", "boolean"].includes(typeof value)) return true;
  if (Array.isArray(value)) return value.every(isSerializable);
  if (typeof value !== "object") return false;
  return Object.values(value).every(isSerializable);
}

type ServerFunction<TInput, TResult> = (options: { data: TInput }) => Promise<TResult>;

export interface IdentityXServerFunctions {
  isSetupDone?: ServerFunction<void, ServerAuthResult>;
  setup?: ServerFunction<{ email: string; password: string }, ServerAuthResult>;
  login: ServerFunction<{ email: string; password: string }, ServerAuthResult>;
  register?: ServerFunction<{ email: string; password: string }, ServerAuthResult>;
  loginWithProvider?: ServerFunction<
    { provider: IdentityXOAuthProvider },
    ServerOperationResult<{ url: string }>
  >;
  completeProviderLogin?: ServerFunction<
    { provider: IdentityXOAuthProvider; code: string },
    ServerAuthResult
  >;
  logout: ServerFunction<void, ServerAuthResult>;
  refresh: ServerFunction<void, ServerAuthResult>;
  restore: ServerFunction<void, ServerSessionSnapshot>;
  introspect?: ServerFunction<
    { apiKey?: string } | undefined,
    ServerProxyResult<BffIntrospectResponse>
  >;
}

function authResponse<T>(result: ServerAuthResult): ApiResponse<T> {
  const base = {
    module: "identityx-bff",
    message: result.message ?? (result.success ? "OK" : "Authentication failed"),
    timestamp: new Date().toISOString(),
    code: result.code,
  };
  if (result.success) {
    return { ...base, success: true, data: undefined as T };
  }
  return {
    ...base,
    success: false,
    error_id: result.error_id ?? "IDENTITYX_BFF_ERROR",
    ...(result.trace ? { trace: result.trace } : {}),
  };
}

export function createTanStackIdentityXAuthProviderAdapter(
  functions: IdentityXServerFunctions,
): AuthProviderAdapter {
  return {
    restoreSession: () => functions.restore({ data: undefined }),
    createAuth({
      callbacks,
      defaultAuth,
      getProfile,
      setProfile,
      setAuthenticated,
    }): AuthService {
      return {
        ...defaultAuth,
        isSetupDone: functions.isSetupDone
          ? async () => authResponse<void>(await functions.isSetupDone!({ data: undefined }))
          : defaultAuth.isSetupDone,
        setup: functions.setup
          ? async (email, password) => {
              const result = await functions.setup!({ data: { email, password } });
              const response = authResponse<AuthTokens>(result);
              if (result.success) {
                setProfile(result.profile ?? null);
                setAuthenticated(true);
                callbacks.onSetup?.(response);
              }
              return response;
            }
          : defaultAuth.setup,
        profile: getProfile,
        login: async (email, password) => {
          const result = await functions.login({ data: { email, password } });
          const response = authResponse<AuthTokens>(result);
          if (result.success) {
            setProfile(result.profile ?? null);
            setAuthenticated(true);
            callbacks.onLogin?.(response);
          }
          return response;
        },
        register: functions.register
          ? async (email, password) => {
              const result = await functions.register!({ data: { email, password } });
              const response = authResponse<void>(result);
              if (result.success) callbacks.onRegister?.(response);
              return response;
            }
          : defaultAuth.register,
        loginWithProvider: functions.loginWithProvider
          ? async (provider) => {
              const result = await functions.loginWithProvider!({ data: { provider } });
              const base = authResponse<{ url: string }>(result);
              return result.success
                ? { ...base, success: true, data: result.data! }
                : base;
            }
          : defaultAuth.loginWithProvider,
        completeProviderLogin: functions.completeProviderLogin
          ? async (provider, code) => {
              const result = await functions.completeProviderLogin!({ data: { provider, code } });
              const response = authResponse<AuthTokens>(result);
              if (result.success) {
                setProfile(result.profile ?? null);
                setAuthenticated(true);
                callbacks.onLogin?.(response);
              }
              return response;
            }
          : defaultAuth.completeProviderLogin,
        logout: async () => {
          const result = await functions.logout({ data: undefined });
          const response = authResponse<void>(result);
          setProfile(null);
          setAuthenticated(false);
          return response;
        },
        refresh: async () => {
          const result = await functions.refresh({ data: undefined });
          const response = authResponse<AuthTokens>(result);
          if (result.success) {
            setProfile(result.profile ?? getProfile());
            setAuthenticated(true);
            callbacks.onRefresh?.(response);
          } else if (result.code === 401 || result.code === 403) {
            setProfile(null);
            setAuthenticated(false);
          }
          return response;
        },
        introspect: functions.introspect
          ? async (apiKey?: string) => {
              const result = await functions.introspect!({ data: apiKey ? { apiKey } : undefined });
              const base = {
                module: "identityx-bff",
                message: result.message ?? (result.success ? "OK" : "Introspect failed"),
                timestamp: new Date().toISOString(),
                code: result.code,
              };
              if (result.success) {
                return { ...base, success: true, data: result.data! };
              }
              return {
                ...base,
                success: false,
                error_id: result.error_id ?? "INTROSPECT_ERROR",
                ...(result.trace ? { trace: result.trace } : {}),
              };
            }
          : defaultAuth.introspect,
      };
    },
  };
}

export function createTanStackServerProxyFetchers(
  proxy: ServerFunction<ServerProxyRequest, ServerProxyResult>,
) {
  const adapter = async (url: string, init?: RequestInit): Promise<Response> => {
    const headers = Object.fromEntries(new Headers(init?.headers).entries());
    let body: SerializableValue | undefined;
    if (typeof init?.body === "string" && init.body.length > 0) {
      try {
        const parsed: unknown = JSON.parse(init.body);
        if (!isSerializable(parsed)) throw new Error("Request body is not JSON serializable");
        body = parsed;
      } catch {
        body = init.body;
      }
    }

    const method = (init?.method ?? "GET").toUpperCase();
    if (!["GET", "POST", "PUT", "PATCH", "DELETE"].includes(method)) {
      throw new Error(`Unsupported proxy method: ${method}`);
    }

    const result = await proxy({
      data: {
        path: url,
        method: method as ProxyHttpMethod,
        ...(body === undefined ? {} : { body }),
        ...(Object.keys(headers).length === 0 ? {} : { headers }),
      },
    });
    const payload = {
      module: "server-proxy",
      message: result.message ?? (result.success ? "OK" : "Request failed"),
      timestamp: new Date().toISOString(),
      code: result.code,
      ...(result.data === undefined ? {} : { data: result.data }),
      ...(result.success
        ? {}
        : { error_id: result.error_id ?? "SERVER_PROXY_ERROR" }),
      ...(result.trace ? { trace: result.trace } : {}),
    };
    // Response.json cannot construct a response with a body for null-body
    // statuses such as 204. The proxy transport still needs a JSON envelope
    // for the fetch client, so use 200 while preserving the upstream code in
    // the payload.
    const transportStatus = [204, 205, 304].includes(result.code)
      ? 200
      : result.code;
    return Response.json(payload, { status: transportStatus });
  };

  const client = createDefaultFetchClient({ adapter });
  return {
    authFetcher: client,
    authQueryFetcher: async <T>(path: string): Promise<T> => {
      const result: DefaultFetchResult<T> = await client.get<T>(path);
      if (!result.success) throw result;
      return result.data;
    },
  };
}

export type { ServerAuthResult, ServerProxyRequest, ServerProxyResult, ServerSessionSnapshot } from "./types";
