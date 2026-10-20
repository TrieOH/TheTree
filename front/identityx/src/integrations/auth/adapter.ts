import { createTanStackIdentityXIntegration } from "@trieoh/front-core-react";
import {
  apiServerFn,
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
    request: ({ data }) =>
      apiServerFn({ data, fetch: (input, init) => fetch(labelApiCall(input, data), init) }),
  },
  {},
);

export const identityXAuthAdapter = identityXIntegration.authAdapter;

/**
 * Every API call goes through the one `api` server function, so in devtools
 * they would all be named "api". This appends the real call to the URL —
 * /fn/api/projects?GET shows up as "projects?GET". The server routes by the
 * first segment after /fn/ and reads the call from the body, so the label
 * never changes what is requested.
 */
function labelApiCall(
  input: RequestInfo | URL,
  data: { path: string; method?: string },
): RequestInfo | URL {
  if (typeof input !== "string") return input;
  const path = data.path.split("?")[0].replace(/^\/+/, "");
  return `${input}/${path}?${data.method ?? "GET"}`;
}
