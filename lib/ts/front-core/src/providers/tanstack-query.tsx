import {
  MutationCache,
  QueryCache,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
import { ApiError } from "@trieoh/identityx-sdk-ts";
import type { ReactNode } from "react";

export interface QueryClientConfig {
  /** Stale time in milliseconds (default: 5 minutes). */
  staleTime?: number;
  /** Maximum retry count (default: 3). */
  maxRetries?: number;
  /** Called when a query or mutation finishes with an error. */
  onError?: (error: unknown) => void;
}

export class QueryError extends Error {
  readonly envelope: { code?: number; message: string };

  constructor(message: string, code?: number) {
    super(message);
    this.name = "QueryError";
    this.envelope = { message, code };
  }
}

function getErrorStatus(error: unknown): number | undefined {
  if (typeof error !== "object" || error === null) return undefined;

  if (error instanceof ApiError || error instanceof QueryError) {
    return error.envelope.code;
  }

  const candidate = error as {
    status?: unknown;
    response?: { status?: unknown };
  };
  const status = [candidate.status, candidate.response?.status].find(
    (value): value is number => typeof value === "number",
  );

  return status;
}

export function queryError(message: string, code?: number): QueryError {
  return new QueryError(message, code);
}

/**
 * Create a new QueryClient with standard TrieOH defaults.
 * The consumer should create this and pass to TanStackQueryProvider
 * to avoid version-mismatch issues with private fields.
 */
export function createQueryClient(config?: QueryClientConfig) {
  const {
    staleTime = 1000 * 60 * 5, // 5 minutes
    maxRetries = 3,
  } = config ?? {};

  return new QueryClient({
    queryCache: new QueryCache({ onError: config?.onError }),
    mutationCache: new MutationCache({ onError: config?.onError }),
    defaultOptions: {
      queries: {
        retry: (failureCount, error) => {
          const status = getErrorStatus(error);
          if (status && status >= 400 && status < 500) return false;
          return failureCount < maxRetries;
        },
        staleTime,
        refetchOnMount: true,
        refetchOnWindowFocus: true,
        refetchOnReconnect: true,
      },
    },
  });
}

/**
 * TanStack Query provider.
 * Accepts a `queryClient` created by the consumer (avoids version-mismatch).
 *
 * Usage:
 *   const queryClient = createQueryClient()
 *   <TanStackQueryProvider queryClient={queryClient}>
 *     {children}
 *   </TanStackQueryProvider>
 */
export function TanStackQueryProvider({
  children,
  queryClient,
}: {
  children: ReactNode;
  queryClient: QueryClient;
}) {
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}
