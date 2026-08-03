import { createTanStackIdentityXAuthProviderAdapter } from "@trieoh/front-core/auth/tanstack/client";
import {
  completeProviderLoginServerFn,
  loginServerFn,
  loginWithProviderServerFn,
  logoutServerFn,
  refreshServerFn,
  restoreSessionServerFn,
} from "./server-functions";

export const identityXAuthAdapter = createTanStackIdentityXAuthProviderAdapter({
  login: loginServerFn,
  loginWithProvider: loginWithProviderServerFn,
  completeProviderLogin: completeProviderLoginServerFn,
  logout: () => logoutServerFn(),
  refresh: () => refreshServerFn(),
  restore: () => restoreSessionServerFn(),
});
