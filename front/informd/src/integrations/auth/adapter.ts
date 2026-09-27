import { createTanStackIdentityXAuthProviderAdapter } from "@trieoh/front-core-react";
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
