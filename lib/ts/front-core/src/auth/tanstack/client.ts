import { createDefaultFetchClient, type DefaultFetchResult } from "@trieoh/envoy-fetch-ts";
import type { AuthProviderAdapter, AuthService } from "@trieoh/identityx-sdk-ts/react";
import { createFetcher, type ApiResponse } from "@trieoh/identityx-sdk-ts";
import type { AuthTokens } from "@trieoh/identityx-sdk-ts";
import type {
  BffIntrospectResponse,
  SerializableValue,
  IdentityXOAuthProvider,
  IdentityXTransportLogEvent,
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
    { provider: IdentityXOAuthProvider; code: string; state: string },
    ServerAuthResult
  >;
  logout: ServerFunction<void, ServerAuthResult>;
  refresh: ServerFunction<void, ServerAuthResult>;
  restore: ServerFunction<void, ServerSessionSnapshot>;
  introspect?: ServerFunction<
    { apiKey?: string } | undefined,
    ServerProxyResult<BffIntrospectResponse>
  >;
  request?: ServerFunction<ServerProxyRequest, ServerProxyResult>;
}

export interface TanStackIdentityXClientOptions {
  mode?: "bff" | "direct";
  apiBaseURL?: string;
  authBaseURL?: string;
  projectId?: string;
  log?: (event: IdentityXTransportLogEvent) => void;
  logSuccesses?: boolean;
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
  options: TanStackIdentityXClientOptions = {},
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
      const proxyResponse = async <T>(
        operation: string,
        path: string,
        method: ProxyHttpMethod = "GET",
        body?: SerializableValue,
      ): Promise<ApiResponse<T>> => {
        if (!functions.request) {
          throw new Error(`IdentityX BFF request transport is not configured for ${operation}`);
        }
        const started = performance.now();
        const result = await functions.request({
          data: {
            path,
            target: "identityx",
            method,
            ...(body === undefined ? {} : { body }),
          },
        });
        const event: IdentityXTransportLogEvent = {
          layer: "bff-client",
          operation,
          method,
          path,
          duration_ms: Math.round(performance.now() - started),
          success: result.success,
          status: result.code,
          ...(result.error_id ? { error_id: result.error_id } : {}),
          ...(result.message ? { message: result.message } : {}),
        };
        if (options.logSuccesses || !result.success) {
          (options.log ?? defaultTransportLogger)(event);
        }
        return proxyResultToApiResponse<T>(result);
      };
      const projectPath = (suffix: string, overrideProjectId?: string) => {
        const projectId = overrideProjectId ?? options.projectId;
        return projectId ? `/projects/${projectId}${suffix}` : suffix;
      };

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
          ? async (provider, code, state) => {
            const result = await functions.completeProviderLogin!({ data: { provider, code, state } });
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
            if (apiKey) return defaultAuth.introspect(apiKey);

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
        ...(functions.request ? {
          resendVerifyEmail: () => proxyResponse<void>("resendVerifyEmail", "/account/verify/resend", "POST"),
          getProjectProfile: (actorId: string, projectId?: string) =>
            proxyResponse("getProjectProfile", projectPath(`/actors/${actorId}/profile`, projectId)),
          upsertProjectProfile: (actorId: string, data: { profile: Record<string, unknown> }, projectId?: string) =>
            proxyResponse("upsertProjectProfile", projectPath(`/actors/${actorId}/profile`, projectId), "PUT", data as SerializableValue),
          getPlatformProfile: (actorId: string) =>
            proxyResponse("getPlatformProfile", `/actors/${actorId}/profile`),
          upsertPlatformProfile: (actorId: string, data: { profile: Record<string, unknown> }) =>
            proxyResponse("upsertPlatformProfile", `/actors/${actorId}/profile`, "PUT", data as SerializableValue),
          getProfileSchema: (projectId?: string) =>
            proxyResponse("getProfileSchema", projectPath("/profile-schema", projectId)),
          upsertProfileSchema: (data: { schema: Record<string, unknown>; active: boolean }, projectId?: string) =>
            proxyResponse("upsertProfileSchema", projectPath("/profile-schema", projectId), "PUT", data as SerializableValue),
          getActorProfile: (actorId: string, projectId?: string) =>
            proxyResponse("getActorProfile", projectPath(`/actors/${actorId}/profile`, projectId)),
          upsertActorProfile: (actorId: string, data: { profile: Record<string, unknown> }, projectId?: string) =>
            proxyResponse("upsertActorProfile", projectPath(`/actors/${actorId}/profile`, projectId), "PUT", data as SerializableValue),
        } : {}),
      };
    },
  };
}

export function createTanStackIdentityXIntegration(
  functions: IdentityXServerFunctions & {
    request: ServerFunction<ServerProxyRequest, ServerProxyResult>;
  },
  options: TanStackIdentityXClientOptions = {},
) {
  if (options.mode === "direct") {
    if (!options.apiBaseURL) {
      throw new Error("IdentityX direct transport requires apiBaseURL");
    }
    const authFetcher = createFetcher({
      baseURL: options.apiBaseURL,
      authBaseURL: options.authBaseURL,
    });
    return {
      mode: "direct" as const,
      authAdapter: undefined,
      authFetcher,
      authQueryFetcher: async <T>(path: string): Promise<T> => {
        const result = await authFetcher.get<T>(path);
        if (!result.success) throw result;
        return result.data;
      },
    };
  }

  const fetchers = createTanStackServerProxyFetchers(functions.request, options);
  return {
    mode: "bff" as const,
    authAdapter: createTanStackIdentityXAuthProviderAdapter(functions, options),
    ...fetchers,
  };
}

function proxyResultToApiResponse<T>(result: ServerProxyResult): ApiResponse<T> {
  const base = {
    module: "identityx-bff",
    message: result.message ?? (result.success ? "OK" : "Request failed"),
    timestamp: new Date().toISOString(),
    code: result.code,
  };
  return result.success
    ? { ...base, success: true, data: result.data as T }
    : {
        ...base,
        success: false,
        error_id: result.error_id ?? "IDENTITYX_BFF_ERROR",
        ...(result.trace ? { trace: result.trace } : {}),
      };
}

function defaultTransportLogger(event: IdentityXTransportLogEvent): void {
  const method = event.success ? "info" : "error";
  console[method]("[identityx-transport]", event);
}

export function createTanStackServerProxyFetchers(
  proxy: ServerFunction<ServerProxyRequest, ServerProxyResult>,
  options: TanStackIdentityXClientOptions = {},
) {
  const adapter = async (url: string, init?: RequestInit): Promise<Response> => {
    const started = performance.now();
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
    const event: IdentityXTransportLogEvent = {
      layer: "bff-client",
      operation: "apiRequest",
      method,
      path: url.split("?", 1)[0] ?? url,
      duration_ms: Math.round(performance.now() - started),
      success: result.success,
      status: result.code,
      ...(result.error_id ? { error_id: result.error_id } : {}),
      ...(result.message ? { message: result.message } : {}),
    };
    if (options.logSuccesses || !result.success) {
      (options.log ?? defaultTransportLogger)(event);
    }
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
