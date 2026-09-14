import { onCleanup } from "solid-js";

export interface HotkeyItem {
  hotkey: string;
  callback: (e: KeyboardEvent) => void;
  options?: {
    enabled?: boolean;
  };
}

export interface CreateHotkeysOptions {
  ignoreInputs?: boolean;
  preventDefault?: boolean;
}

function matchesHotkey(e: KeyboardEvent, hotkey: string): boolean {
  const parts = hotkey.split("+").map((p) => p.trim().toLowerCase());
  const keyPart = parts[parts.length - 1];
  const modifiers = new Set(parts.slice(0, -1));

  const isMod = e.metaKey || e.ctrlKey;
  const needMod = modifiers.has("mod");
  const needCtrl = modifiers.has("ctrl") || modifiers.has("control");
  const needMeta =
    modifiers.has("meta") || modifiers.has("cmd") || modifiers.has("command");
  const needShift = modifiers.has("shift");
  const needAlt = modifiers.has("alt") || modifiers.has("option");

  if (needMod) {
    if (!isMod) return false;
  } else {
    if (needCtrl !== e.ctrlKey) return false;
    if (needMeta !== e.metaKey) return false;
  }

  if (needShift !== e.shiftKey) return false;
  if (needAlt !== e.altKey) return false;

  return e.key.toLowerCase() === keyPart.toLowerCase();
}

/**
 * Native Solid 2.0 hotkeys hook.
 * Avoids third-party Solid 1.x reactivity mismatches.
 */
export function createHotkeys(
  items: () => HotkeyItem[],
  globalOptions: CreateHotkeysOptions = {},
): void {
  if (typeof window === "undefined") return;

  const handleKeyDown = (e: KeyboardEvent) => {
    if (globalOptions.ignoreInputs ?? true) {
      const target = e.target as HTMLElement | null;
      if (
        target &&
        (target.tagName === "INPUT" ||
          target.tagName === "TEXTAREA" ||
          target.tagName === "SELECT" ||
          target.isContentEditable)
      ) {
        return;
      }
    }

    const currentItems = items();
    for (const item of currentItems) {
      if (item.options?.enabled === false) continue;
      if (matchesHotkey(e, item.hotkey)) {
        if (globalOptions.preventDefault ?? true) {
          e.preventDefault();
        }
        item.callback(e);
        break;
      }
    }
  };

  window.addEventListener("keydown", handleKeyDown);

  onCleanup(() => {
    window.removeEventListener("keydown", handleKeyDown);
  });
}
