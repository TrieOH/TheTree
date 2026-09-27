import { redirect } from '@tanstack/solid-router';
import type { RouterSession } from '@trieoh/front-core-solid';

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

  const response = await context.session.service.getActorProfile(actorId);
  const profileExists = response.success && Boolean(response.data);

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
