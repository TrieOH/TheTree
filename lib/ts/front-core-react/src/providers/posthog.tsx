import { useAuth } from "@trieoh/identityx-sdk-ts-react"
import posthog from "posthog-js"
import { PostHogProvider as BasePostHogProvider } from "posthog-js/react"
import { useEffect, type ReactNode } from "react"

export interface PostHogConfig {
  key: string
  host?: string
  /**
   * Pageview capture. Defaults to `"history_change"` - the SDK's own default
   */
  capturePageview?: boolean | "history_change"
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
      capture_pageview: config.capturePageview ?? "history_change",
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
  const subject = isAuthenticated ? auth.profile() : null

  useEffect(() => {
    if (!posthog.__loaded) return

    if (subject?.id) {
      // The IdentityX subject has no display name - that lives in the actor's
      // profile document, so an app that loads it can follow up with
      // `posthog.people.set({ name })`.
      posthog.identify(subject.id, subject.email ? { email: subject.email } : undefined)
    } else if (posthog.get_property("$user_id")) {
      posthog.reset()
    }
  }, [subject?.id, subject?.email])

  return <PostHogProvider config={config}>{children}</PostHogProvider>
}
