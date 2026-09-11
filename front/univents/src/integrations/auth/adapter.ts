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
  { projectId: env.VITE_TRIEOH_AUTH_PROJECT_ID },
);

export const identityXAuthAdapter = identityXIntegration.authAdapter;
