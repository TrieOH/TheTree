import { createTanStackIdentityXAuthProviderAdapter } from "@trieoh/front-core-react";
import {
  loginServerFn,
  logoutServerFn,
  refreshServerFn,
  registerServerFn,
  restoreSessionServerFn,
} from "./server-functions";

export const identityXAuthAdapter = createTanStackIdentityXAuthProviderAdapter({
  login: loginServerFn,
  register: registerServerFn,
  logout: () => logoutServerFn(),
  refresh: () => refreshServerFn(),
  restore: () => restoreSessionServerFn(),
});
