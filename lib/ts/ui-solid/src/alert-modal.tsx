import type { JSX } from "@solidjs/web";
import { Show, createMemo } from "solid-js";
import { Button } from "./button";
import { Dialog, type DialogSize } from "./dialog";
import { cn } from "./lib/cn";

export type AlertModalVariant = "destructive" | "warning" | "info" | "default";

export interface AlertModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: JSX.Element;
  description?: JSX.Element;
  children?: JSX.Element;
  confirmLabel?: string;
  cancelLabel?: string;
  variant?: AlertModalVariant;
  loading?: boolean;
  onConfirm?: () => Promise<void> | void;
  onCancel?: () => void;
  icon?: JSX.Element;
  size?: DialogSize;
  class?: string;
  style?: JSX.CSSProperties;
}

function DefaultVariantIcon(props: { variant: () => AlertModalVariant }): JSX.Element {
  return (
    <Show
      when={props.variant() === "destructive"}
      fallback={
        <Show
          when={props.variant() === "warning"}
          fallback={
            <Show
              when={props.variant() === "info"}
              fallback={
                <div class="flex size-10 shrink-0 items-center justify-center rounded-xl border border-primary/30 bg-primary/10 text-primary shadow-xs">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="size-4.5" aria-hidden="true">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke-linecap="round" stroke-linejoin="round" />
                    <polyline points="14 2 14 8 20 8" stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </div>
              }
            >
              <div class="flex size-10 shrink-0 items-center justify-center rounded-full border border-sky-500/30 bg-sky-500/10 text-sky-600 dark:text-sky-400 shadow-xs">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="size-4.5" aria-hidden="true">
                  <circle cx="12" cy="12" r="10" stroke-linecap="round" stroke-linejoin="round" />
                  <path d="M12 16v-4M12 8h.01" stroke-linecap="round" />
                </svg>
              </div>
            </Show>
          }
        >
          <div class="flex size-10 shrink-0 items-center justify-center rounded-full border border-amber-500/30 bg-amber-500/10 text-amber-600 dark:text-amber-400 shadow-xs">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="size-4.5" aria-hidden="true">
              <circle cx="12" cy="12" r="10" stroke-linecap="round" stroke-linejoin="round" />
              <path d="M12 8v4M12 16h.01" stroke-linecap="round" />
            </svg>
          </div>
        </Show>
      }
    >
      <div class="flex size-10 shrink-0 items-center justify-center rounded-full border border-destructive/30 bg-destructive/10 text-destructive shadow-xs">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="size-4.5" aria-hidden="true">
          <path d="M3 6h18M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M10 11v6M14 11v6" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
      </div>
    </Show>
  );
}

/**
 * Clean, modern AlertModal strictly wired to the application's configured design tokens
 * (`bg-card`, `text-card-foreground`, `border-border`, `bg-destructive`, `bg-primary`, etc.),
 * ensuring perfect contrast in light and dark mode with zero arbitrary color drift.
 * Automatically freezes displayed content during exit transition to prevent parent state teardown flashes.
 */
export function AlertModal(props: AlertModalProps): JSX.Element {
  const confirmLabel = () => props.confirmLabel ?? "Confirmar";
  const cancelLabel = () => props.cancelLabel ?? "Cancelar";

  // Idiomatic Solid: createMemo with accumulator (prev) freezes content when open === false,
  // preventing stale state flashes during exit transitions without mutable `let` references.
  const displayTitle = createMemo<JSX.Element>((prev) => (props.open ? props.title : prev));
  const displayDescription = createMemo<JSX.Element>((prev) => (props.open ? props.description : prev));
  const displayIcon = createMemo<JSX.Element>((prev) => (props.open ? props.icon : prev));
  const displayVariant = createMemo<AlertModalVariant>((prev) => (props.open ? (props.variant ?? "default") : (prev ?? "default")));
  const displayChildren = createMemo<JSX.Element>((prev) => (props.open ? props.children : prev));

  const handleConfirm = async () => {
    await props.onConfirm?.();
  };

  const handleCancel = () => {
    props.onCancel?.();
    props.onOpenChange(false);
  };

  return (
    <Dialog
      open={props.open}
      onOpenChange={(v) => {
        if (!props.loading) {
          props.onOpenChange(v);
        }
      }}
      role="alertdialog"
      size="sm"
      preventClose={props.loading}
      closable={false}
      class={cn(
        "w-full max-w-[340px] rounded-2xl border border-border bg-card text-card-foreground shadow-2xl",
        props.class,
      )}
      style={{
        "max-width": "340px",
        ...(props.style ?? {}),
      }}
      contentClass="p-5 sm:p-6"
    >
      <div class="flex flex-col">
        {/* Status icon badge */}
        <div>
          <Show when={displayIcon()} fallback={<DefaultVariantIcon variant={displayVariant} />}>
            {(customIcon) => <div class="shrink-0">{customIcon()}</div>}
          </Show>
        </div>

        {/* Title and description */}
        <div class="mt-3.5 flex flex-col gap-0.5">
          <h3 class="text-base font-semibold tracking-tight text-foreground leading-tight">
            {displayTitle()}
          </h3>
          <Show when={displayDescription()}>
            <p class="text-xs sm:text-sm text-muted-foreground leading-normal">
              {displayDescription()}
            </p>
          </Show>
        </div>

        {/* Optional custom children */}
        <Show when={displayChildren()}>
          <div class="mt-3 text-sm text-foreground">
            {displayChildren()}
          </div>
        </Show>

        {/* Action buttons strictly using theme tokens */}
        <div class="mt-6 flex items-center justify-end gap-2.5">
          <Button
            variant="ghost"
            size="sm"
            disabled={props.loading}
            onClick={handleCancel}
            class="h-9 px-3.5 text-xs sm:text-sm font-medium text-muted-foreground hover:bg-muted hover:text-foreground rounded-xl"
          >
            {cancelLabel()}
          </Button>
          <Button
            variant={displayVariant() === "destructive" ? "destructive" : "default"}
            size="sm"
            disabled={props.loading}
            onClick={handleConfirm}
            class={cn(
              "h-9 min-w-20 px-4 text-xs sm:text-sm font-medium rounded-xl shadow-xs transition-all active:scale-95",
              displayVariant() === "warning" && "bg-amber-600 text-white hover:bg-amber-700 dark:bg-amber-500 dark:text-zinc-950 dark:hover:bg-amber-400",
              displayVariant() === "info" && "bg-sky-600 text-white hover:bg-sky-700 dark:bg-sky-500 dark:text-zinc-950 dark:hover:bg-sky-400",
            )}
          >
            <Show when={props.loading}>
              <svg
                class="mr-1.5 size-3.5 animate-spin"
                viewBox="0 0 24 24"
                fill="none"
              >
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                />
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                />
              </svg>
            </Show>
            <Show when={props.loading} fallback={confirmLabel()}>
              Processando...
            </Show>
          </Button>
        </div>
      </div>
    </Dialog>
  );
}
