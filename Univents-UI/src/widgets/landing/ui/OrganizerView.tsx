import { motion } from 'motion/react'
import { FAQSection } from './FAQSection'
import { cn } from '@/shared/lib/utils'

const features = [
  {
    title: "Planejamento completo",
    desc: "Organize cronogramas, prazos e checklists com visão completa de cada etapa do evento."
  },
  {
    title: "Gestão de convidados",
    desc: "Controle listas de presença, convites e confirmações em tempo real."
  },
  {
    title: "Ingressos e inscrições",
    desc: "Venda ingressos, gerencie inscrições e acompanhe a receita do seu evento."
  },
  {
    title: "Fornecedores e logística",
    desc: "Centralize contratos, prazos e comunicação com todos os fornecedores."
  },
  {
    title: "Relatórios pós-evento",
    desc: "Analise resultados, feedbacks e métricas para melhorar seus próximos eventos."
  },
  {
    title: "Automações inteligentes",
    desc: "Automatize lembretes, e-mails e tarefas repetitivas do fluxo de trabalho."
  },
]

const steps = [
  { num: "01", title: "Configure seu evento", desc: "Defina detalhes, datas, equipe e fornecedores em poucos minutos com nossos templates prontos." },
  { num: "02", title: "Organize e gerencie", desc: "Acompanhe tarefas, convidados e logística em tempo real com visão unificada de tudo." },
  { num: "03", title: "Execute com confiança", desc: "No dia do evento, tenha controle total com check-ins, timeline ao vivo e comunicação integrada." },
]

const faqs = [
  {
    question: "Como funciona a taxa sobre vendas?",
    answer: "Cobramos apenas 8% sobre cada ingresso vendido. Não há taxas de setup, mensalidade mínima ou custos ocultos. Quanto mais você vende, melhores condições podemos oferecer."
  },
  {
    question: "Preciso pagar alguma mensalidade?",
    answer: "Não obrigatoriamente. Nosso plano Starter é gratuito. Planos pagos (Pro e Enterprise) oferecem recursos avançados, mas você só paga se quiser utilizá-los."
  },
  {
    question: "Quais tipos de eventos posso gerenciar?",
    answer: "Qualquer categoria: corporativos, sociais, culturais, esportivos, educacionais. Desde pequenas reuniões até grandes festivais com milhares de participantes."
  },
  {
    question: "Como recebo o dinheiro das vendas?",
    answer: "O repasse é automático para sua conta em até 2 dias úteis. Oferecemos também antecipação de recebíveis para eventos com fluxo de caixa planejado."
  },
  {
    question: "Posso cancelar a qualquer momento?",
    answer: "Sim. Sem contratos de fidelidade ou multas. Se estiver em um plano pago, você continua com acesso até o final do período contratado."
  },
  {
    question: "A plataforma oferece suporte para check-in no evento?",
    answer: "Sim. Temos aplicativo dedicado para check-in, leitura de QR Code offline, controle de múltiplas entradas e relatório de presença em tempo real."
  },
]

