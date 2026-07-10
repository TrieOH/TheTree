import { motion } from 'motion/react'
import { useNavigate, useRouteContext } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { FAQSection } from './FAQSection'
import { EventCard } from '@/features/events/ui/EventCard'
import { eventsQueryOptions } from '@/features/events/api'
import { Button } from '@/shared/ui/shadcn/button'
import { CardSkeleton, CardsStateGrid, EmptyState } from '@trieoh/ui-base'
import { CalendarX } from 'lucide-react'

const features = [
  {
    title: "Encontre eventos facilmente",
    desc: "Pesquise por categoria, data ou localidade e descubra eventos que combinam com você."
  },
  {
    title: "Compra segura de ingressos",
    desc: "Adquira seus ingressos com pagamento protegido e receba a confirmação na hora."
  },
  {
    title: "Ingresso digital no celular",
    desc: "Acesse seu ingresso direto pelo app. Sem papel, sem filas — basta apresentar o QR Code."
  },
  {
    title: "Notificações do evento",
    desc: "Receba lembretes, atualizações e informações importantes sobre seus eventos."
  },
  {
    title: "Avalie e recomende",
    desc: "Compartilhe sua experiência e ajude outros participantes a escolherem os melhores eventos."
  },
  {
    title: "Garantia de reembolso",
    desc: "Se o evento for cancelado, você recebe o reembolso automaticamente na sua conta."
  },
]

const steps = [
  { num: "01", title: "Crie sua conta", desc: "Cadastre-se em segundos com e-mail ou redes sociais. Totalmente gratuito." },
  { num: "02", title: "Encontre seu evento", desc: "Explore eventos por categoria, local ou data. Compre seu ingresso com segurança." },
  { num: "03", title: "Aproveite a experiência", desc: "Apresente seu QR Code na entrada e curta o evento sem preocupação." },
]

const faqs = [
  {
    question: "Como funciona a taxa sobre vendas?",
    answer: "Não há taxa para participantes. O preço do ingresso é o valor final que você paga. Organizadores pagam uma pequena comissão apenas sobre as vendas realizadas."
  },
  {
    question: "Preciso pagar alguma mensalidade?",
    answer: "Não. Para participantes, o uso é 100% gratuito. Para organizadores, oferecemos planos gratuitos e pagos, mas você só paga se escolher recursos avançados."
  },
  {
    question: "Quais tipos de eventos posso gerenciar?",
    answer: "Qualquer tipo: shows, festivais, conferências, workshops, esportivos, corporativos, sociais. Não há limitação de categoria ou tamanho."
  },
  {
    question: "Como recebo o dinheiro das vendas?",
    answer: "O repasse é feito diretamente para sua conta bancária em até 2 dias úteis após a transação. Para valores maiores, oferecemos antecipação de recebíveis."
  },
  {
    question: "Posso cancelar a qualquer momento?",
    answer: "Sim. Não há contratos de fidelidade ou período mínimo. Cancele quando quiser sem taxas de rescisão."
  },
  {
    question: "A plataforma oferece suporte para check-in no evento?",
    answer: "Sim. Oferecemos app de check-in com leitura de QR Code, lista de convidados offline e controle de entrada em tempo real."
  },
]

