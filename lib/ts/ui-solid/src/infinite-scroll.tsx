import type { JSX } from "@solidjs/web";
import { createEffect, createSignal, onCleanup, Show } from "solid-js";
import { Loading } from "./loading";
import { Button } from "./button";
import { cn } from "./lib/cn";

export interface InfiniteScrollProps {
  children: JSX.Element;
  hasMore: boolean;
  loading?: boolean;
  error?: string;
  /** Promise rejection becomes a retryable error. Parent owns fetched items. */
  onLoadMore: () => void | Promise<void>;
  /** Scroll container. Omit for the document viewport. */
  scrollElement?: HTMLElement;
  /** How far ahead to request the next batch. Defaults to 300px. */
  threshold?: number;
  /** Disable automatic requests while keeping the manual button. */
  automatic?: boolean;
  endMessage?: string;
  class?: string;
}
export function InfiniteScroll(props: InfiniteScrollProps) {
  let sentinel!: HTMLDivElement;
  const [pending, setPending] = createSignal(false);
  const [failure, setFailure] = createSignal<string>();
  const busy = () => !!props.loading || pending();
  const error = () => props.error || failure();
  let disposed = false;
  // Signals flush asynchronously; guard re-entry before the next render.
  let requesting = false;
  const load = async () => {
    if (disposed || requesting || busy() || !props.hasMore) return;
    requesting = true;
    setFailure(undefined);
    setPending(true);
    try {
      await props.onLoadMore();
    } catch (error) {
      if (!disposed)
        setFailure(
          error instanceof Error
            ? error.message
            : "Não foi possível carregar mais itens.",
        );
    } finally {
      requesting = false;
      if (!disposed) setPending(false);
    }
  };
  createEffect(
    () => ({
      root: props.scrollElement,
      margin: props.threshold ?? 300,
      automatic: props.automatic !== false,
      more: props.hasMore,
      busy: busy(),
      error: error(),
    }),
    ({ root, margin, automatic, more, busy, error }) => {
      if (
        !automatic ||
        !more ||
        busy ||
        error ||
        typeof IntersectionObserver === "undefined"
      )
        return;
      const observer = new IntersectionObserver(
        (entries) => {
          if (entries.some((entry) => entry.isIntersecting)) void load();
        },
        { root: root ?? null, rootMargin: `${Math.max(0, margin)}px` },
      );
      observer.observe(sentinel);
      return () => observer.disconnect();
    },
  );
  // Cleanup lives with the mounted instance, independent of observer re-runs.
  onCleanup(() => {
    disposed = true;
  });
  return (
    <div class={cn("w-full min-w-0", props.class)}>
      <div aria-busy={busy() ? "true" : "false"}>{props.children}</div>
      <div
        ref={(element) => {
          sentinel = element;
        }}
        class="flex min-h-16 flex-col items-center justify-center gap-2 py-4"
      >
        <Show when={error()}>
          <p
            role="alert"
            class="max-w-full text-sm text-destructive wrap-anywhere"
          >
            {error()}
          </p>
        </Show>
        <Show
          when={busy()}
          fallback={
            <Show
              when={props.hasMore}
              fallback={
                <p class="text-xs text-muted-foreground">
                  {props.endMessage ?? "Você chegou ao fim da lista."}
                </p>
              }
            >
              <Button size="sm" variant="outline" onClick={() => void load()}>
                {error() ? "Tentar novamente" : "Carregar mais"}
              </Button>
            </Show>
          }
        >
          <Loading size="sm" tone="muted" label="Carregando…" />
        </Show>
      </div>
    </div>
  );
}
