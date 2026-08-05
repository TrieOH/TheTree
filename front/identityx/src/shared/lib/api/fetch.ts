import { configureApiClient, createOrvalTransport } from "@trieoh/api-client";
import { identityXIntegration } from "@/integrations/auth/adapter";

configureApiClient({
  baseURL: "",
  transport: createOrvalTransport(identityXIntegration.authFetcher),
});
