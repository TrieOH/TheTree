import { redirect } from '@tanstack/solid-router';
import type { RouterSession } from '@trieoh/front-core-solid';
import { ownProfileQueryOptions } from '@/features/profile/api';
import type { ActorProfile } from '@/features/profile/model/profile-data';
import { appQueryClient } from '@/shared/lib/query-client';
import { peekQuery } from '@/shared/lib/route-data';

/** How long a navigation waits to learn whether the profile exists. */
const PROFILE_CHECK_WAIT_MS = 1_200;

type GuardArgs = {
  context: { session?: RouterSession };
  location: { href: string; pathname: string; search: Record<string, unknown> };
};

export function requireAuth({ context, location }: GuardArgs) {
  if (context.session?.isAuthenticated !== true) {
    throw redirect({ to: '/auth', search: { redirect: location.href } });
  }
}

/**
 * Sends a signed-in user without a profile to setup. It runs before every
 * route, so it must never be what takes the app down: the answer comes from the
 * cache when there is one, and when the profile service is slow or failing the
 * guard fails open and lets the navigation through. Only a confirmed "no
 * profile" (the query resolves to `null`) redirects.
 */
export async function requireConfiguredProfile({ context, location }: GuardArgs) {
  if (context.session?.isAuthenticated !== true) return;

  const actorId = context.session.service?.profile()?.id;
  if (!actorId) return;

  const profile = await peekQuery<ActorProfile | null>(
    appQueryClient,
    ownProfileQueryOptions(actorId, context.session.service),
    PROFILE_CHECK_WAIT_MS,
  );
  if (profile === undefined) return;

  const profileExists = profile !== null;

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
