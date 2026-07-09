import { createFileRoute } from '@tanstack/react-router'
import { Mail, Phone } from 'lucide-react'

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

      <div className="mx-auto max-w-4xl px-4 py-10 md:px-6 md:py-14">
        <div className="grid gap-4 md:grid-cols-2">
          {[
            { icon: Mail, title: 'E-mail', description: 'suporte@univents.com' },
            { icon: Phone, title: 'Telefone', description: '+55 (11) 0000-0000' },
          ].map((item) => (
            <article key={item.title} className="rounded-2xl border border-border/60 bg-card p-6 shadow-sm">
              <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-muted/70">
                <item.icon className="size-5 text-muted-foreground" />
              </div>
              <h2 className="text-base font-semibold">{item.title}</h2>
              <p className="mt-2 text-sm leading-relaxed text-muted-foreground">{item.description}</p>
            </article>
          ))}
        </div>

        <div className="mt-6 rounded-2xl border border-border/60 bg-muted/30 p-6 shadow-sm md:mt-8 md:p-8">
          <p className="text-xs uppercase tracking-[0.24em] text-muted-foreground">Horário</p>
          <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
            Segunda a sexta, das 9h às 18h. Pedidos e dúvidas enviadas fora do horário são respondidos no próximo dia útil.
          </p>
        </div>
      </div>
    </main>
  )
}
