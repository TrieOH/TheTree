import type { JSX } from "@solidjs/web";
import { Show, createSignal, onCleanup } from "solid-js";

import ArrowLeftIcon from "~icons/lucide/arrow-left";
import MonitorIcon from "~icons/lucide/monitor";

import { Button } from "@trieoh/ui-solid";

const Monitor = MonitorIcon as unknown as (props: { class?: string }) => JSX.Element;
const ArrowLeft = ArrowLeftIcon as unknown as (props: { class?: string }) => JSX.Element;

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
            {/* Clean, sober monitor icon */}
            <div class="mb-5 flex size-12 items-center justify-center rounded-xl border border-border bg-muted/40 text-muted-foreground">
              <Monitor class="size-6" />
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

            {/* Subtle resolution readout */}
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
                <ArrowLeft class="size-3.5" />
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
