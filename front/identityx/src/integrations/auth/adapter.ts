import { createTanStackIdentityXAuthProviderAdapter } from "@trieoh/front-core/auth/tanstack/client";
import {
  completeProviderLoginServerFn,
  isSetupDoneServerFn,
  loginServerFn,
  loginWithProviderServerFn,
  logoutServerFn,
  refreshServerFn,
  registerServerFn,
  restoreSessionServerFn,
  setupServerFn,
} from "./server-functions";

export const identityXAuthAdapter = createTanStackIdentityXAuthProviderAdapter({
  isSetupDone: () => isSetupDoneServerFn(),
  setup: setupServerFn,
  login: loginServerFn,
  register: registerServerFn,
  loginWithProvider: loginWithProviderServerFn,
  completeProviderLogin: completeProviderLoginServerFn,
  logout: () => logoutServerFn(),
  refresh: () => refreshServerFn(),
  restore: () => restoreSessionServerFn(),
});
