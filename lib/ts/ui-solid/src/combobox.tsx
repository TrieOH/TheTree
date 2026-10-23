import { Portal, type JSX } from "@solidjs/web";
import { For, Show, createEffect, createMemo, createSignal } from "solid-js";
import { cn } from "./lib/cn";

export interface ComboboxOption {
  value: string;
  label: string;
  description?: string;
  disabled?: boolean;
}

export interface ComboboxProps {
  id?: string;
  "aria-describedby"?: string;
  "aria-label"?: string;
  value?: string;
  options: readonly ComboboxOption[];
  placeholder: string;
  searchPlaceholder?: string;
  emptyMessage?: string;
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
  const [position, setPosition] = createSignal({ top: 0, left: 0, width: 0 });

  let containerRef: HTMLDivElement | undefined;
  let dropdownRef: HTMLDivElement | undefined;
  let inputRef: HTMLInputElement | undefined;

  const updatePosition = () => {
    if (!containerRef) return;
    const rect = containerRef.getBoundingClientRect();
    const margin = 8;
    const width = Math.min(rect.width, window.innerWidth - margin * 2);
    const height = Math.min(
      dropdownRef?.getBoundingClientRect().height || 280,
      window.innerHeight - margin * 2,
    );
    const below = window.innerHeight - rect.bottom - margin;
    const above = rect.top - margin;
    const top =
      below < height + 4 && above > below
        ? rect.top - height - 4
        : rect.bottom + 4;

    setPosition({
      top: Math.max(
        margin,
        Math.min(top, window.innerHeight - height - margin),
      ),
      left: Math.max(
        margin,
        Math.min(rect.left, window.innerWidth - width - margin),
      ),
      width,
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

      // Measure the mounted portal instead of assuming a 280px menu.
      const observer = new ResizeObserver(updatePosition);
      const frame = requestAnimationFrame(() => {
        updatePosition();
        if (dropdownRef) observer.observe(dropdownRef);
        inputRef?.focus();
      });

      return () => {
        cancelAnimationFrame(frame);
        observer.disconnect();
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
    const q = query().trim().toLocaleLowerCase();
    for (const item of props.options) {
      if (
        `${item.label} ${item.description ?? ""}`
          .toLocaleLowerCase()
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
  const dropdownCls = () => props.dropdownClass ?? props.dropdownClassName;

  return (
    <div
      ref={(el) => (containerRef = el)}
      class={cn("relative min-w-0", rootClass())}
    >
      <button
        type="button"
        disabled={props.disabled}
        id={props.id}
        aria-describedby={props["aria-describedby"]}
        aria-label={props["aria-label"]}
        aria-haspopup="listbox"
        aria-expanded={open() ? "true" : "false"}
        onClick={() => setOpen((v) => !v)}
        class={cn(
          "flex w-full items-center justify-between gap-2 rounded-md border border-border bg-background px-3 text-left text-sm outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 transition-colors hover:bg-muted disabled:pointer-events-none disabled:opacity-50 cursor-pointer",
          triggerCls(),
        )}
      >
        <span
          class={cn(
            "truncate",
            selected() ? "text-foreground" : "text-muted-foreground",
          )}
        >
          {selected()?.label ?? props.placeholder}
        </span>
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="size-3.5 shrink-0 text-muted-foreground"
          aria-hidden="true"
        >
          <path d="m7 15 5 5 5-5M7 9l5-5 5 5" />
        </svg>
      </button>

      <Show when={open()}>
        <Portal>
          <div
            ref={(el) => (dropdownRef = el)}
            style={{
              position: "fixed",
              top: `${position().top}px`,
              left: `${position().left}px`,
              width: `${position().width}px`,
              "max-height": "calc(100dvh - 16px)",
              "overflow-y": "auto",
            }}
            class={cn(
              "fixed z-70 max-w-[calc(100vw-1rem)] overflow-hidden rounded-lg border border-border bg-popover text-popover-foreground shadow-2xl animate-in fade-in-0 zoom-in-95",
              dropdownCls(),
            )}
          >
            <div class="flex items-center gap-2 border-b border-border px-3 py-2">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="size-3.5 shrink-0 text-muted-foreground"
                aria-hidden="true"
              >
                <circle cx="11" cy="11" r="8" />
                <path d="m21 21-4.3-4.3" />
              </svg>
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
                    if (item && !item.disabled) {
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
                    {props.emptyMessage ?? "Nenhum resultado encontrado"}
                  </li>
                }
              >
                <For each={visible()}>
                  {(option, idx) => (
                    <li>
                      <button
                        type="button"
                        role="option"
                        disabled={option.disabled}
                        aria-selected={
                          option.value === props.value ? "true" : "false"
                        }
                        onMouseEnter={() => setHighlightedIndex(idx())}
                        onClick={() => {
                          props.onChange(option.value);
                          setOpen(false);
                          setQuery("");
                        }}
                        class={cn(
                          "flex w-full items-center justify-between gap-2 rounded-sm px-3 py-2 text-left text-sm transition-colors cursor-pointer disabled:pointer-events-none disabled:opacity-50",
                          idx() === highlightedIndex()
                            ? "bg-primary/10 text-primary"
                            : "hover:bg-muted text-foreground",
                        )}
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
                          <svg
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            class="size-3.5 shrink-0"
                            aria-hidden="true"
                          >
                            <path d="M20 6 9 17l-5-5" />
                          </svg>
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
