import { type ParsedLocation, redirect } from '@tanstack/react-router'
import { type useAuth } from '@trieoh/node-auth-sdk/react';

interface BeforeLoadArgs {
  location: ParsedLocation<{}>;
  context: { auth?: ReturnType<typeof useAuth> }
}

export const requireAuth = ({ context, location }: BeforeLoadArgs) => {
  if (!context.auth?.isAuthenticated) {
    throw redirect({
      to: '/auth',
      search: { redirect: location.pathname, }
    })
  }
}

export const requireGuest = ({ context, location }: BeforeLoadArgs) => {
  if (context.auth?.isAuthenticated) {
    throw redirect({
      to: '/projects',
      search: { redirect: location.pathname, }
    })
  }
}