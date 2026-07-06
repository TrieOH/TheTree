import { PostHogProvider as CorePostHogProvider } from "@trieoh/front-core"
import { env } from "@/env"
import type { ReactNode } from "react"

export function PHProvider({ children }: { children: ReactNode }) {
  return (
    <CorePostHogProvider
      config={{
        key: env.VITE_PUBLIC_POSTHOG_KEY,
        host: env.VITE_PUBLIC_POSTHOG_HOST,
      }}
    >
      {children}
    </CorePostHogProvider>
  )
}