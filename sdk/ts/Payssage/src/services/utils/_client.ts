import { createDefaultFetchClient, type DefaultFetchResult } from "@soramux/node-fetch-sdk";

export type ApiResponse<T> = DefaultFetchResult<T>;

const BASE_URL = "https://api.triepayments.trieoh.com";

export function createClient(baseURL = BASE_URL) {
  return createDefaultFetchClient({
    baseURL,
    headers: {
      "X-API-Key": process.env.API_KEY ?? "",
    },
  });
}

export const client = createClient();