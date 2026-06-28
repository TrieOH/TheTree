import { createRouter } from '@tanstack/react-router'
import { setupRouterSsrQueryIntegration } from '@tanstack/react-router-ssr-query'
import { getContext } from './app/providers/tanstack-query/RootProvider'

// Import the generated route tree
import { routeTree } from './routeTree.gen'

// Create a new router instance
export const getRouter = () => {
  const context = getContext()

  const router = createRouter({
    routeTree,
    context: { ...context, auth: undefined },
    defaultPreload: 'intent',
  })

  setupRouterSsrQueryIntegration({ router, queryClient: context.queryClient })

  return router
}
