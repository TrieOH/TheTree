import { For } from "solid-js";
import { FAQSection, type FAQItem } from "./FAQSection";
import { Reveal } from "@/shared/ui/Reveal";
const features = [
  [
    "Planejamento completo",
    "Organize cronogramas, prazos e checklists com visão completa de cada etapa do evento.",
  ],
  [
    "Gestão de convidados",
    "Controle listas de presença, convites e confirmações em tempo real.",
  ],
  [
    "Ingressos e inscrições",
    "Venda ingressos, gerencie inscrições e acompanhe a receita do seu evento.",
  ],
  [
    "Fornecedores e logística",
    "Centralize contratos, prazos e comunicação com todos os fornecedores.",
  ],
  [
    "Relatórios pós-evento",
    "Analise resultados, feedbacks e métricas para melhorar seus próximos eventos.",
  ],
  [
    "Automações inteligentes",
    "Automatize lembretes, e-mails e tarefas repetitivas do fluxo de trabalho.",
  ],
];
const steps = [
  [
    "01",
    "Configure seu evento",
    "Defina detalhes, datas, equipe e fornecedores em poucos minutos com nossos templates prontos.",
  ],
  [
    "02",
    "Organize e gerencie",
    "Acompanhe tarefas, convidados e logística em tempo real com visão unificada de tudo.",
  ],
  [
    "03",
    "Execute com confiança",
    "No dia do evento, tenha controle total com check-ins, timeline ao vivo e comunicação integrada.",
  ],
];
const faqs: FAQItem[] = [
  {
    question: "Como funciona a taxa sobre vendas?",
    answer:
      "Cobramos apenas 8% sobre cada ingresso vendido. Não há taxas de setup, mensalidade mínima ou custos ocultos.",
  },
  {
    question: "Preciso pagar mensalidade?",
    answer:
      "Não obrigatoriamente. Nosso plano Starter é gratuito. Planos pagos oferecem recursos avançados, mas você só paga se quiser utilizá-los.",
  },
  {
    question: "Quais tipos de eventos posso gerenciar?",
    answer:
      "Qualquer categoria: corporativos, sociais, culturais, esportivos e educacionais, desde pequenas reuniões até grandes festivais.",
  },
  {
    question: "Como recebo o dinheiro das vendas?",
    answer:
      "O repasse é automático para sua conta em até 2 dias úteis. Também oferecemos antecipação de recebíveis.",
  },
  {
    question: "Posso cancelar a qualquer momento?",
    answer:
      "Sim. Sem contratos de fidelidade ou multas. Você continua com acesso até o fim do período contratado.",
  },
  {
    question: "A plataforma oferece suporte para check-in no evento?",
    answer:
      "Sim. Temos aplicativo dedicado para check-in, leitura de QR Code offline, múltiplas entradas e relatório de presença em tempo real.",
  },
];
export function OrganizerView() {
  return (
    <Reveal>
      <div class="mx-auto max-w-5xl space-y-20 md:space-y-32">
        <section class="mx-auto max-w-2xl space-y-4 text-center">
          <p class="text-base text-muted-foreground leading-relaxed">
            Planeje, organize e execute eventos de qualquer escala. Do briefing
            ao pós-evento, tudo em uma única plataforma.
          </p>
          <p class="text-xs text-muted-foreground/70">
            Comece sem custo. Pague apenas uma pequena taxa sobre suas vendas.
          </p>
        </section>
        <section class="flex justify-center">
          <div class="flex gap-1 md:gap-2">
            <For each={["Seg", "Ter", "Qua", "Qui", "Sex", "Sáb", "Dom"]}>
              {(day, index) => (
                <div class="text-center">
                  <div
                    class={`mb-1 flex h-10 w-10 items-center justify-center rounded-lg text-xs font-medium md:mb-2 md:h-14 md:w-14 md:rounded-xl md:text-sm ${index() === 2 ? "bg-primary text-primary-foreground" : "bg-muted text-muted-foreground"}`}
                  >
                    {new Date().getDate() + index() - 2}
                  </div>
                  <span class="text-[10px] text-muted-foreground md:text-xs">
                    {day}
                  </span>
                </div>
              )}
            </For>
          </div>
        </section>
        <section class="space-y-8 md:space-y-12">
          <div class="space-y-2 text-center">
            <p class="text-xs uppercase tracking-wider text-muted-foreground md:text-sm">
              Recursos
            </p>
            <h2 class="text-2xl font-semibold md:text-4xl">
              Tudo para criar eventos
              <br class="hidden md:block" /> impecáveis.
            </h2>
          </div>
          <div class="grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
            <For each={features}>
              {(feature) => (
                <div class="reveal-item space-y-2">
                  <h3 class="text-base font-medium md:text-lg">{feature[0]}</h3>
                  <p class="text-sm leading-relaxed text-muted-foreground">
                    {feature[1]}
                  </p>
                </div>
              )}
            </For>
          </div>
        </section>
        <section class="space-y-6 rounded-2xl bg-muted p-6 md:space-y-8 md:rounded-3xl md:p-12">
          <div class="space-y-2 text-center">
            <p class="text-xs uppercase tracking-wider text-muted-foreground md:text-sm">
              Modelo transparente
            </p>
            <h2 class="text-2xl font-semibold md:text-3xl">
              Sem mensalidades.
              <br />
              <span class="text-muted-foreground">
                Você só paga quando vende.
              </span>
            </h2>
          </div>
          <p class="mx-auto max-w-2xl text-center text-sm text-muted-foreground md:text-base">
            Cobramos apenas uma pequena porcentagem sobre cada produto vendido
            dentro da plataforma. Sem surpresas, sem taxas ocultas. Seu lucro
            cresce junto com o nosso.
          </p>
          <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <For
              each={[
                [
                  "Taxa por transação",
                  "Uma porcentagem justa sobre cada venda realizada no seu evento.",
                ],
                [
                  "Zero custo inicial",
                  "Crie sua conta, configure eventos e comece a vender sem pagar nada.",
                ],
                [
                  "Escale sem limites",
                  "Quanto mais você vende, mais a plataforma trabalha por você.",
                ],
                [
                  "Pagamento seguro",
                  "Processamento protegido com repasse direto para sua conta.",
                ],
              ]}
            >
              {(item) => (
                <div class="space-y-2">
                  <h4 class="text-sm font-medium md:text-base">{item[0]}</h4>
                  <p class="text-xs leading-relaxed text-muted-foreground md:text-sm">
                    {item[1]}
                  </p>
                </div>
              )}
            </For>
          </div>
          <div class="mx-auto max-w-md space-y-3 rounded-2xl border border-border bg-card p-4 text-sm md:p-6">
            <p class="text-xs uppercase tracking-wider text-muted-foreground">
              Exemplo de transação
            </p>
            <div class="flex justify-between">
              <span class="text-muted-foreground">
                Ingresso VIP — Show de Verão
              </span>
              <span class="font-medium">R$ 150,00</span>
            </div>
            <div class="flex justify-between text-muted-foreground">
              <span>Taxa da plataforma</span>
              <span>- R$ 12,00</span>
            </div>
            <div class="h-px bg-border" />
            <div class="flex justify-between font-medium">
              <span>Você recebe</span>
              <span>R$ 138,00</span>
            </div>
            <p class="text-right text-xs text-muted-foreground/70">
              92% para você
            </p>
          </div>
        </section>
        <section class="space-y-8 md:space-y-12">
          <div class="text-center">
            <h2 class="text-2xl font-semibold md:text-4xl">Como funciona</h2>
            <p class="text-muted-foreground">
              Três passos para eventos perfeitos
            </p>
          </div>
          <div class="grid gap-6 md:grid-cols-3 md:gap-8">
            <For each={steps}>
              {(step) => (
                <div class="reveal-item space-y-3 md:space-y-4">
                  <span class="text-3xl font-semibold text-muted md:text-4xl">
                    {step[0]}
                  </span>
                  <h3 class="text-lg font-medium md:text-xl">{step[1]}</h3>
                  <p class="text-sm leading-relaxed text-muted-foreground md:text-base">
                    {step[2]}
                  </p>
                </div>
              )}
            </For>
          </div>
        </section>
        <section class="mx-auto max-w-2xl space-y-6">
          <div class="space-y-2 text-center">
            <h2 class="text-center text-2xl font-semibold md:text-3xl">
              Perguntas frequentes
            </h2>
            <p class="text-sm text-muted-foreground">Tire suas dúvidas</p>
          </div>
          <FAQSection items={faqs} />
        </section>
      </div>
    </Reveal>
  );
}
