import { For, createSignal } from "solid-js";
import type { JSX } from "@solidjs/web";
import { cn } from "./lib/cn";

export interface AccordionItem {
  value: string;
  title: JSX.Element;
  content: JSX.Element;
  disabled?: boolean;
}

export interface AccordionProps {
  items: readonly AccordionItem[];
  defaultValue?: string;
  multiple?: boolean;
  class?: string;
}

export function Accordion(props: AccordionProps): JSX.Element {
  const [open, setOpen] = createSignal<string[]>(
    props.defaultValue ? [props.defaultValue] : [],
  );

  const isOpen = (value: string) => open().includes(value);

  const toggle = (value: string) => {
    setOpen((current) => {
      if (props.multiple) {
        return current.includes(value)
          ? current.filter((item) => item !== value)
          : [...current, value];
      }

      return current.includes(value) ? [] : [value];
    });
  };

  return (
    <div class={cn("w-full divide-y divide-border", props.class)}>
      <For each={props.items}>
        {(item) => (
          <div class="border-b border-border last:border-b-0">
            <button
              type="button"
              disabled={item.disabled}
              aria-expanded={isOpen(item.value) ? "true" : "false"}
              class="flex w-full items-center justify-between gap-3 py-4 text-left text-sm font-medium transition-colors hover:text-foreground/80 disabled:pointer-events-none disabled:opacity-50 cursor-pointer"
              onClick={() => toggle(item.value)}
            >
              <span class="min-w-0 flex-1">{item.title}</span>

              <span
                class={cn(
                  "flex shrink-0 items-center justify-center transition-transform duration-200 text-muted-foreground",
                  isOpen(item.value) && "rotate-45",
                )}
              >
                <svg
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  class="size-4"
                  aria-hidden="true"
                >
                  <path d="M5 12h14M12 5v14" />
                </svg>
              </span>
            </button>

            <div
              class={cn(
                "grid transition-[grid-template-rows,opacity] duration-250 ease-out",
                isOpen(item.value)
                  ? "grid-rows-[1fr] opacity-100"
                  : "grid-rows-[0fr] opacity-0",
              )}
            >
              <div class="min-h-0 overflow-hidden">
                <div class="pb-4 text-sm text-muted-foreground">
                  {item.content}
                </div>
              </div>
            </div>
          </div>
        )}
      </For>
    </div>
  );
}
