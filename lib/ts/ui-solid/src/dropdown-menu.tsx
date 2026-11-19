import type { HotkeyInput } from "./hotkey-definition";
import type { JSX } from "@solidjs/web";
import { For, Show, createSignal, onCleanup } from "solid-js";
import { Popover, type PopoverProps } from "./popover";
import { createHotkeys } from "./hotkeys";
import { Kbd } from "./kbd";
import { cn } from "./lib/cn";
export type DropdownMenuItem =
  | {
      id: string;
      label: string;
      icon?: JSX.Element;
      shortcut?: HotkeyInput;
      destructive?: boolean;
      disabled?: boolean;
      onSelect: () => void;
    }
  | { id: string; separator: true; label?: string };
export interface DropdownMenuProps extends Omit<
  PopoverProps,
  "children" | "role"
> {
  items: readonly DropdownMenuItem[];
  /** Shortcuts are active in the open menu by default. Global opts into closed-menu actions. */
  shortcutScope?: "menu" | "global" | "none";
}
export function DropdownMenu(props: DropdownMenuProps) {
  const [internalOpen, setInternalOpen] = createSignal(false);
  const open = () => props.open ?? internalOpen();
  const changeOpen = (value: boolean) => {
    setInternalOpen(value);
    props.onOpenChange?.(value);
  };
  let closeMenu: (() => void) | undefined;
  createHotkeys(() =>
    props.items.flatMap((item) =>
      "separator" in item || !item.shortcut
        ? []
        : [
            {
              hotkey: item.shortcut,
              options: {
                enabled:
                  !props.disabled &&
                  !item.disabled &&
                  props.shortcutScope !== "none" &&
                  (props.shortcutScope === "global" || open()),
              },
              callback: () => {
                if (open()) {
                  if (closeMenu) closeMenu();
                  else changeOpen(false);
                }
                item.onSelect();
              },
            },
          ],
    ),
  );
  let search = "",
    timer: ReturnType<typeof setTimeout> | undefined;
  onCleanup(() => clearTimeout(timer));
  const keyboard = (event: KeyboardEvent) => {
    const menu = event.currentTarget as HTMLElement;
    const items = [
      ...menu.querySelectorAll<HTMLButtonElement>(
        '[role="menuitem"]:not([disabled])',
      ),
    ];
    const index = items.indexOf(document.activeElement as HTMLButtonElement);
    let next: number | undefined;
    if (event.key === "ArrowDown") next = (index + 1) % items.length;
    if (event.key === "ArrowUp")
      next = (index - 1 + items.length) % items.length;
    if (event.key === "Home") next = 0;
    if (event.key === "End") next = items.length - 1;
    if (
      event.key.length === 1 &&
      !event.ctrlKey &&
      !event.metaKey &&
      event.key !== " "
    ) {
      clearTimeout(timer);
      search += event.key.toLocaleLowerCase();
      next = items.findIndex((item) =>
        item.dataset.label?.toLocaleLowerCase().startsWith(search),
      );
      timer = setTimeout(() => {
        search = "";
      }, 500);
    }
    if (next !== undefined && next >= 0 && items[next]) {
      event.preventDefault();
      items[next].focus();
    }
  };
  return (
    <Popover
      {...props}
      open={open()}
      onOpenChange={changeOpen}
      role="menu"
      class={cn("p-1", props.class)}
    >
      {(close) => {
        closeMenu = close;
        return (
          <div
            onKeyDown={(event) => {
              if (event.key === "Tab") close();
              else keyboard(event);
            }}
          >
            <For each={props.items}>
              {(item) =>
                "separator" in item ? (
                  <>
                    <div role="separator" class="my-1 border-t border-border" />
                    <Show when={item.label}>
                      <p class="px-2 py-1 text-xs font-medium text-muted-foreground">
                        {item.label}
                      </p>
                    </Show>
                  </>
                ) : (
                  <button
                    type="button"
                    role="menuitem"
                    tabindex={-1}
                    data-label={item.label}
                    disabled={item.disabled}
                    class={cn(
                      "flex w-full min-w-0 items-center gap-2 rounded-md px-2 py-2 text-left text-sm outline-none hover:bg-muted focus:bg-muted disabled:pointer-events-none disabled:opacity-40",
                      item.destructive && "text-destructive",
                    )}
                    onClick={() => {
                      if (!item.disabled) {
                        close();
                        item.onSelect();
                      }
                    }}
                  >
                    <Show when={item.icon}>
                      <span aria-hidden="true" class="shrink-0 [&>svg]:size-4">
                        {item.icon}
                      </span>
                    </Show>
                    <span class="min-w-0 flex-1 wrap-anywhere">
                      {item.label}
                    </span>
                    <Show when={item.shortcut}>
                      <Kbd keys={item.shortcut!} />
                    </Show>
                  </button>
                )
              }
            </For>
          </div>
        );
      }}
    </Popover>
  );
}
