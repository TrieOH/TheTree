import { configureApiClient, createOrvalTransport } from "@trieoh/api-client";
import { createTanStackServerProxyFetchers } from "@trieoh/front-core/auth/tanstack/client";
import { authenticatedProxyServerFn } from "@/integrations/auth/server-functions";

const { authFetcher } = createTanStackServerProxyFetchers(
  authenticatedProxyServerFn,
);

configureApiClient({
  baseURL: "",
  transport: createOrvalTransport(authFetcher),
});
