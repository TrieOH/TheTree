import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/privacy")({
  component: PrivacyPage,
});

function PrivacyPage() {
  return (
    <main className="min-h-screen bg-background pb-28">
      <section className="border-b border-border/60 bg-linear-to-b from-muted/40 via-background to-background">
        <div className="mx-auto max-w-4xl px-4 py-12 md:px-6 md:py-16">
          <p className="text-xs uppercase tracking-[0.24em] text-muted-foreground">
            Univents
          </p>
          <h1 className="mt-2 text-3xl font-semibold tracking-tight md:text-5xl">
            Privacidade
          </h1>
          <p className="mt-4 max-w-2xl text-sm leading-relaxed text-muted-foreground md:text-base">
            Explicamos quais dados coletamos, como usamos as informações e quais
            controles você tem sobre sua conta.
          </p>
        </div>
      </section>
    </main>
  );
}
