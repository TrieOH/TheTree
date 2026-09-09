import { Link, createFileRoute } from '@tanstack/solid-router';
import Logo from '@/shared/ui/Logo';

export const Route = createFileRoute('/')({
  head: () => ({ meta: [{ title: 'Univents' }] }),
  component: HomePage,
});

function HomePage() {
  return (
    <main
      class="flex min-h-[70vh] flex-col items-center justify-center gap-8 px-6 text-center"
    >
      <Logo variant="complete" priority imgClassName="max-h-24 w-auto" />
      <div class="space-y-2">
        <h1 class="font-heading text-3xl font-semibold">Bem-vindo à Univents</h1>
        <p class="text-muted-foreground">Encontre e participe dos melhores eventos.</p>
      </div>
      <Link
        to="/auth"
        class="rounded-lg bg-primary px-5 py-3 font-medium text-primary-foreground transition-opacity hover:opacity-90"
      >
        Entrar
      </Link>
    </main>
  );
}
