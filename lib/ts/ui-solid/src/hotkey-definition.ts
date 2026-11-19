export const hotkeyKeys = [
  "a",
  "b",
  "c",
  "d",
  "e",
  "f",
  "g",
  "h",
  "i",
  "j",
  "k",
  "l",
  "m",
  "n",
  "o",
  "p",
  "q",
  "r",
  "s",
  "t",
  "u",
  "v",
  "w",
  "x",
  "y",
  "z",
  "0",
  "1",
  "2",
  "3",
  "4",
  "5",
  "6",
  "7",
  "8",
  "9",
  "Enter",
  "Escape",
  "Space",
  "Tab",
  "Backspace",
  "Delete",
  "Insert",
  "Home",
  "End",
  "PageUp",
  "PageDown",
  "ArrowUp",
  "ArrowDown",
  "ArrowLeft",
  "ArrowRight",
  "F1",
  "F2",
  "F3",
  "F4",
  "F5",
  "F6",
  "F7",
  "F8",
  "F9",
  "F10",
  "F11",
  "F12",
  "F13",
  "F14",
  "F15",
  "F16",
  "F17",
  "F18",
  "F19",
  "F20",
  "F21",
  "F22",
  "F23",
  "F24",
  "+",
  "-",
  "=",
  "/",
  ".",
  ",",
  ";",
  "'",
  "[",
  "]",
  "\\",
  "`",
] as const;
export type HotkeyKey = (typeof hotkeyKeys)[number];
export type HotkeyModifier = "Mod" | "Ctrl" | "Meta" | "Alt" | "Shift";
export interface Hotkey {
  key: HotkeyKey;
  modifiers?: readonly HotkeyModifier[];
}
type StringKey = Exclude<HotkeyKey | Uppercase<HotkeyKey>, "+"> | "Plus";
/** Checked legacy notation. Prefer Hotkey objects for autocomplete of each field. */
export type HotkeyString =
  | StringKey
  | `${HotkeyModifier}+${StringKey}`
  | `${HotkeyModifier}+${HotkeyModifier}+${StringKey}`
  | `${HotkeyModifier}+${HotkeyModifier}+${HotkeyModifier}+${StringKey}`;
export type HotkeyInput = Hotkey | HotkeyString;
const modifiers: Record<string, HotkeyModifier> = {
  mod: "Mod",
  ctrl: "Ctrl",
  control: "Ctrl",
  cmd: "Meta",
  command: "Meta",
  meta: "Meta",
  option: "Alt",
  alt: "Alt",
  shift: "Shift",
};
const aliases: Record<string, string> = {
  esc: "Escape",
  spacebar: "Space",
  return: "Enter",
  plus: "+",
  del: "Delete",
};
/** Validate strings at external boundaries (settings, URL or text input); invalid input returns null. */
export function parseHotkey(input: string): Hotkey | null {
  const parts = input.split("+").map((part) => part.trim());
  const last = parts.pop();
  if (!last) return null;
  const name = aliases[last.toLowerCase()] ?? last;
  const key = hotkeyKeys.find(
    (key) => key.toLowerCase() === name.toLowerCase(),
  );
  const selected = parts.map((part) => modifiers[part.toLowerCase()]);
  if (
    !key ||
    selected.some((modifier) => !modifier) ||
    new Set(selected).size !== selected.length
  )
    return null;
  return { key, modifiers: selected };
}
export const resolveHotkey = (input: HotkeyInput): Hotkey | null =>
  typeof input === "string" ? parseHotkey(input) : input;
export const hotkeyTokens = (input: HotkeyInput): readonly string[] => {
  const value = resolveHotkey(input);
  return value ? [...(value.modifiers ?? []), value.key] : [];
};
