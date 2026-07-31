import {
  createTanStackIdentityXAuthProviderAdapter,
} from "@trieoh/front-core/auth/tanstack/client";
import {
  loginServerFn,
  logoutServerFn,
  refreshServerFn,
  restoreSessionServerFn,
} from "./server-functions";

export const identityXAuthAdapter = createTanStackIdentityXAuthProviderAdapter({
  login: loginServerFn,
  logout: () => logoutServerFn(),
  refresh: () => refreshServerFn(),
  restore: () => restoreSessionServerFn(),
});
