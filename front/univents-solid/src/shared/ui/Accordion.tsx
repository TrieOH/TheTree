import { For, createSignal } from "solid-js";
import type { JSX } from "@solidjs/web";
import PlusIcon from "~icons/lucide/plus";

const Plus = PlusIcon as unknown as () => JSX.Element;

export interface AccordionItem {
  value: string;
  title: string;
  content: string;
}

export function Accordion(props: { items: AccordionItem[] }) {
  const [open, setOpen] = createSignal<string>();

  return (
    <div>
      <For each={props.items}>
        {(item) => {
          const isOpen = () => open() === item.value;
          return (
            <div class="border-b border-border last:border-b-0">
              <button
                type="button"
                aria-expanded={isOpen() ? "true" : "false"}
                class="flex w-full items-center justify-between py-4 text-left text-sm font-medium"
                onClick={() => setOpen(isOpen() ? undefined : item.value)}
              >
                {item.title}
                <span
                  class={`transition-transform duration-200 ${isOpen() ? "rotate-45" : ""}`}
                >
                  <Plus />
                </span>
              </button>
              <div
                class={`grid transition-[grid-template-rows,opacity] duration-250 ease-out ${isOpen() ? "grid-rows-[1fr] opacity-100" : "grid-rows-[0fr] opacity-0"}`}
              >
                <div class="min-h-0 overflow-hidden">
                  <p class="pb-4 text-sm text-muted-foreground">
                    {item.content}
                  </p>
                </div>
              </div>
            </div>
          );
        }}
      </For>
    </div>
  );
}
