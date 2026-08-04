import { createTanStackIdentityXIntegration } from "@trieoh/front-core/auth/tanstack/client";
import { env } from "@/env";
import {
  authenticatedProxyServerFn,
  completeProviderLoginServerFn,
  loginServerFn,
  loginWithProviderServerFn,
  logoutServerFn,
  refreshServerFn,
  restoreSessionServerFn,
} from "./server-functions";

export const identityXIntegration = createTanStackIdentityXIntegration(
  {
    login: loginServerFn,
    loginWithProvider: loginWithProviderServerFn,
    completeProviderLogin: completeProviderLoginServerFn,
    logout: () => logoutServerFn(),
    refresh: () => refreshServerFn(),
    restore: () => restoreSessionServerFn(),
    request: authenticatedProxyServerFn,
  },
  {
    mode: env.VITE_AUTH_TRANSPORT,
    apiBaseURL: env.VITE_API_URL,
    authBaseURL: env.VITE_AUTH_API_URL,
    projectId: env.VITE_TRIEOH_AUTH_PROJECT_ID,
  },
);

export const identityXAuthAdapter = identityXIntegration.authAdapter;
