import { createRouter as createTanStackRouter } from "@tanstack/react-router";
import "@/shared/lib/api/fetch";
import { setupRouterSsrQueryIntegration } from "@tanstack/react-router-ssr-query";
import { initBrowserTracing } from "@trieoh/front-core/tracing/browser";
import { getContext } from "./integrations/tanstack-query/root-provider";
import { routeTree } from "./routeTree.gen";
export function getRouter() {
  initBrowserTracing();

  const context = getContext();

  const router = createTanStackRouter({
    routeTree,

    context: {
      ...context,
      auth: undefined,
    },

    scrollRestoration: true,
    defaultPreload: "intent",
    defaultPreloadStaleTime: 0,
  });

  setupRouterSsrQueryIntegration({ router, queryClient: context.queryClient });

  return router;
}

declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof getRouter>;
  }
}
