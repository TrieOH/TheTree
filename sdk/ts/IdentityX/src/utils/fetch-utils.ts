import { logger } from "./logger";

interface SimpleFetchOptions {
  method?: string;
  headers?: Record<string, string>;
  body?: string;
}

export async function simpleFetch<T>(url: string, options?: SimpleFetchOptions): Promise<T> {
  const response = await fetch(url, {
    ...options,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
  });

  const data = await response.json().catch(() => {
    logger.error("Failed to parse response as JSON from", url);
    throw new Error(`Failed to parse response from ${url}`);
  });

  return data as T;
}