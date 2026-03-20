import { createFileRoute } from '@tanstack/react-router'
import { cn } from '@/shared/lib/utils'
import { motion } from 'motion/react'
import {
  Check,
  X,
  Clock,
  Calendar,
  Ticket,
  MapPin,
  Activity,
  Minus,
  Sparkles,
  ShoppingBag,
  ClipboardCheck
} from 'lucide-react'

export const Route = createFileRoute('/comparative')({
  component: ComparativoPage,
})

type FeatureStatus = 'yes' | 'no' | 'soon' | 'partial'

interface Feature {
  name: string
  current: FeatureStatus
  competitor: FeatureStatus
  highlight?: boolean
}

interface Section {
  id: string
  title: string
  icon: React.ElementType
  features: Feature[]
}

const sections: Section[] = [
  {
    id: 'events',
    title: 'Funcionalidade Eventos',
    icon: Calendar,
    features: [
      { name: 'Criar Eventos', current: 'yes', competitor: 'yes', highlight: true },
      { name: 'Configurar Evento', current: 'yes', competitor: 'yes' },
      { name: 'Eventos Sub Institucionais', current: 'yes', competitor: 'yes', highlight: true },
      { name: 'Dashboard do evento', current: 'yes', competitor: 'yes' },
      { name: 'Página do Evento', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Página do Evento Customizável na identidade do Evento', current: 'yes', competitor: 'no' },
      { name: 'Templates de Eventos Simples', current: 'yes', competitor: 'yes', highlight: true },
      { name: 'Templates de Eventos Completamente Customizáveis', current: 'yes', competitor: 'no' },
      { name: 'Análise histórica de dados institucionais', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Customização Completa da Identidade Visual', current: 'yes', competitor: 'no' },
      { name: 'Customização Personalizada da Página do Evento', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Customização Completa das comunicações (email, posts, etc.)', current: 'yes', competitor: 'no' },
      { name: 'Automações IFTTT pré prontas', current: 'no', competitor: 'no', highlight: true },
    ]
  },
  {
    id: 'editions',
    title: 'Funcionalidade Edições',
    icon: Clock,
    features: [
      { name: 'Edições Separadas por evento', current: 'yes', competitor: 'no' },
      { name: 'Configurar Edição', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Página da Edição', current: 'yes', competitor: 'no' },
      { name: 'Página da Edição Customizável na identidade do Evento', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Dashboard da Edição', current: 'yes', competitor: 'no' },
      { name: 'Templates de Edições Simples', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Templates de Eventos Completamente Customizáveis', current: 'yes', competitor: 'no' },
      { name: 'Customização Completa da Identidade Visual', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Customização Personalizada da Página de Edição', current: 'yes', competitor: 'no' },
      { name: 'Customização COmpleta das comunicações (email, posts, etc.)', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Automações IFTTT pré prontas', current: 'no', competitor: 'no', highlight: true },
      { name: 'Formulario de inscrição', current: 'soon', competitor: 'yes', highlight: true },
    ]
  },
  {
    id: 'tickets',
    title: 'Funcionalidade Ingressos',
    icon: Ticket,
    features: [
      { name: 'Criar Ingressos', current: 'yes', competitor: 'yes' },
      { name: 'Transferência de Ingresso', current: 'yes', competitor: 'yes', highlight: true },
      { name: 'Configurar Permissões do Ingresso', current: 'yes', competitor: 'no' },
      { name: 'Ingressos com Visual Customizável', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Presentear Ingresso', current: 'yes', competitor: 'no' },
      { name: 'Ingressos que habilitam descontos', current: 'soon', competitor: 'no', highlight: true },
    ]
  },
  {
    id: 'activities',
    title: 'Funcionalidades Atividades',
    icon: Activity,
    features: [
      { name: 'Criar Atividade', current: 'yes', competitor: 'yes', highlight: true },
      { name: 'Editar Atividade', current: 'yes', competitor: 'yes' },
      { name: 'Inscrições em atividades', current: 'yes', competitor: 'yes', highlight: true },
      { name: 'Atividades pagas', current: 'yes', competitor: 'yes' },
      { name: 'Controle de Presença', current: 'yes', competitor: 'yes', highlight: true },
      { name: 'Credenciamento QR Code', current: 'yes', competitor: 'yes' },
      { name: 'Relatórios Automáticos', current: 'yes', competitor: 'yes', highlight: true },
      { name: 'Certificados Automáticos (Com suporte para assinatura)', current: 'yes', competitor: 'yes' },
      { name: 'Atividades com material', current: 'yes', competitor: 'yes', highlight: true },
      { name: 'Vagas por Atividade', current: 'yes', competitor: 'yes' },
      { name: 'Credenciamento NFC', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Controle de tempo mínimo para presença', current: 'yes', competitor: 'no' },
      { name: 'Controle de Acesso seletivo', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Entrada por tokens', current: 'yes', competitor: 'no' },
      { name: 'Avaliação de Atividades (Privado ao evento)', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Atividades Complexas (Hackathon)', current: 'yes', competitor: 'no' },
      { name: 'Lista de espera', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Lista de interesse', current: 'yes', competitor: 'no' },
      { name: 'Customização da atividade', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Submissão de material de participantes', current: 'soon', competitor: 'yes' },
    ]
  },
  {
    id: 'zones',
    title: 'Funcionalidade Zonas de Controle',
    icon: MapPin,
    features: [
      { name: 'Criar Zona de Controle', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Definir Acessos Zona de Controle', current: 'yes', competitor: 'no' },
      { name: 'Dashboard Para dados', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Exportar planilhas e relatórios', current: 'yes', competitor: 'yes' },
    ]
  },
  {
    id: 'loja',
    title: 'Funcionalidade Loja',
    icon: ShoppingBag,
    features: [
      { name: 'Códigos de Desconto', current: 'yes', competitor: 'yes', highlight: true },
      { name: 'Sistema de Reserva no Check Out', current: 'yes', competitor: 'yes' },
      { name: 'Exportar Relatórios e Planilhas', current: 'yes', competitor: 'yes', highlight: true },
      { name: 'Criar Produtos para venda', current: 'yes', competitor: 'no' },
      { name: 'Definir estoque / ilimitado em produto', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Criar combos (Ex: Camisa, Caneca, Ingresso)', current: 'yes', competitor: 'no' },
      { name: 'Histórico de pagamentos, desistências, inválidos, etc.', current: 'yes', competitor: 'partial', highlight: true },
      { name: 'Sistema de validação de entrega de Produtos', current: 'yes', competitor: 'no' },
      { name: 'Lista de desejos', current: 'soon', competitor: 'no', highlight: true },
      { name: 'Presentear outro usuário', current: 'yes', competitor: 'no' },
      { name: 'Promoções personalizadas', current: 'soon', competitor: 'no', highlight: true },
      { name: 'Limite por usuário', current: 'soon', competitor: 'no' },
    ]
  },
  {
    id: 'staff',
    title: 'Funcionalidade Controle de Staff',
    icon: ClipboardCheck,
    features: [
      { name: 'Painéis de Staff', current: 'yes', competitor: 'yes', highlight: true },
      { name: 'Editor de atividades da Staff (Linha temporal)', current: 'yes', competitor: 'no' },
      { name: 'Avisos de Burnout de Staff', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Controle de Acesso Granular da Staff', current: 'yes', competitor: 'no' },
      { name: 'Resumo do Dia', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Notas de Ocorrência', current: 'yes', competitor: 'no' },
      { name: 'Sugestões de balanceamento de carga da Staff', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Lembrentes de tarefa automáticos', current: 'yes', competitor: 'no' },
      { name: 'Cálculo automático de Carga trabalhada (com validação)', current: 'yes', competitor: 'no', highlight: true },
      { name: 'Gerador de Certificados da Staff (Com carga Horária)', current: 'yes', competitor: 'no' },
      { name: 'Relatórios (Tarefas completas, Horas Trabalhadas, Atividades Cobertas) Por membro da staff', current: 'yes', competitor: 'no', highlight: true },
    ]
  },
]

// const additionalFeatures = [
//   { icon: Shield, label: 'LGPD', desc: 'Compliance total' },
//   { icon: Smartphone, label: 'App', desc: 'iOS & Android' },
//   { icon: Globe, label: 'Idiomas', desc: 'PT, EN, ES' },
//   { icon: BarChart3, label: 'Analytics', desc: 'Real-time' },
//   { icon: Zap, label: 'Automação', desc: 'Workflows' },
//   { icon: Bell, label: 'Comunicação', desc: 'Multi-canal' },
//   { icon: Palette, label: 'Branding', desc: 'White label' },
//   { icon: CreditCard, label: 'Financeiro', desc: 'Split automático' },
//   { icon: Users, label: 'Equipe', desc: 'RBAC' },
//   { icon: QrCode, label: 'Acesso', desc: 'QR dinâmico' },
// ]

function StatusIndicator({ status }: { status: FeatureStatus }) {
  const configs = {
    yes: {
      icon: Check,
      className: 'text-emerald-500',
      bg: 'bg-emerald-500/10',
    },
    no: {
      icon: X,
      className: 'text-muted-foreground/40',
      bg: 'bg-muted',
    },
    soon: {
      icon: Clock,
      className: 'text-amber-500',
      bg: 'bg-amber-500/10',
    },
    partial: {
      icon: Minus,
      className: 'text-orange-500',
      bg: 'bg-orange-500/10',
    }
  }

  const config = configs[status]
  const Icon = config.icon

  return (
    <div className={cn(
      'w-8 h-8 rounded flex items-center justify-center',
      config.bg
    )}>
      <Icon className={cn('w-4 h-4', config.className)} />
    </div>
  )
}

function PlatformCard({
  name,
  tagline,
  isUs = false,
  stats
}: {
  name: string
  tagline: string
  isUs?: boolean
  stats: { label: string; value: string }[]
}) {
  return (
    <div className={cn(
      'relative p-6 border',
      isUs ? 'bg-card border-l-4 border-l-primary' : 'bg-muted/20 border-border'
    )}>
      {isUs && (
        <div className="absolute -top-px left-0 right-0 h-px bg-linear-to-r from-primary via-primary/50 to-transparent" />
      )}

      <div className="space-y-4">
        <div>
          <div className="flex items-center gap-2 mb-1">
            {isUs && <Sparkles className="w-4 h-4 text-primary" />}
            <h3 className="text-lg font-semibold text-foreground tracking-tight">
              {name}
            </h3>
          </div>
          <p className="text-sm text-muted-foreground">
            {tagline}
          </p>
        </div>

        <div className="grid grid-cols-2 gap-3 pt-4 border-t border-border">
          {stats.map((stat, idx) => (
            <div key={idx}>
              <div className={cn(
                'text-2xl font-semibold tracking-tight',
                isUs ? 'text-foreground' : 'text-muted-foreground'
              )}>
                {stat.value}
              </div>
              <div className="text-xs text-muted-foreground uppercase tracking-wide">
                {stat.label}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

function ComparisonTable({ section }: { section: Section }) {
  const Icon = section.icon

  const yesCount = section.features.filter(f => f.current === 'yes').length
  const total = section.features.length

  return (
    <div className="border border-border bg-card">
      {/* Header */}
      <div className="flex items-center justify-between px-5 py-4 border-b border-border bg-muted/30">
        <div className="flex items-center gap-3">
          <Icon className="w-5 h-5 text-primary" />
          <h3 className="font-medium text-foreground">{section.title}</h3>
        </div>
        <div className="text-xs text-muted-foreground">
          {yesCount}/{total} exclusivas
        </div>
      </div>

      {/* Table */}
      <div className="divide-y divide-border/50">
        {section.features.map((feature, idx) => (
          <div
            key={idx}
            className={cn(
              'grid grid-cols-3 gap-4 px-5 py-3 items-center',
              feature.highlight && 'bg-primary/5'
            )}
          >
            <span className={cn(
              'text-sm',
              feature.highlight ? 'font-medium text-foreground' : 'text-muted-foreground'
            )}>
              {feature.name}
            </span>

            <div className="flex justify-center">
              <StatusIndicator status={feature.current} />
            </div>

            <div className="flex justify-center">
              <StatusIndicator status={feature.competitor} />
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

function SummaryBar() {
  const calculateScore = (platform: 'current' | 'competitor') => {
    let score = 0
    sections.forEach(s => {
      s.features.forEach(f => {
        const status = platform === 'current' ? f.current : f.competitor
        if (status === 'yes') score += 2
        if (status === 'partial') score += 1
        if (status === 'soon') score += 0.5
      })
    })
    return score
  }

  const ourScore = calculateScore('current')
  const theirScore = calculateScore('competitor')
  const total = ourScore + theirScore
  const ourPercent = Math.round((ourScore / total) * 100)

  return (
    <div className="border border-border bg-card p-6">
      <div className="flex items-center justify-between mb-4">
        <span className="text-sm font-medium text-foreground">Comparativo geral</span>
        <span className="text-sm text-muted-foreground">{ourPercent}% vs {100 - ourPercent}%</span>
      </div>

      <div className="h-2 bg-muted rounded-full overflow-hidden flex">
        <div
          className="h-full bg-primary transition-all duration-1000"
          style={{ width: `${ourPercent}%` }}
        />
        <div
          className="h-full bg-muted-foreground/20 transition-all duration-1000"
          style={{ width: `${100 - ourPercent}%` }}
        />
      </div>

      <div className="flex justify-between mt-3 text-xs text-muted-foreground">
        <span>Nossa plataforma</span>
        <span>Concorrente</span>
      </div>
    </div>
  )
}

function ComparativoPage() {
  const calculateStats = () => {
    let ourUnique = 0
    let ourTotal = 0
    let competitorTotal = 0
    let competitorUnique = 0
    let ourPartial = 0

    sections.forEach(section => {
      section.features.forEach(feature => {
        ourTotal++
        competitorTotal++

        if (feature.current === 'yes' && (feature.competitor === 'no' || feature.competitor === 'soon')) {
          ourUnique++
        }

        if (feature.competitor === 'yes' && feature.current !== 'yes') {
          competitorUnique++
        }

        // Nós temos parcial
        if (feature.current === 'partial') {
          ourPartial++
        }
      })
    })

    const ourEffective = ourTotal - ourPartial + (ourPartial * 0.5)
    const competitorEffective = competitorTotal - sections.flatMap(s => s.features).filter(f => f.competitor === 'partial').length * 0.5

    return {
      ourTotal,
      ourUnique,
      ourPartial,
      competitorTotal,
      competitorUnique,
      coverage: Math.round((ourEffective / (ourEffective + competitorEffective)) * 100)
    }
  }

  const stats = calculateStats()
  return (
    <div className="min-h-screen bg-background">
      {/* Hero */}
      <header className="border-b border-border">
        <div className="max-w-4xl mx-auto px-6 py-16 md:py-24">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="space-y-6"
          >
            <h1 className="text-3xl md:text-5xl font-semibold tracking-tight text-foreground">
              Comparativo técnico
            </h1>
            <p className="text-base md:text-lg text-muted-foreground max-w-xl leading-relaxed">
              Análise detalhada de funcionalidades entre nossa plataforma e a solução tradicional do mercado.
            </p>
          </motion.div>

          {/* Platform Cards */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="grid md:grid-cols-2 gap-px bg-border border border-border mt-12"
          >
            <PlatformCard
              name="Univents"
              tagline="Plataforma Atual"
              isUs={true}
              stats={[
                { label: 'Total', value: String(stats.ourTotal) },
                { label: 'Recursos Exclusivos', value: `${stats.ourUnique}+` },
              ]}
            />
            <PlatformCard
              name="Even3"
              tagline="Plataforma Concorrente"
              stats={[
                { label: 'Recursos Exclusivos', value: `${stats.competitorUnique}+` },
              ]}
            />
          </motion.div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-4xl mx-auto px-6 py-12 md:py-16 space-y-8">
        {/* Summary */}
        <SummaryBar />

        {/* Legend */}
        <div className="flex flex-wrap gap-4 text-xs text-muted-foreground border-b border-border pb-6">
          <div className="flex items-center gap-2">
            <div className="w-4 h-4 rounded bg-emerald-500/10 flex items-center justify-center">
              <Check className="w-3 h-3 text-emerald-500" />
            </div>
            <span>Disponível</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-4 h-4 rounded bg-orange-500/10 flex items-center justify-center">
              <Minus className="w-3 h-3 text-orange-500" />
            </div>
            <span>Parcial</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-4 h-4 rounded bg-amber-500/10 flex items-center justify-center">
              <Clock className="w-3 h-3 text-amber-500" />
            </div>
            <span>Roadmap</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-4 h-4 rounded bg-muted flex items-center justify-center">
              <X className="w-3 h-3 text-muted-foreground/40" />
            </div>
            <span>Indisponível</span>
          </div>
        </div>

        {/* Sections */}
        <div className="space-y-6">
          {sections.map((section, idx) => (
            <motion.div
              key={section.id}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: idx * 0.05 }}
            >
              <ComparisonTable section={section} />
            </motion.div>
          ))}
        </div>

        {/* Additional Features */}
        {/* <div className="pt-8 border-t border-border">
          <h3 className="text-sm font-medium text-foreground mb-4 uppercase tracking-wide">
            Diferenciais adicionais
          </h3>
          <div className="grid grid-cols-2 md:grid-cols-5 gap-3">
            {additionalFeatures.map((feat, idx) => {
              const Icon = feat.icon
              return (
                <div
                  key={idx}
                  className="flex items-center gap-2 p-3 border border-border bg-card"
                >
                  <Icon className="w-4 h-4 text-primary shrink-0" />
                  <div className="min-w-0">
                    <div className="text-xs font-medium text-foreground truncate">
                      {feat.label}
                    </div>
                    <div className="text-[10px] text-muted-foreground truncate">
                      {feat.desc}
                    </div>
                  </div>
                </div>
              )
            })}
          </div>
        </div> */}
      </main>
    </div>
  )
}