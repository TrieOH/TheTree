import type { JSX } from "@solidjs/web";
import { Portal } from "@solidjs/web";
import { Show, createEffect, createMemo, createSignal, onSettled } from "solid-js";
import { cn } from "./lib/cn";

export type DialogSize = "sm" | "md" | "lg" | "xl" | "full";

export interface DialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title?: JSX.Element;
  description?: JSX.Element;
  children?: JSX.Element;
  footer?: JSX.Element;
  header?: JSX.Element;
  size?: DialogSize;
  /** Extra classes for the panel (width, padding overrides…). */
  class?: string;
  /** Extra classes for the content body. */
  contentClass?: string;
  /** Inline style overrides for the panel. */
  style?: JSX.CSSProperties;
  /** Disables Escape key and backdrop click closing. */
  preventClose?: boolean;
  /** Displays close X button in header. Default true. */
  closable?: boolean;
  /** WAI-ARIA role, defaults to 'dialog'. */
  role?: "dialog" | "alertdialog";
  ariaLabel?: string;
}

const FOCUSABLE =
  'a[href], button:not([disabled]), textarea:not([disabled]), input:not([disabled]), select:not([disabled]), [tabindex]:not([tabindex="-1"])';

const SIZE_CLASSES: Record<DialogSize, string> = {
  sm: "max-w-sm",
  md: "max-w-lg",
  lg: "max-w-2xl",
  xl: "max-w-4xl",
  full: "max-w-[calc(100vw-2rem)]",
};

/**
 * Modern modal dialog hand-crafted for high performance, buttery-smooth scrolling,
 * zero layout shift, centered mobile & desktop presentation with margin clearance,
 * and clean layering above admin sidebars.
 */
export function Dialog(props: DialogProps): JSX.Element {
  const [mounted, setMounted] = createSignal(false);
  const [active, setActive] = createSignal(false);

  createEffect(
    () => props.open,
    (isOpen) => {
      if (isOpen) {
        setMounted(true);
        if (typeof window !== "undefined" && typeof window.requestAnimationFrame === "function") {
          const raf = requestAnimationFrame(() => setActive(true));
          return () => cancelAnimationFrame(raf);
        } else {
          setActive(true);
        }
      } else {
        setActive(false);
        if (typeof window !== "undefined") {
          const timer = setTimeout(() => setMounted(false), 200);
          return () => clearTimeout(timer);
        } else {
          setMounted(false);
        }
      }
    },
  );

  return (
    <Show when={mounted()}>
      <DialogSurface dialogProps={props} active={active()} />
    </Show>
  );
}

