import { createDefaultFetchClient, type DefaultFetchResult } from "@soramux/node-fetch-sdk";

export type ApiResponse<T> = DefaultFetchResult<T>;

const DEFAULT_BASE_URL = "https://api.triepayments.trieoh.com";

export function createClient(
  baseURL = process.env.TRIEOH_PAY_BASE_URL ?? DEFAULT_BASE_URL,
  apiKey = process.env.TRIEOH_PAY_SECRET_KEY ?? ""
) {
  return createDefaultFetchClient({
    baseURL,
    headers: {
      "X-API-Key": apiKey,
    },
  });
}

export const client = createClient();