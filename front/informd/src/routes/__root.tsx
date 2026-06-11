import {
  HeadContent,
  Scripts,
  createRootRouteWithContext,
} from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import { TanStackDevtools } from '@tanstack/react-devtools'

import PostHogProvider from '../integrations/posthog/provider'

import TanStackQueryProvider from '../integrations/tanstack-query/root-provider'

import TanStackQueryDevtools from '../integrations/tanstack-query/devtools'

import appCss from '../styles.css?url'

import type { QueryClient } from '@tanstack/react-query'
import { AuthContextUpdater } from '#/integrations/auth/auth-context-updater'
import { AuthProvider } from '@trieoh/identityx-sdk-ts/react'
import type { useAuth } from '@trieoh/identityx-sdk-ts/react';
import { env } from '#/env'
import { Toaster } from '#/shared/ui/shadcn/sonner'

interface MyRouterContext {
  queryClient: QueryClient
  auth?: ReturnType<typeof useAuth>
}

const THEME_INIT_SCRIPT = `(function(){try{var stored=window.localStorage.getItem('theme');var mode=(stored==='light'||stored==='dark'||stored==='auto')?stored:'auto';var prefersDark=window.matchMedia('(prefers-color-scheme: dark)').matches;var resolved=mode==='auto'?(prefersDark?'dark':'light'):mode;var root=document.documentElement;root.classList.remove('light','dark');root.classList.add(resolved);if(mode==='auto'){root.removeAttribute('data-theme')}else{root.setAttribute('data-theme',mode)}root.style.colorScheme=resolved;}catch(e){}})();`

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      { charSet: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { title: env.VITE_APP_TITLE ?? "Informd" },
      { name: 'apple-mobile-web-app-title', content: env.VITE_APP_TITLE ?? 'Informd' },
      { name: 'mobile-web-app-capable', content: 'yes' },
    ],
    links: [
      { rel: 'stylesheet', href: appCss },
      { rel: 'manifest', href: '/site.webmanifest' },
      { rel: 'icon', type: 'image/png', href: '/favicon-96x96.png', sizes: '96x96' },
      { rel: 'icon', href: '/favicon.svg', type: 'image/svg+xml' },
      { rel: 'shortcut icon', href: '/favicon.ico' },
      { rel: 'apple-touch-icon', href: '/apple-touch-icon.png', sizes: '180x180' },
    ],
  }),
  shellComponent: RootDocument,
})

function RootDocument({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <script dangerouslySetInnerHTML={{ __html: THEME_INIT_SCRIPT }} />
        <HeadContent />
      </head>
      <body className="font-body antialiased wrap-anywhere">
        <PostHogProvider>
          <TanStackQueryProvider>
            <AuthProvider baseURL={env.VITE_AUTH_API_URL}>
              <AuthContextUpdater>
                {children}
                <TanStackDevtools
                  config={{
                    position: 'bottom-right',
                  }}
                  plugins={[
                    {
                      name: 'Tanstack Router',
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
  )
}
