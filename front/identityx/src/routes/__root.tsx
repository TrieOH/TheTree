import type { QueryClient } from "@tanstack/react-query";
import {
  createRootRouteWithContext,
  HeadContent,
  Scripts,
} from "@tanstack/react-router";
import { AuthContextUpdater } from "@trieoh/front-core";
import { AuthProvider, type useAuth } from "@trieoh/identityx-sdk-ts/react";
import { Toaster } from "@trieoh/ui-base/shadcn/sonner";
import { PHProvider } from "@/app/providers/posthog/RootProvider";
import { Provider as TanStackQueryProvider } from "@/app/providers/tanstack-query/RootProvider";
import { env } from "@/env";
import { requireSetup } from "@/features/auth/lib/route-guard";
import { identityXAuthAdapter } from "@/integrations/auth/adapter";
import appCss from "../styles.css?url";

interface MyRouterContext {
  queryClient: QueryClient;
  auth?: ReturnType<typeof useAuth>;
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  beforeLoad: requireSetup,
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      { title: env.VITE_APP_TITLE ?? "IdentityX" },
      {
        name: "apple-mobile-web-app-title",
        content: env.VITE_APP_TITLE ?? "IdentityX",
      },
      { name: "mobile-web-app-capable", content: "yes" },
    ],
    links: [
      { rel: "stylesheet", href: appCss },
      { rel: "manifest", href: "/site.webmanifest" },
      {
        rel: "icon",
        type: "image/png",
        href: "/favicon-96x96.png",
        sizes: "96x96",
      },
      { rel: "icon", href: "/favicon.svg", type: "image/svg+xml" },
      { rel: "shortcut icon", href: "/favicon.ico" },
      {
        rel: "apple-touch-icon",
        href: "/apple-touch-icon.png",
        sizes: "180x180",
      },
    ],
  }),

  shellComponent: RootDocument,
  notFoundComponent: () => {
    return <p>This page doesn't exist!</p>;
  },
});

function RootDocument({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <HeadContent />
      </head>
      <body className="font-body antialiased wrap-anywhere">
        <TanStackQueryProvider>
          <AuthProvider
            baseURL={env.VITE_API_URL}
            isProjectMode={false}
            adapter={identityXAuthAdapter}
          >
            <PHProvider>
              <AuthContextUpdater>{children}</AuthContextUpdater>
            </PHProvider>
          </AuthProvider>
        </TanStackQueryProvider>
        <Toaster />
        <Scripts />
      </body>
    </html>
  );
}