export function ParticipantView() {
  const isAuthenticated = useRouteContext({ from: "/" }).auth?.isAuthenticated ?? false
  const navigate = useNavigate()

  const { data: events = [], isLoading } = useQuery(eventsQueryOptions());

  const handleGetStarted = () => {
    if (isAuthenticated) void navigate({ to: '/events' })
    else void navigate({ to: '/auth' })
  }

  const ctaLabel = isAuthenticated ? 'Continuar explorando' : 'Entrar para ver eventos'
  const footerCtaLabel = isAuthenticated ? 'Ver eventos' : 'Criar conta grátis'

  return (
    <div className="max-w-5xl mx-auto space-y-20 md:space-y-32">
      {/* Hero */}
      <section className="text-center max-w-2xl mx-auto space-y-4 md:space-y-6 px-2">
        <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
          Encontre os melhores eventos perto de você, compre ingressos com segurança
          e aproveite cada momento sem complicação.
        </p>
        <div className="flex justify-center pt-2 md:pt-4">
          <Button
            onClick={handleGetStarted}
            size="xl"
            className="rounded-full px-6"
          >
            {ctaLabel}
          </Button>
        </div>
      </section>

      <section id="trending" className="space-y-6 md:space-y-8 scroll-mt-20">
        <div className="space-y-2">
          <p className="text-xs uppercase tracking-[0.24em] text-muted-foreground">Em alta agora</p>
          <h2 className="text-lg md:text-2xl font-medium text-foreground">Eventos em destaque</h2>
          <p className="max-w-2xl text-sm md:text-base text-muted-foreground">
            Uma seleção curta dos eventos mais recentes e relevantes para você começar sem ruído.
          </p>
        </div>

        <CardsStateGrid
          items={events.slice(0, 4)}
          isLoading={isLoading}
          count={4}
          columns={4}
          className="gap-4 xl:gap-6"
          renderEmpty={
            <EmptyState
              icon={CalendarX}
              eyebrow="Em destaque"
              title="Nenhum evento em destaque agora"
              description="Assim que novos eventos forem publicados, eles aparecem aqui."
              className="sm:col-span-2 xl:col-span-4"
            />
          }
          renderSkeleton={() => (
            <CardSkeleton
              showDescription={false}
              rows={0}
              mediaClassName="aspect-[4/3]"
              footerClassName="space-y-2"
            />
          )}
          renderItem={(event: (typeof events)[number], idx: number) => (
            <EventCard
              key={event.id}
              event={event}
              index={idx}
              className="min-w-0 border border-border/70 bg-card shadow-sm"
            />
          )}
        />
      </section>

      {/* Features */}
      <section className="space-y-8 md:space-y-12">
        <div className="text-center space-y-2">
          <p className="text-xs md:text-sm text-muted-foreground uppercase tracking-wider">Para participantes</p>
          <h2 className="text-2xl md:text-4xl font-semibold text-foreground">
            Tudo para curtir seus<br className="hidden md:block" /> eventos favoritos.
          </h2>
        </div>

        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-x-6 gap-y-8 md:gap-x-8 md:gap-y-12">
          {features.map((f, idx) => (
            <motion.div
              key={idx}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: idx * 0.1 }}
              className="space-y-2"
            >
              <h3 className="text-base md:text-lg font-medium text-foreground">{f.title}</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">{f.desc}</p>
            </motion.div>
          ))}
        </div>
      </section>

      <section className="space-y-8 md:space-y-12">
        <div className="text-center">
          <h2 className="text-2xl md:text-4xl font-semibold text-foreground mb-2">Como funciona</h2>
          <p className="text-muted-foreground">Simples assim</p>
        </div>

        <div className="grid md:grid-cols-3 gap-6 md:gap-8">
          {steps.map((step, idx) => (
            <motion.div
              key={idx}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: idx * 0.15 }}
              className="space-y-3 md:space-y-4"
            >
              <span className="text-3xl md:text-4xl font-semibold text-muted">{step.num}</span>
              <h3 className="text-lg md:text-xl font-medium text-foreground">{step.title}</h3>
              <p className="text-sm md:text-base text-muted-foreground leading-relaxed">{step.desc}</p>
            </motion.div>
          ))}
        </div>
      </section>

      {/* FAQ */}
      <section className="max-w-2xl mx-auto space-y-6 md:space-y-8">
        <div className="text-center space-y-2">
          <h2 className="text-2xl md:text-3xl font-semibold text-foreground">Perguntas frequentes</h2>
          <p className="text-sm text-muted-foreground">Tire suas dúvidas</p>
        </div>

        <FAQSection items={faqs} />
      </section>

      <section className="bg-muted rounded-2xl md:rounded-3xl p-6 md:p-12 lg:p-16 text-center space-y-4 md:space-y-6">
        <h2 className="text-xl md:text-3xl font-semibold text-foreground">
          Seu próximo evento te espera
        </h2>
        <p className="text-sm md:text-base text-muted-foreground max-w-md mx-auto">
          Crie sua conta gratuita e descubra eventos incríveis acontecendo perto de você.
        </p>
        <div className="flex justify-center pt-2">
          <button
            onClick={handleGetStarted}
            className="px-5 py-2.5 md:px-6 md:py-3 bg-primary text-primary-foreground rounded-full text-sm font-medium hover:bg-primary/90 transition-colors"
          >
            {footerCtaLabel}
          </button>
        </div>
        <p className="text-xs text-muted-foreground/70">100% gratuito. Sem cartão de crédito.</p>
      </section>
    </div>
  )
}
