import {
  createDefaultFetchClient,
  type DefaultFetchResult,
} from "@trieoh/envoy-fetch-ts";
import {
  type ApiResponse,
  type AuthCallbacks,
  type AuthService,
  type AuthTokens,
  type TokenSubject,
} from "@trieoh/identityx-sdk-ts";
import {
  addActiveSpanEvent,
  getTraceparent,
  setActiveSpanAttributes,
  withSpan,
} from "../../tracing/browser";
import type {
  BffIntrospectResponse,
  IdentityXTransportLogEvent,
  ProxyHttpMethod,
  SerializableValue,
  ServerAuthResult,
  ServerProxyRequest,
  ServerProxyResult,
  ServerSessionSnapshot,
} from "./types";

/**
 * How an adapter reaches the BFF. `call` is the whole contract: one operation
 * name plus its input. Every transport — TanStack Start server functions, an
 * HTTP endpoint, a stub in a test — plugs in here.
 */
export interface BffTransport {
  call<T>(op: string, input?: unknown): Promise<T>;
  /** Operations this transport implements. Omitted = all of them. */
  has?(op: string): boolean;
}

export interface BffClientOptions {
  projectId?: string;
  log?: (event: IdentityXTransportLogEvent) => void;
  logSuccesses?: boolean;
}

/**
 * Structural twin of the framework adapters (`AuthProviderAdapter` in both the
 * React and Solid SDKs), so one implementation serves either provider.
 */
export interface BffAuthAdapter {
  createAuth(context: {
    callbacks: AuthCallbacks;
    defaultAuth: AuthService;
    getProfile(): TokenSubject | null;
    setProfile(profile: TokenSubject | null): void;
    setAuthenticated(authenticated: boolean): void;
  }): AuthService;
  restoreSession(): Promise<boolean | ServerSessionSnapshot>;
}

function isSerializable(value: unknown): value is SerializableValue {
  if (value === null) return true;
  if (["string", "number", "boolean"].includes(typeof value)) return true;
  if (Array.isArray(value)) return value.every(isSerializable);
  if (typeof value !== "object") return false;
  return Object.values(value).every(isSerializable);
}

const traceOperation = <T>(name: string, call: () => Promise<T>) =>
  withSpan(`bff:${name}`, call);

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

