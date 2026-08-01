import { createAppFetchers } from "@trieoh/api-client";
import { createTanStackServerProxyFetchers } from "@trieoh/front-core/auth/tanstack/client";
import { env } from "#/env";
import { authenticatedProxyServerFn } from "#/integrations/auth/server-functions";

const { publicFetcher } = createAppFetchers({
  apiURL: env.VITE_API_URL,
  authAPIURL: env.VITE_AUTH_API_URL,
});

const { authFetcher, authQueryFetcher } = createTanStackServerProxyFetchers(
  authenticatedProxyServerFn,
);

export { authFetcher, authQueryFetcher as tanstackQueryFetcher, publicFetcher };
