import { HeadContent, Outlet, createRootRouteWithContext } from '@tanstack/solid-router';
import type { RouterSession } from '@trieoh/front-core-solid';
import { OverlayScrollbar } from "@/widgets/ui/OverlayScrollbar";
import { AuthContextUpdater } from '@trieoh/front-core-solid';
import { Toaster } from '@/shared/ui/toast';
import { NavigationDock } from '@/widgets/ui/NavigationDock';
import { UploadQueueProvider } from '@/features/upload-queue';
import { requireConfiguredProfile } from '@/features/auths/lib/route-guard';

// The root route: the site-wide layout every route renders inside, plus the
// not-found boundary. <HeadContent /> renders whatever the matched routes
// declare in their `head` options (titles here).
export const Route = createRootRouteWithContext<{ session?: RouterSession }>()({
  beforeLoad: requireConfiguredProfile,
  head: () => ({ meta: [{ title: 'Univents' }] }),
  component: () => (
    <>
      <HeadContent />
      <AuthContextUpdater>
        <UploadQueueProvider>
          <main id="main-content" tabindex="-1">
            <Outlet />
          </main>
        </UploadQueueProvider>
      </AuthContextUpdater>
      <OverlayScrollbar />
      <NavigationDock />
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
