import { createFetcher, createQueryFetcher } from "@soramux/identityx-sdk-ts";
import { env } from "#/env";

export const authFetcher = createFetcher(
  {
    baseURL: env.VITE_API_URL,
    authBaseURL: env.VITE_AUTH_API_URL,
  }
);

export const tanstackQueryFetcher = createQueryFetcher(
  {
    baseURL: env.VITE_API_URL,
    authBaseURL: env.VITE_AUTH_API_URL,
  }
);
