import type { JSX } from "@solidjs/web";
import { For, Show, createEffect, createSignal } from "solid-js";
import CheckIcon from "~icons/lucide/check";
import ListChecksIcon from "~icons/lucide/list-checks";
import XIcon from "~icons/lucide/x";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideCheck = CheckIcon as unknown as IconComp;
const LucideListChecks = ListChecksIcon as unknown as IconComp;
const LucideX = XIcon as unknown as IconComp;

export interface ChecklistItem {
  id: string | number;
  title: string;
  description?: string;
  completed: boolean;
  action?: { label: string; onClick: () => void; disabled?: boolean };
}

interface StepChecklistProps {
  title?: string;
  items: ChecklistItem[];
  defaultOpen?: boolean;
  mobileInline?: boolean;
  class?: string;
}

export function StepChecklist(props: StepChecklistProps): JSX.Element {
  const [open, setOpen] = createSignal(props.defaultOpen ?? false);
  const pendingCount = () => props.items.filter((item) => !item.completed).length;

  createEffect(
    () => props.mobileInline,
    (mobileInline) => {
      if (mobileInline && typeof window !== "undefined" && window.matchMedia("(max-width: 639px)").matches) {
        setOpen(true);
      }
    },
  );

  return (
    <div class={`relative inline-block ${props.class ?? ""}`}>
      <Show when={!open()}>
        <button
          type="button"
          onClick={() => setOpen(true)}
          aria-label="Abrir checklist"
          class={`${props.mobileInline ? "relative" : "fixed"} ${props.mobileInline ? "right-auto top-auto" : "right-4 top-24"} z-50 flex h-11 w-11 items-center justify-center rounded-full bg-primary text-primary-foreground shadow-lg shadow-primary/20 transition-transform hover:scale-105 hover:bg-primary/90 sm:fixed sm:right-6 sm:top-24`}
        >
          <LucideListChecks class="size-5" />
          <Show when={pendingCount() > 0}>
            <span class="absolute -right-0.5 -top-0.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-destructive px-1 text-[9px] font-bold text-destructive-foreground ring-2 ring-background">
              {pendingCount()}
            </span>
          </Show>
        </button>
      </Show>

      <Show when={open()}>
        <div
          class={`${props.mobileInline ? "relative w-full max-w-full" : "fixed right-4 top-24 w-80 max-w-[calc(100vw-2rem)]"} z-50 rounded-2xl border border-border bg-card/95 p-5 text-card-foreground shadow-xl shadow-primary/5 backdrop-blur-md sm:fixed sm:right-6 sm:top-24 sm:w-80 sm:max-w-[calc(100vw-2rem)]`}
        >
          <div class="mb-4 flex items-start justify-between gap-3">
            <Show when={props.title}>
              <h3 class="text-sm font-semibold text-foreground">
                {props.title}
              </h3>
            </Show>
            <button
              type="button"
              onClick={() => setOpen(false)}
              aria-label="Fechar checklist"
              class={`${props.mobileInline ? "hidden sm:block" : "block"} -mr-1 -mt-1 shrink-0 rounded-md p-1 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground`}
            >
              <LucideX class="size-4" />
            </button>
          </div>

          <ul class="flex flex-col">
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
                        class={`flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-[11px] font-semibold transition-colors ${item.completed ? "bg-primary text-primary-foreground" : "border-2 border-border bg-background text-muted-foreground"}`}
                      >
                        <Show when={item.completed} fallback={index() + 1}>
                          <LucideCheck class="size-3.5 stroke-[3]" />
                        </Show>
                      </span>
                      <Show when={!isLast()}>
                        <span
                          aria-hidden="true"
                          class={`my-1 w-px flex-1 transition-colors ${lineFilled() ? "bg-primary/60" : "bg-border"}`}
                        />
                      </Show>
                    </div>

                    <div
                      class={`flex min-w-0 flex-1 flex-col gap-1 ${isLast() ? "" : "pb-5"}`}
                    >
                      <div class="flex items-center justify-between gap-2">
                        <p
                          class={`text-xs font-medium ${item.completed ? "line-through text-muted-foreground" : "text-foreground"}`}
                        >
                          {item.title}
                        </p>
                        <Show when={item.action}>
                          {(action) => (
                            <button
                              type="button"
                              disabled={action().disabled}
                              onClick={() => action().onClick()}
                              class="text-[11px] font-medium text-primary hover:underline disabled:opacity-50"
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
    </div>
  );
}
