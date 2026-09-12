import { Show, createSignal, createEffect, onCleanup, type Component } from "solid-js";
import { Portal } from "@solidjs/web";
import type { JSX } from "@solidjs/web";
import XRaw from "~icons/lucide/x";
import Loader2Raw from "~icons/lucide/loader-2";
import type { PaymentMethodI } from "./PaymentMethodSelector";

const X = XRaw as Component<{ class?: string }>;
const Loader2 = Loader2Raw as Component<{ class?: string }>;

interface PaymentFormSheetProps {
  open: boolean;
  method: PaymentMethodI | null;
  onClose: () => void;
  onReady: () => void;
  children: JSX.Element;
}

const methodLabel: Record<string, string> = {
  credit_card: 'Cartão de crédito',
  pix: 'Pix',
};

export function PaymentFormSheet(props: PaymentFormSheetProps) {
  const isMobile = () => typeof window !== 'undefined' && window.innerWidth < 768;
  const [ready, setReady] = createSignal(false);

  createEffect(
    () => props.open,
    (isOpen) => {
      if (!isOpen) {
        setReady(false);
        return;
      }

      const timer = setTimeout(() => {
        setReady(true);
        props.onReady();
      }, 500);

      onCleanup(() => clearTimeout(timer));
    },
  );

  const handleClose = () => {
    setReady(false);
    props.onClose();
  };

  return (
    <Show when={props.open}>
      <Portal>
        {/* Backdrop */}
        <div
          class="fixed inset-0 z-60 bg-background/60 backdrop-blur-sm transition-opacity duration-200"
          onClick={handleClose}
        />

        {/* Sheet */}
        <Show
          when={isMobile()}
          fallback={
            /* Desktop: panel from right */
            <div class="fixed top-0 right-0 bottom-0 z-70 w-105 flex flex-col bg-background border-l border-border shadow-2xl">
              <div class="flex items-center justify-between px-6 py-5 border-b border-border shrink-0">
                <span class="text-sm font-bold uppercase tracking-wide">
                  {props.method ? methodLabel[props.method] : ''}
                </span>
                <button
                  type="button"
                  onClick={handleClose}
                  class="p-1 hover:bg-muted rounded-full transition-colors"
                >
                  <X class="w-5 h-5 text-muted-foreground" />
                </button>
              </div>

              <div class="overflow-y-auto px-6 py-5 flex-1">
                <Show
                  when={ready()}
                  fallback={
                    <div class="flex items-center justify-center h-48">
                      <Loader2 class="w-6 h-6 animate-spin text-primary" />
                    </div>
                  }
                >
                  {props.children}
                </Show>
              </div>
            </div>
          }
        >
          {/* Mobile: sheet from bottom */}
          <div class="fixed inset-x-0 bottom-0 z-70 bg-background border-t border-border rounded-t-2xl shadow-2xl max-h-[90dvh] flex flex-col">
            <div class="flex justify-center pt-3 pb-1 shrink-0">
              <div class="w-10 h-1 rounded-full bg-muted-foreground/20" />
            </div>

            <div class="flex items-center justify-between px-5 py-3 border-b border-border shrink-0">
              <span class="text-sm font-bold uppercase tracking-wide">
                {props.method ? methodLabel[props.method] : ''}
              </span>
              <button
                type="button"
                onClick={handleClose}
                class="p-1 hover:bg-muted rounded-full transition-colors"
              >
                <X class="w-4 h-4 text-muted-foreground" />
              </button>
            </div>

            <div class="overflow-y-auto px-5 py-4 flex-1 min-h-75">
              <Show
                when={ready()}
                fallback={
                  <div class="flex items-center justify-center h-48">
                    <Loader2 class="w-6 h-6 animate-spin text-primary" />
                  </div>
                }
              >
                {props.children}
              </Show>
            </div>
          </div>
        </Show>
      </Portal>
    </Show>
  );
}
