import { redirect } from "@tanstack/react-router";
import {
  requireAuth as requireAuthCore,
  requireGuest as requireGuestCore,
} from "@trieoh/front-core";

export { requireSetup, requireSetupNotDone } from "@trieoh/front-core";

type AuthGuardArgs = Parameters<typeof requireAuthCore>[0];
type GuestGuardArgs = Parameters<typeof requireGuestCore>[0];

export function requireAuth(args: AuthGuardArgs) {
  return requireAuthCore(args, {
    onRedirect: (location) => {
      throw redirect({
        to: "/auth",
        search: { redirect: location.href },
      });
    },
  });
}

export function requireGuest(args: GuestGuardArgs) {
  return requireGuestCore(args, {
    onRedirect: () => {
      throw redirect({ to: "/admin" });
    },
  });
}
