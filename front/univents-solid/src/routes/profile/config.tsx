import { createFileRoute, Link } from "@tanstack/solid-router";

export const Route = createFileRoute("/profile/config")({
  head: () => ({ meta: [{ title: "Configurações do perfil - Univents" }] }),
  component: ProfileConfigPage,
});

function ProfileConfigPage() {
  return (
    <main class="mx-auto min-h-dvh max-w-3xl px-4 py-8 md:px-6">
      <div class="rounded-md border border-border bg-card p-5 shadow-md shadow-foreground/5">
        <h1 class="text-xl font-semibold">Configurações do perfil</h1>
        <p class="mt-2 text-sm text-muted-foreground">
          As configurações do perfil estarão disponíveis em breve.
        </p>
        <Link
          to="/profile"
          search={{ tab: "about" }}
          class="mt-5 inline-flex h-10 items-center rounded-md bg-primary px-4 text-sm text-primary-foreground"
        >
          Voltar ao perfil
        </Link>
      </div>
    </main>
  );
}
