import { redirect } from '@tanstack/react-router'
import type { AnySchema, ParsedLocation } from '@tanstack/react-router'
import type { useAuth } from '@trieoh/node-auth-sdk/react';

interface BeforeLoadArgs {
  location: ParsedLocation<AnySchema>;
  context: { auth?: ReturnType<typeof useAuth> }
}

export const requireAuth = ({ context, location }: BeforeLoadArgs) => {
  if (context.auth?.isAuthenticated === false) {
    throw redirect({
      to: '/',
      search: { redirect: location.pathname, }
    })
  }
}

export const requireGuest = ({ context, location }: BeforeLoadArgs) => {
  if (context.auth?.isAuthenticated === true) {
    throw redirect({
      to: '/admin',
      search: { redirect: location.pathname, }
    })
  }
}