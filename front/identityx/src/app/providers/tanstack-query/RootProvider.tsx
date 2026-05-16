import { ApiError } from '@trieoh/identityx-sdk-ts';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

export function getContext() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: (failureCount, error) => {
          if (error instanceof ApiError) {
            const envelope = error.envelope
            if (envelope.code >= 400 && envelope.code < 500) return false
          }
          if (error instanceof Error) {
            const err = error as unknown as {code: number};
            const status = err.code;
            if (status >= 400 && status < 500) return false;
          }
          return failureCount < 3;
        },
        staleTime: 1000 * 60 * 5, // 5 minutes
        refetchOnMount: true,
        refetchOnWindowFocus: true,
        refetchOnReconnect: true,
      },
    },
  });
  return {
    queryClient,
  }
}

export function Provider({
  children,
  queryClient,
}: {
  children: React.ReactNode
  queryClient: QueryClient
}) {
  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  )
}
