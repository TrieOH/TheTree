import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ApiError } from '@soramux/identityx-sdk-ts';
import type { ReactNode } from 'react';

/**
 * Interface matching the SDK's failure envelope structure.
 * Used for type-safe error handling without 'any'.
 */
interface ApiFailureEnvelope {
  code: number;
  message: string;
  error_id?: string;
}

let context: { queryClient: QueryClient } | undefined

export function getContext() {
  if (context) return context

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: (failureCount, error) => {
          if (error instanceof ApiError) {
            const envelope = error.envelope as ApiFailureEnvelope
            if (envelope.code >= 400 && envelope.code < 500) return false
          }

          if (error instanceof Error) {
            const err = error as unknown as ApiFailureEnvelope;
            const status = err.code;
            if (status >= 400 && status < 500) return false;
          }

          return failureCount < 3;
        },
        staleTime: 0, // 0 Seconds for test 1000 * 60 * 5, // 5 minutes
        refetchOnMount: true,
        refetchOnWindowFocus: true,
        refetchOnReconnect: true,
      },
    },
  })

  context = { queryClient }

  return context
}

export default function TanStackQueryProvider({
  children,
}: {
  children: ReactNode
}) {
  const { queryClient } = getContext()

  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}
