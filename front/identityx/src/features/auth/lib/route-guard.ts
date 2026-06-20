import { type AnySchema, type ParsedLocation, redirect } from '@tanstack/react-router'
import type { useAuth } from '@trieoh/identityx-sdk-ts/react';

const SETUP_DONE_KEY = 'trieoh_setup_done';

interface BeforeLoadArgs {
  location: ParsedLocation<AnySchema>;
  context: { auth?: ReturnType<typeof useAuth> }
}

export const requireAuth = ({ context, location }: BeforeLoadArgs) => {
  if (context.auth?.isAuthenticated === false) {
    throw redirect({
      to: '/auth',
      search: { redirect: location.pathname, }
    })
  }
}

export const requireGuest = ({ context, location }: BeforeLoadArgs) => {
  if (context.auth?.isAuthenticated) {
    throw redirect({
      to: '/admin',
      search: { redirect: location.pathname, }
    })
  }
}

/**
 * Cached check: returns true if setup is complete, false if not.
 * Uses localStorage with a permanent marker – once setup is done it never asks again.
 */
function isSetupCached(): boolean {
  if (typeof window === "undefined") return false;
  return localStorage.getItem(SETUP_DONE_KEY) === "true";
}

function markSetupComplete(): void {
  if (typeof window === "undefined") return;
  localStorage.setItem(SETUP_DONE_KEY, "true");
}

/**
 * Guard that redirects to /auth/setup when setup has NOT been done yet.
 * Skips the check for the /auth/setup route itself.
 *
 * The result is cached in localStorage permanently (setup is a one-time event),
 * so the API is called at most once per browser.
 */
export async function requireSetup({ context, location }: BeforeLoadArgs): Promise<void> {
  // Skip the check when we're already on the setup page
  if (location.pathname === '/auth/setup') return;

  // Cache hit – setup is already done, allow navigation
  if (isSetupCached()) return;

  // Cache miss – call the API
  const authService = context.auth?.auth;
  if (!authService) return;
  const res = await authService.isSetupDone();
  console.log(res)
  if (res.success) throw redirect({ to: '/auth/setup' }); // 200 = setup is needed → redirect
  markSetupComplete();
}

/**
 * Guard for the /auth/setup route itself.
 * If setup is already done, redirect to /auth.
 */
export async function requireSetupNotDone({ context }: BeforeLoadArgs): Promise<void> {
  // Cache hit – setup is done, redirect away
  if (isSetupCached()) throw redirect({ to: '/auth' });

  // Cache miss – call the API
  const authService = context.auth?.auth;
  if (!authService) return;
  const res = await authService.isSetupDone();
  if (!res.success) {
    // 503 = setup already done
    markSetupComplete();
    throw redirect({ to: '/auth' });
  }
  // 200 = setup still needed, allow the page
}