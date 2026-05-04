import { createFetcher, createQueryFetcher } from "trieoh/identityx-sdk-ts";
import { env } from "@/env";

export const authFetcher = createFetcher(
  {
    baseURL: env.VITE_API_URL,
    authBaseURL: env.VITE_API_URL,
    clientConfig: {
      timeout: 10_000, // 10 seconds timeout
    }
  }
);

export const tanstackQueryFetcher = createQueryFetcher(
  {
    baseURL: env.VITE_API_URL,
    authBaseURL: env.VITE_API_URL,
    clientConfig: {
      timeout: 10_000, // 10 seconds timeout
    }
  }
);
