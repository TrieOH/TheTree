import { TanStackDevtools } from "@tanstack/react-devtools";
import type { QueryClient } from "@tanstack/react-query";
import {
  createRootRouteWithContext,
  HeadContent,
  Scripts,
} from "@tanstack/react-router";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import { AuthContextUpdater } from "@trieoh/front-core";
import type { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { AuthProvider } from "@trieoh/identityx-sdk-ts/react";
import { env } from "#/env";
import { Toaster } from "#/shared/ui/shadcn/sonner";
import PostHogProvider from "../app/integrations/posthog/provider";
import TanStackQueryDevtools from "../app/integrations/tanstack-query/devtools";
import { Provider as TanStackQueryProvider } from "../app/integrations/tanstack-query/root-provider";
import appCss from "../styles.css?url";

interface MyRouterContext {
  queryClient: QueryClient;
  auth?: ReturnType<typeof useAuth>;
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      { title: env.VITE_APP_TITLE ?? "Payssage" },
      {
        name: "apple-mobile-web-app-title",
        content: env.VITE_APP_TITLE ?? "Payssage",
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
      <body className="min-w-[320px] font-sans antialiased wrap:anywhere selection:bg-primary/10">
        <PostHogProvider>
          <TanStackQueryProvider>
            <AuthProvider baseURL={env.VITE_AUTH_API_URL}>
              <AuthContextUpdater>
                {children}
                <TanStackDevtools
                  config={{
                    position: "bottom-right",
                  }}
                  plugins={[
                    {
                      name: "Tanstack Router",
                      render: <TanStackRouterDevtoolsPanel />,
                    },
                    TanStackQueryDevtools,
                  ]}
                />
              </AuthContextUpdater>
            </AuthProvider>
          </TanStackQueryProvider>
        </PostHogProvider>
        <Toaster />
        <Scripts />
      </body>
    </html>
  );
}
