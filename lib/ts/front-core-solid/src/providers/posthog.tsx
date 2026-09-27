import { useAuth } from "@trieoh/identityx-sdk-ts-solid"
import posthog from "posthog-js"
import { createEffect, onSettled, type ParentProps } from "solid-js"

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
 * Boot PostHog analytics.
 * Usage:
 *   <PostHogProvider config={{ key: import.meta.env.VITE_POSTHOG_KEY, host: import.meta.env.VITE_POSTHOG_HOST }}>
 *     {children}
 *   </PostHogProvider>
 */
export function PostHogProvider(props: ParentProps<{ config: PostHogConfig }>) {
  onSettled(() => initPostHog(props.config))

  return props.children
}

/**
 * {@link PostHogProvider} plus person identification: the IdentityX subject is
 * the distinct id, and its email travels with it as a person property.
 *
 * The subject carries no display name (id, email, project, capabilities and
 * free-form metadata only) — the name lives in the actor's profile document, so
 * an app that loads it can enrich the person afterwards with
 * `posthog.people.set({ name })`.
 */
export function AuthenticatedPostHogProvider(
  props: ParentProps<{ config: PostHogConfig }>,
) {
  onSettled(() => initPostHog(props.config))

  const { auth, isAuthenticated } = useAuth()

  createEffect(
    () => {
      const subject = isAuthenticated() ? auth.profile() : null
      return [subject?.id, subject?.email] as const
    },
    ([id, email]) => {
      if (!posthog.__loaded) return

      if (id) posthog.identify(id, email ? { email } : undefined)
      else if (posthog.get_property("$user_id")) posthog.reset()
    },
  )

  return props.children
}

function initPostHog(config: PostHogConfig) {
  // SSR renders the provider too, and a placeholder key means "not configured"
  if (
    typeof window === "undefined" ||
    !config.key ||
    config.key === "phc_xxx" ||
    posthog.__loaded
  ) {
    return
  }

  posthog.init(config.key, {
    api_host: config.host || "https://us.i.posthog.com",
    person_profiles: config.personProfiles ?? "identified_only",
    capture_pageview: config.capturePageview ?? "history_change",
    defaults: "2026-05-30",
  })
}
