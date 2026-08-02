import {
  configureApiClient,
  createAppFetchers,
  createOrvalTransport,
} from "@trieoh/api-client";
import { createTanStackServerProxyFetchers } from "@trieoh/front-core/auth/tanstack/client";
import { env } from "@/env";
import { authenticatedProxyServerFn } from "@/integrations/auth/server-functions";

const { publicFetcher } = createAppFetchers({
  apiURL: env.VITE_API_URL,
  authAPIURL: env.VITE_AUTH_API_URL,
  timeout: 10_000,
});

const { authFetcher } = createTanStackServerProxyFetchers(
  authenticatedProxyServerFn,
);

configureApiClient({
  baseURL: "",
  transport: createOrvalTransport(authFetcher),
  publicTransport: createOrvalTransport(publicFetcher),
});

// export {
//   authFetcher,
//   authQueryFetcher,
//   publicFetcher,
//   publicQueryFetcher,
//   authQueryFetcher as tanstackQueryFetcher,
// };
