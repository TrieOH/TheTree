import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import type { ReactNode } from 'react'

export function getContext() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 0, // 0 Seconds for test //1000 * 60 * 5, // 5 minutes
        refetchOnMount: true,
        refetchOnWindowFocus: true,
        refetchOnReconnect: true,
      }
    }
  })

  return {
    queryClient,
  }
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

