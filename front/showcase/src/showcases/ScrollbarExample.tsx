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
      <main class="scrollbar-example">
        <header>
          <span class="text-xs text-muted-foreground">
            Documento de exemplo
          </span>
          <h1>Notas de versão</h1>
          <p>Atualizações da biblioteca de componentes.</p>
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
            <section>
              <h2>{title}</h2>
              <p>{description}</p>
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
