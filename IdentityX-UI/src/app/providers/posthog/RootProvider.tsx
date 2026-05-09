import posthog from 'posthog-js';
import { PostHogProvider } from 'posthog-js/react';
import { env } from '@/env';

export function PHProvider({ children }: { children: React.ReactNode }) {
  posthog.init(env.VITE_PUBLIC_POSTHOG_KEY, {
    api_host: env.VITE_PUBLIC_POSTHOG_HOST,
    defaults: '2025-11-30',
  })

  return <PostHogProvider client={posthog}>{children}</PostHogProvider>
}