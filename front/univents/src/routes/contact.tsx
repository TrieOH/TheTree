import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/contact')({
  component: ContactPage,
})

function ContactPage() {
  return (
    <main className="min-h-screen bg-background pb-28 md:pb-40">
      <section className="border-b border-border/60 bg-linear-to-b from-muted/40 via-background to-background">
        <div className="mx-auto max-w-4xl px-4 py-12 md:px-6 md:py-16">
          <p className="text-xs uppercase tracking-[0.24em] text-muted-foreground">Univents</p>
          <h1 className="mt-2 text-3xl font-semibold tracking-tight md:text-5xl">Contato</h1>
          <p className="mt-4 max-w-2xl text-sm leading-relaxed text-muted-foreground md:text-base">
            Fale com a equipe da Univents para suporte, parcerias ou dúvidas sobre a plataforma.
          </p>
        </div>
      </section>
    </main>
  )
}
