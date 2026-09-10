import { redirect } from '@tanstack/solid-router';
import type { RouterAuth } from '@trieoh/front-core/solid';

type GuardArgs = {
  context: { auth?: RouterAuth };
  location: { href: string; pathname: string; search: Record<string, unknown> };
};

export function requireAuth({ context, location }: GuardArgs) {
  if (context.auth?.isAuthenticated !== true) {
    throw redirect({ to: '/auth', search: { redirect: location.href } });
  }
}

export function requireGuest({ context, location }: GuardArgs) {
  if (context.auth?.isAuthenticated === true) {
    const redirectTo =
      typeof location.search.redirect === 'string'
        ? location.search.redirect
        : '/profile';
    throw redirect({ to: redirectTo });
  }
}