export function proxyResultToApiResponse<T>(result: ServerProxyResult): ApiResponse<T> {
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

/**
 * Framework-neutral IdentityX auth adapter over a {@link BffTransport}: it owns
 * the mapping between BFF results and the SDK's `AuthService`, and the local
 * profile/authenticated state, so no app has to re-declare it.
 */
export function createBffAuthAdapter(
  transport: BffTransport,
  options: BffClientOptions = {},
): BffAuthAdapter {
  const supports = (op: string) => transport.has?.(op) ?? true;

  return {
    restoreSession: () =>
      traceOperation("restoreSession", () => transport.call<ServerSessionSnapshot>("restore")),

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
        const started = performance.now();
        const result = await withSpan(`bff:${operation}`, async () => {
          setActiveSpanAttributes({
            "bff.operation": operation,
            "bff.target": "identityx",
            "bff.path": path.split("?", 1)[0] ?? path,
            "bff.method": method,
          });
          addActiveSpanEvent("bff.request.started");
          const traceparent = getTraceparent();
          const response = await transport.call<ServerProxyResult>("request", {
            path,
            target: "identityx",
            method,
            ...(body === undefined ? {} : { body }),
            ...(traceparent ? { headers: { traceparent } } : {}),
          });
          setActiveSpanAttributes({
            "bff.success": response.success,
            "bff.status": response.code,
            ...(response.error_id ? { "bff.error_id": response.error_id } : {}),
          });
          addActiveSpanEvent(
            response.success ? "bff.request.completed" : "bff.request.failed",
          );
          return response;
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
        isSetupDone: supports("isSetupDone")
          ? async () => authResponse<void>(await transport.call<ServerAuthResult>("isSetupDone"))
          : defaultAuth.isSetupDone,
        setup: supports("setup")
          ? async (email, password) => {
            const result = await traceOperation("setup", () =>
              transport.call<ServerAuthResult>("setup", { email, password }));
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
          const result = await traceOperation("login", () =>
            transport.call<ServerAuthResult>("login", { email, password }));
          const response = authResponse<AuthTokens>(result);
          if (result.success) {
            setProfile(result.profile ?? null);
            setAuthenticated(true);
            callbacks.onLogin?.(response);
          }
          return response;
        },
        register: supports("register")
          ? async (email, password) => {
            const result = await traceOperation("register", () =>
              transport.call<ServerAuthResult>("register", { email, password }));
            const response = authResponse<void>(result);
            if (result.success) callbacks.onRegister?.(response);
            return response;
          }
          : defaultAuth.register,
        loginWithProvider: supports("loginWithProvider")
          ? async (provider) => {
            const result = await traceOperation("loginWithProvider", () =>
              transport.call<ServerAuthResult & { data?: { url: string } }>(
                "loginWithProvider",
                { provider },
              ));
            const base = authResponse<{ url: string }>(result);
            return result.success
              ? { ...base, success: true, data: result.data! }
              : base;
          }
          : defaultAuth.loginWithProvider,
        completeProviderLogin: supports("completeProviderLogin")
          ? async (provider, code, state) => {
            const result = await traceOperation("completeProviderLogin", () =>
              transport.call<ServerAuthResult>("completeProviderLogin", {
                provider,
                code,
                state,
              }));
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
          const result = await traceOperation("logout", () =>
            transport.call<ServerAuthResult>("logout"));
          const response = authResponse<void>(result);
          setProfile(null);
          setAuthenticated(false);
          return response;
        },
        refresh: async () => {
          const result = await traceOperation("refresh", () =>
            transport.call<ServerAuthResult>("refresh"));
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
        introspect: supports("introspect")
          ? async (apiKey?: string) => {
            if (apiKey) return defaultAuth.introspect(apiKey);

            const result = await traceOperation("introspect", () =>
              transport.call<ServerProxyResult<BffIntrospectResponse>>(
                "introspect",
                apiKey ? { apiKey } : undefined,
              ));
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
        ...(supports("request") ? {
          sendForgotPassword: (email: string) =>
            proxyResponse<void>("sendForgotPassword", "/auth/forgot-password", "POST", {
              email,
              ...(options.projectId ? { project_id: options.projectId } : {}),
            }),
          resetPassword: (token: string, password: string) =>
            proxyResponse<void>("resetPassword", "/auth/reset-password", "POST", { token, password }),
          verifyEmail: async (token: string) => {
            const response = await proxyResponse<void>(
              "verifyEmail",
              "/auth/verify-email",
              "POST",
              { token },
            );
            return response.success
              ? response
              : {
                ...response,
                message:
                  "Este link de verificação expirou ou não é válido. Solicite um novo link para confirmar seu e-mail.",
              };
          },
          resendVerifyEmail: (email: string) =>
            proxyResponse<void>("resendVerifyEmail", "/auth/resend-verification", "POST", {
              email,
              ...(options.projectId ? { project_id: options.projectId } : {}),
            }),
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
          // These three look public in the SDK, but in BFF mode a
          // direct call would leave the browser for IdentityX with no
          // session (and usually fail CORS), so they go through the
          // proxy like everything else.
          getProfileByHandle: (handle: string) =>
            proxyResponse("getProfileByHandle", `/profiles/by-handle/${encodeURIComponent(handle)}`),
          getOAuthProviders: (overrideProjectId?: string) => {
            const projectId = overrideProjectId ?? options.projectId;
            const query = projectId
              ? `?project_id=${encodeURIComponent(projectId)}`
              : "";
            return proxyResponse("getOAuthProviders", `/auth/oauth-providers${query}`);
          },
          health: () => proxyResponse("health", "/health"),
        } : {}),
      };
    },
  };
}

/**
 * API transport for a BFF-backed app: every request is proxied through the
 * session, so the browser never holds a token.
 */
export function createBffProxyFetchers(
  transport: BffTransport,
  options: BffClientOptions = {},
  onSessionInvalid: () => void = () => undefined,
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

    const result = await withSpan("bff:apiRequest", async () => {
      setActiveSpanAttributes({
        "bff.operation": "apiRequest",
        "bff.target": "identityx",
        "bff.path": url.split("?", 1)[0] ?? url,
        "bff.method": method,
      });
      addActiveSpanEvent("bff.request.started");
      const traceparent = getTraceparent();
      if (traceparent) headers["traceparent"] = traceparent;
      const response = await transport.call<ServerProxyResult>("request", {
        path: url,
        method: method as ProxyHttpMethod,
        ...(body === undefined ? {} : { body }),
        ...(Object.keys(headers).length === 0 ? {} : { headers }),
      } satisfies ServerProxyRequest);
      setActiveSpanAttributes({
        "bff.success": response.success,
        "bff.status": response.code,
        ...(response.error_id ? { "bff.error_id": response.error_id } : {}),
      });
      addActiveSpanEvent(
        response.success ? "bff.request.completed" : "bff.request.failed",
      );
      return response;
    });
    if (result.code === 401) onSessionInvalid();
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

export type {
  ServerAuthResult,
  ServerProxyRequest,
  ServerProxyResult,
  ServerSessionSnapshot,
} from "./types";
