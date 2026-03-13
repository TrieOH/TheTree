import type { QueryClient } from '@tanstack/react-query'
import {
  createRootRouteWithContext,
  HeadContent,
  Scripts,
} from '@tanstack/react-router'
import { AuthProvider, type useAuth } from '@soramux/node-auth-sdk/react'
import { AuthSynchronizer } from '@/app/providers/auth/RouterAuthSync'
import { RouteComponentTemplate, type RouteStaticConfigI } from '@/app/model/route-types'
import appCss from '../styles.css?url'
import Header from '@/widgets/header/ui/Header'
import { Toaster } from 'sonner'
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
      { title: 'TrieAuth' },
    ],
    links: [{ rel: 'stylesheet', href: appCss }],
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
        <AuthProvider baseURL={env.VITE_API_URL} isClient={false}>
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
