import { createEffect, createSignal, onSettled, For } from "solid-js";
import {
  Badge,
  Button,
  Card,
  Link,
  Tooltip,
  Toaster,
  toast,
  buttonVariants,
  cn,
  type Shape,
} from "@trieoh/ui-solid";
import IconSun from "~icons/lucide/sun";
import IconMoon from "~icons/lucide/moon";
import { applyTheme, readThemePreference, saveThemePreference } from "./theme";

const colors = [
  { name: "Background", token: "background", class: "bg-background" },
  { name: "Foreground", token: "foreground", class: "bg-foreground" },
  { name: "Card", token: "card", class: "bg-card" },
  { name: "Primary", token: "primary", class: "bg-primary" },
  { name: "Secondary", token: "secondary", class: "bg-secondary" },
  { name: "Accent", token: "accent", class: "bg-accent" },
  { name: "Muted", token: "muted", class: "bg-muted" },
  { name: "Border", token: "border", class: "bg-border" },
  { name: "Destructive", token: "destructive", class: "bg-destructive" },
  { name: "Ring", token: "ring", class: "bg-ring" },
];
const spacing = [4, 8, 12, 16, 24, 32, 48, 64];
const shapes: Shape[] = ["sharp", "subtle", "soft", "round", "pill"];
const sectionTitle =
  "font-heading text-3xl font-medium tracking-tight sm:text-4xl";
const eyebrow = "text-xs uppercase tracking-[0.14em] text-muted-foreground";

