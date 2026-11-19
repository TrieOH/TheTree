import { resolveHotkey, type HotkeyInput } from "./hotkey-definition";
import { onCleanup } from "solid-js";

export interface HotkeyItem {
  /** Prefer { key: "k", modifiers: ["Mod"] }. Mod accepts Ctrl or Command. */
  hotkey: HotkeyInput;
  callback: (event: KeyboardEvent) => void;
  options?: { enabled?: boolean };
}
export interface CreateHotkeysOptions {
  ignoreInputs?: boolean;
  preventDefault?: boolean;
  allowRepeat?: boolean;
  /** Limit shortcuts to events originating inside this element. */
  scope?: () => HTMLElement | undefined;
}
type Registration = {
  items: () => readonly HotkeyItem[];
  options: CreateHotkeysOptions;
};
const registries = new WeakMap<
  Window,
  { registrations: Set<Registration>; listener: (event: KeyboardEvent) => void }
>();
function matches(event: KeyboardEvent, input: HotkeyInput) {
  const shortcut = resolveHotkey(input);
  if (!shortcut) return false;
  const key = shortcut.key === "Space" ? " " : shortcut.key;
  if (event.key.toLowerCase() !== key.toLowerCase()) return false;
  const modifiers = new Set(shortcut.modifiers ?? []);
  const ctrl = modifiers.has("Ctrl"),
    meta = modifiers.has("Meta");
  if (modifiers.has("Mod")) {
    if (!(event.ctrlKey || event.metaKey)) return false;
    if ((ctrl && !event.ctrlKey) || (meta && !event.metaKey)) return false;
  } else if (ctrl !== event.ctrlKey || meta !== event.metaKey) return false;
  return (
    modifiers.has("Shift") === event.shiftKey &&
    modifiers.has("Alt") === event.altKey
  );
}
function editable(event: KeyboardEvent) {
  return event
    .composedPath()
    .some(
      (target) =>
        target instanceof HTMLElement &&
        (target.matches("input,textarea,select") ||
          target.isContentEditable ||
          target.getAttribute("role") === "textbox" ||
          !!target.closest('[contenteditable]:not([contenteditable="false"])')),
    );
}
/** Lifecycle-bound shortcuts. Latest eligible registration handles a collision once. */
export function createHotkeys(
  items: () => readonly HotkeyItem[],
  options: CreateHotkeysOptions = {},
): void {
  if (typeof window === "undefined") return;
  let registry = registries.get(window);
  if (!registry) {
    const registrations = new Set<Registration>();
    const listener = (event: KeyboardEvent) => {
      if (event.defaultPrevented || event.isComposing || event.keyCode === 229)
        return;
      for (const registration of [...registrations].reverse()) {
        const settings = registration.options;
        if (event.repeat && !settings.allowRepeat) continue;
        if ((settings.ignoreInputs ?? true) && editable(event)) continue;
        const scope = settings.scope?.();
        if (settings.scope && (!scope || !event.composedPath().includes(scope)))
          continue;
        const item = registration
          .items()
          .find(
            (item) =>
              item.options?.enabled !== false && matches(event, item.hotkey),
          );
        if (!item) continue;
        if (settings.preventDefault ?? true) event.preventDefault();
        item.callback(event);
        return;
      }
    };
    registry = { registrations, listener };
    registries.set(window, registry);
    window.addEventListener("keydown", listener);
  }
  const registration = { items, options };
  registry.registrations.add(registration);
  onCleanup(() => {
    registry.registrations.delete(registration);
    if (!registry.registrations.size) {
      window.removeEventListener("keydown", registry.listener);
      registries.delete(window);
    }
  });
}
