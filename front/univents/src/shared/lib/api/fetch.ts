import { createAppFetchers } from "@trieoh/api-client";
import { createTanStackServerProxyFetchers } from "@trieoh/front-core/auth/tanstack/client";
import { env } from "@/env";
import { authenticatedProxyServerFn } from "@/integrations/auth/server-functions";

const { publicFetcher, publicQueryFetcher } = createAppFetchers({
  apiURL: env.VITE_API_URL,
  authAPIURL: env.VITE_AUTH_API_URL,
  timeout: 10_000,
});

const { authFetcher, authQueryFetcher } = createTanStackServerProxyFetchers(
  authenticatedProxyServerFn,
);

export {
  authFetcher,
  authQueryFetcher,
  publicFetcher,
  publicQueryFetcher,
  authQueryFetcher as tanstackQueryFetcher,
};