export function OrganizerView() {
  const now = new Date()
  const currentMonthDay = now.getDate()

  // Calculate Monday of current week
  const diff = now.getDay() === 0 ? -6 : 1 - now.getDay()
  const monday = new Date(now)
  monday.setDate(now.getDate() + diff)

  const weekDays = ['Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb', 'Dom'].map((label, idx) => {
    const d = new Date(monday)
    d.setDate(monday.getDate() + idx)
    return {
      label,
      date: d.getDate(),
      isToday: d.getDate() === currentMonthDay && d.getMonth() === now.getMonth()
    }
  })

  return (
    <div className="max-w-5xl mx-auto space-y-20 md:space-y-32">
      {/* Hero */}
      <section className="text-center max-w-2xl mx-auto space-y-4 md:space-y-6 px-2">
        <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
          Planeje, organize e execute eventos de qualquer escala.
          Do briefing ao pós-evento, tudo em uma única plataforma.
        </p>
        <div className="flex flex-col sm:flex-row gap-3 md:gap-4 justify-center pt-2 md:pt-4">
          <button className="px-5 py-2.5 md:px-6 md:py-3 bg-primary text-primary-foreground rounded-full text-sm font-medium hover:bg-primary/90 transition-colors">
            Começar gratuitamente
          </button>
          <button className="px-5 py-2.5 md:px-6 md:py-3 border border-border text-foreground rounded-full text-sm font-medium hover:border-foreground/50 transition-colors">
            Ver demonstração
          </button>
        </div>
        <p className="text-xs text-muted-foreground/70">Comece sem custo. Pague apenas uma pequena taxa sobre suas vendas.</p>
      </section>

      {/* Calendario Visual */}
      <section className="flex justify-center">
        <div className="flex gap-1 md:gap-2">
          {weekDays.map((day) => (
            <div key={day.label} className="text-center">
              <div className={cn(
                "w-10 h-10 md:w-14 md:h-14 rounded-lg md:rounded-xl flex items-center justify-center text-xs md:text-sm font-medium mb-1 md:mb-2",
                day.isToday ? "bg-primary text-primary-foreground" : "bg-muted text-muted-foreground"
              )}>
                {day.date}
              </div>
              <span className="text-[10px] md:text-xs text-muted-foreground">{day.label}</span>
            </div>
          ))}
        </div>
      </section>

      {/* Features */}
      <section className="space-y-8 md:space-y-12">
        <div className="text-center space-y-2">
          <p className="text-xs md:text-sm text-muted-foreground uppercase tracking-wider">Recursos</p>
          <h2 className="text-2xl md:text-4xl font-semibold text-foreground">
            Tudo para criar eventos<br className="hidden md:block" /> impecáveis.
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
          <p className="text-muted-foreground">Três passos para eventos perfeitos</p>
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

      {/* Modelo de Receita */}
      <section className="bg-muted rounded-2xl md:rounded-3xl p-6 md:p-12 space-y-6 md:space-y-8">
        <div className="text-center space-y-2">
          <p className="text-xs md:text-sm text-muted-foreground uppercase tracking-wider">Modelo transparente</p>
          <h2 className="text-2xl md:text-3xl font-semibold text-foreground">
            Sem mensalidades.<br />
            <span className="text-muted-foreground">Você só paga quando vende.</span>
          </h2>
        </div>

        <p className="text-center text-sm md:text-base text-muted-foreground max-w-2xl mx-auto">
          Cobramos apenas uma pequena porcentagem sobre cada produto vendido dentro da plataforma.
          Sem surpresas, sem taxas ocultas. Seu lucro cresce junto com o nosso.
        </p>

        <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6">
          {[
            { title: "Taxa por transação", desc: "Uma porcentagem justa sobre cada venda realizada no seu evento." },
            { title: "Zero custo inicial", desc: "Crie sua conta, configure eventos e comece a vender sem pagar nada." },
            { title: "Escale sem limites", desc: "Quanto mais você vende, mais a plataforma trabalha por você." },
            { title: "Pagamento seguro", desc: "Processamento protegido com repasse direto para sua conta." },
          ].map((item, idx) => (
            <div key={idx} className="space-y-1 md:space-y-2">
              <h4 className="text-sm md:text-base font-medium text-foreground">{item.title}</h4>
              <p className="text-xs md:text-sm text-muted-foreground leading-relaxed">{item.desc}</p>
            </div>
          ))}
        </div>

        {/* Exemplo de transação */}
        <div className="bg-card border border-border rounded-xl md:rounded-2xl p-4 md:p-6 max-w-md mx-auto">
          <p className="text-xs text-muted-foreground mb-3 md:mb-4 uppercase tracking-wider">Exemplo de transação</p>
          <div className="space-y-2 md:space-y-3 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Ingresso VIP — Show de Verão</span>
              <span className="font-medium text-foreground">R$ 150,00</span>
            </div>
            <div className="flex justify-between text-muted-foreground">
              <span>Taxa da plataforma</span>
              <span>- R$ 12,00</span>
            </div>
            <div className="h-px bg-border my-2" />
            <div className="flex justify-between font-medium text-foreground">
              <span>Você recebe</span>
              <span>R$ 138,00</span>
            </div>
            <p className="text-xs text-muted-foreground/70 text-right">92% para você</p>
          </div>
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
      <section className="bg-foreground rounded-2xl md:rounded-3xl p-6 md:p-12 lg:p-16 text-center space-y-4 md:space-y-6">
        <h2 className="text-xl md:text-3xl font-semibold text-background">
          Seu próximo evento começa aqui
        </h2>
        <p className="text-sm md:text-base text-muted/80 max-w-md mx-auto">
          Crie sua conta gratuita e descubra como simplificar toda a gestão dos seus eventos.
        </p>
        <div className="flex flex-col sm:flex-row gap-3 justify-center pt-2">
          <button className="px-5 py-2.5 md:px-6 md:py-3 bg-background text-foreground rounded-full text-sm font-medium hover:bg-background/90 transition-colors">
            Começar gratuitamente
          </button>
          <button className="px-5 py-2.5 md:px-6 md:py-3 border border-muted/50 text-background rounded-full text-sm font-medium hover:border-muted transition-colors">
            Falar com vendas
          </button>
        </div>
        <p className="text-xs text-muted/60">Sem cartão de crédito. Cancele quando quiser.</p>
      </section>
    </div>
  )
}