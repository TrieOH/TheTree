import type { AuthProviderAdapter } from "@trieoh/identityx-sdk-ts-react";
import {
  createBffAuthAdapter,
  createBffProxyFetchers,
  type BffClientOptions,
  type BffTransport,
} from "@trieoh/front-core/auth/bff/client";
import type {
  BffIntrospectResponse,
  ServerAuthResult,
  ServerOperationResult,
  ServerProxyRequest,
  ServerProxyResult,
  ServerSessionSnapshot,
} from "./types";

type ServerFunction<TInput, TResult> = (options: { data: TInput }) => Promise<TResult>;

export interface IdentityXServerFunctions {
  isSetupDone?: ServerFunction<void, ServerAuthResult>;
  setup?: ServerFunction<{ email: string; password: string }, ServerAuthResult>;
  login: ServerFunction<{ email: string; password: string }, ServerAuthResult>;
  register?: ServerFunction<{ email: string; password: string }, ServerAuthResult>;
  loginWithProvider?: ServerFunction<
    { provider: "github" | "google" },
    ServerOperationResult<{ url: string }>
  >;
  completeProviderLogin?: ServerFunction<
    { provider: "github" | "google"; code: string; state: string },
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

/** Kept as a named alias so consumers have one thing to import. */
export type TanStackIdentityXClientOptions = BffClientOptions;

/**
 * Start-side transport: `{ data }`-shaped server functions behind the neutral
 * `call(op, input)` contract, so the adapter and the fetchers stay shared with
 * every other runtime.
 */
function serverFnTransport(
  functions: Partial<Record<string, ServerFunction<unknown, unknown>>>,
): BffTransport {
  return {
    call: <T>(op: string, input?: unknown) => {
      const fn = functions[op];
      if (!fn) throw new Error(`IdentityX server function "${op}" is not configured`);
      return fn({ data: input }) as Promise<T>;
    },
    has: (op) => typeof functions[op] === "function",
  };
}

export function createTanStackIdentityXAuthProviderAdapter(
  functions: IdentityXServerFunctions,
  options: TanStackIdentityXClientOptions = {},
): AuthProviderAdapter {
  return createBffAuthAdapter(
    serverFnTransport(
      functions as unknown as Partial<Record<string, ServerFunction<unknown, unknown>>>,
    ),
    options,
  );
}

export function createTanStackServerProxyFetchers(
  proxy: ServerFunction<ServerProxyRequest, ServerProxyResult>,
  options: TanStackIdentityXClientOptions = {},
  onSessionInvalid: () => void = () => undefined,
) {
  return createBffProxyFetchers(
    {
      call: <T>(_op: string, input?: unknown) => proxy({ data: input as ServerProxyRequest }) as Promise<T>,
    },
    options,
    onSessionInvalid,
  );
}

export function createTanStackIdentityXIntegration(
  functions: IdentityXServerFunctions & {
    request: ServerFunction<ServerProxyRequest, ServerProxyResult>;
  },
  options: TanStackIdentityXClientOptions = {},
) {
  let invalidateSession = () => undefined;
  const baseAuthAdapter = createTanStackIdentityXAuthProviderAdapter(functions, options);
  const authAdapter: AuthProviderAdapter = {
    createAuth(context) {
      invalidateSession = () => {
        context.setProfile(null);
        context.setAuthenticated(false);
      };
      return baseAuthAdapter.createAuth(context);
    },
    restoreSession: baseAuthAdapter.restoreSession,
  };
  const fetchers = createTanStackServerProxyFetchers(
    functions.request,
    options,
    () => invalidateSession(),
  );
  return {
    authAdapter,
    ...fetchers,
  };
}

export type { ServerAuthResult, ServerProxyRequest, ServerProxyResult, ServerSessionSnapshot } from "./types";
