import { AuthenticatedPostHogProvider } from "@trieoh/front-core";
import type { ReactNode } from "react";
import { env } from "@/env";

export function PHProvider({ children }: { children: ReactNode }) {
  return (
    <AuthenticatedPostHogProvider
      config={{
        key: env.VITE_PUBLIC_POSTHOG_KEY,
        host: env.VITE_PUBLIC_POSTHOG_HOST,
      }}
    >
      {children}
    </AuthenticatedPostHogProvider>
  );
}
