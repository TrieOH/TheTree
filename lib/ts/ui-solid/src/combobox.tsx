import { shapeClass, type Shape } from "./lib/shape";
import { Portal, type JSX } from "@solidjs/web";
import { For, Show, createEffect, createMemo, createSignal } from "solid-js";
import ChevronsUpDownIcon from "~icons/lucide/chevrons-up-down";
import SearchIcon from "~icons/lucide/search";
import CheckIcon from "~icons/lucide/check";
import { cn } from "./lib/cn";

export interface ComboboxOption {
  value: string;
  label: string;
  description?: string;
  disabled?: boolean;
}

export interface ComboboxProps {
  shape?: Shape;
  size?: "default" | "sm";
  id?: string;
  "aria-describedby"?: string;
  "aria-label"?: string;
  value?: string;
  options: readonly ComboboxOption[];
  placeholder: string;
  searchPlaceholder?: string;
  searchable?: boolean;
  emptyMessage?: string;
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

export function Combobox(props: ComboboxProps): JSX.Element {
  const [open, setOpen] = createSignal(false);
  const [query, setQuery] = createSignal("");
  const [highlightedIndex, setHighlightedIndex] = createSignal(0);
  const [position, setPosition] = createSignal({ top: 0, left: 0, width: 0 });

  let containerRef: HTMLDivElement | undefined;
  let dropdownRef: HTMLDivElement | undefined;
  let inputRef: HTMLInputElement | undefined;

  const isSm = () => props.size === "sm";

  const updatePosition = () => {
    if (!containerRef) return;
    const rect = containerRef.getBoundingClientRect();
    const margin = 8;
    const minMenuWidth = props.iconOnly ? 190 : isSm() ? 150 : 0;
    const width = Math.max(
      minMenuWidth,
      Math.min(rect.width, window.innerWidth - margin * 2),
    );
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

      const handleKeyDown = (event: KeyboardEvent) => {
        if (event.key === "Escape") {
          setOpen(false);
        }
      };

      const handleScrollOrResize = () => {
        updatePosition();
      };

      document.addEventListener("pointerdown", handlePointerDown);
      document.addEventListener("keydown", handleKeyDown);
      window.removeEventListener("scroll", handleScrollOrResize, true);
      window.removeEventListener("resize", handleScrollOrResize);
      window.addEventListener("scroll", handleScrollOrResize, true);
      window.addEventListener("resize", handleScrollOrResize);

      setTimeout(() => inputRef?.focus(), 0);

      return () => {
        document.removeEventListener("pointerdown", handlePointerDown);
        document.removeEventListener("keydown", handleKeyDown);
        window.removeEventListener("scroll", handleScrollOrResize, true);
        window.removeEventListener("resize", handleScrollOrResize);
      };
    },
  );

  const selectedOption = createMemo(() =>
    props.options.find((option) => option.value === props.value),
  );

  const visible = createMemo(() => {
    const q = query().trim().toLowerCase();
    if (!q) return props.options;
    return props.options.filter(
      (opt) =>
        opt.label.toLowerCase().includes(q) ||
        (opt.description && opt.description.toLowerCase().includes(q)),
    );
  });

  const showSearch = () =>
    props.searchable ??
    (props.options.length > 5 || Boolean(props.searchPlaceholder));

  const triggerCls = createMemo(() =>
    cn(
      "flex w-full items-center justify-between border border-border bg-background text-sm text-foreground transition-colors outline-none",
      "hover:bg-muted focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
      "disabled:cursor-not-allowed disabled:opacity-50 cursor-pointer",
      shapeClass(props.shape),
      isSm() ? "h-8 px-2.5 text-xs" : "h-9 px-3 text-sm",
      props.triggerClass ?? props.triggerClassName,
    ),
  );

