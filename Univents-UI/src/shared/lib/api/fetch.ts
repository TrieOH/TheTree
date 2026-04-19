import { createFetcher, createQueryFetcher } from "@soramux/node-auth-sdk";
import { createDefaultFetchClient } from "@soramux/node-fetch-sdk";
import { env } from "@/env";

/**
 * Used to handle non authenticated request
 */
export const simpleFetcher = createDefaultFetchClient({
  baseURL: env.VITE_API_URL,
  credentials: "omit",
  timeout: 10_000, // 10 seconds timeout
})

export const authFetcher = createFetcher(
  {
    baseURL: env.VITE_API_URL,
    authBaseURL: env.VITE_AUTH_API_URL,
    clientConfig: {
      timeout: 10_000, // 10 seconds timeout
    }
  }
);

export const tanstackQueryFetcher = createQueryFetcher(
  {
    baseURL: env.VITE_API_URL,
    authBaseURL: env.VITE_AUTH_API_URL,
    clientConfig: {
      timeout: 10_000, // 10 seconds timeout
    }
  }
);
