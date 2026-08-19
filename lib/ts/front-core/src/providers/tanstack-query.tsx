import {
  MutationCache,
  QueryCache,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
import { ApiError } from "@trieoh/identityx-sdk-ts";
import type { ReactNode } from "react";
import { recordCompletedSpan } from "../tracing/browser";

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

  const queryStartedAt = new WeakMap<object, number>();
  const queryCache = new QueryCache({
    onSuccess: (_data, query) => {
      const startedAt = queryStartedAt.get(query);
      if (startedAt === undefined) return;
      queryStartedAt.delete(query);
      recordCompletedSpan(queryOperationName(query.queryKey), startedAt, {
        "operation.type": "query",
        "query.outcome": "success",
        "query.failure_count": query.state.fetchFailureCount,
      });
    },
    onError: (error, query) => {
      config?.onError?.(error);
      const startedAt = queryStartedAt.get(query);
      if (startedAt === undefined) return;
      queryStartedAt.delete(query);
      recordCompletedSpan(
        queryOperationName(query.queryKey),
        startedAt,
        {
          "operation.type": "query",
          "query.outcome": "failure",
          "query.failure_count": query.state.fetchFailureCount,
          ...(getErrorStatus(error) === undefined
            ? {}
            : { "query.status_code": getErrorStatus(error) as number }),
        },
        error,
      );
    },
  });
  queryCache.subscribe((event) => {
    if (event.type === "updated" && event.action.type === "fetch") {
      queryStartedAt.set(event.query, Date.now());
    }
  });

  return new QueryClient({
    queryCache,
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

function queryOperationName(queryKey: readonly unknown[]): string {
  const parts = queryKey
    .filter((part): part is string => typeof part === "string")
    .filter(
      (part) =>
        !part.includes("/") &&
        !/^[0-9a-f]{8}(?:-[0-9a-f]{4}){3}-[0-9a-f]{12}$/i.test(part),
    )
    .slice(0, 2);
  return `query:${parts.join(".") || "anonymous"}`;
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
