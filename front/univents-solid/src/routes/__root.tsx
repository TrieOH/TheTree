import { HeadContent, Outlet, createRootRouteWithContext } from '@tanstack/solid-router';
import type { RouterSession } from '@trieoh/front-core-solid';
import { OverlayScrollbar } from "@/widgets/ui/OverlayScrollbar";
import { AuthContextUpdater } from '@trieoh/front-core-solid';
import { Toaster } from '@/shared/ui/toast';
import { NavigationDock } from '@/widgets/ui/NavigationDock';
import { UploadQueueProvider } from '@/features/upload-queue';
import { requireConfiguredProfile } from '@/features/auths/lib/route-guard';
import { NotFoundPage, RootError } from '@/widgets/ui/RouteStatus';
import { PageFrame } from '@/widgets/ui/PageFrame';
import { ServiceStatusBanner } from '@/widgets/ui/ServiceStatusBanner';

// The root route: the site-wide layout every route renders inside, plus the
// not-found and error boundaries. <HeadContent /> renders whatever the matched
// routes declare in their `head` options (titles here).
export const Route = createRootRouteWithContext<{ session?: RouterSession }>()({
  beforeLoad: requireConfiguredProfile,
  head: () => ({ meta: [{ title: 'Univents' }] }),
  component: () => (
    <>
      <HeadContent />
      <AuthContextUpdater>
        <UploadQueueProvider>
          <PageFrame>
            <Outlet />
          </PageFrame>
        </UploadQueueProvider>
      </AuthContextUpdater>
      <OverlayScrollbar />
      <NavigationDock />
      <ServiceStatusBanner />
      <Toaster />
    </>
  ),
  notFoundComponent: NotFoundPage,
  errorComponent: RootError,
});
