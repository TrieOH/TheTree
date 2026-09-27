import type { JSX } from "@solidjs/web";
import { For, Show, createSignal } from "solid-js";
import CheckIcon from "~icons/lucide/check";
import ListChecksIcon from "~icons/lucide/list-checks";
import XIcon from "~icons/lucide/x";
import { Button, cn } from "@trieoh/ui-solid";
import { createDragToDismiss } from "./hooks/create-drag-to-dismiss";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideCheck = CheckIcon as unknown as IconComp;
const LucideListChecks = ListChecksIcon as unknown as IconComp;
const LucideX = XIcon as unknown as IconComp;

export interface StepChecklistItem {
  id: string | number;
  title: string;
  description?: string;
  completed: boolean;
  action?: {
    label: string;
    onClick: () => void;
    disabled?: boolean;
  };
}

// Alias for compatibility with React naming
export type ChecklistItem = StepChecklistItem;

export interface StepChecklistProps {
  items: StepChecklistItem[];
  title?: string;
  mobileInline?: boolean;
  class?: string;
  className?: string;
}

export function StepChecklist(props: StepChecklistProps): JSX.Element {
  const [open, setOpen] = createSignal(false);

  const drag = createDragToDismiss({
    threshold: 70,
    onDismiss: () => setOpen(false),
  });

  const completedCount = () => props.items.filter((i) => i.completed).length;
  const progressPercent = () =>
    props.items.length === 0
      ? 0
      : Math.round((completedCount() / props.items.length) * 100);

  const isAllCompleted = () =>
    props.items.length > 0 && completedCount() === props.items.length;

  return (
    <>
      {/* Floating Trigger Button - discreet when all completed, prominent when pending steps */}
      <Show when={!open()}>
        <div
          class={cn(
            "fixed bottom-20 right-4 z-40 sm:bottom-6 sm:right-6",
            props.class,
            props.className,
          )}
        >
          <Button
            type="button"
            onClick={() => setOpen(true)}
            class={cn(
              "flex items-center gap-2 rounded-full border backdrop-blur-md transition-all active:scale-95",
              isAllCompleted()
                ? "border-border/50 bg-card/60 px-3 py-1.5 text-muted-foreground opacity-60 shadow-xs hover:border-border hover:bg-card hover:text-foreground hover:opacity-100"
                : "border-border/80 bg-card/95 px-4 py-2.5 text-card-foreground shadow-lg hover:scale-105 hover:bg-accent",
            )}
            aria-label="Abrir checklist de configuração"
          >
            <Show
              when={isAllCompleted()}
              fallback={<LucideListChecks class="size-4 text-primary" />}
            >
              <LucideCheck class="size-3.5 text-emerald-600 dark:text-emerald-400" />
            </Show>
            <span
              class={cn(
                "text-xs",
                isAllCompleted() ? "font-normal" : "font-semibold",
              )}
            >
              {isAllCompleted()
                ? `Configurado (${completedCount()}/${props.items.length})`
                : `Configuração (${completedCount()}/${props.items.length})`}
            </span>
          </Button>
        </div>
      </Show>

      {/* Menu / Panel that appears */}
      <Show when={open()}>
        {/* Mobile Backdrop - z-[65] to ensure it sits safely above the mobile navigation dock (z-50) */}
        <div
          data-testid="backdrop"
          class="fixed inset-0 z-65 bg-black/60 backdrop-blur-xs sm:hidden"
          onClick={() => setOpen(false)}
          aria-hidden="true"
        />

        {/* Panel: Bottom sheet drawer on mobile (z-[70]), floating card at bottom-right on desktop */}
        <div
          style={drag.style()}
          class={cn(
            "fixed inset-x-0 bottom-0 z-70 flex max-h-[85dvh] flex-col rounded-t-3xl border-t border-border bg-card p-5 shadow-2xl backdrop-blur-md",
            "sm:inset-auto sm:bottom-6 sm:right-6 sm:max-h-[calc(100vh-5rem)] sm:w-88 sm:rounded-2xl sm:border sm:shadow-2xl",
            props.class,
            props.className,
          )}
        >
          {/* Mobile swipe/drag handle touch/pointer area */}
          <div
            data-testid="drag-handle"
            class="flex w-full cursor-grab touch-none items-center justify-center py-2 -mt-2 mb-1 active:cursor-grabbing sm:hidden"
            {...drag.dragProps}
          >
            <div class="h-1.5 w-12 shrink-0 rounded-full bg-muted-foreground/30 transition-colors hover:bg-muted-foreground/50" />
          </div>

          {/* Header - also draggable on mobile */}
          <div
            class="mb-3 flex shrink-0 select-none items-start justify-between gap-3 touch-none sm:touch-auto"
            {...drag.dragProps}
          >
            <div class="flex flex-col gap-1">
              <Show when={props.title}>
                <h3 class="text-sm font-semibold text-foreground">
                  {props.title}
                </h3>
              </Show>
              <p
                class={cn(
                  "text-xs",
                  isAllCompleted()
                    ? "font-medium text-emerald-600 dark:text-emerald-400"
                    : "text-muted-foreground",
                )}
              >
                {isAllCompleted()
                  ? "Tudo pronto! Todas as etapas concluídas."
                  : `${completedCount()} de ${props.items.length} concluídos (${progressPercent()}%)`}
              </p>
            </div>
            <button
              type="button"
              onClick={() => setOpen(false)}
              aria-label="Fechar checklist"
              class="-mr-1 -mt-1 shrink-0 rounded-full p-1.5 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
            >
              <LucideX class="size-4" />
            </button>
          </div>

          {/* Progress bar */}
          <div class="mb-4 h-1.5 w-full shrink-0 overflow-hidden rounded-full bg-muted">
            <div
              class={cn(
                "h-full transition-all duration-300",
                isAllCompleted() ? "bg-emerald-500" : "bg-primary",
              )}
              style={{ width: `${progressPercent()}%` }}
            />
          </div>

          {/* Items List - scrollable without excessive padding */}
          <ul class="min-h-0 flex-1 flex flex-col overflow-y-auto pr-1 overscroll-contain">
            <For each={props.items}>
              {(item, index) => {
                const isLast = () => index() === props.items.length - 1;
                const nextCompleted = () =>
                  !isLast() && props.items[index() + 1]?.completed;
                const lineFilled = () => item.completed && nextCompleted();

                return (
                  <li class="flex gap-3">
                    <div class="flex flex-col items-center">
                      <span
                        class={cn(
                          "flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-[11px] font-semibold transition-colors",
                          item.completed
                            ? "bg-primary text-primary-foreground"
                            : "border-2 border-border bg-background text-muted-foreground",
                        )}
                      >
                        <Show when={item.completed} fallback={index() + 1}>
                          <LucideCheck class="size-3.5 stroke-3" />
                        </Show>
                      </span>
                      <Show when={!isLast()}>
                        <span
                          aria-hidden="true"
                          class={cn(
                            "my-1 w-px flex-1 transition-colors",
                            lineFilled() ? "bg-primary/60" : "bg-border",
                          )}
                        />
                      </Show>
                    </div>

                    <div
                      class={cn(
                        "flex min-w-0 flex-1 flex-col gap-1",
                        !isLast() && "pb-5",
                      )}
                    >
                      <div class="flex items-center justify-between gap-2">
                        <p
                          class={cn(
                            "text-xs font-medium",
                            item.completed
                              ? "line-through text-muted-foreground"
                              : "text-foreground",
                          )}
                        >
                          {item.title}
                        </p>
                        <Show when={item.action}>
                          {(action) => (
                            <button
                              type="button"
                              disabled={action().disabled}
                              onClick={() => action().onClick()}
                              class="text-[11px] font-medium text-primary hover:underline disabled:cursor-not-allowed disabled:text-muted-foreground disabled:no-underline disabled:opacity-60"
                            >
                              {action().label}
                            </button>
                          )}
                        </Show>
                      </div>
                      <Show when={item.description}>
                        <p class="text-[11px] text-muted-foreground">
                          {item.description}
                        </p>
                      </Show>
                    </div>
                  </li>
                );
              }}
            </For>
          </ul>
        </div>
      </Show>
    </>
  );
}
