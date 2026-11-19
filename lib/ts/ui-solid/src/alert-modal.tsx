import FileIcon from "~icons/lucide/file";
import InfoIcon from "~icons/lucide/info";
import CircleAlertIcon from "~icons/lucide/circle-alert";
import TrashIcon from "~icons/lucide/trash-2";
import { Spinner } from "./spinner";
import type { Shape } from "./lib/shape";
import type { JSX } from "@solidjs/web";
import { Show, createMemo } from "solid-js";
import { Button } from "./button";
import { Dialog, type DialogSize } from "./dialog";
import { cn } from "./lib/cn";

export type AlertModalVariant = "destructive" | "warning" | "info" | "default";

export interface AlertModalProps {
  shape?: Shape;
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

function DefaultVariantIcon(props: {
  variant: () => AlertModalVariant;
}): JSX.Element {
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
                  <FileIcon class="size-4.5" aria-hidden="true" />
                </div>
              }
            >
              <div class="flex size-10 shrink-0 items-center justify-center rounded-full border border-secondary bg-secondary text-secondary-foreground shadow-xs">
                <InfoIcon class="size-4.5" aria-hidden="true" />
              </div>
            </Show>
          }
        >
          <div class="flex size-10 shrink-0 items-center justify-center rounded-full border border-accent bg-accent text-accent-foreground shadow-xs">
            <CircleAlertIcon class="size-4.5" aria-hidden="true" />
          </div>
        </Show>
      }
    >
      <div class="flex size-10 shrink-0 items-center justify-center rounded-full border border-destructive bg-destructive text-primary shadow-xs">
        <TrashIcon class="size-4.5" aria-hidden="true" />
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
  const displayTitle = createMemo<JSX.Element>((prev) =>
    props.open ? props.title : prev,
  );
  const displayDescription = createMemo<JSX.Element>((prev) =>
    props.open ? props.description : prev,
  );
  const displayIcon = createMemo<JSX.Element>((prev) =>
    props.open ? props.icon : prev,
  );
  const displayVariant = createMemo<AlertModalVariant>((prev) =>
    props.open ? (props.variant ?? "default") : (prev ?? "default"),
  );
  const displayChildren = createMemo<JSX.Element>((prev) =>
    props.open ? props.children : prev,
  );

  const handleConfirm = async () => {
    await props.onConfirm?.();
  };

  const handleCancel = () => {
    props.onCancel?.();
    props.onOpenChange(false);
  };

  return (
    <Dialog
      shape={props.shape}
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
        ...props.style,
      }}
      contentClass="p-5 sm:p-6"
    >
      <div class="flex flex-col">
        {/* Status icon badge */}
        <div>
          <Show
            when={displayIcon()}
            fallback={<DefaultVariantIcon variant={displayVariant} />}
          >
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
          <div class="mt-3 text-sm text-foreground">{displayChildren()}</div>
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
            variant={
              displayVariant() === "destructive" ? "destructive" : "default"
            }
            size="sm"
            disabled={props.loading}
            onClick={handleConfirm}
            class={cn(
              "h-9 min-w-20 px-4 text-xs sm:text-sm font-medium rounded-xl shadow-xs transition-all active:scale-95",
              displayVariant() === "warning" &&
                "bg-accent text-accent-foreground hover:bg-accent/90",
              displayVariant() === "info" &&
                "bg-secondary text-secondary-foreground hover:bg-secondary/90",
            )}
          >
            <Show when={props.loading}>
              <Spinner
                size={14}
                activeColor="currentColor"
                className="mr-1.5"
              />
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
