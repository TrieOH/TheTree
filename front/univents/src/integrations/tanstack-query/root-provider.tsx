import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ApiError } from '@trieoh/identityx-sdk-ts'
import type { DefaultFailureEnvelope } from "@soramux/node-fetch-sdk";
import type { ReactNode } from 'react'

let context: { queryClient: QueryClient } | undefined

export function getContext() {
  if (context) return context

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: (failureCount, error) => {
          if (error instanceof ApiError) {
            const envelope = error.envelope as DefaultFailureEnvelope
            if (envelope.code >= 400 && envelope.code < 500) return false
          }
          if (error instanceof Error) {
            const err = error as unknown as DefaultFailureEnvelope;
            const status = err.code;
            if (status >= 400 && status < 500) return false;
          }
          return failureCount < 3;
        },
        staleTime: 0, // 0 Seconds for test //1000 * 60 * 5, // 5 minutes
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
