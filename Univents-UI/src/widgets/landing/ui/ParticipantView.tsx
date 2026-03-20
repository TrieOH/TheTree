import { motion } from 'motion/react'
import { FAQSection } from './FAQSection'
import type { EventI } from '@/features/events/model'
import { EventCard } from '@/features/events/ui/EventCard'

const mockEvents: EventI[] = [
  {
    id: '1',
    owner_id: 'user-1',
    organization_id: null,
    goauth_scope_id: 'scope-1',
    name: 'Show ao Vivo',
    acronym: null,
    slug: 'show-ao-vivo',
    tagline: 'Uma noite inesquecível',
    description: 'O melhor show do ano com artistas consagrados',
    is_series: false,
    editions_count: 1,
    logo_url: null,
    banner_url: null,
    has_gallery: false,
    gallery_urls: [],
    contact_email: 'contato@show.com',
    social_links: null,
    status: 'active',
    created_by: 'user-1',
    created_at: '2026-03-15T10:00:00Z',
    updated_at: '2026-03-15T10:00:00Z',
    deleted_at: null,
  },
  {
    id: '2',
    owner_id: 'user-2',
    organization_id: 'org-1',
    goauth_scope_id: 'scope-2',
    name: 'Festival de Jazz',
    acronym: 'FJ2026',
    slug: 'festival-jazz-2026',
    tagline: 'Música boa em todo lugar',
    description: 'O maior festival de jazz da região',
    is_series: true,
    editions_count: 5,
    logo_url: null,
    banner_url: null,
    has_gallery: true,
    gallery_urls: [],
    contact_email: 'jazz@festival.com',
    social_links: { instagram: '@jazzfest' },
    status: 'active',
    created_by: 'user-2',
    created_at: '2026-03-22T14:00:00Z',
    updated_at: '2026-03-22T14:00:00Z',
    deleted_at: null,
  },
  {
    id: '3',
    owner_id: 'user-3',
    organization_id: null,
    goauth_scope_id: 'scope-3',
    name: 'Tech Conference',
    acronym: 'TC26',
    slug: 'tech-conference-2026',
    tagline: 'O futuro da tecnologia',
    description: 'Conferência sobre inovação e tecnologia',
    is_series: false,
    editions_count: 1,
    logo_url: null,
    banner_url: null,
    has_gallery: false,
    gallery_urls: [],
    contact_email: 'tech@conf.com',
    social_links: null,
    status: 'active',
    created_by: 'user-3',
    created_at: '2026-04-05T09:00:00Z',
    updated_at: '2026-04-05T09:00:00Z',
    deleted_at: null,
  },
]

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
  return (
    <div className="max-w-5xl mx-auto space-y-20 md:space-y-32">
      {/* Hero */}
      <section className="text-center max-w-2xl mx-auto space-y-4 md:space-y-6 px-2">
        <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
          Encontre os melhores eventos perto de você, compre ingressos com segurança
          e aproveite cada momento sem complicação.
        </p>
        <div className="flex flex-col sm:flex-row gap-3 md:gap-4 justify-center pt-2 md:pt-4">
          <button className="px-5 py-2.5 md:px-6 md:py-3 bg-primary text-primary-foreground rounded-full text-sm font-medium hover:bg-primary/90 transition-colors">
            Criar conta grátis
          </button>
          <button className="px-5 py-2.5 md:px-6 md:py-3 border border-border text-foreground rounded-full text-sm font-medium hover:border-foreground/50 transition-colors">
            Explorar eventos
          </button>
        </div>
      </section>

      {/* Grid de Eventos */}
      <section className="space-y-6 md:space-y-8">
        <div className="flex justify-between items-end border-b border-border pb-3 md:pb-4">
          <h2 className="text-lg md:text-2xl font-medium text-foreground">Em Alta</h2>
          <a href="#" className="text-xs md:text-sm text-muted-foreground hover:text-foreground transition-colors">
            Ver Todos
          </a>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-6">
          {mockEvents.map((event, idx) => (
            <EventCard key={event.id} event={event} index={idx} />
          ))}
        </div>
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

      {/* Como Funciona */}
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

      {/* FAQ com Collapsible */}
      <section className="max-w-2xl mx-auto space-y-6 md:space-y-8">
        <div className="text-center space-y-2">
          <h2 className="text-2xl md:text-3xl font-semibold text-foreground">Perguntas frequentes</h2>
          <p className="text-sm text-muted-foreground">Tire suas dúvidas</p>
        </div>

        <FAQSection items={faqs} />
      </section>

      {/* CTA Final */}
      <section className="bg-muted rounded-2xl md:rounded-3xl p-6 md:p-12 lg:p-16 text-center space-y-4 md:space-y-6">
        <h2 className="text-xl md:text-3xl font-semibold text-foreground">
          Seu próximo evento te espera
        </h2>
        <p className="text-sm md:text-base text-muted-foreground max-w-md mx-auto">
          Crie sua conta gratuita e descubra eventos incríveis acontecendo perto de você.
        </p>
        <div className="flex flex-col sm:flex-row gap-3 justify-center pt-2">
          <button className="px-5 py-2.5 md:px-6 md:py-3 bg-primary text-primary-foreground rounded-full text-sm font-medium hover:bg-primary/90 transition-colors">
            Criar conta grátis
          </button>
          <button className="px-5 py-2.5 md:px-6 md:py-3 border border-border text-foreground rounded-full text-sm font-medium hover:border-foreground/50 transition-colors">
            Explorar eventos
          </button>
        </div>
        <p className="text-xs text-muted-foreground/70">100% gratuito. Sem cartão de crédito.</p>
      </section>
    </div>
  )
}