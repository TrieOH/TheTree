import {
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query"
import { ApiError } from "@trieoh/identityx-sdk-ts"
import type { ReactNode } from "react"

export interface QueryClientConfig {
  /** Stale time in milliseconds (default: 5 minutes). */
  staleTime?: number
  /** Maximum retry count (default: 3). */
  maxRetries?: number
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
  } = config ?? {}

  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: (failureCount, error) => {
          if (error instanceof ApiError) {
            const envelope = error.envelope as { code: number }
            if (envelope.code >= 400 && envelope.code < 500) return false
          }
          if (error instanceof Error) {
            const err = error as unknown as { code: number }
            const status = err.code
            if (status >= 400 && status < 500) return false
          }
          return failureCount < maxRetries
        },
        staleTime,
        refetchOnMount: true,
        refetchOnWindowFocus: true,
        refetchOnReconnect: true,
      },
    },
  })
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
  children: ReactNode
  queryClient: QueryClient
}) {
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}
