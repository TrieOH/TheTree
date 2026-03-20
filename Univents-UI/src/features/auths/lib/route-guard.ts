/* eslint-disable @typescript-eslint/only-throw-error */
import { redirect } from '@tanstack/react-router'
import type { AnySchema, ParsedLocation } from '@tanstack/react-router'
import type { useAuth } from '@soramux/node-auth-sdk/react';

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

export const requireGuest = ({ context }: BeforeLoadArgs) => {
  if (context.auth?.isAuthenticated === true) {
    throw redirect({ to: '/' })
  }
}