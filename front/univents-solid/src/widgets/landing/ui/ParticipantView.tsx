import { useNavigate, useRouter } from "@tanstack/solid-router";
import { For } from "solid-js";
import { FAQSection, type FAQItem } from "./FAQSection";
const features = [
  [
    "Encontre eventos facilmente",
    "Pesquise por categoria, data ou localidade e descubra eventos que combinam com você.",
  ],
  [
    "Compra segura de ingressos",
    "Adquira seus ingressos com pagamento protegido e receba a confirmação na hora.",
  ],
  [
    "Ingresso digital no celular",
    "Acesse seu ingresso direto pelo app. Sem papel, sem filas — basta apresentar o QR Code.",
  ],
  [
    "Notificações do evento",
    "Receba lembretes, atualizações e informações importantes sobre seus eventos.",
  ],
  [
    "Avalie e recomende",
    "Compartilhe sua experiência e ajude outros participantes a escolherem os melhores eventos.",
  ],
  [
    "Garantia de reembolso",
    "Se o evento for cancelado, você recebe o reembolso automaticamente na sua conta.",
  ],
];
const steps = [
  [
    "01",
    "Crie sua conta",
    "Cadastre-se em segundos com e-mail ou redes sociais. Totalmente gratuito.",
  ],
  [
    "02",
    "Encontre seu evento",
    "Explore eventos por categoria, local ou data. Compre seu ingresso com segurança.",
  ],
  [
    "03",
    "Aproveite a experiência",
    "Apresente seu QR Code na entrada e curta o evento sem preocupação.",
  ],
];
const faqs: FAQItem[] = [
  {
    question: "Como funciona a taxa sobre vendas?",
    answer:
      "Não há taxa para participantes. O preço do ingresso é o valor final que você paga. Organizadores pagam uma pequena comissão apenas sobre as vendas realizadas.",
  },
  {
    question: "Preciso pagar alguma mensalidade?",
    answer:
      "Não. Para participantes, o uso é 100% gratuito. Para organizadores, oferecemos planos gratuitos e pagos, mas você só paga se escolher recursos avançados.",
  },
  {
    question: "Quais tipos de eventos posso gerenciar?",
    answer:
      "Qualquer tipo: shows, festivais, conferências, workshops, esportivos, corporativos e sociais. Não há limitação de categoria ou tamanho.",
  },
  {
    question: "Como recebo o dinheiro das vendas?",
    answer:
      "O repasse é feito diretamente para sua conta bancária em até 2 dias úteis após a transação.",
  },
  {
    question: "Posso cancelar a qualquer momento?",
    answer:
      "Sim. Não há contratos de fidelidade ou período mínimo. Cancele quando quiser sem taxas de rescisão.",
  },
  {
    question: "A plataforma oferece suporte para check-in no evento?",
    answer:
      "Sim. Oferecemos app de check-in com leitura de QR Code, lista de convidados offline e controle de entrada em tempo real.",
  },
];
export function ParticipantView() {
  const navigate = useNavigate();
  const router = useRouter();
  const authenticated = () =>
    (router.options.context as { auth?: { isAuthenticated: boolean } }).auth
      ?.isAuthenticated === true;
  const go = () =>
    void navigate({ to: authenticated() ? "/events" : "/auth" } as never);
  return (
    <div class="mx-auto max-w-5xl space-y-20 md:space-y-32">
      <section class="mx-auto max-w-2xl space-y-4 px-2 text-center md:space-y-6">
        <p class="text-base text-muted-foreground leading-relaxed">
          Encontre os melhores eventos perto de você, compre ingressos com
          segurança e aproveite cada momento sem complicação.
        </p>
        <div class="flex justify-center pt-2 md:pt-4">
          <button
            class="rounded-full bg-primary px-6 py-3 text-sm font-medium text-primary-foreground"
            onClick={go}
          >
            {authenticated()
              ? "Continuar explorando"
              : "Entrar para ver eventos"}
          </button>
        </div>
      </section>
      <section id="trending" class="scroll-mt-20 space-y-6 md:space-y-8">
        <div class="space-y-2">
          <p class="text-xs uppercase tracking-[0.24em] text-muted-foreground">
            Em alta agora
          </p>
          <h2 class="text-lg font-medium md:text-2xl">Eventos em destaque</h2>
          <p class="max-w-2xl text-sm text-muted-foreground md:text-base">
            Uma seleção curta dos eventos mais recentes e relevantes para você
            começar sem ruído.
          </p>
        </div>
        <div class="rounded-2xl border border-dashed border-border bg-muted/40 p-8 text-center text-sm text-muted-foreground">
          Assim que novos eventos forem publicados, eles aparecem aqui.
        </div>
      </section>
      <section class="space-y-8 md:space-y-12">
        <div class="space-y-2 text-center">
          <p class="text-xs uppercase tracking-wider text-muted-foreground md:text-sm">
            Para participantes
          </p>
          <h2 class="text-2xl font-semibold md:text-4xl">
            Tudo para curtir seus
            <br class="hidden md:block" /> eventos favoritos.
          </h2>
        </div>
        <div class="grid gap-x-6 gap-y-8 sm:grid-cols-2 md:gap-x-8 md:gap-y-12 lg:grid-cols-3">
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
      <section class="space-y-8 md:space-y-12">
        <div class="text-center">
          <h2 class="text-2xl font-semibold md:text-4xl">Como funciona</h2>
          <p class="text-muted-foreground">Simples assim</p>
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
      <section class="space-y-4 rounded-2xl bg-muted p-6 text-center md:space-y-6 md:rounded-3xl md:p-12 lg:p-16">
        <h2 class="text-xl font-semibold md:text-3xl">
          Seu próximo evento te espera
        </h2>
        <p class="mx-auto max-w-md text-sm text-muted-foreground md:text-base">
          Crie sua conta gratuita e descubra eventos incríveis acontecendo perto
          de você.
        </p>
        <button
          class="rounded-full bg-primary px-6 py-3 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
          onClick={go}
        >
          Criar conta grátis
        </button>
        <p class="text-xs text-muted-foreground/70">
          100% gratuito. Sem cartão de crédito.
        </p>
      </section>
    </div>
  );
}