function DialogSurface(surfaceProps: { dialogProps: DialogProps; active: boolean }): JSX.Element {
  const props = surfaceProps.dialogProps;
  const [panel, setPanel] = createSignal<HTMLDivElement>();
  const idPrefix = Math.random().toString(36).slice(2, 9);
  const titleId = `dialog-title-${idPrefix}`;
  const descriptionId = `dialog-desc-${idPrefix}`;

  // Idiomatic Solid: createMemo with accumulator (prev) freezes content when open === false,
  // preventing stale state flashes during exit transitions without mutable `let` references.
  const displayTitle = createMemo<JSX.Element>((prev) => (props.open ? props.title : prev));
  const displayDescription = createMemo<JSX.Element>((prev) => (props.open ? props.description : prev));
  const displayChildren = createMemo<JSX.Element>((prev) => (props.open ? props.children : prev));
  const displayFooter = createMemo<JSX.Element>((prev) => (props.open ? props.footer : prev));
  const displayHeader = createMemo<JSX.Element>((prev) => (props.open ? props.header : prev));

  const onKeyDown = (event: KeyboardEvent) => {
    if (event.key === "Escape") {
      if (props.preventClose) return;
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

  // Solid 2 canonical lifecycle: onSettled with returned cleanup (no onCleanup outside owner)
  onSettled(() => {
    if (typeof window === "undefined") return;
    const previousOverflow = document.body.style.overflow;
    const previousPaddingRight = document.body.style.paddingRight;
    const previouslyFocused = document.activeElement as HTMLElement | null;

    // Compensate scrollbar width to prevent desktop layout shift
    const scrollbarWidth = window.innerWidth - document.documentElement.clientWidth;
    if (scrollbarWidth > 0) {
      document.body.style.paddingRight = `${scrollbarWidth}px`;
    }
    document.body.style.overflow = "hidden";

    window.addEventListener("keydown", onKeyDown);
    panel()?.focus({ preventScroll: true });

    return () => {
      window.removeEventListener("keydown", onKeyDown);
      document.body.style.overflow = previousOverflow;
      document.body.style.paddingRight = previousPaddingRight;
      previouslyFocused?.focus?.();
    };
  });

  const sizeClass = () => SIZE_CLASSES[props.size ?? "md"];

  const handleBackdropClick = () => {
    if (!props.preventClose) {
      props.onOpenChange(false);
    }
  };

  const handleCloseBtnClick = () => {
    if (!props.preventClose) {
      props.onOpenChange(false);
    }
  };

  return (
    <Portal>
      {/*
        Top-level container at z-90:
        Centered vertically and horizontally on all viewports, with p-4 padding
        so the dialog card is NEVER glued to screen edges on mobile or desktop.
      */}
      <div class="fixed inset-0 z-90 flex items-center justify-center p-4 sm:p-6 overflow-hidden">
        {/*
          Backdrop overlay with subtle, GPU-friendly blur.
          Separated from modal panel so scrolling content never triggers continuous blur re-renders.
        */}
        <div
          aria-hidden="true"
          class={cn(
            "fixed inset-0 bg-black/60 backdrop-blur-xs transition-opacity duration-200 ease-out sm:backdrop-blur-sm",
            surfaceProps.active ? "opacity-100 pointer-events-auto" : "opacity-0 pointer-events-none",
          )}
          onClick={handleBackdropClick}
        />

        {/*
          Modern Dialog Panel:
          - Centered card layout with rounded-2xl on mobile and desktop.
          - max-height: calc(100dvh - 2rem) guarantees 16px safety clearance top and bottom on any screen.
          - Clean border and subtle shadow for depth without performance penalty.
        */}
        <div
          ref={setPanel}
          role={props.role ?? "dialog"}
          aria-modal="true"
          aria-label={props.ariaLabel}
          aria-labelledby={props.title ? titleId : undefined}
          aria-describedby={props.description ? descriptionId : undefined}
          tabindex="-1"
          style={{
            "max-height": "calc(100dvh - 2rem)",
            ...(props.style ?? {}),
          }}
          class={cn(
            "relative z-10 flex w-full flex-col overflow-hidden rounded-2xl border border-border/70 bg-card text-card-foreground shadow-2xl shadow-black/25 ring-1 ring-black/5 transition-all duration-200 ease-out dark:shadow-black/70 dark:ring-white/10",
            sizeClass(),
            surfaceProps.active
              ? "translate-y-0 opacity-100 scale-100"
              : "translate-y-3 opacity-0 scale-[0.97] sm:scale-[0.98]",
            props.class,
          )}
          onClick={(e) => e.stopPropagation()}
        >
          {/* Header */}
          <Show
            when={displayHeader() || displayTitle() || displayDescription()}
            fallback={
              <Show when={props.closable !== false}>
                <div class="absolute right-3.5 top-3.5 z-20">
                  <button
                    type="button"
                    aria-label="Fechar"
                    class="inline-flex size-8 shrink-0 items-center justify-center rounded-xl text-muted-foreground/70 transition-all duration-150 hover:bg-muted hover:text-foreground active:scale-95 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring cursor-pointer"
                    onClick={handleCloseBtnClick}
                  >
                    <svg
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      class="size-4"
                      aria-hidden="true"
                    >
                      <path d="M18 6 6 18M6 6l12 12" stroke-linecap="round" />
                    </svg>
                  </button>
                </div>
              </Show>
            }
          >
            <Show
              when={displayHeader()}
              fallback={
                <header class="flex shrink-0 items-start justify-between gap-3 border-b border-border/40 px-5 py-3.5 sm:px-6 sm:py-4">
                  <div class="flex min-w-0 flex-1 flex-col gap-0.5 pr-2">
                    <Show when={displayTitle()}>
                      <h2
                        id={titleId}
                        class="text-base sm:text-lg font-semibold tracking-tight text-foreground leading-tight"
                      >
                        {displayTitle()}
                      </h2>
                    </Show>
                    <Show when={displayDescription()}>
                      <p
                        id={descriptionId}
                        class="text-xs sm:text-sm text-muted-foreground leading-normal"
                      >
                        {displayDescription()}
                      </p>
                    </Show>
                  </div>

                  <Show when={props.closable !== false}>
                    <button
                      type="button"
                      aria-label="Fechar"
                      class="-mr-1.5 -mt-1 inline-flex size-8 shrink-0 items-center justify-center rounded-xl text-muted-foreground/70 transition-all duration-150 hover:bg-muted hover:text-foreground active:scale-95 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring cursor-pointer"
                      onClick={handleCloseBtnClick}
                    >
                      <svg
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        class="size-4"
                        aria-hidden="true"
                      >
                        <path d="M18 6 6 18M6 6l12 12" stroke-linecap="round" />
                      </svg>
                    </button>
                  </Show>
                </header>
              }
            >
              {displayHeader()}
            </Show>
          </Show>

          {/*
            Scrollable Content Body:
            - min-h-0 allows flex child to shrink properly between pinned header and footer.
            - overscroll-contain keeps momentum scrolling isolated to the dialog.
            - -webkit-overflow-scrolling: touch ensures buttery-smooth mobile momentum scrolling.
          */}
          <div
            class={cn(
              "flex-1 min-h-0 overflow-y-auto overscroll-contain p-5 sm:p-6 isolate",
              props.contentClass,
            )}
            style={{ "-webkit-overflow-scrolling": "touch" }}
          >
            {displayChildren()}
          </div>

          {/* Footer */}
          <Show when={displayFooter()}>
            <footer class="flex shrink-0 flex-wrap items-center justify-end gap-2.5 border-t border-border/40 bg-muted/20 px-5 py-3.5 sm:px-6 sm:py-4">
              {displayFooter()}
            </footer>
          </Show>
        </div>
      </div>
    </Portal>
  );
}
