import { createFetcher, createQueryFetcher } from "@trieoh/identityx-sdk-ts";
import { createDefaultFetchClient } from "@trieoh/envoy-fetch-ts";
import { env } from "#/env";

export const authFetcher = createFetcher(
  {
    baseURL: env.VITE_API_URL,
    authBaseURL: env.VITE_AUTH_API_URL,
  }
);

export const publicFetcher = createDefaultFetchClient({
  baseURL: env.VITE_API_URL,
  credentials: "omit",
  timeout: 10_000, // 10 seconds timeout
});

export const tanstackQueryFetcher = createQueryFetcher(
  {
    baseURL: env.VITE_API_URL,
    authBaseURL: env.VITE_AUTH_API_URL,
  }
);
