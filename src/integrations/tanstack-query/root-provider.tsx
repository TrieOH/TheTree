import { env } from '@/env';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { posthog } from 'posthog-js';
import { PostHogProvider } from 'posthog-js/react';

export function getContext() {
  const queryClient = new QueryClient()
  return {
    queryClient,
  }
}

posthog.init(env.VITE_PUBLIC_POSTHOG_KEY, { 
  api_host: env.VITE_PUBLIC_POSTHOG_HOST, 
  defaults: '2025-11-30', 
}); 

export function Provider({
  children,
  queryClient,
}: {
  children: React.ReactNode
  queryClient: QueryClient
}) {
  return (
    <PostHogProvider client={posthog}>
      <QueryClientProvider client={queryClient}>
        {children}
      </QueryClientProvider>
    </PostHogProvider>
  )
}
