export type Theme = "light" | "dark";
const storageKey = "trieoh-theme";

export function readThemePreference(): Theme | null {
  try {
    const value = localStorage.getItem(storageKey);
    return value === "dark" || value === "light" ? value : null;
  } catch {
    return null;
  }
}

export function saveThemePreference(theme: Theme) {
  try {
    localStorage.setItem(storageKey, theme);
  } catch {
    // Switching themes also works when browser storage is unavailable.
  }
}

export function applyTheme(theme: Theme) {
  const dark = theme === "dark";
  document.documentElement.classList.toggle("dark", dark);
  document
    .querySelector('meta[name="theme-color"]')
    ?.setAttribute("content", dark ? "#191c24" : "#f5f5f5");
}
