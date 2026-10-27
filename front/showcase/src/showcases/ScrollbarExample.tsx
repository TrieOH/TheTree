import {
  createEffect,
  createMemo,
  createSignal,
  onSettled,
  For,
  Show,
} from "solid-js";
import { Button, OverlayScrollbar } from "@trieoh/ui-solid";

/** Separate document so the real page scrollbar can be tested in a compact preview. */
export function ScrollbarExample() {
  const [hash, setHash] = createSignal(window.location.hash);
  const options = createMemo(() => new URLSearchParams(hash().split("?")[1]));
  onSettled(() => {
    const sync = () => {
      setHash(window.location.hash);
    };
    window.addEventListener("hashchange", sync);
    return () => window.removeEventListener("hashchange", sync);
  });
  createEffect(
    () => options().get("dark") === "true",
    (dark) => {
      document.documentElement.classList.toggle("dark", dark);
    },
  );
  return (
    <>
      <Show when={options().get("enabled") !== "false"}>
        <OverlayScrollbar />
      </Show>
      <main class="mx-auto max-w-160 px-7 pt-6 pb-8">
        <header class="pb-6">
          <span class="text-xs text-muted-foreground">
            Documento de exemplo
          </span>
          <h1 class="my-2 text-2xl font-semibold tracking-[-0.6px]">
            Notas de versão
          </h1>
          <p class="text-[13px] leading-[1.8] text-muted-foreground">
            Atualizações da biblioteca de componentes.
          </p>
        </header>
        <For
          each={[
            [
              "Formulários",
              "Campos com rótulos, descrições e mensagens de validação.",
            ],
            [
              "Navegação",
              "Abas de seção, listas expansíveis e seleção com busca.",
            ],
            [
              "Modais",
              "Diálogos de confirmação e formulários divididos em etapas.",
            ],
            [
              "Feedback",
              "Notificações e estados de carregamento para operações assíncronas.",
            ],
            [
              "Acessibilidade",
              "Navegação por teclado e foco visível nos controles.",
            ],
          ]}
        >
          {([title, description]) => (
            <section class="border-t border-border py-6">
              <h2 class="mb-2 text-sm font-semibold">{title}</h2>
              <p class="text-[13px] leading-[1.8] text-muted-foreground">
                {description}
              </p>
            </section>
          )}
        </For>
        <Button
          variant="outline"
          size="sm"
          onClick={() => window.scrollTo({ top: 0, behavior: "smooth" })}
        >
          Voltar ao início
        </Button>
      </main>
    </>
  );
}
