import {
  type AnySchema,
  type ParsedLocation,
  redirect,
} from "@tanstack/react-router"

interface BeforeLoadArgs {
  location: ParsedLocation<AnySchema>
  context: {
    auth?: {
      isAuthenticated: boolean
      auth?: {
        isSetupDone: () => Promise<{ success: boolean }>
      }
    }
  }
}

interface RedirectOptions {
  onRedirect?: (location: ParsedLocation<AnySchema>) => never | void
  authRoute?: string
  setupRoute?: string
}

/**
 * Requires the user to be authenticated.
 * Calls the provided redirect handler if not authenticated.
 */
export const requireAuth = (
  { context, location }: BeforeLoadArgs,
  options: RedirectOptions = {},
) => {
  if (context.auth?.isAuthenticated === false) {
    return options.onRedirect?.(location)
  }
}

/**
 * Requires the user to be a guest (not authenticated).
 * Calls the provided redirect handler if already authenticated.
 */
export const requireGuest = (
  { context, location }: BeforeLoadArgs,
  options: RedirectOptions = {},
) => {
  if (context.auth?.isAuthenticated === true) {
    return options.onRedirect?.(location)
  }
}

const SETUP_DONE_KEY = "trieoh_setup_done"

function isSetupCached(): boolean {
  if (typeof window === "undefined") return false
  return localStorage.getItem(SETUP_DONE_KEY) === "true"
}

function markSetupComplete(): void {
  if (typeof window === "undefined") return
  localStorage.setItem(SETUP_DONE_KEY, "true")
}

/**
 * Guard that redirects to `/auth/setup` when setup has not been done yet.
 * The result is cached in localStorage (setup is a one-time event),
 * so the API is called at most once per browser.
 */
export async function requireSetup({
  context,
  location,
}: BeforeLoadArgs, options: RedirectOptions = {}): Promise<void> {
  const setupRoute = options.setupRoute ?? "/auth/setup"
  if (location.pathname === setupRoute) return
  if (isSetupCached()) return

  const authService = context.auth?.auth
  if (!authService) return

  const res = await authService.isSetupDone()
  if (res.success) throw redirect({ to: setupRoute })
  markSetupComplete()
}

/**
 * Guard for the `/auth/setup` route itself.
 * If setup is already done, redirect to `/auth`.
 */
export async function requireSetupNotDone({
  context,
}: BeforeLoadArgs, options: RedirectOptions = {}): Promise<void> {
  const authRoute = options.authRoute ?? "/auth"
  if (isSetupCached()) throw redirect({ to: authRoute })

  const authService = context.auth?.auth
  if (!authService) return

  const res = await authService.isSetupDone()
  if (!res.success) {
    markSetupComplete()
    throw redirect({ to: authRoute })
  }
}
