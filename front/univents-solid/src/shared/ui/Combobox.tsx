import type { JSX } from "@solidjs/web";
import { For, Show, createMemo, createSignal } from "solid-js";
import CheckIcon from "~icons/lucide/check";
import ChevronsUpDownIcon from "~icons/lucide/chevrons-up-down";
import SearchIcon from "~icons/lucide/search";

const Check = CheckIcon as unknown as () => JSX.Element;
const ChevronsUpDown = ChevronsUpDownIcon as unknown as () => JSX.Element;
const Search = SearchIcon as unknown as () => JSX.Element;
export interface ComboboxOption {
  value: string;
  label: string;
  description?: string;
}

export function Combobox(props: {
  value?: string;
  options: readonly ComboboxOption[];
  placeholder: string;
  searchPlaceholder?: string;
  disabled?: boolean;
  onChange: (value: string) => void;
}) {
  const [open, setOpen] = createSignal(false);
  const [query, setQuery] = createSignal("");
  const selected = createMemo(() =>
    props.options.find((item) => item.value === props.value),
  );
  const visible = createMemo(() => {
    return props.options.filter((item) =>
      `${item.label} ${item.description ?? ""}`
        .toLocaleLowerCase("pt-BR")
        .includes(query().trim().toLocaleLowerCase("pt-BR")),
    );
  });
  return (
    <div class="relative w-full">
      <button
        type="button"
        disabled={props.disabled}
        aria-haspopup="listbox"
        aria-expanded={open() ? "true" : "false"}
        onClick={() => setOpen((value) => !value)}
        class="flex h-10 w-full items-center justify-between gap-2 rounded-md border border-border bg-background px-3 text-left text-sm outline-none hover:bg-muted disabled:pointer-events-none disabled:opacity-50"
      >
        <span class={`truncate ${selected() ? "" : "text-muted-foreground"}`}>
          {selected()?.label ?? props.placeholder}
        </span>
        <ChevronsUpDown />
      </button>
      <Show when={open()}>
        <div class="absolute inset-x-0 top-full z-50 mt-1 overflow-hidden rounded-md border border-border bg-popover text-popover-foreground shadow-xl">
          <div class="flex items-center gap-2 border-b border-border px-3 py-2">
            <Search />
            <input
              value={query()}
              placeholder={props.searchPlaceholder ?? "Buscar…"}
              onInput={(event) => setQuery(event.currentTarget.value)}
              class="min-w-0 flex-1 bg-transparent text-sm outline-none"
            />
          </div>
          <div class="max-h-60 overflow-y-auto p-1" role="listbox">
            <Show
              when={visible().length > 0}
              fallback={
                <p class="px-3 py-2 text-sm text-muted-foreground">
                  Nenhum resultado encontrado
                </p>
              }
            >
              <For each={visible()}>
                {(option) => (
                  <button
                    type="button"
                    role="option"
                    aria-selected={
                      option.value === props.value ? "true" : "false"
                    }
                    onClick={() => {
                      props.onChange(option.value);
                      setOpen(false);
                      setQuery("");
                    }}
                    class="flex w-full items-center justify-between gap-2 rounded-sm px-3 py-2 text-left text-sm hover:bg-muted"
                  >
                    <span class="truncate">{option.label}</span>
                    <Show when={option.value === props.value}>
                      <Check />
                    </Show>
                  </button>
                )}
              </For>
            </Show>
          </div>
        </div>
      </Show>
    </div>
  );
}
