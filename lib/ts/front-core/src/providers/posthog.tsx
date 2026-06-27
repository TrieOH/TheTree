import posthog from "posthog-js"
import { PostHogProvider as BasePostHogProvider } from "posthog-js/react"
import type { ReactNode } from "react"

export interface PostHogConfig {
  key: string
  host?: string
  capturePageview?: boolean
  personProfiles?: "identified_only" | "always" | "never"
}

/**
 * Initialize and provide PostHog analytics.
 *
 * Usage:
 *   <PostHogProvider config={{ key: env.VITE_POSTHOG_KEY, host: env.VITE_POSTHOG_HOST }}>
 *     {children}
 *   </PostHogProvider>
 */
export function PostHogProvider({
  config,
  children,
}: {
  config: PostHogConfig
  children: ReactNode
}) {
  if (typeof window !== "undefined" && config.key && config.key !== "phc_xxx") {
    posthog.init(config.key, {
      api_host: config.host || "https://us.i.posthog.com",
      person_profiles: config.personProfiles ?? "identified_only",
      capture_pageview: config.capturePageview ?? false,
      defaults: "2025-11-30",
    })
  }

  return <BasePostHogProvider client={posthog}>{children}</BasePostHogProvider>
}
