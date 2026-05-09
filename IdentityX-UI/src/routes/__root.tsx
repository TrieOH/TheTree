import type { QueryClient } from '@tanstack/react-query'
import {
  createRootRouteWithContext,
  HeadContent,
  Scripts,
} from '@tanstack/react-router'
import { AuthProvider, type useAuth } from '@trieoh/identityx-sdk-ts/react'
import { AuthSynchronizer } from '@/app/providers/auth/RouterAuthSync'
import { RouteComponentTemplate, type RouteStaticConfigI } from '@/app/model/route-types'
import appCss from '../styles.css?url'
import Header from '@/widgets/header/ui/Header'
import { Toaster } from '@/shared/ui/shadcn/sonner'
import { env } from '@/env'

interface MyRouterContext {
  queryClient: QueryClient
  auth?: ReturnType<typeof useAuth>
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      { charSet: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { title: env.VITE_APP_TITLE ?? "IdentityX" },
      { name: 'apple-mobile-web-app-title', content: env.VITE_APP_TITLE ?? 'IdentityX' },
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
  notFoundComponent: () => { return (<p>This page doesn't exist!</p>) },
  staticData: { components: RouteComponentTemplate }
})

function RootDocument({ children }: { children: React.ReactNode }) {    
  return (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body className='min-w-xs' suppressHydrationWarning>
        <AuthProvider baseURL={env.VITE_API_URL} isProjectMode={false}>
          <AuthSynchronizer>
            {/* <PHProvider> */}
              <Header />
              {children}
            {/* </PHProvider> */}
          </AuthSynchronizer>
        </AuthProvider>
        <Toaster />
        <Scripts />
      </body>
    </html>
  )
}

declare module '@tanstack/react-router' {
  interface StaticDataRouteOption {
    components: RouteStaticConfigI
  }
}