async function copyToken(token: string) {
  try {
    await navigator.clipboard.writeText(`var(--${token})`);
    toast.success("Token copiado", { description: `var(--${token})` });
  } catch {
    toast.error("Não foi possível copiar o token");
  }
}

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
      <Toaster />
      <a
        href="#main"
        class="fixed top-3 left-3 z-50 -translate-y-[200%] bg-foreground p-3 text-background focus:translate-y-0"
      >
        Pular para o conteúdo
      </a>
      <div class="mx-auto max-w-7xl px-6 sm:px-10">
        <header class="flex min-h-22 flex-wrap items-center justify-between gap-4 border-b border-border py-5">
          <Link
            href="https://trieoh.com"
            class="font-heading text-2xl font-bold tracking-tight text-foreground no-underline hover:no-underline"
          >
            TrieOH{" "}
            <span class="ml-1 text-sm font-normal text-muted-foreground">
              Style Guide
            </span>
          </Link>
          <div class="flex items-center gap-4">
            <Link
              href="https://components.trieoh.com"
              class="text-sm text-foreground"
            >
              Components ↗
            </Link>
            <Tooltip label={dark() ? "Usar tema claro" : "Usar tema escuro"}>
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
            <Badge variant="outline">Identidade / Fundamentos</Badge>
            <h1
              id="intro-title"
              class="mt-6 max-w-3xl font-heading text-5xl leading-[1.05] font-medium tracking-[-0.055em] sm:text-7xl"
            >
              Uma linguagem.
              <br />
              <span class="text-muted-foreground">Muitas possibilidades.</span>
            </h1>
            <p class="mt-6 max-w-xl text-base leading-relaxed text-muted-foreground">
              O ponto de partida visual da TrieOH. Cores, tipografia e
              proporções compartilhadas pelos nossos produtos, do primeiro
              detalhe à interface completa.
            </p>
            <nav
              aria-label="Seções do Style Guide"
              class="mt-8 flex flex-wrap gap-x-6 gap-y-3 text-sm"
            >
              <Link href="#cores">Cores</Link>
              <Link href="#tipografia">Tipografia</Link>
              <Link href="#espacamento">Espaçamento</Link>
              <Link href="#formas">Formas</Link>
              <Link href="#composicao">Em uso</Link>
            </nav>
          </section>
          <section
            id="cores"
            class="scroll-mt-8 border-t border-border py-12"
            aria-labelledby="colors-title"
          >
            <p class={eyebrow}>01 / Cor</p>
            <h2 id="colors-title" class={cn(sectionTitle, "mt-3")}>
              Uma paleta que se adapta.
            </h2>
            <p class="mt-4 max-w-xl text-sm leading-relaxed text-muted-foreground">
              Tokens semânticos mantêm a intenção das cores nos dois temas.
              Alterne o tema para ver as cores reais e clique em um token para
              copiá-lo.
            </p>
            <div class="mt-8 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
              <For each={colors}>
                {(color) => (
                  <Card shape="soft" class="overflow-hidden shadow-none">
                    <div
                      class={cn("h-24 border-b border-border", color.class)}
                      aria-hidden="true"
                    />
                    <div class="p-4">
                      <h3 class="text-sm font-medium">{color.name}</h3>
                      <Button
                        variant="ghost"
                        size="sm"
                        class="mt-2 -ml-2 h-auto whitespace-normal px-2 py-1 font-mono text-[11px] text-muted-foreground"
                        aria-label={`Copiar token ${color.token}`}
                        onClick={() => void copyToken(color.token)}
                      >
                        --{color.token}
                      </Button>
                    </div>
                  </Card>
                )}
              </For>
            </div>
          </section>
          <section
            id="tipografia"
            class="scroll-mt-8 border-t border-border py-12"
            aria-labelledby="type-title"
          >
            <p class={eyebrow}>02 / Tipografia</p>
            <h2 id="type-title" class={cn(sectionTitle, "mt-3")}>
              Personalidade com clareza.
            </h2>
            <div class="mt-8 grid gap-6 md:grid-cols-2">
              <Card class="p-6 sm:p-8 shadow-none">
                <Badge variant="outline">Títulos / Space Grotesk</Badge>
                <p
                  class="mt-6 font-heading text-7xl tracking-tight"
                  aria-hidden="true"
                >
                  Aa Bb
                </p>
                <h3 class="mt-6 font-heading text-3xl font-medium tracking-tight">
                  Ideias que ganham forma.
                </h3>
                <p class="mt-4 text-sm leading-relaxed text-muted-foreground">
                  A fonte dos títulos e da marca. Use para criar hierarquia, com
                  espaçamento entre letras mais compacto.
                </p>
              </Card>
              <Card class="p-6 sm:p-8 shadow-none">
                <Badge variant="outline">Texto / Geist</Badge>
                <p class="mt-6 text-7xl tracking-tight" aria-hidden="true">
                  Aa Bb
                </p>
                <h3 class="mt-6 text-xl font-medium">
                  Cada detalhe tem seu lugar.
                </h3>
                <p class="mt-4 text-sm leading-relaxed text-muted-foreground">
                  A fonte das interfaces e da leitura. Use em descrições,
                  campos, navegação e ações, com altura de linha confortável.
                </p>
              </Card>
            </div>
          </section>
          <section
            id="espacamento"
            class="scroll-mt-8 border-t border-border py-12"
            aria-labelledby="spacing-title"
          >
            <p class={eyebrow}>03 / Espaçamento</p>
            <h2 id="spacing-title" class={cn(sectionTitle, "mt-3")}>
              Espaço para respirar.
            </h2>
            <p class="mt-4 text-sm text-muted-foreground">
              Uma escala em múltiplos de 4 px para organizar ritmo, distância e
              hierarquia.
            </p>
            <Card class="mt-8 space-y-4 p-6 shadow-none">
              <For each={spacing}>
                {(size) => (
                  <div class="flex items-center gap-5">
                    <span class="w-12 shrink-0 font-mono text-xs text-muted-foreground">
                      {size} px
                    </span>
                    <div
                      class="h-5 rounded-sm bg-accent"
                      style={{ width: `${size}px` }}
                      aria-hidden="true"
                    />
                    <span class="ml-auto font-mono text-xs text-muted-foreground">
                      {size / 4}
                    </span>
                  </div>
                )}
              </For>
            </Card>
          </section>
          <section
            id="formas"
            class="scroll-mt-8 border-t border-border py-12"
            aria-labelledby="shapes-title"
          >
            <p class={eyebrow}>04 / Formas</p>
            <h2 id="shapes-title" class={cn(sectionTitle, "mt-3")}>
              O mesmo sistema, outras formas.
            </h2>
            <p class="mt-4 text-sm text-muted-foreground">
              As opções de shape disponíveis nos componentes compartilhados.
            </p>
            <div class="mt-8 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
              <For each={shapes}>
                {(shape) => (
                  <Card
                    shape={shape}
                    class="flex min-h-32 items-center justify-center border-border bg-muted/50 p-6 shadow-none"
                  >
                    <span class="font-mono text-xs">{shape}</span>
                  </Card>
                )}
              </For>
            </div>
          </section>
          <section
            id="composicao"
            class="scroll-mt-8 border-t border-border py-12 sm:pb-20"
            aria-labelledby="usage-title"
          >
            <p class={eyebrow}>05 / Em uso</p>
            <h2 id="usage-title" class={cn(sectionTitle, "mt-3")}>
              Os fundamentos, juntos.
            </h2>
            <Card class="mt-8 flex flex-col items-start justify-between gap-8 p-8 shadow-none sm:flex-row sm:items-center">
              <div>
                <Badge variant="accent">TrieOH Components</Badge>
                <h3 class="mt-5 font-heading text-2xl font-medium tracking-tight">
                  Da linguagem à interface.
                </h3>
                <p class="mt-3 max-w-lg text-sm leading-relaxed text-muted-foreground">
                  Explore os componentes interativos e veja como os fundamentos
                  se traduzem em interfaces.
                </p>
              </div>
              <Link
                href="https://components.trieoh.com"
                class={cn(
                  buttonVariants({ size: "lg" }),
                  "shrink-0 no-underline hover:no-underline",
                )}
              >
                Explorar Components ↗
              </Link>
            </Card>
          </section>
        </main>
        <footer class="flex flex-wrap items-center justify-between gap-4 border-t border-border py-8 text-xs text-muted-foreground">
          <span>TrieOH / Style Guide</span>
          <Link href="https://trieoh.com">Voltar à plataforma ↗</Link>
        </footer>
      </div>
    </>
  );
}
