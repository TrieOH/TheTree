import { useAuth } from "@trieoh/identityx-sdk-ts/react"
import posthog from "posthog-js"
import { PostHogProvider as BasePostHogProvider } from "posthog-js/react"
import { useEffect, type ReactNode } from "react"

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
  useEffect(() => {
    if (
      typeof window === "undefined" ||
      !config.key ||
      config.key === "phc_xxx" ||
      posthog.__loaded
    )
      return

    posthog.init(config.key, {
      api_host: config.host || "https://us.i.posthog.com",
      person_profiles: config.personProfiles ?? "identified_only",
      capture_pageview: config.capturePageview ?? false,
      defaults: "2026-05-30",
    })
  }, [config.capturePageview, config.host, config.key, config.personProfiles])

  return <BasePostHogProvider client={posthog}>{children}</BasePostHogProvider>
}

export function AuthenticatedPostHogProvider({
  config,
  children,
}: {
  config: PostHogConfig
  children: ReactNode
}) {
  const { auth, isAuthenticated } = useAuth()
  const userId = isAuthenticated ? auth.profile()?.id : undefined

  useEffect(() => {
    if (!posthog.__loaded) return

    if (userId) {
      posthog.identify(userId)
    } else if (posthog.get_property("$user_id")) {
      posthog.reset()
    }
  }, [userId])

  return <PostHogProvider config={config}>{children}</PostHogProvider>
}
