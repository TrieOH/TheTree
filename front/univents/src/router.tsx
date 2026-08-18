import { createRouter as createTanStackRouter } from "@tanstack/react-router";
import "@/shared/lib/api/fetch";
import { setupRouterSsrQueryIntegration } from "@tanstack/react-router-ssr-query";
import { initBrowserTracing } from "@trieoh/front-core/tracing/browser";
import { env } from "./env";
import { getContext } from "./integrations/tanstack-query/root-provider";
import { routeTree } from "./routeTree.gen";
export function getRouter() {
  initBrowserTracing("univents-web", env.VITE_TRACING_ENABLED, [
    env.VITE_API_URL,
  ]);

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
