import { createFileRoute } from '@tanstack/solid-router';
import { ModernAuth } from '@trieoh/identityx-sdk-ts-solid';
import Logo from '../shared/ui/Logo';
import z from 'zod';

const authSearchSchema = z.object({
  redirect: z.string().optional().catch(""),
});


export const Route = createFileRoute('/auth')({
  head: () => ({ meta: [{ title: 'Entrar - Univents' }] }),
  validateSearch: (search) => authSearchSchema.parse(search),
  component: AuthPage,
});

function AuthPage() {
  // const search = Route.useSearch();
  return (
    <div class="relative [&>main]:py-16 [&>main]:pt-52 md:[&>main]:pt-56">
      <div class="absolute left-1/2 top-16 md:top-20 z-20 w-32 md:w-40 -translate-x-1/2">
        <Logo
          variant="complete"
          priority
          imgClassName="h-auto max-h-16 md:max-h-20"
        />
      </div>

      <ModernAuth
        initialView="signin"
        onLoginSuccess={async () => undefined}
        onSignUpSuccess={async () => undefined}
        onFailed={async (message, trace) => console.error(message, trace)}
      />
    </div>
  );
}
