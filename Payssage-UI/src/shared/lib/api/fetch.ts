import { createFetcher, createQueryFetcher } from "@trieoh/node-auth-sdk";
import { env } from "#/env";

export const authFetcher = createFetcher(
  {
    baseURL: env.VITE_API_URL,
    authBaseURL: "http://localhost:8080",
  }
);

export const tanstackQueryFetcher = createQueryFetcher(
  {
    baseURL: env.VITE_API_URL,
    authBaseURL: "http://localhost:8080",
  }
);
