import { PostHogProvider as CorePostHogProvider } from "@trieoh/front-core";
import type { ReactNode } from "react";
import { env } from "@/env";

export default function PHProvider({ children }: { children: ReactNode }) {
  return (
    <CorePostHogProvider
      config={{
        key: env.VITE_POSTHOG_KEY,
        host: env.VITE_POSTHOG_HOST,
      }}
    >
      {children}
    </CorePostHogProvider>
  );
}
