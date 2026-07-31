import { createTanStackServerProxyFetchers } from "@trieoh/front-core/auth/tanstack/client";
import { authenticatedProxyServerFn } from "@/integrations/auth/server-functions";

const { authFetcher, authQueryFetcher } = createTanStackServerProxyFetchers(
  authenticatedProxyServerFn,
);

export { authFetcher, authQueryFetcher as tanstackQueryFetcher };
