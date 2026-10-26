import { QueryClient } from "@trieoh/front-core-solid";
import { CIRCUIT_OPEN_ERROR_ID, classifyError } from "./service-error";
import { serviceHealth } from "./service-health";

/**
 * Reads retry once, and only when a retry can plausibly help: the server never
 * answered or answered 5xx. Timeouts are not retried, since the next attempt
 * would just wait the same deadline again; a 4xx is a definitive answer; and an
 * open circuit already knows the service is down.
 */
export function shouldRetryRead(failureCount: number, error: unknown): boolean {
  const failure = classifyError(error);
  if (failure.errorId === CIRCUIT_OPEN_ERROR_ID) return false;
  return failure.transient && failure.kind !== "timeout" && failureCount < 1;
}

export function retryReadDelay(failureCount: number): number {
  const base = Math.min(600 * 2 ** failureCount, 4_000);
  return base + Math.random() * base * 0.25;
}

export function createAppQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 1000 * 60 * 5,
        // Keep what was loaded long enough to serve it while the next
        // navigation revalidates, instead of showing a skeleton again.
        gcTime: 1000 * 60 * 30,
        retry: shouldRetryRead,
        retryDelay: retryReadDelay,
        // Try once even when the browser reports offline (its flag is often
        // wrong), and fall back to the cache instead of pausing forever.
        networkMode: "offlineFirst",
        refetchOnReconnect: "always",
      },
    },
  });
}

export const appQueryClient = createAppQueryClient();

/**
 * When a service's circuit cools down, retry the queries on screen that failed
 * because of it, so a page that was showing "couldn't load" heals by itself.
 */
export function refetchFailedQueries(client: QueryClient = appQueryClient) {
  return client.refetchQueries({
    type: "active",
    predicate: (query) => query.state.status === "error",
  });
}

serviceHealth.onProbe(() => void refetchFailedQueries());
