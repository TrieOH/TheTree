import { createEffect, createSignal, onSettled, For } from "solid-js";
import {
  Button,
  Card,
  Badge,
  Link,
  Tooltip,
  buttonVariants,
  cn,
} from "@trieoh/ui-solid";
import IconSun from "~icons/lucide/sun";
import IconMoon from "~icons/lucide/moon";
import { applyTheme, readThemePreference, saveThemePreference } from "./theme";

const resources = [
  {
    name: "Components",
    label: "A interface, peça por peça.",
    description:
      "Explore a biblioteca de componentes e experimente as interações que dão forma aos nossos produtos.",
    href: "https://components.trieoh.com",
    detail: "Biblioteca de componentes",
  },
  {
    name: "Style Guide",
    label: "Uma linguagem em comum.",
    description:
      "Conheça as cores, fontes, espaçamentos e formas que compõem a identidade visual da TrieOH.",
    href: "http://styles.trieoh.com/",
    detail: "Fundamentos visuais",
  },
];
const products = [
  {
    id: "identityx",
    name: "IdentityX",
    href: "https://dev.identityx.com.br",
    category: "Identidade e acesso",
    number: "01",
    title: "O ponto de partida de cada conexão.",
    description:
      "Autenticação e gestão de identidades para organizar usuários, projetos e organizações em uma base comum.",
    features: [
      "Autenticação",
      "Organizações e projetos",
      "Perfis e permissões",
    ],
    color: "bg-secondary/25 text-foreground",
  },
  {
    id: "univents",
    name: "Univents",
    href: "https://univents.com.br",
    category: "Eventos e experiências",
    number: "02",
    title: "Encontros que vão além da inscrição.",
    description:
      "Uma plataforma para organizar eventos e suas edições, da programação e dos ingressos ao credenciamento e aos certificados.",
    features: [
      "Eventos e programação",
      "Inscrições e ingressos",
      "Crachás e certificados",
    ],
    color: "bg-accent/20 text-foreground",
  },
  {
    id: "informd",
    name: "Informd",
    href: "https://dev.informd.trieoh.com",
    category: "Formulários e respostas",
    number: "03",
    title: "Perguntas com espaço para respostas.",
    description:
      "Formulários e pesquisas organizados em etapas e campos, com respostas reunidas por contexto.",
    features: ["Formulários", "Etapas e campos", "Coleta de respostas"],
    color: "bg-muted text-foreground",
  },
  {
    id: "payssage",
    name: "Payssage",
    href: "https://dev.payssage.trieoh.com",
    category: "Pagamentos e transações",
    number: "04",
    title: "Uma base para movimentar ideias.",
    description:
      "Gestão de carteiras, vendedores e intenções de pagamento, com notificações sobre o andamento das transações.",
    features: ["Carteiras e vendedores", "Intenções de pagamento", "Webhooks"],
    color: "bg-secondary/25 text-foreground",
  },
];
const eyebrow =
  "text-[11px] font-medium uppercase tracking-[0.14em] text-muted-foreground";
const wordmark =
  "font-heading text-3xl font-bold tracking-[-1.8px] text-foreground no-underline hover:no-underline";

