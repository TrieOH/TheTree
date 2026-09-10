import { For, createSignal } from "solid-js";
import type { JSX } from "@solidjs/web";
import PlusIcon from "~icons/lucide/plus";

const Plus = PlusIcon as unknown as () => JSX.Element;

export interface AccordionItem {
  value: string;
  title: JSX.Element;
  content: JSX.Element;
}

export function Accordion(props: {
  items: AccordionItem[];
  defaultValue?: string;
  multiple?: boolean;
}) {
  const [open, setOpen] = createSignal<string[]>(
    props.defaultValue ? [props.defaultValue] : [],
  );

  const isOpen = (value: string) =>
    open().includes(value);

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
    <div>
      <For each={props.items}>
        {(item) => (
          <div class="border-b border-border last:border-b-0">
            <button
              type="button"
              aria-expanded={
                isOpen(item.value) ? "true" : "false"
              }
              class="flex w-full items-center justify-between gap-3 py-4 text-left text-sm font-medium"
              onClick={() => toggle(item.value)}
            >
              <span class="min-w-0 flex-1">
                {item.title}
              </span>

              <span
                class={`flex shrink-0 items-center justify-center transition-transform duration-200 [&>svg]:size-4 ${isOpen(item.value)
                    ? "rotate-45"
                    : ""
                  }`}
              >
                <Plus />
              </span>
            </button>

            <div
              class={`grid transition-[grid-template-rows,opacity] duration-250 ease-out ${isOpen(item.value)
                  ? "grid-rows-[1fr] opacity-100"
                  : "grid-rows-[0fr] opacity-0"
                }`}
            >
              <div class="min-h-0 overflow-hidden">
                <div class="pb-4">
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