import {
  configureApiClient,
  createAppFetchers,
  createOrvalTransport,
} from "@trieoh/api-client";
import { env } from "@/env";
import { identityXIntegration } from "@/integrations/auth/adapter";

const { publicFetcher } = createAppFetchers({
  apiURL: env.VITE_API_URL,
  authAPIURL: env.VITE_AUTH_API_URL,
  timeout: 10_000,
});

configureApiClient({
  baseURL: "",
  transport: createOrvalTransport(identityXIntegration.authFetcher),
  publicTransport: createOrvalTransport(publicFetcher),
});

// export {
//   authFetcher,
//   authQueryFetcher,
//   publicFetcher,
//   publicQueryFetcher,
//   authQueryFetcher as tanstackQueryFetcher,
// };
