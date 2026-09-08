import { HeadContent, Outlet, createRootRoute } from '@tanstack/solid-router';
import { OverlayScrollbar } from "@/widgets/ui/OverlayScrollbar";
import { AuthContextUpdater } from '@trieoh/front-core/solid';
import { Toaster } from '@/shared/ui/toast';

// The root route: the site-wide layout every route renders inside, plus the
// not-found boundary. <HeadContent /> renders whatever the matched routes
// declare in their `head` options (titles here).
export const Route = createRootRoute({
  head: () => ({ meta: [{ title: 'Univents' }] }),
  component: () => (
    <>
      <HeadContent />
      <AuthContextUpdater>
        <Outlet />
      </AuthContextUpdater>
      <OverlayScrollbar />
      <Toaster />
    </>
  ),
  notFoundComponent: () => (
    <main>
      <h1>Page Not Found</h1>
      <p>
        Visit{' '}
        <a href="https://docs.solidjs.com" target="_blank" rel="noreferrer">
          docs.solidjs.com
        </a>{' '}
        to learn how to build Solid apps.
      </p>
    </main>
  ),
  errorComponent: (props) => (
    <main>
      <h1>Algo deu errado</h1>
      <p>{props.error instanceof Error ? props.error.message : 'Erro inesperado.'}</p>
    </main>
  ),
});