export default function App() {
  const [dark, setDark] = createSignal(
    document.documentElement.classList.contains("dark"),
  );
  let preference = readThemePreference();

  createEffect(
    () => dark(),
    (value) => applyTheme(value ? "dark" : "light"),
  );

  onSettled(() => {
    const media = window.matchMedia("(prefers-color-scheme: dark)");
    const syncSystemTheme = () => {
      if (preference === null) setDark(media.matches);
    };
    media.addEventListener("change", syncSystemTheme);
    syncSystemTheme();
    return () => media.removeEventListener("change", syncSystemTheme);
  });

  const toggleTheme = () => {
    const next = !dark();
    preference = next ? "dark" : "light";
    saveThemePreference(preference);
    setDark(next);
  };

  return (
    <>
      <a
        href="#main"
        class="fixed top-3 left-3 z-50 translate-y-[-200%] bg-foreground p-3 text-background focus:translate-y-0"
      >
        Pular para o conteúdo
      </a>
      <div class="mx-auto max-w-360 px-6 sm:px-10 lg:px-16">
        <header class="flex min-h-24 flex-wrap items-center justify-between gap-4 border-b border-border py-5">
          <Link href="/" class={wordmark} aria-label="TrieOH — início">
            TrieOH
          </Link>
          <div class="flex items-center gap-4 sm:gap-8">
            <nav
              aria-label="Navegação principal"
              class="flex gap-4 text-xs sm:gap-6 sm:text-sm"
            >
              <Link href="#sistemas" class="text-foreground">
                Sistemas
              </Link>
              <Link href="#recursos" class="text-foreground">
                Para construir
              </Link>
            </nav>
            <Tooltip
              label={dark() ? "Usar tema claro" : "Usar tema escuro"}
              placement="bottom"
              align="end"
            >
              <Button
                variant="ghost"
                size="icon"
                aria-label={dark() ? "Ativar tema claro" : "Ativar tema escuro"}
                onClick={toggleTheme}
              >
                {dark() ? (
                  <IconSun class="size-4" />
                ) : (
                  <IconMoon class="size-4" />
                )}
              </Button>
            </Tooltip>
          </div>
        </header>
        <main id="main">
          <section class="py-16 sm:py-24" aria-labelledby="intro-title">
            <div class="flex items-center gap-2.5">
              <span class="size-2 rounded-full bg-accent" aria-hidden="true" />
              <p class={eyebrow}>Pessoas. Ideias. Sistemas.</p>
            </div>
            <div class="mt-8 grid items-end gap-8 lg:grid-cols-[1.5fr_1fr] lg:gap-16">
              <div>
                <h1
                  id="intro-title"
                  class="font-heading text-[clamp(52px,7.6vw,108px)] leading-[1.04] font-medium tracking-[-0.065em]"
                >
                  Ideias que
                  <br />
                  se{" "}
                  <span class="text-muted-foreground">
                    conectam<span class="text-accent">.</span>
                  </span>
                </h1>
                <p class="mt-7 max-w-xl text-base leading-relaxed text-muted-foreground">
                  Somos a TrieOH. Criamos sistemas para conectar pessoas,
                  organizar experiências e dar forma a novas possibilidades.
                </p>
              </div>
              <div class="lg:pb-2">
                <p class="max-w-sm text-sm leading-relaxed text-muted-foreground">
                  Identidade, eventos, formulários e pagamentos. Produtos com
                  propósitos próprios, construídos sobre uma linguagem em comum.
                </p>
                <Link
                  href="#sistemas"
                  class={cn(
                    buttonVariants({ size: "lg" }),
                    "mt-6 h-auto gap-8 px-5 py-4 no-underline hover:no-underline",
                  )}
                >
                  Conheça nossos sistemas <span aria-hidden="true">↘</span>
                </Link>
              </div>
            </div>
            <div class="mt-12 grid grid-cols-2 overflow-hidden rounded-lg border border-border lg:mt-16 lg:grid-cols-4">
              <For each={products}>
                {(product) => (
                  <Link
                    href={`#${product.id}`}
                    class="group border-border p-5 text-foreground no-underline hover:bg-muted/60 hover:no-underline not-last:border-r sm:p-6"
                  >
                    <span class="flex justify-between text-[10px] uppercase tracking-wider text-muted-foreground">
                      {product.category}
                      <span aria-hidden="true">↘</span>
                    </span>
                    <span class="mt-5 block font-heading text-xl font-medium tracking-tight sm:text-2xl">
                      {product.name}
                    </span>
                  </Link>
                )}
              </For>
            </div>
          </section>
          <section
            id="sistemas"
            class="scroll-mt-8 border-t border-border py-14 sm:py-20"
            aria-labelledby="systems-title"
          >
            <div class="mb-10 flex flex-col justify-between gap-5 sm:flex-row sm:items-end">
              <div>
                <p class={eyebrow}>01 / Ecossistema</p>
                <h2
                  id="systems-title"
                  class="mt-4 font-heading text-4xl font-medium tracking-tight sm:text-5xl"
                >
                  Cada sistema.
                  <br />
                  Uma possibilidade.
                </h2>
              </div>
              <p class="max-w-sm text-sm leading-relaxed text-muted-foreground">
                Conheça os produtos que fazem parte da TrieOH e o papel de cada
                um.
              </p>
            </div>
            <div class="grid gap-5 md:grid-cols-2">
              <For each={products}>
                {(product) => (
                  <Card
                    shape="round"
                    class="scroll-mt-8 overflow-hidden shadow-none"
                  >
                    <div id={product.id} class="scroll-mt-8">
                      <div
                        class={cn(
                          "flex min-h-32 items-center justify-between px-6 py-8 sm:px-8",
                          product.color,
                        )}
                      >
                        <h3 class="font-heading text-3xl font-semibold tracking-tight">
                          {product.name}
                        </h3>
                        <span class="font-mono text-xs opacity-60">
                          /{product.number}
                        </span>
                      </div>
                      <div class="p-6 sm:p-8">
                        <p class={eyebrow}>{product.category}</p>
                        <h4 class="mt-4 font-heading text-2xl font-medium leading-tight tracking-tight">
                          {product.title}
                        </h4>
                        <p class="mt-4 text-sm leading-relaxed text-muted-foreground">
                          {product.description}
                        </p>
                        <ul
                          class="mt-6 flex flex-wrap gap-2"
                          aria-label={`Recursos de ${product.name}`}
                        >
                          <For each={product.features}>
                            {(feature) => (
                              <li>
                                <Badge variant="outline">{feature}</Badge>
                              </li>
                            )}
                          </For>
                        </ul>
                        <Link
                          href={product.href}
                          class={cn(
                            buttonVariants({ variant: "outline", size: "lg" }),
                            "mt-6 gap-6 no-underline hover:no-underline",
                          )}
                        >
                          Acessar {product.name}{" "}
                          <span aria-hidden="true">↗</span>
                        </Link>
                      </div>
                    </div>
                  </Card>
                )}
              </For>
            </div>
          </section>
          <section
            id="recursos"
            class="scroll-mt-8 border-t border-border py-14 sm:py-20"
            aria-labelledby="resources-title"
          >
            <p class={eyebrow}>02 / Nossa linguagem</p>
            <div class="mt-4 grid gap-6 sm:grid-cols-2">
              <h2
                id="resources-title"
                class="font-heading text-4xl font-medium tracking-tight sm:text-5xl"
              >
                Feito de detalhes.
                <br />
                <span class="text-muted-foreground">Pensado como um todo.</span>
              </h2>
              <p class="max-w-md text-sm leading-relaxed text-muted-foreground sm:justify-self-end">
                Compartilhamos os fundamentos e os componentes que dão unidade
                às nossas interfaces. Um ponto de partida para explorar como a
                TrieOH é construída.
              </p>
            </div>
            <div class="mt-10 grid gap-5 md:grid-cols-2">
              <For each={resources}>
                {(resource) => (
                  <Card
                    shape="round"
                    class="flex flex-col items-start p-6 shadow-none sm:p-8"
                  >
                    <p class={eyebrow}>{resource.detail}</p>
                    <h3 class="mt-5 font-heading text-3xl font-medium tracking-tight">
                      {resource.name}
                    </h3>
                    <p class="mt-2 text-sm font-medium">{resource.label}</p>
                    <p class="mt-4 mb-6 max-w-md text-sm leading-relaxed text-muted-foreground">
                      {resource.description}
                    </p>
                    <Link
                      href={resource.href}
                      class={cn(
                        buttonVariants({ variant: "outline", size: "lg" }),
                        "mt-auto gap-8 no-underline hover:no-underline",
                      )}
                    >
                      Explorar {resource.name} <span aria-hidden="true">↗</span>
                    </Link>
                  </Card>
                )}
              </For>
            </div>
          </section>
          <div
            class="mb-14 flex items-center justify-between gap-6 bg-foreground px-6 py-8 font-heading text-background sm:mb-20 sm:px-10 sm:py-12"
            aria-hidden="true"
          >
            <span class="text-3xl font-semibold tracking-tight sm:text-5xl">
              TrieOH
            </span>
            <span class="text-5xl text-accent sm:text-7xl">✳</span>
            <span class="max-w-40 text-lg tracking-tight sm:max-w-none sm:text-3xl">
              Da ideia à experiência.
            </span>
          </div>
        </main>
        <footer class="flex flex-wrap items-center justify-between gap-6 border-t border-border py-8">
          <Link href="/" class={cn(wordmark, "text-2xl")}>
            TrieOH
          </Link>
          <p class="text-xs text-muted-foreground">
            Pessoas, ideias e sistemas em conexão.
          </p>
          <nav aria-label="Recursos" class="flex gap-5 text-xs">
            <For each={resources}>
              {(resource) => (
                <Link href={resource.href} class="text-muted-foreground">
                  {resource.name} ↗
                </Link>
              )}
            </For>
          </nav>
        </footer>
      </div>
    </>
  );
}
