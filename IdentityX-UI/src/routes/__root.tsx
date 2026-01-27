import type { QueryClient } from '@tanstack/react-query'
import {
  createRootRouteWithContext,
  HeadContent,
  Scripts,
  useMatches,
} from '@tanstack/react-router'
import { AuthProvider, type useAuth } from '@trieoh/node-auth-sdk/react'
import { useMemo } from 'react'
import { AuthSynchronizer } from '@/app/providers/auth/RouterAuthSync'
import { RouteComponentTemplate, type RouteStaticConfigI } from '@/shared/types/route-types'
import Header from '@/shared/ui/navigation/header/Header'
import appCss from '../styles.css?url'

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
  const matches = useMatches();
  const routeConfig = useMemo(() => getRouteConfig(matches), [matches]);
  
  return (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body className='min-w-xs' suppressHydrationWarning>
        <AuthProvider baseURL="http://localhost:8080">
          <AuthSynchronizer>
            {/* <PHProvider> */}
              <Header {...routeConfig.components.header} />
              {children}
            {/* </PHProvider> */}
          </AuthSynchronizer>
        </AuthProvider>
        <Scripts />
      </body>
    </html>
  )
}

function getRouteConfig(matches: ReturnType<typeof useMatches>) {
  return matches.reduce((acc, match) => {
    const route = match.staticData;
    if(!route) return acc;
    return route;
  }, { components: RouteComponentTemplate });
}

declare module '@tanstack/react-router' {
  interface StaticDataRouteOption {
    components: RouteStaticConfigI
  }
}
