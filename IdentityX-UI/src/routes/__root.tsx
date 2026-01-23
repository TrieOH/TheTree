import {
  HeadContent,
  Scripts,
  createRootRouteWithContext,
  useMatches,
} from '@tanstack/react-router'

import appCss from '../styles.css?url'

import type { QueryClient } from '@tanstack/react-query'
import Header from '@/components/Header'
import { RouteComponentTemplate, type RouteStaticConfigI } from '@/types/route-types'

interface MyRouterContext {
  queryClient: QueryClient
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
      <body className='min-w-xs'>
        <Header {...routeConfig.components.header}/>
        {children}
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