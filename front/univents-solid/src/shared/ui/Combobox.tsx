import { Portal, type JSX } from "@solidjs/web";
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

const Check = CheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const ChevronsUpDown = ChevronsUpDownIcon as unknown as (props: { class?: string }) => JSX.Element;
const Search = SearchIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface ComboboxOption {
  value: string;
  label: string;
  description?: string;
}

export interface ComboboxProps {
  value?: string;
  options: readonly ComboboxOption[];
  placeholder: string;
  searchPlaceholder?: string;
  disabled?: boolean;
  class?: string;
  className?: string;
  triggerClass?: string;
  triggerClassName?: string;
  dropdownClass?: string;
  dropdownClassName?: string;
  onChange: (value: string) => void;
}

export function Combobox(props: ComboboxProps): JSX.Element {
  const [open, setOpen] = createSignal(false);
  const [query, setQuery] = createSignal("");
  const [highlightedIndex, setHighlightedIndex] = createSignal(0);
  const [position, setPosition] = createSignal({ top: 0, left: 0, minWidth: 0 });

  let containerRef: HTMLDivElement | undefined;
  let dropdownRef: HTMLDivElement | undefined;
  let inputRef: HTMLInputElement | undefined;

  const updatePosition = () => {
    if (!containerRef) return;
    const rect = containerRef.getBoundingClientRect();
    const spaceBelow = window.innerHeight - rect.bottom;
    const estimatedHeight = 280;

    let top = rect.bottom + 4;
    if (spaceBelow < estimatedHeight && rect.top > estimatedHeight) {
      top = rect.top - estimatedHeight - 4;
    }

    setPosition({
      top: Math.max(top, 8),
      left: Math.max(rect.left, 8),
      minWidth: rect.width,
    });
  };

  createEffect(
    () => open(),
    (isOpen) => {
      if (!isOpen) return;

      updatePosition();

      const handlePointerDown = (event: MouseEvent) => {
        const target = event.target as Node | null;
        if (
          (containerRef && containerRef.contains(target)) ||
          (dropdownRef && dropdownRef.contains(target))
        ) {
          return;
        }
        setOpen(false);
      };

      const handleScrollOrResize = () => {
        updatePosition();
      };

      document.addEventListener("mousedown", handlePointerDown);
      window.addEventListener("resize", handleScrollOrResize);
      window.addEventListener("scroll", handleScrollOrResize, true);

      // Focus input
      setTimeout(() => inputRef?.focus(), 10);

      return () => {
        document.removeEventListener("mousedown", handlePointerDown);
        window.removeEventListener("resize", handleScrollOrResize);
        window.removeEventListener("scroll", handleScrollOrResize, true);
      };
    },
  );

  createEffect(
    () => query(),
    () => {
      setHighlightedIndex(0);
    },
  );

  const selected = createMemo(() =>
    props.options.find((item) => item.value === props.value),
  );

  const visible = createMemo(() => {
    const list: ComboboxOption[] = [];
    const q = query().trim().toLocaleLowerCase("pt-BR");
    for (const item of props.options) {
      if (
        `${item.label} ${item.description ?? ""}`
          .toLocaleLowerCase("pt-BR")
          .includes(q)
      ) {
        list.push(item);
      }
    }
    return list;
  });

  const rootClass = () => props.class ?? props.className ?? "w-full";
  const triggerCls = () =>
    props.triggerClass ?? props.triggerClassName ?? "h-10";
  const dropdownCls = () =>
    props.dropdownClass ?? props.dropdownClassName ?? "w-full";

  return (
    <div
      ref={(el) => (containerRef = el)}
      class={`relative min-w-0 ${rootClass()}`}
    >
      <button
        type="button"
        disabled={props.disabled}
        aria-haspopup="listbox"
        aria-expanded={open() ? "true" : "false"}
        onClick={() => setOpen((v) => !v)}
        class={`flex w-full items-center justify-between gap-2 rounded-md border border-border bg-background px-3 text-left text-sm outline-none hover:bg-muted disabled:pointer-events-none disabled:opacity-50 ${triggerCls()}`}
      >
        <span class={`truncate ${selected() ? "" : "text-muted-foreground"}`}>
          {selected()?.label ?? props.placeholder}
        </span>
        <ChevronsUpDown class="size-3.5 shrink-0 text-muted-foreground" />
      </button>

      <Show when={open()}>
        <Portal>
          <div
            ref={(el) => (dropdownRef = el)}
            style={{
              position: "fixed",
              top: `${position().top}px`,
              left: `${position().left}px`,
              "min-width": `${position().minWidth}px`,
            }}
            class={`fixed z-70 w-max max-w-[calc(100vw-1rem)] overflow-hidden rounded-lg border border-border bg-popover text-popover-foreground shadow-2xl ${dropdownCls()}`}
          >
            <div class="flex items-center gap-2 border-b border-border px-3 py-2">
              <Search class="size-3.5 shrink-0 text-muted-foreground" />
              <input
                ref={(el) => (inputRef = el)}
                value={query()}
                placeholder={props.searchPlaceholder ?? "Buscar…"}
                onInput={(event) => setQuery(event.currentTarget.value)}
                onKeyDown={(event) => {
                  if (event.key === "Escape") {
                    setOpen(false);
                  } else if (event.key === "ArrowDown") {
                    event.preventDefault();
                    setHighlightedIndex((idx) =>
                      Math.min(idx + 1, Math.max(visible().length - 1, 0)),
                    );
                  } else if (event.key === "ArrowUp") {
                    event.preventDefault();
                    setHighlightedIndex((idx) => Math.max(idx - 1, 0));
                  } else if (event.key === "Enter") {
                    event.preventDefault();
                    const item = visible()[highlightedIndex()];
                    if (item) {
                      props.onChange(item.value);
                      setOpen(false);
                      setQuery("");
                    }
                  }
                }}
                class="min-w-0 flex-1 bg-transparent text-sm outline-none placeholder:text-muted-foreground"
              />
            </div>
            <ul class="max-h-56 overflow-y-auto p-1" role="listbox">
              <Show
                when={visible().length > 0}
                fallback={
                  <li class="px-3 py-2 text-xs text-muted-foreground">
                    Nenhum resultado encontrado
                  </li>
                }
              >
                <For each={visible()}>
                  {(option, idx) => (
                    <li>
                      <button
                        type="button"
                        role="option"
                        aria-selected={
                          option.value === props.value ? "true" : "false"
                        }
                        onMouseEnter={() => setHighlightedIndex(idx())}
                        onClick={() => {
                          props.onChange(option.value);
                          setOpen(false);
                          setQuery("");
                        }}
                        class={`flex w-full items-center justify-between gap-2 rounded-sm px-3 py-2 text-left text-sm transition-colors ${
                          idx() === highlightedIndex()
                            ? "bg-primary/10"
                            : "hover:bg-muted"
                        }`}
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
        </Portal>
      </Show>
    </div>
  );
}
