import type { JSX } from "@solidjs/web";
import { Show, createSignal, onCleanup } from "solid-js";
import { Button } from "./button";

export interface DesktopOnlyProps {
  children?: JSX.Element;
  /**
   * Width threshold in pixels. Defaults to 1200.
   */
  threshold?: number;
  /**
   * Custom title for the mobile/tablet blocking notice.
   */
  title?: string;
  /**
   * Custom description text.
   */
  description?: string;
  /**
   * Label for the back/exit button. Defaults to "Voltar à página anterior".
   */
  actionLabel?: string;
  /**
   * Callback when the user clicks the action button.
   * Defaults to window.history.back().
   */
  onAction?: () => void;
}

/**
 * Reusable guard component that restricts access to desktop screens.
 * When screen width is below the specified threshold (default 1200px),
 * renders a clean, distraction-free page informing the user to open on desktop.
 */
export function DesktopOnly(props: DesktopOnlyProps): JSX.Element {
  const thresholdPx = () => props.threshold ?? 1200;
  const [currentWidth, setCurrentWidth] = createSignal(
    typeof window !== "undefined" ? window.innerWidth : 1200,
  );
  const [isDesktop, setIsDesktop] = createSignal(
    typeof window !== "undefined" ? window.innerWidth >= thresholdPx() : true,
  );

  if (typeof window !== "undefined") {
    const updateMatch = () => {
      const w = window.innerWidth;
      setCurrentWidth(w);
      setIsDesktop(w >= thresholdPx());
    };

    window.addEventListener("resize", updateMatch, { passive: true });
    onCleanup(() => {
      window.removeEventListener("resize", updateMatch);
    });
  }

  const handleAction = () => {
    if (props.onAction) {
      props.onAction();
    } else if (typeof window !== "undefined") {
      window.history.back();
    }
  };

  return (
    <Show
      when={isDesktop()}
      fallback={
        <div class="flex min-h-dvh w-full flex-col items-center justify-center bg-background px-6 py-12 select-none">
          <div class="flex w-full max-w-md flex-col items-center text-center">
            {/* Monitor icon */}
            <div class="mb-5 flex size-12 items-center justify-center rounded-xl border border-border bg-muted/40 text-muted-foreground">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="size-6"
                aria-hidden="true"
              >
                <rect width="20" height="14" x="2" y="3" rx="2" />
                <line x1="8" x2="16" y1="21" y2="21" />
                <line x1="12" x2="12" y1="17" y2="21" />
              </svg>
            </div>

            {/* Title */}
            <h1 class="text-lg font-semibold tracking-tight text-foreground">
              {props.title ?? "Disponível apenas no computador"}
            </h1>

            {/* Description */}
            <p class="mt-2 text-sm text-muted-foreground leading-normal">
              {props.description ??
                "Esta ferramenta foi projetada para telas a partir de 1200px de largura para permitir o gerenciamento da grade e organização dos horários."}
            </p>

            {/* Resolution readout */}
            <div class="mt-4 flex items-center gap-2 text-xs text-muted-foreground font-mono">
              <span>{currentWidth()}px</span>
              <span>/</span>
              <span>mínimo de {thresholdPx()}px</span>
            </div>

            {/* Back action */}
            <div class="mt-6">
              <Button
                type="button"
                variant="outline"
                onClick={handleAction}
                class="gap-2 text-xs font-medium cursor-pointer"
              >
                <svg
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  class="size-3.5"
                  aria-hidden="true"
                >
                  <path d="m12 19-7-7 7-7" />
                  <path d="M19 12H5" />
                </svg>
                {props.actionLabel ?? "Voltar à página anterior"}
              </Button>
            </div>
          </div>
        </div>
      }
    >
      {props.children}
    </Show>
  );
}
