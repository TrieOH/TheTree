import type { HotkeyInput } from "./hotkey-definition";
import {
  For,
  Show,
  createEffect,
  createMemo,
  createSignal,
  createUniqueId,
} from "solid-js";
import type { JSX } from "@solidjs/web";
import Search from "~icons/lucide/search";
import { Kbd } from "./kbd";
import { Dialog } from "./dialog";
import { createHotkeys } from "./hotkeys";
import { cn } from "./lib/cn";
export interface CommandItem {
  id: string;
  label: string;
  description?: string;
  keywords?: readonly string[];
  group?: string;
  icon?: JSX.Element;
  shortcut?: HotkeyInput;
  disabled?: boolean;
  onSelect: () => void;
}
export interface CommandPaletteProps {
  items: readonly CommandItem[];
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  hotkey?: boolean;
  shortcut?: HotkeyInput;
  /** Stable result area in pixels; clamped to the viewport. Defaults to 320. */
  resultsHeight?: number;
  title?: string;
  placeholder?: string;
  emptyMessage?: string;
  class?: string;
}
const normalize = (value: string) =>
  value
    .normalize("NFD")
    .replace(/\p{Diacritic}/gu, "")
    .toLocaleLowerCase();
export function CommandPalette(props: CommandPaletteProps) {
  const id = createUniqueId();
  const [internal, setInternal] = createSignal(false),
    [query, setQuery] = createSignal(""),
    [active, setActive] = createSignal(0);
  const open = () => props.open ?? internal();
  const change = (value: boolean) => {
    setInternal(value);
    props.onOpenChange?.(value);
  };
  const visible = createMemo(() =>
    props.items.filter((item) =>
      normalize(
        [item.label, item.description, ...(item.keywords ?? [])].join(" "),
      ).includes(normalize(query().trim())),
    ),
  );
  const enabled = () => visible().filter((item) => !item.disabled);
  const selected = () =>
    enabled()[Math.min(active(), Math.max(0, enabled().length - 1))];
  let input: HTMLInputElement | undefined;
  createHotkeys(
    () => [
      {
        hotkey: props.shortcut ?? { key: "k", modifiers: ["Mod"] },
        callback: () => change(!open()),
        options: { enabled: props.hotkey !== false },
      },
    ],
    { ignoreInputs: false },
  );
  createEffect(
    () => open(),
    (value) => {
      if (value) {
        setQuery("");
        setActive(0);
        const timer = setTimeout(() => input?.focus(), 0);
        return () => clearTimeout(timer);
      }
    },
  );
  createEffect(
    () => selected()?.id,
    (value) => {
      if (value)
        document
          .getElementById(`${id}-${value}`)
          ?.scrollIntoView?.({ block: "nearest" });
    },
  );
  const run = (item?: CommandItem) => {
    if (!item || item.disabled) return;
    change(false);
    item.onSelect();
  };
  return (
    <Dialog
      open={open()}
      onOpenChange={change}
      title={props.title ?? "Buscar páginas e ações"}
      size="md"
      class={props.class}
      contentClass="p-0 sm:p-0"
      footer={
        <span class="text-xs text-muted-foreground">
          ↑ ↓ navegar · Enter selecionar · Esc fechar
        </span>
      }
    >
      <div class="flex items-center gap-2 border-b border-border px-4">
        <Search class="size-4 text-muted-foreground" />
        <input
          ref={(element) => {
            input = element;
          }}
          role="combobox"
          aria-autocomplete="list"
          aria-expanded="true"
          aria-controls={`${id}-results`}
          aria-activedescendant={
            selected() ? `${id}-${selected()!.id}` : undefined
          }
          aria-label={props.placeholder ?? "Buscar comando"}
          placeholder={props.placeholder ?? "Buscar páginas e ações…"}
          value={query()}
          onInput={(e) => {
            setQuery(e.currentTarget.value);
            setActive(0);
          }}
          class="h-12 min-w-0 flex-1 bg-transparent text-sm outline-none"
          onKeyDown={(e) => {
            if (e.key === "ArrowDown" || e.key === "ArrowUp") {
              e.preventDefault();
              setActive(
                (index) =>
                  (index +
                    (e.key === "ArrowDown" ? 1 : -1) +
                    enabled().length) %
                  Math.max(1, enabled().length),
              );
            } else if (e.key === "Enter") {
              e.preventDefault();
              run(selected());
            }
          }}
        />
      </div>
      <div
        id={`${id}-results`}
        role="listbox"
        aria-label="Resultados"
        class="overflow-y-auto overscroll-contain p-2 [scrollbar-gutter:stable]"
        style={{
          height: `min(${Math.max(128, props.resultsHeight ?? 320)}px, 55dvh)`,
        }}
      >
        <Show
          when={visible().length}
          fallback={
            <p
              role="status"
              class="flex h-full items-center justify-center px-3 py-8 text-center text-sm text-muted-foreground"
            >
              {props.emptyMessage ?? "Nenhum resultado encontrado."}
            </p>
          }
        >
          <For each={visible()}>
            {(item, index) => (
              <>
                <Show
                  when={
                    item.group && item.group !== visible()[index() - 1]?.group
                  }
                >
                  <p class="px-2 pt-3 pb-1 text-xs font-medium text-muted-foreground">
                    {item.group}
                  </p>
                </Show>
                <div
                  id={`${id}-${item.id}`}
                  role="option"
                  aria-selected={selected()?.id === item.id ? "true" : "false"}
                  aria-disabled={item.disabled ? "true" : undefined}
                  onPointerMove={() => {
                    const i = enabled().findIndex(
                      (value) => value.id === item.id,
                    );
                    if (i >= 0) setActive(i);
                  }}
                  onClick={() => run(item)}
                  class={cn(
                    "flex min-w-0 cursor-pointer items-center gap-3 rounded-lg p-2.5",
                    selected()?.id === item.id && "bg-muted text-foreground",
                    item.disabled && "cursor-not-allowed opacity-40",
                  )}
                >
                  <Show when={item.icon}>
                    <span
                      aria-hidden="true"
                      class="shrink-0 text-muted-foreground [&>svg]:size-4"
                    >
                      {item.icon}
                    </span>
                  </Show>
                  <span class="min-w-0 flex-1">
                    <span class="block text-sm wrap-anywhere">
                      {item.label}
                    </span>
                    <Show when={item.description}>
                      <span class="block text-xs text-muted-foreground wrap-anywhere">
                        {item.description}
                      </span>
                    </Show>
                  </span>
                  <Show when={item.shortcut}>
                    <Kbd keys={item.shortcut!} />
                  </Show>
                </div>
              </>
            )}
          </For>
        </Show>
      </div>
    </Dialog>
  );
}
