import { hotkeyTokens, type HotkeyInput } from "./hotkey-definition";
import { For } from "solid-js";
import { cn } from "./lib/cn";
export interface KbdProps {
  keys: HotkeyInput | null;
  class?: string;
}
const labels: Record<string, string> = {
  mod: "Ctrl/⌘",
  ctrl: "Ctrl",
  control: "Ctrl",
  cmd: "⌘",
  meta: "⌘",
  command: "⌘",
  shift: "⇧",
  alt: "Alt",
  option: "Alt",
  escape: "Esc",
  esc: "Esc",
  enter: "↵",
  arrowup: "↑",
  arrowdown: "↓",
  arrowleft: "←",
  arrowright: "→",
  space: "Space",
  plus: "+",
};
/** Visual shortcut hint; register behavior separately with createHotkeys. */
export function Kbd(props: KbdProps) {
  return (
    <span
      aria-label={
        props.keys ? hotkeyTokens(props.keys).join(" + ") : "Atalho inválido"
      }
      class={cn(
        "inline-flex shrink-0 items-center gap-1 text-muted-foreground",
        props.class,
      )}
    >
      <For each={props.keys ? hotkeyTokens(props.keys) : ["Atalho inválido"]}>
        {(key) => (
          <kbd class="inline-flex min-w-4 items-center justify-center rounded border border-border bg-muted/50 px-1 py-0.5 font-mono text-[10px] leading-none">
            {labels[key.toLowerCase()] ?? key.toUpperCase()}
          </kbd>
        )}
      </For>
    </span>
  );
}
