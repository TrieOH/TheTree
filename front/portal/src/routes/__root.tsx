import { HeadContent, Outlet, createRootRoute } from "@tanstack/solid-router";
import { AppShell } from "@/widgets/shell/AppShell";
import { NotFoundPage, RouteError } from "@/widgets/shell/RouteStatus";

export const Route = createRootRoute({
  head: () => ({ meta: [{ title: "TrieOH dev" }] }),
  component: () => (
    <>
      <HeadContent />
      <AppShell>
        <Outlet />
      </AppShell>
    </>
  ),
  notFoundComponent: NotFoundPage,
  errorComponent: RouteError,
});
