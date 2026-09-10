import { For, Show, createSignal, untrack } from "solid-js";
import { Portal } from "@solidjs/web";
import type { JSX } from "@solidjs/web";
import MinusIcon from "~icons/lucide/minus";
import PlusIcon from "~icons/lucide/plus";
import ShoppingCartIcon from "~icons/lucide/shopping-cart";
import TrashIcon from "~icons/lucide/trash-2";
import XIcon from "~icons/lucide/x";
import { useCart } from "../hooks/use-cart";

const Minus = MinusIcon as unknown as () => JSX.Element;
const Plus = PlusIcon as unknown as () => JSX.Element;
const ShoppingCart = ShoppingCartIcon as unknown as () => JSX.Element;
const Trash = TrashIcon as unknown as () => JSX.Element;
const X = XIcon as unknown as () => JSX.Element;
const money = (cents: number) =>
  new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(cents / 100);

export function EventCart(props: { editionId: string; checkoutHref?: string }) {
  const [open, setOpen] = createSignal(false);
  const editionId = untrack(() => props.editionId);
  const cart = useCart(editionId);

  return (
    <>
      <Portal>
        <div
          class={`fixed inset-0 z-50 ${open() ? "visible" : "invisible pointer-events-none"}`}
        >
          <button
            type="button"
            aria-label="Fechar carrinho"
            class={`absolute inset-0 bg-black/30 backdrop-blur-sm transition-opacity ${open() ? "opacity-100" : "opacity-0"}`}
            onClick={() => setOpen(false)}
          />
          <aside
            role="dialog"
            aria-modal="true"
            aria-label="Seu carrinho"
            class={`absolute bottom-0 right-0 top-0 flex w-full max-w-100 flex-col border-l border-border bg-background shadow-2xl transition-transform duration-300 ${open() ? "translate-x-0" : "translate-x-full"}`}
          >
            <header class="flex items-center justify-between border-b border-border bg-card px-5 py-4">
              <div class="flex items-center gap-3">
                <span class="size-5 text-primary">
                  <ShoppingCart />
                </span>
                <h2 class="text-sm font-semibold uppercase tracking-wide">
                  Seu Carrinho
                </h2>
              </div>
              <button
                type="button"
                class="size-8"
                onClick={() => setOpen(false)}
              >
                <X />
              </button>
            </header>
            <div class="flex-1 overflow-y-auto">
              <Show
                when={cart.items().length > 0}
                fallback={
                  <div class="flex h-full flex-col items-center justify-center px-6 text-center text-muted-foreground">
                    <span class="mb-4 size-16 opacity-20">
                      <ShoppingCart />
                    </span>
                    <p class="text-sm font-medium uppercase tracking-wide text-foreground">
                      Carrinho vazio
                    </p>
                    <p class="mt-1 text-xs">Adicione produtos para começar</p>
                  </div>
                }
              >
                <For each={cart.items()}>
                  {(item) => (
                    <div class="flex gap-3 border-b border-border/60 bg-card p-4">
                      <div class="min-w-0 flex-1">
                        <p class="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
                          {item.type === "ticket" ? "Ingresso" : "Produto"}
                        </p>
                        <h3 class="truncate text-sm font-medium">
                          {item.name}
                        </h3>
                        <p class="mt-0.5 text-xs text-muted-foreground">
                          {money(item.price_cents)} un
                        </p>
                        <div class="mt-3 flex items-center gap-1.5">
                          <button
                            type="button"
                            class="flex size-8 items-center justify-center border border-border"
                            disabled={item.quantity <= 1}
                            onClick={() =>
                              cart.quantity(item, item.quantity - 1)
                            }
                          >
                            <Minus />
                          </button>
                          <span class="flex h-8 w-12 items-center justify-center border border-border text-sm font-semibold">
                            {item.quantity}
                          </span>
                          <button
                            type="button"
                            class="flex size-8 items-center justify-center border border-border"
                            disabled={
                              item.type === "ticket" ||
                              item.quantity >= (item.stock ?? 999)
                            }
                            onClick={() =>
                              cart.quantity(item, item.quantity + 1)
                            }
                          >
                            <Plus />
                          </button>
                        </div>
                      </div>
                      <div class="flex flex-col items-end justify-between">
                        <button
                          type="button"
                          aria-label="Remover item"
                          class="size-8 text-muted-foreground hover:text-destructive"
                          onClick={() => cart.remove(item)}
                        >
                          <Trash />
                        </button>
                        <strong class="text-sm">
                          {money(item.price_cents * item.quantity)}
                        </strong>
                      </div>
                    </div>
                  )}
                </For>
              </Show>
            </div>
            <Show when={cart.items().length > 0}>
              <footer class="border-t border-border bg-card p-5">
                <div class="mb-4 flex items-center justify-between">
                  <span class="text-sm font-semibold uppercase">Total</span>
                  <strong class="text-2xl text-primary">
                    {money(cart.total())}
                  </strong>
                </div>
                <Show when={props.checkoutHref}>
                  {(href) => (
                    <a
                      href={href()}
                      class="flex h-10 w-full items-center justify-center rounded-md bg-primary font-semibold text-primary-foreground"
                    >
                      Comprar
                    </a>
                  )}
                </Show>
                <button
                  type="button"
                  class="mt-2 h-9 w-full border border-border text-xs font-medium uppercase"
                  onClick={() => cart.clear()}
                >
                  Limpar
                </button>
              </footer>
            </Show>
          </aside>
        </div>
      </Portal>
      <button
        type="button"
        aria-label="Abrir carrinho"
        class="fixed bottom-24 right-4 z-40 flex h-13 items-center rounded-full bg-primary px-5 text-primary-foreground shadow-md transition-transform hover:scale-105 md:right-8"
        onClick={() => setOpen(true)}
      >
        <span class="mr-2 size-5">
          <ShoppingCart />
        </span>
        <span class="hidden sm:inline">Carrinho</span>
        <Show when={cart.count() > 0}>
          <span class="ml-2 rounded-full bg-background px-2 py-0.5 text-xs font-bold text-foreground">
            {cart.count()}
          </span>
        </Show>
      </button>
    </>
  );
}
