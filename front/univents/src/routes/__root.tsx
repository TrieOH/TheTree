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
import { Toaster } from "@trieoh/ui-base/shadcn/sonner";
import { ThemeProvider } from "next-themes";
import { env } from "@/env";
import { requireConfiguredProfile } from "@/features/auths/lib/route-guard";
import { UploadQueueProvider } from "@/features/upload-queue/ui/upload-queue-provider";
import "@/features/upload-queue/associations";
import WaveSpinnerLoading from "@/shared/ui/loader/WaveSpinnerLoading";
import NotFound from "@/widgets/feedback/ui/NotFound";
import { NavigationDock } from "@/widgets/ui/navigation-dock";
import { identityXAuthAdapter } from "../integrations/auth/adapter";
import PostHogProvider from "../integrations/posthog/provider";
import TanStackQueryDevtools from "../integrations/tanstack-query/devtools";
import { Provider as TanStackQueryProvider } from "../integrations/tanstack-query/root-provider";
import appCss from "../styles.css?url";

interface MyRouterContext {
  queryClient: QueryClient;
  auth?: ReturnType<typeof useAuth>;
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  beforeLoad: requireConfiguredProfile,
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      { title: env.VITE_APP_TITLE ?? "Univents" },
      {
        name: "apple-mobile-web-app-title",
        content: env.VITE_APP_TITLE ?? "Univents",
      },
      { name: "mobile-web-app-capable", content: "yes" },
    ],
    links: [
      { rel: "stylesheet", href: appCss },
      { rel: "manifest", href: "/site.webmanifest" },
      { rel: "icon", href: "/favicon.svg?v=2", type: "image/svg+xml" },
      {
        rel: "apple-touch-icon",
        sizes: "180x180",
        href: "/apple-touch-icon.png?v=2",
      },
    ],
  }),
  shellComponent: RootDocument,
  notFoundComponent: NotFound,
});

function RootDocument({ children }: { children: React.ReactNode }) {
  return (
    <html lang="pt-BR" suppressHydrationWarning>
      <head>
        <script
          crossOrigin="anonymous"
          src="//unpkg.com/react-scan/dist/auto.global.js"
        />
        <HeadContent />
      </head>
      <body className="min-w-[320px] font-sans antialiased wrap:anywhere selection:bg-primary/10">
        <PostHogProvider>
          <TanStackQueryProvider>
            <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
              <AuthProvider
                adapter={identityXAuthAdapter}
                baseURL={env.VITE_AUTH_API_URL}
                fallback={
                  <div className="h-screen w-screen flex items-center justify-center">
                    <WaveSpinnerLoading text="Carregando..." />
                  </div>
                }
              >
                <AuthContextUpdater>
                  <UploadQueueProvider>
                    {children}
                    <NavigationDock className="print:hidden" />
                  </UploadQueueProvider>
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
            </ThemeProvider>
          </TanStackQueryProvider>
        </PostHogProvider>
        <Toaster />
        <Scripts />
      </body>
    </html>
  );
}
