import {
  HeadContent,
  Scripts,
  createRootRouteWithContext,
} from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import { TanStackDevtools } from '@tanstack/react-devtools'
import { AuthProvider } from '@soramux/node-auth-sdk/react'

import StoreDevtools from '../lib/demo-store-devtools'

import PostHogProvider from '../integrations/posthog/provider'

import TanStackQueryProvider from '../integrations/tanstack-query/root-provider'

import TanStackQueryDevtools from '../integrations/tanstack-query/devtools'

import appCss from '../styles.css?url'
import type { useAuth } from '@soramux/node-auth-sdk/react';

import type { QueryClient } from '@tanstack/react-query'
import { env } from '@/env'
import { AuthContextUpdater } from '@/integrations/auth/auth-context-updater'
import { NavigationDock } from '@/widgets/ui/navigation-dock'
import NotFound from '@/widgets/feedback/ui/NotFound'
import { Toaster } from '@/shared/ui/shadcn/sonner'
import WaveSpinnerLoading from '@/shared/ui/loader/WaveSpinnerLoading'

interface MyRouterContext {
  queryClient: QueryClient
  auth?: ReturnType<typeof useAuth>
}

const THEME_INIT_SCRIPT = `(function(){try{var stored=window.localStorage.getItem('theme');var mode=(stored==='light'||stored==='dark'||stored==='auto')?stored:'auto';var prefersDark=window.matchMedia('(prefers-color-scheme: dark)').matches;var resolved=mode==='auto'?(prefersDark?'dark':'light'):mode;var root=document.documentElement;root.classList.remove('light','dark');root.classList.add(resolved);if(mode==='auto'){root.removeAttribute('data-theme')}else{root.setAttribute('data-theme',mode)}root.style.colorScheme=resolved;}catch(e){}})();`

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      {
        charSet: 'utf-8',
      },
      {
        name: 'viewport',
        content: 'width=device-width, initial-scale=1',
      },
      {
        title: 'Univents',
      },
    ],
    links: [
      {
        rel: 'stylesheet',
        href: appCss,
      },
    ],
  }),
  shellComponent: RootDocument,
  notFoundComponent: NotFound
})

function RootDocument({ children }: { children: React.ReactNode }) {
  return (
    <html lang="pt-BR" suppressHydrationWarning>
      <head>
        <script dangerouslySetInnerHTML={{ __html: THEME_INIT_SCRIPT }} />
        <script
          crossOrigin="anonymous"
          src="//unpkg.com/react-scan/dist/auto.global.js"
        />
        <HeadContent />
      </head>
      <body>
        <PostHogProvider>
          <TanStackQueryProvider>
            <AuthProvider
              baseURL={env.VITE_AUTH_API_URL}
              exchangeURL={env.VITE_EXCHANGE_API_URL}
              clientConfig={{
                timeout: 5_000
              }}
              fallback={
                <div className='h-screen w-screen flex items-center justify-center'>
                  <WaveSpinnerLoading text='Carregando...' />
                </div>
              }
            >
              <AuthContextUpdater>
                {children}
                <NavigationDock />
                <TanStackDevtools
                  config={{
                    position: 'bottom-right',
                  }}
                  plugins={[
                    {
                      name: 'Tanstack Router',
                      render: <TanStackRouterDevtoolsPanel />,
                    },
                    StoreDevtools,
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
