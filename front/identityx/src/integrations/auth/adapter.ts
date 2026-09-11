import { createTanStackIdentityXIntegration } from "@trieoh/front-core/auth/tanstack/client";
import {
  authenticatedProxyServerFn,
  completeProviderLoginServerFn,
  introspectServerFn,
  isSetupDoneServerFn,
  loginServerFn,
  loginWithProviderServerFn,
  logoutServerFn,
  refreshServerFn,
  registerServerFn,
  restoreSessionServerFn,
  setupServerFn,
} from "./server-functions";

export const identityXIntegration = createTanStackIdentityXIntegration(
  {
    isSetupDone: () => isSetupDoneServerFn(),
    setup: setupServerFn,
    login: loginServerFn,
    register: registerServerFn,
    loginWithProvider: loginWithProviderServerFn,
    completeProviderLogin: completeProviderLoginServerFn,
    logout: () => logoutServerFn(),
    refresh: () => refreshServerFn(),
    restore: () => restoreSessionServerFn(),
    introspect: introspectServerFn,
    request: authenticatedProxyServerFn,
  },
  {},
);

export const identityXAuthAdapter = identityXIntegration.authAdapter;
