import { redirect } from "@tanstack/react-router"
import type { AnySchema, ParsedLocation } from "@tanstack/react-router";
import {
  requireAuth as requireAuthBase,
  requireGuest as requireGuestBase,
} from "@trieoh/front-core"

type GuardArgs = Parameters<typeof requireAuthBase>[0]

function redirectToLogin(location: ParsedLocation<AnySchema>) {
  throw redirect({
    to: '/',
    search: { redirect: location.pathname },
  })
}

function redirectToAdmin() {
  throw redirect({
    to: '/admin',
  })
}

export const requireAuth = (args: GuardArgs) =>
  requireAuthBase(args, { onRedirect: redirectToLogin })

export const requireGuest = (args: Parameters<typeof requireGuestBase>[0]) =>
  requireGuestBase(args, { onRedirect: redirectToAdmin })
