import type { JSX } from "@solidjs/web";
import { Portal } from "@solidjs/web";
import { Show, createSignal, onCleanup, onSettled } from "solid-js";
import { cn } from "./lib/cn";

export interface DialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: JSX.Element;
  description?: JSX.Element;
  children?: JSX.Element;
  footer?: JSX.Element;
  /** Extra classes for the panel (width, padding overrides…). */
  class?: string;
}

const FOCUSABLE =
  'a[href], button:not([disabled]), textarea:not([disabled]), input:not([disabled]), select:not([disabled]), [tabindex]:not([tabindex="-1"])';

/**
 * Modal dialog, hand-rolled: no shadcn, no headless kit.
 *
 * Owns the four things that are easy to get wrong and tedious to repeat —
 * escape to close, a focus trap, scroll lock, and focus restore — and nothing
 * else. Rendered only while `open`, so all of that is bound to the dialog's
 * lifetime instead of leaking listeners.
 */
export function Dialog(props: DialogProps) {
  return (
    <Show when={props.open}>
      <DialogSurface {...props} />
    </Show>
  );
}

function DialogSurface(props: DialogProps): JSX.Element {
  const [panel, setPanel] = createSignal<HTMLDivElement>();
  const titleId = `dialog-title-${Math.random().toString(36).slice(2, 9)}`;
  const descriptionId = `${titleId}-description`;

  let previousOverflow: string | undefined;
  let previouslyFocused: HTMLElement | null = null;

  const onKeyDown = (event: KeyboardEvent) => {
    if (event.key === "Escape") {
      event.stopPropagation();
      props.onOpenChange(false);
      return;
    }
    if (event.key !== "Tab") return;

    const element = panel();
    if (!element) return;

    const focusable = Array.from(element.querySelectorAll<HTMLElement>(FOCUSABLE));
    const first = focusable[0];
    const last = focusable[focusable.length - 1];
    if (!first || !last) {
      event.preventDefault();
      return;
    }

    const active = document.activeElement;
    if (event.shiftKey && (active === first || active === element)) {
      event.preventDefault();
      last.focus();
    } else if (!event.shiftKey && active === last) {
      event.preventDefault();
      first.focus();
    }
  };

  // Imperative work lives in `onSettled` (Solid 2's `onMount`, and the idiom
  // this codebase uses) while the cleanup is registered on the component owner.
  // `onCleanup` *inside* a callback is a `NO_OWNER_CLEANUP` diagnostic here.
  onSettled(() => {
    previousOverflow = document.body.style.overflow;
    previouslyFocused = document.activeElement as HTMLElement | null;
    document.body.style.overflow = "hidden";
    window.addEventListener("keydown", onKeyDown);
    panel()?.focus();
  });

  onCleanup(() => {
    if (typeof window === "undefined") return;
    window.removeEventListener("keydown", onKeyDown);
    if (previousOverflow !== undefined) {
      document.body.style.overflow = previousOverflow;
    }
    previouslyFocused?.focus?.();
  });

  return (
    <Portal>
      <div class="fixed inset-0 z-50 flex items-end justify-center sm:items-center">
        <button
          type="button"
          aria-label="Fechar"
          class="absolute inset-0 cursor-default bg-black/50 backdrop-blur-sm"
          onClick={() => props.onOpenChange(false)}
        />

        <div
          ref={setPanel}
          role="dialog"
          aria-modal="true"
          aria-labelledby={titleId}
          aria-describedby={props.description ? descriptionId : undefined}
          tabindex="-1"
          class={cn(
            "relative z-10 flex max-h-[90vh] w-full flex-col overflow-hidden rounded-t-2xl border border-border bg-card shadow-2xl sm:max-w-lg sm:rounded-2xl",
            props.class,
          )}
        >
          <header class="flex items-start justify-between gap-4 border-b border-border p-5">
            <div class="flex flex-col gap-1">
              <h2 id={titleId} class="text-base font-semibold text-foreground">
                {props.title}
              </h2>
              <Show when={props.description}>
                <p id={descriptionId} class="text-sm text-muted-foreground">
                  {props.description}
                </p>
              </Show>
            </div>

            <button
              type="button"
              aria-label="Fechar"
              class="flex size-8 shrink-0 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
              onClick={() => props.onOpenChange(false)}
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="size-4" aria-hidden="true">
                <path d="M18 6 6 18M6 6l12 12" stroke-linecap="round" />
              </svg>
            </button>
          </header>

          <div class="flex-1 overflow-y-auto p-5">{props.children}</div>

          <Show when={props.footer}>
            <footer class="flex flex-wrap items-center justify-end gap-2 border-t border-border p-5">
              {props.footer}
            </footer>
          </Show>
        </div>
      </div>
    </Portal>
  );
}
