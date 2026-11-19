import { VirtualizedContainer } from "@trieoh/ui-solid";
import {
  createEffect,
  createMemo,
  createSignal,
  For,
  onSettled,
  Show,
} from "solid-js";

export function VirtualizedExample() {
  const [hash, setHash] = createSignal(window.location.hash);
  const query = () => new URLSearchParams(hash().split("?")[1] ?? "");
  const mode = () => (query().get("mode") === "Normal" ? "normal" : "virtual");
  const count = () => {
    const requested = Number(query().get("count"));
    return [1000, 10000, 50000].includes(requested) ? requested : 10_000;
  };
  const [mounted, setMounted] = createSignal(0);
  const items = createMemo(() =>
    Array.from({ length: count() }, (_, index) => index + 1),
  );
  let content!: HTMLDivElement;
  onSettled(() => {
    let frame = 0;
    const update = () => {
      frame = 0;
      setMounted(content.querySelectorAll("[data-comparison-item]").length);
    };
    const schedule = () => {
      if (!frame) frame = requestAnimationFrame(update);
    };
    const observer = new MutationObserver(schedule);
    observer.observe(content, { childList: true, subtree: true });
    schedule();
    return () => {
      observer.disconnect();
      cancelAnimationFrame(frame);
    };
  });
  function Item(props: { item: number }) {
    return (
      <div
        data-comparison-item
        class="rounded-md border border-border bg-card p-4 text-xs"
        style={{ "min-height": `${64 + (props.item % 3) * 16}px` }}
      >
        Item {props.item}
      </div>
    );
  }
  const columns = () =>
    Math.max(1, Math.min(3, Number(query().get("columns")) || 1));
  onSettled(() => {
    const sync = () => {
      window.scrollTo(0, 0);
      setHash(window.location.hash);
    };
    window.addEventListener("hashchange", sync);
    return () => window.removeEventListener("hashchange", sync);
  });
  createEffect(
    () => query().get("dark") === "true",
    (dark) => {
      document.documentElement.classList.toggle("dark", dark);
    },
  );
  return (
    <main class="bg-background p-4 text-foreground">
      <header class="sticky top-0 z-10 mb-4 border-b border-border bg-background py-3">
        <p class="text-xs" aria-live="polite">
          {mounted().toLocaleString("pt-BR")} de{" "}
          {count().toLocaleString("pt-BR")} itens montados no DOM
        </p>
      </header>
      <div
        ref={(element) => {
          content = element;
        }}
      >
        <Show
          when={mode() === "virtual"}
          fallback={
            <div
              style={{
                display: "grid",
                "grid-template-columns": `repeat(${columns()}, minmax(0, 1fr))`,
                gap: "8px",
              }}
            >
              <For each={items()}>{(item) => <Item item={item} />}</For>
            </div>
          }
        >
          <VirtualizedContainer
            items={items()}
            columns={columns()}
            gap={8}
            estimatedRowHeight={80}
            overscan={Number(query().get("overscan")) || 2}
          >
            {(item) => <Item item={item} />}
          </VirtualizedContainer>
        </Show>
      </div>
    </main>
  );
}
