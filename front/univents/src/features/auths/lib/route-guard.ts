import type { QueryClient } from "@tanstack/react-query";
import { redirect } from "@tanstack/react-router";
import {
  requireAuth as requireAuthCore,
  requireGuest as requireGuestCore,
} from "@trieoh/front-core";
import type { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { profileKeys } from "@/features/profile/api/query-keys";

type AuthGuardArgs = Parameters<typeof requireAuthCore>[0];
type AuthGuardOptions = Parameters<typeof requireAuthCore>[1];
type GuestGuardArgs = Parameters<typeof requireGuestCore>[0];
type GuestGuardOptions = Parameters<typeof requireGuestCore>[1];

export function requireAuth(
  args: AuthGuardArgs,
  options: AuthGuardOptions = {},
) {
  return requireAuthCore(args, {
    ...options,
    onRedirect:
      options.onRedirect ??
      ((location) => {
        throw redirect({
          to: "/auth",
          search: { redirect: location.href },
        });
      }),
  });
}

export function requireGuest(
  args: GuestGuardArgs,
  options: GuestGuardOptions = {},
) {
  return requireGuestCore(args, {
    ...options,
    onRedirect:
      options.onRedirect ??
      (() => {
        throw redirect({ to: "/profile", search: { tab: "about" } });
      }),
  });
}

export async function requireConfiguredProfile({
  context,
  location,
}: {
  context: {
    auth?: ReturnType<typeof useAuth>;
    queryClient: QueryClient;
  };
  location: { pathname: string };
}) {
  if (context.auth?.isAuthenticated !== true) {
    return;
  }

  const actorId = context.auth.auth.profile()?.id;
  if (!actorId) return;

  const profile = await context.queryClient.ensureQueryData({
    queryKey: profileKeys.detail(actorId),
    queryFn: async () => {
      const response = await context.auth?.auth.getActorProfile(actorId);
      return response?.success ? (response.data ?? null) : null;
    },
  });

  if (!profile && location.pathname !== "/profile/setup") {
    throw redirect({ to: "/profile/setup", search: { returnTo: undefined } });
  }
  if (profile && location.pathname === "/profile/setup") {
    throw redirect({ to: "/profile", search: { tab: "about" } });
  }
}