  return (
    <div
      ref={(el) => (containerRef = el)}
      class={cn(
        "relative inline-block w-full",
        props.class ?? props.className,
      )}
    >
      <Show
        when={!props.iconOnly}
        fallback={
          <button
            type="button"
            id={props.id}
            aria-describedby={props["aria-describedby"]}
            aria-label={props["aria-label"] ?? props.placeholder}
            aria-expanded={open() ? "true" : "false"}
            aria-haspopup="listbox"
            disabled={props.disabled}
            onClick={() => setOpen(!open())}
            class={cn(
              "flex items-center justify-center border border-border bg-background text-foreground transition-colors outline-none",
              "hover:bg-muted focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
              "disabled:cursor-not-allowed disabled:opacity-50 cursor-pointer",
              shapeClass(props.shape),
              isSm() ? "size-8" : "size-9",
              props.triggerClass ?? props.triggerClassName,
            )}
          >
            {props.icon}
          </button>
        }
      >
        <button
          type="button"
          id={props.id}
          aria-describedby={props["aria-describedby"]}
          aria-label={props["aria-label"]}
          aria-expanded={open() ? "true" : "false"}
          aria-haspopup="listbox"
          disabled={props.disabled}
          onClick={() => setOpen(!open())}
          class={triggerCls()}
        >
          <span class="flex min-w-0 items-center gap-2 truncate">
            <Show when={props.icon}>
              <span class="shrink-0">{props.icon}</span>
            </Show>
            <span
              class={cn(
                "block truncate",
                !selectedOption() && "text-muted-foreground",
              )}
            >
              {selectedOption()?.label ?? props.placeholder}
            </span>
          </span>
          <ChevronsUpDownIcon class="size-4 shrink-0 text-muted-foreground opacity-50" />
        </button>
      </Show>

      <Show when={open()}>
        <Portal>
          <div
            ref={(el) => (dropdownRef = el)}
            style={{
              position: "fixed",
              top: `${position().top}px`,
              left: `${position().left}px`,
              width: `${position().width}px`,
            }}
            class={cn(
              "z-50 overflow-hidden border border-border bg-popover text-popover-foreground shadow-md animate-in fade-in-0 zoom-in-95",
              shapeClass(props.shape),
              props.dropdownClass ?? props.dropdownClassName,
            )}
          >
            <Show when={showSearch()}>
              <div class="flex items-center border-b border-border px-2.5 py-1.5">
                <SearchIcon class="mr-2 size-3.5 shrink-0 text-muted-foreground opacity-50" />
                <input
                  ref={(el) => (inputRef = el)}
                  value={query()}
                  onInput={(event) => {
                    setQuery(event.currentTarget.value);
                    setHighlightedIndex(0);
                  }}
                  placeholder={props.searchPlaceholder ?? "Buscar…"}
                  onKeyDown={(event) => {
                    if (event.key === "ArrowDown") {
                      event.preventDefault();
                      setHighlightedIndex((idx) =>
                        Math.min(idx + 1, visible().length - 1),
                      );
                    } else if (event.key === "ArrowUp") {
                      event.preventDefault();
                      setHighlightedIndex((idx) => Math.max(idx - 1, 0));
                    } else if (event.key === "Enter") {
                      event.preventDefault();
                      const list = visible();
                      const item = list[highlightedIndex()];
                      if (item && !item.disabled) {
                        props.onChange(item.value);
                        setOpen(false);
                        setQuery("");
                      }
                    }
                  }}
                  class="h-7 min-w-0 flex-1 bg-transparent text-sm outline-none placeholder:text-muted-foreground"
                />
              </div>
            </Show>
            <ul
              role="listbox"
              class="max-h-60 overflow-y-auto p-1 text-popover-foreground"
            >
              <Show
                when={visible().length > 0}
                fallback={
                  <li class="px-3 py-2 text-center text-xs text-muted-foreground">
                    {props.emptyMessage ?? "Nenhum resultado"}
                  </li>
                }
              >
                <For each={visible()}>
                  {(option, idx) => (
                    <li
                      role="option"
                      aria-selected={option.value === props.value ? "true" : "false"}
                      aria-disabled={option.disabled ? "true" : undefined}
                      onMouseEnter={() => setHighlightedIndex(idx())}
                      onClick={() => {
                        if (option.disabled) return;
                        props.onChange(option.value);
                        setOpen(false);
                        setQuery("");
                      }}
                      class={cn(
                        "relative flex w-full items-center justify-between gap-2 rounded-sm cursor-pointer select-none",
                        isSm() ? "px-2 py-1.5 text-xs" : "px-3 py-2 text-sm",
                        option.disabled && "pointer-events-none opacity-50",
                        option.value === props.value
                          ? "bg-accent text-accent-foreground"
                          : idx() === highlightedIndex()
                            ? "bg-muted text-foreground"
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
                        <CheckIcon class="size-3.5 shrink-0" />
                      </Show>
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
