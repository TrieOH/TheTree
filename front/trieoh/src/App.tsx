import { createEffect, createSignal, onSettled } from "solid-js";
import IconSun from "~icons/lucide/sun";
import IconMoon from "~icons/lucide/moon";
import { applyTheme, readThemePreference, saveThemePreference } from "./theme";
import {
  Button,
  Card,
  Link,
  Tooltip,
  buttonVariants,
  cn,
} from "@trieoh/ui-solid";

const componentsUrl = "https://components.trieoh.com";
const eyebrow =
  "flex items-center gap-2.5 text-[11px] font-medium uppercase tracking-[0.09em]";
const wordmark =
  "font-heading text-3xl font-bold tracking-[-1.8px] text-foreground no-underline hover:no-underline";
const navigation = "flex items-center gap-6 text-[13px] text-foreground";

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
        class="fixed top-3 left-3 z-10 translate-y-[-200%] bg-foreground p-3 text-background focus:translate-y-0 focus-visible:outline-2 focus-visible:outline-ring"
      >
        Pular para o conteúdo
      </a>
      <div class="mx-auto max-w-360 px-6 min-[761px]:px-16 max-[380px]:px-4.5">
        <header class="flex h-20 items-center justify-between border-b border-border min-[761px]:h-26">
          <Link href="/" class={wordmark} aria-label="TrieOH — início">
            TrieOH
          </Link>
          <div class="flex items-center gap-4 min-[761px]:gap-6">
            <Link href={componentsUrl} class={navigation}>
              Components <span aria-hidden="true">↗</span>
            </Link>
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
          <section
            class="pt-14.5 min-[761px]:pt-22.5"
            aria-labelledby="intro-title"
          >
            <p class={eyebrow}>
              <span
                class="size-1.75 rounded-full bg-accent"
                aria-hidden="true"
              />
              Um começo. Muitas possibilidades.
            </p>
            <h1
              id="intro-title"
              class="mt-9 mb-11 font-heading text-[clamp(58px,11vw,84px)] leading-[1.02] font-medium tracking-[-0.075em] min-[761px]:text-[clamp(68px,8.8vw,124px)]"
            >
              Ideias que
              <br />
              ganham{" "}
              <span class="text-muted-foreground">
                forma<span class="text-accent">.</span>
              </span>
            </h1>
            <div class="flex flex-col items-start justify-between gap-6 pb-10 min-[761px]:flex-row min-[761px]:items-end min-[761px]:pb-15">
              <p class="m-0 text-[15px] leading-relaxed text-muted-foreground min-[761px]:text-[17px]">
                Somos a TrieOH.
                <br />
                Um espaço para criar, conectar e construir
                <br class="hidden min-[761px]:block" /> o que vem a seguir.
              </p>
              <span class="text-[10px] uppercase tracking-[0.08em]">
                Ideias / Interfaces / Código
              </span>
            </div>
            <div
              class="flex min-h-25 items-center justify-between gap-5 bg-foreground p-5 font-heading text-background min-[761px]:min-h-31.5 min-[761px]:px-8.5"
              aria-hidden="true"
            >
              <span class="text-[32px] font-semibold tracking-[-2px] min-[761px]:text-[45px] min-[761px]:tracking-[-3px]">
                TrieOH
              </span>
              <span class="text-[42px] leading-none text-accent min-[761px]:text-7xl max-[380px]:text-[32px]">
                ✳
              </span>
              <span class="text-[clamp(16px,2.2vw,30px)] tracking-[-0.035em] max-[380px]:max-w-25">
                Da ideia à interface.
              </span>
            </div>
          </section>
          <section
            class="py-14 min-[761px]:pt-21.25 min-[761px]:pb-22.5"
            aria-labelledby="components-title"
          >
            <div class="mb-8 flex items-center justify-between border-t border-border pt-4.5 text-[11px] uppercase tracking-[0.09em] min-[761px]:mb-11.5">
              <span>01 / Explore</span>
              <span class="text-[22px]" aria-hidden="true">
                ↘
              </span>
            </div>
            <div class="grid items-center gap-10 min-[761px]:grid-cols-2 min-[761px]:gap-17.5">
              <div>
                <p class={eyebrow}>Nosso jeito de construir</p>
                <h2
                  id="components-title"
                  class="my-5.5 font-heading text-[clamp(30px,3.3vw,45px)] leading-[1.12] font-medium tracking-[-0.055em]"
                >
                  Pequenas peças.
                  <br />
                  Uma linguagem em comum.
                </h2>
                <p class="max-w-87.5 text-sm leading-[1.8] text-muted-foreground">
                  Cores, tipografia e componentes que dão forma às nossas
                  interfaces. Conheça a biblioteca e experimente cada detalhe.
                </p>
                <Link
                  href={componentsUrl}
                  class={cn(
                    buttonVariants({ size: "lg" }),
                    "mt-4 h-auto justify-between gap-9.5 rounded px-5 py-4 text-[13px] no-underline hover:no-underline",
                  )}
                >
                  Explorar Components <span aria-hidden="true">↗</span>
                </Link>
              </div>
              <Card shape="soft" class="p-6.5 shadow-none">
                <div
                  role="img"
                  aria-label="A linguagem visual TrieOH: tons neutros, azul, laranja e tipografia Space Grotesk e Geist"
                >
                  <div class="flex items-center justify-between text-[10px] text-muted-foreground">
                    <span>TrieOH / Visual language</span>
                    <span class="text-2xl text-accent" aria-hidden="true">
                      ✳
                    </span>
                  </div>
                  <div
                    class="flex justify-between gap-4 pt-5.5 pb-4 font-heading text-[clamp(60px,12vw,90px)] leading-none tracking-[-0.075em] min-[761px]:text-[clamp(40px,5.5vw,78px)]"
                    aria-hidden="true"
                  >
                    Aa<span class="text-muted-foreground">Bb</span>
                    <span class="font-sans text-accent">Cc</span>
                  </div>
                  <div class="flex items-center justify-between pb-6.5 text-[10px] text-muted-foreground">
                    <span>Space Grotesk</span>
                    <span>Geist</span>
                  </div>
                  <div
                    class="grid h-17 grid-cols-4 overflow-hidden rounded"
                    aria-hidden="true"
                  >
                    <span class="bg-foreground" />
                    <span class="bg-border" />
                    <span class="bg-secondary" />
                    <span class="bg-accent" />
                  </div>
                  <div class="flex items-center justify-between pt-5.5 text-[10px]">
                    <span>Feito de detalhes.</span>
                    <span class="text-xl" aria-hidden="true">
                      ↗
                    </span>
                  </div>
                </div>
              </Card>
            </div>
          </section>
        </main>
        <footer class="flex flex-wrap items-center justify-between gap-6 border-t border-border py-7.5 min-[761px]:flex-nowrap">
          <Link href="/" class={cn(wordmark, "text-[23px]")}>
            TrieOH
          </Link>
          <p class="order-3 m-0 basis-full text-[11px] text-muted-foreground min-[761px]:order-0 min-[761px]:basis-auto">
            O próximo capítulo está em construção.
          </p>
          <Link href={componentsUrl} class={navigation}>
            Components <span aria-hidden="true">↗</span>
          </Link>
        </footer>
      </div>
    </>
  );
}
