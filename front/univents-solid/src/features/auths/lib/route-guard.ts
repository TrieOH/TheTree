import { redirect } from '@tanstack/solid-router';
import type { RouterSession } from '@trieoh/front-core-solid';
import { ownProfileQueryOptions } from '@/features/profile/api';
import { appQueryClient } from '@/shared/lib/query-client';

type GuardArgs = {
  context: { session?: RouterSession };
  location: { href: string; pathname: string; search: Record<string, unknown> };
};

export function requireAuth({ context, location }: GuardArgs) {
  if (context.session?.isAuthenticated !== true) {
    throw redirect({ to: '/auth', search: { redirect: location.href } });
  }
}

export async function requireConfiguredProfile({ context, location }: GuardArgs) {
  if (context.session?.isAuthenticated !== true) return;

  const actorId = context.session.service?.profile()?.id;
  if (!actorId) return;

  const profile = await appQueryClient.fetchQuery(
    ownProfileQueryOptions(actorId, context.session.service),
  );
  const profileExists = Boolean(profile);

  if (!profileExists && location.pathname !== "/profile/setup" && location.pathname !== "/auth/verify") {
    throw redirect({
      to: "/profile/setup",
      search: { returnTo: location.href },
    });
  }

  if (profileExists && location.pathname === "/profile/setup") {
    throw redirect({ to: "/profile", search: { tab: "about" } });
  }
}

export function requireGuest({ context, location }: GuardArgs) {
  if (context.session?.isAuthenticated === true) {
    const redirectTo =
      typeof location.search.redirect === 'string'
        ? location.search.redirect
        : '/profile';
    throw redirect({ to: redirectTo });
  }
}
