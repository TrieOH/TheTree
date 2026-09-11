import { useSession, getRequest } from "@tanstack/react-start/server";
import type { AuthTokens } from "@trieoh/identityx-sdk-ts";
import { createIdentityXBff, type BffSession } from "../bff/core";
import type { IdentityXTransportLogEvent } from "../bff/types";

interface IdentityXSessionData {
  tokens: AuthTokens;
}

interface TanStackIdentityXBffRuntime {
  getRequest?: () => Request;
  useSession?: () => Promise<any>;
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
  runtime?: TanStackIdentityXBffRuntime;
  apiBaseURL: string;
  observability?: {
    log?: (event: IdentityXTransportLogEvent) => void;
    logSuccesses?: boolean;
  };
}

/**
 * TanStack Start binding for `createIdentityXBff`: it only supplies the runtime
 * pieces (encrypted session cookie, incoming request). The IdentityX contract —
 * endpoints, envelopes, refresh, proxy guard — lives in the shared core.
 */
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

  return createIdentityXBff({
    identityX: config.identityX,
    apiBaseURL: config.apiBaseURL,
    observability: config.observability,
    request: () => config.runtime?.getRequest?.() ?? getRequest(),
    session: async (): Promise<BffSession> => {
      const current = await (config.runtime?.useSession?.() ??
        useSession<IdentityXSessionData>(sessionConfig));

      return {
        read: () => current.data ?? null,
        write: (tokens) => current.update({ tokens }),
        clear: () => current.clear(),
      };
    },
  });
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
} from "../bff/types";
