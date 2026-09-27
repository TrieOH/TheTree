import type { JSX } from "@solidjs/web";
import {
  For,
  Show,
  createEffect,
  createMemo,
  createSignal,
} from "solid-js";
import CheckIcon from "~icons/lucide/check";
import ChevronsUpDownIcon from "~icons/lucide/chevrons-up-down";
import SearchIcon from "~icons/lucide/search";
import { cn } from "@trieoh/ui-solid";

const Check = CheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const ChevronsUpDown = ChevronsUpDownIcon as unknown as (props: {
  class?: string;
}) => JSX.Element;
const Search = SearchIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface ToolbarComboboxOption {
  value: string;
  label: string;
  description?: string;
}

export interface ToolbarComboboxProps {
  value?: string;
  options: readonly ToolbarComboboxOption[];
  placeholder: string;
  searchPlaceholder?: string;
  disabled?: boolean;
  icon?: JSX.Element;
  iconOnly?: boolean;
  class?: string;
  className?: string;
  triggerClass?: string;
  triggerClassName?: string;
  dropdownClass?: string;
  dropdownClassName?: string;
  onChange: (value: string) => void;
}

export function ToolbarCombobox(props: ToolbarComboboxProps): JSX.Element {
  let containerRef: HTMLDivElement | undefined;
  const [open, setOpen] = createSignal(false);
  const [query, setQuery] = createSignal("");
  const [highlightedIndex, setHighlightedIndex] = createSignal(0);

  const normalizedQuery = createMemo(() => query().trim().toLowerCase());
  const selectedOption = createMemo(() =>
    props.options.find((option) => option.value === props.value),
  );

  const visibleOptions = createMemo(() => {
    if (!normalizedQuery()) return props.options;
    return props.options.filter((option) =>
      `${option.label} ${option.description ?? ""}`
        .toLowerCase()
        .includes(normalizedQuery()),
    );
  });

  createEffect(
    () => open(),
    (isOpen) => {
      if (!isOpen) return;

      function closeOnOutsideClick(event: MouseEvent) {
        if (!containerRef?.contains(event.target as Node)) {
          setOpen(false);
        }
      }

      document.addEventListener("mousedown", closeOnOutsideClick);
      return () => {
        document.removeEventListener("mousedown", closeOnOutsideClick);
      };
    },
  );

  createEffect(
    () => [open(), query()],
    () => {
      queueMicrotask(() => {
        setHighlightedIndex(0);
      });
    },
  );

  function select(option: ToolbarComboboxOption) {
    props.onChange(option.value);
    setOpen(false);
    setQuery("");
  }

  return (
    <div
      ref={(el) => (containerRef = el)}
      class={cn(
        "relative min-w-0 shrink-0 rounded-md border border-border",
        props.class,
        props.className,
      )}
    >
      <button
        type="button"
        disabled={props.disabled}
        aria-haspopup="listbox"
        aria-expanded={open() ? "true" : "false"}
        aria-label={props.iconOnly ? props.placeholder : undefined}
        title={props.iconOnly ? props.placeholder : undefined}
        class={cn(
          "flex h-7 items-center rounded-md bg-background text-xs transition-colors hover:bg-muted disabled:pointer-events-none disabled:opacity-45 cursor-pointer",
          props.triggerClass,
          props.triggerClassName,
          props.iconOnly
            ? "w-7 justify-center p-0"
            : "w-full justify-between gap-2 px-2",
        )}
        onClick={() => setOpen((current) => !current)}
      >
        <Show
          when={!props.iconOnly}
          fallback={props.icon}
        >
          <span
            class={cn(
              "truncate",
              !selectedOption() && "text-muted-foreground",
            )}
          >
            {selectedOption()?.label ?? props.placeholder}
          </span>
          <ChevronsUpDown class="size-3.5 shrink-0 opacity-50" />
        </Show>
      </button>

      <Show when={open()}>
        <div
          class={cn(
            "absolute top-[calc(100%+4px)] left-0 z-50 min-w-44 max-w-72 rounded-md border border-border bg-popover text-popover-foreground shadow-md",
            props.dropdownClass,
            props.dropdownClassName,
          )}
        >
          <div class="flex items-center gap-1.5 border-b border-border px-2 py-1.5">
            <Search class="size-3.5 text-muted-foreground" />
            <input
              type="text"
              value={query()}
              placeholder={props.searchPlaceholder ?? "Buscar…"}
              class="h-7 min-w-0 flex-1 bg-transparent text-sm outline-none placeholder:text-muted-foreground"
              onInput={(event) => setQuery(event.currentTarget.value)}
              onKeyDown={(event) => {
                if (event.key === "ArrowDown") {
                  event.preventDefault();
                  setHighlightedIndex((index) =>
                    Math.min(index + 1, Math.max(visibleOptions().length - 1, 0)),
                  );
                } else if (event.key === "ArrowUp") {
                  event.preventDefault();
                  setHighlightedIndex((index) => Math.max(index - 1, 0));
                } else if (event.key === "Enter") {
                  event.preventDefault();
                  const list = visibleOptions();
                  if (list.length > 0 && list[highlightedIndex()]) {
                    select(list[highlightedIndex()]);
                  }
                } else if (event.key === "Escape") {
                  setOpen(false);
                }
              }}
            />
          </div>
          <ul class="max-h-56 overflow-y-auto py-1">
            <Show
              when={visibleOptions().length > 0}
              fallback={
                <li class="px-3 py-2 text-xs text-muted-foreground">
                  Nenhum resultado encontrado
                </li>
              }
            >
              <For each={visibleOptions()}>
                {(option, index) => (
                  <li>
                    <button
                      type="button"
                      role="option"
                      aria-selected={option.value === props.value ? "true" : "false"}
                      class={cn(
                        "flex w-full items-center justify-between gap-2 px-3 py-2.5 text-left text-sm transition-colors cursor-pointer",
                        index() === highlightedIndex()
                          ? "bg-primary/10"
                          : "hover:bg-muted",
                      )}
                      onMouseEnter={() => setHighlightedIndex(index())}
                      onClick={() => select(option)}
                    >
                      <span class="min-w-0">
                        <span class="block truncate">{option.label}</span>
                        <Show when={option.description}>
                          <span class="block truncate text-xs text-muted-foreground">
                            {option.description}
                          </span>
                        </Show>
                      </span>
                      <Show when={option.value === props.value}>
                        <Check class="size-3.5 shrink-0" />
                      </Show>
                    </button>
                  </li>
                )}
              </For>
            </Show>
          </ul>
        </div>
      </Show>
    </div>
  );
}
