import {
  HeadContent,
  Scripts,
  createRootRouteWithContext,
  useMatches,
  useRouter,
} from '@tanstack/react-router'

import appCss from '../styles.css?url'

import type { QueryClient } from '@tanstack/react-query'
import Header from '@/components/Header'
import { RouteComponentTemplate, type RouteStaticConfigI } from '@/types/route-types'
import { AuthProvider, useAuth } from '@trieoh/node-auth-sdk/react'
import { useEffect, useState } from 'react'

interface MyRouterContext {
  queryClient: QueryClient
  auth?: ReturnType<typeof useAuth>
}

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
        title: 'TrieAuth',
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
  notFoundComponent: () => { return (<p>This page doesn't exist!</p>) },
  staticData: { components: RouteComponentTemplate }
})

function RootDocument({ children }: { children: React.ReactNode }) {
  const matches = useMatches();
  const routeConfig = getRouteConfig(matches);
  
  return (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body className='min-w-xs' suppressHydrationWarning>
        <AuthProvider baseURL="http://localhost:8080">
          <AuthSyncronizer>
            {/* <PHProvider> */}
              <Header {...routeConfig.components.header} />
              {children}
            {/* </PHProvider> */}
          </AuthSyncronizer>
        </AuthProvider>
        <Scripts />
      </body>
    </html>
  )
}

function AuthSyncronizer({ children }: { children: React.ReactNode }) {
  const auth = useAuth() 
  const router = useRouter()
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if(router.options.context.auth?.isAuthenticated !== auth.isAuthenticated) {
      router.update({
        context: { ...router.options.context, auth: auth },
      })
      setIsLoading(false);
      router.invalidate()
    }
  }, [auth.isAuthenticated, router])

  if(isLoading) return null; // change to guard
  
  return <>{children}</>
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
