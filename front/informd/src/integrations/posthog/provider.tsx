import { AuthenticatedPostHogProvider } from "@trieoh/front-core";
import type { ReactNode } from "react";
import { env } from "#/env";

export default function PostHogProvider({ children }: { children: ReactNode }) {
  return (
    <AuthenticatedPostHogProvider
      config={{
        key: env.VITE_POSTHOG_KEY,
        host: env.VITE_POSTHOG_HOST,
      }}
    >
      {children}
    </AuthenticatedPostHogProvider>
  );
}
