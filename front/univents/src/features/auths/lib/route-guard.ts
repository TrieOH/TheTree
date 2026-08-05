import { redirect } from "@tanstack/react-router";
import {
  requireAuth as requireAuthCore,
  requireGuest as requireGuestCore,
} from "@trieoh/front-core";

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
        throw redirect({ to: "/profile" });
      }),
  });
}
