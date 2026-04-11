import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import type { ReactNode } from 'react'

export function getContext() {
  const queryClient = new QueryClient()

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

